# Backend

This is a backend service built using PocketBase and Echo.

## How to Start the Backend

1. Clone the repository.
2. Make sure Go is installed on your machine.
3. Use the following commands to start the backend:

   - **For Development**:

     ```sh
     go run main.go --dev serve
     ```

   - **For Production**:

     ```sh
     go run main.go serve
     ```

4. The server will be running and accessible at [http://localhost:8090](http://localhost:8090).
   - REST API: [http://127.0.0.1:8090/api/](http://127.0.0.1:8090/api/)
   - Admin UI: [http://127.0.0.1:8090/_/](http://127.0.0.1:8090/_/)

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