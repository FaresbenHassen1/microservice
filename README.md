Description:
A microservice project with the objective of learning how to communicate a GO project with Database.
the project consist of a small microservice that communicate with SQL database which is in Postgres that is contained within a docker container and also exposes some rest API.

Architecture:
the next figure will present the general architecture of the project with its key figures:
![alt text](/docs/architecture.png)

Entity Relationship Diagram (ERD):
The ERD consist of 3 tables which are Wallet,Transaction and user.
![alt Entity Relationship Diagram](/docs/erd.png)

CMD:
to run REST api:
go run main.go

to run GRPC :

- server: go run grpc_server/main.go
- client: go run grpc_client/main.go
