package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/asn1"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"time"
)

// TSS-related types - these would typically come from the TSS library
// In a real implementation, you'd import the actual library (e.g., binance-chain/tss-lib)

// PartyID represents a participant in the TSS protocol
type PartyID struct {
	ID      string
	KeyInt  *big.Int
	Moniker string
	Index   int
}

// Message represents a message in the TSS protocol
type Message struct {
	From        *PartyID
	To          []*PartyID   // nil for broadcast
	IsBroadcast bool
	Type        string
	Content     []byte
}

// ProtocolMessage represents a protocol message that will be sent between nodes
type ProtocolMessage struct {
	SessionID string   `json:"sessionId"`
	From      string   `json:"from"`
	To        []string `json:"to,omitempty"` // nil for broadcast
	Type      string   `json:"type"`
	Content   []byte   `json:"content"`
}

// TssSessionData holds all TSS-related session data
type TssSessionData struct {
	PartyID   *PartyID
	OutChan   chan Message
	ErrChan   chan error
	EndChan   interface{} // This will be cast to the appropriate type based on session type
	PublicKey *ecdsa.PublicKey
	Signature []byte      // To store signature results
}

// These are mock implementations of the actual TSS library types
// In a real implementation, you would import these from the TSS library

// LocalPartySaveData represents the output of a key generation
type LocalPartySaveData struct {
	ECDSAPub     *ecdsa.PublicKey
	ShareID      *big.Int
	SharePolynomial []byte // Simplified for mock
	ShareIndex   *big.Int
	Threshold    int
	PartyCount   int
}

// SignatureData represents the output of a signing operation
type SignatureData struct {
	R *big.Int
	S *big.Int
}

// KeyMetadata stores information about a distributed key
type KeyMetadata struct {
	KeyID         string   `json:"keyId"`
	Identity      string   `json:"identity"`
	IdentityCurve string   `json:"identityCurve"`
	KeyCurve      string   `json:"keyCurve"`
	Signers       []string `json:"signers"`
	Threshold     int      `json:"threshold"`
	CreatedAt     int64    `json:"createdAt"`
	PublicKey     *ecdsa.PublicKey `json:"publicKey,omitempty"`
}

// Initialize a TSS keygen party
func initTssKeygenParty(session *Session) (*PartyID, error) {
	// Create our party ID
	ourIndex := determineSignerIndex(enclaveID, session.Signers)
	if ourIndex == -1 {
		return nil, fmt.Errorf("enclave ID %s not found in signers list", enclaveID)
	}
	
	// Generate a key for the party ID
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate party ID key: %w", err)
	}
	
	// Create the party ID
	partyID := &PartyID{
		ID:      enclaveID,
		KeyInt:  privateKey.D,
		Moniker: fmt.Sprintf("node-%d", ourIndex),
		Index:   ourIndex,
	}
	
	return partyID, nil
}

// Initialize a TSS signing party
func initTssSigningParty(session *Session, keyData *LocalPartySaveData) (*PartyID, error) {
	// Create our party ID - similar to keygen but we need to use the same ID as in keygen
	ourIndex := determineSignerIndex(enclaveID, session.Signers)
	if ourIndex == -1 {
		return nil, fmt.Errorf("enclave ID %s not found in signers list", enclaveID)
	}
	
	// Create the party ID (in real implementation we'd reuse the one from keygen)
	partyID := &PartyID{
		ID:      enclaveID,
		KeyInt:  keyData.ShareID,
		Moniker: fmt.Sprintf("node-%d", ourIndex),
		Index:   ourIndex,
	}
	
	return partyID, nil
}

