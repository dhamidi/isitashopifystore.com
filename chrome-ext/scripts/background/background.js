import { checkDomain } from './api.js';
import { getCachedResult, setCachedResult } from './cache.js';
import { setupTabListeners } from './tabs.js';

async function handleDomainChange(tabId, domain) {
  let result = await getCachedResult(domain);
  
  if (!result) {
    result = await checkDomain(domain);
    await setCachedResult(domain, result);
  }

  try {
    await chrome.tabs.sendMessage(tabId, {
      type: 'SHOPIFY_STATUS',
      data: result
    });
  } catch (error) {
    console.debug('Could not send message to tab, content script may not be ready:', error);
  }
}

setupTabListeners(handleDomainChange); 