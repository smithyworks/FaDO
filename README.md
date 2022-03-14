# Master Thesis

## Data-Aware Function Scheduling on a Multi-Serverless Platform

### Author: Christopher Smith

### Advisor: M.Sc. Anshul Jindal

### Supervisor: Prof. Dr. Michael Gerndt

### School: Technical University of Munich, Department of Informatics

## Abstract

Function-as-a-Service (FaaS) is an attractive cloud computing model that simplifies application development and deployment. However, current FaaS technologies do not consider data placement when scheduling tasks. With the growing demand for multi-cloud, multi-serverless applications, this flaw means serverless technologies are still ill-suited to latency-sensitive operations like media streaming.

This thesis proposes a solution by presenting **FaDO**, the *Function and Data Orchestrator*, which is a proof-of-concept application designed to allow data-aware function scheduling on a multi-serverless platform.

The application comprises a back-end server and API, along with a high-performance load balancer, a database, and a frontend browser client. These components allow users to interact with the application easily and seamlessly schedule functions onto the multi-serverless platform according to their data requirements.

FaDO further provides users with an abstraction of the platform's storage, allowing users to interact with data across different storage services through a unified interface. In addition, users can configure automatic and policy-aware granular data replications, causing the application to spread data across the platform while respecting location constraints.

The implementation thus enables users to distribute functions across a heterogeneous platform through data replication, balancing location constraints and performance requirements, and optimizing throughput using different load balancing policies.

The application fulfills its requirements, and load testing results show that it is capable of load balancing high-throughput workloads, placing tasks near their data without contributing any significant performance overhead. A qualitative evaluation of the system's design further indicates that FaDO has the ingredients necessary to make a reliable and performant network application.

## Structure of the Repository

Directories:

- `documents`: Documents submitted for the thesis.
- `database`: SQL files defining the data model and initial data (PostgreSQL).
- `server`: FaDO's backend server code (Golang).
- `client`: FaDO's frontend browser client (React.js).
- `mock-faas`: A mock FaaS server used for development (Node.js).
- `bin`: Convenience scripts to run FaDO's development environment.
- `deployment`: Files related to FaDO's local deployment.

## Running the Application

FaDO is meant to run on a complex cloud platform containing MinIO storage deployments and FaaS endpoints.

This repository provides a development environment that can be run locally using Docker Compose at the root of the project. For instance:

```
# To bring the application up and run it in the background.
$ docker compose up -d

# To bring the application down and remove locally build Docker images.
$ docker compose down -v --rmi local
```

The environment is composed of:
- 3 MinIO deployments that are each made up of 3 MinIO servers,
- An NGINX servers used to interface with the 3 MinIO deployments.
- 3 mock FaaS servers that can respond to load-balanced requests.
- The PostgreSQL database containing the application state.
- The Caddy load balancer that serves to distribute function invocations to the FaaS services.
- The FaDO backend server which orchestrates the different resources, exposes an API, and serves the frontend browser client.

The environment can also be spun up piecewise using the scripts inside the `bin` folder:
- `bin/dc-minio.sh` is a Docker Compose stub to run the MinIO deployments on ports 9010, 9020, and 9030, with the management consoles on ports 9011, 9021, and 9031 (e.g. to start them, use `$ bin/dc-minio.sh up`).
- `bin/dc-faas.sh` is a Docker Compose stub to run the mock FaaS servers on ports 9101, 9102, and 9103.
- `bin/run-db.sh` runs FaDO's database on port 5454.
- `bin/run-caddy.sh` runs FaDO's Caddy load balancer on port 6000 with the administration endpoint on port 2019.
- `bin/run-server.sh`\* runs FaDO's backend server with the API on port 9090.
- `bin/run-client.sh` runs FaDO's frontend client on port 3000.

\* The `bin/run-server.sh` script runs the backend server on the host machine, and depends on the presence of the Go language tools and MinIO's `mc` command line tool.

## Useful Resources

- [MinIO](https://min.io/)
- [PostgreSQL](https://www.postgresql.org/)
- [Caddy Server](https://caddyserver.com/)
- [Go](https://go.dev/)
- [Node.js](https://nodejs.org/en/)
- [React.js](https://reactjs.org/)
- [Docker](https://www.docker.com/)
