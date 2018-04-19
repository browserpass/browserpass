"use strict";

var Tldjs = require("tldjs");

var app = "com.dannyvankooten.browserpass";

var tabInfos = {};
var authListeners = {};

chrome.runtime.onMessage.addListener(onMessage);
chrome.tabs.onUpdated.addListener(onTabUpdated);
chrome.runtime.onInstalled.addListener(onExtensionInstalled);

// fill login form & submit
function fillLoginForm(login, tab, autoSubmit) {
  const loginParam = JSON.stringify(login);
  chrome.tabs.executeScript(
    tab.id,
    {
      allFrames: true,
      file: "/inject.js"
    },
    function() {
      chrome.tabs.executeScript({
        allFrames: true,
        code: `browserpassFillForm(${loginParam}, ${autoSubmit});`
      });
    }
  );

  if (login.digits) {
    tabInfos[tab.id] = {
      login: loginParam,
      hostname: getHostname(tab.url)
    };
    displayOTP(tab.id);
  }
}

function displayOTP(tabId) {
  chrome.tabs.executeScript(
    tabId,
    {
      file: "/inject_otp.js"
    },
    function() {
      chrome.tabs.executeScript(tabId, {
        code: `browserpassDisplayOTP(${tabInfos[tabId].login});`
      });
    }
  );
}

