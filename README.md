# IsItAShopifyStore.com

A web application that analyzes websites to determine if they are built on the Shopify platform. The service accepts a URL or domain name, performs an analysis of the website's content, and reports whether it detects Shopify indicators.

## Features

- URL/Domain analysis for Shopify detection
- Real-time analysis status updates
- Event logging for analysis tracking
- Simple and intuitive web interface
- Chrome extension for instant website analysis

## Prerequisites

- Go 1.x or later
- SQLite3
- For Chrome extension development: Chrome browser

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
# For local development
./tools/build local

# For DigitalOcean deployment
./tools/build digitalocean
```

## Running Locally

Start the server:
```bash
./isitashopifystore
```

The application will be available at `http://localhost:8080`

## Chrome Extension

The Chrome extension allows users to instantly check if their current website is a Shopify store. 

### Development Installation

1. Navigate to the chrome extension directory:
```bash
cd chrome-ext
```

2. Open Chrome and go to `chrome://extensions/`
3. Enable "Developer mode" in the top right corner
4. Click "Load unpacked" and select the `chrome-ext` directory

### Building for Production

```bash
cd chrome-ext
./build.sh
```

This will create a zip file in the `dist` directory ready for Chrome Web Store submission.

## Deployment

The application is deployed on DigitalOcean:

1. Build the application for the target platform:
```bash
./tools/build digitalocean
```

2. Deploy to DigitalOcean:
```bash
./tools/deploy/push-to-do
```

The deployment uses Caddy for HTTPS and reverse proxying.

## Project Structure

- `main.go` - Application entry point and server setup
- `handler.go` - HTTP request handlers
- `analysis.go` - Shopify detection logic
- `db.go` - Database operations and event logging
- `chrome-ext/` - Chrome extension files
- `html/` - HTML templates
- `assets/` - Static assets
- `tools/` - Build and deployment tools

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

[Add your license information here] 