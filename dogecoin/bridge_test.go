package dogecoin

import (
	"testing"

	"github.com/test-go/testify/require"
)

func TestWithdrawDogeNativeGetSignature(t *testing.T) {
	// Test valid withdrawal
	serializedTx, dataToSign, err := WithdrawDogeNativeGetSignature("rpcURL", "naPvQVSd2YGVjuvc1xDFeTAWP4qixahRr6", "1000000", "nVKnX46PTjyTaVXc6m9uLncgGXqC6ZDUww")
	require.NoError(t, err)
	require.NotEmpty(t, serializedTx)
	require.Len(t, dataToSign, 64)

	// Test invalid recipient address
	_, _, err = WithdrawDogeNativeGetSignature("rpcURL", "test_account", "1000000", "invalid_address")
	require.Error(t, err)
}
