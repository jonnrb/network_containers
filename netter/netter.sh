#!/bin/bash

function rm-tables () {
	echo "[util] rm-tables"

	ip rule |grep docker-netter |sed 's/^[0-9]*\:[ \t]*//g' |while read rule; do
		ip rule del ${rule}
	done

	cat /etc/iproute2/rt_tables |grep -v docker-netter > /etc/iproute2/rt_tables.2
	mv /etc/iproute2/rt_tables{.2,}
}

function create-table () {
	local name=${1}

	echo "[util] create-table.${name}"

	for i in {100..200}; do
		if ! cat /etc/iproute2/rt_tables |cut -f 1 -d '	' |grep "${i}" >/dev/null; then

			printf '%d\t%s\n' "${i}" "${name}" >> /etc/iproute2/rt_tables

			# Clean up existing default route if one exists
			ip route del default table "${name}"

			return

		fi
	done
}

function flush-chains () {
	echo "[util] flush-chains"

	iptables -t filter -F docker-netter
	iptables -t mangle -F docker-netter
}

function router () {

	local container=${1}
	local network=${2}


	echo "[cmd] router:container.${container}:docker.${network}"

	local br=$(/go/bin/get_bridge_name "${network}")
	if [ -z "${br}" ]; then
		echo "[error] bridge for ${network} doesn't exist"
		return -1
	fi

	local router_ip=$(/go/bin/get_container_ip "${container}" "${network}")
	if [ -z "${router_ip}" ]; then
		echo "[error] ${container} has no ip on ${network}"
		return -1
	fi
	
	local mark="0x$(echo "${br}" |cut -b 4-11)"

	echo "[debug] br=${br} router_ip=${router_ip} mark=${mark}"


	create-table "docker-netter.${network}"
	iptables -t filter -A docker-netter -j ACCEPT -i "${br}" -o "${br}" -m comment --comment "netter.router:docker.${network}"
	iptables -t mangle -A docker-netter -j MARK -i "${br}" --set-mark "${mark}" -m comment --comment "netter.router:docker.${network}"
	ip rule add fwmark "${mark}" lookup "docker-netter.${network}"
	ip route add default via "${router_ip}" table "docker-netter.${network}"

}

function allow-x-routing () {
	local network=${1}
	local if=${2}

	echo "[cmd] allow-x-routing:docker.${network}:if.${if}"

	local br=$(/go/bin/get_bridge_name "${network}")
	if [ -z "${br}" ]; then
		echo "[error] bridge for ${network} doesn't exist"
		return -1
	fi

	iptables -t filter -A docker-netter -j ACCEPT -i "${br}" -o "${if}" -m comment --comment "netter.allow-x-routing:docker.${network}:if.${if}"
}

function allow-routing () {
	local network=${1}

	echo "[cmd] allow-routing:docker.${network}"

	local br=$(/go/bin/get_bridge_name "${network}")
	if [ -z "${br}" ]; then
		echo "[error] bridge for ${network} doesn't exist"
		return -1
	fi

	iptables -t filter -A docker-netter -j ACCEPT -i "${br}" -o "${br}" -m comment --comment "netter.allow-routing:docker.${network}"
}

function net-slave () {
	local network=${1}
	local dev=${2}

	echo "[cmd] net-slave:docker.${network}:phys.${dev}"

	local br=$(/go/bin/get_bridge_name "${network}")
	if [ -z "${br}" ]; then
		echo "[error] bridge for ${network} doesn't exist"
		return -1
	fi

	ip link set "${dev}" master "${br}"
}

function main () {
	flush-chains
	rm-tables

	for f in $@; do
		cat $f |grep -Ev '^\w*(\#|$)' |while read cmd args; do
			case $cmd in
			router)
				router $args
			;;
			allow-routing)
				allow-routing $args
			;;
			allow-x-routing)
				allow-x-routing $args
			;;
			net-slave)
				net-slave $args
			;;
			esac
		done
	done
}


main $(find /data/*.conf)
