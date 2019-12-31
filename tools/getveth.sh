#!/usr/bin/env bash

veth=""

get_veth () {
    # This function expects docker container ID as the first argument
    networkmode=$(docker inspect -f "{{.HostConfig.NetworkMode}}" $1)
    if [ "$networkmode" == "host" ]; then
        veth="host"
    else
        pid=$(docker inspect --format '{{.State.Pid}}' "$1")
        ifindex=$(nsenter -t $pid -n ip link | sed -n -e 's/.*eth0@if\([0-9]*\):.*/\1/p')
        if [ -z "$ifindex" ]; then
            veth="not_found"
        else
            veth=$(ip -o link | grep ^$ifindex | sed -n -e 's/.*\(veth[[:alnum:]]*@if[[:digit:]]*\).*/\1/p')
        fi
    fi
}

get_veth $1

echo "$veth"
