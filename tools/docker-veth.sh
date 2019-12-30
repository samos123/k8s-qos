#!/bin/bash

set -xe

veth=""

containers=$(docker ps --format '{{.ID}} {{.Names}}' "$@")

get_veth () {
    # This function expects docker container ID as the first argument
    veth=""
    networkmode=$(docker inspect -f "{{.HostConfig.NetworkMode}}" $1)
    echo $networkmode
    if [ "$networkmode" == "host" ]; then
        veth="host"
    else
        pid=$(docker inspect --format '{{.State.Pid}}' "$1")
        ifindex=$(nsenter -t $pid -n ip link | sed -n -e 's/.*eth0@if\([0-9]*\):.*/\1/p')
        veth=$(ip -o link | grep ^$ifindex | sed -n -e 's/.*\(veth[[:alnum:]]*@if[[:digit:]]*\).*/\1/p')
    fi
}

while IFS= read -r line
do
    containerid=$(echo $line | awk '{ print $1 }')
    get_veth $containerid
    echo "$veth $line"
done <<< "$containers"
