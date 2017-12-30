function save_options() {
  var autoSubmit = document.getElementById("auto-submit").checked;
  localStorage.setItem("autoSubmit", autoSubmit);
  window.close();
}

function restore_options() {
  var autoSubmit = localStorage.getItem("autoSubmit");
  if (autoSubmit == "false") {
    autoSubmit = false;
  } else {
    // Load the default
    autoSubmit = false;
  }
  document.getElementById("auto-submit").checked = autoSubmit;
}

document.addEventListener("DOMContentLoaded", restore_options);
document.getElementById("save").addEventListener("click", save_options);
