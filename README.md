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

#### Docker Compose establishment with all microservices on docker
* Containerization of microservices
* Containerization of databases
* Override Environment variables

## Run The Project

1. After the cloning this repository run below command at the root directory which include docker-compose.yml files;

```go
docker-compose up -d
```

2. You can launch microservices as below urls:
* **Order API -> http://host.docker.internal:8011/swagger/index.html**
* **User API -> http://host.docker.internal:8012/swagger/index.html**
