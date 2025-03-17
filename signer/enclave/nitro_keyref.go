package main

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"sync"
)

// KeyReferenceStore manages references to sealed keys for AWS Nitro Enclaves
type KeyReferenceStore struct {
	references map[string][]byte // Maps key IDs to encrypted key references
	mutex      sync.RWMutex
}

var (
	// Global key reference store
	keyRefStore = &KeyReferenceStore{
		references: make(map[string][]byte),
	}
)

// getKeyReference is implemented below with additional nitroKeyStore support

// storeKeyReferenceSimple stores a key reference without address tracking
// This is the simple version without address tracking used by non-keygen components
func storeKeyReferenceSimple(keyID string, keyRef []byte) {
	keyRefStore.mutex.Lock()
	defer keyRefStore.mutex.Unlock()
	
	keyRefStore.references[keyID] = keyRef
	log.Printf("Stored key reference for %s", keyID)
}

// removeKeyReference removes a key reference
func removeKeyReference(keyID string) {
	keyRefStore.mutex.Lock()
	defer keyRefStore.mutex.Unlock()
	
	delete(keyRefStore.references, keyID)
	log.Printf("Removed key reference for %s", keyID)
}

// generateKeyID creates a unique key ID based on identity and curves
func generateKeyID(identity, identityCurve, keyCurve string) string {
	return identity + "_" + identityCurve + "_" + keyCurve
}

// generateSealedKeyID creates a hash-based sealed key ID for a particular key
// This allows us to have multiple sealed versions of the same key (for different enclaves)
func generateSealedKeyID(keyID string, enclaveID string) string {
	// Create a hash of the key ID and enclave ID
	hasher := sha256.New()
	hasher.Write([]byte(keyID))
	hasher.Write([]byte(enclaveID))
	hash := hasher.Sum(nil)
	
	return hex.EncodeToString(hash)
}

// getKeyReference checks if a key exists and returns its sealed form
func getKeyReference(keyID string) ([]byte, bool) {
	keyRefStore.mutex.RLock()
	defer keyRefStore.mutex.RUnlock()
	
	keyRef, exists := keyRefStore.references[keyID]
	if !exists {
		// Check if we have it in the nitro key store
		keyBytes, keyExists := nitroKeyStore[keyID]
		if keyExists {
			return keyBytes, true
		}
		return nil, false
	}
	
	return keyRef, true
}
