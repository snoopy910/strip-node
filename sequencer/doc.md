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

### Create Intent

```
curl --location 'localhost:8082/createIntent' \
--header 'Content-Type: application/json' \
--data '{
    "signature": "0x01",
    "identity": "0x02",
    "identityCurve": "0x03",
    "operations": [{
        "serializedTxn": "0xeb808477359400825208941c3bb7529e6a484db8f6d7f30e8e840c68dcd13788016345785d8a000080808080",
        "dataToSign": "0xcdf74de9e6b07952da2e4c6f38e14c1536abefeb9e1f37240e9392e195785c20",
        "chainId": "1337",
        "keyCurve": "ecdsa"
    },{
        "serializedTxn": "0xeb018477359400825208941c3bb7529e6a484db8f6d7f30e8e840c68dcd13788016345785d8a000080808080",
        "dataToSign": "0x341d703a42c3a62818d5ec709c2935411cbf72836c9277fcdfbfdc216d212a6d",
        "chainId": "1337",
        "keyCurve": "ecdsa"
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

