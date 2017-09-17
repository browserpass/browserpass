Vagrant.configure("2") do |config|
  config.vm.box = "minimal/xenial64"
  config.vm.provision "shell" do |s|
    s.inline = %q{
      apt-get update
      apt-get install -y nodejs npm golang cmdtest
      ln -sfn $(which nodejs) /usr/local/bin/node

      grep -q GOPATH /home/vagrant/.profile || echo 'export GOPATH="${HOME}/go"' >> /home/vagrant/.profile
      grep -q browserify /home/vagrant/.profile || echo 'PATH="$PATH:/home/vagrant/go/src/github.com/dannyvankooten/browserpass/node_modules/.bin"' >> /home/vagrant/.profile
      mkdir -p go/src/github.com/dannyvankooten/
      ln -sfn /vagrant go/src/github.com/dannyvankooten/browserpass
    }
  end
end
