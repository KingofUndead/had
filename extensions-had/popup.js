async function loadSavedSettings() {
  const result = await chrome.storage.local.get(["proxyHost", "proxyPort", "proxyScheme"]);
  if (result.proxyHost) document.getElementById("host").value = result.proxyHost;
  if (result.proxyPort) document.getElementById("port").value = result.proxyPort;
  if (result.proxyScheme) document.getElementById("scheme").value = result.proxyScheme;
}

async function saveSettings(host, port, scheme) {
  await chrome.storage.local.set({
    proxyHost: host,
    proxyPort: port,
    proxyScheme: scheme
  });
}

async function updateStatus() {
  const response = await chrome.runtime.sendMessage({ action: "getStatus" });
  const statusDiv = document.getElementById("status");
  if (response.enabled) {
    statusDiv.className = "status status-on";
    statusDiv.innerHTML = "🔥 ACTIVE 🔥";
  } else {
    statusDiv.className = "status status-off";
    statusDiv.innerHTML = "💀 INACTIVE 💀";
  }
}

document.getElementById("enableBtn").addEventListener("click", async () => {
  const host = document.getElementById("host").value.trim();
  const port = document.getElementById("port").value;
  const scheme = document.getElementById("scheme").value;
  
  if (!host || !port) {
    alert("[HAD] ERROR: Host & Port required!");
    return;
  }
  
  await saveSettings(host, port, scheme);
  await chrome.runtime.sendMessage({
    action: "setProxy",
    host: host,
    port: port,
    scheme: scheme
  });
  await updateStatus();
});

document.getElementById("disableBtn").addEventListener("click", async () => {
  await chrome.runtime.sendMessage({ action: "disableProxy" });
  await updateStatus();
});

loadSavedSettings();
updateStatus();
