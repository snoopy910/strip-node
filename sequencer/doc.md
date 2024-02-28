# Start Sequencer

First run the sequencer DB:

```
docker run --name sequencer-postgres -p 5434:5432 -e POSTGRES_PASSWORD=password -d postgres
```

Then start the sequencer using:

```
go run main.go --httpPort="8082" -postgresHost="localhost:5434" --isSequencer=true
```

# Sequener APIs

### Create Intent

```
curl --location 'localhost:8082/createIntent' \
--header 'Content-Type: application/json' \
--data '{
    "signature": "0x01",
    "identity": "0x02",
    "identityCurve": "0x03",
    "operations": [{
        "serializedTxn": "0x04",
        "dataToSign": "0x05",
        "chainId": "0x06",
        "keyCurve": "0x07"
    },{
        "serializedTxn": "0x08",
        "dataToSign": "0x09",
        "chainId": "0x10",
        "keyCurve": "0x11"
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