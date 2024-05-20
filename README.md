# YoungAstrologer
The service of a young astrologer, which receives metadata and an image of the APOD day once a day and, upon request, gives all album recordings and recordings for the selected day.

## Getting Started

### Configuration

Ensure you have the following environment variables set:

- `POSTGRES_URL` - The URL for connecting to the PostgreSQL database.
- `NASA_API_KEY` - Your NASA API key for fetching APOD data.
- `SERVER_PORT` - The port on which the server will run.


### Docker Compose

The service can be built and run using Docker Compose.

#### 1. Build the Docker images

This command builds the Docker images for both the application and the database.

```sh
docker-compose build
```

#### 2. Start the service

This command starts the application and the PostgreSQL database in the background.

```sh
docker-compose up -d
```

#### 3. Stop the service

This command stops the running containers for the application and the database.

```sh
docker-compose down
```

## Accessing the Service

Once the service is running, you can access it via the following endpoints:

    Retrieve all images.
    GET /images

    Retrieve an image by date.
    GET /images/date?date=YYYY-MM-DD