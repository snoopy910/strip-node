package main

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net"
    "os"
    "sync"
    "time"
    "math/big"
)

const (
    // Default vsock port to listen on
    DefaultPort = 5000
    
    // Environment variables
    EnvEnclavePort = "ENCLAVE_PORT"
    EnvKMSKeyID    = "KMS_KEY_ID"
    EnvRegion      = "AWS_REGION"
    EnvEnclaveID   = "ENCLAVE_ID"
    EnvLogLevel    = "LOG_LEVEL"
)

// Global variables
var (
    // Enclave ID (unique identifier for this enclave)
    enclaveID string
    
    // KMS configuration
    region   string
    
    // Active connections
    connections     = make(map[net.Conn]bool)
    connectionsMu   sync.Mutex
)

// Request defines the structure for requests to the enclave
type Request struct {
    RequestID     uint64          `json:"requestId"`
    Operation     string          `json:"operation"`
    Identity      string          `json:"identity,omitempty"`
    IdentityCurve string          `json:"identityCurve,omitempty"`
    KeyCurve      string          `json:"keyCurve,omitempty"`
    Signers       []string        `json:"signers,omitempty"`
    Hash          []byte          `json:"hash,omitempty"`
    SessionID     string          `json:"sessionId,omitempty"`
    From          string          `json:"from,omitempty"`
    Message       []byte          `json:"message,omitempty"`
    Attestation   json.RawMessage `json:"attestation,omitempty"`
    KeyID         string          `json:"keyId,omitempty"`
    MetaData      interface{}     `json:"metadata,omitempty"`
    LogLevel      string          `json:"logLevel,omitempty"`
}

// Response defines the structure for responses from the enclave
type Response struct {
    RequestID     uint64          `json:"requestId"`
    Success       bool            `json:"success"`
    Error         string          `json:"error,omitempty"`
    KeyShare      string          `json:"keyShare,omitempty"`
    Signature     []byte          `json:"signature,omitempty"`
    Address       string          `json:"address,omitempty"`
    SessionID     string          `json:"sessionId,omitempty"`
    SessionState  int             `json:"sessionState,omitempty"`
    Attestation   json.RawMessage `json:"attestation,omitempty"`
    Result        interface{}     `json:"result,omitempty"`
    KeyID         string          `json:"keyId,omitempty"`
}

func main() {
    // Configure logging
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    
    // Initialize the enclave
    if err := initialize(); err != nil {
        log.Fatalf("Failed to initialize enclave: %v", err)
    }
    
    // Get port from environment or use default
    port := DefaultPort
    if portEnv := os.Getenv(EnvEnclavePort); portEnv != "" {
        fmt.Sscanf(portEnv, "%d", &port)
    }
    
    // Listen on vsock
    listenAddr := fmt.Sprintf("vsock:%d", port)
    listener, err := net.Listen("vsock", listenAddr)
    if err != nil {
        log.Fatalf("Failed to listen on vsock: %v", err)
    }
    defer listener.Close()
    
    log.Printf("Enclave server listening on port %d", port)
    
    // Accept connections
    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Printf("Failed to accept connection: %v", err)
            continue
        }
        
        // Register connection
        connectionsMu.Lock()
        connections[conn] = true
        connectionsMu.Unlock()
        
        // Handle connection in a goroutine
        go handleConnection(conn)
    }
}

// initialize sets up the enclave environment
func initialize() error {
    // Get enclave ID from environment
    enclaveID = os.Getenv(EnvEnclaveID)
    if enclaveID == "" {
        // Generate a random ID if not provided
        enclaveID = fmt.Sprintf("enclave-%d", time.Now().UnixNano())
    }
    
    // Get KMS configuration from environment
    region = os.Getenv(EnvRegion)
    if region == "" {
        region = "us-east-1" // Default
    }
    
    // Initialize KMS client
    if err := initKMS(); err != nil {
        return fmt.Errorf("failed to initialize KMS: %w", err)
    }
    
    // Initialize key store
    if err := initKeyStore(); err != nil {
        return fmt.Errorf("failed to initialize key store: %w", err)
    }
    
    // Generate initial attestation document
    if _, err := GetAttestationDocument(); err != nil {
        return fmt.Errorf("failed to generate initial attestation document: %w", err)
    }
    
    log.Printf("Enclave initialized with ID: %s", enclaveID)
    
    return nil
}

