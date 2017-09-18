(function(d) {
  const USERNAME_FIELDS =
    "input[id=username i], input[id=user_name i], input[id=userid i], input[id=user_id i], input[id=login i], input[id=email i], input[type=email i], input[type=text i]";
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
      d,
      PASSWORD_FIELDS + ", " + USERNAME_FIELDS,
      undefined
    );
    return field && field.form ? field.form : undefined;
  }

  function field(selector) {
    return (
      queryFirstVisible(d, selector, form()) || document.createElement("input")
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

  if (typeof form() === "undefined") {
    return;
  }

  update(field(PASSWORD_FIELDS), login.p);
  update(field(USERNAME_FIELDS), login.u);

  if (login.digits) {
    alert(login.digits);
  }

  var password_inputs = queryAllVisible(form(), PASSWORD_FIELDS);
  if (autoSubmit == "false" || password_inputs.length > 1) {
    password_inputs[1].select();
  } else {
    window.requestAnimationFrame(function() {
      field("[type=submit]").click();
    });
  }
})(document);
