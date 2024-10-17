# go to configfs directory for USB gadgets
CONFIGFS_ROOT=/sys/kernel/config # adapt to your machine
cd "${CONFIGFS_ROOT}"/usb_gadget

# create gadget directory and enter it
mkdir g1
cd g1

# USB ids
echo 0x1d6b > idVendor
echo 0x104 > idProduct

# USB strings, optional
mkdir strings/0x409 # US English, others rarely seen
echo "Collabora" > strings/0x409/manufacturer
echo "ECM" > strings/0x409/product

# create the (only) configuration
mkdir configs/c.1

# create the (only) function
mkdir functions/ecm.usb0

# assign function to configuration
ln -s functions/ecm.usb0/ configs/c.1/

# bind!
echo ff400000.usb > UDC # ls /sys/class/udc to see available UDCs

ifconfig usb0 192.168.10.1 netmask 255.255.255.0 up