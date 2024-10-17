package main

var left_arrow = []byte{ 0x00, 0x00, 0x18, 0x3c, 0x7e, 0xff, 0x00, 0x00 }

// 'error', 24x24px
var error_img = []byte{0x00, 0x00, 0x00, 0x80, 0xe0, 0xf0, 0xf0, 0xf8, 0xf8, 0xfc, 0xfc, 0x3c, 0x3c, 0xfc, 0xfc, 0xf8,
    0xf8, 0xf0, 0xf0, 0xe0, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x7e, 0xff, 0xff, 0xff, 0xff, 0xff,
    0xff, 0xff, 0xff, 0xc0, 0xc0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7e, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x01, 0x07, 0x0f, 0x0f, 0x1f, 0x1f, 0x3f, 0x3f, 0x3c, 0x3c, 0x3f, 0x3f, 0x1f,
    0x1f, 0x0f, 0x0f, 0x07, 0x01, 0x00, 0x00, 0x00 }

// 'question/warning', 24x24px
var warning_img = []byte{0x00, 0x00, 0x00, 0x80, 0xe0, 0xf0, 0xf0, 0xf8, 0xf8, 0x7c, 0x7c, 0x7c, 0x7c, 0x7c, 0x7c, 0xf8,
    0xf8, 0xf0, 0xf0, 0xe0, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x7e, 0xff, 0xff, 0xff, 0xff, 0xff,
    0xf8, 0xf8, 0xfe, 0x9e, 0x8e, 0xe6, 0xe0, 0xf0, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7e, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x01, 0x07, 0x0f, 0x0f, 0x1f, 0x1f, 0x3f, 0x3f, 0x3c, 0x3c, 0x3f, 0x3f, 0x1f,
    0x1f, 0x0f, 0x0f, 0x07, 0x01, 0x00, 0x00, 0x00 }

// 'speaker', 21x21px
var speaker_img = []byte{
    0xc0, 0xe0, 0xe0, 0xe0, 0xe0, 0xe0, 0xe0, 0xf0, 0xf8, 0xfc, 0xfe, 0xff, 0xff, 0x00, 0x20, 0x70, 
    0x38, 0x38, 0x18, 0x00, 0x00, 0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 
    0xff, 0xff, 0x00, 0x80, 0xc0, 0x8e, 0x8e, 0x0e, 0x0e, 0x0e, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
    0x00, 0x01, 0x03, 0x07, 0x0f, 0x1f, 0x1f, 0x00, 0x00, 0x01, 0x03, 0x03, 0x03, 0x00, 0x00 }

// 'speaker1', 16x16px
var speaker1_img = []byte{
    0xf0, 0xf0, 0xf0, 0xf0, 0xf0, 0xf8, 0xfc, 0xfe, 0xff, 0xff, 0x00, 0x18, 0x98, 0x8c, 0x80, 0x80, 
    0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x1f, 0x3f, 0x7f, 0xff, 0xff, 0x00, 0x18, 0x19, 0x31, 0x01, 0x01 }

var wifi_off_img = []byte{// 'wifi-off-svgrepo-com', 18x18px
    0x00, 0x00, 0x80, 0xc8, 0x70, 0x20, 0xc0, 0x80, 0x10, 0x10, 0x10, 0xb0, 0x20, 0x60, 0xc0, 0x80,
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x03, 0x01, 0x0c, 0x25, 0x26, 0x0c, 0x08, 0x11, 0x22,
    0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00 }

var wifi_on_img = []byte{
    0x00, 0x00, 0x80, 0xc0, 0x60, 0x20, 0xb0, 0x90, 0x90, 0x90, 0x90, 0xb0, 0x20, 0x60, 0xc0, 0x80,
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x03, 0x01, 0x0c, 0x24, 0x24, 0x0c, 0x09, 0x03, 0x02,
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00 }

// 'wifiap', 18x18px
var wifi_ap_img = []byte{
    0x00, 0x00, 0xf0, 0x18, 0xe4, 0x1a, 0x0b, 0xc5, 0xe5, 0xe5, 0xc5, 0x0b, 0x1a, 0xe4, 0x18, 0xf0, 
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xc0, 0x70, 0x0f, 0x03, 0x03, 0x1d, 0x70, 0xc0, 0x00, 
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x03, 0x02, 0x02, 0x02, 0x02, 0x02, 0x02, 
    0x03, 0x02, 0x00, 0x00, 0x00, 0x00 }

var amp_on_img = []byte{
    0x00, 0x00, 0xc0, 0x60, 0x60, 0x60, 0x30, 0x18, 0x08, 0xfc, 0x00, 0x00, 0xc0, 0x00, 0x30, 0xe0, 
    0x00, 0x00, 0x00, 0x00, 0x0f, 0x18, 0x18, 0x10, 0x30, 0x60, 0x40, 0xff, 0x00, 0x00, 0x0f, 0x00, 
    0x30, 0x1f, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00 }

var amp_off_img = []byte{
    0x00, 0x00, 0xc0, 0x60, 0x60, 0x60, 0x30, 0x18, 0x08, 0xfc, 0x00, 0x00, 0xc0, 0x80, 0x00, 0x80, 
    0xc0, 0x00, 0x00, 0x00, 0x0f, 0x18, 0x18, 0x10, 0x30, 0x60, 0x40, 0xff, 0x00, 0x00, 0x0c, 0x07, 
    0x03, 0x07, 0x0c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00 }


var amp2_off_img = []byte{
// 'amp2-off', 18x18px
    0x00, 0xc0, 0xc0, 0xc0, 0xc0, 0xc0, 0xe0, 0xf0, 0xf8, 0xfc, 0xfc, 0x00, 0x00, 0x80, 0x00, 0x00, 
    0x80, 0x00, 0x00, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x1f, 0x3f, 0x7f, 0xff, 0xff, 0x00, 0x00, 0x04, 
    0x03, 0x03, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00 }

var amp2_on_img = []byte{
// 'amp2-on', 18x18px
    0x00, 0xc0, 0xc0, 0xc0, 0xc0, 0xc0, 0xe0, 0xf0, 0xf8, 0xfc, 0xfc, 0x00, 0x00, 0x80, 0x00, 0xc0, 
    0x00, 0x00, 0x00, 0x0f, 0x0f, 0x0f, 0x0f, 0x0f, 0x1f, 0x3f, 0x7f, 0xff, 0xff, 0x00, 0x00, 0x07, 
    0x00, 0x0f, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00 }