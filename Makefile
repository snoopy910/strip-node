restart:
	docker compose down
	docker volume prune -f
	# ./runNetwork.sh
	docker compose up --force-recreate --build -d

test:
	cd ../strip-sdk && yarn test src/__tests__/ripple.intent.test.ts

go-mod-tidy:
	cd bootnode && go mod tidy
	cd strip-validator && go mod tidy
	go mod tidy