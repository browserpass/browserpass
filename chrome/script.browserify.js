"use strict";

var m = require("mithril");
var app = "com.dannyvankooten.browserpass";
var activeTab;
var searching = false;
var logins = null;
var domain, urlDuringSearch;

m.mount(document.getElementById("mount"), { view: view, oncreate: oncreate });

chrome.browserAction.setIcon({ path: "icon-lock.svg" });
chrome.tabs.onActivated.addListener(init);
chrome.tabs.query({ currentWindow: true, active: true }, function(tabs) {
  init(tabs[0]);
});

function view() {
  var results = "";

  if (searching) {
    results = m("div.loader");
  } else if (logins !== null) {
    if (typeof logins === "undefined") {
      results = m("div.status-text", "Error talking to Browserpass host");
    } else if (logins.length === 0) {
      results = m(
        "div.status-text",
        m.trust(`No passwords found for <strong>${domain}</strong>.`)
      );
    } else if (logins.length > 0) {
      results = logins.map(function(l) {
        var faviconUrl = getFaviconUrl(domain);
        return m(
          "button.login",
          {
            onclick: getLoginData.bind(l),
            style: `background-image: url('${faviconUrl}')`
          },
          l
        );
      });
    }
  }

  return m("div.container", { onkeydown: keyHandler }, [
    // search form
    m("div.search", [
      m(
        "form",
        {
          onsubmit: submitSearchForm
        },
        [
          m("input", {
            type: "text",
            id: "search-field",
            name: "s",
            placeholder: "Search password..",
            autocomplete: "off",
            autofocus: "on"
          }),
          m("input", {
            type: "submit",
            value: "Search",
            style: "display: none;"
          })
        ]
      )
    ]),

    // results
    m("div.results", results)
  ]);
}

function submitSearchForm(e) {
  e.preventDefault();

  // don't search without input.
  if (!this.s.value.length) {
    return;
  }

  searchPassword(this.s.value);
}

function init(tab) {
  // do nothing if called from a non-tab context
  if (!tab || !tab.url) {
    return;
  }

  activeTab = tab;
  var activeDomain = parseDomainFromUrl(tab.url);
  searchPassword(activeDomain);
}

function searchPassword(_domain) {
  searching = true;
  logins = null;
  domain = _domain;
  urlDuringSearch = activeTab.url;
  m.redraw();

  chrome.runtime.sendNativeMessage(
    app,
    { action: "search", domain: _domain },
    function(response) {
      if (chrome.runtime.lastError) {
        console.log(chrome.runtime.lastError);
      }

      searching = false;
      logins = response;
      m.redraw();
    }
  );
}

function parseDomainFromUrl(url) {
  var a = document.createElement("a");
  a.href = url;
  return a.hostname;
}

function getFaviconUrl(domain) {
  // use current favicon when searching for current tab
  if (
    activeTab &&
    activeTab.favIconUrl &&
    activeTab.favIconUrl.indexOf(domain) > -1
  ) {
    return activeTab.favIconUrl;
  }

  return "icon-key.svg";
}

function getLoginData() {
  searching = true;
  logins = null;
  m.redraw();

  chrome.runtime.sendMessage(
    { action: "login", entry: this, urlDuringSearch: urlDuringSearch },
    function(response) {
      searching = false;
      window.close();
    }
  );
}

// This function uses regular DOM
// therefore there is no need for redraw calls
function keyHandler(e) {
  switch (e.key) {
    case "ArrowUp":
      switchFocus("button.login:last-child", "previousElementSibling");
      break;

    case "ArrowDown":
      switchFocus("button.login:first-child", "nextElementSibling");
      break;
  }
}

function switchFocus(firstSelector, nextNodeAttr) {
  var searchField = document.getElementById('search-field');
  var newActive = document.activeElement === searchField
      ? document.querySelector(firstSelector)
      : document.activeElement[nextNodeAttr];

  if (newActive) {
    newActive.focus();
  } else {
    searchField.focus();
  }
}

// The oncreate(vnode) hook is called after a DOM element is created and attached to the document.
// see https://mithril.js.org/lifecycle-methods.html#oncreate for mor informations
function oncreate() {
  var searchField = document.getElementById('search-field');

  // FireFox probably prevents `focus()` calls for some time
  // after extension is rendered.
  searchField.focus();
  window.setTimeout(function() {
    if (searchField) {
      searchField.focus();
    }
  }, 100);
}
