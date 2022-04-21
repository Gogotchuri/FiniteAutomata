#!/bin/bash

### Variables
BIN_DIR=./bin

### Functions
Help() {
  echo "This script will build the two binary files: *run* and *build*. Additional options can be specified:"
  echo "Usage: build.sh [-h|c|i]"
  echo "  -h: Print this help"
  echo "  -c: Clean the build directory"
  echo "  -i: Install the dependencies"
  exit 0
}

Clean() {
  echo "Cleaning the build directory..."
  rm -rf $BIN_DIR
  echo "Done Cleaning!"
}

InstallGolang() {
  CurrentLinuxID=$(cat "/etc/os-release" | grep '^ID=' | cut -d '=' -f 2)
  echo "Current Linux Distribution: $CurrentLinuxID"
  if [ "$CurrentLinuxID" == "ubuntu" ]; then
    echo "Installing Golang... Requires super user privileges (sudo)"
    # Install Latest version of Golang https://github.com/golang/go/wiki/Ubuntu
    sudo add-apt-repository -y ppa:longsleep/golang-backports
    sudo apt update -y
    sudo apt install -y golang-go
  elif [ "$CurrentLinuxID" == "centos" ]; then
    echo "Installing Golang... Requires super user privileges (sudo)"
    sudo yum install -y golang
  elif [ "$CurrentLinuxID" == "manjaro" ]; then
    echo "Installing Golang... Requires super user privileges (sudo)"
    sudo pacman -S --noconfirm --needed go
  else
    echo "Unsupported Linux Distribution. Install go manually!"
    exit 1
  fi
  echo "Golang installed!"
}

### Main
while getopts ":hic" option; do
   case $option in
      h) # display Help
         Help;;
      i) # Install Golang
         InstallGolang;;
      c) # Clean the build directory
         Clean;;
     \?) # Invalid option
         echo "Error: Invalid option"
         Help;;
   esac
done

# Build the project binaries
echo "Building the project binaries..."
mkdir -p $BIN_DIR
go build -o $BIN_DIR/build ./build.go
go build -o $BIN_DIR/run ./run.go

echo "Everything Done!"
echo "Built Binaries are located in $BIN_DIR"