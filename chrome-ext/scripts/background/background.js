import { checkDomain } from './api.js';
import { getCachedResult, setCachedResult } from './cache.js';
import { setupTabListeners } from './tabs.js';

// Track which tabs have content scripts ready
const readyTabs = new Set();
// Queue messages for tabs not ready yet
const messageQueue = new Map();

// Handle content script ready signals
chrome.runtime.onMessage.addListener((message, sender) => {
  if (message.type === 'CONTENT_SCRIPT_READY' && sender.tab) {
    console.log(`[Background] Content script ready in tab ${sender.tab.id}`);
    readyTabs.add(sender.tab.id);
    
    // Send any queued messages for this tab
    if (messageQueue.has(sender.tab.id)) {
      const queuedMessages = messageQueue.get(sender.tab.id);
      messageQueue.delete(sender.tab.id);
      queuedMessages.forEach(msg => sendMessageToTab(sender.tab.id, msg));
    }
  }
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