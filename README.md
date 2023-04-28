# Go Microservice Project


There are three microservices which are **Order, User** and **OrderElastic** microservices with using **MongoDB** and **Elasticsearch**.

![alt text](https://i.ibb.co/QfdgZRZ/Order-elastic.jpg)

## Whats Including In This Repository

#### Order microservice
* Web API application 
* REST API principles, CRUD operations
* Repository Pattern Implementation
* **MongoDB** connection and containerization
* Using **Echo Framework**
* Using **Go-playground/Validator** and **Mongo-Driver**
* Using **Custom Response, Middleware and Exceptions** with Shared Library

#### User microservice
* Web API application 
* REST API principles, CRUD operations
* Repository Pattern Implementation
* **MongoDB** connection and containerization
* Using **Echo Framework**
* Using **Go-playground/Validator** and **Mongo-Driver**
* Using **Custom Response, Middleware and Exceptions** with Shared Library

#### OrderElastic microservice
* Fix job application 
* **Elasticsearch** connection and containerization
* Using **Custom Response, Middleware and Exceptions** with Shared Library

#### Asynchronous Communication of Microservices
* Using **Confluent-kafka** for **Kafka** Message-Broker system
* Publishing Order Create-Update-Delete event from Order microservices and Subscribing this message from OrderElastic microservices

#### Docker Compose establishment with on docker
* Containerization of databases
* Override Environment variables

#### Kuberntes establishment with on kubernetes
* Containerization of application with kubernetes

## Run The Project

1. After the cloning this repository run below command at the root directory which include docker-compose.yml files;

```go
docker-compose up -d
```

2. Applications will run with kubernetes so before running you have to do some regulations;
```go
docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' <container_name>
```
  Get MongoDB, Kafka and Elasticsearch ip addresses with this command. Then enter these addresses instead of the ip addresses in the config "development"

3. To create docker image run below command
```go
docker build -t order-user-project/user-api:V01 -f internal/apps/user-api/Dockerfile .
```
```go
docker build -t order-user-project/order-api:V01 -f internal/apps/order-api/Dockerfile .
```
```go
docker build -t order-user-project/order-elastic:V01 -f internal/apps/order-elastic/Dockerfile .
```

4. Run below command to start applications with kubernetes
```go
kubectl apply -f project-deployment.yaml
```

5. You can launch microservices as below urls:
* **Order API -> http://localhost:30001/swagger/index.html**
* **User API -> http://localhost:30002/swagger/index.html**
