package main

import (
    "bytes"
    "context"
    "crypto/sha256"
    "crypto/sha512"
    "encoding/base32"
    "encoding/base64"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "log"
    "math/big"
    "net"
    "sync"
    "time"
    
    "github.com/StripChain/strip-node/bitcoin"
    "github.com/StripChain/strip-node/dogecoin"
    "github.com/StripChain/strip-node/ripple"
    ecdsaKeygen "github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
    eddsaKeygen "github.com/bnb-chain/tss-lib/v2/eddsa/keygen"
    "github.com/bnb-chain/tss-lib/v2/tss"
    "github.com/decred/dcrd/dcrec/edwards/v2"
    "github.com/ethereum/go-ethereum/crypto"
    "github.com/mr-tron/base58"
    "github.com/stellar/go/strkey"
    "golang.org/x/crypto/blake2b"
)

// metaData stores additional data beyond keys
var metaData = struct {
    data map[string][]byte
    mu   sync.RWMutex
}{
    data: make(map[string][]byte),
}

// KeygenSession represents an active key generation operation
type KeygenSession struct {
    identity      string
    identityCurve string
    keyCurve      string
    signers       []string
    party         tss.Party
    outChan       chan tss.Message
    saveChans     map[string]interface{} // Different save channels for different curves
    startTime     time.Time
    index         int                    // Our index in signers list
}

// KeygenSessions tracks ongoing key generation operations
var (
    keygenSessions     = map[string]*KeygenSession{}
    keygenSessionsMutex sync.RWMutex
)

// SessionTimeout is the maximum time a keygen session can be active
const SessionTimeout = 10 * time.Minute

// storeKeyReference is implemented in both nitro_keyref.go and here for backward compatibility
// This version includes the address parameter which is specific to key generation
func storeKeyReference(keyID string, sealedKey []byte, address string) {
    // Also store in the global keyRefStore
    keyRefStore.mutex.Lock()
    keyRefStore.references[keyID] = sealedKey
    keyRefStore.mutex.Unlock()
    
    // Also store in the keyStore for backward compatibility
    keyStore.AddKey(keyID, sealedKey, address)
    log.Printf("Stored key reference for %s with address %s", keyID, address)
}

// storeSignersData stores signers information for a key
func storeSignersData(keyID string, signersData []byte) {
    metaData.mu.Lock()
    defer metaData.mu.Unlock()
    
    signersKey := keyID + "_signers"
    metaData.data[signersKey] = signersData
    
    log.Printf("Stored signers data for key %s", keyID)
}

// getKeyReference is now implemented in nitro_keyref.go
// This is a legacy wrapper that will be removed in a future version

// getSignersData retrieves signers data for a key ID
func getSignersData(keyID string) ([]byte, bool) {
    metaData.mu.RLock()
    defer metaData.mu.RUnlock()
    
    signersKey := keyID + "_signers"
    signersData, exists := metaData.data[signersKey]
    return signersData, exists
}

// handleKeyGen handles incoming key generation requests from the parent instance
func handleKeyGen(conn net.Conn, req Request) {
    log.Printf("Handling key generation for %s_%s_%s", req.Identity, req.IdentityCurve, req.KeyCurve)
    
    // Validate request
    if req.Identity == "" || req.IdentityCurve == "" || req.KeyCurve == "" || len(req.Signers) == 0 {
        sendError(conn, "missing required parameters")
        return
    }
    
    // Generate session ID
    sessionID := fmt.Sprintf("%s_%s_%s", req.Identity, req.IdentityCurve, req.KeyCurve)
    
    // Check if key already exists
    sealedKey, exists := getKeyReference(sessionID)
    if exists {
        log.Printf("Key already exists for %s", sessionID)
        sendError(conn, "key already exists")
        return
    }
    
    // Determine our signer index
    enclaveID := getEnclaveID()
    signerIndex := determineSignerIndex(enclaveID, req.Signers)
    if signerIndex == -1 {
        log.Printf("Enclave is not in signers list for %s", sessionID)
        sendError(conn, "enclave is not in signers list")
        return
    }
    
    log.Printf("Our index in signers: %d", signerIndex)
    
    // Check total signers
    if len(req.Signers) == 0 {
        sendError(conn, "no signers provided")
        return
    }
    
    // Check for existing session
    keygenSessionsMutex.RLock()
    existingSession := keygenSessions[sessionID]
    keygenSessionsMutex.RUnlock()
    
    if existingSession != nil {
        // We're already processing this key generation
        log.Printf("Keygen session %s already active", sessionID)
        sendSuccess(conn, Response{
            Status: "in_progress",
        })
        return
    }
    
    // Start a new key generation session
    session, err := startKeygenSession(sessionID, req.Identity, req.IdentityCurve, req.KeyCurve, 
                                      req.Signers, signerIndex)
    if err != nil {
        log.Printf("Failed to start keygen session: %v", err)
        sendError(conn, fmt.Sprintf("failed to start keygen session: %v", err))
        return
    }
    
    // Register the session
    keygenSessionsMutex.Lock()
    keygenSessions[sessionID] = session
    keygenSessionsMutex.Unlock()
    
    // Start session cleanup timer
    go func() {
        time.Sleep(SessionTimeout)
        keygenSessionsMutex.Lock()
        delete(keygenSessions, sessionID)
        keygenSessionsMutex.Unlock()
        log.Printf("Keygen session %s timed out", sessionID)
    }()
    
    // Tell the parent that we've started the key generation process
    sendSuccess(conn, Response{
        Status: "started",
    })
    
    // Handle session messages in the background
    go handleKeygenSession(sessionID, session)
}

