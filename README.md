# Silent TSS CLI

## Run using Shell

Here is the steps to run the CLI:

### Build CLI

```
go build
```

### Generate Account with Two Signers

First run ganache to simulate a local ethereum network:

```
ganache --port 8545 -h=0.0.0.0 -m="rifle cloud amused end pyramid swarm anxiety kitchen ceiling cotton rib gain" 
```

Then run postgres:

```
docker run --name node1-postgres -p 5432:5432 -e POSTGRES_PASSWORD=password -d postgres
docker run --name node2-postgres -p 5433:5432 -e POSTGRES_PASSWORD=password -d postgres
```

Then deploy the IntentOperatorsRegistry contract:

```
go run main.go --isDeployIntentOperatorsRegistry=true --privateKey="76163f58a83febacfdef93e2142591d4d676432fa6c737ce1dd90a51083c461a"
```

Then add two signers to the IntentOperatorsRegistry contract:

```
go run main.go --isAddsigner=true --signerPublicKey="0x26d1556a83c01a9d2b1cce29b32cb520238efc602f86481d2d0b9af8a2fff0cf" --privateKey="76163f58a83febacfdef93e2142591d4d676432fa6c737ce1dd90a51083c461a" --signerNodeURL="http://localhost:8080"

go run main.go --isAddsigner=true --signerPublicKey="0x54455a1f7f4244ef645ac62baa8bd90af0cc18cdb0eae369766b7b58134edf35" --privateKey="76163f58a83febacfdef93e2142591d4d676432fa6c737ce1dd90a51083c461a" --signerNodeURL="http://localhost:8081"
```

Then run the bootnode:

```
go run main.go --port=30303 --isBootstrap=true --keyPath="./static-bootnode"
```

And then run the following commands to spin two signers:

```
go run main.go --signerPublicKey="0x26d1556a83c01a9d2b1cce29b32cb520238efc602f86481d2d0b9af8a2fff0cf" --signerPrivateKey="0xb0a0aa1369839ffbf2778fcedcad2ba70b0237e6071b791a80a6f9e11380ffa2" --bootnode="/ip4/0.0.0.0/tcp/30303/p2p/QmTfM73oQxzx6DVyjCm5AECW3hVbXJiSLYtosNauaX9gJR" --httpPort="8080" --port=30304 --keyPath="./keys/node1"

go run main.go --signerPublicKey="0x54455a1f7f4244ef645ac62baa8bd90af0cc18cdb0eae369766b7b58134edf35" --signerPrivateKey="0x4d539b1896a8f7064a7207fa005b13b64f90eff78564e278c14b1089d2d5f8de" --bootnode="/ip4/0.0.0.0/tcp/30303/p2p/QmTfM73oQxzx6DVyjCm5AECW3hVbXJiSLYtosNauaX9gJR" --httpPort="8081" --port=30305 --keyPath="./keys/node2" --postgresHost="localhost:5433"
```

Then finally start keygen generation round using this command:

```
curl --location 'http://localhost:8080/keygen' \
--header 'Content-Type: application/json' \
--data '{
    "identity": "0x40310390beF6D518f62Cc41a60a4E9b4a8b98730",
    "identityCurve":"ecdsa",
    "keyCurve":"eddsa",
    "signers": [
        "0x26d1556a83c01a9d2b1cce29b32cb520238efc602f86481d2d0b9af8a2fff0cf",
        "0x54455a1f7f4244ef645ac62baa8bd90af0cc18cdb0eae369766b7b58134edf35"
    ]
}'

curl --location 'http://localhost:8080/keygen' \
--header 'Content-Type: application/json' \
--data '{
    "identity": "0x40310390beF6D518f62Cc41a60a4E9b4a8b98730",
    "identityCurve":"ecdsa",
    "keyCurve":"ecdsa",
    "signers": [
        "0x26d1556a83c01a9d2b1cce29b32cb520238efc602f86481d2d0b9af8a2fff0cf",
        "0x54455a1f7f4244ef645ac62baa8bd90af0cc18cdb0eae369766b7b58134edf35"
    ]
}'
```

### Sign Hash

```
curl http://localhost:8080/signature?message=87PYr1rWKA1AWiAcp27tjJbsW6fFWD6FkfXm57ior88XhFuEdwKrZWDnzMvNvvjpUX26JnAJZ2Znoa9nfzSAnc5fNuiqHVjqufDqberL4WW11eRHXB28iUZnbYvodhZhYRXodx6tXbSy6QLcmXdRWgG6EVt8K8i6qbiEmAwKnyG1wfyYdSLoWD1VLV2R7hc92rACtt6NE7Zq&identity=0x40310390beF6D518f62Cc41a60a4E9b4a8b98730&identityCurve=ecdsa&keyCurve=eddsa

curl http://localhost:8080/signature?message=97250d83d64402e2a545ec59594743b6bf8e132395de3074392bbf34987bf675&identity=0x40310390beF6D518f62Cc41a60a4E9b4a8b98730&identityCurve=ecdsa&keyCurve=ecdsa
```

### Get Address

```
curl http://localhost:8080/address?identity=0x40310390beF6D518f62Cc41a60a4E9b4a8b98730&identityCurve=ecdsa&keyCurve=eddsa

curl http://localhost:8080/address?identity=0x40310390beF6D518f62Cc41a60a4E9b4a8b98730&identityCurve=ecdsa&keyCurve=ecdsa
```

### E2E Tests

To run e2e tests use the following command:

```
./e2e_test.sh
```