package main

import (
	"github.com/StripChain/strip-node/common"
	"github.com/StripChain/strip-node/util/logger"
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

func InitialiseDB(host string, database string, username string, password string) {

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

func AddKeyShare(identity string, identityCurve common.Curve, keyCurve common.Curve, key string) error {
	logger.Sugar().Infof("Adding key share to postgres %s_%s_%s", identity, identityCurve, keyCurve)
	kvStore := &KVStore{
		Key:   identity + "_" + string(identityCurve) + "_" + string(keyCurve),
		Value: key,
	}

	_, err := client.Model(kvStore).Insert()
	return err
}

func GetKeyShare(identity string, identityCurve common.Curve, keyCurve common.Curve) (string, error) {
	var keys []KVStore
	err := client.Model(&keys).Where("key = ?", identity+"_"+string(identityCurve)+"_"+string(keyCurve)).Select()

	if err != nil {
		return "", err
	}

	if len(keys) == 0 {
		return "", nil
	}

	return keys[0].Value, nil
}

func AddSignersForKeyShare(identity string, identityCurve common.Curve, keyCurve common.Curve, signers string) error {
	logger.Sugar().Infof("Adding signers to postgres %s_%s_%s", identity, identityCurve, keyCurve)
	kvStore := &KVStore{
		Key:   identity + "_" + string(identityCurve) + "_" + string(keyCurve) + "_" + "signers",
		Value: signers,
	}

	_, err := client.Model(kvStore).Insert()
	return err
}

func GetSignersForKeyShare(identity string, identityCurve common.Curve, keyCurve common.Curve) (string, error) {
	var keys []KVStore
	err := client.Model(&keys).Where("key = ?", identity+"_"+string(identityCurve)+"_"+string(keyCurve)+"_signers").Select()

	if err != nil {
		return "", err
	}

	if len(keys) == 0 {
		return "", nil
	}

	return keys[0].Value, nil
}
