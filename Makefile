restart:
	docker compose down -v
	docker compose up --force-recreate --build -d

test:
	cd ../strip-sdk && yarn test src/__tests__/ripple.intent.test.ts
	# cd ../strip-sdk && yarn test src/__tests__/solana.intent.test.ts
	cd ../strip-sdk && yarn test src/__tests__/evm/user_transaction.test.ts
	# cd ../strip-sdk && yarn test src/__tests__/cardano.intent.test.ts
	# cd ../strip-sdk && yarn test src/__tests__/dogecoin.intent.test.ts
	# cd ../strip-sdk && yarn test src/__tests__/bitcoin.intent.test.ts
	# cd ../strip-sdk && yarn test src/__tests__/aptos.intent.test.ts
	cd ../strip-sdk && yarn test src/__tests__/stellar.intent.test.ts
	cd ../strip-sdk && yarn test src/__tests__/sui.intent.test.ts
	# cd ../strip-sdk && yarn test src/__tests__/algorand.intent.test.ts
	

go-mod-tidy:
	cd bootnode && go mod tidy
	cd strip-validator && go mod tidy
	go mod tidy

go-update:
	cd bootnode && go get -u ./...
	cd strip-validator && go get -u ./...
	go get -u ./...

gen-bootnode-key:
	cd bootnode && go run bootnode.go

gen-proto:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	--experimental_allow_proto3_optional \
    libs/proto/validator.proto