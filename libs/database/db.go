package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/StripChain/strip-node/libs"
	"github.com/StripChain/strip-node/util/logger"
	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var (
	dbClient *pg.DB
	once     sync.Once
)

// GetDB returns the singleton database client instance.
// It assumes InitialiseDB has been called successfully at least once.
func GetDB() *pg.DB {
	if dbClient == nil {
		panic("database client is not initialized. Call InitialiseDB first.")
	}
	return dbClient
}

type IntentSchema struct {
	tableName     struct{}          `pg:"intents"` //lint:ignore U1000 ok
	Id            uuid.UUID         `pg:",type:uuid,notnull,default:gen_random_uuid()"`
	Signature     string            `pg:",notnull"`
	Identity      string            `pg:",notnull"`
	IdentityCurve string            `pg:",notnull"`
	Status        libs.IntentStatus `pg:",notnull"`
	Expiry        time.Time         `pg:",notnull"`
	CreatedAt     time.Time         `pg:",notnull,default:CURRENT_TIMESTAMP"`
}

type OperationSchema struct {
	tableName        struct{} `pg:"operations"` //lint:ignore U1000 ok
	Id               int64
	IntentId         uuid.UUID     `pg:",type:uuid,notnull"`
	Intent           *IntentSchema `pg:"rel:has-one"`
	SerializedTxn    string        `pg:",notnull"`
	DataToSign       string        `pg:",notnull"`
	ChainId          string        `pg:",notnull"`
	GenesisHash      string
	KeyCurve         string               `pg:",notnull"`
	Status           libs.OperationStatus `pg:",notnull"`
	Result           string
	Type             libs.OperationType `pg:",notnull"`
	Solver           string
	SolverMetadata   string `pg:",type:jsonb"`
	SolverDataToSign string
	SolverOutput     string    `pg:",type:jsonb"`
	CreatedAt        time.Time `pg:",notnull,default:CURRENT_TIMESTAMP"`
}

type WalletSchema struct {
	tableName                struct{} `pg:"wallets"` //lint:ignore U1000 ok
	Id                       int64    `json:"id"`
	Identity                 string   `json:"identity" pg:",notnull"`
	IdentityCurve            string   `json:"identityCurve" pg:",notnull"`
	EDDSAPublicKey           string   `json:"eddsaPublicKey"`
	AptosEDDSAPublicKey      string   `json:"aptosEddsaPublicKey"`
	ECDSAPublicKey           string   `json:"ecdsaPublicKey"`
	BitcoinMainnetPublicKey  string   `json:"bitcoinMainnetPublicKey"`
	BitcoinTestnetPublicKey  string   `json:"bitcoinTestnetPublicKey"`
	BitcoinRegtestPublicKey  string   `json:"bitcoinRegtestPublicKey"`
	StellarPublicKey         string   `json:"stellarPublicKey"`
	DogecoinMainnetPublicKey string   `json:"dogecoinMainnetPublicKey"`
	DogecoinTestnetPublicKey string   `json:"dogecoinTestnetPublicKey"`
	SuiPublicKey             string   `json:"suiPublicKey"`
	AlgorandEDDSAPublicKey   string   `json:"algorandEddsaPublicKey"`
	RippleEDDSAPublicKey     string   `json:"rippleEddsaPublicKey"`
	CardanoPublicKey         string   `json:"cardanoPublicKey"`
	Signers                  []string `json:"signers" pg:",type:jsonb"`
}

type LockSchema struct {
	tableName     struct{} `pg:"locks"` //lint:ignore U1000 ok
	Id            int64    `json:"id"`
	Identity      string   `json:"identity" pg:",notnull"`
	IdentityCurve string   `json:"identityCurve" pg:",notnull"`
	Locked        bool     `json:"locked" pg:",notnull,default:false"`
}

