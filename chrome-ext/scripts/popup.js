document.addEventListener('DOMContentLoaded', function() {
  const clearCacheButton = document.getElementById('clearCache');

  clearCacheButton.addEventListener('click', function() {
    chrome.runtime.sendMessage({ action: 'clearCache' }, function(response) {
      if (chrome.runtime.lastError) {
        console.error('Error sending message:', chrome.runtime.lastError);
        // Optionally display an error message to the user
      } else {
        console.log('Cache clear message sent, response:', response);
        // Optionally display a success message
        window.close(); // Close the popup after action
      }
    });
  });
}); 