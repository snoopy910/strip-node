package sequencer

import (
	"context"
	"encoding/base32"
	"fmt"
	"time"

	"github.com/StripChain/strip-node/common"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/indexer"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/types"
)

// GetAlgorandTransfers retrieves transfer information from an Algorand transaction
// It handles both native ALGO transfers and ASA (Algorand Standard Asset) transfers
func GetAlgorandTransfers(chainId string, txnHash string) ([]common.Transfer, error) {
    chain, err := common.GetChain(chainId)
    if err != nil {
        return nil, err
    }

    // Create an indexer client (no API key needed for AlgoNode)
    indexerClient, err := indexer.MakeClient(chain.IndexerUrl, "")
    if err != nil {
        return nil, fmt.Errorf("failed to create indexer client: %v", err)
    }

    // Create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Look up the transaction
    txnResponse, err := indexerClient.LookupTransaction(txnHash).Do(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to lookup transaction: %v", err)
    }

    // Validate transaction
    if txnResponse.Transaction.ConfirmedRound == 0 {
        return nil, fmt.Errorf("transaction not found or not confirmed")
    }

    txn := txnResponse.Transaction
    var transfers []common.Transfer

    // Handle native ALGO transfer
    if txn.Type == types.PaymentTx {
        // Validate addresses
        if txn.Sender.IsZero() || txn.Receiver.IsZero() {
            return nil, fmt.Errorf("invalid sender or receiver address")
        }

        // ALGO amounts are in microAlgos (1 ALGO = 1,000,000 microAlgos)
        amount := float64(txn.Amount) / 1_000_000

        transfers = append(transfers, common.Transfer{
            From:   txn.Sender.String(),
            To:     txn.Receiver.String(),
            Amount: fmt.Sprintf("%.6f", amount),
            Token:  "ALGO", // Native token
        })
    }

    // Handle ASA (Algorand Standard Asset) transfer
    if txn.Type == types.AssetTransferTx {
        // Validate addresses
        if txn.Sender.IsZero() || txn.AssetReceiver.IsZero() {
            return nil, fmt.Errorf("invalid sender or receiver address")
        }

        // Get asset info for decimals
        assetInfo, err := indexerClient.LookupAssetByID(txn.XferAsset).Do(ctx)
        if err != nil {
            return nil, fmt.Errorf("failed to lookup asset info: %v", err)
        }

        // Calculate amount with proper decimals
        decimals := assetInfo.Asset.Params.Decimals
        amount := float64(txn.AssetAmount) / math.Pow10(int(decimals))

        transfers = append(transfers, common.Transfer{
            From:   txn.Sender.String(),
            To:     txn.AssetReceiver.String(),
            Amount: fmt.Sprintf("%.*f", decimals, amount),
            Token:  fmt.Sprintf("%d", assetInfo.Asset.Index), // ASA ID as token identifier
        })
    }

    return transfers, nil
}

// CheckAlgorandTransactionConfirmed checks if an Algorand transaction is confirmed
// It first tries the Algod API, then falls back to the Indexer if needed
func CheckAlgorandTransactionConfirmed(chainId string, txnHash string) (bool, error) {
    chain, err := common.GetChain(chainId)
    if err != nil {
        return false, err
    }

    // First try using native Algod API (Priority 1)
    algodClient, err := algod.MakeClient(chain.ChainUrl, "")
    if err == nil {
        // Get pending transaction information
        pendingTxn, _, err := algodClient.PendingTransactionInformation(txnHash).Do(context.Background())
        if err == nil {
            // If confirmed round is non-zero, transaction is confirmed
            if pendingTxn.ConfirmedRound > 0 {
                return true, nil
            }
            // If pool error is empty and confirmed round is zero, transaction is still pending
            if pendingTxn.PoolError == "" {
                return false, nil
            }
        }
    }

    // Fallback to Indexer if Algod fails or transaction not found (Priority 2)
    indexerClient, err := indexer.MakeClient(chain.IndexerUrl, "")
    if err != nil {
        return false, fmt.Errorf("failed to create indexer client: %v", err)
    }

    // Look up the transaction
    txnResponse, err := indexerClient.LookupTransaction(txnHash).Do(context.Background())
    if err != nil {
        return false, fmt.Errorf("failed to lookup transaction: %v", err)
    }

    // If we can find the transaction in the indexer, it means it's confirmed
    // The indexer only indexes confirmed transactions
    return txnResponse.Transaction.ConfirmedRound > 0, nil
}

// SendAlgorandTransaction sends a signed Algorand transaction to the network
func SendAlgorandTransaction(serializedTxn string, chainId string, keyCurve string, dataToSign string, signatureBase64 string) (string, error) {
    chain, err := common.GetChain(chainId)
    if err != nil {
        return "", err
    }

    // Create an algod client (no API key needed for AlgoNode)
    algodClient, err := algod.MakeClient(chain.ChainUrl, "")
    if err != nil {
        return "", fmt.Errorf("failed to create algod client: %v", err)
    }

    // Decode the serialized transaction (base32 encoded)
    txnBytes, err := base32.StdEncoding.DecodeString(serializedTxn)
    if err != nil {
        return "", fmt.Errorf("failed to decode serialized transaction: %v", err)
    }

    // Decode the signature (base32 encoded)
    sigBytes, err := base32.StdEncoding.DecodeString(signatureBase64)
    if err != nil {
        return "", fmt.Errorf("failed to decode signature: %v", err)
    }

    // Deserialize the transaction using msgpack
    var txn types.Transaction
    err = msgpack.Decode(txnBytes, &txn)
    if err != nil {
        return "", fmt.Errorf("failed to deserialize transaction: %v", err)
    }

    // Create a signed transaction with the provided signature
    signedTxn := types.SignedTxn{
        Txn: txn,
        Sig: types.Signature(sigBytes),
    }

    // Encode the signed transaction using msgpack
    stxnBytes, err := msgpack.Encode(signedTxn)
    if err != nil {
        return "", fmt.Errorf("failed to encode signed transaction: %v", err)
    }

    // Send the transaction
    txid, err := algodClient.SendRawTransaction(stxnBytes).Do(context.Background())
    if err != nil {
        return "", fmt.Errorf("failed to send transaction: %v", err)
    }

    return txid, nil
}
