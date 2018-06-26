#!/bin/sh

function process () {
  f="$1"
  out="/tmp/$(basename "$f" | sed 's/\.tmpl$//')"

  if ! [ -f "$f" ]; then
    echo "[init.sh] [error] no $f found; exiting" >&2
    exit 1
  fi

  echo "[init.sh] [info] reading configuration template" >&2
  sed_script=$(cat "$f" |grep -oE 'docker-network:.*$' |while read target; do
    # Get the network name from the second half of the directive.
    network_name="$(echo "${target}" |cut -b 16-)"

    # Use utils to get the bridge name.
    bridge_name="$(/bin/get_bridge_name "${network_name}")"
    if [ -z "${bridge_name}" ]; then
      echo "[init.sh] [error] couldn't get bridge name for \"${network_name}\"; exiting" >&2
      exit 1
    fi

    # Add a replacement commmand to our list of sed commands.
    printf "; s/${target}/${bridge_name}/g"
  done)

  # Create the hostapd config file
  echo "[init.sh] [info] writing configuration $out" >&2
  cat "$f" |sed -e "${sed_script}" > "$out"

  echo "$out"
}

if [ -z "$1" ]; then
  tmpl="/data/hostapd.conf.tmpl"
else
  tmpl="$@"
fi

conf="$(for f in $tmpl; do
  process "$f"
done)"

exec hostapd $conf
