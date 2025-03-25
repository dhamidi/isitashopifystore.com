function log(message, ...args) {
  const logMessage = `[Content] ${message}`;
  chrome.runtime.sendMessage({ 
    type: 'LOG', 
    data: { message: logMessage, args } 
  });
}

function injectShopifyIcon() {
  log('Attempting to inject Shopify icon');
  const existingIcon = document.querySelector('.shopify-detector-icon');
  if (existingIcon) {
    log('Icon already exists, skipping injection');
    return;
  }

  const icon = document.createElement('img');
  icon.src = chrome.runtime.getURL('assets/shopify-icon-512.png');
  icon.className = 'shopify-detector-icon';
  icon.title = 'This is a Shopify store';
  
  document.body.appendChild(icon);
  log('Successfully injected Shopify icon');
}

// Send ready signal when content script loads
log('Sending ready signal');
chrome.runtime.sendMessage({ type: 'CONTENT_SCRIPT_READY' });

chrome.runtime.onMessage.addListener((message) => {
  log('Received message:', message);
  if (message.type === 'SHOPIFY_STATUS' && message.data.is_shopify) {
    log('Detected Shopify store, injecting icon');
    injectShopifyIcon();
  }
}); 