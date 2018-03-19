"use strict";

var m = require("mithril");
var FuzzySort = require("fuzzysort");
var app = "com.dannyvankooten.browserpass";
var activeTab;
var searching = false;
var resultLogins = [];
var logins = [];
var fillOnSubmit = false;
var error;
var domain, urlDuringSearch;

m.mount(document.getElementById("mount"), { view: view, oncreate: oncreate });

chrome.tabs.onActivated.addListener(init);
chrome.tabs.query({ lastFocusedWindow: true, active: true }, function(tabs) {
  init(tabs[0]);
});

function view() {
  var results = "";

  if (searching) {
    results = m("div.loader");
  } else if (error) {
    results = m("div.status-text", "Error: " + error);
    error = undefined;
  } else if (logins) {
    if (logins.length === 0 && domain && domain.length > 0) {
      results = m(
        "div.status-text",
        m.trust(`No matching passwords found for <strong>${domain}</strong>.`)
      );
    } else if (logins.length > 0) {
      results = logins.map(function(login) {
        let selector = "button.login";
        let options = { onclick: getLoginData.bind(login) };

        let faviconUrl = getFaviconUrl(domain);
        if (faviconUrl) {
          selector += ".favicon";
          options.style = `background-image: url('${faviconUrl}')`;
        }

        return m("div.entry", [
          m(selector, options, login),
          m("button.copy.username", {
            onclick: loginToClipboard.bind({ entry: login, what: "username" }),
            tabindex: -1
          }),
          m("button.copy.password", {
            onclick: loginToClipboard.bind({ entry: login, what: "password" }),
            tabindex: -1
          })
        ]);
      });
    }
  }

  return m("div.container", { onkeydown: keyHandler }, [
    // search form
    m("div.search", [
      m(
        "form",
        {
          onsubmit: submitSearchForm,
          onkeydown: searchKeyHandler
        },
        [
          m("div", {
            "id": "filter-search"
          }),
          m("div", [
            m("input", {
              type: "text",
              id: "search-field",
              name: "s",
              placeholder: "Search passwords..",
              autocomplete: "off",
              autofocus: "on",
              oninput: filterLogins
            }),
            m("input", {
              type: "submit",
              value: "Search",
              style: "display: none;"
            })
          ])
        ]
      )
    ]),

    // results
    m("div.results", results)
  ]);
}

function filterLogins(e) {
  // use fuzzy search to filter results
  var filter = e.target.value.trim().split(/[\s\/]+/);
  if (filter.length > 0) {
    logins = resultLogins.slice(0);
    filter.forEach(function(word) {
      if (word.length > 0) {
        var refine = [];
        FuzzySort.go(word, logins, {allowTypo: false}).forEach(function(result) {
          refine.push(result.target);
        });
        logins = refine.slice(0);
      }
    });

    // fill login forms on submit rather than initiating a search
    fillOnSubmit = logins.length > 0;
  } else {
    // reset the result list if the filter is empty
    logins = resultLogins.slice(0);
  }

  // redraw the list
  m.redraw();

  // show / hide the filter hint
  showFilterHint(logins.length);
}

function searchKeyHandler(e) {
  // switch to search mode if backspace is pressed and no filter text has been entered
  if (e.code == "Backspace" && logins.length > 0 && e.target.value.length == 0) {
    e.preventDefault();
    logins = resultLogins = [];
    e.target.value = fillOnSubmit ? '' : domain;
    domain = '';
    showFilterHint(false);
  }
}

function showFilterHint(show=true) {
  var filterHint = document.getElementById("filter-search");
  var searchField = document.getElementById("search-field");
  if (show) {
    filterHint.style.display = "block";
    searchField.setAttribute("placeholder", "Refine search...");
  } else {
    filterHint.style.display = "none";
    searchField.setAttribute("placeholder", "Search passwords...");
  }
}

function submitSearchForm(e) {
  e.preventDefault();

  if (fillOnSubmit && logins.length > 0) {
    // fill using the first result
    getLoginData.bind(logins[0])();
  } else {
    // don't search without input.
    if (!this.s.value.length) {
      return;
    }

    // search for matching entries
    searchPassword(this.s.value, "search", false);
  }
}

