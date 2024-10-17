#!/bin/sh

IN="in2_direct"

case "$1" in
    on)
        echo "SPDIF(RCA3) capture"
        {
            sleep 0.5
            alsaloop -C $IN --rate=96000 -f S32_LE -P default -t 50000 -S 3 -c 2
        } >/dev/null 2>&1 &
        exit 0
        ;;

    off)
        killall alsaloop
        #killall `basename $0`
        exit 0
        ;;

    stop)
        killall alsa_signal
        killall alsaloop
        #killall `basename $0`
        exit 0
        ;;

    listen)
        killall alsa_signal
        {
            sleep 0.5
            ./alsa_signal $IN s32le:44100:2:128:500000 14.0 2:32 "./$(basename $0) on" "./$(basename $0) off"
        } >/dev/null 2>&1 &
        exit 0
        ;;
esac