// startKeygenSession initializes a new key generation session
func startKeygenSession(sessionID, identity, identityCurve, keyCurve string, 
                       signers []string, signerIndex int) (*KeygenSession, error) {
    
    // Create TSS parties
    totalSigners := len(signers)
    parties, partyIDs := getParties(totalSigners)
    
    // Calculate threshold
    threshold := calculateThreshold(totalSigners)
    
    // Create TSS peer context
    ctx := tss.NewPeerContext(parties)
    
    // Channels for communication
    outChan := make(chan tss.Message)
    saveChans := make(map[string]interface{})
    
    // Create a new keygen session
    session := &KeygenSession{
        identity:      identity,
        identityCurve: identityCurve,
        keyCurve:      keyCurve,
        signers:       signers,
        outChan:       outChan,
        saveChans:     saveChans,
        startTime:     time.Now(),
        index:         signerIndex,
    }
    
    // Initialize the party based on the key curve
    params := tss.NewParameters(getCurveType(keyCurve), ctx, partyIDs[signerIndex], len(parties), threshold)
    
    var err error
    
    switch keyCurve {
    case "ecdsa":
        // Standard ECDSA
        saveChan := make(chan *ecdsaKeygen.LocalPartySaveData)
        saveChans["ecdsa"] = &saveChan
        
        preParams, err := ecdsaKeygen.GeneratePreParams(2 * time.Minute)
        if err != nil {
            return nil, fmt.Errorf("failed to generate ECDSA pre-params: %w", err)
        }
        
        localParty := ecdsaKeygen.NewLocalParty(params, outChan, saveChan, *preParams)
        session.party = localParty
        
        // Start the party
        go localParty.Start()
        
    case "eddsa", "aptos_eddsa", "sui_eddsa", "stellar_eddsa", "algorand_eddsa", "ripple_eddsa", "cardano_eddsa":
        // All EdDSA-based curves
        saveChan := make(chan *eddsaKeygen.LocalPartySaveData)
        saveChans[keyCurve] = &saveChan
        
        localParty := eddsaKeygen.NewLocalParty(params, outChan, saveChan)
        session.party = localParty
        
        // Start the party
        go localParty.Start()
        
    case "bitcoin_ecdsa", "secp256k1":
        // Bitcoin and Secp256k1 (like Dogecoin)
        saveChan := make(chan *ecdsaKeygen.LocalPartySaveData)
        saveChans[keyCurve] = &saveChan
        
        preParams, err := ecdsaKeygen.GeneratePreParams(2 * time.Minute)
        if err != nil {
            return nil, fmt.Errorf("failed to generate Bitcoin/Secp256k1 pre-params: %w", err)
        }
        
        localParty := ecdsaKeygen.NewLocalParty(params, outChan, saveChan, *preParams)
        session.party = localParty
        
        // Start the party
        go localParty.Start()
        
    default:
        return nil, fmt.Errorf("unsupported key curve: %s", keyCurve)
    }
    
    return session, nil
}

// handleKeygenSession processes the key generation session
// This function is already defined in kms.go

