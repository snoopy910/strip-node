restart:
	docker compose down
	docker volume prune -f
	./runNetwork.sh

test:
	cd ../strip-sdk && yarn test src/__tests__/cardano.intent.test.ts