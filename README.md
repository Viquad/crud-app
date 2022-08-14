# crud-app

Simple CRUD application for bank's users and their accounts.

## Install

```sh
    git clone https://github.com/Viquad/crud-app.git
```

## Setup 

Set your `POSTGRESS_PASSWORD` to `.env` - environment file. It will be used for docker `backend` and `postgres` containers. 

## Run
```sh
    make run
```
Builds `backend`-container with app. Starts `postgres`-container and `backend`-container. After successful connection `backend`-container to `postgres`-container will applied database migration.

## Stop
```sh
    make stop
```
Stops and removes `backend` and `postgres` containers.

## Clean
```sh
    make clean
```
Use to remove `backend` image.

# REST API

The REST API to the crud app is described below.

## Get list of accounts

### Request

`GET /account`

### Response

```json
[
    {
        "firstName": "Mr",
        "lastName": "Nobody",
        "balance": 666,
        "currency": "USD"
    }
]
```

## Create account

### Request

`POST or PUT /account` 

```json
{
    "firstName": "Mr",
    "lastName": "Nobody",
    "balance": 666,
    "currency": "USD"
}
```

### Response

    TODO

## Get account by id

### Request

`GET /account/:id`

    TODO

### Response

```json
{
    "id": 1,
    "firstName": "Mr",
    "lastName": "Nobody",
    "balance": 666,
    "currency": "USD",
    "lastUpdate": "2022-08-14T13:46:09.236194Z"
}
```

## Edit account by id

### Request

`POST or PUT /account/:id`

```json
{
    "firstName": "Mr",
    "lastName": "Somebody",
    "balance": 1000,
    "currency": "UAH"
}
```

### Response

    TODO

## Delete account by id

### Request

`DELETE /account/:id`

    TODO

### Response

    TODO
    
