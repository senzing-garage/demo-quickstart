#!/bin/zsh

export DYLD_LIBRARY_PATH=/opt/senzing/g2/lib:/opt/senzing/g2/lib/macos
export LD_LIBRARY_PATH=${DYLD_LIBRARY_PATH}

"$@"
