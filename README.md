# Schwarz
### Development Challenge

We need a service which can create, update and delete a Kubernetes resource which is able to bootstrap a simple Postgres database.
The service itself should provide a gRPC interface.
Also, the service should provide metrics and health checks to ensure that everything is working and can be monitored.
The key features of this service are:

- Kubernetes resource generation - secure authentication
- code generation
- metrics
- health-check

#### Task 1
Implement the service in Go.
#### Task 2
Create a design diagram (e.g. UML) how your system could interact in a real environment with other services.
#### Task 3
Package your service inside a Docker image, so it can be rolled out in a Kubernetes environment.

## Instructions
- First of all, it's necessary to set the right *kubeconfig* file as **config** to run locally or **docker_config** to run inside docker container
- On the other hand is necessary to set the proper environmental variables that are defined inside **.env** file to run it locally
```
set -a
. ./.env
set +a
```
- To set up the project to work on it:
```
make setup-local
```
- The following instruction does: **generate rpc files -> linter -> formatted code -> compile -> unitary test** 
```
make all
```
- Executable file is located in **.\bin** folder
- To run dockerized project:
```
make docker
```

## Design Diagram
![Alt text](doc/diagram.png?raw=true "Diagram")
