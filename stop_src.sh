#!/bin/sh
echo 0 > /tmp/mute
killall rca1_src.sh
killall rca2_src.sh
killall spdif_src.sh
killall alsa_signal
killall alsaloop
echo 1 > /tmp/mute
exit 0