func handleKeygenSession(sessionID string, session *KeygenSession) {
    completed := false
    
    for !completed {
        select {
        case msg := <-session.outChan:
            // TSS message that needs to be sent to other parties
            dest := msg.GetTo()
            wireBytes, _, _ := msg.WireBytes()
            
            to := 0
            if dest == nil {
                to = -1 // Broadcast
            } else {
                to = dest[0].Index
            }
            
            // Forward this message to the parent instance for distribution
            keygenMessage := TssMessage{
                SessionID:    sessionID,
                From:         msg.GetFrom().Index,
                To:           to,
                IsBroadcast:  msg.IsBroadcast(),
                WireMessage:  wireBytes,
                KeyCurve:     session.keyCurve,
                Signers:      session.signers,
            }
            
            messageBytes, _ := json.Marshal(keygenMessage)
            sendParentMessage("keygen_message", messageBytes)
            log.Printf("Sent keygen message from party %d to %d (broadcast: %v)", 
                      keygenMessage.From, keygenMessage.To, keygenMessage.IsBroadcast)
            
        // Handle save channels for different curve types
        case save := <-*(session.saveChans["ecdsa"].(*chan *ecdsaKeygen.LocalPartySaveData)):
            handleECDSAKeygenCompletion(sessionID, session, save)
            completed = true
            
        case save := <-*(session.saveChans["eddsa"].(*chan *eddsaKeygen.LocalPartySaveData)):
            handleEdDSAKeygenCompletion(sessionID, session, save)
            completed = true
            
        case save := <-*(session.saveChans["bitcoin_ecdsa"].(*chan *ecdsaKeygen.LocalPartySaveData)):
            handleBitcoinKeygenCompletion(sessionID, session, save)
            completed = true
            
        case save := <-*(session.saveChans["secp256k1"].(*chan *ecdsaKeygen.LocalPartySaveData)):
            handleSecp256k1KeygenCompletion(sessionID, session, save)
            completed = true
            
        case save := <-*(session.saveChans["sui_eddsa"].(*chan *eddsaKeygen.LocalPartySaveData)):
            handleSuiKeygenCompletion(sessionID, session, save)
            completed = true
            
        case save := <-*(session.saveChans["stellar_eddsa"].(*chan *eddsaKeygen.LocalPartySaveData)):
            handleStellarKeygenCompletion(sessionID, session, save)
            completed = true
            
        case save := <-*(session.saveChans["algorand_eddsa"].(*chan *eddsaKeygen.LocalPartySaveData)):
            handleAlgorandKeygenCompletion(sessionID, session, save)
            completed = true
            
        case save := <-*(session.saveChans["ripple_eddsa"].(*chan *eddsaKeygen.LocalPartySaveData)):
            handleRippleKeygenCompletion(sessionID, session, save)
            completed = true
            
        case save := <-*(session.saveChans["aptos_eddsa"].(*chan *eddsaKeygen.LocalPartySaveData)):
            handleAptosKeygenCompletion(sessionID, session, save)
            completed = true
            
        case save := <-*(session.saveChans["cardano_eddsa"].(*chan *eddsaKeygen.LocalPartySaveData)):
            handleCardanoKeygenCompletion(sessionID, session, save)
            completed = true
            
        case <-time.After(SessionTimeout):
            // Timeout case
            completed = true
            log.Printf("Keygen session %s timed out while processing", sessionID)
            
            // Clean up the session
            keygenSessionsMutex.Lock()
            delete(keygenSessions, sessionID)
            keygenSessionsMutex.Unlock()
        }
    }
}

// handleTssMessage processes TSS messages forwarded from the parent instance
func handleTssKeygenMessage(message []byte) {
    var keygenMessage TssMessage
    if err := json.Unmarshal(message, &keygenMessage); err != nil {
        log.Printf("Failed to unmarshal keygen message: %v", err)
        return
    }
    
    // Find the session
    keygenSessionsMutex.RLock()
    session := keygenSessions[keygenMessage.SessionID]
    keygenSessionsMutex.RUnlock()
    
    if session == nil {
        log.Printf("Received message for unknown session %s", keygenMessage.SessionID)
        return
    }
    
    // Process the message
    parties, _ := getParties(len(session.signers))
    
    // Parse the wire message
    pMsg, err := tss.ParseWireMessage(keygenMessage.WireMessage, 
                                     parties[keygenMessage.From], keygenMessage.IsBroadcast)
    if err != nil {
        log.Printf("Failed to parse wire message: %v", err)
        return
    }
    
    // Get the party
    party := session.party
    
    // Check if this message is for us
    ourIndex := party.PartyID().Index
    if keygenMessage.To != -1 && keygenMessage.To != ourIndex {
        // Message not intended for us
        return
    }
    
    // Prevent message loops
    if keygenMessage.From == ourIndex {
        // Message from ourselves
        return
    }
    
    // Update the party with the message
    log.Printf("Processing keygen message from party %d to %d (broadcast: %v)", 
              keygenMessage.From, keygenMessage.To, keygenMessage.IsBroadcast)
    
    _, err = party.Update(pMsg)
    if err != nil {
        log.Printf("Failed to update party with message: %v", err)
        return
    }
    
    log.Printf("Successfully processed keygen message for session %s", keygenMessage.SessionID)
}

// Handle completed key generation for different curve types
func handleECDSAKeygenCompletion(sessionID string, session *KeygenSession, save *ecdsaKeygen.LocalPartySaveData) {
    log.Printf("ECDSA key generation completed for %s", sessionID)
    
    // Format Ethereum address
    x := toHexInt(save.ECDSAPub.X())
    y := toHexInt(save.ECDSAPub.Y())
    publicKeyStr := "04" + x + y
    publicKeyBytes, _ := hex.DecodeString(publicKeyStr)
    address := publicKeyToAddress(publicKeyBytes)
    
    // Serialize the key data
    keyBytes, err := json.Marshal(save)
    if err != nil {
        log.Printf("Failed to marshal key data: %v", err)
        notifyKeygenComplete(sessionID, nil, "", nil, err)
        return
    }
    
    // Seal the key data
    sealedKey, err := sealKeyShare(context.Background(), keyBytes, sessionID)
    if err != nil {
        log.Printf("Failed to seal key data: %v", err)
        notifyKeygenComplete(sessionID, nil, "", nil, err)
        return
    }
    
    // Serialize signers information
    signersBytes, err := json.Marshal(session.signers)
    if err != nil {
        log.Printf("Failed to marshal signers data: %v", err)
        notifyKeygenComplete(sessionID, nil, "", nil, err)
        return
    }
    
    // Store key reference and signers information
    storeKeyReference(sessionID, sealedKey, address)
    storeSignersData(sessionID, signersBytes)
    
    // Clean up session
    keygenSessionsMutex.Lock()
    delete(keygenSessions, sessionID)
    keygenSessionsMutex.Unlock()
    
    // Notify parent
    notifyKeygenComplete(sessionID, sealedKey, address, session.signers, nil)
    
    log.Printf("ECDSA key generation successfully completed for %s with address %s", 
              sessionID, address)
}

