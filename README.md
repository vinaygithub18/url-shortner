# URL Shortener

A simple URL shortener service built with Go, Fiber, and Redis.

## Features

- Shorten long URLs to easy-to-share links
- Custom short URLs (optional)
- Expiration time for links
- Rate limiting to prevent abuse
- Redis-based storage

## Tech Stack

- Go (Golang)
- Fiber web framework
- Redis for data storage
- Docker and Docker Compose for containerization

## Project Structure

```
url-shortner/
├── api/                # Go application
│   ├── database/       # Redis client
│   ├── helpers/        # Utility functions
│   ├── routes/         # API endpoints
│   ├── .env            # Environment variables
│   ├── Dockerfile      # API container definition
│   ├── go.mod          # Go modules
│   └── main.go         # Entry point
├── db/                 # Redis configuration
│   └── Dockerfile      # Redis container definition
├── data/               # Redis persistence
└── docker-compose.yml  # Container orchestration
```

## Installation

### Prerequisites

- Docker and Docker Compose
- Go 1.16+ (for local development)

### Running with Docker

1. Clone the repository
   ```
   git clone https://github.com/yourusername/url-shortener.git
   cd url-shortener
   ```

2. Start the application
   ```
   docker-compose up -d
   ```

3. The service will be available at http://localhost:3000

### Local Development

1. Install dependencies
   ```
   cd api
   go mod download
   ```

2. Set up Redis locally or use Docker for Redis only
   ```
   docker-compose up -d db
   ```

3. Configure environment variables in `.env`

4. Run the application
   ```
   go run main.go
   ```

## API Usage

### Shorten a URL

```
POST /api/v1
Content-Type: application/json

{
  "url": "https://example.com/very/long/url/that/needs/shortening",
  "custom_short": "custom",  // Optional
  "expiry": 3600            // Optional (in seconds)
}
```

Response:
```json
{
  "url": "https://example.com/very/long/url/that/needs/shortening",
  "custom_short": "http://localhost:3000/custom",
  "expiry": 3600,
  "x_rate_remaining": 9,
  "x_rate_limit_reset": 30
}
```

### Access a shortened URL

```
GET /:shortURL
```

This will redirect to the original URL.

## Environment Variables

- `DB_ADDR`: Redis address (default: "db:6379")
- `DB_PASS`: Redis password
- `APP_PORT`: Application port (default: ":3000")
- `DOMAIN`: Domain for shortened URLs (default: "http://localhost:3000")
- `API_QUOTA`: Rate limit per IP (default: 10)

## License

MIT