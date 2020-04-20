#!/bin/bash

pgrep app && pkill app || echo 'app not runnning'
