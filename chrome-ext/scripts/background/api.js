export async function checkDomain(domain) {
  const maxAttempts = 60;
  let attempts = 0;

  while (attempts < maxAttempts) {
    try {
      const response = await fetch(`https://isitashopifystore.com/status/${domain}`);
      if (!response.ok) {
        throw new Error(`API request failed: ${response.status}`);
      }
      const result = await response.json();
      if (result.status !== 'in_progress' || attempts === maxAttempts - 1) {
        return result;
      }
      await new Promise(resolve => setTimeout(resolve, 1000));
      attempts++;
    } catch (error) {
      console.error('Error checking domain:', error);
      return { isShopify: false, reason: 'API error', status: 'error' };
    }
  }
} 