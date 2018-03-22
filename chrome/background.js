"use strict";

var app = "com.dannyvankooten.browserpass";

var tabInfos = {};

chrome.runtime.onMessage.addListener(onMessage);
chrome.tabs.onUpdated.addListener(onTabUpdated);
chrome.runtime.onInstalled.addListener(onExtensionInstalled);

// fill login form & submit
function fillLoginForm(login, tab) {
  const loginParam = JSON.stringify(login);
  const autoSubmit = localStorage.getItem("autoSubmit");
  const autoSubmitParam = autoSubmit == "true";
  if (autoSubmit === null) {
    localStorage.setItem("autoSubmit", autoSubmitParam);
  }

  chrome.tabs.executeScript(
    tab.id,
    {
      allFrames: true,
      file: "/inject.js"
    },
    function() {
      chrome.tabs.executeScript({
        allFrames: true,
        code: `browserpassFillForm(${loginParam}, ${autoSubmitParam});`
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
  if (request.action == "login") {
    chrome.runtime.sendNativeMessage(
      app,
      { action: "get", entry: request.entry },
      function(response) {
        if (chrome.runtime.lastError) {
          var error = chrome.runtime.lastError.message;
          console.error(error);
          sendResponse({ status: "ERROR", error: error });
        }

        chrome.tabs.query({ lastFocusedWindow: true, active: true }, function(
          tabs
        ) {
          // do not send login data to page if URL changed during search.
          if (tabs[0].url == request.urlDuringSearch) {
            fillLoginForm(response, tabs[0]);
          }
        });

        sendResponse({ status: "OK" });
      }
    );

    // Need to return true if we are planning to sendResponse asynchronously
    return true;
  }

  if (request.action == "dismissOTP" && sender.tab.id in tabInfos) {
    delete tabInfos[sender.tab.id];
  }

  // allows the local communication to request settings. Returns an
  // object that has current settings. Update this as new settings
  // are added (or old ones removed)
  if (request.action == "getSettings") {
    const use_fuzzy_search =
      localStorage.getItem("use_fuzzy_search") != "false";
    sendResponse({ use_fuzzy_search: use_fuzzy_search });
  }

  // spawn a new tab with pre-provided credentials
  if (request.action == "launch") {
    chrome.tabs.create({ url: request.url }, function(tab) {
      var tabLoaded = false;
      var authAttempted = false;

      // listener function for authentication interception
      function authListener(requestDetails) {
        // only supply credentials if this is the first time for this tab, and the tab is not loaded
        if (authAttempted) {
          return {};
        }
        authAttempted = true;

        // don't supply credentials for loaded tabs
        if (tabLoaded) {
          return {};
        }

        // ask the user before sending credentials to a different domain
        var launchHost = request.url.match(/:\/\/([^\/]+)/)[1];
        if (launchHost !== requestDetails.challenger.host) {
          var message =
            "You are about to send login credentials to a domain that is different than " +
            "the one you lauched from the browserpass extension. Do you wish to proceed?\n\n" +
            "Launched URL: " +
            request.url +
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
            username: request.username,
            password: request.password
          }
        };
      }

      // intercept requests for authentication
      chrome.webRequest.onAuthRequired.addListener(
        authListener,
        { urls: ["*://*/*"], tabId: tab.id },
        ["blocking"]
      );

      // notice when the tab has been loaded
      chrome.tabs.onUpdated.addListener(function tabCompleteListener(
        tabId,
        info
      ) {
        if (tabId == tab.id && info.status == "complete") {
          // remove listeners
          chrome.tabs.onUpdated.removeListener(tabCompleteListener);
          chrome.webRequest.onAuthRequired.removeListener(authListener);

          // mark tab as loaded & wipe credentials
          tabLoaded = true;
          request.username = request.password = null;
        }
      });
    });
  }
}

function onTabUpdated(tabId, changeInfo, tab) {
  if (changeInfo.url && tabId in tabInfos) {
    if (getHostname(changeInfo.url) != tabInfos[tabId].hostname) {
      delete tabInfos[tabId];
    }
  }

  if (changeInfo.status == "complete" && tabId in tabInfos) {
    displayOTP(tabId);
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
