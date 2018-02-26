function save_options() {
  var autoSubmit = document.getElementById("auto-submit").checked;
  localStorage.setItem("autoSubmit", autoSubmit);

  // Options related to fuzzy finding.
  //  use_fuzzy indicates if any fuzzy finding should be used, or if
  //  the glob search should be used when a user inputs a string to search
  var use_fuzzy = document.getElementById("use-fuzzy").checked;
  localStorage.setItem("use_fuzzy", use_fuzzy);

  window.close();
}

function restore_options() {
  var autoSubmit = localStorage.getItem("autoSubmit") == "true";
  document.getElementById("auto-submit").checked = autoSubmit;

  // Restore the view to show the settings described above
  var use_fuzzy = localStorage.getItem("use_fuzzy") == "true";
  document.getElementById("use-fuzzy").checked = use_fuzzy;
}

document.addEventListener("DOMContentLoaded", restore_options);
document.getElementById("save").addEventListener("click", save_options);
