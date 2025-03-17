package main

import (
    "bytes"
    "context"
    "crypto/ed25519"
    "crypto/sha512"
    "encoding/base32"
    "encoding/base64"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "log"
    "math/big"
    "net"
    "strings"
    "sync"
    "time"
    
    "github.com/StripChain/strip-node/bitcoin"
    "github.com/StripChain/strip-node/dogecoin"
    "github.com/StripChain/strip-node/ripple"
    "github.com/bnb-chain/tss-lib/v2/common"
    ecdsaKeygen "github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
    ecdsaSigning "github.com/bnb-chain/tss-lib/v2/ecdsa/signing"
    eddsaKeygen "github.com/bnb-chain/tss-lib/v2/eddsa/keygen"
    eddsaSigning "github.com/bnb-chain/tss-lib/v2/eddsa/signing"
    "github.com/bnb-chain/tss-lib/v2/tss"
    "github.com/coming-chat/go-sui/v2/lib"
    "github.com/coming-chat/go-sui/v2/sui_types"
    "github.com/decred/dcrd/dcrec/edwards/v2"
    "github.com/ethereum/go-ethereum/crypto"
    "github.com/mr-tron/base58"
    "github.com/stellar/go/strkey"
    "golang.org/x/crypto/blake2b"
)

// Curve constants to match parent implementation
const (
    ECDSA_CURVE      = "ecdsa"
    EDDSA_CURVE      = "eddsa"
    BITCOIN_CURVE    = "bitcoin_ecdsa"
    SUI_EDDSA_CURVE  = "sui_eddsa"
    SECP256K1_CURVE  = "secp256k1"
    APTOS_EDDSA_CURVE = "aptos_eddsa"
    STELLAR_CURVE    = "stellar_eddsa"
    ALGORAND_CURVE   = "algorand_eddsa"
    RIPPLE_CURVE     = "ripple_eddsa"
    CARDANO_CURVE    = "cardano_eddsa"
)

// SigningSession represents an active signing operation
type SigningSession struct {
    identity      string
    identityCurve string
    keyCurve      string
    hash          []byte
    party         *tss.Party
    outChan       chan tss.Message
    saveChan      chan *common.SignatureData
    signers       []string
    startTime     time.Time
}

// ActiveSessions tracks ongoing signing operations
var (
    signingSession     = map[string]*SigningSession{}
    signingSessionMutex sync.RWMutex
)

// SessionTimeout is the maximum time a signing session can be active
const SessionTimeout = 5 * time.Minute

