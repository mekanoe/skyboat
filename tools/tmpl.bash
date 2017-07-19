#!/bin/bash
set -e
# usage tmpl.bash dockerfile

SCRIPT_DIR=$(dirname $0)
TMPL_DIR="$SCRIPT_DIR/templates"

cmd () {
    echo "> $@"
    "$@"
}

exec_tmpl () {
    cmd bash "$1"
}

main () {
    if [[ -e "$TMPL_DIR/$1.tmpl.bash" ]]; then
        exec_tmpl "$TMPL_DIR/$1.tmpl.bash"
    fi
}

main "$@"