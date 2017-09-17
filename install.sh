#!/usr/bin/env bash

set -e

DIR="$( cd "$( dirname "$0" )" && pwd )"
APP_NAME="com.dannyvankooten.browserpass"
HOST_FILE="$DIR/browserpass"

# Find target dirs for various browsers & OS'es
# https://developer.chrome.com/extensions/nativeMessaging#native-messaging-host-location
# https://wiki.mozilla.org/WebExtensions/Native_Messaging
OPERATING_SYSTEM=$(uname -s)

case $OPERATING_SYSTEM in
Linux)
  HOST_FILE="$DIR/browserpass-linux64"
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
  HOST_FILE="$DIR/browserpass-darwinx64"
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
  HOST_FILE="$DIR/browserpass-openbsd64"
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
  HOST_FILE="$DIR/browserpass-freebsd64"
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

if [ -e "$DIR/browserpass" ]; then
  echo "Detected development binary"
  HOST_FILE="$DIR/browserpass"
fi

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

# Set target dir from user input
case $BROWSER in
1)
  BROWSER_NAME="Chrome"
  TARGET_DIR="$TARGET_DIR_CHROME"
  ;;
2)
  BROWSER_NAME="Chromium"
  TARGET_DIR="$TARGET_DIR_CHROMIUM"
  ;;
3)
  BROWSER_NAME="Firefox"
  TARGET_DIR="$TARGET_DIR_FIREFOX"
  ;;
4)
  BROWSER_NAME="Vivaldi"
  TARGET_DIR="$TARGET_DIR_VIVALDI"
  ;;
*)
  echo "Invalid selection.  Please select 1-4."
  exit 1
  ;;
esac

echo "Installing $BROWSER_NAME host config"

# Create config dir if not existing
mkdir -p "$TARGET_DIR"

# Escape host file
ESCAPED_HOST_FILE=${HOST_FILE////\\/}

# Copy manifest host config file
if [ "$BROWSER_NAME" == "Chrome" ] || \
   [ "$BROWSER_NAME" == "Chromium" ] || \
   [ "$BROWSER_NAME" == "Vivaldi" ]; then
  if [ ! -f "$DIR/chrome-host.json" ]; then
    echo "ERROR: '$DIR/chrome-host.json' is missing."
    echo "If you are running './install.sh' from a release archive, please file a bug."
    echo "If you are running './install.sh' from the source code, make sure to follow CONTRIBUTING.md on how to build first."
    exit 1
  fi
  cp "$DIR/chrome-host.json" "$TARGET_DIR/$APP_NAME.json"
else
  if [ ! -f "$DIR/firefox-host.json" ]; then
    echo "ERROR: '$DIR/firefox-host.json' is missing."
    echo "If you are running './install.sh' from a release archive, please file a bug."
    echo "If you are running './install.sh' from the source code, make sure to follow CONTRIBUTING.md on how to build first."
    exit 1
  fi
  cp "$DIR/firefox-host.json" "$TARGET_DIR/$APP_NAME.json"
fi

# Replace path to host
sed -i -e "s/%%replace%%/$ESCAPED_HOST_FILE/" "$TARGET_DIR/$APP_NAME.json"

# Set permissions for the manifest so that all users can read it.
chmod o+r "$TARGET_DIR/$APP_NAME.json"

echo "Native messaging host for $BROWSER_NAME has been installed to $TARGET_DIR."
