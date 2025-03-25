export async function checkDomain(domain) {
  try {
    const response = await fetch(`https://isitashopifystore.com/status/${domain}`);
    if (!response.ok) {
      throw new Error(`API request failed: ${response.status}`);
    }
    return await response.json();
  } catch (error) {
    console.error('Error checking domain:', error);
    return { isShopify: false, reason: 'API error', status: 'error' };
  }
} 