var otpInput = document.getElementById("otp");
var otpLabel = document.getElementById("label");
var otpCopy = document.getElementById("copy");
var otpDismiss = document.getElementById("dismiss");

window.addEventListener("message", receiveMessage, false);
function receiveMessage(event) {
  otpInput.value = event.data.digits;
  otpInput.setAttribute("size", event.data.digits.length);
  otpLabel.innerText = (event.data.label || "OTP") + ":";
  var message = {
    action: "resize",
    payload: {
      width: document.body.scrollWidth,
      height: document.body.scrollHeight
    }
  };
  window.parent.postMessage(message, "*");
}

window.onload = function() {
  window.parent.postMessage({ action: "load" }, "*");
}

document.body.onclick = function() {
  otpInput.select();
}

otpCopy.onclick = function() {
  otpInput.select();
  document.execCommand("copy");
}

otpDismiss.onclick = function() {
  chrome.runtime.sendMessage({ action: "dismissOTP" });
  window.parent.postMessage({ action: "dismiss" }, "*");
}
