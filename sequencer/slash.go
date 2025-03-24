package sequencer

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	intentoperatorsregistry "github.com/StripChain/strip-node/intentOperatorsRegistry"
	"github.com/StripChain/strip-node/util/logger"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

const (
	CHECK_INTERVAL    = 5 * time.Minute
	HEARTBEAT_TIMEOUT = 2 * time.Hour
)

func UpdateSignersList() error {
	// NOTE: ganache is not supporting notification and not able to subscribe new events
	instance := intentoperatorsregistry.GetIntentOperatorsRegistryContract(RPC_URL, IntentOperatorsRegistryContractAddress)
	filterOpts := &bind.FilterOpts{Context: context.Background()}
	itr, err := instance.FilterSignerUpdated(filterOpts, nil)
	if err != nil {
		logger.Sugar().Errorw("Failed to filter old signerUpdated events", "error", err)
		return err
	}
	heartbeats, err := GetHeartbeats()
	if err != nil {
		logger.Sugar().Errorw("Failed to get heartbeats for already registered signer", "error", err)
		return err
	}
	for itr.Next() {
		event := itr.Event
		if event.Added {
			if len(heartbeats) == 0 {
				logger.Sugar().Infow("Add heartbeat for newly registered signer", "publickey", "0x"+hex.EncodeToString(event.Publickey[:]), "url", event.Url)
				if err := AddHeartbeat("0x" + hex.EncodeToString(event.Publickey[:])); err != nil {
					logger.Sugar().Errorw("Failed to add heartbeat for already registered signer", "error", err)
				}
			} else {
				exist := false
				for _, h := range heartbeats {
					if hex.EncodeToString(event.Publickey[:]) == h.PublicKey {
						exist = true
						break
					}
				}
				if !exist {
					logger.Sugar().Infow("Add heartbeat for newly registered signer", "publickey", "0x"+hex.EncodeToString(event.Publickey[:]), "url", event.Url)
					if err := AddHeartbeat("0x" + hex.EncodeToString(event.Publickey[:])); err != nil {
						logger.Sugar().Errorw("Failed to add heartbeat for already registered signer", "error", err)
					}
				}
			}
		}
	}
	return nil
}

func CheckSignersStatus() {
	signers, err := SignersList()
	if err != nil {
		logger.Sugar().Errorw("Failed to get signers", "error", err)
		return
	}

	for _, signer := range signers {
		logger.Sugar().Infow("Checking signer", "url", signer.URL)
		url := fmt.Sprintf("%s/health", signer.URL)
		resp, err := http.Get(url)
		if err != nil {
			logger.Sugar().Errorf("Signer %s is not working well, error: %v", signer, err)
			return
		}
		if resp != nil && resp.StatusCode == http.StatusOK {
			err := UpdateHeartbeat(signer.PublicKey)
			if err != nil {
				logger.Sugar().Errorf("Failed to update heartbeat for signer %s: %v", signer, err)
				return
			} else {
				logger.Sugar().Infow("Updated heartbeat successfully", "heartbeat", signer.PublicKey, "timestamp", time.Now())
			}
		}
		defer resp.Body.Close()
	}
}

func startCheckingSigner() {
	go func() {
		for {
			UpdateSignersList()
			CheckSignersStatus()

			activeSigners, err := GetActiveSigners()
			if err != nil {
				logger.Sugar().Errorw("Failed to get active signers", "error", err)
				return
			}

			activeSignersMap := make(map[string]bool)
			for _, signer := range activeSigners {
				activeSignersMap[signer.PublicKey] = true
			}

			// instance := intentoperatorsregistry.GetIntentOperatorsRegistryContract(RPC_URL, IntentOperatorsRegistryContractAddress)
			signers, err := SignersList()
			if err != nil {
				logger.Sugar().Errorw("Failed to get all signers", "error", err)
				return
			}

			for _, signer := range signers {
				if !activeSignersMap[signer.PublicKey] {
					//TODO: Call the contract to apply detailed slashing logic
					logger.Sugar().Infow("Slash signer", "signer", signer)
				} else {
					logger.Sugar().Infow("Signer is acting properly", "signer", signer.URL)
				}
			}

			time.Sleep(CHECK_INTERVAL)
		}
	}()
}