// Route TSS messages to other parties
func routeTssMessages(session *Session, excludeIDs map[string]bool) {
	if session.TssData == nil || session.TssData.OutChan == nil {
		log.Printf("Warning: Cannot route TSS messages, no TSS data for session %s", session.ID)
		return
	}
	
	// Get a connection to the parent
	conn, err := getParentConnection()
	if err != nil {
		log.Printf("Failed to get parent connection for message routing: %v", err)
		session.TssData.ErrChan <- fmt.Errorf("failed to connect to parent for message routing: %w", err)
		return
	}
	
	// Message counter for logging
	msgCount := 0
	
	// Start listening for messages on the outChan
	for msg := range session.TssData.OutChan {
		msgCount++
		
		// Prepare the protocol message
		protocolMsg := ProtocolMessage{
			SessionID: session.ID,
			From:      enclaveID,
			To:        nil, // Broadcast by default
			Type:      msg.Type,
			Content:   msg.Content,
		}
		
		// If the message has specific recipients, set them
		if !msg.IsBroadcast && msg.To != nil {
			recipients := make([]string, 0, len(msg.To))
			for _, party := range msg.To {
				// Skip excluded IDs
				if excludeIDs != nil && excludeIDs[party.ID] {
					continue
				}
				recipients = append(recipients, party.ID)
			}
			protocolMsg.To = recipients
		}
		
		// Marshal the message
		msgBytes, err := json.Marshal(protocolMsg)
		if err != nil {
			log.Printf("Failed to marshal protocol message: %v", err)
			continue
		}
		
		// Prepare the request to the parent
		req := Request{
			RequestID:  generateRequestID(),
			Operation:  "protocol_message",
			SessionID:  session.ID,
			Message:    msgBytes,
		}
		
		// Marshal the request
		reqBytes, err := json.Marshal(req)
		if err != nil {
			log.Printf("Failed to marshal request: %v", err)
			continue
		}
		
		// Send the request to the parent
		if _, err := conn.Write(reqBytes); err != nil {
			log.Printf("Failed to send message to parent: %v", err)
			continue
		}
		
		// For logging only
		if msgCount%10 == 0 {
			log.Printf("Sent %d TSS messages for session %s", msgCount, session.ID)
		}
	}
	
	log.Printf("TSS message routing complete for session %s, sent %d messages", session.ID, msgCount)
}

// handleTssProtocolMessage processes a TSS protocol message in the TSS library context
func handleTssProtocolMessage(session *Session, from string, message []byte) error {
	// Verify session exists and has TSS data
	if session == nil {
		return fmt.Errorf("session is nil")
	}
	
	if session.TssData == nil || session.TssData.OutChan == nil {
		return fmt.Errorf("no TSS data for session %s", session.ID)
	}
	
	// Parse the protocol message
	var protocolMsg ProtocolMessage
	if err := json.Unmarshal(message, &protocolMsg); err != nil {
		return fmt.Errorf("failed to unmarshal protocol message: %w", err)
	}
	
	// Determine the sender's index
	senderIndex := determineSignerIndex(from, session.Signers)
	if senderIndex == -1 {
		return fmt.Errorf("sender %s not found in signers list", from)
	}
	
	// Create a party ID for the sender
	senderPartyID := &PartyID{
		ID:      from,
		KeyInt:  nil, // We don't have this, but it's not needed for receiving
		Moniker: fmt.Sprintf("node-%d", senderIndex),
		Index:   senderIndex,
	}
	
	// Create a message object to pass to the TSS library
	// Message would be passed to the TSS library in a real implementation
	// We're logging the important information instead for now
	log.Printf("Would process TSS message: type=%s, from=%s, to=<broadcast>, content_len=%d",
		protocolMsg.Type, senderPartyID.ID, len(protocolMsg.Content))
	
	// Send the message to the TSS party
	// In a real implementation, this would be provided to the TSS library
	log.Printf("Processing protocol message of type %s from %s for session %s", 
	          protocolMsg.Type, from, session.ID)
	
	// TODO: In real implementation, update this based on the TSS library's API
	// For now, we're just sending the message to a channel
	
	// Update the session's last activity time
	session.LastActivityAt = time.Now()
	
	return nil
}