// Handle EdDSA key generation completion
func handleEdDSAKeygenCompletion(sessionID string, session *KeygenSession, save *eddsaKeygen.LocalPartySaveData) {
    log.Printf("EdDSA key generation completed for %s", sessionID)
    
    // Get EdDSA public key
    pk := edwards.PublicKey{
        Curve: save.EDDSAPub.Curve(),
        X:     save.EDDSAPub.X(),
        Y:     save.EDDSAPub.Y(),
    }
    
    // Format Solana-style address
    address := base58.Encode(pk.Serialize())
    
    // Serialize the key data
    keyBytes, err := json.Marshal(save)
    if err != nil {
        log.Printf("Failed to marshal key data: %v", err)
        notifyKeygenComplete(sessionID, nil, "", nil, err)
        return
    }
    
    // Seal the key data
    sealedKey, err := sealKeyShare(context.Background(), keyBytes, sessionID)
    if err != nil {
        log.Printf("Failed to seal key data: %v", err)
        notifyKeygenComplete(sessionID, nil, "", nil, err)
        return
    }
    
    // Serialize signers information
    signersBytes, err := json.Marshal(session.signers)
    if err != nil {
        log.Printf("Failed to marshal signers data: %v", err)
        notifyKeygenComplete(sessionID, nil, "", nil, err)
        return
    }
    
    // Store key reference and signers information
    storeKeyReference(sessionID, sealedKey, address)
    storeSignersData(sessionID, signersBytes)
    
    // Clean up session
    keygenSessionsMutex.Lock()
    delete(keygenSessions, sessionID)
    keygenSessionsMutex.Unlock()
    
    // Notify parent
    notifyKeygenComplete(sessionID, sealedKey, address, session.signers, nil)
    
    log.Printf("EdDSA key generation successfully completed for %s with address %s", 
              sessionID, address)
}

// Handler for Bitcoin key generation completion
func handleBitcoinKeygenCompletion(sessionID string, session *KeygenSession, save *ecdsaKeygen.LocalPartySaveData) {
    log.Printf("Bitcoin key generation completed for %s", sessionID)
    
    // Format Bitcoin address
    x := toHexInt(save.ECDSAPub.X())
    y := toHexInt(save.ECDSAPub.Y())
    publicKeyStr := "04" + x + y
    publicKeyBytes, _ := hex.DecodeString(publicKeyStr)
    
    // Use Bitcoin address derivation logic
    bitcoinAddress, _, _ := bitcoin.PublicKeyToBitcoinAddresses(publicKeyBytes)
    
    // Serialize the key data
    keyBytes, err := json.Marshal(save)
    if err != nil {
        log.Printf("Failed to marshal key data: %v", err)
        notifyKeygenComplete(sessionID, nil, "", nil, err)
        return
    }
    
    // Seal the key data
    sealedKey, err := sealKeyShare(context.Background(), keyBytes, sessionID)
    if err != nil {
        log.Printf("Failed to seal key data: %v", err)
        notifyKeygenComplete(sessionID, nil, "", nil, err)
        return
    }
    
    // Serialize signers information
    signersBytes, err := json.Marshal(session.signers)
    if err != nil {
        log.Printf("Failed to marshal signers data: %v", err)
        notifyKeygenComplete(sessionID, nil, "", nil, err)
        return
    }
    
    // Store key reference and signers information
    storeKeyReference(sessionID, sealedKey, bitcoinAddress)
    storeSignersData(sessionID, signersBytes)
    
    // Clean up session
    keygenSessionsMutex.Lock()
    delete(keygenSessions, sessionID)
    keygenSessionsMutex.Unlock()
    
    // Notify parent
    notifyKeygenComplete(sessionID, sealedKey, bitcoinAddress, session.signers, nil)
    
    log.Printf("Bitcoin key generation successfully completed for %s with address %s", 
              sessionID, bitcoinAddress)
}