// handleSign handles incoming signing requests from the parent instance
func handleSign(conn net.Conn, req Request) {
    log.Printf("Handling signing for %s_%s_%s", req.Identity, req.IdentityCurve, req.KeyCurve)
    
    // Validate request
    if req.Identity == "" || req.IdentityCurve == "" || req.KeyCurve == "" || len(req.Hash) == 0 {
        sendError(conn, "missing required parameters")
        return
    }
    
    // Get key ID
    keyID := fmt.Sprintf("%s_%s_%s", req.Identity, req.IdentityCurve, req.KeyCurve)
    
    // Get sealed key
    sealedKey, exists := getKeyReference(keyID)
    if !exists {
        sendError(conn, "key not found")
        return
    }
    
    // Get signers information
    signersBytes, exists := getSignersDataForSigning(keyID)
    if !exists {
        sendError(conn, "signers information not found")
        return
    }
    
    var signers []string
    if err := json.Unmarshal(signersBytes, &signers); err != nil {
        log.Printf("Failed to unmarshal signers: %v", err)
        sendError(conn, "invalid signers format")
        return
    }
    
    // Determine our signer index (using enclave ID as a placeholder)
    enclaveID := getEnclaveID()
    signerIndex := determineSignerIndex(enclaveID, signers)
    
    // Unseal the key
    keyShare, err := unsealKeyShare(context.Background(), sealedKey, keyID)
    if err != nil {
        log.Printf("Failed to unseal key share: %v", err)
        sendError(conn, fmt.Sprintf("failed to unseal key share: %v", err))
        return
    }
    
    // Initialize the TSS signing party
    sessionID := fmt.Sprintf("%s_%x", keyID, req.Hash[:8])
    
    // Check if we already have an active session
    signingSessionMutex.RLock()
    existingSession := signingSession[sessionID]
    signingSessionMutex.RUnlock()
    
    if existingSession != nil {
        // We're already processing this signing request
        log.Printf("Signing session %s already active", sessionID)
        sendSuccess(conn, Response{
            Status: "in_progress",
        })
        return
    }
    
    // Start a new TSS signing session
    session, err := startSigningSession(keyID, req.Identity, req.IdentityCurve, req.KeyCurve, 
                                       req.Hash, keyShare, signers, signerIndex)
    if err != nil {
        log.Printf("Failed to start signing session: %v", err)
        sendError(conn, fmt.Sprintf("failed to start signing session: %v", err))
        return
    }
    
    // Register the session
    signingSessionMutex.Lock()
    signingSession[sessionID] = session
    signingSessionMutex.Unlock()
    
    // Start session cleanup timer
    go func() {
        time.Sleep(SessionTimeout)
        signingSessionMutex.Lock()
        delete(signingSession, sessionID)
        signingSessionMutex.Unlock()
        log.Printf("Signing session %s timed out", sessionID)
    }()
    
    // Tell the parent that we've started the signing process
    sendSuccess(conn, Response{
        Status: "started",
    })
    
    // Handle session messages in the background
    go handleSigningSession(sessionID, session)
}

// handleSigningSession processes the signing session and generates the signature
func handleSigningSession(sessionID string, session *SigningSession) {
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
            signMessage := TssMessage{
                SessionID:    sessionID,
                From:         msg.GetFrom().Index,
                To:           to,
                IsBroadcast:  msg.IsBroadcast(),
                WireMessage:  wireBytes,
                KeyCurve:     session.keyCurve,
                MessageHash:  session.hash,
            }
            
            messageBytes, _ := json.Marshal(signMessage)
            sendParentMessage("tss_message", messageBytes)
            log.Printf("Sent TSS message from party %d to %d (broadcast: %v)", 
                      signMessage.From, signMessage.To, signMessage.IsBroadcast)
            
        case saveData := <-session.saveChan:
            // Signing completed successfully
            completed = true
            
            // Determine the final signature format based on the key curve
            signature, address, err := formatSignature(session, saveData)
            if err != nil {
                log.Printf("Failed to format signature: %v", err)
                notifySigningComplete(sessionID, nil, "", err)
                return
            }
            
            // Notify about successful signing
            log.Printf("Signature generated successfully for %s", sessionID)
            notifySigningComplete(sessionID, signature, address, nil)
            
            // Clean up the session
            signingSessionMutex.Lock()
            delete(signingSession, sessionID)
            signingSessionMutex.Unlock()
            
        case <-time.After(SessionTimeout):
            // Timeout case
            completed = true
            log.Printf("Signing session %s timed out while processing", sessionID)
            notifySigningComplete(sessionID, nil, "", fmt.Errorf("signing timed out"))
            
            // Clean up the session
            signingSessionMutex.Lock()
            delete(signingSession, sessionID)
            signingSessionMutex.Unlock()
        }
    }
}

