# -*- mode: ruby -*-
# vi: set ft=ruby :

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  # All Vagrant configuration is done here. The most common configuration
  # options are documented and commented below. For a complete reference,
  # please see the online documentation at vagrantup.com.

  config.vm.define "ubuntu" do |ubuntu|
    ubuntu.vm.box = "phusion/ubuntu-14.04-amd64"
    ubuntu.vm.provision :shell, path: "tests/data/provisioning/vagrant-ubuntu.sh"
  end

  config.vm.define "centos" do |centos|
    centos.vm.box = "chef/centos-7.0"
    centos.vm.provision :shell, path: "tests/data/provisioning/vagrant-centos.sh"
  end

  config.ssh.insert_key = 'true'

  config.ssh.forward_agent = true

  config.vm.synced_folder ENV['GOPATH'], "/go"

  # Provider-specific configuration so you can fine-tune various
  # backing providers for Vagrant. These expose provider-specific options.
  # Example for VirtualBox:
  #
  # config.vm.provider "virtualbox" do |vb|
  #   # Don't boot with headless mode
  #   vb.gui = true
  #
  #   # Use VBoxManage to customize the VM. For example to change memory:
  #   vb.customize ["modifyvm", :id, "--memory", "1024"]
  # end
  #
  # View the documentation for the provider you're using for more
  # information on available options.
end
