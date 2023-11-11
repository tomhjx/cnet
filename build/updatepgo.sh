#!/bin/bash
root=$(dirname $(dirname "$0"))
wget "http://127.0.0.1:6060/debug/pprof/profile" -O ${root}/cmd/cnet/default.pgo