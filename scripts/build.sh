#!/bin/bash
export GOPATH=$GOPATH:$(pwd)

echo "aip_food_lookup"
#GOOS=linux GOARCH=amd64 go build -o ./docker/aip_food_lookup ./cmd/aip_food_lookup/*.go
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./docker/aip_food_lookup ./cmd/aip_food_lookup/*.go

cd ./docker
sudo docker build -t aip_food_lookup .

#sudo docker run -p 8080:8080 aip_food_lookup

#curl http://localhost:8080