// handleTssMessage processes TSS messages forwarded from the parent instance
func handleTssMessage(message []byte) {
    var tssMessage TssMessage
    if err := json.Unmarshal(message, &tssMessage); err != nil {
        log.Printf("Failed to unmarshal TSS message: %v", err)
        return
    }
    
    // Find the session
    signingSessionMutex.RLock()
    session := signingSession[tssMessage.SessionID]
    signingSessionMutex.RUnlock()
    
    if session == nil {
        log.Printf("Received message for unknown session %s", tssMessage.SessionID)
        return
    }
    
    // Process the message
    parties, _ := getParties(len(session.signers))
    
    // Parse the wire message
    pMsg, err := tss.ParseWireMessage(tssMessage.WireMessage, 
                                     parties[tssMessage.From], tssMessage.IsBroadcast)
    if err != nil {
        log.Printf("Failed to parse wire message: %v", err)
        return
    }
    
    // Get the party
    party := session.party
    
    // Check if this message is for us
    ourIndex := party.PartyID().Index
    if tssMessage.To != -1 && tssMessage.To != ourIndex {
        // Message not intended for us
        return
    }
    
    // Prevent message loops
    if tssMessage.From == ourIndex {
        // Message from ourselves
        return
    }
    
    // Update the party with the message
    log.Printf("Processing TSS message from party %d to %d (broadcast: %v)", 
              tssMessage.From, tssMessage.To, tssMessage.IsBroadcast)
    
    _, err = party.Update(pMsg)
    if err != nil {
        log.Printf("Failed to update party with message: %v", err)
        return
    }
    
    log.Printf("Successfully processed TSS message for session %s", tssMessage.SessionID)
}

// startSigningSession initializes a new signing session using the key share
func startSigningSession(keyID, identity, identityCurve, keyCurve string, 
                        hash []byte, keyShare []byte, signers []string, signerIndex int) (*SigningSession, error) {
    
    // Create TSS parties
    totalSigners := len(signers)
    parties, partyIDs := getParties(totalSigners)
    
    // Calculate threshold
    threshold := calculateThreshold(totalSigners)
    
    // Create TSS peer context
    ctx := tss.NewPeerContext(parties)
    
    // Channels for communication
    outChan := make(chan tss.Message)
    saveChan := make(chan *common.SignatureData)
    
    // Create a new signing session
    session := &SigningSession{
        identity:      identity,
        identityCurve: identityCurve,
        keyCurve:      keyCurve,
        hash:          hash,
        outChan:       outChan,
        saveChan:      saveChan,
        signers:       signers,
        startTime:     time.Now(),
    }
    
    // Initialize the party based on the key curve
    var err error
    switch keyCurve {
    case "ecdsa":
        // Standard ECDSA for Ethereum, etc.
        var ecdsaData ecdsaKeygen.LocalPartySaveData
        if err := json.Unmarshal(keyShare, &ecdsaData); err != nil {
            return nil, fmt.Errorf("failed to unmarshal ECDSA key data: %w", err)
        }
        
        msg, _ := new(big.Int).SetString(string(hash), 16)
        params := tss.NewParameters(tss.S256(), ctx, partyIDs[signerIndex], len(parties), threshold)
        
        localParty := ecdsaSigning.NewLocalParty(msg, params, ecdsaData, outChan, saveChan)
        session.party = localParty
        
        // Start the party
        go localParty.Start()
        
    case "bitcoin_ecdsa":
        // Bitcoin-specific ECDSA
        var ecdsaData ecdsaKeygen.LocalPartySaveData
        if err := json.Unmarshal(keyShare, &ecdsaData); err != nil {
            return nil, fmt.Errorf("failed to unmarshal Bitcoin key data: %w", err)
        }
        
        msg, _ := new(big.Int).SetString(string(hash), 16)
        params := tss.NewParameters(tss.S256(), ctx, partyIDs[signerIndex], len(parties), threshold)
        
        localParty := ecdsaSigning.NewLocalParty(msg, params, ecdsaData, outChan, saveChan)
        session.party = localParty
        
        // Start the party
        go localParty.Start()
        
    case "eddsa":
        // EdDSA (Solana, etc.)
        var eddsaData eddsaKeygen.LocalPartySaveData
        if err := json.Unmarshal(keyShare, &eddsaData); err != nil {
            return nil, fmt.Errorf("failed to unmarshal EdDSA key data: %w", err)
        }
        
        msg := new(big.Int).SetBytes(hash)
        params := tss.NewParameters(tss.Edwards(), ctx, partyIDs[signerIndex], len(parties), threshold)
        
        localParty := eddsaSigning.NewLocalParty(msg, params, eddsaData, outChan, saveChan)
        session.party = localParty
        
        // Start the party
        go localParty.Start()
        
    case "secp256k1":
        // Secp256k1 (Dogecoin, etc.)
        var ecdsaData ecdsaKeygen.LocalPartySaveData
        if err := json.Unmarshal(keyShare, &ecdsaData); err != nil {
            return nil, fmt.Errorf("failed to unmarshal Secp256k1 key data: %w", err)
        }
        
        msg := new(big.Int).SetBytes(crypto.Keccak256(hash))
        params := tss.NewParameters(tss.S256(), ctx, partyIDs[signerIndex], len(parties), threshold)
        
        localParty := ecdsaSigning.NewLocalParty(msg, params, ecdsaData, outChan, saveChan)
        session.party = localParty
        
        // Start the party
        go localParty.Start()
        
    case "sui_eddsa":
        // Sui Ed25519
        var eddsaData eddsaKeygen.LocalPartySaveData
        if err := json.Unmarshal(keyShare, &eddsaData); err != nil {
            return nil, fmt.Errorf("failed to unmarshal Sui EdDSA key data: %w", err)
        }
        
        msg := new(big.Int).SetBytes(hash)
        params := tss.NewParameters(tss.Edwards(), ctx, partyIDs[signerIndex], len(parties), threshold)
        
        localParty := eddsaSigning.NewLocalParty(msg, params, eddsaData, outChan, saveChan)
        session.party = localParty
        
        // Start the party
        go localParty.Start()
        
    case "aptos_eddsa", "stellar_eddsa", "ripple_eddsa", "cardano_eddsa", "algorand_eddsa":
        // Other EdDSA curves
        var eddsaData eddsaKeygen.LocalPartySaveData
        if err := json.Unmarshal(keyShare, &eddsaData); err != nil {
            return nil, fmt.Errorf("failed to unmarshal EdDSA key data: %w", err)
        }
        
        msg := new(big.Int).SetBytes(hash)
        params := tss.NewParameters(tss.Edwards(), ctx, partyIDs[signerIndex], len(parties), threshold)
        
        localParty := eddsaSigning.NewLocalParty(msg, params, eddsaData, outChan, saveChan)
        session.party = localParty
        
        // Start the party
        go localParty.Start()
        
    default:
        return nil, fmt.Errorf("unsupported key curve: %s", keyCurve)
    }
    
    return session, nil
}

