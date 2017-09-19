# CONTRIBUTING

You will need Node, [Yarn](https://yarnpkg.com/), Golang and [dep](https://github.com/golang/dep) installed.

## To build
- Open `makefile` and, if needed, change `google-chrome` to the appropriate name of your Google Chrome executable (in Linux, it could be google-chrome-stable)
- Run `make deps` to install Node/GO dependencies
- Run `make release`

The command above will generate packed extensions for both Firefox and Chrome and compile the Go binaries for Linux and MacOSX.

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
