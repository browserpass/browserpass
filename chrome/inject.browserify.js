window.browserpassFillForm = function(login, autoSubmit) {
  const USERNAME_FIELDS =
    "input[id*=user i], input[id*=login i], input[id*=email i], input[type=email i], input[type=text i]";
  const PASSWORD_FIELDS = "input[type=password i]";

  function queryAllVisible(parent, selector, form) {
    var result = [];
    var selectors = selector.split(",");
    for (var i = 0; i < selectors.length; i++) {
      var selector = selectors[i].trim();
      var elems = parent.querySelectorAll(selector);
      for (var j = 0; j < elems.length; j++) {
        // Elem or its parent has a style 'display: none',
        // or it is just too narrow to be a real field (a trap for spammers?).
        if (elems[j].offsetWidth < 50 || elems[j].offsetHeight < 10) {
          continue;
        }
        // Select only elements from specified form
        if (form && form != elems[j].form) {
          continue;
        }
        var style = window.getComputedStyle(elems[j]);
        // Elem takes space on the screen, but it or its parent is hidden with a visibility style.
        if (style.visibility == "hidden") {
          continue;
        }
        // This element is visible, will use it.
        result.push(elems[j]);
      }
    }
    return result;
  }

  function queryFirstVisible(parent, selector, form) {
    var elems = queryAllVisible(parent, selector, form);
    return elems.length > 0 ? elems[0] : undefined;
  }

  function form() {
    var field = queryFirstVisible(
      document,
      PASSWORD_FIELDS + ", " + USERNAME_FIELDS,
      undefined
    );
    return field && field.form ? field.form : undefined;
  }

  function field(selector) {
    return (
      queryFirstVisible(document, selector, form()) ||
      document.createElement("input")
    );
  }

  function update(el, value) {
    if (!value.length) {
      return false;
    }
    el.setAttribute("value", value);
    el.value = value;
    var eventNames = [
      "click",
      "focus",
      "keypress",
      "keydown",
      "keyup",
      "input",
      "blur",
      "change"
    ];
    eventNames.forEach(function(eventName) {
      el.dispatchEvent(new Event(eventName, { bubbles: true }));
    });
    return true;
  }

  update(field(PASSWORD_FIELDS), login.p);
  update(field(USERNAME_FIELDS), login.u);

  if (login.digits) {
    alert((login.label || "OTP") + ": " + login.digits);
  }

  var password_inputs = queryAllVisible(document, PASSWORD_FIELDS, form());
  if (password_inputs.length > 1) {
    // There is likely a field asking for OTP code, so do not submit form just yet
    password_inputs[1].select();
  } else {
    window.requestAnimationFrame(function() {
      if (autoSubmit == "false") {
        field("[type=submit]").focus();
      } else {
        field("[type=submit]").click();
      }
    });
  }
};
