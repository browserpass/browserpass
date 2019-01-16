# Browserpass

Browserpass is a Chrome & Firefox extension for [zx2c4's pass](https://www.passwordstore.org/), a UNIX based password manager. It retrieves your decrypted passwords for the current domain and allows you to auto-fill login forms, as well as copy it to clipboard. If you have multiple logins for the current site, the extension shows you a list of usernames to choose from.

![Browserpass in the Chrome menu](https://user-images.githubusercontent.com/1177900/38048732-79007146-32c6-11e8-98d5-557ba3f8b262.gif)

It uses a [native binary written in Golang](https://github.com/dannyvankooten/browserpass/blob/master/browserpass.go) to do the interfacing with your password store. Secure communication between the binary and the browser extension is handled through [native messaging](https://developer.chrome.com/extensions/nativeMessaging).

## Table of Contents

- [Requirements](#requirements)
- [Installation](#installation)
- [Updates](#updates)
- [Usage](#usage)
- [Options](#options)
- [Security](#security)
- [FAQ](#faq)
- [Contributing](#contributing)
- [License](#license)

## Requirements

- A recent version of Chrome, Chromium or Firefox 50+.
- Pass (on UNIX)
- Your password filename must match your username **or** your file must have a line starting with `login:`, `user:` or `username:`, followed by your username.

_Examples_

```bash
$ pass website.com/johndoe
the-password

$ pass website.com
the-password
login: johndoe
```

## Installation

In order to install browserpass correctly, you have to install two of its components:

- host application
- browser extension(s).

### Installing the host application

The following OS have a browserpass package that can be installed via package manager:

- [Arch Linux](https://aur.archlinux.org/packages/browserpass/)
- [Debian GNU/Linux](https://tracker.debian.org/pkg/browserpass)
- [macOS](https://github.com/dustinwilson/homebrew-tap/blob/master/browserpass.rb) - make sure to read [these instructions](#installing-browserpass-on-macos-with-homebrew)
- [NixOS](https://github.com/NixOS/nixpkgs/blob/master/pkgs/tools/security/browserpass/default.nix) - make sure to read [these instructions](#configuring-browserpass-on-nixos--for-nix)
- [Ubuntu](https://launchpad.net/ubuntu/+source/browserpass)

If your OS is not listed above, proceed with the manual installation steps below.

#### Download the latest Github release.

Start out by downloading the [latest release package](https://github.com/dannyvankooten/browserpass/releases) for your operating system.

#### Verifying authenticity of the releases

All release files are signed with [this PGP key](https://keybase.io/maximbaz). To verify the signature of a given file, use `$ gpg --verify <file>.sig`.

It should report:

```
gpg: Signature made ...
gpg:                using RSA key 8053EB88879A68CB4873D32B011FDC52DA839335
gpg: Good signature from "Maxim Baz <...>"
gpg:                 aka ...
Primary key fingerprint: EB4F 9E5A 60D3 2232 BB52  150C 12C8 7A28 FEAC 6B20
     Subkey fingerprint: 8053 EB88 879A 68CB 4873  D32B 011F DC52 DA83 9335
```

#### Installing the host application

1. Extract the package to where you would like to have the binary.
1. Run `./install.sh` (`.\install.ps1` on Windows) to install the native messaging host. If you want a system-wide installation, run the script with `sudo`. For Windows, system-wide installation can be done by running `.\install.ps1` as Administrator and specifying "yes" at the "Install for all users?" prompt.
   - If you desire a non-interactive installation on a Unix system, pass the name of the browser to the script (e.g. `./install.sh chrome`)

Installing the binary & registering it with your browser through the installation script is required to allow the browser extension to talk to the local binary application.

#### Installing the host application on Windows through WSL

If you already use `pass` under WSL and prefer to have a single copy of your password store, you can use browserpass through WSL as well.

1. Install the Windows host application (see previous section) as well as the Linux host application (under WSL).
2. Create `%localappdata%\browserpass\browserpass-wsl.bat` with the following contents:

```
@echo off
bash -c ~/.browserpass/browserpass-linux64
```

If you installed the Linux host application in a location different from `~/.browserpass`, replace that path in the above script.

3. Change the path in `%localappdata%\browserpass\browserpass-firefox.json` (or `-chrome.json`) to point to `browserpass-wsl.bat`

If your GPG key has a password, the host application running under WSL won't be able to unlock it since it can't interactively prompt for the password. This means you can't decrypt any passwords unless you've already got the key loaded in gpg-agent.
As a workaround, you can use the key (`pass website.com`) in a WSL terminal to load the key into gpg-agent. Then browserpass will work until gpg-agent times out (it is possible to configure larger timeouts, check manual for gpg-agent).

### Installing the Chrome extension

You can either [install the Chrome extension from the Chrome Web Store](https://chrome.google.com/webstore/detail/browserpass/naepdomgkenhinolocfifgehidddafch) or drag the `chrome-browserpass.crx` file from the release package into the [Chrome Extensions](chrome://extensions) (`chrome://extensions`) page.

### Installing the Firefox extension

You can [install the Firefox extension from the Mozilla add-ons site](https://addons.mozilla.org/en-US/firefox/addon/browserpass-ce/). Please note that you will need Firefox 50 or higher.

## Updates

**IMPORTANT**: Majority of the improvements require changing code in both browser extensions and the host application. While we are trying to maintain backwards compatibility, it is expected that you will make sure to keep both components up to date.

### Updating the host application

If you installed the host application via a package manager for your OS, you will likely update it in the the same way.

If not, repeat the installation instructions for your OS.

### Updating browser extensions

If you installed the extension from a webstore, you will receive updates automatically.

If not, repeat the installation instructions for the extension.

## Usage

Click the lock icon or use <kbd>Ctrl+Shift+L</kbd> to open browserpass with the entries that match current domain.

- Chrome allows changing the shortcut via chrome://extensions > Keyboard shortcuts.
- Firefox unfortunately does not allow changing the default shortcut.
- Firefox supports the keyboard shortcut only since version 53.

### Filter and search modes

Browserpass has two modes for working with password entries: filter and search.

When opened, browserpass automatically switches to the filter mode if at least one matching entries exists.

**Filter mode** is designed to quickly refine a few search results, for example to choose one of several accounts that you have on a given domain. This is done on client side, the filter is always fuzzy and always works in real time. When browserpass is in the filter mode, you will see a domain name in the input field. To exit filter mode, press <kbd>Backspace</kbd>.

**Search mode** is designed to search password entries on your disk, this is much more expensive operation (especially visible on Windows) that's why it is **not** real time, and instead searches only when <kbd>Enter</kbd> is pressed. The search is fuzzy by default, but can be changed to glob algorithm in the options. If you want to search everything interactively, just search for `/` or `.` and then use the filter mode to refine the search in real time.

### Fill (and submit) the login form

Click or select the entry that you want to submit, and the login form will be filled with the selected credentials (injected directly into the DOM, browserpass does not use clipboard for this). When the focus is in the input field, hitting <kbd>Enter</kbd> will submit the first entry in the list (this is useful in combination with filter mode).

If the login button is found, it will be focused so that you can just hit <kbd>Enter</kbd> to submit the form. If you enable `Automatically submit forms after filling` in the options, the login button will be pressed instead.

If your password entry has OTP configuration, browserpass will use it at this point to display the code.

### Navigating the entries

Navigate through the list of available credentials with <kbd>Tab</kbd> and <kbd>Shift+Tab</kbd> or with arrow keys.

### Copy to clipboard

Click on the username or password buttons to copy them to clipboard. Keyboard shortcuts are also available, use <kbd>Ctrl+C</kbd> to copy password of the selected entry and <kbd>Shift+C</kbd> to copy the username.

### Open URL

Click on the globe button or use the <kbd>g</kbd> shortcut to navigate to the URL in the current tab, hold <kbd>Shift</kbd> while doing so to open a new tab instead. You can also specify one of the following metadata fields in your pass file to control exactly which URL is navigated to: `url:`, `link:`, `website:`, `web:` or `site:`.

Keep in mind that browserpass can only fill HTTP basic auth credentials _if you open this URL using browserpass_.

### Manual search

To prevent phishing attacks, browserpass prefills the list of passwords with only those entries that match the current domain. If you want search for credentials across the entire password store, exit the filter mode with <kbd>Backspace</kbd> (domain name in the input field will disappear), type the search request and hit <kbd>Enter</kbd> to start the search. Instead of using <kbd>Backspace</kbd>, you can also type your search query while in the filter mode, as soon as there are no matching results left browserpass will automatically switch to the search mode and will await <kbd>Enter</kbd> to initiate the search.

### Password store location(s)

When deciding where to look for the password store, browserpass uses `PASSWORD_STORE_DIR` environment variable, and if it is not defined, checks the `~/.password-store` folder. However, using the `Custom store locations` setting in the options of the browser extension you can configure a different location for browserpass to look for, or even multiple locations. There are no restrictions, you can define subfolders in the password store, gopass mounts or any other folder that has pass entries.

When you have more than one password store configured and enabled, in order to help you distinguish the password entries from different locations (e.g. between passwords for work and personal GitHub accounts), a green badge next to each password entry will appear indicating its origin (the name of its password store).

## Options

Open settings to configure browserpass:

- Right click on the lock icon > "Options".
- Find the browserpass in the list of extensions in your browser > "Options".

The list of currently available options:

- `Automatically submit forms after filling`: make browserpass automatically submit the login form for you.
- `Use fuzzy search`: whether the _manual search mode_ should be fuzzy or not (filter mode is always fuzzy).
- `Custom store locations`: allows configuring multiple password store locations and toggle them on the fly.

## Security

Browserpass aims to protect your passwords and computer from malicious or fraudulent websites.

- To protect against phishing, only passwords matching the origin hostname are suggested or selected without an explicit search term.
- To minimize attack surface, the website is not allowed to trigger any extension action without user invocation.
- Only data from the selected password is made available to the website.
- Given full control of the non-native component of the extension, the attacker can extract passwords stored in the configured repository, but can not obtain files elsewhere on the filesystem or reach code execution.

## FAQ

### Does not work on MacOS: "Native host has exited"

First install required dependencies:

```
$ brew install gnupg pinentry-mac
```

It is important that you have the `gpg` binary at `/usr/local/bin/gpg`. If you have your `gpg` in another location, create a symlink:

```
$ sudo ln -s /path/to/your/gpg /usr/local/bin/gpg
```

If you don't have admin rights to create the symlink, the workaround is to [patch browser launcher](https://github.com/dannyvankooten/browserpass/issues/197#issuecomment-354586602).

Now edit `~/.gnupg/gpg.conf`:

```
# Comment out or remove this line if it's there:
# pinentry-mode loopback

# and add this line:
use-agent
```

Add the following line to `~/.gnupg/gpg-agent.conf`:

```
pinentry-program /usr/local/bin/pinentry-mac
```

Then restart `gpg-agent`:

```
$ gpgconf --kill gpg-agent
```

And finally restart your browser.

If you still experience the issue, try starting your browser from terminal. If this helps, the issue is likely due to the absence of `/usr/local/bin/gpg`, follow the steps above to make sure it exists.

### Configuring Browserpass on NixOS / for Nix

#### On NixOS

If you wish to have a stateless setup, make sure you have this in your `/etc/nixos/configuration.nix` and rebuild your system:

```nix
{ pkgs, ... }: {
  programs.browserpass.enable = true;
  environment.systemPackages = with pkgs; [
    # All of these browsers will work with it
    chromium
    firefox
    google-chrome
    vivaldi
  ];
}
```

Note: firefox*-bin versions do _not_ work statelessly. If you require such firefox versions, use the stateful setup in the following section.

#### For Nix / stateful

Install browserpass native messaging host with

```
nix-env -iA nixpkgs.browserpass
```

And install the browser extension like normal. Then link the necessary files

```
# For firefox
mkdir -p ~/.mozilla/native-messaging-hosts && \
  ln -s ~/.nix-profile/lib/mozilla/native-messaging-hosts/com.dannyvankooten.browserpass.json ~/.mozilla/native-messaging-hosts
# For chrome
mkdir -p ~/.config/google-chrome/NativeMessagingHosts && \
  ln -s ~/.nix-profile/etc/chrome-host.json ~/.config/google-chrome/NativeMessagingHosts/com.dannyvankooten.browserpass.json
# For chromium
mkdir -p ~/.config/chromium/NativeMessagingHosts && \
  ln -s ~/.nix-profile/etc/chrome-host.json ~/.config/chromium/NativeMessagingHosts/com.dannyvankooten.browserpass.json
# For vivaldi
mkdir -p ~/.config/vivaldi/NativeMessagingHosts && \
  ln -s ~/.nix-profile/etc/chrome-host.json ~/.config/vivaldi/NativeMessagingHosts/com.dannyvankooten.browserpass.json
```

All versions of firefox are supported with this way

### Installing Browserpass on macOS with Homebrew

Browserpass isn't included in the main Homebrew repository, so it must be installed by adding a third party "tap". That only requires one additional step.

```
$ brew tap dustinwilson/tap
$ brew install browserpass
```

Instead of running `install.sh` Homebrew supplies an additional command called `browserpass-setup` to handle this and works the same way as `install.sh` above. For example this will install the native host files for Firefox:

```
$ browserpass-setup firefox
```

You must install the browser extensions manually using conventional methods for each browser. All of this information is supplied when running `brew install browserpass`.

### How to configure OTP?

The easiest way to add OTP in your password entries is to use [pass-otp](https://github.com/tadfisher/pass-otp). You don't have to configure anything extra, browserpass will automatically detect if an OTP is configured and show you the code after filling the form.

## Contributing

Check out [Contributing](CONTRIBUTING.md) for details on how to build browser extension and host app from sources, and how to load browserpass as an unpacked extension into your browser.

## License

MIT Licensed.
