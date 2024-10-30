# Backend

This is a backend service built using PocketBase and Echo.

## How to Start the Backend

1. Clone the repository.
2. Make sure Go is installed on your machine.
3. Run the following command to start the backend:

   ```sh
   go run main.go
   ```

4. The server will be running and accessible at `http://localhost:8090`.
   - REST API: `http://127.0.0.1:8090/api/`
   - Admin UI: `http://127.0.0.1:8090/_/`

### API Endpoint
- `/api/test`: Fetches all user records from the database.

