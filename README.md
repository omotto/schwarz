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
Create a design diagram (e.g. UML) how your system could interact in a real environment with other services .
#### Task 3
Package your service inside a Docker image so it can be rolled out in a Kubernetes environment.
