# Payment gateway

This is a monorepo for the payment gateway. It contains the following services:

* Payment gateway 
* Bank simulator

# Payment gateway

Service that accepts a request from a merchant, sends it to the bank simulator, and returns the response to the user.
Service store information about payments in the database. Service match bank statuses to the payment gateway statuses.

Project follows DDD approach, has isolated data access layer and stateless business logic. It has bank adapter
to isolate bank interaction from the business logic.

Payment gateway workflow
* Merchant sends a request to the payment gateway
* Payment gateway validates request
* Payment gateway store the payment in the database with *INIT* status
* Payment gateway store card info in encrypted form
* Payment gateway sends a request to the bank
* Payment gateway matches the bank status code to the payment gateway statuses
* Payment gateway updates the payment in the database according to the bank response
* Payment gateway returns the response to the merchant


Service has 2 endpoints:


* POST /payments - create payment
* GET /payments/{id} - get payment by id

Detailed API description is in the [openapi.yaml](./payment_gateway/cmd/web_server/docs/swagger.yaml) file.
Api docs are served by the gateway service on [the swagger endpoint](http://localhost:8080/swagger/index.html#/).

Bank card details are stored seperatly from payment statuses in encrypted form. It's done to avoid storing sensitive data in the database.
See CardRepository for details.

# Bank simulator
Bank simulator is a service that accepts a request from the payment gateway, and returns the response to the payment gateway.
Bank simulator is a mock service that emulates bank behaviour.
It has cards to emulate different bank behaviour.


    4242424242424242 - success
    4000000000000002 - insufficient funds
    4000000000009995 - processing

Bank simulator returns code for different statuses:

    StatusCodeSucceed = "0"
	StatusCodeFailed  = "100"
	StatusCodePending = "200"

Bank simulator api is described in the [deposit.http](./bank_simulator/api/deposit.http) file.


# Building and running the project

Docker is required to run the project.

## Running in docker compose

### 1. Run services
```bash
docker-compose up -d
 ```

check that all services are running

### 2. Validate that the services are running

```bash
docker ps
```

you should see 3 running services

```bash

```

### (first run only) 3. Run migrations for the db

```bash

```

### 4. Call api
Api uses bearer auth, please use test tokens
    
    first_merchant_mock_token
    second_merchant_mock_token

#### a. Use api file (recommended)
Requests are described in the [payment.http](./payment_gateway/api/payment.http) file.
Call it from IDEA http client or use curl


#### b. You can send api requests with similar way

Open swagger docs and call api from there
http://localhost:8080/swagger/index.html#/

Put token in the auth field in format
    
    Bearer <token>
    
    for example:
    Bearer first_merchant_mock_token    

Fill in test data, you can use any test card from bank simulator

```json
{
  "card": {
    "number": "4242424242424242",
    "exp_month": 12,
    "exp_year": 2022,
    "cvv": "123"
  },
  "amount": 100,
  "currency": "USD"
}
```




# Assumptions and Improvements


* ~~Main BL~~
* ~~Enrypt card info~~
* ~~Add swagger docs~~
* ~~Add funcitonal tests for api~~
* ~~Add middleware auth to add merchant information based on their key~~
* Payment gateway supports only api integration, hosted page integration can be added later
* Payment gateway and bank simulator doesn't support 3d secure payments and emulation
* Add configs
* Improve validators in payment gateway
* Add tests that affects the database layer with running db in docker for every test run