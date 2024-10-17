#!/bin/sh

KG_HOME="/opt/kirogaz"
AIR2="false"
DLNA="true"

CDSP="false"

WPA_CONF=/etc/wpa_supplicant/wpa_supplicant.conf

aoutcfg_30="out_dac4.conf"
aoutcfg_31="out_cdsp_2ch.conf"
aoutcfg_32="out_ldsp_2to21.conf"
aoutcfg_33="out_ldsp_2to31.conf"
aoutcfg_34="out_ldsp_2to41.conf"

OUT_ACTIVE=$(jq -r .OUT_active $KG_HOME/config.json)
WIFI_AP=$(jq -r .WIFI_AP_Mode $KG_HOME/config.json)

alsa_conf() {
    #configure ALSA output
    echo "Do ALSA output configuration"
    eval alsa_outcfg="\$aoutcfg_${OUT_ACTIVE}"
    echo "cfg: $alsa_outcfg"
    cp $KG_HOME/alsa/$alsa_outcfg /tmp/asound.conf

    if [ "$OUT_ACTIVE" == "31" ]; then
        CDSP="true"
    fi

}

start_services() {
    echo "Starting services.."
    if [ "$CDSP" = "true" ]; then
        echo "CamillaDSP"
        /etc/init.d/cdsp start
    fi
    # $KG_HOME/udc.sh start
    if [ "$AIR2" = "true" ]; then
        echo "AirPlay2"
        /etc/init.d/shairport2 start
    else
        echo "AirPlay1"
        /etc/init.d/shairport start
    fi
    if [ "$DLNA" = "true" ]; then
        echo "GMediaRender"
        (sleep 10 ; /etc/init.d/gmrenderer restart) &
    fi
    echo "Starting services done.."
}

stop_services() {
    echo "Stopping services.."
    /etc/init.d/cdsp stop
    #$RIVA_HOME/udc.sh stop
    /etc/init.d/gmrenderer stop
    killall -9 /usr/local/bin/gmediarender
    if [ "$AIR2" = "true" ]; then
        /etc/init.d/shairport2 stop
    else
        /etc/init.d/shairport stop
    fi
    echo "Stopping services done.."
}

wifi_on() {
    ifconfig wlan0 up
    wpa_supplicant -B -i wlan0 -c /etc/wpa_supplicant/wpa_supplicant.conf
    udhcpc -i wlan0 &
}

wifi_off() {
    killall udhcpc
    killall wpa_supplicant
    ifconfig wlan0 down
}

wifi_ap_on() {
    hostapd -B $KG_HOME/hostapd.conf
    udhcpd $KG_HOME/dhcpdwlan.conf
    ifconfig wlan0 192.168.0.1 netmask 255.255.255.0 up
}

wifi_ap_off() {
    killall hostapd
    killall dhcpdwlan
    ifconfig wlan0 down
}

case "$1" in
  start)
    alsa_conf
    start_services
    ;;

  stop)
    stop_services
    ;;

  alsa_update_cfg)
    alsa_conf
    "$0" restart
    ;;

  restart)
    "$0" stop
    "$0" start
    ;;

  cdsp_stop)
    /etc/init.d/cdsp stop
    ;;

  terminate)
    stop_services
    killall alsa_signal
    killall alsaloop
    ;;

  wifiap_start)
    wifi_off
    sleep 1
    wifi_ap_on
    exit 0
    ;;

  wifiap_stop)
    wifi_ap_off
    wifi_on
    exit 0
    ;;

  net_start)
    if [ "$WIFI_AP" == "true" ]; then
        wifi_ap_on
    else
        wifi_on
    fi
    ;;

  *)
    echo "Usage: $0 {start|stop|restart|terminate|cdsp_stop|wifiap_start|wifiap_stop"
    exit 1
esac

exit $?