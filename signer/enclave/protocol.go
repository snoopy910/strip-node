package main

import (
	"encoding/json"
	"fmt"
	"log"
	
	// Note: The tss-lib imports are placeholders for actual implementations
	// In a real implementation, you would use a proper TSS library with Nitro Enclave support
)

// ContributionMessage represents a message to be sent to other signers
type ContributionMessage struct {
	Type    string          `json:"type"`
	From    string          `json:"from"`
	To      []string        `json:"to,omitempty"`
	Payload json.RawMessage `json:"payload"`
}

// Note: prepareKeygenContribution was moved to tss.go to avoid duplication

// Note: prepareSigningContribution was moved to tss.go to avoid duplication

// processProtocolMessage processes messages from other signers in the TSS protocol
func processProtocolMessage(session *Session, from string, messageBytes []byte) error {
	// Parse the message
	var message ContributionMessage
	if err := json.Unmarshal(messageBytes, &message); err != nil {
		return fmt.Errorf("failed to unmarshal protocol message: %w", err)
	}
	
	log.Printf("Processing %s message from %s for session %s", message.Type, from, session.ID)
	
	// Process the message in the TSS context
	if err := handleTssProtocolMessage(session, from, messageBytes); err != nil {
		return fmt.Errorf("failed to handle TSS protocol message: %w", err)
	}
	
	// Process based on message type and session type
	switch message.Type {
	case "keygen_init", "keygen_round1", "keygen_round2", "keygen_round3":
		// Key generation message
		if session.Type != SessionTypeKeygen {
			return fmt.Errorf("received keygen message for a non-keygen session")
		}
		
		// Process the keygen message at the protocol level
		return processKeygenMessage(session, from, message.Payload)
		
	case "sign_init", "sign_round1", "sign_round2", "sign_round3":
		// Signing message
		if session.Type != SessionTypeSign {
			return fmt.Errorf("received signing message for a non-signing session")
		}
		
		// Process the signing message at the protocol level
		return processSigningMessage(session, from, message.Payload)
		
	default:
		return fmt.Errorf("unknown message type: %s", message.Type)
	}
}

// Note: completeKeygenProtocol was moved to tss.go to avoid duplication

// Note: completeSigningProtocol was moved to tss.go to avoid duplication
