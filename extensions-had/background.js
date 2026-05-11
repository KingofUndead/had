let proxyEnabled = false;

async function setProxy(host, port, scheme) {
  const config = {
    mode: "fixed_servers",
    rules: {
      singleProxy: {
        scheme: scheme || "http",
        host: host,
        port: parseInt(port)
      },
      bypassList: ["localhost", "127.0.0.1", "::1"]
    }
  };
  
  chrome.proxy.settings.set(
    { value: config, scope: "regular" },
    () => console.log("[HAD] Proxy activated on:", host + ":" + port)
  );
}

async function disableProxy() {
  chrome.proxy.settings.clear(
    { scope: "regular" },
    () => console.log("[HAD] Proxy disabled - The darkness fades")
  );
}

chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
  if (request.action === "setProxy") {
    setProxy(request.host, request.port, request.scheme);
    proxyEnabled = true;
    sendResponse({ status: "ok" });
  }
  else if (request.action === "disableProxy") {
    disableProxy();
    proxyEnabled = false;
    sendResponse({ status: "ok" });
  }
  else if (request.action === "getStatus") {
    sendResponse({ enabled: proxyEnabled });
  }
  return true;
});

chrome.runtime.onInstalled.addListener(async () => {
  const result = await chrome.storage.local.get(["proxyHost", "proxyPort", "proxyScheme"]);
  if (result.proxyHost && result.proxyPort) {
    setProxy(result.proxyHost, result.proxyPort, result.proxyScheme);
    proxyEnabled = true;
  }
});
