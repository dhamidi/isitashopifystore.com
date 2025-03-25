export function extractDomain(url) {
  try {
    const { hostname } = new URL(url);
    return hostname;
  } catch {
    return null;
  }
}

export function setupTabListeners(onDomainChange) {
  chrome.tabs.onUpdated.addListener((tabId, changeInfo) => {
    if (changeInfo.status === 'loading') {
      chrome.tabs.get(tabId).then(tab => {
        if (tab.url) {
          const domain = extractDomain(tab.url);
          if (domain) {
            onDomainChange(tabId, domain);
          }
        }
      });
    }
  });
} 