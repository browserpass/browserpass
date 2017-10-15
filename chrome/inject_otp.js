window.browserpassDisplayOTP = function(login) {
  iframe = document.createElement("iframe");
  iframe.id = "browserpass-otp-iframe";
  iframe.src = chrome.runtime.getURL("otp.html");
  iframe.style = `
    position: fixed;
    top: 0;
    right: 0;
    background-color: white;
    border-bottom-left-radius: 4px;
    border-left: 1px solid #888888;
    border-bottom: 1px solid #888888;
    border-top: none;
    border-right: none;
    box-sizing: content-box;
    z-index: 1000000;
    display: none;
  `;

  window.addEventListener("message", receiveMessage, false);
  function receiveMessage(event) {
    if (event.data.action == "load") {
      iframe.contentWindow.postMessage(login, "*");
    }

    if (event.data.action == "resize") {
      iframe.style.display = "block";

      iframe.width = event.data.payload.width;
      iframe.height = event.data.payload.height;
    }

    if (event.data.action == "dismiss") {
      iframe.remove();
      window.removeEventListener("message", receiveMessage);
    }
  };

  var oldIframe = document.getElementById("browserpass-otp-iframe")
  if (oldIframe != null) {
    oldIframe.remove();
  }
  document.body.appendChild(iframe);
}
