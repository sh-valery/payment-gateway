# Payment gateway

This is a monorepo for the payment gateway. It contains the following services:

* Payment gateway 
* Bank simulator

# Payment gateway

Service that accepts a request from a user, sends it to the bank simulator, and returns the response to the user.
Service store information about payments in the database. Service match bank statuses to the payment gateway statuses.

Project follows DDD approach, has isolated data access layer and stateless business logic. It has bank adapter
to isolate bank interaction from the business logic.



Service has 2 endpoints:


* POST /payments - create payment
* GET /payments/{id} - get payment by id

Detailed API description is in the [openapi.yaml](./payment_gateway/cmd/web_server/docs/swagger.yaml) file.
It served by the gateway service on the /openapi.yaml endpoint.


# Bank simulator
Bank simulator is a service that accepts a request from the payment gateway, and returns the response to the payment gateway.
Bank simulator is a mock service that emulates bank behaviour.
It has cards to emulate different bank behaviour.


    4242424242424242 - success
    4000000000000002 - insufficient funds
    4000000000009995 - processing error

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

you should see something like this

```bash

```

### (first run only) 3. Run migrations for the db

```bash

```

### 4. Call api

#### a. You can send api requests with similar way

Open swagger docs and call api from there

http://localhost:8080/swagger/index.html#/

```bash 


```
#### b. Use api file
Requests are described in the [payment.http](./payment_gateway/api/payment.http) file.