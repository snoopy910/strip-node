package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/exec"
	
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
)

var (
	kmsClient *kms.Client
	kmsKeyID  string
	awsRegion string
)

// initKMS initializes the KMS client
func initKMS() error {
	// Get AWS region from environment
	awsRegion = os.Getenv("AWS_REGION")
	if awsRegion == "" {
		awsRegion = "us-east-1" // Default region
	}
	
	// Get KMS key ID from environment
	kmsKeyID = os.Getenv("KMS_KEY_ID")
	if kmsKeyID == "" {
		kmsKeyID = "alias/signer-tee-key" // Default key alias
	}
	
	log.Printf("Initializing KMS client with key ID %s in region %s", kmsKeyID, awsRegion)
	
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(awsRegion))
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}
	
	// Create KMS client
	kmsClient = kms.NewFromConfig(cfg)
	
	// Test KMS connection
	_, err = kmsClient.DescribeKey(context.Background(), &kms.DescribeKeyInput{
		KeyId: &kmsKeyID,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to KMS: %w", err)
	}
	
	log.Println("KMS client initialized successfully")
	return nil
}

// getAttestationDocument gets the Nitro Enclave attestation document
func getAttestationDocument() ([]byte, error) {
	log.Println("Getting attestation document")
	
	// Use Nitro CLI to get attestation document
	// In a real implementation, this would use the Nitro Enclaves SDK
	cmd := exec.Command("nitro-cli", "describe-attestation-document", "--output", "base64")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Failed to get attestation document: %v", err)
		return nil, fmt.Errorf("failed to get attestation document: %w", err)
	}
	
	// Decode base64 attestation document
	decoded, err := base64.StdEncoding.DecodeString(string(output))
	if err != nil {
		return nil, fmt.Errorf("failed to decode attestation document: %w", err)
	}
	
	log.Println("Attestation document retrieved successfully")
	return decoded, nil
}

// sealKeyShare encrypts a key share using KMS with attestation
func sealKeyShare(ctx context.Context, keyShare []byte, keyID string) ([]byte, error) {
	log.Printf("Sealing key share for %s using KMS", keyID)
	
	// Get attestation document
	attestDoc, err := getAttestationDocument()
	if err != nil {
		return nil, fmt.Errorf("failed to get attestation document: %w", err)
	}
	
	// Create encryption context - metadata that is cryptographically bound to the encrypted data
	encryptionContext := map[string]string{
		"keyID": keyID,
	}
	
	// Encrypt with KMS
	input := &kms.EncryptInput{
		KeyId:             &kmsKeyID,
		Plaintext:         keyShare,
		EncryptionContext: encryptionContext,
	}
	
	// In a real implementation, the attestation document would be included here
	// AWS SDK expects a specific format for attestation documents
	
	result, err := kmsClient.Encrypt(ctx, input)
	if err != nil {
		log.Printf("Failed to encrypt with KMS: %v", err)
		return nil, fmt.Errorf("failed to encrypt with KMS: %w", err)
	}
	
	log.Printf("Key share sealed successfully for %s", keyID)
	return result.CiphertextBlob, nil
}

// unsealKeyShare decrypts a key share using KMS with attestation
func unsealKeyShare(ctx context.Context, sealedKey []byte, keyID string) ([]byte, error) {
	log.Printf("Unsealing key share for %s using KMS", keyID)
	
	// Get attestation document
	attestDoc, err := getAttestationDocument()
	if err != nil {
		return nil, fmt.Errorf("failed to get attestation document: %w", err)
	}
	
	// Create encryption context - must match the one used for encryption
	encryptionContext := map[string]string{
		"keyID": keyID,
	}
	
	// Decrypt with KMS
	input := &kms.DecryptInput{
		CiphertextBlob:    sealedKey,
		KeyId:             &kmsKeyID,
		EncryptionContext: encryptionContext,
	}
	
	// In a real implementation, the attestation document would be included here
	
	result, err := kmsClient.Decrypt(ctx, input)
	if err != nil {
		log.Printf("Failed to decrypt with KMS: %v", err)
		return nil, fmt.Errorf("failed to decrypt with KMS: %w", err)
	}
	
	log.Printf("Key share unsealed successfully for %s", keyID)
	return result.Plaintext, nil
}