// Handler for Secp256k1 (e.g., Dogecoin) key generation completion
func handleSecp256k1KeygenCompletion(sessionID string, session *KeygenSession, save *ecdsaKeygen.LocalPartySaveData) {
    log.Printf("Secp256k1 key generation completed for %s", sessionID)
    
    // Format Dogecoin address
    x := toHexInt(save.ECDSAPub.X())
    y := toHexInt(save.ECDSAPub.Y())
    publicKeyStr := "04" + x + y
    
    // Use Dogecoin address derivation logic
    dogecoinAddress, err := dogecoin.PublicKeyToAddress(publicKeyStr)
    if err != nil {
        log.Printf("Error generating Dogecoin address: %v", err)
        notifyKeygenComplete(sessionID, nil, "", nil, err)
        return
    }
    
    // Serialize the key data
    keyBytes, err := json.Marshal(save)
    if err != nil {
        log.Printf("Failed to marshal key data: %v", err)
        notifyKeygenComplete(sessionID, nil, "", nil, err)
        return
    }
    
    // Seal the key data
    sealedKey, err := sealKeyShare(context.Background(), keyBytes, sessionID)
    if err != nil {
        log.Printf("Failed to seal key data: %v", err)
        notifyKeygenComplete(sessionID, nil, "", nil, err)
        return
    }
    
    // Serialize signers information
    signersBytes, err := json.Marshal(session.signers)
    if err != nil {
        log.Printf("Failed to marshal signers data: %v", err)
        notifyKeygenComplete(sessionID, nil, "", nil, err)
        return
    }
    
    // Store key reference and signers information
    storeKeyReference(sessionID, sealedKey, dogecoinAddress)
    storeSignersData(sessionID, signersBytes)
    
    // Clean up session
    keygenSessionsMutex.Lock()
    delete(keygenSessions, sessionID)
    keygenSessionsMutex.Unlock()
    
    // Notify parent
    notifyKeygenComplete(sessionID, sealedKey, dogecoinAddress, session.signers, nil)
    
    log.Printf("Secp256k1 key generation successfully completed for %s with address %s", 
              sessionID, dogecoinAddress)
}

// Handler for Sui key generation completion
func handleSuiKeygenCompletion(sessionID string, session *KeygenSession, save *eddsaKeygen.LocalPartySaveData) {
    log.Printf("Sui key generation completed for %s", sessionID)
    
    // Get the EdDSA public key
    pk := edwards.PublicKey{
        Curve: save.EDDSAPub.Curve(),
        X:     save.EDDSAPub.X(),
        Y:     save.EDDSAPub.Y(),
    }
    
    // Serialize the Ed25519 public key
    pkBytes := pk.Serialize()
    
    // Hash the public key with Blake2b-256 to get Sui address
    hasher := blake2b.Sum256(pkBytes)
    suiAddress := "0x" + hex.EncodeToString(hasher[:])
    
    // Serialize the key data
    keyBytes, err := json.Marshal(save)
    if err != nil {
        log.Printf("Failed to marshal key data: %v", err)
        notifyKeygenComplete(sessionID, nil, "", nil, err)
        return
    }
    
    // Seal the key data
    sealedKey, err := sealKeyShare(context.Background(), keyBytes, sessionID)
    if err != nil {
        log.Printf("Failed to seal key data: %v", err)
        notifyKeygenComplete(sessionID, nil, "", nil, err)
        return
    }
    
    // Serialize signers information
    signersBytes, err := json.Marshal(session.signers)
    if err != nil {
        log.Printf("Failed to marshal signers data: %v", err)
        notifyKeygenComplete(sessionID, nil, "", nil, err)
        return
    }
    
    // Store key reference and signers information
    storeKeyReference(sessionID, sealedKey, suiAddress)
    storeSignersData(sessionID, signersBytes)
    
    // Clean up session
    keygenSessionsMutex.Lock()
    delete(keygenSessions, sessionID)
    keygenSessionsMutex.Unlock()
    
    // Notify parent
    notifyKeygenComplete(sessionID, sealedKey, suiAddress, session.signers, nil)
    
    log.Printf("Sui key generation successfully completed for %s with address %s", 
              sessionID, suiAddress)
}

// Handler for Stellar key generation completion
func handleStellarKeygenCompletion(sessionID string, session *KeygenSession, save *eddsaKeygen.LocalPartySaveData) {
    log.Printf("Stellar key generation completed for %s", sessionID)
    
    // Get the EdDSA public key
    pk := edwards.PublicKey{
        Curve: save.EDDSAPub.Curve(),
        X:     save.EDDSAPub.X(),
        Y:     save.EDDSAPub.Y(),
    }
    
    // Get the public key bytes
    pkBytes := pk.Serialize()
    
    if len(pkBytes) != 32 {
        log.Printf("Invalid public key length")
        notifyKeygenComplete(sessionID, nil, "", nil, fmt.Errorf("invalid public key length"))
        return
    }
    
    // Version byte for ED25519 public key in Stellar
    versionByte := strkey.VersionByteAccountID // 6 << 3, or 48
    
    // Use Stellar SDK's strkey package to encode
    stellarAddress, err := strkey.Encode(versionByte, pkBytes)
    if err != nil {
        log.Printf("Error encoding Stellar address: %v", err)
        notifyKeygenComplete(sessionID, nil, "", nil, err)
        return
    }
    
    // Serialize the key data
    keyBytes, err := json.Marshal(save)
    if err != nil {
        log.Printf("Failed to marshal key data: %v", err)
        notifyKeygenComplete(sessionID, nil, "", nil, err)
        return
    }
    
    // Seal the key data
    sealedKey, err := sealKeyShare(context.Background(), keyBytes, sessionID)
    if err != nil {
        log.Printf("Failed to seal key data: %v", err)
        notifyKeygenComplete(sessionID, nil, "", nil, err)
        return
    }
    
    // Serialize signers information
    signersBytes, err := json.Marshal(session.signers)
    if err != nil {
        log.Printf("Failed to marshal signers data: %v", err)
        notifyKeygenComplete(sessionID, nil, "", nil, err)
        return
    }
    
    // Store key reference and signers information
    storeKeyReference(sessionID, sealedKey, stellarAddress)
    storeSignersData(sessionID, signersBytes)
    
    // Clean up session
    keygenSessionsMutex.Lock()
    delete(keygenSessions, sessionID)
    keygenSessionsMutex.Unlock()
    
    // Notify parent
    notifyKeygenComplete(sessionID, sealedKey, stellarAddress, session.signers, nil)
    
    log.Printf("Stellar key generation successfully completed for %s with address %s", 
              sessionID, stellarAddress)
}

