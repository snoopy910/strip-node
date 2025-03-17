package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

// KeyMetadata stores metadata about a distributed key
type KeyMetadata struct {
	KeyID         string   `json:"keyId"`
	Identity      string   `json:"identity"`
	IdentityCurve string   `json:"identityCurve"`
	KeyCurve      string   `json:"keyCurve"`
	Signers       []string `json:"signers"`
	Threshold     int      `json:"threshold"`
	CreatedAt     int64    `json:"createdAt"`
	PublicKey     interface{} `json:"publicKey,omitempty"`
}

var (
	// In-memory storage for AWS Nitro Enclaves
	// Since Nitro Enclaves don't persist data across restarts,
	// we need to store keys in memory and rely on parent process for persistence
	nitroKeyStore     = make(map[string][]byte)
	nitroKeyMetadata  = make(map[string]*KeyMetadata)
	nitroSessionKeys  = make(map[string][]byte) // Session ID -> Key bytes
	nitroKeystoreMutex sync.RWMutex
)

// storeNitroKeyInMemory stores a key in the Nitro enclave's memory
func storeNitroKeyInMemory(sessionID string, keyBytes []byte) {
	nitroKeystoreMutex.Lock()
	defer nitroKeystoreMutex.Unlock()
	
	// Store the key in the session-specific storage
	nitroSessionKeys[sessionID] = keyBytes
	log.Printf("Stored key in memory for session %s", sessionID)
}

// getNitroSessionKey retrieves a key for a specific session
func getNitroSessionKey(sessionID string) ([]byte, bool) {
	nitroKeystoreMutex.RLock()
	defer nitroKeystoreMutex.RUnlock()
	
	key, exists := nitroSessionKeys[sessionID]
	return key, exists
}

// storeKeyShare stores a key share in the Nitro enclave's memory
func storeKeyShare(keyID string, keyBytes []byte, metadata *KeyMetadata) error {
	nitroKeystoreMutex.Lock()
	defer nitroKeystoreMutex.Unlock()
	
	// Store the key bytes
	nitroKeyStore[keyID] = keyBytes
	
	// Store the metadata
	nitroKeyMetadata[keyID] = metadata
	
	log.Printf("Stored key share for %s in Nitro memory storage", keyID)
	return nil
}

// retrieveKeyShare retrieves a key share from the Nitro enclave's memory
func retrieveKeyShare(keyID string) ([]byte, error) {
	nitroKeystoreMutex.RLock()
	defer nitroKeystoreMutex.RUnlock()
	
	keyBytes, exists := nitroKeyStore[keyID]
	if !exists {
		return nil, fmt.Errorf("key %s not found in Nitro memory storage", keyID)
	}
	
	return keyBytes, nil
}

// getKeyMetadata retrieves metadata for a key
func getKeyMetadata(keyID string) (*KeyMetadata, error) {
	nitroKeystoreMutex.RLock()
	defer nitroKeystoreMutex.RUnlock()
	
	metadata, exists := nitroKeyMetadata[keyID]
	if !exists {
		return nil, fmt.Errorf("metadata for key %s not found", keyID)
	}
	
	return metadata, nil
}

// saveKeyMetadata saves metadata for a key
func saveKeyMetadata(metadata *KeyMetadata) error {
	if metadata == nil || metadata.KeyID == "" {
		return fmt.Errorf("invalid key metadata")
	}
	
	nitroKeystoreMutex.Lock()
	defer nitroKeystoreMutex.Unlock()
	
	nitroKeyMetadata[metadata.KeyID] = metadata
	return nil
}

// deleteKey removes a key from memory
// This should be used when migrating keys out of the enclave
func deleteKey(keyID string) {
	nitroKeystoreMutex.Lock()
	defer nitroKeystoreMutex.Unlock()
	
	delete(nitroKeyStore, keyID)
	delete(nitroKeyMetadata, keyID)
	
	log.Printf("Deleted key %s from Nitro memory storage", keyID)
}

// syncKeyStoreToParent notifies the parent process that it should
// persist the key store (e.g., after key generation or key updates)
// This is a crucial function for AWS Nitro since memory is not persistent
func syncKeyStoreToParent(conn net.Conn, keyID string) error {
	// Get the key data
	nitroKeystoreMutex.RLock()
	keyBytes, keyExists := nitroKeyStore[keyID]
	metadata, metadataExists := nitroKeyMetadata[keyID]
	nitroKeystoreMutex.RUnlock()
	
	if !keyExists || !metadataExists {
		return fmt.Errorf("key %s not found for sync", keyID)
	}
	
	// Prepare the sync data
	syncData := struct {
		KeyID     string       `json:"keyId"`
		KeyBytes  []byte       `json:"keyBytes"`
		Metadata  *KeyMetadata `json:"metadata"`
	}{
		KeyID:     keyID,
		KeyBytes:  keyBytes,
		Metadata:  metadata,
	}
	
	// Serialize to JSON
	syncDataBytes, err := json.Marshal(syncData)
	if err != nil {
		return fmt.Errorf("failed to marshal key sync data: %w", err)
	}
	
	// Create sync request for parent
	syncRequest := Request{
		RequestID: generateRequestID(),
		Operation: "sync_key",
		Message:   syncDataBytes,
	}
	
	// Send to parent
	requestBytes, err := json.Marshal(syncRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal sync request: %w", err)
	}
	
	if _, err := conn.Write(requestBytes); err != nil {
		return fmt.Errorf("failed to send key sync request: %w", err)
	}
	
	log.Printf("Sent key sync request to parent for key %s", keyID)
	return nil
}

// loadKeyShare loads a key share sent from the parent process
func loadKeyShare(keyID string, keyBytes []byte, metadata *KeyMetadata) error {
	if keyBytes == nil || len(keyBytes) == 0 {
		return fmt.Errorf("invalid key bytes")
	}
	
	if metadata == nil || metadata.KeyID == "" {
		return fmt.Errorf("invalid key metadata")
	}
	
	return storeKeyShare(keyID, keyBytes, metadata)
}

// generateRequestID creates a unique request ID
func generateRequestID() uint64 {
	// Simple implementation - in a real system this would be more sophisticated
	return uint64(time.Now().UnixNano())
}
