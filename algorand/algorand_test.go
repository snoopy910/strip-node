package algorand

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"testing"

	"github.com/StripChain/strip-node/common"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/indexer"
	"github.com/algorand/go-algorand-sdk/types"
	"github.com/test-go/testify/assert"
	"github.com/test-go/testify/mock"
)

func TestConnectToTestnetAlgorandChain(t *testing.T) {
	chain, err := common.GetChain("SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI=")
	if err != nil {
		t.Fatal(err)
	}

	// Create an algod client (no API key needed for AlgoNode/ Nodely)
	client, err := algod.MakeClient(chain.ChainUrl, "")
	if err != nil {
		t.Fatal(err)
	}
	respStatus, err := client.Status().Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Client status: ", respStatus)
	assert.NotNil(t, respStatus)
	responseHealth := client.HealthCheck()
	fmt.Println("Client health check: ", responseHealth)
	assert.NotNil(t, responseHealth)
	// check indexer
	indexerClient, err := indexer.MakeClient(chain.IndexerUrl, "")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Indexer status: ", indexerClient)
	assert.NotNil(t, indexerClient)
}

func TestConnectToMainnetAlgorandChain(t *testing.T) {
	chain, err := common.GetChain("wGHE2Pwdvd7S12BL5FaOP20EGYesN73ktiC1qzkkit8=")
	if err != nil {
		t.Fatal(err)
	}

	// Create an algod client (no API key needed for AlgoNode/ Nodely)
	client, err := algod.MakeClient(chain.ChainUrl, "")
	if err != nil {
		t.Fatal(err)
	}
	respStatus, err := client.Status().Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Client status: ", respStatus)
	assert.NotNil(t, respStatus)
	responseHealth := client.HealthCheck()
	fmt.Println("Client health check: ", responseHealth)
	assert.NotNil(t, responseHealth)
	// check indexer
	indexerClient, err := indexer.MakeClient(chain.IndexerUrl, "")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Indexer status: ", indexerClient)
	assert.NotNil(t, indexerClient)

}

func TestSendAlgorandTransaction(t *testing.T) {
	// Testnet

	tests := []struct {
		Name            string
		SerializedTxn   string
		GenesisHash     string
		SignatureBase64 string
		MockClient      bool
		IsError         bool
		Error           string
		TxId            string
	}{
		{
			Name:            "Valid transaction",
			SerializedTxn:   "RGRWC3LUMSRWMZLFZYAAHGPAUJTHMZFDM5SW5LDUMVZXI3TFOQWXMMJOGCRGO2GEEBEGHNIYUSZ4QTWICDZC2TYQQHFQ64PQLGT2YIG6YYXX64HFBE5CFITMO3GMRI3SMN3MIICAJDHOGBAHNB3EKVOZ34L6AOSCOPB2CCHMVXBAYVPLSJJYKSCILGRXG3TEYQQEASGO4MCAO2DWIVK5TXYX4A5EE46DUEEOZLOCBRK6XESTQVEEQWNEOR4XAZNDOBQXS===",
			GenesisHash:     "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI=",
			SignatureBase64: "KA/XmHsJjQipAXl4gqKwZtrb1WmfvClUQkHxzyzrmL1cPEPWEVMwkzF1re9E7+1iQdp0BQLEMgsAA8QM2p65Dg==",
			IsError:         false,
			MockClient:      true,
			TxId:            "P5RRK4SQHBC5QZUCKNMCQ6MDMAJSEPCQKEZJCPRV4R3MX3YWICUQ",
		},
		{
			Name:            "Invalid request",
			SerializedTxn:   "RGRWC3LUMSRWMZLFZYAAHGPAUJTHMZFDM5SW5LDUMVZXI3TFOQWXMMJOGCRGO2GEEBEGHNIYUSZ4QTWICDZC2TYQQHFQ64PQLGT2YIG6YYXX64HFBE5CFITMO3GMRI3SMN3MIICAJDHOGBAHNB3EKVOZ34L6AOSCOPB2CCHMVXBAYVPLSJJYKSCILGRXG3TEYQQEASGO4MCAO2DWIVK5TXYX4A5EE46DUEEOZLOCBRK6XESTQVEEQWNEOR4XAZNDOBQXS===",
			GenesisHash:     "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI=",
			SignatureBase64: "KA/XmHsJjQipAXl4gqKwZtrb1WmfvClUQkHxzyzrmL1cPEPWEVMwkzF1re9E7+1iQdp0BQLEMgsAA8QM2p65Dg==",
			IsError:         true,
			MockClient:      true,
			Error:           "failed to send transaction: error in send raw transaction",
		},
		{
			Name:            "Invalid serialized transaction",
			SerializedTxn:   "0xRGRWC3L",
			GenesisHash:     "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI=",
			SignatureBase64: "2I3C36T4ZTOJ4NSRD3HLXJSD4G3WBP64HTPXMN73AJN7RF2HFWAQ",
			IsError:         true,
			Error:           "failed to decode serialized transaction: illegal base32 data at input byte 0",
		},
		{
			Name:            "Invalid signature",
			SerializedTxn:   "RGRWC3LUMSRWMZLFZYAAHGPAUJTHMZFDM5SW5LDUMVZXI3TFOQWXMMJOGCRGO2GEEBEGHNIYUSZ4QTWICDZC2TYQQHFQ64PQLGT2YIG6YYXX64HFBE5CFITMO3GMRI3SMN3MIICAJDHOGBAHNB3EKVOZ34L6AOSCOPB2CCHMVXBAYVPLSJJYKSCILGRXG3TEYQQEASGO4MCAO2DWIVK5TXYX4A5EE46DUEEOZLOCBRK6XESTQVEEQWNEOR4XAZNDOBQXS===",
			GenesisHash:     "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI=",
			SignatureBase64: "2I3C36T4ZTOJ4NSRD3HLXJS",
			IsError:         true,
			Error:           "failed to decode signature: illegal base64 data at input byte 20",
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			clients := &MockClients{mockAlgod: new(MockAlgodClient), mockIndexer: new(MockIndexerClient)}
			if tt.MockClient && !tt.IsError {
				mockSendRawTransactionRequester := new(MockSendRawTransactionRequester)
				clients.mockAlgod.On("SendRawTransaction", mock.Anything).Return(mockSendRawTransactionRequester)
				mockSendRawTransactionRequester.On("Do", mock.Anything).Return(tt.TxId, nil)
			}
			if tt.MockClient && tt.IsError {
				mockSendRawTransactionRequester := new(MockSendRawTransactionRequester)
				clients.mockAlgod.On("SendRawTransaction", mock.Anything).Return(mockSendRawTransactionRequester)
				mockSendRawTransactionRequester.On("Do", mock.Anything).Return("", fmt.Errorf("error in send raw transaction"))
			}
			_, err := clients.SendAlgorandTransaction(tt.SerializedTxn, tt.GenesisHash, tt.SignatureBase64)
			if !tt.IsError && err != nil {
				t.Errorf("WithdrawAlgorandTokenGetSignature() error = %v", err)
			} else if tt.IsError && err.Error() != tt.Error {
				t.Errorf("expected error: %v", err)
			}

		})
	}

}