// handleConnection processes a connection from the parent instance
func handleConnection(conn net.Conn) {
    defer func() {
        // Clean up connection
        conn.Close()
        connectionsMu.Lock()
        delete(connections, conn)
        connectionsMu.Unlock()
    }()
    
    log.Printf("New connection from parent instance")
    
    // Buffer for reading requests
    buffer := make([]byte, 65536) // 64KB buffer
    
    for {
        // Read request
        n, err := conn.Read(buffer)
        if err != nil {
            if err == io.EOF {
                log.Printf("Connection closed by parent")
                return
            }
            log.Printf("Failed to read from connection: %v", err)
            return
        }
        
        // Parse request
        var req Request
        if err := json.Unmarshal(buffer[:n], &req); err != nil {
            log.Printf("Failed to parse request: %v", err)
            sendError(conn, req.RequestID, "invalid request format")
            continue
        }
        
        // Handle request
        go handleRequest(conn, req)
    }
}

// handleRequest processes a request from the parent instance
func handleRequest(conn net.Conn, req Request) {
    log.Printf("Handling request: %s", req.Operation)
    
    switch req.Operation {
    case "get_address":
        handleGetAddress(conn, req)
    case "create_keygen_session":
        handleCreateKeygenSession(conn, req)
    case "create_signing_session":
        handleCreateSigningSession(conn, req)
    case "get_session_status":
        handleGetSessionStatus(conn, req)
    case "process_session_message":
        handleProcessSessionMessage(conn, req)
    case "get_attestation":
        handleGetAttestation(conn, req)
    case "verify_attestation":
        handleVerifyAttestation(conn, req)
    default:
        sendError(conn, req.RequestID, "unknown operation")
    }
}

// handleGetAddress handles address retrieval requests
func handleGetAddress(conn net.Conn, req Request) {
    // Validate request
    if req.Identity == "" || req.IdentityCurve == "" || req.KeyCurve == "" {
        sendError(conn, req.RequestID, "missing required parameters")
        return
    }
    
    // Generate key ID
    keyID := fmt.Sprintf("%s_%s_%s", req.Identity, req.IdentityCurve, req.KeyCurve)
    
    // Get address from storage
    address, exists := getAddress(keyID)
    if !exists {
        sendError(conn, req.RequestID, "key not found")
        return
    }
    
    // Send response
    sendSuccess(conn, Response{
        RequestID: req.RequestID,
        Address:   address,
    })
}

// handleCreateKeygenSession handles key generation session creation
func handleCreateKeygenSession(conn net.Conn, req Request) {
    // Validate request
    if req.Identity == "" || req.IdentityCurve == "" || req.KeyCurve == "" || len(req.Signers) == 0 {
        sendError(conn, req.RequestID, "missing required parameters")
        return
    }
    
    // Create session
    session, err := CreateSession(SessionTypeKeygen, req.Identity, req.IdentityCurve, req.KeyCurve, req.Signers, nil, nil) // No sealed key for keygen
    if err != nil {
        sendError(conn, req.RequestID, fmt.Sprintf("failed to create session: %v", err))
        return
    }
    
    // Send response
    sendSuccess(conn, Response{
        RequestID:    req.RequestID,
        SessionID:    session.ID,
        SessionState: int(session.State),
    })
    
    // In the background, initiate the key generation protocol
    go initiateKeygenProtocol(session)
}

// handleCreateSigningSession handles signing session creation
func handleCreateSigningSession(conn net.Conn, req Request) {
    // Validate request
    if req.Identity == "" || req.IdentityCurve == "" || req.KeyCurve == "" || len(req.Hash) == 0 {
        sendError(conn, req.RequestID, "missing required parameters")
        return
    }
    
    // Check if key exists
    keyID := fmt.Sprintf("%s_%s_%s", req.Identity, req.IdentityCurve, req.KeyCurve)
    sealedKey, exists := getKeyReference(keyID)
    if !exists {
        sendError(conn, req.RequestID, "key not found")
        return
    }
    
    // Get signers from key metadata
    signers, err := getSignersForKey(keyID)
    if err != nil {
        sendError(conn, req.RequestID, fmt.Sprintf("failed to get signers: %v", err))
        return
    }
    
    // Create session with the sealed key for AWS Nitro
    session, err := CreateSession(SessionTypeSign, req.Identity, req.IdentityCurve, req.KeyCurve, signers, req.Hash, sealedKey)
    if err == nil {
        // Register the key in memory store for AWS Nitro
        // This ensures the key is accessible throughout the session lifecycle
        storeNitroKeyInMemory(session.ID, sealedKey)
    }
    if err != nil {
        sendError(conn, req.RequestID, fmt.Sprintf("failed to create session: %v", err))
        return
    }
    
    // Send response
    sendSuccess(conn, Response{
        RequestID:    req.RequestID,
        SessionID:    session.ID,
        SessionState: int(session.State),
    })
    
    // In the background, initiate the signing protocol
    go initiateSigningProtocol(session)
}

