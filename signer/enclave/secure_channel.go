package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/ecdsa"
    "crypto/elliptic"
    "crypto/rand"
    "crypto/sha256"
    "crypto/tls"
    "crypto/x509"
    "encoding/binary"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "math/big"
    "net"
    "time"
)

// SecureChannel represents an established secure channel
type SecureChannel struct {
    conn          net.Conn
    sharedSecret  []byte
    encryptCipher cipher.AEAD
    decryptCipher cipher.AEAD
    peerID        string
    attestation   *AttestationDocument
}

// EstablishSecureChannel establishes a secure channel with another enclave
func EstablishSecureChannel(conn net.Conn, peerID string) (*SecureChannel, error) {
    log.Printf("Establishing secure channel with peer %s", peerID)
    
    // 1. Generate ephemeral ECDH key pair
    privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    if err != nil {
        return nil, fmt.Errorf("failed to generate ECDH key: %w", err)
    }
    
    // 2. Get our attestation document
    attestation, err := GetAttestationDocument()
    if err != nil {
        return nil, fmt.Errorf("failed to get attestation document: %w", err)
    }
    
    // 3. Send our public key and attestation
    handshake := struct {
        PublicKey   []byte               `json:"publicKey"`
        Attestation *AttestationDocument `json:"attestation"`
    }{
        PublicKey:   elliptic.Marshal(elliptic.P256(), privateKey.X, privateKey.Y),
        Attestation: attestation,
    }
    
    handshakeBytes, err := json.Marshal(handshake)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal handshake: %w", err)
    }
    
    // Send handshake size and data
    sizeBytes := make([]byte, 4)
    binary.BigEndian.PutUint32(sizeBytes, uint32(len(handshakeBytes)))
    if _, err := conn.Write(sizeBytes); err != nil {
        return nil, fmt.Errorf("failed to send handshake size: %w", err)
    }
    
    if _, err := conn.Write(handshakeBytes); err != nil {
        return nil, fmt.Errorf("failed to send handshake: %w", err)
    }
    
    // 4. Receive peer's public key and attestation
    sizeBytes = make([]byte, 4)
    if _, err := io.ReadFull(conn, sizeBytes); err != nil {
        return nil, fmt.Errorf("failed to read peer handshake size: %w", err)
    }
    
    size := binary.BigEndian.Uint32(sizeBytes)
    if size > 1024*1024 { // 1MB limit
        return nil, fmt.Errorf("peer handshake too large: %d bytes", size)
    }
    
    peerHandshakeBytes := make([]byte, size)
    if _, err := io.ReadFull(conn, peerHandshakeBytes); err != nil {
        return nil, fmt.Errorf("failed to read peer handshake: %w", err)
    }
    
    var peerHandshake struct {
        PublicKey   []byte               `json:"publicKey"`
        Attestation *AttestationDocument `json:"attestation"`
    }
    
    if err := json.Unmarshal(peerHandshakeBytes, &peerHandshake); err != nil {
        return nil, fmt.Errorf("failed to unmarshal peer handshake: %w", err)
    }
    
    // 5. Verify peer's attestation
    valid, err := VerifyAttestationDocument(peerHandshake.Attestation)
    if err != nil {
        return nil, fmt.Errorf("failed to verify peer attestation: %w", err)
    }
    
    if !valid {
        return nil, fmt.Errorf("peer attestation verification failed")
    }
    
    // 6. Derive the shared secret
    x, y := elliptic.Unmarshal(elliptic.P256(), peerHandshake.PublicKey)
    if x == nil {
        return nil, fmt.Errorf("invalid peer public key")
    }
    
    sharedX, _ := privateKey.ScalarMult(x, y, privateKey.D.Bytes())
    sharedSecret := sha256.Sum256(sharedX.Bytes())
    
    // 7. Set up encryption and decryption
    encryptKey := sha256.Sum256(append(sharedSecret[:], []byte("encrypt")...))
    decryptKey := sha256.Sum256(append(sharedSecret[:], []byte("decrypt")...))
    
    encryptBlock, err := aes.NewCipher(encryptKey[:])
    if err != nil {
        return nil, fmt.Errorf("failed to create encrypt cipher: %w", err)
    }
    
    decryptBlock, err := aes.NewCipher(decryptKey[:])
    if err != nil {
        return nil, fmt.Errorf("failed to create decrypt cipher: %w", err)
    }
    
    encryptCipher, err := cipher.NewGCM(encryptBlock)
    if err != nil {
        return nil, fmt.Errorf("failed to create encrypt GCM: %w", err)
    }
    
    decryptCipher, err := cipher.NewGCM(decryptBlock)
    if err != nil {
        return nil, fmt.Errorf("failed to create decrypt GCM: %w", err)
    }
    
    log.Printf("Secure channel established with peer %s", peerID)
    
    return &SecureChannel{
        conn:          conn,
        sharedSecret:  sharedSecret[:],
        encryptCipher: encryptCipher,
        decryptCipher: decryptCipher,
        peerID:        peerID,
        attestation:   peerHandshake.Attestation,
    }, nil
}

// SendMessage sends an encrypted message over the secure channel
func (c *SecureChannel) SendMessage(message []byte) error {
    // Create a nonce
    nonce := make([]byte, c.encryptCipher.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return fmt.Errorf("failed to generate nonce: %w", err)
    }
    
    // Encrypt the message
    ciphertext := c.encryptCipher.Seal(nil, nonce, message, nil)
    
    // Prepare the message: nonce + ciphertext
    fullMessage := append(nonce, ciphertext...)
    
    // Send the message size
    sizeBytes := make([]byte, 4)
    binary.BigEndian.PutUint32(sizeBytes, uint32(len(fullMessage)))
    if _, err := c.conn.Write(sizeBytes); err != nil {
        return fmt.Errorf("failed to send message size: %w", err)
    }
    
    // Send the message
    if _, err := c.conn.Write(fullMessage); err != nil {
        return fmt.Errorf("failed to send message: %w", err)
    }
    
    return nil
}

// ReceiveMessage receives and decrypts a message from the secure channel
func (c *SecureChannel) ReceiveMessage() ([]byte, error) {
    // Read message size
    sizeBytes := make([]byte, 4)
    if _, err := io.ReadFull(c.conn, sizeBytes); err != nil {
        return nil, fmt.Errorf("failed to read message size: %w", err)
    }
    
    size := binary.BigEndian.Uint32(sizeBytes)
    if size > 1024*1024 { // 1MB limit
        return nil, fmt.Errorf("message too large: %d bytes", size)
    }
    
    // Read full message
    fullMessage := make([]byte, size)
    if _, err := io.ReadFull(c.conn, fullMessage); err != nil {
        return nil, fmt.Errorf("failed to read message: %w", err)
    }
    
    // Extract nonce and ciphertext
    nonceSize := c.decryptCipher.NonceSize()
    if len(fullMessage) < nonceSize {
        return nil, fmt.Errorf("message too short")
    }
    
    nonce := fullMessage[:nonceSize]
    ciphertext := fullMessage[nonceSize:]
    
    // Decrypt the message
    plaintext, err := c.decryptCipher.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt message: %w", err)
    }
    
    return plaintext, nil
}

// Close closes the secure channel
func (c *SecureChannel) Close() error {
    return c.conn.Close()
}

// PeerID returns the ID of the peer
func (c *SecureChannel) PeerID() string {
    return c.peerID
}

// PeerAttestation returns the attestation document of the peer
func (c *SecureChannel) PeerAttestation() *AttestationDocument {
    return c.attestation
}