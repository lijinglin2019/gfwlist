#!/bin/sh

function error() {
    echo "[$(date +'%Y-%m-%dT%H:%M:%S%z')]: $@" >&2
}

GFW_LIST_BASE64=$(curl -s https://raw.githubusercontent.com/gfwlist/gfwlist/master/gfwlist.txt)
if [[ $? != 0 ]]; then
    error "get gfwlist failed"
    exit -1
fi

GFW_LIST=$(echo "$GFW_LIST_BASE64" | base64 -d)
if [[ $? != 0 ]]; then
    error "decode gfwlist failed"
    exit -1
fi

PROXY_LIST=$(cat proxy.list | sort | uniq)
if [[ $? != 0 ]]; then
    error "get proxylist failed"
    exit -1
fi
PROXY_LIST=$(echo "$PROXY_LIST" | sort | uniq)

echo -e "$GFW_LIST\n\n$PROXY_LIST" | base64 >proxylist.txt
echo "$PROXY_LIST" >proxy.list
