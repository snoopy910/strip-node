package main

import (
	"log"
	"sync"
)

// KeyStore manages key storage in the enclave
type KeyStore struct {
	mu         sync.RWMutex
	sealedKeys map[string][]byte // Map of keyID -> sealed key share
	addresses  map[string]string // Map of keyID -> address
}

var keyStore *KeyStore

// initKeyStore initializes the key store
func initKeyStore() error {
	log.Println("Initializing key store")
	
	keyStore = &KeyStore{
		sealedKeys: make(map[string][]byte),
		addresses:  make(map[string]string),
	}
	
	// In a production implementation, this might load state from persistent storage
	// within the enclave, although Nitro Enclaves don't have persistent storage by default
	
	log.Println("Key store initialized successfully")
	return nil
}

// AddKey adds a key to the store
func (ks *KeyStore) AddKey(keyID string, sealedKey []byte, address string) {
	ks.mu.Lock()
	defer ks.mu.Unlock()
	
	ks.sealedKeys[keyID] = sealedKey
	ks.addresses[keyID] = address
	
	log.Printf("Added key %s with address %s to key store", keyID, address)
}

// GetKey retrieves a key from the store
func (ks *KeyStore) GetKey(keyID string) ([]byte, bool) {
	ks.mu.RLock()
	defer ks.mu.RUnlock()
	
	key, exists := ks.sealedKeys[keyID]
	return key, exists
}

// GetAddress retrieves an address from the store
func (ks *KeyStore) GetAddress(keyID string) (string, bool) {
	ks.mu.RLock()
	defer ks.mu.RUnlock()
	
	address, exists := ks.addresses[keyID]
	return address, exists
}

// GetAllKeys returns a list of all key IDs in the store
func (ks *KeyStore) GetAllKeys() []string {
	ks.mu.RLock()
	defer ks.mu.RUnlock()
	
	keys := make([]string, 0, len(ks.sealedKeys))
	for k := range ks.sealedKeys {
		keys = append(keys, k)
	}
	
	return keys
}

// GetStatistics returns statistics about the key store
func (ks *KeyStore) GetStatistics() map[string]interface{} {
	ks.mu.RLock()
	defer ks.mu.RUnlock()
	
	return map[string]interface{}{
		"totalKeys": len(ks.sealedKeys),
	}
}