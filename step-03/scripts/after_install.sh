#!/bin/bash

# set golang path
echo "export PATH=$PATH:/usr/local/go/bin" >> $HOME/.bashrc

# activate changes
source $HOME/.bashrc

# go to golang app
cd /webapps/goapp/

# install golang dependency
go mod tidy

# make binary
go build -o goappbinary
go install