// handleGetSessionStatus handles session status retrieval
func handleGetSessionStatus(conn net.Conn, req Request) {
    // Validate request
    if req.SessionID == "" {
        sendError(conn, req.RequestID, "missing session ID")
        return
    }
    
    // Get session
    session, err := GetSession(req.SessionID)
    if err != nil {
        sendError(conn, req.RequestID, fmt.Sprintf("session not found: %v", err))
        return
    }
    
    // Send response
    sendSuccess(conn, Response{
        RequestID:    req.RequestID,
        SessionID:    session.ID,
        SessionState: int(session.State),
    })
}

// handleProcessSessionMessage handles session message processing
func handleProcessSessionMessage(conn net.Conn, req Request) {
    // Validate request
    if req.SessionID == "" || req.From == "" || len(req.Message) == 0 {
        sendError(conn, req.RequestID, "missing required parameters")
        return
    }
    
    // Process message
    err := ProcessSessionMessage(req.SessionID, req.From, req.Message)
    if err != nil {
        sendError(conn, req.RequestID, fmt.Sprintf("failed to process message: %v", err))
        return
    }
    
    // Send response
    sendSuccess(conn, Response{
        RequestID: req.RequestID,
    })
}

// handleGetAttestation handles attestation document retrieval
func handleGetAttestation(conn net.Conn, req Request) {
    // Get attestation document
    doc, err := GetAttestationDocument()
    if err != nil {
        sendError(conn, req.RequestID, fmt.Sprintf("failed to get attestation document: %v", err))
        return
    }
    
    // Serialize attestation document
    attestationBytes, err := json.Marshal(doc)
    if err != nil {
        sendError(conn, req.RequestID, fmt.Sprintf("failed to marshal attestation document: %v", err))
        return
    }
    
    // Send response
    sendSuccess(conn, Response{
        RequestID:   req.RequestID,
        Attestation: attestationBytes,
    })
}

// handleVerifyAttestation handles attestation verification
func handleVerifyAttestation(conn net.Conn, req Request) {
    // Validate request
    if len(req.Attestation) == 0 {
        sendError(conn, req.RequestID, "missing attestation document")
        return
    }
    
    // Parse attestation document
    var doc AttestationDocument
    if err := json.Unmarshal(req.Attestation, &doc); err != nil {
        sendError(conn, req.RequestID, fmt.Sprintf("failed to parse attestation document: %v", err))
        return
    }
    
    // Verify attestation
    valid, err := VerifyAttestationDocument(&doc)
    if err != nil {
        sendError(conn, req.RequestID, fmt.Sprintf("failed to verify attestation: %v", err))
        return
    }
    
    // Send response
    sendSuccess(conn, Response{
        RequestID: req.RequestID,
        Success:   valid,
    })
}

// sendError sends an error response
func sendError(conn net.Conn, requestID uint64, message string) {
    resp := Response{
        RequestID: requestID,
        Success:   false,
        Error:     message,
    }
    
    data, _ := json.Marshal(resp)
    conn.Write(data)
    log.Printf("Error response sent: %s", message)
}

// sendSuccess sends a successful response to the client
func sendSuccess(conn net.Conn, resp Response) {
    // Ensure success flag is set
    resp.Success = true
    
    data, _ := json.Marshal(resp)
    conn.Write(data)
    log.Printf("Success response sent")
}

