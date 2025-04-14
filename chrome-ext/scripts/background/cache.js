const CACHE_EXPIRY = 24 * 60 * 60; // 24 hours in seconds

export async function getCachedResult(domain) {
  const data = await chrome.storage.local.get(domain);
  if (!data[domain]) return null;

  const { result, timestamp } = data[domain];
  if (Math.floor(Date.now() / 1000) - timestamp > CACHE_EXPIRY) {
    await chrome.storage.local.remove(domain);
    return null;
  }

  return result;
}

export function setCachedResult(domain, result) {
  if (result.status === 'in_progress') {
    return;
  }
  
  chrome.storage.local.set({
    [domain]: {
      result,
      timestamp: Math.floor(Date.now() / 1000)
    }
  });
} 