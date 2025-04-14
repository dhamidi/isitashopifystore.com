import { checkDomain } from './api.js';
import { getCachedResult, setCachedResult } from './cache.js';
import { setupTabListeners } from './tabs.js';

// Track which tabs have content scripts ready
const readyTabs = new Set();
// Queue messages for tabs not ready yet
const messageQueue = new Map();

// Function to clear the entire cache
async function clearAllCache() {
  try {
    await chrome.storage.local.clear();
    console.log('[Background] Cache cleared successfully.');
    return { success: true };
  } catch (error) {
    console.error('[Background] Error clearing cache:', error);
    return { success: false, error: error.message };
  }
}

// Handle messages from content scripts and popup
chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  if (message.type === 'CONTENT_SCRIPT_READY' && sender.tab) {
    console.log(`[Background] Content script ready in tab ${sender.tab.id}`);
    readyTabs.add(sender.tab.id);
    
    // Send any queued messages for this tab
    if (messageQueue.has(sender.tab.id)) {
      const queuedMessages = messageQueue.get(sender.tab.id);
      messageQueue.delete(sender.tab.id);
      queuedMessages.forEach(msg => sendMessageToTab(sender.tab.id, msg));
    }
    return false; // No async response needed for this message type
  } else if (message.type === 'LOG') {
    console.log(...[message.data.message, ...(message.data.args || [])]);
    return false; // No async response needed for this message type
  } else if (message.action === 'clearCache') {
    console.log('[Background] Received clear cache request');
    // Handle cache clearing asynchronously
    clearAllCache().then(response => {
      sendResponse(response);
    }).catch(error => {
      console.error('[Background] Failed to clear cache:', error);
      sendResponse({ success: false, error: error.message });
    });
    return true; // Indicate that the response will be sent asynchronously
  }
  // Default case if no message type matches
  return false;
});

// Clean up when tabs are closed
chrome.tabs.onRemoved.addListener((tabId) => {
  readyTabs.delete(tabId);
  messageQueue.delete(tabId);
});

async function sendMessageToTab(tabId, message) {
  if (!readyTabs.has(tabId)) {
    console.log(`[Background] Tab ${tabId} not ready, queueing message:`, message);
    if (!messageQueue.has(tabId)) {
      messageQueue.set(tabId, []);
    }
    messageQueue.get(tabId).push(message);
    return;
  }

  try {
    console.log(`[Background] Sending message to tab ${tabId}:`, message);
    await chrome.tabs.sendMessage(tabId, message);
  } catch (error) {
    console.error('[Background] Error sending message to tab:', error);
  }
}

async function handleDomainChange(tabId, domain) {
  console.log(`[Background] Checking domain ${domain} for tab ${tabId}`);
  let result = await getCachedResult(domain);
  
  if (!result) {
    console.log(`[Background] Cache miss for ${domain}, fetching from API`);
    result = await checkDomain(domain);
    console.log(`[Background] API result for ${domain}:`, result);
    await setCachedResult(domain, result);
  } else {
    console.log(`[Background] Cache hit for ${domain}:`, result);
  }

  await sendMessageToTab(tabId, {
    type: 'SHOPIFY_STATUS',
    data: result
  });
}

setupTabListeners(handleDomainChange); 