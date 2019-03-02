# -*- mode: ruby -*-
# vi: set ft=ruby :

ENV['VAGRANT_DEFAULT_PROVIDER'] = 'docker'

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.ssh.username = "develop"
  config.ssh.password = "bita123"

  config.vm.provision "shell" do |s|
    s.path   = "scripts/provision.sh"
    s.args   = [%x(ip addr | grep inet | grep docker0 | awk -F" " '{print $2}'| sed -e 's/\\/.*$//')]
  end

  config.vm.network "forwarded_port", guest: 8090,    	host: 8070     # nginx
  config.vm.network "forwarded_port", guest: 8080,    	host: 8080     # goconvey
  config.vm.network "forwarded_port", guest: 2345,    	host: 2345     # dlv
  config.vm.network "forwarded_port", guest: 15672,     host: 15672    # rabbitmq management
  config.vm.network "forwarded_port", guest: 22,        host: 5555     # ssh server
  config.vm.network "forwarded_port", guest: 9898,      host: 9898     # ssh server
  config.vm.network "forwarded_port", guest: 8090,      host: 8090     # console
  config.vm.network "forwarded_port", guest: 3500,      host: 3500     # supplier
  config.vm.synced_folder ".", "/home/develop/go/src/clickyab.com/crane", owner: "develop", group: "develop", create: true

  config.vm.provider "docker" do |d|
    d.image = "registry.3rdad.com/clickyab/baseimage-go"
    d.has_ssh = true
    d.privileged = true
    d.cmd = ["/bin/bash", "/home/develop/go/src/clickyab.com/crane/scripts/init.sh"]
  end
end
