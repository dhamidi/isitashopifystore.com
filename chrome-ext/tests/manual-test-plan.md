# Manual Test Plan for Shopify Store Detector Extension

## Test Cases

### 1. Known Shopify Stores
- [ ] Visit allbirds.com - should show icon
- [ ] Visit gymshark.com - should show icon
- [ ] Visit kylie.com - should show icon

### 2. Non-Shopify Stores
- [ ] Visit amazon.com - should not show icon
- [ ] Visit ebay.com - should not show icon
- [ ] Visit walmart.com - should not show icon

### 3. Cache Behavior
- [ ] Visit a Shopify store, verify icon appears
- [ ] Close and reopen browser
- [ ] Revisit same store within 24h - should use cached result
- [ ] Wait 24h and revisit - should make new API call

### 4. Edge Cases
- [ ] Test on page with iframes
- [ ] Test on page with multiple redirects
- [ ] Test on local development Shopify stores
- [ ] Test on password-protected Shopify stores
- [ ] Test on Shopify checkout pages
- [ ] Test when API is down

### 5. Visual Tests
- [ ] Icon appears in top-right corner
- [ ] Icon doesn't overlap with page elements
- [ ] Icon is visible on light backgrounds
- [ ] Icon is visible on dark backgrounds
- [ ] Hover effects work correctly

### 6. Performance
- [ ] Extension doesn't noticeably impact page load time
- [ ] Icon appears within 1 second of page load
- [ ] Cached results return instantly

### 7. Browser States
- [ ] Works in normal browsing mode
- [ ] Works in incognito mode (if permitted)
- [ ] Works across multiple windows
- [ ] Works across multiple tabs

## Test Execution Log

| Date | Tester | Test Case | Result | Notes |
|------|--------|-----------|---------|-------| 