function init(tab) {
  // do nothing if called from a non-tab context
  if (!tab || !tab.url) {
    return;
  }

  activeTab = tab;
  var activeDomain = parseDomainFromUrl(tab.url);
  searchPassword(activeDomain, "match_domain");
}

function searchPassword(_domain, action="search", useFillOnSubmit=true) {
  searching = true;
  logins = resultLogins = [];
  domain = _domain;
  urlDuringSearch = activeTab.url;
  m.redraw();

  // First get the settings needed by the browserpass native client
  // by requesting them from the background script (which has localStorage access
  // to the settings). Then construct the message to send to browserpass and
  // send that via sendNativeMessage.
  chrome.runtime.sendMessage(
    { action: "getSettings" },
    function(response) {
      chrome.runtime.sendNativeMessage(
        app,
        { action: action, domain: _domain, settings: response},
        function(response) {
          if (chrome.runtime.lastError) {
            error = chrome.runtime.lastError.message;
            console.error(error);
          }

          searching = false;
          logins = resultLogins = response ? response : [];
          document.getElementById("filter-search").textContent = domain;
          fillOnSubmit = useFillOnSubmit && logins.length > 0;
          if (logins.length > 0) {
            showFilterHint(true);
            document.getElementById("search-field").value = '';
          }
          m.redraw();
        }
      );
    });
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

  return null;
}

function getLoginData() {
  searching = true;
  logins = resultLogins = [];
  m.redraw();

  chrome.runtime.sendMessage(
    { action: "login", entry: this, urlDuringSearch: urlDuringSearch },
    function(response) {
      searching = false;
      fillOnSubmit = false;

      if (response.error) {
        error = response.error;
        m.redraw();
      } else {
        window.close();
      }
    }
  );
}

function loginToClipboard() {
  var what = this.what;
  chrome.runtime.sendNativeMessage(
    app,
    { action: "get", entry: this.entry },
    function(response) {
      if (chrome.runtime.lastError) {
        error = chrome.runtime.lastError.message;
        m.redraw();
      } else {
        if (what === "password") {
          toClipboard(response.p);
        } else if (what === "username") {
          toClipboard(response.u);
        }
        window.close();
      }
    }
  );
}

function toClipboard(s) {
  var clipboardContainer = document.getElementById("clipboard-container");
  var clipboard = document.createElement("input");
  clipboardContainer.appendChild(clipboard);
  clipboard.value = s;
  clipboard.select();
  document.execCommand("copy");
  clipboard.blur();
  clipboardContainer.removeChild(clipboard);
}

// This function uses regular DOM
// therefore there is no need for redraw calls
function keyHandler(e) {
  switch (e.key) {
    case "ArrowUp":
      switchFocus("div.entry:last-child > .login", "previousElementSibling");
      break;

    case "ArrowDown":
      switchFocus("div.entry:first-child > .login", "nextElementSibling");
      break;
    case "c":
      if (e.target.id != "search-field" && e.ctrlKey) {
        document.activeElement["nextElementSibling"][
          "nextElementSibling"
        ].click();
      }
      break;
    case "C":
      if (e.target.id != "search-field") {
        document.activeElement["nextElementSibling"].click();
      }
      break;
  }
}

function switchFocus(firstSelector, nextNodeAttr) {
  var searchField = document.getElementById("search-field");
  var newActive = searchField;

  if (document.activeElement === searchField) {
    newActive = document.querySelector(firstSelector);
  } else {
    let tmp = document.activeElement["parentElement"][nextNodeAttr];
    if (tmp !== null) {
      newActive = tmp["firstElementChild"];
    }
  }

  newActive.focus();
}

// The oncreate(vnode) hook is called after a DOM element is created and attached to the document.
// see https://mithril.js.org/lifecycle-methods.html#oncreate for mor informations
function oncreate() {
  // FireFox probably prevents `focus()` calls for some time
  // after extension is rendered.
  window.setTimeout(function() {
    document.getElementById("search-field").focus();
  }, 100);
}
