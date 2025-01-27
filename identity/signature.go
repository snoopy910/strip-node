package identity

import (
	"bytes"
	"encoding/json"

	"github.com/StripChain/strip-node/sequencer"
	// Add this import
)

type OperationForSigning struct {
	SerializedTxn  string `json:"serializedTxn"`
	DataToSign     string `json:"dataToSign"`
	ChainId        string `json:"chainId"`
	KeyCurve       string `json:"keyCurve"`
	Type           string `json:"type"`
	Solver         string `json:"solver"`
	SolverMetadata string `json:"solverMetadata"`
}

type IntentForSigning struct {
	Operations    []OperationForSigning `json:"operations"`
	Identity      string                `json:"identity"`
	IdentityCurve string                `json:"identityCurve"`
	Expiry        uint64                `json:"expiry"`
}

func SanitiseIntent(intent sequencer.Intent) (string, error) {
	intentForSigning := IntentForSigning{
		Identity:      intent.Identity,
		IdentityCurve: intent.IdentityCurve,
		Expiry:        intent.Expiry,
	}

	for _, operation := range intent.Operations {
		operationForSigning := OperationForSigning{
			SerializedTxn:  operation.SerializedTxn,
			DataToSign:     operation.DataToSign,
			ChainId:        operation.ChainId,
			KeyCurve:       operation.KeyCurve,
			Type:           operation.Type,
			Solver:         operation.Solver,
			SolverMetadata: operation.SolverMetadata,
		}

		intentForSigning.Operations = append(intentForSigning.Operations, operationForSigning)
	}

	jsonBytes, err := json.Marshal(intentForSigning)
	if err != nil {
		return "", err
	}

	dst := &bytes.Buffer{}
	if err := json.Compact(dst, jsonBytes); err != nil {
		panic(err)
	}

	return dst.String(), nil
}
