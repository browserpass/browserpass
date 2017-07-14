# CONTRIBUTING

You will need Node, [Yarn](https://yarnpkg.com/) and Golang installed.

## To build
- `npm install -g browserify`
- Open `makefile` and, if needed, change `google-chrome` to the appropriate name of your Google Chrome executable (in Linux, it could be google-chrome-stable)
- Run `make js` to install Node dependencies and generate a JS bundle
- Run `make release`

The command above will generate packed extensions for both Firefox and Chrome and compile the Go binaries for Linux and MacOSX.

## To contribute

1. Fork [the repo](https://github.com/dannyvankooten/browserpass)
2. Create your feature branch
   * `git checkout -b my-new-feature`
3. Commit your changes
   * `git commit -am 'Add some feature'`
4. Push to the branch
   * `git push origin my-new-feature`
5. Create new pull Request
