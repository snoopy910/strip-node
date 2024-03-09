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
	Id            int64
	IntentId      int64
	Intent        *IntentSchema `pg:"rel:has-one"`
	SerializedTxn string
	DataToSign    string
	ChainId       string
	KeyCurve      string
	Status        string
	TxnHash       string
}

func createSchemas(db *pg.DB) error {
	models := []interface{}{
		(*IntentSchema)(nil),
		(*OperationSchema)(nil),
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
			IntentId:      intentSchema.Id,
			SerializedTxn: operation.SerializedTxn,
			DataToSign:    operation.DataToSign,
			ChainId:       operation.ChainId,
			KeyCurve:      operation.KeyCurve,
			Status:        OPERATION_STATUS_PENDING,
			TxnHash:       "",
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
			SerializedTxn: operationSchema.SerializedTxn,
			DataToSign:    operationSchema.DataToSign,
			ChainId:       operationSchema.ChainId,
			KeyCurve:      operationSchema.KeyCurve,
			Status:        operationSchema.Status,
			TxnHash:       operationSchema.TxnHash,
		}
	}

	intent := &Intent{
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
				SerializedTxn: operationSchema.SerializedTxn,
				DataToSign:    operationSchema.DataToSign,
				ChainId:       operationSchema.ChainId,
				KeyCurve:      operationSchema.KeyCurve,
				Status:        operationSchema.Status,
				TxnHash:       operationSchema.TxnHash,
			}
		}

		intent.Operations = operations
		intents[i] = intent
	}

	return intents, nil
}
