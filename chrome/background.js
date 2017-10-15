"use strict";

var app = "com.dannyvankooten.browserpass";

var tabInfos = {};

chrome.runtime.onMessage.addListener(onMessage);
chrome.tabs.onUpdated.addListener(onTabUpdated);

// fill login form & submit
function fillLoginForm(login, tab) {
  const loginParam = JSON.stringify(login);
  const autoSubmitParam = JSON.stringify(localStorage.getItem("autoSubmit"));

  chrome.tabs.executeScript(
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
      chrome.tabs.executeScript(
        tabId,
        {
          code: `browserpassDisplayOTP(${tabInfos[tabId].login});`
        }
      );
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

        chrome.tabs.query({ currentWindow: true, active: true }, function(
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
