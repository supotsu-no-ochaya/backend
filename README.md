# Backend

This is a backend service built using PocketBase and Echo.

---

## Run the backend
Requirements: 
* Docker

### Run with testdata
To run the docker container with testdata you need to run the ```docker-compose.yml``` 
and need to have the ```pb_data``` testdata folder. This folder will be mounted into the container to supply it with the testdata.
To download the testdata run the ```download-testdata.sh``` script by supplying the latest testdatat version e.g.: 
```shell
./download-testdata.sh v2
```
You can now run the docker container. 
```shell
docker compose up -d
```

### Run without testdata
To run the docker container without testdata simply run the ```docker-compose.prod.yml```.
The ```docker-compose.prod.yml``` serves as an example configuration.
For security you should provide your own credentials as environment variables e.g. via an env file.

```shell
docker compose -f docker-compose.prod.yml up -d
```
---
The server will be running and accessible at [http://localhost:8090](http://localhost:8090).
- REST API: [http://127.0.0.1:8090/api/](http://127.0.0.1:8090/api/)
- Admin UI: [http://127.0.0.1:8090/_/](http://127.0.0.1:8090/_/)

---

## Contribute
1. Clone the repository.
2. Make sure Go is installed on your machine.
3. Use the following commands to start the backend:

   - **For Development (run from source code)**:
     - Ensure the **`download-testdata.sh`** script has been executed if there is no test data yet. This script writes the test data into the current working directory in a folder named `pb_data`, which both the binary and Docker Compose setup expect to be present.
       ```sh
       ./download-testdata.sh v1  
       ```
       ```sh
       go run cmd/app/main.go --dev serve
       ```

   - **For Development (run via docker compose)**:
     - Ensure the **`download-testdata.sh`** script has been executed if there is no test data yet. This script writes the test data into the current working directory in a folder named `pb_data`, which both the binary and Docker Compose setup expect to be present.
       ```sh
       ./download-testdata.sh v1  
       ```
       ```sh
       docker compose up -d
       ```
4. In case you need to build the docker image localy run:
      ```sh
   make docker-local
      ```
   The docker image will be named ```supotsu-backend:local```
---
## Release
Releases of the docker image happens automatically via the ```container-release``` workflow in github actions. The pipeline runs automatically on every new commit on main. So on every direct push, merge or rebase etc..
---

## API Endpoint

### `/api/test`
Fetches all user records from the database.
- **Method**: `GET`
- **Response**:
    - `200 OK` with a JSON array of user records.
    - `500 Internal Server Error` if unable to fetch users.
- **Example**:
    ```sh
    curl -X GET http://localhost:8090/api/test
    ```

### `/api/export-json`
Exports product and event data within a specified datetime range as a JSON file.
- **Method**: `GET`
- **Query Parameters**:
    - `start`: (required): Start datetime in RFC3339 format.
    - `end`: (required): End datetime in RFC3339 format.
- **Response**:
    - `200 OK` with a downloadable JSON file containing the export data.
    - `400 Bad Request` if query parameters are missing or invalid.
    - `500 Internal Server Error` if an error occurs during data fetching or processing.
- **Example**:
    ```sh
    curl -o export.json "http://localhost:8090/api/export-json?start=2023-01-01T00:00:00Z&end=2023-12-31T23:59:59Z"
    ```
- **Note**:
    - The exported JSON includes products enriched with attributes, categories, and station details.
    - Events are filtered based on the `created` timestamp within the provided datetime range.
