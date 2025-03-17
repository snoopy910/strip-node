package main

import (
    "encoding/json"
    "fmt"
    "log"
    "sync"
    "time"
)

// SessionType identifies the type of session
type SessionType int

const (
    SessionTypeKeygen SessionType = iota + 1
    SessionTypeSign
)

// SessionState tracks the state of a session
type SessionState int

const (
    SessionStateInitialized SessionState = iota
    SessionStateInProgress
    SessionStateCompleted
    SessionStateFailed
)

// Session represents an active TEE operation session
type Session struct {
    ID             string               `json:"id"`
    Type           SessionType          `json:"type"`
    State          SessionState         `json:"state"`
    Identity       string               `json:"identity"`
    IdentityCurve  string               `json:"identityCurve"`
    KeyCurve       string               `json:"keyCurve"`
    Signers        []string             `json:"signers"`
    Participants   map[string]bool      `json:"participants"` // Map of signer ID to participation status
    Hash           []byte               `json:"hash,omitempty"`
    CreatedAt      time.Time            `json:"createdAt"`
    LastActivityAt time.Time            `json:"lastActivityAt"`
    Result         interface{}          `json:"result,omitempty"`
    Error          string               `json:"error,omitempty"`
    Channels       map[string]*SecureChannel `json:"-"` // Not serialized
    SealedKey      []byte               `json:"-"` // AWS Nitro in-memory key storage - not serialized
    TssData        *TssSessionData      `json:"-"` // TSS protocol data - not serialized
}

var (
    // Active sessions
    sessions     = make(map[string]*Session)
    sessionsMu   sync.RWMutex
    
    // Session timeout in minutes
    sessionTimeout = 10 * time.Minute
)

// CreateSession creates a new session for multi-party operations
func CreateSession(sessionType SessionType, identity, identityCurve, keyCurve string, 
                 signers []string, hash []byte, sealedKey []byte) (*Session, error) {
    // Generate a unique session ID
    sessionID := generateSessionID()
    
    // Create the session
    session := &Session{
        ID:             sessionID,
        Type:           sessionType,
        State:          SessionStateInitialized,
        Identity:       identity,
        IdentityCurve:  identityCurve,
        KeyCurve:       keyCurve,
        Signers:        signers,
        Participants:   make(map[string]bool),
        Hash:           hash,
        CreatedAt:      time.Now(),
        LastActivityAt: time.Now(),
        Channels:       make(map[string]*SecureChannel),
        SealedKey:      sealedKey, // Store key in memory for AWS Nitro
    }
    
    // Mark ourselves as a participant
    session.Participants[enclaveID] = true
    
    // Store the session
    sessionsMu.Lock()
    sessions[sessionID] = session
    sessionsMu.Unlock()
    
    // Start a cleanup goroutine
    go func() {
        time.Sleep(sessionTimeout)
        cleanupSession(sessionID)
    }()
    
    log.Printf("Created session %s of type %d", sessionID, sessionType)
    
    return session, nil
}

// GetSession retrieves a session by ID
func GetSession(sessionID string) (*Session, error) {
    sessionsMu.RLock()
    session, exists := sessions[sessionID]
    sessionsMu.RUnlock()
    
    if !exists {
        return nil, fmt.Errorf("session not found: %s", sessionID)
    }
    
    return session, nil
}

// UpdateSession updates a session's state and activity time
func UpdateSession(sessionID string, state SessionState, result interface{}, errMsg string) error {
    sessionsMu.Lock()
    defer sessionsMu.Unlock()
    
    session, exists := sessions[sessionID]
    if !exists {
        return fmt.Errorf("session not found: %s", sessionID)
    }
    
    session.State = state
    session.LastActivityAt = time.Now()
    
    if result != nil {
        session.Result = result
    }
    
    if errMsg != "" {
        session.Error = errMsg
    }
    
    return nil
}

// ProcessSessionMessage processes a message for a session
func ProcessSessionMessage(sessionID string, from string, message []byte) error {
    sessionsMu.RLock()
    session, exists := sessions[sessionID]
    sessionsMu.RUnlock()
    
    if !exists {
        return fmt.Errorf("session not found: %s", sessionID)
    }
    
    // Update activity time
    session.LastActivityAt = time.Now()
    
    // Handle message based on session type
    switch session.Type {
    case SessionTypeKeygen:
        return processKeygenMessage(session, from, message)
    case SessionTypeSign:
        return processSigningMessage(session, from, message)
    default:
        return fmt.Errorf("unsupported session type: %d", session.Type)
    }
}

// processKeygenMessage processes a message for a key generation session
func processKeygenMessage(session *Session, from string, message []byte) error {
    // In a real implementation, this would:
    // 1. Parse the message
    // 2. Update the TSS protocol state
    // 3. Generate and send responses if needed
    
    // This is a placeholder implementation
    log.Printf("Processing keygen message from %s for session %s", from, session.ID)
    
    var tssMessage struct {
        Type      string          `json:"type"`
        Data      json.RawMessage `json:"data"`
    }
    
    if err := json.Unmarshal(message, &tssMessage); err != nil {
        return fmt.Errorf("failed to parse message: %w", err)
    }
    
    // Mark the sender as a participant
    session.Participants[from] = true
    
    // Check if all participants have joined
    if len(session.Participants) >= len(session.Signers) {
        log.Printf("All participants joined session %s", session.ID)
        
        // In a real implementation, this would trigger the next step in the TSS protocol
        session.State = SessionStateInProgress
    }
    
    return nil
}

// processSigningMessage processes a message for a signing session
func processSigningMessage(session *Session, from string, message []byte) error {
    // Similar to processKeygenMessage, but for signing
    log.Printf("Processing signing message from %s for session %s", from, session.ID)
    
    var tssMessage struct {
        Type      string          `json:"type"`
        Data      json.RawMessage `json:"data"`
    }
    
    if err := json.Unmarshal(message, &tssMessage); err != nil {
        return fmt.Errorf("failed to parse message: %w", err)
    }
    
    // Mark the sender as a participant
    session.Participants[from] = true
    
    // Check if all participants have joined
    if len(session.Participants) >= len(session.Signers) {
        log.Printf("All participants joined session %s", session.ID)
        
        // In a real implementation, this would trigger the next step in the TSS protocol
        session.State = SessionStateInProgress
    }
    
    return nil
}

// cleanupSession removes a session after timeout
func cleanupSession(sessionID string) {
    sessionsMu.Lock()
    defer sessionsMu.Unlock()
    
    session, exists := sessions[sessionID]
    if !exists {
        return
    }
    
    // Check if the session is still active
    if time.Since(session.LastActivityAt) < sessionTimeout {
        // Session is still active, schedule another cleanup
        go func() {
            time.Sleep(sessionTimeout)
            cleanupSession(sessionID)
        }()
        return
    }
    
    // Session has timed out
    log.Printf("Session %s timed out", sessionID)
    
    // Close all channels
    for _, channel := range session.Channels {
        channel.Close()
    }
    
    // Remove the session
    delete(sessions, sessionID)
}

// generateSessionID creates a unique session ID
func generateSessionID() string {
    // In a real implementation, this would generate a secure random ID
    return fmt.Sprintf("session-%d", time.Now().UnixNano())
}