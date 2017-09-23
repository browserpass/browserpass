"use strict";

var m = require("mithril");
var app = "com.dannyvankooten.browserpass";
var activeTab;
var searching = false;
var logins = null;
var highlightedItem = null;
var domain, urlDuringSearch;

m.mount(document.getElementById("mount"), { view: view });

chrome.browserAction.setIcon({ path: "icon-lock.svg" });
chrome.tabs.onActivated.addListener(init);
chrome.tabs.query({ currentWindow: true, active: true }, function(tabs) {
  init(tabs[0]);
});

function view() {
  var results = "";

  function style(faviconUrl, l, i) {
    var style = `background-image: url('${faviconUrl}');`;
    return i == highlightedItem ? style + ` background-color: yellow;` : style;
  }

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
      results = logins.map(function(l, i) {
        var faviconUrl = getFaviconUrl(domain);
        return m(
          "button.login",
          {
            onclick: getLoginData.bind(l),
            style: style(faviconUrl, l, i)
          },
          l
        );
      });
    }
  }

  return [
    // search form
    m("div.search", [
      m(
        "form",
        {
          onsubmit: submitSearchForm(results)
        },
        [
          m("input", {
            type: "text",
            name: "s",
            placeholder: "Search password..",
            autocomplete: "off",
            autofocus: "on",
            onkeydown: keyHandler.bind(results),
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
  ];
}

function submitSearchForm(resultsArray) {
  return function (e) {
    e.preventDefault();

    // don't search without input.
    if (!this.s.value.length) {

      // open highlighted item
      if (highlightedItem !== null) {
        var item = resultsArray[highlightedItem];
        if (item) {
          return getLoginData.call(item.children[0]);
        }
      }
      return;
    }

    searchPassword(this.s.value);
  }
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
  highlightedItem = null;
  m.redraw();

  chrome.runtime.sendMessage(
    { action: "login", entry: this, urlDuringSearch: urlDuringSearch },
    function(response) {
      searching = false;
      window.close();
    }
  );
}

function rotateHighlight(n) {
  // if items are empty
  if (this.length < 1) {
    return null;
  }

  // if nothing is highlighted at the moment
  if (highlightedItem === null) {
    return n > 0 ? 0 : this.length - 1;
  }

  var newHighlightedItem = highlightedItem + n;

  // if newHighlighted greater than number of items
  if (newHighlightedItem >= this.length) {
    return 0;
  }

  // if newHighlighter lesser than zero
  if (newHighlightedItem < 0) {
    return this.length - 1;
  }

  return newHighlightedItem;
}

function keyHandler(e) {
  switch (e.key) {
  case 'ArrowUp':
    highlightedItem = rotateHighlight.call(this, -1);
    m.redraw();
    break;
  case 'ArrowDown':
    highlightedItem = rotateHighlight.call(this, 1);
    m.redraw();
    break;
  }
}
