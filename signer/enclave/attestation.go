package main

import (
    "crypto/ecdsa"
    "crypto/elliptic"
    "crypto/rand"
    "crypto/ed25519"
    "crypto/x509"
    "encoding/json"
    "fmt"
    "log"
    "os/exec"
    "time"
)

// AttestationDocument represents a TEE attestation proof
type AttestationDocument struct {
    Data        []byte            `json:"data"`
    PCRs        map[int][]byte    `json:"pcrs"`
    TEEType     string            `json:"teeType"`
    Timestamp   time.Time         `json:"timestamp"`
    Signature   []byte            `json:"signature"`
    Certificate []byte            `json:"certificate"`
    EnclaveID   string            `json:"enclaveId"`
    PublicKey   []byte            `json:"publicKey"`
}

// NitroAttestationDocument represents the AWS Nitro attestation format
type NitroAttestationDocument struct {
    ModuleID    string            `json:"moduleId"`
    Timestamp   string            `json:"timestamp"`
    PCRs        map[string]string `json:"pcrs"`
    Certificate string            `json:"certificate"`
    CabundleURL string            `json:"cabundleUrl"`
    PublicKey   string            `json:"publicKey"`
    UserData    string            `json:"userData"`
    Nonce       string            `json:"nonce"`
    Signature   string            `json:"signature"`
}

var (
    // Cached attestation document
    cachedAttestation     *AttestationDocument
    cachedAttestationTime time.Time
    
    // Cached trusted certificates for verification
    trustedCerts []*x509.Certificate
    
    // Cached attestation keypair
    cachedKeyPair        interface{}
)

// GetAttestationDocument retrieves an attestation document for this enclave
func GetAttestationDocument() (*AttestationDocument, error) {
    // Check if we have a recent cached attestation (less than 1 minute old)
    now := time.Now()
    if cachedAttestation != nil && now.Sub(cachedAttestationTime) < time.Minute {
        return cachedAttestation, nil
    }
    
    // Generate a fresh attestation document
    log.Println("Generating fresh attestation document")
    
    // Use nitro-cli to generate an attestation (in AWS Nitro Enclaves)
    cmd := exec.Command("nitro-cli", "describe-attestation-document", "--output", "json")
    output, err := cmd.Output()
    if err != nil {
        return nil, fmt.Errorf("failed to generate attestation document: %w", err)
    }
    
    // Parse the Nitro attestation document
    var nitroDoc NitroAttestationDocument
    if err := json.Unmarshal(output, &nitroDoc); err != nil {
        return nil, fmt.Errorf("failed to parse attestation document: %w", err)
    }
    
    // Convert to our standard attestation format
    doc := &AttestationDocument{
        Data:        []byte(nitroDoc.UserData),
        PCRs:        make(map[int][]byte),
        TEEType:     "aws-nitro",
        Timestamp:   time.Now(), // We'll parse the actual timestamp later
        Signature:   []byte(nitroDoc.Signature),
        Certificate: []byte(nitroDoc.Certificate),
        EnclaveID:   enclaveID,
        PublicKey:   []byte(nitroDoc.PublicKey),
    }
    
    // Parse PCR values
    for pcr, value := range nitroDoc.PCRs {
        var pcrIndex int
        fmt.Sscanf(pcr, "PCR%d", &pcrIndex)
        doc.PCRs[pcrIndex] = []byte(value)
    }
    
    // Parse timestamp
    if t, err := time.Parse(time.RFC3339, nitroDoc.Timestamp); err == nil {
        doc.Timestamp = t
    }
    
    // Cache the document
    cachedAttestation = doc
    cachedAttestationTime = now
    
    return doc, nil
}

// VerifyAttestationDocument verifies the authenticity of an attestation document
func VerifyAttestationDocument(doc *AttestationDocument) (bool, error) {
    // First, ensure we have trusted certificates
    if err := ensureTrustedCerts(); err != nil {
        return false, fmt.Errorf("failed to load trusted certificates: %w", err)
    }
    
    // Verify based on TEE type
    switch doc.TEEType {
    case "aws-nitro":
        return verifyNitroAttestation(doc)
    case "mock":
        // For mock attestations, always return true
        return true, nil
    default:
        return false, fmt.Errorf("unsupported TEE type: %s", doc.TEEType)
    }
}