// initiateKeygenProtocol starts the key generation protocol
func initiateKeygenProtocol(session *Session) {
    log.Printf("Initiating distributed key generation protocol for session %s", session.ID)
    
    // Update session state to show it's in progress
    UpdateSession(session.ID, SessionStateInProgress, nil, "")
    
    // For AWS Nitro, we need to use the parent process for external communication
    // Initialize the TSS keygen party
    tssParams, err := initTssKeygenParty(session)
    if err != nil {
        log.Printf("Failed to initialize TSS keygen party: %v", err)
        UpdateSession(session.ID, SessionStateFailed, nil, fmt.Sprintf("Failed to initialize TSS keygen party: %v", err))
        return
    }
    
    // Create message handlers for incoming TSS messages
    errCh := make(chan error, len(session.Signers))
    outCh := make(chan Message, len(session.Signers)*3)
    endCh := make(chan *LocalPartySaveData, 1)
    
    // Store these channels in our session for later use
    session.TssData = &TssSessionData{
        ErrChan:   errCh,
        OutChan:   outCh,
        EndChan:   endCh,
        PartyID:   tssParams,
        PublicKey: nil,
    }
    
    // Start the TSS party in a goroutine
    go func() {
        // Create the keygen party using our local implementation
        party, err := createLocalKeygenParty(session, tssParams, outCh, endCh)
        if err != nil {
            log.Printf("Failed to create keygen party: %v", err)
            errCh <- err
            UpdateSession(session.ID, SessionStateFailed, nil, fmt.Sprintf("Failed to create keygen party: %v", err))
            return
        }
        
        // Start the party computation
        if err := party.Start(); err != nil {
            log.Printf("Failed to start keygen party: %v", err)
            errCh <- err
            UpdateSession(session.ID, SessionStateFailed, nil, fmt.Sprintf("Failed to start keygen party: %v", err))
            return
        }
        
        log.Printf("TSS keygen party started for session %s", session.ID)
    }()
    
    // Start message router goroutine to handle outgoing messages
    go routeTssMessages(session, nil)
    
    // Start a goroutine to handle save data when key generation is complete
    go func() {
        select {
        case keyData := <-endCh:
            // Key generation completed successfully
            log.Printf("Key generation completed for session %s", session.ID)
            
            // Calculate the threshold based on signers count (typically t = ⌊(n+1)/2⌋ for t-of-n threshold)
            threshold := (len(session.Signers) + 1) / 2
            
            // Save key data
            keyBytes, err := json.Marshal(keyData)
            if err != nil {
                log.Printf("Failed to marshal key data: %v", err)
                UpdateSession(session.ID, SessionStateFailed, nil, fmt.Sprintf("Failed to marshal key data: %v", err))
                return
            }
            
            // Generate key ID
            keyID := fmt.Sprintf("%s_%s_%s", session.Identity, session.IdentityCurve, session.KeyCurve)
            
            // Create key metadata
            keyMeta := &KeyMetadata{
                KeyID:         keyID,
                Identity:      session.Identity,
                IdentityCurve: session.IdentityCurve,
                KeyCurve:      session.KeyCurve,
                Signers:       session.Signers,
                Threshold:     threshold,
                CreatedAt:     time.Now().Unix(),
                PublicKey:     keyData.ECDSAPub, // Store public key for verification
            }
            
            // Store the key in Nitro enclave memory
            if err := storeKeyShare(keyID, keyBytes, keyMeta); err != nil {
                log.Printf("Failed to store key share: %v", err)
                UpdateSession(session.ID, SessionStateFailed, nil, fmt.Sprintf("Failed to store key share: %v", err))
                return
            }
            
            // Sync key to parent process for persistence
            conn, err := getParentConnection()
            if err != nil {
                log.Printf("Warning: Failed to get parent connection for key sync: %v", err)
            } else {
                if err := syncKeyStoreToParent(conn, keyID); err != nil {
                    log.Printf("Warning: Failed to sync key to parent: %v", err)
                }
            }
            
            // Update session state to completed
            UpdateSession(session.ID, SessionStateCompleted, keyData.ECDSAPub, "")
            
        case err := <-errCh:
            // Key generation failed
            log.Printf("Key generation failed for session %s: %v", session.ID, err)
            UpdateSession(session.ID, SessionStateFailed, nil, fmt.Sprintf("Key generation failed: %v", err))
            
        case <-time.After(10 * time.Minute):
            // Timeout after 10 minutes
            log.Printf("Key generation timed out for session %s", session.ID)
            UpdateSession(session.ID, SessionStateFailed, nil, "Key generation protocol timed out")
        }
    }()
    
    log.Printf("Waiting for keygen messages from other signers for session %s", session.ID)
}

