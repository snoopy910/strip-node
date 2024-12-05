# Silent TSS CLI

To create a network with bootnode, 2 SIOs and sequencer 

```
docker build -t strip-contracts .
./runNetwork.sh
```

### Test

To run the testDB for unit testing, 
```
docker compose -f docker-compose.test.yml up -d
```