// verifyNitroAttestation verifies an AWS Nitro attestation document
func verifyNitroAttestation(doc *AttestationDocument) (bool, error) {
    // In a real implementation, we would:
    // 1. Parse the certificate chain
    // 2. Verify the certificate against trusted roots
    // 3. Extract the public key from the certificate
    // 4. Verify the signature on the attestation data
    // 5. Check the PCR values against expected values
    
    // For now, this is a placeholder implementation
    log.Println("Verifying Nitro attestation document")
    
    // Trusted PCR values for our enclave
    // These would be configured based on the actual enclave
    trustedPCRs := map[int][]byte{
        // Example PCR values
        0: []byte("pcr0-value"),
        1: []byte("pcr1-value"),
        // ...
    }
    
    // Check PCR values (simplified)
    for pcr, expectedValue := range trustedPCRs {
        actualValue, exists := doc.PCRs[pcr]
        if !exists {
            return false, fmt.Errorf("missing PCR%d in attestation", pcr)
        }
        
        // In reality, we'd use a constant-time comparison
        if string(actualValue) != string(expectedValue) {
            return false, fmt.Errorf("PCR%d value mismatch", pcr)
        }
    }
    
    return true, nil
}

// ensureTrustedCerts loads the trusted certificates if not already loaded
func ensureTrustedCerts() error {
    if len(trustedCerts) > 0 {
        return nil
    }
    
    // In a real implementation, we would load trusted certificates from a 
    // well-known location like the AWS Nitro Attestation CA bundle
    
    // For now, this is a placeholder implementation
    log.Println("Loading trusted certificates")
    
    // We would download the CA bundle from the provided URL and parse the certificates
    // For example:
    /*
    resp, err := http.Get("https://aws-nitro-enclaves.amazonaws.com/ca-bundle")
    if err != nil {
        return fmt.Errorf("failed to download CA bundle: %w", err)
    }
    defer resp.Body.Close()
    
    caData, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("failed to read CA bundle: %w", err)
    }
    
    trustedCerts, err = x509.ParseCertificates(caData)
    if err != nil {
        return fmt.Errorf("failed to parse CA bundle: %w", err)
    }
    */
    
    return nil
}

// Embed attestation in a KMS request
func attachAttestationToKMSRequest(keyID string, encryptionContext map[string]string) (map[string]string, error) {
    // Get attestation document
    doc, err := GetAttestationDocument()
    if err != nil {
        return nil, fmt.Errorf("failed to get attestation document: %w", err)
    }
    
    // Serialize attestation document
    attestationBytes, err := json.Marshal(doc)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal attestation document: %w", err)
    }
    
    // Add attestation to encryption context
    newContext := make(map[string]string)
    for k, v := range encryptionContext {
        newContext[k] = v
    }
    
    // In a real implementation, we would format the attestation according to KMS requirements
    // This is a placeholder to illustrate the concept
    newContext["attestation"] = string(attestationBytes)
    
    return newContext, nil
}

// getAttestationReport gets the current attestation report for the enclave
func getAttestationReport() (*AttestationDocument, error) {
    // Use the existing GetAttestationDocument function
    return GetAttestationDocument()
}

// getAttestationKeyPair gets or generates the attestation key pair
func getAttestationKeyPair() (interface{}, error) {
    // If we already have a cached key pair, return it
    if cachedKeyPair != nil {
        return cachedKeyPair, nil
    }
    
    // In a real SGX implementation, this would either:
    // 1. Generate a new key pair within the enclave
    // 2. Retrieve a key pair from sealed storage
    // 3. Use a platform-provided key (like the SGX sealing key)
    
    // For this implementation, we'll generate an ECDSA key pair
    // This is compatible with the tss-lib we're using for the threshold signature scheme
    privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    if err != nil {
        return nil, fmt.Errorf("failed to generate attestation key pair: %w", err)
    }
    
    // Cache the key pair for future use
    cachedKeyPair = privateKey
    
    return privateKey, nil
}

// extractPublicKeyBytes extracts the raw bytes of a public key from a key pair
func extractPublicKeyBytes(keyPair interface{}) ([]byte, error) {
    switch k := keyPair.(type) {
    case *ecdsa.PrivateKey:
        // For ECDSA keys, we'll use the standard SEC1 encoding (compressed)
        return elliptic.MarshalCompressed(k.PublicKey.Curve, k.PublicKey.X, k.PublicKey.Y), nil
        
    case *ed25519.PrivateKey:
        // For Ed25519 keys, we can directly use the public key bytes
        return k.Public().(ed25519.PublicKey), nil
        
    default:
        return nil, fmt.Errorf("unsupported key type: %T", keyPair)
    }
}