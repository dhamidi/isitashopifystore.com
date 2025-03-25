# Shopify Store Detector Chrome Extension

This Chrome extension detects whether a website is powered by Shopify by checking against the isitashopifystore.com API.

## Development Installation

1. Clone this repository
2. Open Chrome and navigate to `chrome://extensions/`
3. Enable "Developer mode" in the top right corner
4. Click "Load unpacked" and select the `chrome-ext` directory

## Building for Production

```bash
./build.sh
```

This will create a zip file in the `dist` directory ready for Chrome Web Store submission.

## Updating the Extension

### In Development Mode
1. Make your changes to the extension code
2. Go to `chrome://extensions/`
3. Click the refresh icon on the extension card

### Production Version
The extension will auto-update when a new version is published to the Chrome Web Store. 