// Handler for Algorand key generation completion
func handleAlgorandKeygenCompletion(sessionID string, session *KeygenSession, save *eddsaKeygen.LocalPartySaveData) {
    log.Printf("Algorand key generation completed for %s", sessionID)
    
    // Get the EdDSA public key
    pk := edwards.PublicKey{
        Curve: save.EDDSAPub.Curve(),
        X:     save.EDDSAPub.X(),
        Y:     save.EDDSAPub.Y(),
    }
    
    // Get the public key bytes
    pkBytes := pk.Serialize()
    
    // Calculate checksum (last 4 bytes of SHA512/256 hash)
    hasher := sha512.New512_256()
    hasher.Write(pkBytes)
    checksum := hasher.Sum(nil)[28:] // Last 4 bytes
    
    // Concatenate public key and checksum
    addressBytes := append(pkBytes, checksum...)
    
    // Encode in base32 without padding
    algorandAddress := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(addressBytes)
    
    // Serialize the key data
    keyBytes, err := json.Marshal(save)
    if err != nil {
        log.Printf("Failed to marshal key data: %v", err)
        notifyKeygenComplete(sessionID, nil, "", nil, err)
        return
    }
    
    // Seal the key data
    sealedKey, err := sealKeyShare(context.Background(), keyBytes, sessionID)
    if err != nil {
        log.Printf("Failed to seal key data: %v", err)
        notifyKeygenComplete(sessionID, nil, "", nil, err)
        return
    }
    
    // Serialize signers information
    signersBytes, err := json.Marshal(session.signers)
    if err != nil {
        log.Printf("Failed to marshal signers data: %v", err)
        notifyKeygenComplete(sessionID, nil, "", nil, err)
        return
    }
    
    // Store key reference and signers information
    storeKeyReference(sessionID, sealedKey, algorandAddress)
    storeSignersData(sessionID, signersBytes)
    
    // Clean up session
    keygenSessionsMutex.Lock()
    delete(keygenSessions, sessionID)
    keygenSessionsMutex.Unlock()
    
    // Notify parent
    notifyKeygenComplete(sessionID, sealedKey, algorandAddress, session.signers, nil)
    
    log.Printf("Algorand key generation successfully completed for %s with address %s", 
              sessionID, algorandAddress)
}

// Handler for Ripple key generation completion
func handleRippleKeygenCompletion(sessionID string, session *KeygenSession, save *eddsaKeygen.LocalPartySaveData) {
    log.Printf("Ripple key generation completed for %s", sessionID)
    
    // Use Ripple-specific address derivation
    rippleAddress := ripple.PublicKeyToAddress(save)
    
    // Serialize the key data
    keyBytes, err := json.Marshal(save)
    if err != nil {
        log.Printf("Failed to marshal key data: %v", err)
        notifyKeygenComplete(sessionID, nil, "", nil, err)
        return
    }
    
    // Seal the key data
    sealedKey, err := sealKeyShare(context.Background(), keyBytes, sessionID)
    if err != nil {
        log.Printf("Failed to seal key data: %v", err)
        notifyKeygenComplete(sessionID, nil, "", nil, err)
        return
    }
    
    // Serialize signers information
    signersBytes, err := json.Marshal(session.signers)
    if err != nil {
        log.Printf("Failed to marshal signers data: %v", err)
        notifyKeygenComplete(sessionID, nil, "", nil, err)
        return
    }
    
    // Store key reference and signers information
    storeKeyReference(sessionID, sealedKey, rippleAddress)
    storeSignersData(sessionID, signersBytes)
    
    // Clean up session
    keygenSessionsMutex.Lock()
    delete(keygenSessions, sessionID)
    keygenSessionsMutex.Unlock()
    
    // Notify parent
    notifyKeygenComplete(sessionID, sealedKey, rippleAddress, session.signers, nil)
    
    log.Printf("Ripple key generation successfully completed for %s with address %s", 
              sessionID, rippleAddress)
}

