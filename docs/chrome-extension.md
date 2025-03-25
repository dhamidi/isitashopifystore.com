# Chrome Extension Specification

## Overview
A Chrome extension that automatically checks if the current website is a Shopify store using the isitashopifystore.com API.

## Location
The extension code should be placed in the `./chrome-ext` directory.

## Functionality
1. The extension automatically checks if the current website the user is visiting is a Shopify store.
2. It makes an API call to `https://isitashopifystore.com/status/{domain}` to determine the status.
3. If the website is a Shopify store, the extension displays a small Shopify icon on the page.
4. The icon is only displayed when the result is confirmed (after the API call completes), not while waiting for the response.

## API Integration
- The extension will use the `/status/{domain}` endpoint.
- The API returns an `AnalysisResult` with the following structure:
  ```json
  {
    "status": "string",
    "reason": "string",
    "is_shopify": boolean,
    "domain": "string"
  }
  ```
- The extension should only show the Shopify icon when `is_shopify` is `true`.

## Technical Implementation
1. Create a manifest.json file for extension configuration.
2. Implement a background script to handle tab URL changes and make API requests.
3. Create a content script to display the Shopify icon on the page when appropriate.
4. Include necessary assets like the Shopify icon.
5. Ensure good performance by caching results when possible. 