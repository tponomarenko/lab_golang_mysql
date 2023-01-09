# Lab – Phonebook service in Go

The goal of the lab is to build and run a simple phonebook service
that requires a relational database to store phonebook records. Another
goal is to practice interaction with a service that provides a REST API
over the HTTP protocol.

The service provides a typical **CRUD** API for managing a phone book. CRUD
stands for **C**reate, **R**ead, **U**pdate, **D**elete. APIs like this
one are common in production environments.


## Table of contents

- Tasks:
  - [Task 1 – Build the service](#task-1--build-the-service)
  - [Task 2 – Start the service without a database](#task-2--start-the-service-without-a-database)
  - [Task 3 – Build a docker image for the service](#task-3--build-a-docker-image-for-the-service)
  - [Task 4 – Prepare MySQL server](#task-4--prepare-mysql-server)
  - [Task 5 – Configure database connection for the service](#task-5--configure-database-connection-for-the-service)
  - [Task 6 – Different methods of running HTTP queries](#task-6--different-methods-of-running-http-queries)
  - [Task 7 – Authentication](#task-7--authentication)
- Documentation
  - [Service configuration](#service-configuration)
  - [HTTP endpoints](#http-endpoints)
  - [Record format](#record-format)
  - [Authentication](#authentication)


## Task 1 – Build the service

The service is written in Go thus requires to be built with Go compiler
to create an executable binary file (program).

1. Install Go compiler. Ubuntu package name is `golang`.
2. Compile the source code with the following command:

   ```go build -o /path/to/output /path/to/source/code/directory```

3. Make sure that the source code compiles without errors.
4. Ensure the produced file specified `/path/to/output` could be executed.


## Task 2 – Start the service without a database

1. Execute the file, produced after Task 1.
2. Did the service start. If no - what error did you see?
3. Refer to the [Service configuration](#service-configuration) section to configure the
   port  and to fix the problem. Please note, that only ports **80**
   and **443** are open on the server.
4. Run the service again. Did it start?
5. Explore some of service endpoints. Please refer to the
   [HTTP endpoints](#http-endpoints) section.


## Task 3 – Build a docker image for the service

1. Use Docker multistage build to compile service in one container
   and then copy the built file into a small and clean image.
2. When compiling the code, please set `CGO_ENABLED` environment
   variable to `0`.
3. Run the container and test the service the same way as in Task 2.
4. Create a docker compose file to start the container using docker-compose.

## Task 4 – Prepare MySQL server

1. Using official MySQL image on Docker Hub start a MySQL 8 server.
2. By setting appropriate environment variables configure the following
   * Username
   * Password
   * Database name
   * Enable random root password
3. Connect to MySQL server using `mysql` command line tool.
   Note: you may need to install `mysql-client` ubuntu package for that.
4. Add MySQL server to the docker compose file created in **Task 3** and
   start it using docker-compose.
5. Connect to MySQL server with `mysql` client and create the phonebook
   table by running the following SQL query:
   ```sql
   CREATE TABLE records
   (
       id           VARCHAR(36) NOT NULL,
       first_name   VARCHAR(50) NOT NULL,
       last_name    VARCHAR(50) NOT NULL,
       phone_number VARCHAR(20) NOT NULL,
       CONSTRAINT records_pk
           PRIMARY KEY (id)
   );
   ```
6. Ensure the table exists by running the `SHOW TABLES;` query.


## Task 5 – Configure database connection for the service

1. Using the **Service configuration** section configure database
  connection for the service.
2. Restart the service by running `docker-compose up -d` again. Note how
   the container gets re-created after docker-compose notices changes in the
   compose file.
3. Using **HTTP endpoints** test the service endpoints again.

## Task 6 – Different methods of running HTTP queries

Try running HTTP queries using the following tools:

1. Browser
2. `curl` command line tool – please explore the manual to find out the options
3. Postman
4. (optional) python requests

What are the benefits and limitations of every tool

## Task 7 – Authentication

1. Following the [Service configuration](#service-configuration) section set configure
   an authentication token for the service.
2. Recreate the container with docker-compose to reflect the changes.
3. Try sending requests as you did before.
4. What HTTP status code are you getting?
5. Following the [Authentication](#authentication) section understand and fix the error.

## Service configuration

The service reads it's configuration from environment variables.
The following variables could be set:

* `SERVICE_PORT` – TCP port on which the service will listen for requests.
* `AUTH_TOKEN` – if set to any non-empty value, tells the service to check the
                 specified token on every request. Otherwise, all requests will be
                 allowed unauthenticated. Details in the [Authentication](#authentication)
                 section.
* `DB_HOST` – host on which host a MySQL database server is running.
* `DB_PORT` – port on which a MySQL server is listening.
* `DB_USERNAME` – name of a user on the MySQL server.
* `DB_PASSWORD` – password for the user on the MySQL server.
* `DB_NAME` – name of the database on the MySQL server.
* `DB_ENGINE` – type of the database engine. Should be either `mysql` or `postgresql`.


## HTTP endpoints

* `/records/` - Perform operations with the list of phonebook records.
  The following HTTP methods are accepted: 
  * **GET** – Returns list of all records. 
              Expected status on success: 200 - OK.
              Expected body on success: list of records.
  * **POST** - Creates a new record. A record is required in the
               request body. For details refer to the [Record format](#record-format)
               section.
               Expected status on success: 201 - Created.
               Expected body on success: created record.
* `/records/{record_id}` – Perform operations with a specific record
  identified by the `record_id`.
  The following HTTP methods are accepted:
   * **GET** – Returns a single record.
     Expected status on success: 200 - OK.
     Expected body on success: Specified record.
   * **DELETE** – Deletes specified record.
     Expected status on success: 204 - No Content.
   * **PUT** - Updates specified record. New record data is required
     in the request body. It is impossible to update record's id.
     For details refer to the [Record format](#record-format) section.
     Expected status on success: 200 - OK.
     Expected body on success: updated record.


## Record format

Every record has the following JSON attributes:

```json
{
   "id": "bda6174e-8bc5-11ed-833b-acde48001123",
   "first_name": "John",
   "last_name": "Smith",
   "phone_number": "9379992"
}
```

The `id` attribute should not be passed when creating a new record.
Otherwise, all attributes are mandatory.


## Authentication

Authentication is a process of proving the identity of a user. For this service, it
is an optional feature which may be turned on or off by specifying or not specifying
an authentication token. As mentioned in the [Service configuration](#service-configuration)
section the token may be specified with the `AUTH_TOKEN` environment variable.

If no token is set, the service will run in no-auth mode and will not authenticate
requests  at all. Any value set to that variable will set the service to authenticate
requests  before doing anything.

After the service is set to authenticate requests the token must be specified with
every request in the `Authentication` header. The value of the header must be in the
`TOKEN {YOUR_SECRET_TOKEN}` format (without the curly braces). Example:

```
Authentication: TOKEN my-super-secret-token
```

When the service is authenticating a request and finds an `Authentication` header with
a valid token, it then proceeds to process that request as usual. Otherwise, if the
header was not set, or the format is wrong, or the token is invalid, the service stops
processing that request and returns `401 - Unauthirized` status code.

