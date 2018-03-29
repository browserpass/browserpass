Browserpass
=======

Browserpass is a Chrome & Firefox extension for [zx2c4's pass](https://www.passwordstore.org/), a UNIX based password manager. It retrieves your decrypted passwords for the current domain and allows you to auto-fill login forms, as well as copy it to clipboard. If you have multiple logins for the current site, the extension shows you a list of usernames to choose from.

![Browserpass in the Chrome menu](https://user-images.githubusercontent.com/1177900/38048732-79007146-32c6-11e8-98d5-557ba3f8b262.gif)

It uses a [native binary written in Golang](https://github.com/dannyvankooten/browserpass/blob/master/browserpass.go) to do the interfacing with your password store. Secure communication between the binary and the browser extension is handled through [native messaging](https://developer.chrome.com/extensions/nativeMessaging).

## Table of Contents

- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
- [Options](#options)
- [Security](#security)
- [FAQ](#faq)
- [Contribute](#contributing)
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

Start out by downloading the [latest release package](https://github.com/dannyvankooten/browserpass/releases) for your operating system. Prebuilt binaries for 64-bit OSX & Linux and Windows are available. Arch users can install browserpass [from the AUR](https://aur.archlinux.org/packages/browserpass/).

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
	* If you desire a non-interactive installation on a Unix system, pass the name of the browser to the script (i.e. `./install.sh chrome`)

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

#### Installing the Chrome extension

You can either [install the Chrome extension from the Chrome Web Store](https://chrome.google.com/webstore/detail/browserpass/naepdomgkenhinolocfifgehidddafch) or drag the `chrome-browserpass.crx` file from the release package into the [Chrome Extensions](chrome://extensions) (`chrome://extensions`) page.

#### Installing the Firefox extension

You can [install the Firefox extension from the Mozilla add-ons site](https://addons.mozilla.org/en-US/firefox/addon/browserpass-ce/). Please note that you will need Firefox 50 or higher.

## Usage

Click the lock icon or use <kbd>Ctrl+Shift+L</kbd> to fill & submit your login info for the current site.

- Chrome allows changing the shortcut via chrome://extensions > Keyboard shortcuts.
- Firefox unfortunately does not allow changing the default shortcut.
- Firefox supports the keyboard shortcut only since version 53.

Navigate through the list of available credentials with <kbd>Tab</kbd> / <kbd>Shift+Tab</kbd> or with arrow keys.

Click on the username or password buttons to copy them to clipboard. Keyboard shortcuts are also available, use <kbd>Ctrl+C</kbd> to copy password of the selected entry and <kbd>Shift+C</kbd> to copy the username.

Click on the globe button to navigate to the URL or use the <kbd>g</kbd> shortcut. You can also specify one of the following metadata fields in your pass file to control exactly which URL is navigated to: `url:`, `link:`, `website:`, `web:` or `site:`.

Keep in mind that browserpass can only fill HTTP basic auth credentials _if you open this URL using browserpass_.

## Options

Open settings to configure browserpass:

- Right click on the lock icon > "Options".
- Find the browserpass in the list of extensions in your browser > "Options".

The list of currently available options:

- `Auto-submit the form`: make browserpass automatically submit the login form for you.
- `Use fuzzy search`: whether the manual search should be fuzzy or not (refining search results is always fuzzy).

## Security

Browserpass aims to protect your passwords and computer from malicious or fraudulent websites.

* To protect against phishing, only passwords matching the origin hostname are suggested or selected without an explicit search term.
* To minimize attack surface, the website is not allowed to trigger any extension action without user invocation.
* Only data from the selected password is made available to the website.
* Given full control of the non-native component of the extension, the attacker can extract passwords stored in the configured repository, but can not obtain files elsewhere on the filesystem or reach code execution.

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

## Contributing

Check out [Contributing](CONTRIBUTING.md) for details on how to build browser extension and host app from sources, and how to load browserpass as an unpacked extension into your browser.

## License

MIT Licensed.
