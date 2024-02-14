docker-compose -f docker-compose.yaml up -d --build

sleep 5

# generate keygen
echo "Generating keygen"
curl  --location 'http://localhost:8080/keygen?networkId=1337'
sleep 120
echo "Completed generating keygen"

if curl --location 'http://localhost:8080/signature?hash=1c8aff950685c2ed4bc3174f3472287b56d9517b9c948127319a09a7a36deac8&networkId=1337'; then
  docker-compose -f docker-compose.yaml down
  echo "e2e tests passed"
else
  docker-compose -f docker-compose.yaml down
  echo "e2e tests failed"
  exit;
fi;