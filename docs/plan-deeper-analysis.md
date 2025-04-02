# Plan for Implementing Deeper Shopify Analysis

## Overview

Currently, our system detects Shopify stores by analyzing the main domain for Shopify-specific indicators. However, some Shopify stores use custom domains that don't have obvious Shopify markers on their main pages. This plan outlines how to enhance our detection by checking the checkout page, which often reveals Shopify usage through specific headers even when the main site doesn't.

## Implementation Steps

### 1. Modify the Analysis Logic

#### Files to Modify:
- `analysis.go`: Add a secondary check when initial Shopify detection fails

#### Changes Required:
1. Update the domain analysis function to include a fallback checkout page check
2. Construct the checkout URL using the pattern: `https://checkout.[domain]/checkout/cn`
3. Make an HTTP request to this checkout URL
4. Check the response headers for the presence of `x-shopid`
5. If the header exists, mark the site as a Shopify store regardless of the initial analysis

### 2. Implementation Details

#### Checkout URL Construction
- For a domain like `vanmoof.com`, construct `https://checkout.vanmoof.com/checkout/cn`
- Handle various domain formats correctly (subdomains, etc.)

#### HTTP Request Handling
- Set appropriate timeouts for the request
- Handle various response codes (especially 404, which is expected)
- Focus only on examining the response headers, not the body

#### Response Analysis
- Check specifically for the `x-shopid` header
- Update the analysis result based on this header's presence

### 3. Error Handling

- Handle network errors gracefully
- Consider rate limiting to avoid overloading servers
- Log detailed information about the checkout request for debugging

### 4. Testing Plan

#### Manual Testing Scenarios:
1. Test with known Shopify stores that don't show obvious Shopify markers
2. Test with non-Shopify stores to ensure no false positives
3. Test with edge cases (redirects, unusual domain structures)

#### Test Domains:
- `vanmoof.com` - Expected to be identified as Shopify via checkout page
- [Add other test domains here]

### 5. Integration with Existing System

- Ensure the checkout check only runs when the initial analysis doesn't find Shopify markers
- Update the result format to indicate when a site was identified via checkout page vs. main page
- Maintain backward compatibility with existing API responses

## Future Considerations

- Add automated testing for this feature in a future iteration
- Consider caching checkout results to improve performance
- Monitor for any changes in Shopify's checkout page structure or headers

## Success Criteria

The implementation will be considered successful if:
1. It correctly identifies Shopify stores that use the checkout.domain.com pattern
2. It maintains the accuracy of the existing detection system
3. It doesn't significantly increase analysis time for non-Shopify stores
