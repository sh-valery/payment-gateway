# Payment gateway

This is a monorepo for the payment gateway. It contains the following services:

* Payment gateway
* Bank simulator

# Payment gateway

Payment gateway is a service that accepts a request from a merchant, sends it to the bank, and returns the response to the user.
Service store information about payments in the database. Service match bank statuses to the payment gateway statuses.

Project follows DDD approach, has isolated data access layer and stateless business logic. It has bank adapter
to isolate bank interaction from the business logic. All sensitive information are stored in encrypted form.

## Payment gateway workflow

* Merchant sends a request to the payment gateway
* Payment gateway validates request
* Payment gateway store the payment in the database with *INIT* status
* Payment gateway store card info in encrypted form
* Payment gateway sends a request to the bank
* Payment gateway matches the bank status code to the payment gateway statuses
* Payment gateway updates the payment info in the database according to the bank response
* Payment gateway returns the response to the merchant

## API

Service has 2 endpoints:

* POST /payments - create payment
* GET /payments/{id} - get payment by id

Detailed API description is in the [openapi.yaml](./payment_gateway/cmd/web_server/docs/swagger.yaml) file.
Api docs are served by the gateway service on [the swagger endpoint](http://localhost:8080/swagger/index.html#/).

Bank card details are stored separately from payment statuses in encrypted form. It's done to avoid storing sensitive
data in the database.
See CardRepository for details.

# Bank simulator

Bank simulator is a service that accepts a request from the payment gateway, and returns the response to the payment
gateway.
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

## Run tests
Required GoLang

```bash
cd ./payment_gateway/
go test -cover ./internal/payment/
```

tests covered business logic and bank adaptor, they don't cover
db layer and web server.

## Run services
Required Docker. 
    
*It was tested on MacOS 12.04 and Ubuntu 20.04. Go 1.19 and Docker version 20.10.17*


### 1. Run services
```bash
cd ./payment-gateway
docker-compose up -d
 ```

check that all services are running

### 2. Validate that the services are running

```bash
docker ps
```

you should see 3 running containers

```bash
CONTAINER ID   IMAGE                                       COMMAND                  CREATED         STATUS         PORTS                               NAMES
e3aa7afcffbf   payment-gateway_payment_gateway             "/bin/sh -c '/app/we…"   7 seconds ago   Up 7 seconds   0.0.0.0:8080->8080/tcp              payment-gateway_payment_gateway_1
3a1ede2f9400   mysql:5.7                                   "docker-entrypoint.s…"   5 minutes ago   Up 2 minutes   0.0.0.0:3306->3306/tcp, 33060/tcp   payment-gateway_db_1
99f80d54c9ce   payment-gateway_bank_simulator              "/bin/sh -c '/app/we…"   5 minutes ago   Up 2 minutes   8080/tcp                            payment-gateway_bank_simulator_1
```

### 3. (first run only) Run migrations for the db

```bash
 docker exec -i payment-gateway_db_1 sh -c 'exec mysql -u root -ppass' < ./payment_gateway/db/migration/0001_create_payments_tables.sql
```

check that tables were created
```bash
 docker exec  payment-gateway_db_1 sh -c "mysql -u root -ppass -t -e 'use payment_gateway; show tables;'"
```  
you should see
    
```
+---------------------------+
| Tables_in_payment_gateway |
+---------------------------+
| cards                     |
| payments                  |
+---------------------------+
```
### 4. Call api

Api uses bearer auth, please use test merchant tokens

    first_merchant_mock_token
    second_merchant_mock_token

#### a. Use api file (recommended)

Requests are described in the [payment.http](./payment_gateway/api/payment.http) file.
Call it from IDEA http client or use curl.
You can call check last status after any payment request. It's not necessary any manual actions. 


#### b. You can send api requests from swagger fe

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

#### c. You can use curl

```bash
curl -X POST --location "http://localhost:8080/api/v1/payment" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer first_merchant_mock_token" \
    -d "{
          \"card\": {
            \"number\": \"4242424242424242\",
            \"expiry_month\": \"12\",
            \"expiry_year\": \"2042\",
            \"cvv\": \"123\",
            \"holder_name\": \"John Doe\"
          },
          \"amount\": 1000,
          \"currency\": \"USD\"
        }"
```

# Running in production with k8s
## Run app in the production
It's a short version for experienced users

It requires:
* k8s cluster and docker registry
* Static files served by S3 + CDN or Nginx
* Production database (MySQL, Postgres, etc)

### 1. Setup BE
prepare the image, set the correct registry name in docker-compose.yaml
```shell
    image: {your_registry_user}/bank_simulator
```

```shell
docker login
# user=shvalery
user="{fill in your dockerhub username}"
docker compose push 
```

setup k8s Secrets to store variables and sensitive data.
prepare the file with variables
```shell    
vi secrets

DATABASE_NAME=your_db_name
DATABASE_USERNAME=your_db_user
DATABASE_PASSWORD=your_db_user_password
DATABASE_HOST=your_db_host

```
generate secrets
```shell
kubectl create secret generic payment-gateway-secrets --from-env-file=secrets
```
check results
```shell
kubectl describe secret payment-gateway-secrets
```

deploy to k8s, change shvalery to your registry username in deployment.yaml
```shell

kubectl apply -f ./deployments/bank-simulator-deployment.yaml
kubectl apply -f ./deployments/bank_simulator-service.yaml

kubectl apply -f ./deployments/payment-gateway-deployment.yaml
kubectl apply -f ./deployments/payment_gateway-service.yaml

kubectl apply -f ./deployments/payment-gateway-backend-networkpolicy.yaml
kubectl apply -f ./deployments/log-volume-persistentvolumeclaim.yaml
```

check status
```shell
kubectl get deployments
```

```shell
 kubectl get pods
```


you should see the following output:
```
NAME               READY   UP-TO-DATE   AVAILABLE   AGE
payment-gateway-app   3/3     3            3           38s
```

allow external access
```shell
kubectl apply -f ./payment-gateway-svc.yaml
```

check status
```shell
kubectl get svc payment-gateway
```

you should see the following output:
```
NAME           TYPE       CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
payment-gateway   NodePort   10.101.133.32   <none>        8000:32690/TCP   32s
```
whereas 32690 is external port with api


# Assumptions and Improvements

* ~~Main BL and Tests~~
* ~~Encrypt card info~~
* ~~Add swagger docs~~
* ~~Add functional tests for api~~
* ~~Add middleware auth to add merchant information based on their key~~
* Payment gateway supports only api integration, hosted page integration can be added later
* Payment gateway and bank simulator doesn't support 3d secure payments and emulation
* ~~Add configs~~
* Improve validators in payment gateway
* Add tests that affects the database layer with running db in docker for every test run
* Add metrics and a time series db

# Usefull scripts
 * See make file for swagger gen
 * See interface definition for mock gen
 * Check db payments (btw db port is mapped to 3306 local and available from host)
```bash
docker exec  payment-gateway_db_1 sh -c "mysql -u root -ppass -t -e 'select * from payment_gateway.payments limit 30;'"
```  
```bash
+--------------------------------------+--------------+--------------------------------------+--------------------------------------+--------+----------+-----------+-------------+---------------------+---------------------+
| uuid                                 | merchant_id  | tracking_id                          | card_id                              | amount | currency | status    | status_code | created_at          | updated_at          |
+--------------------------------------+--------------+--------------------------------------+--------------------------------------+--------+----------+-----------+-------------+---------------------+---------------------+
| d09c4ca0-99b2-44d9-b1fa-24e1d9f1bcef | merchant_one | 3254ac93-4ff1-4022-b115-6ac1b57e50e2 | 87599901-6966-4379-84a4-0e255a6d443e |   1000 | USD      | succeeded | 0           | 2023-03-22 12:24:18 | 2023-03-22 12:24:18 |
+--------------------------------------+--------------+--------------------------------------+--------------------------------------+--------+----------+-----------+-------------+---------------------+---------------------+
```
* Check db cards
```bash
  docker exec  payment-gateway_db_1 sh -c "mysql -u root -ppass -t -e 'select * from payment_gateway.cards limit 30;'"
```
```bash
+--------------------------------------+----------------------------------------------+--------------------------+--------------+-------------+------------------------------+---------------------+---------------------+
| uuid                                 | card_number                                  | card_holder              | expiry_month | expiry_year | cvv                          | created_at          | updated_at          |
+--------------------------------------+----------------------------------------------+--------------------------+--------------+-------------+------------------------------+---------------------+---------------------+
| 87599901-6966-4379-84a4-0e255a6d443e | UH2oENhTFfrZ7pmLsFfq91qIOWGvErxXybZNu8QX7ts= | RujnAU8iQAwcd-qv7JQaaA== | 12           | 2042        | HeWhUm0ppuisbXbrXqYRDmoCRg== | 2023-03-22 12:24:18 | 2023-03-22 12:24:18 |
+--------------------------------------+----------------------------------------------+--------------------------+--------------+-------------+------------------------------+---------------------+---------------------+

```