'use strict';

var m = require('mithril');
var parseDomain = require('parse-domain');
var app = 'com.dannyvankooten.browserpass';
var activeTab;
var searching = false;
var logins = null;
var domain, urlDuringSearch;

m.mount(document.getElementById('mount'), { "view": view });

chrome.browserAction.setIcon({ path: 'icon-lock.svg' });
chrome.tabs.onActivated.addListener(init);
chrome.tabs.query({ currentWindow: true, active: true }, function (tabs) {
  init(tabs[0]);
});

function view() {
  var results = '';

  if( searching ) {
    results = m('div.loader');
  } else if( logins !== null ) {
    if( typeof(logins) === "undefined" ) {
      results = m('div.status-text', "Error talking to Browserpass host");
    } else if( logins.length === 0 ) {
      results = m('div.status-text',  m.trust(`No passwords found for <strong>${domain}</strong>.`));
    } else if( logins.length > 0 ) {
        results = logins.map(function(l) {
          var faviconUrl = getFaviconUrl(domain);
          return m('button.login', {
            "onclick": getLoginData.bind(l),
            "style": `background-image: url('${faviconUrl}')`
          },
          l)
        });
    }
  }

  return [
    // search form
    m('div.search', [
      m('form', {
        "onsubmit": submitSearchForm
      }, [
        m('input', {
          "type": "text",
          "name": "s",
          "placeholder": "Search password..",
          "autocomplete": "off",
          "autofocus": "on"
        }),
        m('input', {
          "type": "submit",
          "value": "Search",
          "style": "display: none;"
        })
      ])
    ]),

    // results
    m('div.results', results)
  ];
}

function submitSearchForm(e) {
  e.preventDefault();

  // don't search without input.
  if( ! this.s.value.length ) {
      return;
  }

  searchPassword(this.s.value);
}

function init(tab) {
  // do nothing if called from a non-tab context
  if( ! tab || ! tab.url ) {
    return;
  }

  activeTab = tab;
  var parsedDomain = parseDomain(tab.url, {
    customTlds:/localhost/,
  });

  if( parsedDomain ) {
    var searchDomain = [parsedDomain.domain, parsedDomain.tld]
      .filter(function (x) { return x; })
      .join('.');

    if( searchDomain ) {
      searchPassword(searchDomain);
    }
  }
}

function searchPassword(_domain) {
  searching = true;
  logins = null;
  domain = _domain;
  urlDuringSearch = activeTab.url;
  m.redraw();

  chrome.runtime.sendNativeMessage(app, { "action": "search", "domain": _domain }, function(response) {
    if( chrome.runtime.lastError ) {
      console.log(chrome.runtime.lastError);
    }

    searching = false;
    logins = response;
    m.redraw();
  });
}

function getFaviconUrl(domain){

  // use current favicon when searching for current tab
  if(activeTab && activeTab.favIconUrl && activeTab.favIconUrl.indexOf(domain) > -1) {
    return activeTab.favIconUrl;
  }

  return 'icon-key.png';
}

function getLoginData() {
  searching = true;
  logins = null;
  m.redraw();

  chrome.runtime.sendNativeMessage(app, { "action": "get", "entry": this }, function(response) {
    if( chrome.runtime.lastError) {
      console.log(chrome.runtime.lastError);
    }

    searching = false;
    fillLoginForm(response);
    if (response.digits !== "") {
      alert(response.digits)
    }
  });
}

// fill login form & submit
function fillLoginForm(login) {
  // do not send login data to page if URL changed during search.
  if( activeTab.url != urlDuringSearch ) {
    return false;
  }

  var code = `
  (function(d) {
    function form() {
      return d.querySelector('input[type=password]').form || document.createElement('form');
    }

    function field(selector) {
      return form().querySelector(selector) || document.createElement('input');
    }

    function update(el, value) {
      if( ! value.length ) {
        return false;
      }

      el.setAttribute('value', value);
      el.value = value;

      var eventNames = [ 'click', 'focus', 'keyup', 'keydown', 'change', 'blur' ];
      eventNames.forEach(function(eventName) {
        el.dispatchEvent(new Event(eventName, {"bubbles":true}));
      });
      return true;
    }

    update(field('input[type=password]'), ${JSON.stringify(login.p)});
    update(field('input[type=email], input[type=text]'), ${JSON.stringify(login.u)});

    var password_inputs = document.querySelectorAll('input[type=password]');
    if (password_inputs.length > 1) {
      password_inputs[1].select();
    } else {
      field('[type=submit]').click();
    }
  })(document);
  `;
  chrome.tabs.executeScript({ code: code });
  window.close();
}
