#!/bin/sh

function error() {
    echo "[$(date +'%Y-%m-%dT%H:%M:%S%z')]: $@" >&2
}

GFWLIST=$(curl -s https://raw.githubusercontent.com/gfwlist/gfwlist/master/gfwlist.txt | base64 -d)
if [[ $? != 0 ]]; then
    error "get gfwlist failed"
    exit -1
fi

PROXY=$(cat proxy.list | sort | uniq)
if [[ $? != 0 ]]; then
    error "get proxylist failed"
    exit -1
fi

echo -e "$GFWLIST\n\n$PROXY" | base64 >gfwlist.txt
echo "$PROXY" >proxy.list
