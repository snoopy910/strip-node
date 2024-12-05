package sequencer

import (
	"os"
	"testing"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/stretchr/testify/assert"
)

var testDB *pg.DB

func TestMain(m *testing.M) {
	// Setup test database connection
	testDB = pg.Connect(&pg.Options{
		User:     "test_user",
		Password: "test_password",
		Database: "test_db",
		Addr:     "localhost:5433",
	})

	// Store original client
	originalClient := client

	// Set test client
	client = testDB

	// Create tables
	models := []interface{}{
		&IntentSchema{},
		&OperationSchema{},
		&WalletSchema{},
		&LockSchema{},
	}

	for _, model := range models {
		err := testDB.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			panic(err)
		}
	}

	// Run tests
	code := m.Run()

	// Cleanup
	client = originalClient
	testDB.Close()

	os.Exit(code)
}

func clearTables(t *testing.T) {
	tables := []interface{}{
		&IntentSchema{},
		&OperationSchema{},
		&WalletSchema{},
		&LockSchema{},
	}

	for _, table := range tables {
		_, err := testDB.Model(table).Where("1 = 1").Delete()
		assert.NoError(t, err)
	}
}

func TestAddIntent(t *testing.T) {
	clearTables(t)

	intent := &Intent{
		Signature:     "test_sig",
		Identity:      "test_identity",
		IdentityCurve: "test_curve",
		Expiry:        uint64(time.Now().Add(time.Hour).Unix()),
		Operations: []Operation{
			{
				SerializedTxn: "test_txn",
				DataToSign:    "test_data",
				ChainId:       "test_chain",
				KeyCurve:      "test_key_curve",
				Type:          "test_type",
				Solver:        "test_solver",
			},
		},
	}

	// Test adding intent
	id, err := AddIntent(intent)
	assert.NoError(t, err)
	assert.Greater(t, id, int64(0))

	// Test retrieving intent
	retrievedIntent, err := GetIntent(id)
	assert.NoError(t, err)
	assert.Equal(t, intent.Identity, retrievedIntent.Identity)
	assert.Equal(t, intent.IdentityCurve, retrievedIntent.IdentityCurve)
	assert.Equal(t, len(intent.Operations), len(retrievedIntent.Operations))
}

func TestLockOperations(t *testing.T) {
	clearTables(t)

	// Test adding lock
	id, err := AddLock("test_identity", "test_curve")
	assert.NoError(t, err)
	assert.Greater(t, id, int64(0))

	// Test locking
	err = LockIdentity(id)
	assert.NoError(t, err)

	// Test getting lock
	lock, err := GetLock("test_identity", "test_curve")
	assert.NoError(t, err)
	assert.True(t, lock.Locked)

	// Test unlocking
	err = UnlockIdentity(id)
	assert.NoError(t, err)

	lock, err = GetLock("test_identity", "test_curve")
	assert.NoError(t, err)
	assert.False(t, lock.Locked)
}

func TestWalletOperations(t *testing.T) {
	clearTables(t)

	wallet := &WalletSchema{
		Identity:       "test_identity",
		IdentityCurve:  "test_curve",
		EDDSAPublicKey: "test_eddsa",
		ECDSAPublicKey: "test_ecdsa",
		Signers:        "test_signers",
	}

	// Test adding wallet
	id, err := AddWallet(wallet)
	assert.NoError(t, err)
	assert.Greater(t, id, int64(0))

	// Test getting wallet
	retrievedWallet, err := GetWallet("test_identity", "test_curve")
	assert.NoError(t, err)
	assert.Equal(t, wallet.Identity, retrievedWallet.Identity)
	assert.Equal(t, wallet.EDDSAPublicKey, retrievedWallet.EDDSAPublicKey)
	assert.Equal(t, wallet.ECDSAPublicKey, retrievedWallet.ECDSAPublicKey)
}

func TestGetIntentsWithPagination(t *testing.T) {
	clearTables(t)

	// Add some test intents
	for i := 0; i < 5; i++ {
		intent := &Intent{
			Signature:     "test_sig",
			Identity:      "test_identity",
			IdentityCurve: "test_curve",
			Expiry:        uint64(time.Now().Add(time.Hour).Unix()),
			Operations:    []Operation{},
		}
		_, err := AddIntent(intent)
		assert.NoError(t, err)
	}

	// Test pagination
	intents, count, err := GetIntentsWithPagination(2, 0)
	assert.NoError(t, err)
	assert.Equal(t, 5, count)
	assert.Equal(t, 2, len(intents))
}

func TestUpdateOperations(t *testing.T) {
	clearTables(t)

	// Add test intent with operation
	intent := &Intent{
		Signature:     "test_sig",
		Identity:      "test_identity",
		IdentityCurve: "test_curve",
		Expiry:        uint64(time.Now().Add(time.Hour).Unix()),
		Operations: []Operation{
			{
				SerializedTxn: "test_txn",
				DataToSign:    "test_data",
				ChainId:       "test_chain",
				KeyCurve:      "test_key_curve",
				Type:          "test_type",
				Solver:        "test_solver",
			},
		},
	}

	intentId, err := AddIntent(intent)
	assert.NoError(t, err)

	retrievedIntent, err := GetIntent(intentId)
	assert.NoError(t, err)
	operationId := retrievedIntent.Operations[0].ID

	// Test updating operation status
	err = UpdateOperationStatus(operationId, "completed")
	assert.NoError(t, err)

	// Test updating operation result
	err = UpdateOperationResult(operationId, "completed", "success")
	assert.NoError(t, err)

	// Test updating intent status
	err = UpdateIntentStatus(intentId, "completed")
	assert.NoError(t, err)
}
