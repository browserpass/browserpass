function save_options() {
  var autoSubmit = document.getElementById("auto-submit").checked;
  localStorage.setItem("autoSubmit", autoSubmit);

  // Options related to fuzzy finding.
  //  use_fuzzy_search indicates if fuzzy finding or glob searching should
  //  be used in manual searches
  var use_fuzzy = document.getElementById("use-fuzzy").checked;
  localStorage.setItem("use_fuzzy_search", use_fuzzy);

  window.close();
}

function restore_options() {
  var autoSubmit = localStorage.getItem("autoSubmit") == "true";
  document.getElementById("auto-submit").checked = autoSubmit;

  // Restore the view to show the settings described above
  var use_fuzzy = localStorage.getItem("use_fuzzy_search") != "false";
  document.getElementById("use-fuzzy").checked = use_fuzzy;
}

document.addEventListener("DOMContentLoaded", restore_options);
document.getElementById("save").addEventListener("click", save_options);