function onMessage(request, sender, sendResponse) {
  switch (request.action) {
    case "login": {
      chrome.runtime.sendNativeMessage(
        app,
        { action: "get", entry: request.entry, settings: getSettings() },
        function(response) {
          if (chrome.runtime.lastError) {
            var error = chrome.runtime.lastError.message;
            console.error(error);
            sendResponse({ status: "ERROR", error: error });
            return;
          }

          if (typeof response == "string") {
            console.error(response);
            sendResponse({ status: "ERROR", error: response });
            return;
          }

          chrome.tabs.query({ lastFocusedWindow: true, active: true }, function(
            tabs
          ) {
            // do not send login data to page if URL changed during search.
            if (tabs[0].url == request.urlDuringSearch) {
              var autoSubmit = response.hasOwnProperty('autoSubmit') ? response.autoSubmit : getSettings().autoSubmit;
              fillLoginForm(response, tabs[0], autoSubmit);
            }
          });

          sendResponse({ status: "OK" });
        }
      );

      // Must return true when sendResponse is being called asynchronously
      return true;
    }

    case "copyToClipboard": {
      chrome.runtime.sendNativeMessage(
        app,
        { action: "get", entry: request.entry, settings: getSettings() },
        function(response) {
          if (chrome.runtime.lastError) {
            var error = chrome.runtime.lastError.message;
            console.error(error);
            sendResponse({ status: "ERROR", error: error });
            return;
          }

          if (typeof response == "string") {
            console.error(response);
            sendResponse({ status: "ERROR", error: response });
            return;
          }

          let text = "";
          if (request.what === "password") {
            text = response.p;
          } else if (request.what === "username") {
            text = response.u;
          } else if (request.what === "otp") {
            if (!response.digits) {
              sendResponse({status: "ERROR", error:
                            "Unable to determine the OTP code for this entry."
                           });
              return;
            }
            text = response.digits;
          }

          try {
            sendResponse({ status: "OK", text: text });
          } catch (e) {
            // If unable to send text to the popup to let it copy the text to clipboard,
            // try to copy to clipboard from the background page itself.
            // This only works in Chrome. See #241
            copyToClipboard(text);
          }
        }
      );

      // Must return true when sendResponse is being called asynchronously
      return true;
    }

    case "dismissOTP": {
      if (request.action == "dismissOTP" && sender.tab.id in tabInfos) {
        delete tabInfos[sender.tab.id];
      }
      break;
    }

    // allows the local communication to request settings. Returns an
    // object that has current settings. Update this as new settings
    // are added (or old ones removed)
    case "getSettings": {
      sendResponse(getSettings());
      break;
    }

    // spawn a new tab with credentials from the password file
    case "launch": {
      chrome.runtime.sendNativeMessage(
        app,
        { action: "get", entry: request.entry, settings: getSettings() },
        function(response) {
          if (chrome.runtime.lastError) {
            var error = chrome.runtime.lastError.message;
            console.error(error);
            sendResponse({ status: "ERROR", error: error });
            return;
          }

          if (typeof response == "string") {
            console.error(response);
            sendResponse({ status: "ERROR", error: response });
            return;
          }

          if (!response.hasOwnProperty("url") || response.url.length == 0) {
            // guess url from login path if not available in the host app response
            response.url = parseUrlFromEntry(request.entry);

            // if url is not available at this point, send an error
            if (!response.hasOwnProperty("url") || response.url.length == 0) {
              sendResponse({
                status: "ERROR",
                error:
                  "Unable to determine the URL for this entry. If you have defined one in the password file, " +
                  "your host application must be at least v2.0.14 for this to be usable."
              });
              return;
            }
          }

          var url = response.url.match(/^([a-z]+:)?\/\//i)
            ? response.url
            : "http://" + response.url;

          var tabsMethod = request.openInNewTab ? "create" : "update";
          chrome.tabs[tabsMethod]({ url: url }, function(tab) {
            var authAttempted = false;

            authListeners[tab.id] = function(requestDetails) {
              // only supply credentials if this is the first time for this tab, and the tab is not loaded
              if (authAttempted) {
                return {};
              }
              authAttempted = true;
              return onAuthRequired(url, requestDetails, response);
            };

            // intercept requests for authentication
            chrome.webRequest.onAuthRequired.addListener(
              authListeners[tab.id],
              { urls: ["*://*/*"], tabId: tab.id },
              ["blocking"]
            );
          });

          sendResponse({ status: "OK" });
        }
      );

      // Must return true when sendResponse is being called asynchronously
      return true;
    }
  }
}

function parseUrlFromEntry(entry) {
  var parts =
    entry.indexOf(":") > 0 ? entry.substr(entry.indexOf(":") + 1) : entry;
  parts = parts.split(/\//).reverse();
  for (var i in parts) {
    var part = parts[i];
    var info = Tldjs.parse(part);
    if (
      info.isValid &&
      info.tldExists &&
      info.domain !== null &&
      info.hostname === part
    ) {
      return part;
    }
  }
  return "";
}

function copyToClipboard(s) {
  var clipboard = document.createElement("input");
  document.body.appendChild(clipboard);
  clipboard.value = s;
  clipboard.select();
  document.execCommand("copy");
  clipboard.blur();
  document.body.removeChild(clipboard);
}

function getSettings() {
  // default settings
  var settings = {
    autoSubmit: false,
    use_fuzzy_search: true,
    customStores: []
  };

  // load settings from local storage
  for (var key in settings) {
    var value = localStorage.getItem(key);
    if (value !== null) {
      settings[key] = JSON.parse(value);
    }
  }

  // filter custom stores by enabled & path length, and ensure they are named
  settings.customStores = settings.customStores
    .filter(store => store.enabled && store.path.length > 0)
    .map(function(store) {
      if (!store.name) {
        store.name = store.path;
      }
      return store;
    });

  return settings;
}

// listener function for authentication interception
function onAuthRequired(url, requestDetails, response) {
  // ask the user before sending credentials to a different domain
  var launchHost = url.match(/:\/\/([^\/]+)/)[1];
  if (launchHost !== requestDetails.challenger.host) {
    var message =
      "You are about to send login credentials to a domain that is different than " +
      "the one you lauched from the browserpass extension. Do you wish to proceed?\n\n" +
      "Launched URL: " +
      url +
      "\n" +
      "Authentication URL: " +
      requestDetails.url;
    if (!confirm(message)) {
      return {};
    }
  }

  // ask the user before sending credentials over an insecure connection
  if (!requestDetails.url.match(/^https:/i)) {
    var message =
      "You are about to send login credentials via an insecure connection!\n\n" +
      "Are you sure you want to do this? If there is an attacker watching your " +
      "network traffic, they may be able to see your username and password.\n\n" +
      "URL: " +
      requestDetails.url;
    if (!confirm(message)) {
      return {};
    }
  }

  // supply credentials
  return {
    authCredentials: {
      username: response.u,
      password: response.p
    }
  };
}

function onTabUpdated(tabId, info, tab) {
  if (info.url && tabId in tabInfos) {
    if (getHostname(info.url) != tabInfos[tabId].hostname) {
      delete tabInfos[tabId];
    }
  }

  if (info.status != "complete") {
    return;
  }

  if (tabId in tabInfos) {
    displayOTP(tabId);
  }

  if (tabId in authListeners) {
    chrome.webRequest.onAuthRequired.removeListener(authListeners[tabId]);
    delete authListeners[tabId];
  }
}

function getHostname(url) {
  // Manipulate the browser into parsing the URL for us
  var a = document.createElement("a");
  a.href = url;
  return a.hostname;
}

function onExtensionInstalled(details) {
  // No permissions
  if (!chrome.notifications) {
    return;
  }

  if (details.reason != "update") {
    return;
  }

  var changelog = {
    2012: "Breaking change: please update the host app to at least v2.0.12"
  };

  var parseVersion = version => parseInt(version.replace(/\./g, ""));
  var newVersion = parseVersion(chrome.runtime.getManifest().version);
  var prevVersion = parseVersion(details.previousVersion);

  Object.keys(changelog)
    .sort()
    .forEach(function(version) {
      if (version > prevVersion && version <= newVersion) {
        chrome.notifications.create(version, {
          title: "browserpass: Important changes",
          message: changelog[version],
          iconUrl: "icon-lock.png",
          type: "basic"
        });
      }
    });
}
