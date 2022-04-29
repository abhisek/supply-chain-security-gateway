#!/bin/bash

python3 -m pip config --user set global.index http://localhost:10000/pypi
python3 -m pip config --user set global.index-url http://localhost:10000/pypi/simple
python3 -m pip config --user set global.trusted-host localhost
