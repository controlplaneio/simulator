#!/bin/bash

set -Eeuo pipefail

## Flag 1
mkdir -p /opt/treasure/.pirates-loot/{🏴,☠️,🦜}/{davy/jones\'/locker,ᶜᵃᵖᵗᵃⁱⁿ/Hλ$ħ𝔍Ⱥ¢ks/quarters,quarterdeck,crows-nest}
mkdir -p /opt/treasure/.pirates-loot/{🏴,☠️,🦜}/boltholes/{st-marys,tortuga}

DIR=$(find /opt/treasure/.pirates-loot -type d | sort -R | head -n 1)
echo "flag_ctf{It'sAPiratesLifeForMe}" > "$DIR/captains-booty.hash"
chmod 444 "$DIR/captains-booty.hash"

## Flag 2
echo "flag_ctf{YarrArrghGarrrAhoyMeHarties}" > /root/captains-prize.hash
chmod 400 /root/captains-prize.hash
