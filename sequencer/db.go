package sequencer

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

var client *pg.DB

type IntentSchema struct {
	Id            int64
	Signature     string
	Identity      string
	IdentityCurve string
	Status        string
}

type OperationSchema struct {
	Id             int64
	IntentId       int64
	Intent         *IntentSchema `pg:"rel:has-one"`
	SerializedTxn  string
	DataToSign     string
	ChainId        string
	KeyCurve       string
	Status         string
	Result         string
	Type           string
	Solver         string
	SolverMetadata string
}

type WalletSchema struct {
	Id             int64  `json:"id"`
	Identity       string `json:"identity"`
	IdentityCurve  string `json:"identityCurve"`
	EDDSAPublicKey string `json:"eddsaPublicKey"`
	ECDSAPublicKey string `json:"ecdsaPublicKey"`
	Signers        string `json:"signers"`
}

func createSchemas(db *pg.DB) error {
	models := []interface{}{
		(*IntentSchema)(nil),
		(*OperationSchema)(nil),
		(*WalletSchema)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func InitialiseDB(host string, database string, username string, password string) {

	client = pg.Connect(&pg.Options{
		User:     username,
		Password: password,
		Database: database,
		Addr:     host,
	})

	err := createSchemas(client)
	if err != nil {
		panic(err)
	}
}

func AddIntent(
	Intent *Intent,
) (int64, error) {
	intentSchema := &IntentSchema{
		Signature:     Intent.Signature,
		Identity:      Intent.Identity,
		IdentityCurve: Intent.IdentityCurve,
		Status:        INTENT_STATUS_PROCESSING,
	}

	_, err := client.Model(intentSchema).Insert()
	if err != nil {
		return 0, err
	}

	for _, operation := range Intent.Operations {
		operationSchema := &OperationSchema{
			IntentId:       intentSchema.Id,
			SerializedTxn:  operation.SerializedTxn,
			DataToSign:     operation.DataToSign,
			ChainId:        operation.ChainId,
			KeyCurve:       operation.KeyCurve,
			Status:         OPERATION_STATUS_PENDING,
			Result:         "",
			Type:           operation.Type,
			Solver:         operation.Solver,
			SolverMetadata: operation.SolverMetadata,
		}

		_, err := client.Model(operationSchema).Insert()
		if err != nil {
			return 0, err
		}
	}

	return intentSchema.Id, nil
}

func GetIntent(intentId int64) (*Intent, error) {
	var intentSchema IntentSchema
	err := client.Model(&intentSchema).Where("id = ?", intentId).Select()
	if err != nil {
		return nil, err
	}

	var operationsSchema []OperationSchema
	err = client.Model(&operationsSchema).Where("intent_id = ?", intentSchema.Id).Select()
	if err != nil {
		return nil, err
	}

	operations := make([]Operation, len(operationsSchema))
	for i, operationSchema := range operationsSchema {
		operations[i] = Operation{
			ID:             operationSchema.Id,
			SerializedTxn:  operationSchema.SerializedTxn,
			DataToSign:     operationSchema.DataToSign,
			ChainId:        operationSchema.ChainId,
			KeyCurve:       operationSchema.KeyCurve,
			Status:         operationSchema.Status,
			Result:         operationSchema.Result,
			Type:           operationSchema.Type,
			Solver:         operationSchema.Solver,
			SolverMetadata: operationSchema.SolverMetadata,
		}
	}

	intent := &Intent{
		ID:            intentSchema.Id,
		Operations:    operations,
		Signature:     intentSchema.Signature,
		Identity:      intentSchema.Identity,
		IdentityCurve: intentSchema.IdentityCurve,
		Status:        intentSchema.Status,
	}

	return intent, nil
}

func GetIntents(status string) ([]*Intent, error) {
	var intentSchemas []IntentSchema
	err := client.Model(&intentSchemas).Where("status = ?", status).Select()
	if err != nil {
		return nil, err
	}

	intents := make([]*Intent, len(intentSchemas))
	for i, intentSchema := range intentSchemas {
		intent, err := GetIntent(intentSchema.Id)
		if err != nil {
			return nil, err
		}

		var operationsSchema []OperationSchema
		err = client.Model(&operationsSchema).Where("intent_id = ?", intentSchema.Id).Select()

		if err != nil {
			return nil, err
		}

		operations := make([]Operation, len(operationsSchema))
		for i, operationSchema := range operationsSchema {
			operations[i] = Operation{
				SerializedTxn:  operationSchema.SerializedTxn,
				DataToSign:     operationSchema.DataToSign,
				ChainId:        operationSchema.ChainId,
				KeyCurve:       operationSchema.KeyCurve,
				Status:         operationSchema.Status,
				Result:         operationSchema.Result,
				Type:           operationSchema.Type,
				Solver:         operationSchema.Solver,
				SolverMetadata: operationSchema.SolverMetadata,
			}
		}

		intent.Operations = operations
		intents[i] = intent
	}

	return intents, nil
}

func UpdateOperationResult(operationId int64, status string, result string) error {
	operationSchema := &OperationSchema{
		Id:     operationId,
		Status: status,
		Result: result,
	}

	_, err := client.Model(operationSchema).Column("status", "txn_hash").WherePK().Update()
	if err != nil {
		return err
	}

	return nil
}

func UpdateOperationStatus(operationId int64, status string) error {
	operationSchema := &OperationSchema{
		Id:     operationId,
		Status: status,
	}

	_, err := client.Model(operationSchema).Column("status").WherePK().Update()
	if err != nil {
		return err
	}

	return nil
}

func UpdateIntentStatus(intentId int64, status string) error {
	intentSchema := &IntentSchema{
		Id:     intentId,
		Status: status,
	}

	_, err := client.Model(intentSchema).Column("status").WherePK().Update()
	if err != nil {
		return err
	}

	return nil
}

func GetWallet(identity string, identityCurve string) (*WalletSchema, error) {
	var walletSchema WalletSchema
	err := client.Model(&walletSchema).Where("identity = ? AND identity_curve = ?", identity, identityCurve).Select()
	if err != nil {
		return nil, err
	}

	return &walletSchema, nil
}

func AddWallet(wallet *WalletSchema) (int64, error) {
	_, err := client.Model(wallet).Insert()
	if err != nil {
		return 0, err
	}

	return wallet.Id, nil
}
