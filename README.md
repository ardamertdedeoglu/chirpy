# Chirpy

Chirpy is a simple social media API built in Go, allowing users to create accounts, log in, and post short messages called "chirps." It features user authentication with JWT tokens and refresh tokens, password hashing, and a clean REST API.

## Why Care

Chirpy demonstrates modern Go development practices, including database interactions with sqlc, secure authentication, and API design. It's a great starting point for learning backend development or building your own social platform.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/ardamertdedeoglu/chirpy.git
   cd chirpy
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Set up the database (assuming PostgreSQL):
   - Create a database named `chirpy`.
   - Run the schema files in `sql/schema/` in order.

4. Generate database code:
   ```bash
   sqlc generate
   ```

## Running the Project

1. Set environment variables (e.g., database URL, JWT secret).

2. Run the server:
   ```bash
   go run .
   ```

The API will be available at `http://localhost:8080` (or your configured port).

## Environment Variables

Set the following environment variables in a `.env` file or your environment:

- `DB_URL`: PostgreSQL database connection string (e.g., `postgres://user:password@localhost:5432/chirpy`)
- `SECRET`: JWT secret key for token signing
- `PLATFORM`: Platform identifier (e.g., "dev")
- `POLKA_KEY`: API key for Polka webhooks

## API Endpoints

- `GET /api/healthz`: Health check
- `POST /api/users`: Create a new user
- `POST /api/login`: User login
- `POST /api/refresh`: Refresh access token
- `POST /api/revoke`: Revoke refresh token
- `PUT /api/users`: Update user information
- `POST /api/chirps`: Create a new chirp
- `GET /api/chirps`: Get all chirps
- `GET /api/chirps/{chirpID}`: Get a specific chirp
- `DELETE /api/chirps/{chirpID}`: Delete a chirp
- `POST /api/polka/webhooks`: Handle Polka webhooks
- `POST /admin/reset`: Reset database (admin)
- `GET /admin/metrics`: Get metrics (admin)
