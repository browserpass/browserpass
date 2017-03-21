# CONTRIBUTING

You will need Node, [Yarn](https://yarnpkg.com/) and Golang installed.

## To build
- Open `makefile` and, if needed, change `google-chrome` to the appropriate name of your Google Chrome executable (in Linux, it could be google-chrome-stable)
- Run `make js` to install Node dependencies and generate a JS bundle
- Run `make release`

The command above will generate packed extensions for both Firefox and Chrome and compile the Go binaries for Linux and MacOSX.

## Setting up a Vagrant build environment and building the Linux binary

These instructions will walk you through the process of setting up a build environment for browserpass using [Vagrant](https://www.vagrantup.com/) on Debian/Ubuntu. These instructions were valid for an Ubuntu 16.04 host. This only addresses building the Linux 64-bit binary - you'll need to faff around a bit to do other things, but this should provide you with a good starting point.

Install vagrant:
```shell
$ sudo apt-get install vagrant
```

Make a temporary build directory, set up a copy of Ubuntu 16.04 in it, and start a VM.
```shell
$ mkdir browserpass-build
$ cd browserpass-build
~/browserpass-build $ vagrant init minimal/xenial64
[ A message from Vagrant about how it wrote a Vagrantfile in this directory ]
~/browserpass-build $ vagrant up
[ Messages from Vagrant as it downloads, installs, and boots this image in a new VM ]
```
(You can find alternate Vagrant images in the [Atlas Vagrant image repository](https://atlas.hashicorp.com/boxes/search?order=desc&page=1&provider=&q=xenial64&sort=downloads&utf8=%E2%9C%93))

ssh to the new build environment.
```shell
~/browserpass-build $ vagrant ssh
vagrant@minimal-xenial:~$
```

Inside the VM, install browserpass build dependencies.
```shell
vagrant@minimal-xenial:~$ sudo apt-get update
vagrant@minimal-xenial:~$ sudo apt-get install nodejs golang cmdtest
```

Set up the go and browserpass build directories.
```shell
vagrant@minimal-xenial:~$ export GOPATH=$HOME/go
vagrant@minimal-xenial:~$ mkdir -p go/src/github.com/dannyvankooten/
vagrant@minimal-xenial:~$ cd go/src/github.com/dannyvankooten/
vagrant@minimal-xenial:~/go/src/github.com/dannyvankooten$ git clone https://github.com/dannyvankooten/browserpass.git
mal-xenial:~/go/src/github.com/dannyvankooten$ cd browserpass/
vagrant@minimal-xenial:~/go/src/github.com/dannyvankooten/browserpass$
```
(The git clone happens over https because you won't have your ssh key in the dev VM unless you set up agent forwarding, etc)

Build the Linux browserpass binary!
```shell
vagrant@minimal-xenial:~/go/src/github.com/dannyvankooten/browserpass$ make browserpass-linux64
```

Copy the resulting go binary (or whatever else you need) out of the vagrant VM:
```shell
vagrant@minimal-xenial:~/go/src/github.com/dannyvankooten/browserpass$ cp browserpass-linux64 /vagrant/
```

Exit the build environment, clean up the vagrant image, and pick up the previously-copied output files.
```shell
vagrant@minimal-xenial:~/go/src/github.com/dannyvankooten/browserpass$ exit
~/browserpass-build $ vagrant destroy
[ Vagrant tells you about stopping and removing the VM ]
~/browserpass-build $ ls
browserpass-linux64  Vagrantfile
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
