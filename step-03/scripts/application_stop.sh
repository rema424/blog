#!/bin/bash

pgrep goappbinary && kill $(pgrep goappbinary) || echo 'goapp not runnning'
