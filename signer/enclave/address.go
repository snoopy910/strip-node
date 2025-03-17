package main

// getAddress retrieves a stored address
func getAddress(keyID string) (string, bool) {
	keyStore.mu.RLock()
	defer keyStore.mu.RUnlock()
	
	address, exists := keyStore.addresses[keyID]
	return address, exists
}