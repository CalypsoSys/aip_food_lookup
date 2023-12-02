#!/bin/bash
export GOPATH=$GOPATH:$(pwd)

echo "aip_food_lookup"
#GOOS=linux GOARCH=amd64 go build -o ./docker/aip_food_lookup ./cmd/aip_food_lookup/*.go
cd ./cmd/aip_food_lookup/
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ../../docker/aip_food_lookup .

cd ../../docker
sudo docker build -t aip_food_lookup .

sudo docker rmi -f $(sudo docker images -f "dangling=true" -q)
sudo docker save aip_food_lookup > aip_food_lookup.tar
mkdir -p /mnt/c/transfer/aip/data
cp aip_food_lookup.tar /mnt/c/transfer/aip/.
cp -r ../cmd/aip_food_lookup/data/ /mnt/c/transfer/aip/

#sudo docker run -p 8080:8080 aip_food_lookup

#curl http://localhost:8080

#1. copy aip_food_lookup.tar /home/calypso/docker
#2. cd /home/calypso/docker/aip_food_lookup_go
#3. sudo docker-compose down
#4. cd /home/calypso/docker/
#5. sudo docker load < aip_food_lookup.tar
#6. sudo docker rmi -f $(sudo docker images -f "dangling=true" -q)
#7. cd /home/calypso/docker/aip_food_lookup_go
#9. sudo docker-compose up -d
