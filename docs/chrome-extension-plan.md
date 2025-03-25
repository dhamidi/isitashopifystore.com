# Chrome Extension Implementation Plan

This document outlines the step-by-step plan for implementing the Shopify store detector Chrome extension as specified in `chrome-extension.md`.

## 1. Project Setup - DONE

### Step 1.1: Create Extension Directory Structure
- Create the `chrome-ext` directory at the project root
- Create subdirectories for assets, scripts, and styles
- **Definition of Done**: The following directory structure exists:
  ```
  chrome-ext/
  ├── assets/
  ├── scripts/
  │   ├── background/
  │   └── content/
  └── styles/
  ```

### Step 1.2: Copy and Optimize Icon Assets
- Copy `assets/favicon-512.png` to `chrome-ext/assets/shopify-icon-512.png`
- Create resized versions at 128px, 48px, 32px, and 16px for different contexts
- **Definition of Done**: All icon sizes exist in the `chrome-ext/assets/` directory

## 2. Extension Configuration - DONE

### Step 2.1: Create Manifest File
- Create `chrome-ext/manifest.json` with the following configurations:
  - Basic extension metadata (name, version, description)
  - Permissions (activeTab, storage, host permissions for API)
  - Background script registration
  - Content script registration with appropriate matches pattern
  - Extension icons at different sizes
- **Definition of Done**: Valid `manifest.json` file that passes Chrome's validation

### Step 2.2: Setup Extension Icons
- Configure the manifest to use the Shopify icon for the extension
- **Definition of Done**: Extension uses the Shopify icon when installed

## 3. Background Script Implementation - DONE

### Step 3.1: Create API Service
- Create `chrome-ext/scripts/background/api.js`
- Implement function to check if a domain is a Shopify store via the API
- Add error handling and response parsing
- **Definition of Done**: Function successfully calls the API and returns parsed results

### Step 3.2: Create Domain Cache
- Create `chrome-ext/scripts/background/cache.js`
- Implement functions to store/retrieve domain results using Chrome's storage API
- Add cache expiration logic (24 hours)
- **Definition of Done**: Cache functions correctly store and retrieve domain data

### Step 3.3: Create Tab Handler
- Create `chrome-ext/scripts/background/tabs.js`
- Implement event listeners for tab updates and navigation
- Extract domains from URLs
- **Definition of Done**: Background script correctly detects when a user navigates to a new domain

### Step 3.4: Create Main Background Script
- Create `chrome-ext/scripts/background/background.js`
- Integrate API service, cache, and tab handler
- Implement messaging between background and content scripts
- **Definition of Done**: Background script orchestrates the detection process and communicates results to content scripts

## 4. Content Script Implementation - DONE

### Step 4.1: Create Content Script
- Create `chrome-ext/scripts/content/content.js`
- Implement message listeners to receive results from background script
- Implement function to create and inject the Shopify icon into the page
- Style the icon to appear in the top-right corner of the page
- Add hover effects and positioning logic
- Call icon injection when a positive result is received
- **Definition of Done**: Content script successfully listens for messages and injects a well-styled icon when appropriate

### Step 4.2: Create Content Styles
- Create `chrome-ext/styles/content.css`
- Define styles for the injected Shopify icon
- Ensure proper z-index to appear above page content
- **Definition of Done**: CSS file includes all necessary styles for the icon

## 5. Testing

### Step 5.1: Manual Test Plan
- Create `chrome-ext/tests/manual-test-plan.md`
- Define test cases for:
  - Known Shopify stores
  - Non-Shopify stores
  - Edge cases (redirects, frames, etc.)
- **Definition of Done**: Comprehensive test plan document exists

### Step 5.2: Execute Manual Tests
- Test the extension on various sites according to the test plan
- Verify proper icon display only on Shopify stores
- Test caching behavior
- **Definition of Done**: Extension passes all test cases in the manual test plan

## 6. Packaging and Distribution

### Step 6.1: Create Build Script
- Create `chrome-ext/build.sh` script
- Implement logic to:
  - Clean output directory
  - Copy necessary files
  - Zip contents for Chrome Web Store submission
- **Definition of Done**: Script successfully creates a deployable zip file

### Step 6.2: Document Installation Process
- Create `chrome-ext/README.md`
- Document how to:
  - Install the extension in developer mode
  - Build for production
  - Update the extension
- **Definition of Done**: Clear and complete documentation exists

## 7. Implementation Details

### Background Script (background.js) Logic:
```javascript
// 1. Listen for tab updates
// 2. When a tab navigates to a new URL:
//    a. Extract the domain from URL
//    b. Check if domain is in cache
//    c. If cached and valid, send result to content script
//    d. If not cached or expired, query API
//    e. Store result in cache
//    f. Send result to content script
```

### Content Script (content.js) Logic:
```javascript
// 1. Listen for messages from background script
// 2. When a "isShopify" message is received:
//    a. If true, inject Shopify icon
//    b. If false, do nothing
```

### API Service Logic:
```javascript
// Function: checkDomain(domain)
// 1. Construct API URL: https://isitashopifystore.com/status/{domain}
// 2. Make fetch request
// 3. Parse JSON response
// 4. Return { isShopify, reason, status }
```

## 8. Final Review and Submission

### Step 8.1: Code Review
- Conduct internal code review
- Check for:
  - Code quality and organization
  - Performance considerations
  - Security best practices
- **Definition of Done**: All code review comments addressed

### Step 8.2: Final Testing
- Test the extension in different Chrome versions
- Verify functionality across operating systems
- **Definition of Done**: Extension works consistently across environments

### Step 8.3: Prepare for Submission
- Create screenshots for Chrome Web Store
- Write promotional description
- **Definition of Done**: All Chrome Web Store assets prepared

## 9. Maintenance Plan

- Monitor API endpoint for changes
- Establish version update strategy
- Document future feature ideas
- **Definition of Done**: Maintenance plan document created 