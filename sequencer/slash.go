package sequencer

import (
	"fmt"
	"net/http"
	"time"

	intentoperatorsregistry "github.com/StripChain/strip-node/intentOperatorsRegistry"
	"github.com/StripChain/strip-node/util/logger"
	"github.com/go-pg/pg/v10"
)

var db *pg.DB

type Heartbeat struct {
	PublicKey string    `pg:"publickey,pk"`
	Timestamp time.Time `pg:"timestamp"`
}

const (
	CHECK_INTERVAL    = 5 * time.Minute
	HEARTBEAT_TIMEOUT = 5 * time.Minute
)

func AddHeartbeat(publicKey string) error {
	heartbeat := &Heartbeat{
		PublicKey: publicKey,
		Timestamp: time.Now(),
	}
	_, err := db.Model(heartbeat).Insert()
	if err != nil {
		return err
	}
	return nil
}

func UpdateHeartbeat(publicKey string) error {
	heartbeat := &Heartbeat{
		PublicKey: publicKey,
		Timestamp: time.Now(),
	}
	_, err := db.Model(heartbeat).
		Set("timestamp = ?timestamp").
		Where("publickey = ?publickey").
		Update()
	if err != nil {
		return err
	}
	return nil
}

func GetHeartbeat(publicKey string) (Heartbeat, error) {
	heartbeat := &Heartbeat{
		PublicKey: publicKey,
	}
	err := db.Model(heartbeat).
		Where("publickey = ?publickey").
		Select()
	if err != nil {
		return Heartbeat{}, err
	}
	return *heartbeat, nil
}

func DeleteHeartbeat(publicKey string) error {
	heartbeat := &Heartbeat{
		PublicKey: publicKey,
	}
	_, err := db.Model(heartbeat).Delete()
	if err != nil {
		return err
	}
	return nil
}

func OnSignerRegistered(publicKey string) error {
	// TODO: retrieve event when new signer is registered
	signers, err := SignersList()
	if err != nil {
		logger.Sugar().Errorw("Failed to get active signers", "error", err)
	}
	newSigners := []string{}
	for _, signer := range signers {
		if _, err := GetHeartbeat(signer.PublicKey); err == nil {
			newSigners = append(newSigners, signer.PublicKey)
		}
	}
	for _, publicKey := range newSigners {
		if err := AddHeartbeat(publicKey); err != nil {
			logger.Sugar().Errorw("Failed to add heartbeat", "error", err)
		}
	}
	return nil
}

func IsSignerAlive(publicKey string) bool {
	heartbeat := &Heartbeat{
		PublicKey: publicKey,
	}
	err := db.Model(heartbeat).Last()
	if err != nil {
		return false
	}
	if time.Since(heartbeat.Timestamp) > HEARTBEAT_TIMEOUT {
		return false
	}
	return true
}

func GetActiveSigners() ([]Signer, error) {
	var signers []Signer
	err := db.Model((*Heartbeat)(nil)).
		ColumnExpr("distinct publickey").
		Where("timestamp > ?", time.Now().Add(-HEARTBEAT_TIMEOUT)).
		Select(&signers)
	if err != nil {
		return nil, err
	}
	return signers, nil
}

func CheckSignersStatus() {
	signers, err := SignersList()
	if err != nil {
		logger.Sugar().Errorw("Failed to get signers", "error", err)
		return
	}

	for _, signer := range signers {
		go func(signer Signer) {
			url := fmt.Sprintf("http://%s/health", signer.URL)
			resp, err := http.Get(url)
			if err != nil {
				logger.Sugar().Errorf("Failed to update heartbeat for signer %s: %v", signer, err)
				return
			}
			if resp != nil && resp.StatusCode == http.StatusOK {
				UpdateHeartbeat(signer.PublicKey)
			}
			defer resp.Body.Close()
		}(signer)
	}
}

func startCheckingSigner() {
	go func() {
		time.Sleep(CHECK_INTERVAL)

		CheckSignersStatus()

		signers, err := GetActiveSigners()
		if err != nil {
			logger.Sugar().Errorw("Failed to get active signers", "error", err)
			return
		}
		activeSigners := make(map[string]bool)
		for _, signer := range signers {
			activeSigners[signer.PublicKey] = true
		}

		instance := intentoperatorsregistry.GetIntentOperatorsRegistryContract(RPC_URL, IntentOperatorsRegistryContractAddress)
		signers, err = SignersList()
		if err != nil {
			logger.Sugar().Errorw("Failed to get all signers", "error", err)
			return
		}

		for _, signer := range signers {
			if !activeSigners[signer.PublicKey] {
				//TODO: Call the contract to apply detailed slashing logic
				logger.Sugar().Infow("Slash signer", "signer", signer)
			}
		}
	}()
}
