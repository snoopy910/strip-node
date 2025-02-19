package algorand

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/StripChain/strip-node/common"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
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

func TestGetAlgorandTransfers(t *testing.T) {

	txA := models.Transaction{
		Id:     "P5RRK4SQHBC5QZUCKNMCQ6MDMAJSEPCQKEZJCPRV4R3MX3YWICUQ",
		Sender: "IBEM5YYEA5UHMRKV3HPRPYB2IJZ4HIII5SW4EDCV5OJFHBKIJBMQWNME6U",
		PaymentTransaction: models.TransactionPayment{
			Receiver: "IBEM5YYEA5UHMRKV3HPRPYB2IJZ4HIII5SW4EDCV5OJFHBKIJBMQWNME6U",
			Amount:   100,
		},
		ConfirmedRound: 12345,
		// Note:           "test note",
		// AssetID:        0,
		Type: string(types.PaymentTx),
	}

	respA := models.TransactionResponse{
		CurrentRound: 12350,
		Transaction:  txA,
	}

	txB := models.Transaction{
		Id:     "P5RRK4SQHBC5QZUCKNMCQ6MDMAJSEPCQKEZJCPRV4R3MX3YWICUQ",
		Sender: "IBEM5YYEA5UHMRKV3HPRPYB2IJZ4HIII5SW4EDCV5OJFHBKIJBMQWNME6U",
		AssetTransferTransaction: models.TransactionAssetTransfer{
			Receiver: "IBEM5YYEA5UHMRKV3HPRPYB2IJZ4HIII5SW4EDCV5OJFHBKIJBMQWNME6U",
			Amount:   100,
			AssetId:  123,
		},
		ConfirmedRound: 12345,
		Type:           string(types.AssetTransferTx),
	}

	respB := models.TransactionResponse{
		CurrentRound: 12350,
		Transaction:  txB,
	}

	tests := []struct {
		Name                 string
		GenesisHash          string
		TxHash               string
		Error                string
		IsLookUpTxnError     bool
		IsLookUpAssetIdError bool
		IsError              bool
		Response             models.TransactionResponse
	}{
		{
			Name:             "Valid payment transaction",
			GenesisHash:      "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI=",
			TxHash:           "P5RRK4SQHBC5QZUCKNMCQ6MDMAJSEPCQKEZJCPRV4R3MX3YWICUQ",
			Response:         respA,
			Error:            "",
			IsLookUpTxnError: false,
			IsError:          false,
		},
		{
			Name:             "Valid asset transfer transaction",
			GenesisHash:      "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI=",
			TxHash:           "P5RRK4SQHBC5QZUCKNMCQ6MDMAJSEPCQKEZJCPRV4R3MX3YWICUQ",
			Response:         respB,
			Error:            "",
			IsLookUpTxnError: false,
			IsError:          false,
		},
		{
			Name:             "Look up transactioninfo error",
			GenesisHash:      "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI=",
			TxHash:           "P5RRK4SQHBC5QZUCKNMCQ6MDMAJSEPCQKEZJCPRV4R3MX3YWICUQ",
			Error:            "failed to lookup transaction: error in look up transaction",
			IsLookUpTxnError: true,
			IsError:          true,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			clients := &MockClients{mockAlgod: new(MockAlgodClient), mockIndexer: new(MockIndexerClient)}
			mockLookupTransactionRequester := new(MockLookupTransactionRequester)
			clients.mockIndexer.On("LookupTransaction", mock.Anything).Return(mockLookupTransactionRequester)
			if !tt.IsLookUpTxnError {
				mockLookupTransactionRequester.On("Do", mock.Anything).Return(tt.Response, nil)
			} else {
				mockLookupTransactionRequester.On("Do", mock.Anything).Return(models.TransactionResponse{}, fmt.Errorf("error in look up transaction"))
			}

			mockLookupAssetIdRequester := new(MockLookupAssetByIDRequester)
			clients.mockIndexer.On("LookupAssetByID", mock.Anything).Return(mockLookupAssetIdRequester)
			if !tt.IsLookUpAssetIdError {
				mockLookupAssetIdRequester.On("Do", mock.Anything).Return(uint64(334), models.Asset{Index: 123}, nil)
			} else {
				mockLookupAssetIdRequester.On("Do", mock.Anything).Return(uint64(0), models.Asset{}, fmt.Errorf("error in look up asset"))
			}

			serializedTxn, err := clients.GetAlgorandTransfers(tt.GenesisHash, tt.TxHash)
			fmt.Println(serializedTxn)
			fmt.Println(err)
			if !tt.IsError && err != nil {
				t.Errorf("WithdrawAlgorandTxn() error = %v", err)
			} else if tt.IsError && err.Error() != tt.Error {
				t.Errorf("expected error = %v", err)
			}
		})
	}

}

