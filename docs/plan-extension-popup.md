# Extension Popup Implementation Plan

## Overview

This document outlines the detailed implementation plan for adding a popup interface to the Shopify Store
Detector Chrome extension. The popup will appear when users click the extension icon in the browser toolbar
and will include:

1. A clear visual indicator showing whether the current site is a Shopify store (green) or not (red)
2. A button that allows users to clear the extension's cache
3. Basic styling to match the extension's design language

This functionality will improve the user experience by providing immediate feedback about the current site
and allowing users to force a fresh check when needed.

## Implementation Steps

### 1. Create the Popup HTML File

**File to create:** `chrome-ext/popup/popup.html`

This file will define the structure of our popup interface.

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Shopify Store Detector</title>
    <link rel="stylesheet" href="popup.css" />
  </head>
  <body>
    <div class="container">
      <div class="header">
        <img
          src="../assets/shopify-icon-48.png"
          alt="Shopify Icon"
          class="logo"
        />
        <h1>Shopify Detector</h1>
      </div>

      <div class="status-container">
        <div id="status-indicator" class="status-indicator">
          <span id="status-text">Checking...</span>
        </div>
      </div>

      <div class="actions">
        <button id="clear-cache-btn" class="btn">Clear Cache</button>
      </div>

      <div class="footer">
        <p>Detects if a website is running on Shopify</p>
      </div>
    </div>

    <script src="popup.js" type="module"></script>
  </body>
</html>
```

### 2. Create the Popup CSS File

**File to create:** `chrome-ext/popup/popup.css`

This file will style our popup interface.

```css
body {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen,
    Ubuntu, Cantarell, "Open
Sans", "Helvetica Neue", sans-serif;
  margin: 0;
  padding: 0;
  width: 300px;
  color: #333;
}

.container {
  padding: 16px;
}

.header {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
}

.logo {
  width: 24px;
  height: 24px;
  margin-right: 8px;
}

h1 {
  font-size: 18px;
  margin: 0;
  font-weight: 500;
}

.status-container {
  margin-bottom: 16px;
}

.status-indicator {
  padding: 12px;
  border-radius: 4px;
  text-align: center;
  font-weight: 500;
}

.status-indicator.shopify {
  background-color: #95bf47;
  color: white;
}

.status-indicator.not-shopify {
  background-color: white;
  color: #ff0000;
  border: 1px solid #ff0000;
}

.status-indicator.loading {
  background-color: #f5f5f5;
  color: #666;
}

.actions {
  margin-bottom: 16px;
}

.btn {
  width: 100%;
  padding: 8px 12px;
  background-color: #95bf47;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: background-color 0.2s;
}

.btn:hover {
  background-color: #7ea039;
}

.footer {
  font-size: 12px;
  color: #666;
  text-align: center;
}
```

### 3. Create the Popup JavaScript File

**File to create:** `chrome-ext/popup/popup.js`

This file will handle the popup's functionality.

```javascript
// Get DOM elements
const statusIndicator = document.getElementById("status-indicator");
const statusText = document.getElementById("status-text");
const clearCacheBtn = document.getElementById("clear-cache-btn");

// Initialize popup
async function initPopup() {
  // Set initial loading state
  setLoadingState();

  // Get current tab
  const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });

  if (!tab || !tab.url || tab.url.startsWith("chrome://")) {
    setNotApplicableState();
    return;
  }

  // Get domain from URL
  const domain = extractDomain(tab.url);
  if (!domain) {
    setNotApplicableState();
    return;
  }

  // Check if domain is in cache
  const cachedResult = await getCachedResult(domain);
  if (cachedResult) {
    updateStatusDisplay(cachedResult);
  } else {
    // If not in cache, we'll show loading state until background script updates
    setLoadingState();
    // Send message to background script to check domain
    chrome.runtime.sendMessage({
      type: "CHECK_DOMAIN",
      data: { domain },
    });
  }
}

// Extract domain from URL
function extractDomain(url) {
  try {
    const { hostname } = new URL(url);
    return hostname;
  } catch {
    return null;
  }
}

// Get cached result for domain
async function getCachedResult(domain) {
  const data = await chrome.storage.local.get(domain);
  return data[domain]?.result || null;
}

// Update status display based on result
function updateStatusDisplay(result) {
  statusIndicator.classList.remove("loading", "shopify", "not-shopify");

  if (result.is_shopify || result.isShopify) {
    statusIndicator.classList.add("shopify");
    statusText.textContent = "This is a Shopify store";
  } else {
    statusIndicator.classList.add("not-shopify");
    statusText.textContent = "Not a Shopify store";
  }
}