// formatSignature prepares the final signature format based on the key curve
func formatSignature(session *SigningSession, saveData *common.SignatureData) ([]byte, string, error) {
    // Retrieve key data
    keyID := fmt.Sprintf("%s_%s_%s", session.identity, session.identityCurve, session.keyCurve)
    sealedKey, exists := getKeyReference(keyID)
    if !exists {
        return nil, "", fmt.Errorf("key not found")
    }
    
    keyShare, err := unsealKeyShare(context.Background(), sealedKey, keyID)
    if err != nil {
        return nil, "", fmt.Errorf("failed to unseal key: %w", err)
    }
    
    switch session.keyCurve {
    case ECDSA_CURVE:
        // Standard ECDSA (Ethereum, etc.)
        var ecdsaData ecdsaKeygen.LocalPartySaveData
        if err := json.Unmarshal(keyShare, &ecdsaData); err != nil {
            return nil, "", fmt.Errorf("failed to unmarshal ECDSA key data: %w", err)
        }
        
        // Format signature for Ethereum (r || s || v)
        signature := append(saveData.R.Bytes(), saveData.S.Bytes()...)
        signature = append(signature, byte(saveData.Recovery))
        
        // Derive Ethereum address
        x := toHexInt(ecdsaData.ECDSAPub.X())
        y := toHexInt(ecdsaData.ECDSAPub.Y())
        publicKeyBytes, _ := hex.DecodeString("04" + x + y)
        address := publicKeyToAddress(publicKeyBytes)
        
        return signature, address, nil
        
    case BITCOIN_CURVE:
        // Bitcoin-specific ECDSA
        var ecdsaData ecdsaKeygen.LocalPartySaveData
        if err := json.Unmarshal(keyShare, &ecdsaData); err != nil {
            return nil, "", fmt.Errorf("failed to unmarshal Bitcoin key data: %w", err)
        }
        
        // Format signature for Bitcoin
        signature := hex.EncodeToString(saveData.Signature)
        
        // Get Bitcoin address
        x := toHexInt(ecdsaData.ECDSAPub.X())
        y := toHexInt(ecdsaData.ECDSAPub.Y())
        uncompressedPubKeyStr := "04" + x + y
        compressedPubKeyStr, err := bitcoin.ConvertToCompressedPublicKey(uncompressedPubKeyStr)
        if err != nil {
            return nil, "", fmt.Errorf("failed to get Bitcoin address: %w", err)
        }
        
        return []byte(signature), compressedPubKeyStr, nil
        
    case EDDSA_CURVE:
        // EdDSA (Solana, etc.)
        var eddsaData eddsaKeygen.LocalPartySaveData
        if err := json.Unmarshal(keyShare, &eddsaData); err != nil {
            return nil, "", fmt.Errorf("failed to unmarshal EdDSA key data: %w", err)
        }
        
        // Get EdDSA public key
        pk := edwards.PublicKey{
            Curve: tss.Edwards(),
            X:     eddsaData.EDDSAPub.X(),
            Y:     eddsaData.EDDSAPub.Y(),
        }
        
        // Get Solana-style address
        address := base58.Encode(pk.Serialize())
        
        return saveData.Signature, address, nil
        
    case SECP256K1_CURVE:
        // Secp256k1 (Dogecoin, etc.)
        var ecdsaData ecdsaKeygen.LocalPartySaveData
        if err := json.Unmarshal(keyShare, &ecdsaData); err != nil {
            return nil, "", fmt.Errorf("failed to unmarshal Secp256k1 key data: %w", err)
        }
        
        // Format signature
        signature := hex.EncodeToString(saveData.Signature) + hex.EncodeToString([]byte{byte(saveData.Recovery)})
        
        // Get address
        x := toHexInt(ecdsaData.ECDSAPub.X())
        y := toHexInt(ecdsaData.ECDSAPub.Y())
        publicKeyStr := "04" + x + y
        
        // Get chain information from metadata if available
        var address string
        // Check if we need testnet address format
        if strings.HasSuffix(keyID, "_testnet") {
            address, err = dogecoin.PublicKeyToTestnetAddress(publicKeyStr)
        } else {
            address, err = dogecoin.PublicKeyToAddress(publicKeyStr)
        }
        
        if err != nil {
            return nil, "", fmt.Errorf("failed to get Dogecoin address: %w", err)
        }
        
        return []byte(signature), address, nil
        
    case SUI_EDDSA_CURVE:
        // Sui Ed25519
        var eddsaData eddsaKeygen.LocalPartySaveData
        if err := json.Unmarshal(keyShare, &eddsaData); err != nil {
            return nil, "", fmt.Errorf("failed to unmarshal Sui EdDSA key data: %w", err)
        }
        
        // Get Ed25519 public key
        pk := edwards.PublicKey{
            Curve: tss.Edwards(),
            X:     eddsaData.EDDSAPub.X(),
            Y:     eddsaData.EDDSAPub.Y(),
        }
        
        // Serialize the full Ed25519 public key
        pkBytes := pk.Serialize()
        
        // Convert to Sui address format
        flag := byte(0x00)
        hasher, _ := blake2b.New256(nil)
        hasher.Write([]byte{flag})
        hasher.Write(pkBytes)
        arr := hasher.Sum(nil)
        suiAddress := "0x" + hex.EncodeToString(arr)
        
        // Format for Sui with scheme flag + signature + public key
        var signatureBytes [ed25519.PublicKeySize + ed25519.SignatureSize + 1]byte
        signatureBuffer := bytes.NewBuffer([]byte{})
        scheme := sui_types.SignatureScheme{ED25519: &lib.EmptyEnum{}}
        signatureBuffer.WriteByte(scheme.Flag())
        signatureBuffer.Write(saveData.Signature)
        signatureBuffer.Write(pkBytes[:])
        copy(signatureBytes[:], signatureBuffer.Bytes())
        
        signatureBase64 := base64.StdEncoding.EncodeToString(signatureBytes[:])
        
        return []byte(signatureBase64), suiAddress, nil
        
    case APTOS_EDDSA_CURVE:
        // Aptos EdDSA
        var eddsaData eddsaKeygen.LocalPartySaveData
        if err := json.Unmarshal(keyShare, &eddsaData); err != nil {
            return nil, "", fmt.Errorf("failed to unmarshal Aptos EdDSA key data: %w", err)
        }
        
        // Get Ed25519 public key
        pk := edwards.PublicKey{
            Curve: tss.Edwards(),
            X:     eddsaData.EDDSAPub.X(),
            Y:     eddsaData.EDDSAPub.Y(),
        }
        
        // Format for Aptos
        publicKeyStr := hex.EncodeToString(pk.Serialize())
        address := "0x" + publicKeyStr
        
        return saveData.Signature, address, nil
        
    case STELLAR_CURVE:
        // Stellar EdDSA
        var eddsaData eddsaKeygen.LocalPartySaveData
        if err := json.Unmarshal(keyShare, &eddsaData); err != nil {
            return nil, "", fmt.Errorf("failed to unmarshal Stellar EdDSA key data: %w", err)
        }
        
        // Get Ed25519 public key
        pk := edwards.PublicKey{
            Curve: tss.Edwards(),
            X:     eddsaData.EDDSAPub.X(),
            Y:     eddsaData.EDDSAPub.Y(),
        }
        
        // Get public key bytes
        pkBytes := pk.Serialize()
        if len(pkBytes) != 32 {
            return nil, "", fmt.Errorf("invalid public key length for Stellar")
        }
        
        // Version byte for ED25519 public key in Stellar
        versionByte := strkey.VersionByteAccountID
        
        // Encode address
        address, err := strkey.Encode(versionByte, pkBytes)
        if err != nil {
            return nil, "", fmt.Errorf("failed to encode Stellar address: %w", err)
        }
        
        return saveData.Signature, address, nil
        
    case ALGORAND_CURVE:
        // Algorand EdDSA
        var eddsaData eddsaKeygen.LocalPartySaveData
        if err := json.Unmarshal(keyShare, &eddsaData); err != nil {
            return nil, "", fmt.Errorf("failed to unmarshal Algorand EdDSA key data: %w", err)
        }
        
        // Get Ed25519 public key
        pk := edwards.PublicKey{
            Curve: tss.Edwards(),
            X:     eddsaData.EDDSAPub.X(),
            Y:     eddsaData.EDDSAPub.Y(),
        }
        
        // Get public key bytes
        pkBytes := pk.Serialize()
        
        // Calculate checksum
        hasher := sha512.New512_256()
        hasher.Write(pkBytes)
        checksum := hasher.Sum(nil)[28:] // Last 4 bytes
        
        // Concatenate public key and checksum
        addressBytes := append(pkBytes, checksum...)
        
        // Encode in base32 without padding
        address := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(addressBytes)
        
        return saveData.Signature, address, nil
        
    case RIPPLE_CURVE:
        // Ripple EdDSA
        var eddsaData eddsaKeygen.LocalPartySaveData
        if err := json.Unmarshal(keyShare, &eddsaData); err != nil {
            return nil, "", fmt.Errorf("failed to unmarshal Ripple EdDSA key data: %w", err)
        }
        
        // Format signature
        address := ripple.PublicKeyToAddress(&eddsaData)
        
        return saveData.Signature, address, nil
        
    case CARDANO_CURVE:
        // Cardano EdDSA
        var eddsaData eddsaKeygen.LocalPartySaveData
        if err := json.Unmarshal(keyShare, &eddsaData); err != nil {
            return nil, "", fmt.Errorf("failed to unmarshal Cardano EdDSA key data: %w", err)
        }
        
        // Get Ed25519 public key
        pk := edwards.PublicKey{
            Curve: tss.Edwards(),
            X:     eddsaData.EDDSAPub.X(),
            Y:     eddsaData.EDDSAPub.Y(),
        }
        
        // Format for Cardano
        publicKeyStr := hex.EncodeToString(pk.Serialize())
        
        return saveData.Signature, publicKeyStr, nil
        
    default:
        return nil, "", fmt.Errorf("unsupported key curve: %s", session.keyCurve)
    }
}