// Handler for Aptos key generation completion
func handleAptosKeygenCompletion(sessionID string, session *KeygenSession, save *eddsaKeygen.LocalPartySaveData) {
    log.Printf("Aptos key generation completed for %s", sessionID)
    
    // Get the EdDSA public key
    pk := edwards.PublicKey{
        Curve: save.EDDSAPub.Curve(),
        X:     save.EDDSAPub.X(),
        Y:     save.EDDSAPub.Y(),
    }
    
    // For Aptos, use hex encoding of the public key
    aptosAddress := hex.EncodeToString(pk.Serialize())
    
    // Serialize the key data
    keyBytes, err := json.Marshal(save)
    if err != nil {
        log.Printf("Failed to marshal key data: %v", err)
        notifyKeygenComplete(sessionID, nil, "", nil, err)
        return
    }
    
    // Seal the key data
    sealedKey, err := sealKeyShare(context.Background(), keyBytes, sessionID)
    if err != nil {
        log.Printf("Failed to seal key data: %v", err)
        notifyKeygenComplete(sessionID, nil, "", nil, err)
        return
    }
    
    // Serialize signers information
    signersBytes, err := json.Marshal(session.signers)
    if err != nil {
        log.Printf("Failed to marshal signers data: %v", err)
        notifyKeygenComplete(sessionID, nil, "", nil, err)
        return
    }
    
    // Store key reference and signers information
    storeKeyReference(sessionID, sealedKey, aptosAddress)
    storeSignersData(sessionID, signersBytes)
    
    // Clean up session
    keygenSessionsMutex.Lock()
    delete(keygenSessions, sessionID)
    keygenSessionsMutex.Unlock()
    
    // Notify parent
    notifyKeygenComplete(sessionID, sealedKey, aptosAddress, session.signers, nil)
    
    log.Printf("Aptos key generation successfully completed for %s with address %s", 
              sessionID, aptosAddress)
}

// Handler for Cardano key generation completion
func handleCardanoKeygenCompletion(sessionID string, session *KeygenSession, save *eddsaKeygen.LocalPartySaveData) {
    log.Printf("Cardano key generation completed for %s", sessionID)
    
    // Get the EdDSA public key
    pk := edwards.PublicKey{
        Curve: save.EDDSAPub.Curve(),
        X:     save.EDDSAPub.X(),
        Y:     save.EDDSAPub.Y(),
    }
    
    // For Cardano, use hex encoding of the public key
    cardanoAddress := hex.EncodeToString(pk.Serialize())
    
    // Serialize the key data
    keyBytes, err := json.Marshal(save)
    if err != nil {
        log.Printf("Failed to marshal key data: %v", err)
        notifyKeygenComplete(sessionID, nil, "", nil, err)
        return
    }
    
    // Seal the key data
    sealedKey, err := sealKeyShare(context.Background(), keyBytes, sessionID)
    if err != nil {
        log.Printf("Failed to seal key data: %v", err)
        notifyKeygenComplete(sessionID, nil, "", nil, err)
        return
    }
    
    // Serialize signers information
    signersBytes, err := json.Marshal(session.signers)
    if err != nil {
        log.Printf("Failed to marshal signers data: %v", err)
        notifyKeygenComplete(sessionID, nil, "", nil, err)
        return
    }
    
    // Store key reference and signers information
    storeKeyReference(sessionID, sealedKey, cardanoAddress)
    storeSignersData(sessionID, signersBytes)
    
    // Clean up session
    keygenSessionsMutex.Lock()
    delete(keygenSessions, sessionID)
    keygenSessionsMutex.Unlock()
    
    // Notify parent
    notifyKeygenComplete(sessionID, sealedKey, cardanoAddress, session.signers, nil)
    
    log.Printf("Cardano key generation successfully completed for %s with address %s", 
              sessionID, cardanoAddress)
}

// TssMessage represents a TSS message to be exchanged between enclaves
type TssMessage struct {
    SessionID    string   `json:"sessionId"`
    From         int      `json:"from"`
    To           int      `json:"to"`
    IsBroadcast  bool     `json:"isBroadcast"`
    WireMessage  []byte   `json:"wireMessage"`
    KeyCurve     string   `json:"keyCurve"`
    Signers      []string `json:"signers,omitempty"`
}

// sendParentMessage sends a message to the parent instance
func sendParentMessage(messageType string, data []byte) {
    // In a real implementation, this would use vsock to send a message to the parent
    parentMsg := ParentMessage{
        Type: messageType,
        Data: data,
    }
    
    messageBytes, _ := json.Marshal(parentMsg)
    
    // Get a vsock connection to parent
    conn, err := getParentConnection()
    if err != nil {
        log.Printf("Failed to connect to parent: %v", err)
        return
    }
    defer conn.Close()
    
    // Send the message
    conn.Write(messageBytes)
}

// ParentMessage represents a message to the parent instance
type ParentMessage struct {
    Type string `json:"type"`
    Data []byte `json:"data"`
}

// getParentConnection is now centralized in tss.go to avoid duplication

// notifyKeygenComplete notifies the parent instance that key generation is complete
func notifyKeygenComplete(sessionID string, keyData []byte, address string, signers []string, err error) {
    var result ParentKeygenResult
    
    if err != nil {
        result = ParentKeygenResult{
            SessionID: sessionID,
            Success:   false,
            Error:     err.Error(),
        }
    } else {
        result = ParentKeygenResult{
            SessionID: sessionID,
            Success:   true,
            Address:   address,
            Signers:   signers,
        }
    }
    
    resultBytes, _ := json.Marshal(result)
    sendParentMessage("keygen_result", resultBytes)
}