// Set loading state
function setLoadingState() {
  statusIndicator.classList.remove("shopify", "not-shopify");
  statusIndicator.classList.add("loading");
  statusText.textContent = "Checking...";
}

// Set not applicable state
function setNotApplicableState() {
  statusIndicator.classList.remove("loading", "shopify");
  statusIndicator.classList.add("not-shopify");
  statusText.textContent = "Cannot check this page";
}

// Clear cache
async function clearCache() {
  await chrome.storage.local.clear();

  // Show confirmation
  clearCacheBtn.textContent = "Cache Cleared!";
  setTimeout(() => {
    clearCacheBtn.textContent = "Clear Cache";
  }, 1500);

  // Re-initialize popup to refresh status
  initPopup();
}

// Listen for messages from background script
chrome.runtime.onMessage.addListener((message) => {
  if (message.type === "SHOPIFY_STATUS") {
    updateStatusDisplay(message.data);
  }
});

// Add event listener for clear cache button
clearCacheBtn.addEventListener("click", clearCache);

// Initialize popup when DOM is loaded
document.addEventListener("DOMContentLoaded", initPopup);
```

### 4. Update the Manifest File

**File to modify:** `chrome-ext/manifest.json`

We need to update the manifest to include the popup.

Add the following to the manifest:

```json
"action": {
  "default_popup": "popup/popup.html",
  "default_icon": {
    "16": "assets/shopify-icon-16.png",
    "32": "assets/shopify-icon-32.png",
    "48": "assets/shopify-icon-48.png",
    "128": "assets/shopify-icon-128.png"
  }
}
```

### 5. Update Background Script to Handle Popup Requests

**File to modify:** `chrome-ext/scripts/background/background.js`

We need to add a message listener for the popup's "CHECK_DOMAIN" message.

Add the following to the existing message listener:

```javascript
// Inside the existing chrome.runtime.onMessage.addListener function
if (message.type === "CHECK_DOMAIN" && message.data.domain) {
  handleDomainChange(sender.tab.id, message.data.domain);
}
```

### 6. Testing Plan

1. **Load the extension:**

   - Open Chrome and navigate to `chrome://extensions/`
   - Enable "Developer mode"
   - Click "Load unpacked" and select the `chrome-ext` directory

2. **Test on a Shopify site:**

   - Navigate to a known Shopify site (e.g., allbirds.com)
   - Click the extension icon
   - Verify the green indicator appears with "This is a Shopify store" text

3. **Test on a non-Shopify site:**

   - Navigate to a non-Shopify site (e.g., google.com)
   - Click the extension icon
   - Verify the red indicator appears with "Not a Shopify store" text

4. **Test cache clearing:**

   - Navigate to a Shopify site
   - Click the extension icon
   - Click "Clear Cache"
   - Verify the "Cache Cleared!" message appears
   - Verify the status updates correctly after clearing

5. **Test on invalid pages:**
   - Navigate to a chrome:// page
   - Click the extension icon
   - Verify the "Cannot check this page" message appears

### 7. Directory Structure

After implementation, the extension directory structure should look like:

```
chrome-ext/
├── assets/
│   ├── shopify-icon-16.png
│   ├── shopify-icon-32.png
│   ├── shopify-icon-48.png
│   ├── shopify-icon-128.png
│   └── shopify-icon-512.png
├── popup/
│   ├── popup.html
│   ├── popup.css
│   └── popup.js
├── scripts/
│   ├── background/
│   │   ├── api.js
│   │   ├── background.js
│   │   ├── cache.js
│   │   └── tabs.js
│   └── content/
│       └── content.js
├── styles/
│   └── content.css
└── manifest.json
```

## Implementation Notes

1. **Color Scheme:**

   - Use Shopify green (#95BF47) for positive indicators and buttons
   - Use red (#FF0000) for negative indicators

2. **Error Handling:**

   - All network requests should have proper error handling
   - The UI should never be in a broken state

3. **Performance Considerations:**

   - Cache results to minimize API calls
   - Clear specific domains rather than the entire cache when possible

4. **Accessibility:**
   - Ensure proper color contrast for text
   - Use semantic HTML elements
   - Include proper alt text for images

## Future Enhancements

1. Add ability to view detailed analysis results
2. Add option to clear cache for specific domains only
3. Add settings page for customizing extension behavior
4. Add keyboard shortcuts for common actions

## Resources

- [Chrome Extension Documentation](https://developer.chrome.com/docs/extensions/)
- [Chrome Storage API](https://developer.chrome.com/docs/extensions/reference/storage/)
- [Chrome Tabs API](https://developer.chrome.com/docs/extensions/reference/tabs/)