func TestCheckAlgorandTransactionConfirmed(t *testing.T) {

	respA := models.PendingTransactionInfoResponse{
		ConfirmedRound: 12350,
		PoolError:      "sss",
	}

	respB := models.TransactionResponse{
		Transaction: models.Transaction{
			Id:     "P5RRK4SQHBC5QZUCKNMCQ6MDMAJSEPCQKEZJCPRV4R3MX3YWICUQ",
			Sender: "IBEM5YYEA5UHMRKV3HPRPYB2IJZ4HIII5SW4EDCV5OJFHBKIJBMQWNME6U",
			PaymentTransaction: models.TransactionPayment{
				Receiver: "IBEM5YYEA5UHMRKV3HPRPYB2IJZ4HIII5SW4EDCV5OJFHBKIJBMQWNME6U",
				Amount:   100,
			},
			ConfirmedRound: 12345,
			// Note:           "test note",
			// AssetID:        0,
			Type: string(types.PaymentTx),
		},
	}

	tests := []struct {
		Name                  string
		GenesisHash           string
		TxHash                string
		Error                 string
		IsErrorPendingTxnInfo bool
		IsErrorLookupTxnInfo  bool
		IsError               bool
	}{
		{
			Name:                  "Pending transaction confirmed",
			GenesisHash:           "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI=",
			TxHash:                "P5RRK4SQHBC5QZUCKNMCQ6MDMAJSEPCQKEZJCPRV4R3MX3YWICUQ",
			IsErrorPendingTxnInfo: false,
			IsErrorLookupTxnInfo:  false,
			IsError:               false,
		},
		{
			Name:                  "Pending transaction not confirmed",
			GenesisHash:           "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI=",
			TxHash:                "P5RRK4SQHBC5QZUCKNMCQ6MDMAJSEPCQKEZJCPRV4R3MX3YWICUQ",
			IsErrorPendingTxnInfo: true,
			IsErrorLookupTxnInfo:  false,
			IsError:               false,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			clients := &MockClients{mockAlgod: new(MockAlgodClient), mockIndexer: new(MockIndexerClient)}
			mockPendingTransactionInformationRequester := new(MockPendingTransactionInformationRequester)
			clients.mockAlgod.On("PendingTransactionInformation", mock.Anything).Return(mockPendingTransactionInformationRequester)
			if !tt.IsErrorPendingTxnInfo {
				mockPendingTransactionInformationRequester.On("Do", mock.Anything).Return(respA, types.SignedTxn{}, nil)
			} else {
				mockPendingTransactionInformationRequester.On("Do", mock.Anything).Return(models.PendingTransactionInfoResponse{}, types.SignedTxn{}, errors.New("error pending transaction information"))
			}

			mockLookupTransactionRequester := new(MockLookupTransactionRequester)
			clients.mockIndexer.On("LookupTransaction", mock.Anything).Return(mockLookupTransactionRequester)
			if !tt.IsErrorLookupTxnInfo {
				mockLookupTransactionRequester.On("Do", mock.Anything).Return(respB, nil)
			} else {
				mockLookupTransactionRequester.On("Do", mock.Anything).Return(models.TransactionResponse{}, fmt.Errorf("error in look up transaction"))
			}

			serializedTxn, err := clients.CheckAlgorandTransactionConfirmed(tt.GenesisHash, tt.TxHash)
			fmt.Println(serializedTxn)
			fmt.Println(err)
			if !tt.IsError && err != nil {
				t.Errorf("WithdrawAlgorandTxn() error = %v", err)
			} else if tt.IsError && err.Error() != tt.Error {
				t.Errorf("expected error = %v", err)
			}
		})
	}

}