func TestWithdrawAlgorandNativeGetSignature(t *testing.T) {
	// Testnet

	genesisHashStr := "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI="

	// Decode the Base64 string to get the raw bytes.
	genesisHashBytes, err := base64.StdEncoding.DecodeString(genesisHashStr)
	if err != nil {
		log.Fatalf("failed to decode genesis hash: %v", err)
	}

	tests := []struct {
		Name            string
		Account         string
		Amount          string
		Recipient       string
		TokenAddr       string
		IsError         bool
		Error           string
		SuggestedParams types.SuggestedParams
	}{
		{
			Name: "Valid transaction",
			SuggestedParams: types.SuggestedParams{
				Fee:             1000,
				FirstRoundValid: 100,
				LastRoundValid:  200,
				GenesisID:       "testnet-v1.0",
				// GenesisHash is normally a 32-byte array encoded in base64.
				GenesisHash: genesisHashBytes,
			},
			Account:   "IBEM5YYEA5UHMRKV3HPRPYB2IJZ4HIII5SW4EDCV5OJFHBKIJBMQWNME6U",
			Amount:    "100",
			Recipient: "ZQK6ICG4FMGMV35J5JPI7HPBVV6P4SBWY2PK35QSNQQ3SEADVHXYX3G7OQ",
			IsError:   false,
		},
		{
			Name: "Error in suggested params",
			SuggestedParams: types.SuggestedParams{
				Fee:             1000,
				FirstRoundValid: 100,
				LastRoundValid:  200,
				GenesisID:       "testnet-v1.0",
				// GenesisHash is normally a 32-byte array encoded in base64.
				GenesisHash: genesisHashBytes,
			},
			Account:   "IBEM5YYEA5UHMRKV3HPRPYB2IJZ4HIII5SW4EDCV5OJFHBKIJBMQWNME6U",
			Amount:    "100",
			Recipient: "ZQK6ICG4FMGMV35J5JPI7HPBVV6P4SBWY2PK35QSNQQ3SEADVHXYX3G7OQ",
			IsError:   true,
			Error:     "failed to get suggested params: error in suggested params",
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			clients := &MockClients{mockAlgod: new(MockAlgodClient), mockIndexer: new(MockIndexerClient)}
			mockSuggestedParamsRequester := new(MockSuggestedParamsRequester)
			clients.mockAlgod.On("SuggestedParams", mock.Anything).Return(mockSuggestedParamsRequester)
			if tt.IsError {
				mockSuggestedParamsRequester.On("Do", mock.Anything).Return(types.SuggestedParams{}, fmt.Errorf("error in suggested params"))
			} else {
				mockSuggestedParamsRequester.On("Do", mock.Anything).Return(tt.SuggestedParams, nil)
			}
			serializedTxn, dataToSign, err := clients.WithdrawAlgorandNativeGetSignature(tt.Account, tt.Amount, tt.Recipient)
			fmt.Println(serializedTxn)
			fmt.Println(dataToSign)
			if !tt.IsError && err != nil {
				t.Errorf("WithdrawAlgorandTokenGetSignature() error = %v", err)
			} else if tt.IsError && err.Error() != tt.Error {
				t.Errorf("expected error = %v", err)
			}
		})
	}

}