// prepareKeygenContribution prepares the initial contribution for key generation
func prepareKeygenContribution(session *Session) ([]byte, error) {
	log.Printf("Preparing keygen contribution for session %s", session.ID)
	
	// Get our index in the signers list
	ourIndex := determineSignerIndex(enclaveID, session.Signers)
	if ourIndex == -1 {
		return nil, fmt.Errorf("this enclave (%s) is not in the signers list", enclaveID)
	}
	
	// Create threshold setting - typically 2/3 of participants
	threshold := (len(session.Signers) * 2) / 3 
	if threshold < 1 {
		threshold = 1
	}
	
	// Create keygen parameters based on the key curve type
	var initialParams interface{}
	
	switch session.KeyCurve {
	case "ed25519":
		// EdDSA parameters for key generation
		initialParams = map[string]interface{}{
			"curve":      "ed25519",
			"actor_index": ourIndex,
			"threshold":  threshold,
			"party_count": len(session.Signers),
		}
		
	case "secp256k1", "bitcoin_ecdsa":
		// ECDSA parameters for key generation
		initialParams = map[string]interface{}{
			"curve":      "secp256k1",
			"actor_index": ourIndex,
			"threshold":  threshold,
			"party_count": len(session.Signers),
		}
		
	default:
		return nil, fmt.Errorf("unsupported key curve: %s", session.KeyCurve)
	}
	
	// Prepare the message data
	msgBytes, err := json.Marshal(initialParams)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal keygen parameters: %w", err)
	}
	
	// Create a contribution message to be sent to other signers
	contribution := ContributionMessage{
		Type:    "keygen_init",
		From:    enclaveID,
		To:      nil, // To all participants
		Payload: msgBytes,
	}
	
	// Serialize the contribution
	contributionBytes, err := json.Marshal(contribution)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal contribution: %w", err)
	}
	
	// Update session to include this message
	UpdateSession(session.ID, SessionStateInProgress, nil, "")
	
	return contributionBytes, nil
}

// prepareSigningContribution prepares the initial contribution for signing
func prepareSigningContribution(session *Session, keyShare []byte) ([]byte, error) {
	log.Printf("Preparing signing contribution for session %s", session.ID)
	
	// Get our index in the signers list
	ourIndex := determineSignerIndex(enclaveID, session.Signers)
	if ourIndex == -1 {
		return nil, fmt.Errorf("this enclave (%s) is not in the signers list", enclaveID)
	}
	
	// Create threshold setting
	threshold := (len(session.Signers) * 2) / 3 // Default to 2/3 threshold
	if threshold < 1 {
		threshold = 1
	}
	
	// Parse the key share data to use for signing
	var keyData map[string]interface{}
	if err := json.Unmarshal(keyShare, &keyData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal key share: %w", err)
	}
	
	// Create signing parameters based on the key curve
	var signingParams interface{}
	
	switch session.KeyCurve {
	case "ed25519":
		// Create signing parameters for EdDSA
		signingParams = map[string]interface{}{
			"curve":       "ed25519",
			"actor_index": ourIndex,
			"hash":        session.Hash,
			"threshold":   threshold,
			"party_count": len(session.Signers),
		}
		
	case "secp256k1", "bitcoin_ecdsa":
		// Create signing parameters for ECDSA
		signingParams = map[string]interface{}{
			"curve":       "secp256k1",
			"actor_index": ourIndex,
			"hash":        session.Hash,
			"threshold":   threshold,
			"party_count": len(session.Signers),
		}
		
	default:
		return nil, fmt.Errorf("unsupported key curve: %s", session.KeyCurve)
	}
	
	// Prepare the message
	msgBytes, err := json.Marshal(signingParams)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal signing parameters: %w", err)
	}
	
	// Create a contribution message to be sent to other signers
	contribution := ContributionMessage{
		Type:    "sign_init",
		From:    enclaveID,
		To:      nil, // To all participants
		Payload: msgBytes,
	}
	
	// Serialize the contribution
	contributionBytes, err := json.Marshal(contribution)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal contribution: %w", err)
	}
	
	// Update session to include this message
	UpdateSession(session.ID, SessionStateInProgress, nil, "")
	
	return contributionBytes, nil
}

