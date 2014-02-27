#!/bin/bash
#
# Installs/updates the local dependencies of this project.
#
# Usage:
# ./install_dependencies.sh
#
# Authors:
# - Florin Patan <florin.patan@motain.de>
# - Caio Moritz Ronchi <caiomoritz.ronchi@motain.de>

go_get() {
    echo "$1"
    go get -u "$1"
}

dependencies="\
github.com/motain/mux
github.com/motain/mysql
github.com/motain/gorp

while read -a dependency; do
    go_get "$dependency"
done <<< "$dependencies"