// initiateSigningProtocol starts the signing protocol
func initiateSigningProtocol(session *Session) {
    log.Printf("Initiating distributed signing protocol for session %s", session.ID)
    
    // Update session state to show it's in progress
    UpdateSession(session.ID, SessionStateInProgress, nil, "")
    
    // 1. Retrieve the key for this signing operation from in-memory storage (AWS Nitro specific)
    keyID := fmt.Sprintf("%s_%s_%s", session.Identity, session.IdentityCurve, session.KeyCurve)
    keyShareBytes, err := retrieveKeyShare(keyID)
    if err != nil {
        log.Printf("Failed to retrieve key share: %v", err)
        UpdateSession(session.ID, SessionStateFailed, nil, fmt.Sprintf("Failed to retrieve key share: %v", err))
        return
    }
    
    // Get key metadata for threshold and signers information
    _, err = getKeyMetadata(keyID)
    if err != nil {
        log.Printf("Failed to retrieve key metadata: %v", err)
        UpdateSession(session.ID, SessionStateFailed, nil, fmt.Sprintf("Failed to retrieve key metadata: %v", err))
        return
    }
    
    // Verify hash presence
    if session.Hash == nil || len(session.Hash) == 0 {
        log.Printf("No hash provided for signing")
        UpdateSession(session.ID, SessionStateFailed, nil, "No hash provided for signing")
        return
    }
    
    // Convert key data back to LocalPartySaveData
    var keyData LocalPartySaveData
    if err := json.Unmarshal(keyShareBytes, &keyData); err != nil {
        log.Printf("Failed to unmarshal key data: %v", err)
        UpdateSession(session.ID, SessionStateFailed, nil, fmt.Sprintf("Failed to unmarshal key data: %v", err))
        return
    }
    
    // Initialize TSS signing party
    tssParams, err := initTssSigningParty(session, &keyData)
    if err != nil {
        log.Printf("Failed to initialize TSS signing party: %v", err)
        UpdateSession(session.ID, SessionStateFailed, nil, fmt.Sprintf("Failed to initialize TSS signing party: %v", err))
        return
    }
    
    // Create message handlers for incoming TSS messages
    errCh := make(chan error, len(session.Signers))
    outCh := make(chan Message, len(session.Signers)*3)
    endCh := make(chan SignatureData, 1)
    
    // Store these channels in our session for later use
    session.TssData = &TssSessionData{
        ErrChan:   errCh,
        OutChan:   outCh,
        EndChan:   endCh,
        PartyID:   tssParams,
        PublicKey: keyData.ECDSAPub,
    }
    
    // Get big.Int representation of hash
    hashBigInt := new(big.Int).SetBytes(session.Hash)
    
    // Start the TSS party in a goroutine
    go func() {
        // Create the signing party using our local implementation
        party, err := createLocalSigningParty(session, hashBigInt, tssParams, &keyData, outCh, endCh)
        if err != nil {
            log.Printf("Failed to create signing party: %v", err)
            errCh <- err
            UpdateSession(session.ID, SessionStateFailed, nil, fmt.Sprintf("Failed to create signing party: %v", err))
            return
        }
        
        // Start the party computation
        if err := party.Start(); err != nil {
            log.Printf("Failed to start signing party: %v", err)
            errCh <- err
            UpdateSession(session.ID, SessionStateFailed, nil, fmt.Sprintf("Failed to start signing party: %v", err))
            return
        }
        
        log.Printf("TSS signing party started for session %s", session.ID)
    }()
    
    // Start message router goroutine to handle outgoing messages
    go routeTssMessages(session, nil)
    
    // Start a goroutine to handle the signature when signing is complete
    go func() {
        select {
        case sigData := <-endCh:
            // Signing completed successfully
            log.Printf("Signing completed for session %s", session.ID)
            
            // Convert the signature to DER format as expected by most clients
            // This is especially important for ECDSA signatures
            r := sigData.R
            s := sigData.S
            signature, err := createDERSignature(r, s, session.KeyCurve)
            if err != nil {
                log.Printf("Failed to create DER signature: %v", err)
                UpdateSession(session.ID, SessionStateFailed, nil, fmt.Sprintf("Failed to create DER signature: %v", err))
                return
            }
            
            // Update session state to completed with the signature
            UpdateSession(session.ID, SessionStateCompleted, signature, "")
            
        case err := <-errCh:
            // Signing failed
            log.Printf("Signing failed for session %s: %v", session.ID, err)
            UpdateSession(session.ID, SessionStateFailed, nil, fmt.Sprintf("Signing failed: %v", err))
            
        case <-time.After(5 * time.Minute):
            // Timeout after 5 minutes
            log.Printf("Signing timed out for session %s", session.ID)
            UpdateSession(session.ID, SessionStateFailed, nil, "Signing protocol timed out")
        }
    }()
    
    log.Printf("Waiting for signing messages from other signers for session %s", session.ID)
}

