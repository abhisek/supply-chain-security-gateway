#!/bin/bash

set -e

cd $(dirname $0)
mkdir -p .cache

if [ ! -f ".cache/plantuml.jar" ]; then
  wget -O .cache/plantuml.jar \
    https://github.com/plantuml/plantuml/releases/download/v1.2022.4/plantuml-1.2022.4.jar
fi;

java -jar .cache/plantuml.jar -o ./images/ *.plantuml
