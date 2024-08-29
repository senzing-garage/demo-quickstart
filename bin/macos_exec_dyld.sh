#!/bin/zsh

export DYLD_LIBRARY_PATH=/opt/senzing/er/lib:/opt/senzing/er/lib/macos
export LD_LIBRARY_PATH=${DYLD_LIBRARY_PATH}

"$@"