func TestWithdrawAlgorandASAGetSignature(t *testing.T) {
	// Testnet

	genesisHashStr := "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI="

	// Decode the Base64 string to get the raw bytes.
	genesisHashBytes, err := base64.StdEncoding.DecodeString(genesisHashStr)
	if err != nil {
		log.Fatalf("failed to decode genesis hash: %v", err)
	}

	tests := []struct {
		Name            string
		Account         string
		Amount          string
		Recipient       string
		AssetId         string
		TransferType    string
		IsError         bool
		Error           string
		SuggestedParams types.SuggestedParams
	}{
		{
			Name: "Valid transaction",
			SuggestedParams: types.SuggestedParams{
				Fee:             1000,
				FirstRoundValid: 100,
				LastRoundValid:  200,
				GenesisID:       "testnet-v1.0",
				// GenesisHash is normally a 32-byte array encoded in base64.
				GenesisHash: genesisHashBytes,
			},
			Account:      "IBEM5YYEA5UHMRKV3HPRPYB2IJZ4HIII5SW4EDCV5OJFHBKIJBMQWNME6U",
			Amount:       "100",
			Recipient:    "ZQK6ICG4FMGMV35J5JPI7HPBVV6P4SBWY2PK35QSNQQ3SEADVHXYX3G7OQ",
			AssetId:      "10458941",
			TransferType: "axfer",
			IsError:      false,
		},
		{
			Name: "Error in suggested params",
			SuggestedParams: types.SuggestedParams{
				Fee:             1000,
				FirstRoundValid: 100,
				LastRoundValid:  200,
				GenesisID:       "testnet-v1.0",
				// GenesisHash is normally a 32-byte array encoded in base64.
				GenesisHash: genesisHashBytes,
			},
			Account:      "IBEM5YYEA5UHMRKV3HPRPYB2IJZ4HIII5SW4EDCV5OJFHBKIJBMQWNME6U",
			Amount:       "100",
			Recipient:    "ZQK6ICG4FMGMV35J5JPI7HPBVV6P4SBWY2PK35QSNQQ3SEADVHXYX3G7OQ",
			AssetId:      "10458941",
			TransferType: "axfer",
			IsError:      true,
			Error:        "failed to get suggested params: error in suggested params",
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			clients := &MockClients{mockAlgod: new(MockAlgodClient), mockIndexer: new(MockIndexerClient)}
			mockSuggestedParamsRequester := new(MockSuggestedParamsRequester)
			clients.mockAlgod.On("SuggestedParams", mock.Anything).Return(mockSuggestedParamsRequester)
			if tt.IsError {
				mockSuggestedParamsRequester.On("Do", mock.Anything).Return(types.SuggestedParams{}, fmt.Errorf("error in suggested params"))
			} else {
				mockSuggestedParamsRequester.On("Do", mock.Anything).Return(tt.SuggestedParams, nil)
			}
			serializedTxn, dataToSign, err := clients.WithdrawAlgorandASAGetSignature(tt.Account, tt.Amount, tt.Recipient, tt.AssetId)
			fmt.Println(serializedTxn)
			fmt.Println(dataToSign)
			fmt.Println(err)
			if !tt.IsError && err != nil {
				t.Errorf("WithdrawAlgorandTokenGetSignature() error = %v", err)
			} else if tt.IsError && err.Error() != tt.Error {
				t.Errorf("expected error = %v", err)
			}
			if !tt.IsError {
				assert.Equal(t, string(dataToSign.Type), tt.TransferType)
			}
		})
	}

}
