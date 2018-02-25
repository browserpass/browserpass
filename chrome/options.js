function save_options() {
  var autoSubmit = document.getElementById("auto-submit").checked;
  localStorage.setItem("autoSubmit", autoSubmit);

  // Options related to fuzzy finding.
  //  use_fuzzy indicates if any fuzzy finding should be used, or if
  //  the glob search should be used when a user inputs a string to search
  var use_fuzzy = document.getElementById("use-fuzzy").checked;
  localStorage.setItem("use_fuzzy", use_fuzzy);

  // fuzzy_algorithm - since there are two libraries being evaluated
  // provide a way to select between them withoug needing rebuild
  // of the browserpass native client.
  var fuzzy_algorithm = document.querySelector('input[name="fuzzyalgo"]:checked').value;
  localStorage.setItem("fuzzy_algorithm", fuzzy_algorithm);

  window.close();
}

function restore_options() {
  var autoSubmit = localStorage.getItem("autoSubmit") == "true";
  document.getElementById("auto-submit").checked = autoSubmit;

  // Restore the view to show the settings described above
  var use_fuzzy = localStorage.getItem("use_fuzzy") == "true";
  document.getElementById("use-fuzzy").checked = use_fuzzy;

  var fuzzy_algorithm = localStorage.getItem("fuzzy_algorithm");
  if (fuzzy_algorithm != null) {
    document.getElementsByName("fuzzyalgo").forEach(function(elem) {
      if (elem.value == fuzzy_algorithm) {
        elem.checked = true;
      } else {
        elem.checked = false;
      }
    });
  }
}

document.addEventListener("DOMContentLoaded", restore_options);
document.getElementById("save").addEventListener("click", save_options);
