var settings = {
  autoSubmit: {
    type: "checkbox",
    title: "Automatically submit forms after filling",
    value: false
  },
  use_fuzzy_search: {
    type: "checkbox",
    title: "Use fuzzy search",
    value: true
  },
  customStores: {
    title: "Custom password store locations",
    value: [{ enabled: true, name: "", path: "" }]
  }
};

// load settings & create render tree
loadSettings();
var tree = {
  view: function() {
    var nodes = [m("h3", "Basic Settings")];
    for (var key in settings) {
      var type = settings[key].type;
      if (type == "checkbox") {
        nodes.push(createCheckbox(key, settings[key]));
      }
    }
    nodes.push(m("h3", "Custom Store Locations"));
    for (var key in settings.customStores.value) {
      nodes.push(createCustomStore(key, settings.customStores.value[key]));
    }
    nodes.push(
      m(
        "button.add-store",
        {
          onclick: function() {
            settings.customStores.value.push({
              enabled: true,
              name: "",
              path: ""
            });
            saveSetting("customStores");
          }
        },
        "Add Store"
      )
    );
    return nodes;
  }
};

// attach tree
var m = require("mithril");
m.mount(document.body, tree);

// load settings from local storage
function loadSettings() {
  for (var key in settings) {
    var value = localStorage.getItem(key);
    if (value !== null) {
      settings[key].value = JSON.parse(value);
    }
  }
}

// save settings to local storage
function saveSetting(name) {
  var value = settings[name].value;
  if (Array.isArray(value)) {
    value = value.filter(item => item !== null);
  }
  value = JSON.stringify(value);
  localStorage.setItem(name, value);
}

// create a checkbox option
function createCheckbox(name, option) {
  return m("div.option", { class: name }, [
    m("input[type=checkbox]", {
      name: name,
      title: option.title,
      checked: option.value,
      onchange: function(e) {
        settings[name].value = e.target.checked;
        saveSetting(name);
      }
    }),
    m("label", { for: name }, option.title)
  ]);
}

// create a custom store option
function createCustomStore(key, store) {
  return m("div.option.custom-store", { class: "store-" + store.name }, [
    m("input[type=checkbox]", {
      title: "Whether to enable this password store",
      checked: store.enabled,
      onchange: function(e) {
        store.enabled = e.target.checked;
        saveSetting("customStores");
      }
    }),
    m("input[type=text].name", {
      title: "The name for this password store",
      value: store.name,
      placeholder: "name",
      onchange: function(e) {
        store.name = e.target.value;
        saveSetting("customStores");
      }
    }),
    m("input[type=text].path", {
      title: "The full path to this password store",
      value: store.path,
      placeholder: "/path/to/store",
      onchange: function(e) {
        store.path = e.target.value;
        saveSetting("customStores");
      }
    }),
    m(
      "a.remove",
      {
        title: "Remove this password store",
        onclick: function() {
          delete settings.customStores.value[key];
          saveSetting("customStores");
        }
      },
      "[X]"
    )
  ]);
}