type HeartbeatSchema struct {
	tableName struct{}  `pg:"heartbeats"` //lint:ignore U1000 ok
	PublicKey string    `pg:"publickey,pk,notnull"`
	UpdatedAt time.Time `pg:"updated_at,notnull,default:CURRENT_TIMESTAMP"`
}

// Add these constants for pool configuration
const (
	minPoolSize     = 2
	maxPoolSize     = 10
	maxConnIdleTime = 30 * time.Minute
)

// func createSchemas(db *pg.DB) error {
// 	models := []interface{}{
// 		(*IntentSchema)(nil),
// 		(*OperationSchema)(nil),
// 		(*WalletSchema)(nil),
// 		(*LockSchema)(nil),
// 		(*HeartbeatSchema)(nil),
// 	}

// 	for _, model := range models {
// 		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
// 			IfNotExists: true,
// 		})
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

func InitialiseDB(host string, database string, username string, password string) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", username, password, host, database)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(fmt.Errorf("error opening database connection for migration: %v", err))
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(fmt.Errorf("error pinging database for migration: %v", err))
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{
		MigrationsTable: "migrations",
		DatabaseName:    database,
	})
	if err != nil {
		panic(fmt.Errorf("error creating the migration database driver: %v", err))
	}

	ex, err := os.Executable()
	if err != nil {
		panic(fmt.Errorf("error getting executable path: %v", err))
	}
	exPath := filepath.Dir(ex)
	migrationsPath := filepath.Join(exPath, "migrations")

	logger.Sugar().Infof("Looking for migrations in: %s", migrationsPath)

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		database,
		driver)
	if err != nil {
		panic(fmt.Errorf("error creating migration instance: %v", err))
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		panic(fmt.Errorf("error applying migrations: %v", err))
	} else if err == migrate.ErrNoChange {
		logger.Sugar().Info("Database migrations: No changes detected.")
	} else {
		logger.Sugar().Info("Database migrations applied successfully.")
	}

	sourceErr, dbErr := m.Close()
	if sourceErr != nil {
		logger.Sugar().Warnw("Warning: error closing migration source", "error", sourceErr)
	}
	if dbErr != nil {
		logger.Sugar().Warnw("Warning: error closing migration database connection", "error", dbErr)
	}

	once.Do(func() {
		opts := &pg.Options{
			User:                  username,
			Password:              password,
			Database:              database,
			Addr:                  host,
			MinIdleConns:          minPoolSize,
			MaxConnAge:            maxConnIdleTime,
			PoolSize:              maxPoolSize,
			PoolTimeout:           30 * time.Second,
			IdleTimeout:           maxConnIdleTime,
			MaxRetries:            3,
			RetryStatementTimeout: true,
		}
		dbClient = pg.Connect(opts)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := dbClient.Ping(ctx); err != nil {
			dbClient = nil
			panic(fmt.Sprintf("Error connecting main DB client (go-pg): %v", err))
		}
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := GetDB().Ping(ctx); err != nil {
		panic(fmt.Sprintf("Error pinging main DB client after initialization: %v", err))
	}

	logger.Sugar().Info("Database initialised successfully.")
}

func AddLock(identity string, identityCurve string) (int64, error) {
	lock := &LockSchema{
		Identity:      identity,
		IdentityCurve: identityCurve,
		Locked:        false,
	}

	_, err := GetDB().Model(lock).Insert()
	if err != nil {
		return 0, err
	}

	return lock.Id, nil
}

func LockIdentity(id int64) error {
	lockSchema := LockSchema{
		Id:     id,
		Locked: true,
	}

	_, err := GetDB().Model(&lockSchema).Column("locked").WherePK().Update()
	if err != nil {
		return err
	}

	return nil
}

func GetLock(identity string, identityCurve string) (*LockSchema, error) {
	var lockSchema LockSchema
	err := GetDB().Model(&lockSchema).Where("identity = ? AND identity_curve = ?", identity, identityCurve).Select()
	if err != nil {
		return nil, err
	}

	return &lockSchema, nil
}

func UnlockIdentity(id int64) error {
	lockSchema := LockSchema{
		Id:     id,
		Locked: false,
	}

	_, err := GetDB().Model(&lockSchema).Column("locked").WherePK().Update()
	if err != nil {
		return err
	}

	return nil
}

func AddIntent(
	Intent *libs.Intent,
) (uuid.UUID, error) {
	intentSchema := &IntentSchema{
		Signature:     Intent.Signature,
		Identity:      Intent.Identity,
		IdentityCurve: Intent.IdentityCurve,
		Status:        libs.IntentStatusProcessing,
		Expiry:        Intent.Expiry,
	}

	_, err := GetDB().Model(intentSchema).Insert()
	if err != nil {
		return uuid.Nil, err
	}

	for _, operation := range Intent.Operations {
		operationSchema := &OperationSchema{
			IntentId:         intentSchema.Id,
			SerializedTxn:    operation.SerializedTxn,
			DataToSign:       operation.DataToSign,
			ChainId:          operation.ChainId,
			GenesisHash:      operation.GenesisHash,
			KeyCurve:         operation.KeyCurve,
			Status:           libs.OperationStatusPending,
			Result:           "",
			Type:             operation.Type,
			Solver:           operation.Solver,
			SolverMetadata:   operation.SolverMetadata,
			SolverDataToSign: operation.SolverDataToSign,
		}

		_, err := GetDB().Model(operationSchema).Insert()
		if err != nil {
			return uuid.Nil, err
		}
	}

	return intentSchema.Id, nil
}

func GetIntent(id uuid.UUID) (*libs.Intent, error) {
	var intentSchema IntentSchema
	err := GetDB().Model(&intentSchema).Where("id = ?", id).Select()
	if err != nil {
		return nil, err
	}

	var operationsSchema []OperationSchema
	err = GetDB().Model(&operationsSchema).Where("intent_id = ?", intentSchema.Id).Select()
	if err != nil {
		return nil, err
	}

	operations := make([]libs.Operation, len(operationsSchema))
	for i, operationSchema := range operationsSchema {
		operations[i] = libs.Operation{
			ID:               operationSchema.Id,
			SerializedTxn:    operationSchema.SerializedTxn,
			DataToSign:       operationSchema.DataToSign,
			ChainId:          operationSchema.ChainId,
			GenesisHash:      operationSchema.GenesisHash,
			KeyCurve:         operationSchema.KeyCurve,
			Status:           operationSchema.Status,
			Result:           operationSchema.Result,
			Type:             operationSchema.Type,
			Solver:           operationSchema.Solver,
			SolverMetadata:   operationSchema.SolverMetadata,
			SolverDataToSign: operationSchema.SolverDataToSign,
			SolverOutput:     operationSchema.SolverOutput,
			CreatedAt:        operationSchema.CreatedAt,
		}
	}

	sort.Slice(operations, func(i, j int) bool {
		a := operations[i]
		b := operations[j]
		return a.ID < b.ID
	})

	intent := &libs.Intent{
		ID:            intentSchema.Id,
		Operations:    operations,
		Signature:     intentSchema.Signature,
		Identity:      intentSchema.Identity,
		IdentityCurve: intentSchema.IdentityCurve,
		Status:        intentSchema.Status,
		Expiry:        intentSchema.Expiry,
		CreatedAt:     intentSchema.CreatedAt,
	}

	return intent, nil
}

func GetOperation(intentId uuid.UUID, operationIndex int64) (*libs.Operation, error) {
	var intentSchema IntentSchema
	err := GetDB().Model(&intentSchema).Where("id = ?", intentId).Select()
	if err != nil {
		return nil, err
	}

	var operationsSchema []OperationSchema
	err = GetDB().Model(&operationsSchema).Where("intent_id = ?", intentSchema.Id).Select()
	if err != nil {
		return nil, err
	}

	operations := make([]libs.Operation, len(operationsSchema))
	for i, operationSchema := range operationsSchema {
		operations[i] = libs.Operation{
			ID:               operationSchema.Id,
			SerializedTxn:    operationSchema.SerializedTxn,
			DataToSign:       operationSchema.DataToSign,
			ChainId:          operationSchema.ChainId,
			GenesisHash:      operationSchema.GenesisHash,
			KeyCurve:         operationSchema.KeyCurve,
			Status:           operationSchema.Status,
			Result:           operationSchema.Result,
			Type:             operationSchema.Type,
			Solver:           operationSchema.Solver,
			SolverMetadata:   operationSchema.SolverMetadata,
			SolverDataToSign: operationSchema.SolverDataToSign,
			SolverOutput:     operationSchema.SolverOutput,
		}
	}

	sort.Slice(operations, func(i, j int) bool {
		a := operations[i]
		b := operations[j]
		return a.ID < b.ID
	})

	return &operations[operationIndex], nil
}

func getIntents(intentSchemas *([]IntentSchema)) ([]*libs.Intent, error) {
	intents := make([]*libs.Intent, len(*intentSchemas))
	for i, intentSchema := range *intentSchemas {
		intent, err := GetIntent(intentSchema.Id)
		if err != nil {
			return nil, err
		}

		var operationsSchema []OperationSchema
		err = GetDB().Model(&operationsSchema).Where("intent_id = ?", intentSchema.Id).Select()

		if err != nil {
			return nil, err
		}

		operations := make([]libs.Operation, len(operationsSchema))
		for i, operationSchema := range operationsSchema {
			operations[i] = libs.Operation{
				SerializedTxn:    operationSchema.SerializedTxn,
				DataToSign:       operationSchema.DataToSign,
				ChainId:          operationSchema.ChainId,
				GenesisHash:      operationSchema.GenesisHash,
				KeyCurve:         operationSchema.KeyCurve,
				Status:           operationSchema.Status,
				Result:           operationSchema.Result,
				Type:             operationSchema.Type,
				Solver:           operationSchema.Solver,
				SolverMetadata:   operationSchema.SolverMetadata,
				SolverDataToSign: operationSchema.SolverDataToSign,
				SolverOutput:     operationSchema.SolverOutput,
			}
		}

		intent.Operations = operations
		intents[i] = intent
	}

	return intents, nil
}

func GetSolverIntents(solver string, limit, skip int) ([]*libs.Intent, int, error) {
	// max limit is 100
	if limit > 100 {
		return nil, 0, errors.New("limit cannot be greater than 100")
	}

	var operationSchemas []OperationSchema
	count, err := GetDB().Model(&operationSchemas).Where("solver = ?", solver).DistinctOn("intent_id").Count()
	if err != nil {
		return nil, 0, err
	}

	err = GetDB().Model(&operationSchemas).Limit(limit).Offset(skip).Where("solver = ?", solver).Order("intent_id DESC").DistinctOn("intent_id").Select()
	if err != nil {
		return nil, 0, err
	}

	var intents []*libs.Intent

	for _, operationSchema := range operationSchemas {
		intent, err := GetIntent(operationSchema.IntentId)
		if err != nil {
			return nil, 0, err
		}

		intents = append(intents, intent)
	}

	return intents, count, nil
}

func GetTotalIntents() (int, error) {
	count, err := GetDB().Model(&IntentSchema{}).Count()
	if err != nil {
		return 0, err
	}

	return count, nil
}

func GetIntentsOfAddress(address string, limit, skip int) ([]*libs.Intent, int, error) {
	// max limit is 100
	if limit > 100 {
		return nil, 0, errors.New("limit cannot be greater than 100")
	}

	var intentSchemas []IntentSchema

	// first search for identity. If length is 0, search for ecdsa, if length is 0, then search for eddsa
	err := GetDB().Model(&intentSchemas).Limit(limit).Offset(skip).Where("identity = ?", address).Order("id DESC").Select()
	if err != nil {
		return nil, 0, err
	}

	if len(intentSchemas) != 0 {
		count, err := GetDB().Model(&intentSchemas).Where("identity = ?", address).Count()

		if err != nil {
			return nil, 0, err
		}

		intents, err := getIntents(&intentSchemas)
		if err != nil {
			return nil, 0, err
		}

		return intents, count, nil
	}

	var walletSchemas []WalletSchema
	err = GetDB().Model(&walletSchemas).Where("eddsa_public_Key = ?", address).Select()

	if err != nil {
		return nil, 0, err
	}

	if len(walletSchemas) != 0 {
		err = GetDB().Model(&intentSchemas).Limit(limit).Offset(skip).Where("identity = ?", walletSchemas[0].Identity).Order("id DESC").Select()
		if err != nil {
			return nil, 0, err
		}

		count, err := GetDB().Model(&intentSchemas).Where("identity = ?", walletSchemas[0].Identity).Count()

		if err != nil {
			return nil, 0, err
		}

		intents, err := getIntents(&intentSchemas)

		if err != nil {
			return nil, 0, err
		}

		return intents, count, nil
	}

	err = GetDB().Model(&walletSchemas).Where("ecdsa_public_Key = ?", address).Select()

	if err != nil {
		return nil, 0, err
	}

	if len(walletSchemas) != 0 {
		err = GetDB().Model(&intentSchemas).Limit(limit).Offset(skip).Where("identity = ?", walletSchemas[0].Identity).Order("id DESC").Select()
		if err != nil {
			return nil, 0, err
		}

		count, err := GetDB().Model(&intentSchemas).Where("identity = ?", walletSchemas[0].Identity).Count()

		if err != nil {
			return nil, 0, err
		}

		intents, err := getIntents(&intentSchemas)

		if err != nil {
			return nil, 0, err
		}

		return intents, count, nil
	}

	intents, err := getIntents(&intentSchemas)
	if err != nil {
		return nil, 0, err
	}

	return intents, 0, nil
}

func GetIntentsWithPagination(limit, skip int) ([]*libs.Intent, int, error) {
	// max limit is 100
	if limit > 100 {
		return nil, 0, errors.New("limit cannot be greater than 100")
	}

	var intentSchemas []IntentSchema
	count, err := GetDB().Model(&intentSchemas).Count()

	if err != nil {
		return nil, 0, err
	}

	err = GetDB().Model(&intentSchemas).Limit(limit).Offset(skip).Order("id DESC").Select()
	if err != nil {
		return nil, 0, err
	}

	intents, err := getIntents(&intentSchemas)
	if err != nil {
		return nil, 0, err
	}

	return intents, count, nil
}

func GetIntentsWithStatus(status libs.IntentStatus) ([]*libs.Intent, error) {
	var intentSchemas []IntentSchema
	err := GetDB().Model(&intentSchemas).Where("status = ?", status).Select()
	if err != nil {
		return nil, err
	}

	return getIntents(&intentSchemas)
}

func UpdateOperationResult(operationId int64, status libs.OperationStatus, result string) error {
	operationSchema := &OperationSchema{
		Id:     operationId,
		Status: status,
		Result: result,
	}

	_, err := GetDB().Model(operationSchema).Column("status", "result").WherePK().Update()
	if err != nil {
		return err
	}

	return nil
}

func UpdateOperationStatus(operationId int64, status libs.OperationStatus) error {
	operationSchema := &OperationSchema{
		Id:     operationId,
		Status: status,
	}

	_, err := GetDB().Model(operationSchema).Column("status").WherePK().Update()
	if err != nil {
		return err
	}

	return nil
}

func UpdateOperationSolverOutput(operationId int64, result string) error {
	operationSchema := &OperationSchema{
		Id:           operationId,
		SolverOutput: result,
	}

	_, err := GetDB().Model(operationSchema).Column("solver_output").WherePK().Update()
	if err != nil {
		return err
	}

	return nil
}

func UpdateOperationSolverDataToSign(operationId int64, result string) error {
	operationSchema := &OperationSchema{
		Id:               operationId,
		SolverDataToSign: result,
	}

	_, err := GetDB().Model(operationSchema).Column("solver_data_to_sign").WherePK().Update()
	if err != nil {
		return err
	}

	return nil
}

func UpdateIntentStatus(id uuid.UUID, status libs.IntentStatus) error {
	intentSchema := &IntentSchema{
		Id:     id,
		Status: status,
	}

	_, err := GetDB().Model(intentSchema).Column("status").WherePK().Update()
	if err != nil {
		return err
	}

	return nil
}

func GetWallet(identity string, identityCurve string) (*WalletSchema, error) {
	var walletSchema WalletSchema
	err := GetDB().Model(&walletSchema).Where("identity = ? AND identity_curve = ?", identity, identityCurve).Select()
	if err != nil {
		return nil, err
	}

	return &walletSchema, nil
}

var AddWallet = func(wallet *WalletSchema) (int64, error) {
	_, err := GetDB().Model(wallet).Insert()
	if err != nil {
		return 0, err
	}

	return wallet.Id, nil
}

func AddHeartbeat(publicKey string) error {
	heartbeat := &HeartbeatSchema{
		PublicKey: publicKey,
		UpdatedAt: time.Now(),
	}
	_, err := GetDB().Model(heartbeat).Insert()
	if err != nil {
		return err
	}
	return nil
}

func UpdateHeartbeat(publicKey string) error {
	heartbeat := &HeartbeatSchema{
		PublicKey: publicKey,
		UpdatedAt: time.Now(),
	}
	_, err := GetDB().Model(heartbeat).
		Set("updated_at = ?", heartbeat.UpdatedAt).
		Where("publickey = ?", heartbeat.PublicKey).
		Update()
	if err != nil {
		return err
	}
	return nil
}

func GetHeartbeat(publicKey string) (HeartbeatSchema, error) {
	heartbeat := &HeartbeatSchema{
		PublicKey: publicKey,
	}
	err := GetDB().Model(heartbeat).
		Where("publickey = ?", heartbeat.PublicKey).
		Select()
	if err != nil {
		return HeartbeatSchema{}, err
	}
	return *heartbeat, nil
}

func GetHeartbeats() ([]HeartbeatSchema, error) {
	var heartbeats []HeartbeatSchema
	err := GetDB().Model(&heartbeats).Select()
	if err != nil {
		return nil, err
	}
	return heartbeats, nil
}

func DeleteHeartbeat(publicKey string) error {
	heartbeat := &HeartbeatSchema{
		PublicKey: publicKey,
	}
	_, err := GetDB().Model(heartbeat).Delete()
	if err != nil {
		return err
	}
	return nil
}

func IsSignerAlive(publicKey string) bool {
	heartbeat := &HeartbeatSchema{
		PublicKey: publicKey,
	}
	err := GetDB().Model(heartbeat).Last()
	if err != nil {
		return false
	}
	if time.Since(heartbeat.UpdatedAt) > libs.HEARTBEAT_TIMEOUT {
		return false
	}
	return true
}

func GetActiveSigners() ([]HeartbeatSchema, error) {
	var heartbeats []HeartbeatSchema
	err := GetDB().Model((*HeartbeatSchema)(nil)).
		ColumnExpr("distinct publickey").
		Where("updated_at > ?", time.Now().Add(-libs.HEARTBEAT_TIMEOUT)).
		Select(&heartbeats)
	if err != nil {
		return nil, err
	}
	return heartbeats, nil
}
