package db

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var client *redis.Client

func Initialise(host string, db int, username string, password string) {
	client = redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       db,
		Username: username,
	})
}

func AddKeyShare(identity string, identityCurve string, key string, keyCurve string) error {
	err := client.Set(context.Background(), identity+"_"+identityCurve+"_"+keyCurve, key, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetKeyShare(identity string, identityCurve string, keyCurve string) (string, error) {
	val, err := client.Get(context.Background(), identity+"_"+identityCurve+"_"+keyCurve).Result()
	if err != nil {
		return "", err
	}

	return val, nil
}
