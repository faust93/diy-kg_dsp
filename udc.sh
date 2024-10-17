#!/bin/sh

KG_HOME="/opt/kirogaz"

readonly GADGET_BASE_DIR="/sys/kernel/config/usb_gadget/g_audio"
readonly DEV_ETH_ADDR="48:6f:73:74:50:43"
readonly HOST_ETH_ADDR="42:61:64:55:53:42"

USB_ETH="false"
UAC="2"

[[ -z "$UAC" ]] && UAC="1"

AUDIO_CHANNEL_MASK=3
#AUDIO_SAMPLE_RATES=44100,48000,88200,96000
AUDIO_SAMPLE_RATES=96000
AUDIO_SAMPLE_SIZE=2
[[ "$UAC" == "2" ]] && AUDIO_SAMPLE_SIZE=4

case "$1" in
    start)
            modprobe libcomposite

            USB_IP="192.168.0.1"
            USB_MASK="255.255.255.0"

            cwd=$(pwd)

            mkdir "${GADGET_BASE_DIR}"
            cd "${GADGET_BASE_DIR}"

            echo 0x1d6b > idVendor # Linux Foundation
            echo 0x0104 > idProduct # Multifunction Composite Gadget
            echo 0x0100 > bcdDevice # v1.0.0
            echo 0x0200 > bcdUSB # USB2

            #echo 0x01   > bDeviceClass
            #echo 0x02   > bDeviceSubClass
            #echo 0x20   > bDeviceProtocol

            mkdir -p strings/0x409
            echo "2023080101F93" > strings/0x409/serialnumber
            echo "DIYAudio" > strings/0x409/manufacturer
            echo "Kirogaz" > strings/0x409/product

            mkdir -p configs/c.1/strings/0x409
            echo 120 > configs/c.1/MaxPower
            echo "Audio" > configs/c.1/strings/0x409/configuration
            #echo 1       > os_desc/use
            #echo 0xcd    > os_desc/b_vendor_code
            #echo MSFT100 > os_desc/qw_sign

            mkdir -p functions/uac${UAC}.usb0

            echo $AUDIO_CHANNEL_MASK > functions/uac${UAC}.usb0/c_chmask
            echo $AUDIO_SAMPLE_RATES > functions/uac${UAC}.usb0/c_srate
            echo $AUDIO_SAMPLE_SIZE > functions/uac${UAC}.usb0/c_ssize
            echo $AUDIO_CHANNEL_MASK > functions/uac${UAC}.usb0/p_chmask
            echo $AUDIO_SAMPLE_RATES > functions/uac${UAC}.usb0/p_srate
            echo $AUDIO_SAMPLE_SIZE > functions/uac${UAC}.usb0/p_ssize
            ln -s functions/uac${UAC}.usb0 configs/c.1

            # Ethernet
            if [ "$USB_ETH" == "true" ]; then
                mkdir functions/ecm.usb0
                echo "${DEV_ETH_ADDR}" > functions/ecm.usb0/dev_addr
                echo "${HOST_ETH_ADDR}" > functions/ecm.usb0/host_addr
                ln -s functions/ecm.usb0 configs/c.1/
            fi

            ls /sys/class/udc > UDC
            [ $? = 0 ] || (cd ${cwd} && $0 stop; exit 1)

            if [ "$USB_ETH" == "true" ]; then
                ifconfig usb0 $USB_IP netmask $USB_MASK up
            fi
            exit 0
            ;;
    stop)
            [ -d $GADGET_BASE_DIR ] || exit

            echo '' > $GADGET_BASE_DIR/UDC

            echo "Removing strings from configurations"
            for dir in $GADGET_BASE_DIR/configs/*/strings/*; do
                [ -d $dir ] && rmdir $dir
            done

            echo "Removing functions from configurations"
            for func in $GADGET_BASE_DIR/configs/*.*/*.*; do
                [ -e $func ] && rm $func
            done

            echo "Removing configurations"
            for conf in $GADGET_BASE_DIR/configs/*; do
                [ -d $conf ] && rmdir $conf
            done

            echo "Removing functions"
            for func in $GADGET_BASE_DIR/functions/*.*; do
                [ -d $func ] && rmdir $func
            done

            echo "Removing strings"
            for str in $GADGET_BASE_DIR/strings/*; do
                [ -d $str ] && rmdir $str
            done

            echo "Removing gadget"
            rmdir $GADGET_BASE_DIR
            ;;

    restart|reload)
        "$0" stop
        "$0" start
        ;;
    *)
        echo "Usage: $0 {start|stop|restart}"
        exit 1
esac
