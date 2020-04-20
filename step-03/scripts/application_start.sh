#!/bin/bash

# go to golang app
cd /webapps/goapp/

# start application
# sudo -E ./app
./goapp > /dev/null 2> /dev/null < /dev/null &