// ParentKeygenResult is sent to the parent instance when key generation completes
type ParentKeygenResult struct {
    SessionID string   `json:"sessionId"`
    Success   bool     `json:"success"`
    Address   string   `json:"address,omitempty"`
    Signers   []string `json:"signers,omitempty"`
    Error     string   `json:"error,omitempty"`
}

// storeSignersData is defined above

// getCurveType returns the appropriate tss.CurveType for a given key curve
func getCurveType(keyCurve string) tss.CurveType {
    switch keyCurve {
    case "ecdsa", "bitcoin_ecdsa", "secp256k1":
        return tss.S256()
    case "eddsa", "sui_eddsa", "stellar_eddsa", "algorand_eddsa", "ripple_eddsa", "aptos_eddsa", "cardano_eddsa":
        return tss.Edwards()
    default:
        return tss.S256() // Default to secp256k1
    }
}

// getEnclaveID returns the enclave's unique identifier
func getEnclaveID() string {
    // Use the report data from SGX to generate a unique ID
    // This ensures the ID is tied to this specific enclave instance
    reportData, err := getAttestationReport()
    if err != nil {
        log.Printf("Error getting attestation report: %v", err)
        return "unknown-enclave"
    }
    
    // If the attestation document already has an EnclaveID, use it
    if reportData.EnclaveID != "" {
        return reportData.EnclaveID
    }
    
    // Otherwise, generate an ID based on PCR values, which are unique to the enclave
    // PCR0 typically contains the hash of the initial code and data loaded
    var idBytes []byte
    for _, pcrValue := range reportData.PCRs {
        idBytes = append(idBytes, pcrValue...)
        break // Just use the first PCR value for simplicity
    }
    
    if len(idBytes) == 0 {
        // If no PCR values are available, use the TEE type and timestamp as a fallback
        idStr := reportData.TEEType + "-" + reportData.Timestamp.Format(time.RFC3339)
        return fmt.Sprintf("enclave-%x", sha256.Sum256([]byte(idStr)))
    }
    
    // Convert the binary ID to a hex string
    enclaveID := hex.EncodeToString(idBytes)
    
    return enclaveID
}

// getEnclavePublicKey retrieves the enclave's public key from the attestation system
func getEnclavePublicKey() ([]byte, error) {
    // Get the attestation key pair or create one if it doesn't exist
    keyPair, err := getAttestationKeyPair()
    if err != nil {
        return nil, fmt.Errorf("failed to get attestation key pair: %v", err)
    }
    
    // Extract the public key
    publicKey, err := extractPublicKeyBytes(keyPair)
    if err != nil {
        return nil, fmt.Errorf("failed to extract public key: %v", err)
    }
    
    return publicKey, nil
}

// encodePublicKey encodes a public key in the format used in the signers list
func encodePublicKey(publicKey []byte) string {
    // For TSS, we typically use base64 or hex encoding for keys in the signers list
    // Using a standard format helps with interoperability
    encodedKey := base64.StdEncoding.EncodeToString(publicKey)
    
    return encodedKey
}

// determineSignerIndex finds our index in the signers list
func determineSignerIndex(enclaveID string, signers []string) int {
    // Get the enclave's public key from the attestation
    enclavePublicKey, err := getEnclavePublicKey()
    if err != nil {
        log.Printf("Error getting enclave public key: %v", err)
        return -1
    }
    
    // Encode the public key in the same format used in signers list
    encodedPublicKey := encodePublicKey(enclavePublicKey)
    
    // Find our position in the signers list
    for i, signerKey := range signers {
        if signerKey == encodedPublicKey {
            return i
        }
    }
    
    log.Printf("Warning: Enclave public key not found in signers list")
    return -1 // Not found
}

// publicKeyToAddress converts a public key to an Ethereum address
func publicKeyToAddress(pubkey []byte) string {
    var buf []byte
    hash := crypto.Keccak256(pubkey[1:]) // Remove the '04' prefix
    address := "0x" + hex.EncodeToString(hash[12:])
    return address
}

// toHexInt converts a big.Int to a hex string without 0x prefix
func toHexInt(n *big.Int) string {
    return fmt.Sprintf("%x", n)
}

// calculateThreshold calculates the threshold based on the number of parties
func calculateThreshold(numParties int) int {
    if numParties <= 2 {
        return 1
    }
    return (numParties / 2) + 1
}

// getParties creates TSS parties and IDs
func getParties(numParties int) (tss.SortedPartyIDs, []*tss.PartyID) {
    partyIDs := make([]*tss.PartyID, numParties)
    
    for i := 0; i < numParties; i++ {
        partyIDs[i] = tss.NewPartyID(fmt.Sprintf("%d", i), "", big.NewInt(int64(i+1)))
    }
    
    return tss.SortPartyIDs(partyIDs), partyIDs
}