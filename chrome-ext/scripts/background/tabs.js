export function extractDomain(url) {
  try {
    if (url.startsWith('chrome://')) return null;
    const { hostname } = new URL(url);
    return hostname;
  } catch {
    return null;
  }
}

async function getCurrentTab() {
  const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });
  return tab;
}

export function setupTabListeners(onDomainChange) {
  chrome.tabs.onActivated.addListener(async (activeInfo) => {
    const tab = await getCurrentTab();
    if (tab?.url) {
      const domain = extractDomain(tab.url);
      if (domain) {
        onDomainChange(tab.id, domain);
      }
    }
  });

  chrome.tabs.onUpdated.addListener(async (tabId, changeInfo) => {
    if (changeInfo.status === 'complete') {
      const tab = await getCurrentTab();
      if (tab?.id === tabId && tab.url) {
        const domain = extractDomain(tab.url);
        if (domain) {
          onDomainChange(tabId, domain);
        }
      }
    }
  });
} 