// getSignersForKey retrieves the signers for a key
func getSignersForKey(keyID string) ([]string, error) {
    // Retrieve key metadata from our AWS Nitro enclave storage
    keyMetadata, err := getKeyMetadata(keyID)
    if err != nil {
        // Try to fetch it from the parent process if not found locally
        conn, err := getParentConnection()
        if err != nil {
            return nil, fmt.Errorf("no local key metadata and failed to connect to parent: %w", err)
        }
        
        // Request key metadata from parent
        req := Request{
            RequestID: generateRequestID(),
            Operation: "get_key_metadata",
            KeyID:     keyID,
        }
        
        // Send the request
        reqBytes, err := json.Marshal(req)
        if err != nil {
            return nil, fmt.Errorf("failed to marshal key metadata request: %w", err)
        }
        
        if _, err := conn.Write(reqBytes); err != nil {
            return nil, fmt.Errorf("failed to request key metadata from parent: %w", err)
        }
        
        // Wait for response with a timeout
        respCh := make(chan []byte, 1)
        errCh := make(chan error, 1)
        
        go func() {
            buffer := make([]byte, 8192) // 8KB buffer
            n, err := conn.Read(buffer)
            if err != nil {
                errCh <- err
                return
            }
            respCh <- buffer[:n]
        }()
        
        select {
        case respBytes := <-respCh:
            // Parse the response
            var resp Response
            if err := json.Unmarshal(respBytes, &resp); err != nil {
                return nil, fmt.Errorf("failed to unmarshal key metadata response: %w", err)
            }
            
            if resp.Error != "" {
                return nil, fmt.Errorf("parent returned error for key metadata: %s", resp.Error)
            }
            
            // Extract the key metadata
            var metaBytes []byte
            if resp.Result != nil {
                metaBytes, err = json.Marshal(resp.Result)
                if err != nil {
                    return nil, fmt.Errorf("failed to marshal key metadata from response: %w", err)
                }
            } else {
                return nil, fmt.Errorf("no key metadata in response")
            }
            
            // Unmarshal the key metadata
            if err := json.Unmarshal(metaBytes, &keyMetadata); err != nil {
                return nil, fmt.Errorf("failed to unmarshal key metadata from response: %w", err)
            }
            
            // Cache the metadata locally
            if err := saveKeyMetadata(keyMetadata); err != nil {
                log.Printf("Warning: failed to cache key metadata: %v", err)
                // Continue anyway since we have the metadata now
            }
            
        case err := <-errCh:
            return nil, fmt.Errorf("error reading key metadata response: %w", err)
            
        case <-time.After(5 * time.Second):
            return nil, fmt.Errorf("timeout waiting for key metadata response")
        }
    }
    
    // Check if we have signer information even after potential parent fetch
    if keyMetadata == nil || len(keyMetadata.Signers) == 0 {
        // If we don't have any signer information, return just this enclave as the signer
        // This handles backward compatibility with single-signer keys
        log.Printf("No signers found for key %s, using enclave ID only", keyID)
        return []string{enclaveID}, nil
    }
    
    // Validate that this enclave is one of the signers
    isSigner := false
    for _, signer := range keyMetadata.Signers {
        if signer == enclaveID {
            isSigner = true
            break
        }
    }
    
    if !isSigner {
        return nil, fmt.Errorf("this enclave (%s) is not authorized for key %s", enclaveID, keyID)
    }
    
    return keyMetadata.Signers, nil
}

// TssParty represents a local party in the TSS protocol
type TssParty struct {
    // Fields needed for the party implementation
}

// Start begins the TSS party's computation
func (p *TssParty) Start() error {
    // Implementation for starting the party computation
    return nil
}

// createLocalKeygenParty creates a local keygen party
func createLocalKeygenParty(session *Session, tssParams *PartyID, outCh chan Message, endCh chan *LocalPartySaveData) (*TssParty, error) {
    // Create a new local party implementation
    threshold := (len(session.Signers)*2)/3
    if threshold < 1 {
        threshold = 1
    }
    
    log.Printf("Starting TSS keygen with threshold %d out of %d participants", 
              threshold, len(session.Signers))
              
    // For now, return a minimal implementation
    return &TssParty{}, nil
}

// createLocalSigningParty creates a local signing party
func createLocalSigningParty(session *Session, hashBigInt *big.Int, tssParams *PartyID, keyData *LocalPartySaveData, outCh chan Message, endCh chan SignatureData) (*TssParty, error) {
    log.Printf("Creating TSS signing party for session %s", session.ID)
    
    // For now, return a minimal implementation
    return &TssParty{}, nil
}