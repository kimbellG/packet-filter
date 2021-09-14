#!/bin/bash

[ $# -ge 1 ] || { echo >&2 "Usage: $0 {up|down};" exit 1; }
[ $UID -eq 0 ] || { echo >&2 "This script requires root."; exit 1; }

action="$1"
netns=xdp-test
local=xdp-local
remote=xdp-remote


case "$action" in 
up)
	ip netns add ${netns}
	ip link add ${remote} type veth peer name ${local}
	ip link set ${remote} netns ${netns}
	ip netns exec ${netns} ip address add 192.0.2.2/24 dev ${remote}
	ip netns exec ${netns} ip link set ${remote} up
	ip address add 192.0.2.1/24 dev ${local}
	ip link set ${local} up
	;;
down)
	ip netns delete ${netns}
	ip link delete ${local}
	;;
esac

