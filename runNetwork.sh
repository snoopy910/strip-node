sudo docker rm -f $(sudo docker ps -a -q)
sudo service postgresql stop
sudo docker compose -f docker-compose.yaml up -d --build