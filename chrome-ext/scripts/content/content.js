function injectShopifyIcon() {
  const existingIcon = document.querySelector('.shopify-detector-icon');
  if (existingIcon) return;

  const icon = document.createElement('img');
  icon.src = chrome.runtime.getURL('assets/shopify-icon-32.png');
  icon.className = 'shopify-detector-icon';
  icon.title = 'This is a Shopify store';
  
  document.body.appendChild(icon);
}

chrome.runtime.onMessage.addListener((message) => {
  if (message.type === 'SHOPIFY_STATUS' && message.data.isShopify) {
    injectShopifyIcon();
  }
  return true;
}); 