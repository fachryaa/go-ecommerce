# REST API Ecommerce Application

This is REST Api for simple e-commerce application.

API Docs : [here](https://www.postman.com/lively-escape-882117/workspace/cbb6a13d-41bc-4dfe-b5b6-5152bbd45512/api/2868ac9b-a983-4417-a904-7e65b778157c/documentation/20020606-c6f9156d-6221-4152-a738-84b0024c3a3c?version=1b1e8d58-b165-4861-8c89-11ab014f3eb1
)

## Technologies
- Golang
- Mux Router
- GORM
- MySql
- Docker

## Install
- Clone the repository
```
git clone https://github.com/fachryaa/go-ecommerce.git
```
- Download all the dependencies
```
go mod download
```

## Run
- Create `.env` file and setup with [the following variables](./.env_example)
- Run the app
```
go build main.go
```

## Run with docker
- Run the app
```
docker-compose up
```

## Features
- Register
- Login
- Logout
- Product
- Cart
- Order