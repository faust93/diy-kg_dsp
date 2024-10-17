#!/bin/sh
# BT sink

case "$1" in
    start)
            echo -n "Starting bluetooth sink"
            /etc/init.d/bluetooth start
            sleep 1
            #/usr/bin/hciconfig hci0 piscan
            #/usr/bin/hciconfig hci0 sspmode 1
            sleep 1
            /usr/bin/bluetoothctl power on
            /usr/bin/bluetoothctl discoverable on
            #bluetoothctl --agent=NoInputNoOutput
            #echo "agent NoInputNoOutput" | bluetoothctl
            #echo "default-agent" | bluetoothctl
            /usr/local/bin/a2dp-agent3 &
            a2dp_pid=$!
            sleep 1
            #bluealsa -p a2dp-source -p a2dp-sink &
            /usr/bin/bluealsa -i hci0 -p a2dp-sink &
            bluealsa_pid=$!
            sleep 1
            /usr/bin/bluealsa-aplay --profile-a2dp --single-audio --pcm-buffer-time=1000000 --mixer-name="Line Out" 00:00:00:00:00:00 &
            aplay_pid=$!

            kill -0 $a2dp_pid
            [ $? = 0 ] || ($0 stop; exit 1)
            kill -0 $bluealsa_pid
            [ $? = 0 ] || ($0 stop; exit 1)
            kill -0 $aplay_pid
            [ $? = 0 ] || ($0 stop; exit 1)
            exit 0
            ;;

    stop)
            echo -n "Stopping bluetooth sink"
            killall bluealsa-aplay
            killall bluealsa
            killall a2dp-agent3
            /etc/init.d/bluetooth stop
            ;;

    *)
        echo "Usage: $0 {start|stop}"
        exit 1
esac
