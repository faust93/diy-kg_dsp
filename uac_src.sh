#!/bin/sh

IN=`aplay -l | grep UAC | cut -f2 -d' ' | tr -d :`

case "$1" in
    on)
        echo "UAC capture"
        {
            alsaloop -C uac_direct --rate=96000 -f S32_LE -P default -t 50000 -S 3 -c 2
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
        ./udc.sh stop
        exit 0
        ;;

    listen)
        killall alsa_signal
        ./udc.sh start
        {
            ./alsa_signal uac_direct s32le:44100:2 18.0 1:5 "./$(basename $0) on" "./$(basename $0) off"
        } >/dev/null 2>&1 &
        exit 0
        ;;
esac