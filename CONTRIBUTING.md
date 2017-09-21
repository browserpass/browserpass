# CONTRIBUTING

You will need Docker or Node, [Yarn](https://yarnpkg.com/), Golang and [dep](https://github.com/golang/dep) installed.

## To build
- Run `make` to fetch all dependencies and compile both front-end and back-end code

OR

- Run `make deps` to download all dependencies (you don't need to run this very often)
- Run `make js` to compile only front-end code
- Run `make browserpass` to compile only back-end code

The commands above will generate unpacked extensions for both Firefox and Chrome and compile the Go binaries for all supported platforms.

## To load an unpacked extension
- Run `./install.sh` or `./install.ps1` to install the compiled Go binary
- For Chrome:
    - Go to `chrome://extensions`
    - Enable `Developer mode`
    - Click `Load unpacked extension`
    - Select `browserpass/chrome` directory
- For Firefox:
    - Go to `about:debugging#addons`
    - Click `Load temporary add-on`
    - Select `browserpass/firefox` directory

## Build using Docker

The [Dockerfile](Dockerfile) will set up a docker image suitable for running basic make targets such as building frontend and backend.

The `crx` target is not supported by now (therefore `release` target will not work).

To build the docker image run the following command in project root:
```shell
docker build -t browserpass-dev .
```

To build browserpass (frontend and backend) via docker, run the following from project root:
```shell
docker run --rm -v "$(pwd)":/browserpass browserpass-dev
```

If you only want to download dependencies, build front-end or backend code, run one of the following:
```shell
docker run --rm -v "$(pwd)":/browserpass browserpass-dev deps
docker run --rm -v "$(pwd)":/browserpass browserpass-dev js
docker run --rm -v "$(pwd)":/browserpass browserpass-dev browserpass
```

## Setting up a Vagrant build environment and building the Linux binary

Vagrant will set up a virtual machine with all dependencies installed for you. Your local working directory is shared into the VM.
These instructions will walk you through the process of setting up a build environment for browserpass using [Vagrant](https://www.vagrantup.com/) on Debian/Ubuntu. These instructions were valid for an Ubuntu 16.04 host. This only addresses building the Linux 64-bit binary - you'll need to faff around a bit to do other things, but this should provide you with a good starting point.

Install vagrant:
```shell
$ sudo apt-get install vagrant
```

Start the VM:
```shell
$ vagrant up
```

Jump into the VM and build the project:
```shell
$ vagrant ssh
vagrant@minimal-xenial:~$ cd go/src/github.com/dannyvankooten/browserpass
vagrant@minimal-xenial:~/go/src/github.com/dannyvankooten/browserpass$ make js
vagrant@minimal-xenial:~/go/src/github.com/dannyvankooten/browserpass$ make browserpass-linux64
```

Exit the build environment, clean up the vagrant image.
```shell
vagrant@minimal-xenial:~/go/src/github.com/dannyvankooten/browserpass$ exit
$ vagrant destroy
[ Vagrant tells you about stopping and removing the VM ]
```

## To contribute

1. Fork [the repo](https://github.com/dannyvankooten/browserpass)
2. Create your feature branch
   * `git checkout -b my-new-feature`
3. Commit your changes
   * `git commit -am 'Add some feature'`
4. Push to the branch
   * `git push origin my-new-feature`
5. Create new pull Request
