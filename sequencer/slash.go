package sequencer

import (
	"context"
	"encoding/hex"
	"errors"
	"time"

	intentoperatorsregistry "github.com/StripChain/strip-node/intentOperatorsRegistry"
	"github.com/StripChain/strip-node/libs"
	db "github.com/StripChain/strip-node/libs/database"
	"github.com/StripChain/strip-node/util/logger"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

func UpdateSignersList() error {
	// NOTE: ganache is not supporting notification and not able to subscribe new events
	instance, err := intentoperatorsregistry.GetIntentOperatorsRegistryContract(RPC_URL, IntentOperatorsRegistryContractAddress)
	if err != nil {
		logger.Sugar().Errorw("Failed to get intent operators registry contract", "error", err)
		return err
	}
	filterOpts := &bind.FilterOpts{Context: context.Background()}
	itr, err := instance.FilterSignerUpdated(filterOpts, nil)
	if err != nil {
		logger.Sugar().Errorw("Failed to filter old signerUpdated events", "error", err)
		return err
	}
	heartbeats, err := db.GetHeartbeats()
	if err != nil {
		logger.Sugar().Errorw("Failed to get heartbeats for already registered signer", "error", err)
		return err
	}
	for itr.Next() {
		event := itr.Event
		if event.Added {
			if len(heartbeats) == 0 {
				logger.Sugar().Infow("Add heartbeat for newly registered signer", "publickey", "0x"+hex.EncodeToString(event.Publickey[:]), "url", event.Url)
				if err := db.AddHeartbeat("0x" + hex.EncodeToString(event.Publickey[:])); err != nil {
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
					if err := db.AddHeartbeat("0x" + hex.EncodeToString(event.Publickey[:])); err != nil {
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
		serviceClient, err := validatorClientManager.FindOrCreateClient(signer.URL)
		if err != nil {
			logger.Sugar().Errorf("Failed to get managed client for signer %s: %v", signer.URL, err)
			return
		}

		conn := serviceClient.GetClientConn()
		healthClient := grpc_health_v1.NewHealthClient(conn)
		healthCheckRequest := &grpc_health_v1.HealthCheckRequest{
			Service: "",
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		resp, err := healthClient.Check(ctx, healthCheckRequest)
		if err != nil {
			statusErr, ok := status.FromError(err)
			target := conn.Target()
			if ok {
				logger.Sugar().Errorf("gRPC health check call failed for target %s: Code=%s, Msg=%s", target, statusErr.Code(), statusErr.Message())
			} else {
				if errors.Is(err, context.DeadlineExceeded) {
					logger.Sugar().Warnf("gRPC health check timed out for target %s after 5 seconds", target)
				} else {
					logger.Sugar().Errorf("gRPC health check call failed for target %s: %v", target, err)
				}
			}
			return
		}
		if resp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
			logger.Sugar().Warnf("Signer %s is not serving (status: %s)", signer.URL, resp.Status)
			return
		}

		err = db.UpdateHeartbeat(signer.PublicKey)
		if err != nil {
			logger.Sugar().Errorw("Failed to update heartbeat for signer", "signer", signer, "publicKey", signer.PublicKey, "error", err)
			return
		} else {
			logger.Sugar().Infow("Updated heartbeat successfully", "signerURL", signer.URL, "publicKey", signer.PublicKey)
		}
	}
}

func startCheckingSigner() {
	go func() {
		for {
			UpdateSignersList()
			CheckSignersStatus()

			activeSigners, err := db.GetActiveSigners()
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

			time.Sleep(libs.SLASHING_CHECK_INTERVAL)
		}
	}()
}
