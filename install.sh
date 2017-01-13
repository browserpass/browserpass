#!/usr/bin/env bash

set -e

DIR="$( cd "$( dirname "$0" )" && pwd )"
APP_NAME="com.dannyvankooten.browserpass"
HOST_FILE="$DIR/browserpass"

# Find target dirs for various browsers & OS'es
# https://developer.chrome.com/extensions/nativeMessaging#native-messaging-host-location
# https://wiki.mozilla.org/WebExtensions/Native_Messaging
if [ $(uname -s) == 'Darwin' ]; then
  HOST_FILE="$DIR/browserpass-darwinx64"
  if [ "$(whoami)" == "root" ]; then
    TARGET_DIR_CHROME="/Library/Google/Chrome/NativeMessagingHosts"
    TARGET_DIR_CHROMIUM="/Library/Application Support/Chromium/NativeMessagingHosts"
    TARGET_DIR_FIREFOX="/Library/Application Support/Mozilla/NativeMessagingHosts"
  else
    TARGET_DIR_CHROME="$HOME/Library/Application Support/Google/Chrome/NativeMessagingHosts"
    TARGET_DIR_CHROMIUM="$HOME/Library/Application Support/Chromium/NativeMessagingHosts"
    TARGET_DIR_FIREFOX="$HOME/Library/Application Support/Mozilla/NativeMessagingHosts"
  fi
else
  HOST_FILE="$DIR/browserpass-linux64"
  if [ "$(whoami)" == "root" ]; then
    TARGET_DIR_CHROME="/etc/opt/chrome/native-messaging-hosts"
    TARGET_DIR_CHROMIUM="/etc/chromium/native-messaging-hosts"
    TARGET_DIR_FIREFOX="/usr/lib/mozilla/native-messaging-hosts"
  else
    TARGET_DIR_CHROME="$HOME/.config/google-chrome/NativeMessagingHosts"
    TARGET_DIR_CHROMIUM="$HOME/.config/chromium/NativeMessagingHosts/"
    TARGET_DIR_FIREFOX="$HOME/.mozilla/native-messaging-hosts"
  fi
fi

if [ -e "$DIR/browserpass" ]; then
  echo "Detected development binary"
  HOST_FILE="$DIR/browserpass"
fi

# Escape host file
ESCAPED_HOST_FILE=${HOST_FILE////\\/}

echo ""
echo "Select your browser:"
echo "===================="
echo "1) Chrome"
echo "2) Chromium"
echo "3) Firefox"
echo -n "1-3: "
read BROWSER
echo ""

# Set target dir from user input
if [[ "$BROWSER" == "1" ]]; then
  TARGET_DIR="$TARGET_DIR_CHROME"
fi

if [[ "$BROWSER" == "2" ]]; then
  TARGET_DIR="$TARGET_DIR_CHROMIUM"
fi

if [[ "$BROWSER" == "3" ]]; then
  TARGET_DIR="$TARGET_DIR_FIREFOX"
fi

# Create config dir if not existing
mkdir -p "$TARGET_DIR"

if [ "$BROWSER" == "1" ] || [ "$BROWSER" == "2" ]; then
  echo "Installing Chrome / Chromium host config"
  cp "$DIR/chrome-host.json" "$TARGET_DIR/$APP_NAME.json"
else
  echo "Installing Firefox host config"
  cp "$DIR/firefox-host.json" "$TARGET_DIR_FIREFOX/$APP_NAME.json"
fi

# Replace path to host
sed -i -e "s/%%replace%%/$ESCAPED_HOST_FILE/" "$TARGET_DIR/$APP_NAME.json"

# Set permissions for the manifest so that all users can read it.
chmod o+r "$TARGET_DIR/$APP_NAME.json"

echo "Native messaging host has been installed to $TARGET_DIR."
