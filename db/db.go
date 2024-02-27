package db

import (
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

var client *pg.DB

type KVStore struct {
	Id    int64
	Key   string
	Value string
}

func createKeyValueSchema(db *pg.DB) error {
	models := []interface{}{
		(*KVStore)(nil),
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

func Initialise(host string, database string, username string, password string) {

	client = pg.Connect(&pg.Options{
		User:     username,
		Password: password,
		Database: database,
		Addr:     host,
	})

	err := createKeyValueSchema(client)
	if err != nil {
		panic(err)
	}
}

func AddKeyShare(identity string, identityCurve string, keyCurve string, key string) error {
	fmt.Println("Adding key share to postgres", identity+"_"+identityCurve+"_"+keyCurve)
	kvStore := &KVStore{
		Key:   identity + "_" + identityCurve + "_" + keyCurve,
		Value: key,
	}

	_, err := client.Model(kvStore).Insert()
	return err
}

func GetKeyShare(identity string, identityCurve string, keyCurve string) (string, error) {
	var keys []KVStore
	err := client.Model(&keys).Where("key = ?", identity+"_"+identityCurve+"_"+keyCurve).Select()

	if err != nil {
		return "", err
	}

	if len(keys) == 0 {
		return "", nil
	}

	return keys[0].Value, nil
}

func AddSignersForKeyShare(identity string, identityCurve string, keyCurve string, signers string) error {
	fmt.Println("Adding signers to postgres", identity+"_"+identityCurve+"_"+keyCurve)
	kvStore := &KVStore{
		Key:   identity + "_" + identityCurve + "_" + keyCurve + "_" + "signers",
		Value: signers,
	}

	_, err := client.Model(kvStore).Insert()
	return err
}

func GetSignersForKeyShare(identity string, identityCurve string, keyCurve string) (string, error) {
	var keys []KVStore
	err := client.Model(&keys).Where("key = ?", identity+"_"+identityCurve+"_"+keyCurve+"_signers").Select()

	if err != nil {
		return "", err
	}

	if len(keys) == 0 {
		return "", nil
	}

	return keys[0].Value, nil
}
