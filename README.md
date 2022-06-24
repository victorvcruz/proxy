# This project is a proxy in golang

To run project you need Golang, [Docker](https://docs.docker.com/engine/install/) and [Docker-compose](https://docs.docker.com/compose/install/),  installed in your pc

## How to run project

1. run `sudo docker-compose up -d` in root directory
2. run `go mod download` in root directory
3. set the ports and hosts in `.env`
4. run `go run main.go` in root directory
5. in http://localhost:8080 insert your request

To stop execution run `sudo docker-compose down`