// getSignersDataForSigning retrieves the signers data for a key
func getSignersDataForSigning(keyID string) ([]byte, bool) {
    // Use the metaData variable defined in keygen.go
    metaData.mu.RLock()
    defer metaData.mu.RUnlock()
    
    signersKey := keyID + "_signers"
    signersData, exists := metaData.data[signersKey]
    return signersData, exists
}

// TssMessage represents a TSS message to be exchanged between enclaves via parent instances
type TssMessage struct {
    SessionID    string `json:"sessionId"`
    From         int    `json:"from"`
    To           int    `json:"to"`
    IsBroadcast  bool   `json:"isBroadcast"`
    WireMessage  []byte `json:"wireMessage"`
    KeyCurve     string `json:"keyCurve"`
    MessageHash  []byte `json:"messageHash"`
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

// getParentConnection gets a vsock connection to the parent instance
func getParentConnection() (net.Conn, error) {
    // This is a simplified version - in reality, you'd use a connection pool
    // or maintain a long-lived connection to the parent
    conn, err := net.Dial("vsock", "parent:5001")
    if err != nil {
        return nil, err
    }
    return conn, nil
}

// getEnclaveID returns the enclave's unique identifier
func getEnclaveID() string {
    // In a real implementation, this would return a unique ID for this enclave
    // For now, we'll return a placeholder
    return "enclave-1"
}

// determineSignerIndex finds our index in the signers list
func determineSignerIndex(enclaveID string, signers []string) int {
    // In a real implementation, this would map the enclave ID to the signer index
    // For now, we'll return a placeholder index
    return 0
}

// notifySigningComplete notifies the parent instance that signing is complete
func notifySigningComplete(sessionID string, signature []byte, address string, err error) {
    var result ParentSigningResult
    
    if err != nil {
        result = ParentSigningResult{
            SessionID: sessionID,
            Success:   false,
            Error:     err.Error(),
        }
    } else {
        result = ParentSigningResult{
            SessionID: sessionID,
            Success:   true,
            Signature: signature,
            Address:   address,
        }
        
        // Set Algorand flags if this is an Algorand signature
        if strings.Contains(sessionID, ALGORAND_CURVE) {
            // Default flag for normal transactions
            result.AlgorandFlags = 1
        }
    }
    
    resultBytes, _ := json.Marshal(result)
    sendParentMessage("signing_result", resultBytes)
}

// ParentSigningResult is sent to the parent instance when signing completes
type ParentSigningResult struct {
    SessionID     string `json:"sessionId"`
    Success       bool   `json:"success"`
    Signature     []byte `json:"signature,omitempty"`
    Address       string `json:"address,omitempty"`
    AlgorandFlags byte   `json:"algorandFlags,omitempty"`
    Error         string `json:"error,omitempty"`
}

// calculateThreshold calculates the threshold based on the number of parties
func calculateThreshold(numParties int) int {
    if numParties <= 2 {
        return 1
    }
    return (numParties / 2) + 1
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