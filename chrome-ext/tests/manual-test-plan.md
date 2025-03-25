# Manual Test Plan for Shopify Store Detector Extension

## Test Cases

### 1. Known Shopify Stores
- [x] Visit allbirds.com - should show icon
- [x] Visit gymshark.com - should show icon
- [x] Visit kylie.com - should show icon
    - the URL changed in the meantime: https://store.kylie.com/ does show the icon

### 2. Non-Shopify Stores
- [x] Visit amazon.com - should not show icon
- [x] Visit ebay.com - should not show icon
- [x] Visit walmart.com - should not show icon

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
- [x] Icon appears in top-right corner
- [x] Icon doesn't overlap with page elements
- [x] Icon is visible on light backgrounds
- [x] Icon is visible on dark backgrounds
- [x] Hover effects work correctly

### 6. Performance
- [x] Extension doesn't noticeably impact page load time
- [x] Icon appears within 1 second of page load
- [x] Cached results return instantly

### 7. Browser States
- [x] Works in normal browsing mode
- [x] Works in incognito mode (if permitted)
- [x] Works across multiple windows
- [x] Works across multiple tabs

## Test Execution Log

| Date | Tester | Test Case | Result | Notes |
|------|--------|-----------|---------|-------| 