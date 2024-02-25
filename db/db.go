package db

import (
	"context"
	"fmt"

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

func AddKeyShare(identity string, identityCurve string, keyCurve string, key string) error {
	fmt.Println("Adding key share to redis", identity+"_"+identityCurve+"_"+keyCurve)
	err := client.Set(context.Background(), identity+"_"+identityCurve+"_"+keyCurve, key, 0).Err()
	return err
}

func GetKeyShare(identity string, identityCurve string, keyCurve string) (string, error) {
	fmt.Println("Getting key share from redis", identity+"_"+identityCurve+"_"+keyCurve)
	val, err := client.Get(context.Background(), identity+"_"+identityCurve+"_"+keyCurve).Result()
	return val, err
}

func AddSignersForKeyShare(identity string, identityCurve string, keyCurve string, signers string) error {
	err := client.Set(context.Background(), identity+"_"+identityCurve+"_"+keyCurve+"_"+"signers", signers, 0).Err()
	return err
}

func GetSignersForKeyShare(identity string, identityCurve string, keyCurve string) (string, error) {
	val, err := client.Get(context.Background(), identity+"_"+identityCurve+"_"+keyCurve+"_"+"signers").Result()
	return val, err
}
