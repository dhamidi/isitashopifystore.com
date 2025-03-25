function injectShopifyIcon() {
  console.log('[Content] Attempting to inject Shopify icon');
  const existingIcon = document.querySelector('.shopify-detector-icon');
  if (existingIcon) {
    console.log('[Content] Icon already exists, skipping injection');
    return;
  }

  const icon = document.createElement('img');
  icon.src = chrome.runtime.getURL('assets/shopify-icon-32.png');
  icon.className = 'shopify-detector-icon';
  icon.title = 'This is a Shopify store';
  
  document.body.appendChild(icon);
  console.log('[Content] Successfully injected Shopify icon');
}

// Send ready signal when content script loads
console.log('[Content] Sending ready signal');
chrome.runtime.sendMessage({ type: 'CONTENT_SCRIPT_READY' });

chrome.runtime.onMessage.addListener((message) => {
  console.log('[Content] Received message:', message);
  if (message.type === 'SHOPIFY_STATUS' && message.data.isShopify) {
    console.log('[Content] Detected Shopify store, injecting icon');
    injectShopifyIcon();
  }
}); 