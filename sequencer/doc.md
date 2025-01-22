# Start Sequencer

First run the sequencer DB:

```
docker run --name sequencer-postgres -p 5434:5432 -e POSTGRES_PASSWORD=password -d postgres
```

Then start the sequencer using:

```
go run main.go --httpPort="8082" -postgresHost="localhost:5434" --isSequencer=true
```

# Status

Intent Status changes in the following way:

```
Processing -> Failed -> Completed
```

Operations Status changes in the following way:

```
Pending -> Waiting -> Failed or Completed
```

Whenever an intent is created a goroutine is created to handle the intent. It will process the intent operations one by one. In case process is exited while an intent is ongoing then after process is started all the process intents are handled. 

The goroutine that handles the intents, will loop through the operations and handle them one by one.

# Sequener APIs

### Create Wallet

```
curl http://localhost:8082/createWallet?identity=0x40310390beF6D518f62Cc41a60a4E9b4a8b98730&identityCurve=ecdsa
```

### Get Wallet

```
curl http://localhost:8082/getWallet?identity=0x40310390beF6D518f62Cc41a60a4E9b4a8b98730&identityCurve=ecdsa
```

### Create Intent

```
curl --location 'localhost:8082/createIntent' \
--header 'Content-Type: application/json' \
--data '{
    "signature": "0x01",
    "identity": "0x40310390beF6D518f62Cc41a60a4E9b4a8b98730",
    "identityCurve": "ecdsa",
    "operations": [{
        "serializedTxn": "eb808477359400825208941c3bb7529e6a484db8f6d7f30e8e840c68dcd13788016345785d8a000080808080",
        "dataToSign": "cdf74de9e6b07952da2e4c6f38e14c1536abefeb9e1f37240e9392e195785c20",
        "chainId": "1337",
        "keyCurve": "ecdsa"
    },{
        "serializedTxn": "187PYuMkJ1sTeui7KoDted7fTc7BimkvmbK5VR4f6vwCBdotDjKZbMtmNn976iWqPoqFesmMEKDdj3SqmgWvQQhD4wsuov477Wd6r2yjinZ3BiBKupWYTvi2U9iqq6eRTUnKKdJV7rkywHeHimYR1jtAqzjsrjM9gt6aWppRjjULx94PCYqZbkEcbp3xeyHT7aBmPMAU2xRg3",
        "dataToSign": "187PYuMkJ1sTeui7KoDted7fTc7BimkvmbK5VR4f6vwCBdotDjKZbMtmNn976iWqPoqFesmMEKDdj3SqmgWvQQhD4wsuov477Wd6r2yjinZ3BiBKupWYTvi2U9iqq6eRTUnKKdJV7rkywHeHimYR1jtAqzjsrjM9gt6aWppRjjULx94PCYqZbkEcbp3xeyHT7aBmPMAU2xRg3",
        "chainId": "901",
        "keyCurve": "eddsa"
    }]
}'
```

### Get Intent by Id

```
curl --location 'localhost:8082/getIntent?id=1'
```

### Get Intents by Status

```
curl --location 'localhost:8082/getIntents?status=processing'
```

# Google oauth2 login


### Setup


### Endpoints

`/oauth/login`: Login with Google

`/oauth/callback`: Callback from Google - checks if the user is authenticated and have a valid access token otherwise it will authenticate the user

`/oauth/verifySignature`: Verifies the signature of the predefined message signed by the user's identity wallet and creates a new id,access and refresh tokens including the identity and identity curve.

`/oauth/accessToken`: Get a new access token when the refresh token is still valid.


### Environment Variables