// Complete the key generation protocol
func completeKeygenProtocol(session *Session, keyData *LocalPartySaveData) error {
	// Update the session state
	session.State = SessionStateCompleted
	session.LastActivityAt = time.Now()
	
	// Generate the key ID
	keyID := fmt.Sprintf("%s_%s_%s", session.Identity, session.IdentityCurve, session.KeyCurve)
	
	// Marshal the key data
	keyBytes, err := json.Marshal(keyData)
	if err != nil {
		return fmt.Errorf("failed to marshal key data: %w", err)
	}
	
	// Create key metadata
	keyMeta := &KeyMetadata{
		KeyID:         keyID,
		Identity:      session.Identity,
		IdentityCurve: session.IdentityCurve,
		KeyCurve:      session.KeyCurve,
		Signers:       session.Signers,
		Threshold:     keyData.Threshold,
		CreatedAt:     time.Now().Unix(),
		PublicKey:     keyData.ECDSAPub,
	}
	
	// Store the key in Nitro enclave memory
	if err := storeKeyShare(keyID, keyBytes, keyMeta); err != nil {
		return fmt.Errorf("failed to store key share: %w", err)
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
	
	return nil
}

// Complete the signing protocol
func completeSigningProtocol(session *Session, sigData *SignatureData) error {
	// Convert the signature to DER format
	signature, err := createDERSignature(sigData.R, sigData.S, session.KeyCurve)
	if err != nil {
		return fmt.Errorf("failed to create DER signature: %w", err)
	}
	
	// Update the session
	session.State = SessionStateCompleted
	session.Signature = signature
	session.LastActivityAt = time.Now()
	
	return nil
}

// Create a DER-encoded signature from R and S values
func createDERSignature(r, s *big.Int, curve string) ([]byte, error) {
	// ECDSA signature is encoded as ASN.1 DER format for most applications
	type ecdsaSignature struct {
		R, S *big.Int
	}
	
	// Create the signature structure
	sig := ecdsaSignature{
		R: r,
		S: s,
	}
	
	// Marshal to DER format
	return asn1.Marshal(sig)
}

// Connect to the parent process
func getParentConnection() (net.Conn, error) {
	// Get the enclave port from environment
	port := getEnvOrDefault("ENCLAVE_PORT", "8000")
	
	// Connect to the parent
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%s", port))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to parent: %w", err)
	}
	
	return conn, nil
}

// Sync key data to parent process for persistence
func syncKeyStoreToParent(conn net.Conn, keyID string) error {
	// Get the key metadata
	keyMeta, err := getKeyMetadata(keyID)
	if err != nil {
		return fmt.Errorf("failed to get key metadata: %w", err)
	}
	
	// Create a sync request
	req := Request{
		RequestID:  generateRequestID(),
		Operation:  "sync_key",
		KeyID:      keyID,
		MetaData:   keyMeta,
	}
	
	// Marshal the request
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal sync request: %w", err)
	}
	
	// Send the request
	if _, err := conn.Write(reqBytes); err != nil {
		return fmt.Errorf("failed to send sync request: %w", err)
	}
	
	// Wait for response with a timeout
	respCh := make(chan []byte, 1)
	errCh := make(chan error, 1)
	
	go func() {
		buffer := make([]byte, 4096)
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
			return fmt.Errorf("failed to unmarshal sync response: %w", err)
		}
		
		if resp.Error != "" {
			return fmt.Errorf("parent returned error for key sync: %s", resp.Error)
		}
		
		log.Printf("Successfully synced key %s to parent", keyID)
		return nil
		
	case err := <-errCh:
		return fmt.Errorf("error reading sync response: %w", err)
		
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout waiting for sync response")
	}
}

// Helper function to get environment variable with default
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
