#!/usr/bin/env bash

set -e

assert_file_exists() {
  if [ ! -f "$1" ]; then
    echo "ERROR: '$1' is missing."
    echo "If you are running './install.sh' from a release archive, please file a bug."
    echo "If you are running './install.sh' from the source code, make sure to follow CONTRIBUTING.md on how to build first."
    exit 1
  fi
}

BIN_DIR="$( cd "$( dirname "$0" )" && pwd )"
JSON_DIR="$BIN_DIR"
APP_NAME="com.dannyvankooten.browserpass"
HOST_FILE="$BIN_DIR/browserpass"
BROWSER="$1"

# Find target dirs for various browsers & OS'es
# https://developer.chrome.com/extensions/nativeMessaging#native-messaging-host-location
# https://wiki.mozilla.org/WebExtensions/Native_Messaging
OPERATING_SYSTEM=$(uname -s)

case $OPERATING_SYSTEM in
Linux)
  HOST_FILE="$BIN_DIR/browserpass-linux64"
  if [ "$(whoami)" == "root" ]; then
    TARGET_DIR_CHROME="/etc/opt/chrome/native-messaging-hosts"
    TARGET_DIR_CHROMIUM="/etc/chromium/native-messaging-hosts"
    TARGET_DIR_FIREFOX="/usr/lib/mozilla/native-messaging-hosts"
    TARGET_DIR_VIVALDI="/etc/chromium/native-messaging-hosts"
  else
    TARGET_DIR_CHROME="$HOME/.config/google-chrome/NativeMessagingHosts"
    TARGET_DIR_CHROMIUM="$HOME/.config/chromium/NativeMessagingHosts"
    TARGET_DIR_FIREFOX="$HOME/.mozilla/native-messaging-hosts"
    TARGET_DIR_VIVALDI="$HOME/.config/vivaldi/NativeMessagingHosts"
  fi
  ;;
Darwin)
  HOST_FILE="$BIN_DIR/browserpass-darwinx64"
  if [ "$(whoami)" == "root" ]; then
    TARGET_DIR_CHROME="/Library/Google/Chrome/NativeMessagingHosts"
    TARGET_DIR_CHROMIUM="/Library/Application Support/Chromium/NativeMessagingHosts"
    TARGET_DIR_FIREFOX="/Library/Application Support/Mozilla/NativeMessagingHosts"
    TARGET_DIR_VIVALDI="/Library/Application Support/Vivaldi/NativeMessagingHosts"
  else
    TARGET_DIR_CHROME="$HOME/Library/Application Support/Google/Chrome/NativeMessagingHosts"
    TARGET_DIR_CHROMIUM="$HOME/Library/Application Support/Chromium/NativeMessagingHosts"
    TARGET_DIR_FIREFOX="$HOME/Library/Application Support/Mozilla/NativeMessagingHosts"
    TARGET_DIR_VIVALDI="$HOME/Library/Application Support/Vivaldi/NativeMessagingHosts"
  fi
  ;;
OpenBSD)
  HOST_FILE="$BIN_DIR/browserpass-openbsd64"
  if [ "$(whoami)" == "root" ]; then
    echo "Installing as root not supported."
    exit 1
  fi
  TARGET_DIR_CHROME="$HOME/.config/google-chrome/NativeMessagingHosts"
  TARGET_DIR_CHROMIUM="$HOME/.config/chromium/NativeMessagingHosts"
  TARGET_DIR_FIREFOX="$HOME/.mozilla/native-messaging-hosts"
  TARGET_DIR_VIVALDI="$HOME/.config/vivaldi/NativeMessagingHosts"
  ;;
FreeBSD)
  HOST_FILE="$BIN_DIR/browserpass-freebsd64"
  if [ "$(whoami)" == "root" ]; then
    echo "Installing as root not supported"
    exit 1
  fi
  TARGET_DIR_CHROME="$HOME/.config/google-chrome/NativeMessagingHosts"
  TARGET_DIR_CHROMIUM="$HOME/.config/chromium/NativeMessagingHosts"
  TARGET_DIR_FIREFOX="$HOME/.mozilla/native-messaging-hosts"
  TARGET_DIR_VIVALDI="$HOME/.config/vivaldi/NativeMessagingHosts"
  ;;
*)
  echo "$OPERATING_SYSTEM is not supported"
  exit 1
  ;;
esac

if [ -e "$BIN_DIR/browserpass" ]; then
  echo "Detected development binary"
  HOST_FILE="$BIN_DIR/browserpass"
fi

if [ -z "$BROWSER" ]; then
  echo ""
  echo "Select your browser:"
  echo "===================="
  echo "1) Chrome"
  echo "2) Chromium"
  echo "3) Firefox"
  echo "4) Vivaldi"
  echo -n "1-4: "
  read BROWSER
  echo ""
fi

# Set target dir from user input
case $BROWSER in
1|[Cc]hrome)
  BROWSER_NAME="Chrome"
  TARGET_DIR="$TARGET_DIR_CHROME"
  ;;
2|[Cc]hromium)
  BROWSER_NAME="Chromium"
  TARGET_DIR="$TARGET_DIR_CHROMIUM"
  ;;
3|[Ff]irefox)
  BROWSER_NAME="Firefox"
  TARGET_DIR="$TARGET_DIR_FIREFOX"
  ;;
4|[Vv]ivaldi)
  BROWSER_NAME="Vivaldi"
  TARGET_DIR="$TARGET_DIR_VIVALDI"
  ;;
*)
  echo "Invalid selection. Please select 1-4 or one of the browser names."
  exit 1
  ;;
esac

echo "Installing $BROWSER_NAME host config"

# Create config dir if not existing
mkdir -p "$TARGET_DIR"

if [ "$BROWSER_NAME" == "Firefox" ]; then
  MANIFEST="$JSON_DIR/firefox-host.json"
else
  MANIFEST="$JSON_DIR/chrome-host.json"
  POLICY="$JSON_DIR/chrome-policy.json"
fi

# Copy native host manifest, filling in binary path
assert_file_exists "$MANIFEST"
sed "s/%%replace%%/${HOST_FILE////\\/}/" "$MANIFEST" \
  > "$TARGET_DIR/$APP_NAME.json"

# Set permissions for the manifest so that all users can read it.
chmod o+r "$TARGET_DIR/$APP_NAME.json"

# Copy policy file, if any
if [ -n "$POLICY" ]; then
  assert_file_exists "$POLICY"
  POLICY_DIR="$TARGET_DIR"/../policies/managed/
  mkdir -p "$POLICY_DIR"
  cp "$POLICY" "$POLICY_DIR/$APP_NAME.json"
fi

echo "Native messaging host for $BROWSER_NAME has been installed to $TARGET_DIR."
