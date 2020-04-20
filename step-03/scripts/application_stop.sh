#!/bin/bash

pgrep goapp && pkill goapp || echo 'goapp not runnning'
