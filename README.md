# movierecserver 

## Overview
This is a Go-based backend for a movie application that provides user management, movie information retrieval, and personalized features such as watchlists and favorites. 

## Architecture

The application follows a layered architecture pattern:

1. **Router Layer**: Handles incoming HTTP requests
2. **Service Layer**: Implements core business logic
3. **Model Layer**: Defines data structures
4. **Database Layer**: Manages data persistence using PostgreSQL and GORM

### Environment Variables

A template `.env` file is provided in the repository with the following variables:
GOOGLE_CLIENT_ID=""
GOOGLE_CLIENT_SECRET=""
TMDB_API_KEY=""
AUTH_KEY=""
POSTGRES_DB=""
POSTGRES_USER=""
POSTGRES_PASSWORD=""
DB_HOST=""
DB_PORT=""
DB_SSLMODE=""
DB_TIMEZONE=""

To set up your environment:

1. Copy the `.env` file to `.env.local`:
cp .env .env.local

2. Open `.env.local` and fill in the appropriate values for each variable.

3. The `.env.local` file is gitignored, ensuring your sensitive information isn't committed to the repository.

**Note:** Never commit your actual environment values to the repository. The `.env` file in the repository should only contain empty variables as a template.

### Running the Application

1. Clone the repository:
git clone https://github.com/abhishekdas600/movierecserver.git

2. Navigate to the project directory:
cd movie-application

3. Start the PostgreSQL database using Docker:
docker-compose up -d

4. Run the server:
go run main.go

5. The application should now be running and accessible at `http://localhost:8080`.

## API Documentation

API documentation is available via the OpenAPI specification file (`openapi-movie-spec.yaml`, `openapi-user-spec.yaml`).



## Acknowledgments

- The Movie Database (TMDb) for providing movie data
- Google for authentication services
