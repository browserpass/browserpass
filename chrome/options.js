function save_options() {
  let autoSubmit = document.getElementById("auto-submit").checked;
  localStorage.setItem("autoSubmit", autoSubmit);

  // Options related to fuzzy finding.
  //  use_fuzzy_search indicates if fuzzy finding or glob searching should
  //  be used in manual searches
  let use_fuzzy = document.getElementById("use-fuzzy").checked;
  localStorage.setItem("use_fuzzy_search", use_fuzzy);

  let groups = document.getElementsByClassName("password-store-group");
  let paths = [];
  for (let i = 0; i < groups.length; i++) {
    let group = groups[i];
    let enabled = group.querySelector(".password-store-path-enabled").checked;
    let path = group.querySelector(".password-store-path").value;
    if (path) {
      paths.push({ path: path, enabled: enabled });
    }
  }
  localStorage.setItem("paths", JSON.stringify(paths));

  window.close();
}

function restore_options() {
  let autoSubmit = localStorage.getItem("autoSubmit") == "true";
  document.getElementById("auto-submit").checked = autoSubmit;

  // Restore the view to show the settings described above
  let use_fuzzy = localStorage.getItem("use_fuzzy_search") != "false";
  document.getElementById("use-fuzzy").checked = use_fuzzy;

  let groups = document.getElementsByClassName("password-store-group");
  let paths = JSON.parse(localStorage.getItem("paths") || "[]");
  for (let i = 0; i < paths.length; i++) {
    let path = paths[i];
    let group = groups[i];
    group.querySelector(".password-store-path-enabled").checked = path.enabled;
    group.querySelector(".password-store-path").value = path.path;
  }
}

document.addEventListener("DOMContentLoaded", restore_options);
document.getElementById("save").addEventListener("click", save_options);
