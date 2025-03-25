# IsItAShopifyStore.com

A web application that analyzes websites to determine if they are built on the Shopify platform. The service accepts a URL or domain name, performs an analysis of the website's content, and reports whether it detects Shopify indicators.

## Features

- URL/Domain analysis for Shopify detection
- Real-time analysis status updates
- Event logging for analysis tracking
- Simple and intuitive web interface

## Prerequisites

- Go 1.x or later
- SQLite3

## Building

1. Clone the repository:
```bash
git clone https://github.com/yourusername/isitashopifystore.com.git
cd isitashopifystore.com
```

2. Install dependencies:
```bash
go mod download
```

3. Build the application:
```bash
go build
```

## Running Locally

Start the server:
```bash
./isitashopifystore.com
```

The application will be available at `http://localhost:8080`

## Deployment

The application is designed to be deployed on any platform that supports Go applications. Here's a basic deployment guide:

1. Build the application for your target platform
2. Ensure SQLite3 is installed on the server
3. Set up environment variables if needed
4. Run the application with appropriate process management (e.g., systemd, supervisor)

### Docker Deployment

1. Build the Docker image:
```bash
docker build -t isitashopifystore .
```

2. Run the container:
```bash
docker run -p 8080:8080 isitashopifystore
```

## Project Structure

- `main.go` - Application entry point and server setup
- `handler.go` - HTTP request handlers
- `analysis.go` - Shopify detection logic
- `db.go` - Database operations and event logging

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

[Add your license information here] 