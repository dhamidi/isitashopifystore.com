const CACHE_EXPIRY = 24 * 60 * 60 * 1000; // 24 hours

export async function getCachedResult(domain) {
  const data = await chrome.storage.local.get(domain);
  if (!data[domain]) return null;

  const { result, timestamp } = data[domain];
  if (Date.now() - timestamp > CACHE_EXPIRY) {
    await chrome.storage.local.remove(domain);
    return null;
  }

  return result;
}

export async function setCachedResult(domain, result) {
  await chrome.storage.local.set({
    [domain]: {
      result,
      timestamp: Date.now()
    }
  });
} 