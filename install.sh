#!/usr/bin/env bash

set -e

DIR="$( cd "$( dirname "$0" )" && pwd )"
APP_NAME="com.dannyvankooten.browserpass"
HOST_FILE="$DIR/browserpass"

install_host_config () {
    browser_name=$1
    target_dir=$2

    echo "Installing $browser_name host config"

    # Create config dir if not existing
    mkdir -p "$target_dir"

    # Escape host file
    ESCAPED_HOST_FILE=${HOST_FILE////\\/}


    # Copy manifest host config file
    if [ "$browser_name" == "Chrome" ] || \
        [ "$browser_name" == "Chromium" ] || \
        [ "$browser_name" == "Vivaldi" ]; then
        if [ ! -f "$DIR/chrome-host.json" ]; then
            echo "ERROR: '$DIR/chrome-host.json' is missing."
            echo "If you are running './install.sh' from a release archive, please file a bug."
            echo "If you are running './install.sh' from the source code, make sure to follow CONTRIBUTING.md on how to build first."
            exit 1
        fi
        cp "$DIR/chrome-host.json" "$target_dir/$APP_NAME.json"
    else
        if [ ! -f "$DIR/firefox-host.json" ]; then
            echo "ERROR: '$DIR/firefox-host.json' is missing."
            echo "If you are running './install.sh' from a release archive, please file a bug."
            echo "If you are running './install.sh' from the source code, make sure to follow CONTRIBUTING.md on how to build first."
            exit 1
        fi
        cp "$DIR/firefox-host.json" "$target_dir/$APP_NAME.json"
    fi

    # Replace path to host
    sed -i -e "s/%%replace%%/$ESCAPED_HOST_FILE/" "$target_dir/$APP_NAME.json"

    # Set permissions for the manifest so that all users can read it.
    chmod o+r "$target_dir/$APP_NAME.json"

    echo "Native messaging host for $browser_name has been installed to $target_dir."
}

# Find target dirs for various browsers & OS'es
# https://developer.chrome.com/extensions/nativeMessaging#native-messaging-host-location
# https://wiki.mozilla.org/WebExtensions/Native_Messaging
OPERATING_SYSTEM=$(uname -s)

case $OPERATING_SYSTEM in
    Linux)
        HOST_FILE="$DIR/browserpass-linux64"
        if [ "$(whoami)" == "root" ]; then
            BROWSER_PATHS=( "Chrome:/etc/opt/chrome/native-messaging-hosts"
            "Chromium:/etc/chromium/native-messaging-hosts"
            "Firefox:/usr/lib/mozilla/native-messaging-hosts"
            "Vivaldi:/etc/vivaldi/native-messaging-hosts" )
        else
            BROWSER_PATHS=( "Chrome:$HOME/.config/google-chrome/NativeMessagingHosts"
            "Chromium:$HOME/.config/chromium/NativeMessagingHosts"
            "Firefox:$HOME/.mozilla/native-messaging-hosts"
            "Vivaldi:$HOME/.config/vivaldi/NativeMessagingHosts" )
        fi
        ;;
    Darwin)
        HOST_FILE="$DIR/browserpass-darwinx64"
        if [ "$(whoami)" == "root" ]; then
            BROWSER_PATHS=( "Chrome:/Library/Google/Chrome/NativeMessagingHosts"
            "Chromium:/Library/Application Support/Chromium/NativeMessagingHosts"
            "Firefox:/Library/Application Support/Mozilla/NativeMessagingHosts"
            "Vivaldi:/Library/Application Support/Vivaldi/NativeMessagingHosts" )
        else
            BROWSER_PATHS=( "Chrome:$HOME/Library/Application Support/Google/Chrome/NativeMessagingHosts"
            "Chromium:$HOME/Library/Application Support/Chromium/NativeMessagingHosts"
            "Firefox:$HOME/Library/Application Support/Mozilla/NativeMessagingHosts"
            "Vivaldi:$HOME/Library/Application Support/Vivaldi/NativeMessagingHosts" )
        fi
        ;;
    OpenBSD)
        HOST_FILE="$DIR/browserpass-openbsd64"
        if [ "$(whoami)" == "root" ]; then
            echo "Installing as root not supported."
            exit 1
        fi
        BROWSER_PATHS=( "Chrome:$HOME/.config/google-chrome/NativeMessagingHosts"
        "Chromium:$HOME/.config/chromium/NativeMessagingHosts"
        "Firefox:$HOME/.mozilla/native-messaging-hosts"
        "Vivaldi:$HOME/.config/vivaldi/NativeMessagingHosts" )
        ;;
    FreeBSD)
        HOST_FILE="$DIR/browserpass-freebsd64"
        if [ "$(whoami)" == "root" ]; then
            echo "Installing as root not supported"
            exit 1
        fi
        BROWSER_PATHS=( "Chrome:$HOME/.config/google-chrome/NativeMessagingHosts"
        "Chromium:$HOME/.config/chromium/NativeMessagingHosts"
        "Firefox:$HOME/.mozilla/native-messaging-hosts"
        "Vivaldi:$HOME/.config/vivaldi/NativeMessagingHosts" )
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

no_browser_found=true
rejected_installs=()

for browser in "${BROWSER_PATHS[@]}" ; do
    browser_name=${browser%%:*}
    browser_path=${browser#*:}

    if [ -d $(dirname "$browser_path") ]; then
        no_browser_found=false
        read -r -p "Detected $browser_name. Install for this browser? [Y/n] " response
        case "$response" in
            [yY][eE][sS]|[yY])
                install_host_config $browser_name $browser_path
                ;;
            *)
                rejected_installs+=("$(dirname "$browser_path")")
                ;;
        esac
    fi
done

if [ "$no_browser_found" = true ]; then
    echo "No compatible browsers found."
    echo "If you do actually have a browser installed which you think should be supported,"
    echo "please submit a bug report to https://github.com/dannyvankooten/browserpass, detailing"
    echo "the exact name of the browser, and (if you know it) the location of its configuration"
    echo "directory."
    exit 1
fi

if [ ${#rejected_installs[@]} -ne 0 ]; then
    echo "Note: If this script wrongly detected a browser you no longer have (but once had)"
    echo "on your system, it might be because the configuration directory was left behind"
    echo "after uninstalling."
    echo
    echo "To fix this, try issuing the following commands (only if you don't need the old"
    echo "config anymore, obviously):"
    echo
    for path in "${rejected_installs[@]}" ; do
        echo "$ rm -rf $path"
    done
fi

exit 0
