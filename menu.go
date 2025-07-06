// menu handlers
package main

import (
    "log"
    "strings"
    "strconv"
    "math"
    "fmt"
    "time"
    )

const (
    I_BACK = 0
    I_SCREEN = 1
    I_FUNC = 2

    I_SWVOL = 10
    I_CVOL = 11
    I_FRONTVOL = 12
    I_SIDEVOL = 13
    I_DAC4VOL = 14

    I_ADC_RCA1VOL = 15
    I_ADC_RCA2VOL = 16
    I_ADC_SPDIFVOL = 17

    I_RCA1_IN = 19
    I_RCA2_IN = 20
    I_SPDIF_IN = 21
    I_STREAM_IN = 22
    I_UAC_IN = 23

    I_LINE_OUT = 30
    I_CDSP_OUT = 31
    I_LDSP21_OUT = 32
    I_LDSP31_OUT = 33
    I_LDSP41_OUT = 34

    I_CDSP_CFILE = 40
    I_CDSP_STAT = 41
    I_CDSP_SAVE = 42
    I_CDSP_CFG_SAVE = 43
    I_CDSP_SIGHUP = 44

    I_WIFI_AP = 50
    I_USB_ETH = 51
    // amp idle timeout
    I_AMP_MUTE = 52
    // OLED Signal Meter
    I_SIGNAL_METR = 53

    I_REBOOT  = 80
    I_HALT    = 81

    MAIN_MODE = 0x7F
    NAV_MODE  = 0x80
    CTL_MODE  = 0x81
)

// hw mixer controls
const (
    ALSA_CARD = 1
    DAC1_VOL_CTL = 14
    DAC2_VOL_CTL = 15
    DAC3_VOL_CTL = 16
    DAC4_VOL_CTL = 17
    ADC1_VOL_CTL = 18
    ADC2_VOL_CTL = 19
    ADC3_VOL_CTL = 11

    )

// amp mute timeouts (seconds)
const (
    AMP_MUTE_T1 = 60
    AMP_MUTE_T2 = 300
    AMP_MUTE_T3 = 600
)

var evtTime time.Time

var ackMsg bool = false

var cdsp_configs []string

type menuItem struct {
    name string
    x int16
    y int16
    itype int
    iarg1 int
    iarg2 int
    ifunc func()
}

type menuPage struct {
    id int
    items []menuItem
}

type MainMenu struct {
    screens []menuPage

    active_screen int
    prev_screen int
    cursor_pos int

    mode int
    ctl_id int
}

var menu MainMenu

func MenuInit() {
    menu = MainMenu{
        screens: []menuPage{
                {
                    id: 0,
                    items: []menuItem{
                        { name: "DAC volume",
                          x: 11, y: 1,
                          itype: I_SCREEN,
                          iarg1: 1, //screen id
                        },
                        { name: "ADC volume",
                          x: 11, y: 11,
                          itype: I_SCREEN,
                          iarg1: 2, //screen id
                        },
                        { name: "IN select",
                          x: 11, y: 21,
                          itype: I_SCREEN,
                          iarg1: 3,
                        },
                        { name: "OUT select",
                          x: 11, y: 31,
                          itype: I_SCREEN,
                          iarg1: 4,
                        },
                        { name: "CDSP config",
                          x: 11, y: 41,
                          itype: I_SCREEN,
                          iarg1: 5,
                        },
                        { name: "Options",
                          x: 11, y: 51,
                          itype: I_SCREEN,
                          iarg1: 6,
                        },
                        { name: "Reboot",
                          x: 11, y: 1,
                          itype: I_REBOOT,
                          iarg1: 7,
                        },
                        { name: "Power Off",
                          x: 11, y: 11,
                          itype: I_HALT,
                          iarg1: 7,
                        },
                        { name: "< Back",
                          x: 1, y: 21,
                          itype: I_BACK,
                          iarg1: MAIN_MODE,
                        },
                    },
                },
                {
                    id: 1, //DAC vol
                    items: []menuItem{
                        { name: "SW",
                          x: 11, y: 1,
                          itype: I_SWVOL,
                        },
                        { name: "CEN",
                          x: 11, y: 11,
                          itype: I_CVOL,
                        },
                        { name: "FL/FR",
                          x: 11, y: 21,
                          itype: I_FRONTVOL,
                        },
                        { name: "SL/SR",
                          x: 11, y: 31,
                          itype: I_SIDEVOL,
                        },
                        { name: "DAC4",
                          x: 11, y: 41,
                          itype: I_DAC4VOL,
                        },
                        { name: "< Back",
                          x: 1, y: 51,
                          itype: I_BACK,
                        },
                    },
                },
                {
                    id: 2, //ADC vol
                    items: []menuItem{
                        { name: "RCA1 In",
                          x: 11, y: 1,
                          itype: I_ADC_RCA1VOL,
                        },
                        { name: "RCA2 In",
                          x: 11, y: 11,
                          itype: I_ADC_RCA2VOL,
                        },
                        { name: "SPDIF In",
                          x: 11, y: 21,
                          itype: I_ADC_SPDIFVOL,
                        },
                        { name: "< Back",
                          x: 1, y: 31,
                          itype: I_BACK,
                        },

                    },
               },
                {
                    id: 3, //INPUTS
                    items: []menuItem{
                        { name: "RCA1",
                          x: 11, y: 1,
                          itype: I_RCA1_IN,
                        },
                        { name: "RCA2",
                          x: 11, y: 11,
                          itype: I_RCA2_IN,
                        },
                        { name: "SPDIF",
                          x: 11, y: 21,
                          itype: I_SPDIF_IN,
                        },
                        { name: "USB UAC",
                          x: 11, y: 31,
                          itype: I_UAC_IN,
                        },
                        { name: "STREAMER ONLY",
                          x: 11, y: 41,
                          itype: I_STREAM_IN,
                        },
                        { name: "< Back",
                          x: 1, y: 51,
                          itype: I_BACK,
                        },

                    },
                },
                {
                    id: 4, //OUTPUTS
                    items: []menuItem{
                        { name: "CDSP",
                          x: 11, y: 1,
                          itype: I_CDSP_OUT,
                        },
                        { name: "LineOut (DAC4)",
                          x: 11, y: 11,
                          itype: I_LINE_OUT,
                        },
                        { name: "LDSP 2.1ch",
                          x: 11, y: 21,
                          itype: I_LDSP21_OUT,
                        },
                        { name: "LDSP 3.1ch",
                          x: 11, y: 31,
                          itype: I_LDSP31_OUT,
                        },
                        { name: "LDSP 4.1ch",
                          x: 11, y: 41,
                          itype: I_LDSP41_OUT,
                        },
                        { name: "< Back",
                          x: 1, y: 51,
                          itype: I_BACK,
                        },

                    },
                },
                {
                    id: 5, //CDSP
                    items: []menuItem{
                        { name: "Set DSP Config",
                          x: 11, y: 1,
                          itype: I_CDSP_CFILE,
                        },
                        { name: "Save active config",
                          x: 11, y: 11,
                          itype: I_CDSP_CFG_SAVE,
                        },
                        { name: "Save on exit",
                          x: 11, y: 21,
                          itype: I_CDSP_SAVE,
                        },
                        { name: "CDSP Stats",
                          x: 11, y: 31,
                          itype: I_CDSP_STAT,
                        },
                        { name: "Reload",
                          x: 11, y: 41,
                          itype: I_CDSP_SIGHUP,
                        },
                        { name: "< Back",
                          x: 1, y: 51,
                          itype: I_BACK,
                        },

                    },
                },
                {
                    id: 6, //Opts
                    items: []menuItem{
                        { name: "AMP Mute Timeout",
                          x: 11, y: 1,
                          itype: I_SCREEN,
                          iarg1: 7,
                        },
                        { name: "Signal Meter",
                          x: 11, y: 11,
                          itype: I_SIGNAL_METR,
                        },
                        { name: "WIFI AP Mode",
                          x: 11, y: 21,
                          itype: I_WIFI_AP,
                        },
                        { name: "USB ETH Mode",
                          x: 11, y: 31,
                          itype: I_USB_ETH,
                        },
                        { name: "< Back",
                          x: 1, y: 41,
                          itype: I_BACK,
                          iarg1: MAIN_MODE,
                        },

                    },
                },
                {
                    id: 7, //Amp timeout
                    items: []menuItem{
                        { name: "1 min",
                          x: 11, y: 1,
                          itype: I_AMP_MUTE,
                          iarg1: AMP_MUTE_T1,
                        },
                        { name: "5 min",
                          x: 11, y: 11,
                          itype: I_AMP_MUTE,
                          iarg1: AMP_MUTE_T2,
                        },
                        { name: "10 min",
                          x: 11, y: 21,
                          itype: I_AMP_MUTE,
                          iarg1: AMP_MUTE_T3,
                        },
                        { name: "< Back",
                          x: 1, y: 31,
                          itype: I_BACK,
                        },
                    },
                },


            },
        }
    menu.active_screen = 0
    menu.cursor_pos = 0
    menu.mode = MAIN_MODE
    menu.ctl_id = I_SCREEN
}

func MainScreen() {
    var aIN,aOUT string
    switch conf.IN_active {
        case I_RCA1_IN:
            aIN = "RCA1"
        case I_RCA2_IN:
            aIN = "RCA2"
        case I_SPDIF_IN:
            aIN = "SPDIF"
        case I_STREAM_IN:
            aIN = "NWS"
        case I_UAC_IN:
            aIN = "USB"

    }
    switch conf.OUT_active {
        case I_CDSP_OUT:
            aOUT = "CDSP"
        case I_LINE_OUT:
            aOUT = "LINE"
        case I_LDSP21_OUT:
            aOUT = "LDSP2.1"
        case I_LDSP31_OUT:
            aOUT = "LDSP3.1"
        case I_LDSP41_OUT:
            aOUT = "LDSP4.1"
    }
    disp.oledSet_font_inverted(true)
    disp.oledDraw_string(1, 1, aIN, NORMAL_SIZE, WHITE)
    disp.oledDraw_string(int16(132-(len(aOUT)*8)), 1, aOUT, NORMAL_SIZE, WHITE)
    disp.oledSet_font_inverted(false)

    if conf.WIFI_AP_Mode {
        disp.oledDraw_bitmap(105, 10, 18, 18, &wifi_ap_img, WHITE)
    } else {
        if WifiState {
            disp.oledDraw_bitmap(105, 10, 18, 18, &wifi_on_img, WHITE)
        } else {
            disp.oledDraw_bitmap(105, 10, 18, 18, &wifi_off_img, WHITE)
        }
    }

    if AMP_Mute == 1 {
        disp.oledDraw_bitmap(1, 9, 18, 18, &amp_on_img, WHITE)
    } else {
        disp.oledDraw_bitmap(1, 9, 18, 18, &amp_off_img, WHITE)
    }

    if len(DacPRate) != 0 {
        dpRate := DacPRate[:2] + "kHz"
        disp.oledDraw_string(25, 11, dpRate, DOUBLE_SIZE, WHITE)
        disp.oledDraw_string(84, 18, DacPFmt[:3], NORMAL_SIZE, WHITE)
        dpCh := DacPCh + "Ch"
        disp.oledDraw_string(48, 1, dpCh, NORMAL_SIZE, WHITE)

    }

    if conf.OUT_active == I_CDSP_OUT {
        disp.oledDraw_bitmap(3, 47, 16, 16, &speaker1_img, WHITE)
        vol, err := cdspGetVolume()
        if err {
            msgDisplay(27, 52, "CDSP error", 0)
        } else {
            disp.oledDraw_hbar(23, 52, 6, int16(50 - math.Abs(vol)))
            dbvol := fmt.Sprintf("%.1fdB", vol)
            disp.oledDraw_string(76, 52, dbvol, NORMAL_SIZE, WHITE)
        }

        disp.oledDraw_string(1, 33, "C:", NORMAL_SIZE, WHITE)
        cfg_title, err := cdspGetConfigTitle()
        if err {
            msgDisplay(11, 33, "--", 0)
        } else {
            disp.oledDraw_string(11, 33, cfg_title, NORMAL_SIZE, WHITE)
        }

        pl_x := (len(cfg_title)+1)*8
        pload,err := cdspGetPLoad()
        if err {
            pload = 0.0
        }
        pload_s := fmt.Sprintf("L:%.1f", pload)
        disp.oledDraw_string(int16(pl_x), 33, pload_s, NORMAL_SIZE, WHITE)
    }
    disp.oledDisplay()
}

func MenuDisplay() {
    if OLED_off || ackMsg { return }
    disp.oledClear(BLACK)
    switch menu.mode {
        case MAIN_MODE:
            MainScreen()
        case NAV_MODE:
            n_items := len(menu.screens[menu.active_screen].items)
            var offset int = 0
            if menu.cursor_pos > 5 {
                offset = 6 * (menu.cursor_pos/6)
            }
            its := 0
            for i := offset; i < n_items; i++ {
                if its <= 5 {  //max onscreen items
                    disp.oledDraw_string(menu.screens[menu.active_screen].items[i].x,
                             menu.screens[menu.active_screen].items[i].y,
                             menu.screens[menu.active_screen].items[i].name,
                             NORMAL_SIZE, WHITE)
                    switch menu.screens[menu.active_screen].items[i].itype {
                        case I_SWVOL:
                            ivol := fmt.Sprintf("[%d]", conf.Sw_volume)
                            disp.oledDraw_string(94, menu.screens[menu.active_screen].items[i].y,
                                 ivol, NORMAL_SIZE, WHITE)
                        case I_CVOL:
                            ivol := fmt.Sprintf("[%d]", conf.Center_volume)
                            disp.oledDraw_string(94, menu.screens[menu.active_screen].items[i].y,
                                 ivol, NORMAL_SIZE, WHITE)
                        case I_FRONTVOL:
                            ivol := fmt.Sprintf("[%d-%d]", conf.FL_volume, conf.FR_volume)
                            disp.oledDraw_string(70, menu.screens[menu.active_screen].items[i].y,
                                 ivol, NORMAL_SIZE, WHITE)
                        case I_SIDEVOL:
                            ivol := fmt.Sprintf("[%d-%d]", conf.SL_volume, conf.SR_volume)
                            disp.oledDraw_string(70, menu.screens[menu.active_screen].items[i].y,
                                 ivol, NORMAL_SIZE, WHITE)
                        case I_DAC4VOL:
                            ivol := fmt.Sprintf("[%d-%d]", conf.DL_volume, conf.DR_volume)
                            disp.oledDraw_string(70, menu.screens[menu.active_screen].items[i].y,
                                 ivol, NORMAL_SIZE, WHITE)

                        case I_ADC_RCA1VOL:
                            ivol := fmt.Sprintf("[%d-%d]", conf.ADC1_L_volume, conf.ADC1_R_volume)
                            disp.oledDraw_string(70, menu.screens[menu.active_screen].items[i].y,
                                 ivol, NORMAL_SIZE, WHITE)
                        case I_ADC_RCA2VOL:
                            ivol := fmt.Sprintf("[%d-%d]", conf.ADC2_L_volume, conf.ADC2_R_volume)
                            disp.oledDraw_string(70, menu.screens[menu.active_screen].items[i].y,
                                 ivol, NORMAL_SIZE, WHITE)
                        case I_ADC_SPDIFVOL:
                            ivol := fmt.Sprintf("[%d-%d]", conf.ADC3_L_volume, conf.ADC3_R_volume)
                            disp.oledDraw_string(70, menu.screens[menu.active_screen].items[i].y,
                                 ivol, NORMAL_SIZE, WHITE)

                        case I_RCA1_IN, I_RCA2_IN, I_SPDIF_IN, I_STREAM_IN, I_UAC_IN:
                            if menu.screens[menu.active_screen].items[i].itype == conf.IN_active {
                                disp.oledDraw_bitmap(100, menu.screens[menu.active_screen].items[i].y, 8, 8, &left_arrow, WHITE)
                            }

                        case I_LINE_OUT, I_CDSP_OUT, I_LDSP21_OUT, I_LDSP31_OUT, I_LDSP41_OUT:
                            if menu.screens[menu.active_screen].items[i].itype == conf.OUT_active {
                                disp.oledDraw_bitmap(100, menu.screens[menu.active_screen].items[i].y, 8, 8, &left_arrow, WHITE)
                            }

                        case I_AMP_MUTE:
                            if conf.AMP_Mute_Timeout == menu.screens[menu.active_screen].items[i].iarg1 {
                                disp.oledDraw_bitmap(100, menu.screens[menu.active_screen].items[i].y, 8, 8, &left_arrow, WHITE)
                            }

                        case I_WIFI_AP:
                            if conf.WIFI_AP_Mode {
                                disp.oledDraw_bitmap(100, menu.screens[menu.active_screen].items[i].y, 8, 8, &left_arrow, WHITE)
                            }

                        case I_SIGNAL_METR:
                            if conf.Signal_Meter {
                                disp.oledDraw_bitmap(100, menu.screens[menu.active_screen].items[i].y, 8, 8, &left_arrow, WHITE)
                            }

                        case I_USB_ETH:
                            if USB_ETH_Mode {
                                disp.oledDraw_bitmap(100, menu.screens[menu.active_screen].items[i].y, 8, 8, &left_arrow, WHITE)
                            }
                        case I_CDSP_SAVE:
                            if conf.CDSP_Save {
                                disp.oledDraw_bitmap(100, menu.screens[menu.active_screen].items[i].y, 8, 8, &left_arrow, WHITE)
                            }

                    }
                    if i == menu.cursor_pos {
                        disp.oledRect_invert(0, menu.screens[menu.active_screen].items[i].y-1, disp.width, 9)
                    }
                    its++
                }
            }
        case CTL_MODE:
            switch menu.ctl_id {
                case I_SWVOL:
                    dac3VolGet(1)
                case I_CVOL:
                    dac3VolGet(0)
                case I_FRONTVOL:
                    dac1VolGet()
                case I_SIDEVOL:
                    dac2VolGet()
                case I_DAC4VOL:
                    dac4VolGet()
                case I_ADC_RCA1VOL:
                    adc1VolGet()
                case I_ADC_RCA2VOL:
                    adc2VolGet()
                case I_ADC_SPDIFVOL:
                    adc3VolGet()

                case I_CDSP_CFILE:
                    showCDSPConfigs()
            }
    }
    disp.oledDisplay()
}

func MenuSwitch(dir bool) {     // true = CW, false = CCW
    OLED_time = conf.OLED_Timeout
    if OLED_off { return }
    if OLED_LowContrast {
        OLED_LowContrast = false
        disp.oledSet_contrast(conf.OLED_Contrast)
    }
    if ackMsg { ackMsg = false }

    switch menu.mode {
        case NAV_MODE:
            n_items := len(menu.screens[menu.active_screen].items) - 1
            if dir {
                menu.cursor_pos++
                if menu.cursor_pos > n_items { menu.cursor_pos = 0 }
            } else {
                menu.cursor_pos--
                if menu.cursor_pos < 0 { menu.cursor_pos = n_items }
            }
        case CTL_MODE:
            switch menu.ctl_id {
                case I_SWVOL:
                    dac3VolSet(dir, 1)
                case I_CVOL:
                    dac3VolSet(dir, 0)
                case I_FRONTVOL:
                    dac1VolSet(dir)
                case I_SIDEVOL:
                    dac2VolSet(dir)
                case I_DAC4VOL:
                    dac4VolSet(dir)
                case I_ADC_RCA1VOL:
                    adc1VolSet(dir)
                case I_ADC_RCA2VOL:
                    adc2VolSet(dir)
                case I_ADC_SPDIFVOL:
                    adc3VolSet(dir)

                case I_CDSP_CFILE:
                    n_items := len(cdsp_configs) - 1
                    if dir {
                        menu.cursor_pos++
                        if menu.cursor_pos > n_items { menu.cursor_pos = 0 }
                    } else {
                        menu.cursor_pos--
                        if menu.cursor_pos < 0 { menu.cursor_pos = n_items }
                    }

            }
        case MAIN_MODE:
            if conf.OUT_active == I_CDSP_OUT {
                vol, err := cdspGetVolume()
                if err { return }
                if dir {
                    if vol > 0 {
                        vol = 0.0
                    } else {
                        vol += 1.0
                    }
                } else {
                    if vol < -50.0 {
                        vol = -50.0
                    } else {
                        vol -= 1.0
                    }
                }
                _ = cdspSetVolume(vol)
            }
    }
}

func MenuClick() {
    OLED_time = conf.OLED_Timeout
    if OLED_off { return }
    if OLED_LowContrast {
        OLED_LowContrast = false
        disp.oledSet_contrast(conf.OLED_Contrast)
    }
    if ackMsg { ackMsg = false }

    switch menu.mode {
        case MAIN_MODE:
            menu.mode = NAV_MODE
        case NAV_MODE:
            switch menu.screens[menu.active_screen].items[menu.cursor_pos].itype {
                case I_SCREEN:
                    menu.prev_screen = menu.active_screen
                    menu.active_screen = menu.screens[menu.active_screen].items[menu.cursor_pos].iarg1
                    menu.cursor_pos = 0
                case I_BACK:
                    if menu.screens[menu.active_screen].items[menu.cursor_pos].iarg1 == MAIN_MODE {
                        menu.mode = MAIN_MODE
                        menu.cursor_pos = 0
                        menu.active_screen = 0
                    } else {
                        menu.active_screen = menu.prev_screen
                        menu.cursor_pos = 0
                    }
                case I_SWVOL:
                    menu.mode = CTL_MODE
                    menu.ctl_id = I_SWVOL
                case I_CVOL:
                    menu.mode = CTL_MODE
                    menu.ctl_id = I_CVOL
                case I_FRONTVOL:
                    menu.mode = CTL_MODE
                    menu.ctl_id = I_FRONTVOL
                case I_SIDEVOL:
                    menu.mode = CTL_MODE
                    menu.ctl_id = I_SIDEVOL
                case I_DAC4VOL:
                    menu.mode = CTL_MODE
                    menu.ctl_id = I_DAC4VOL
                case I_ADC_RCA1VOL:
                    menu.mode = CTL_MODE
                    menu.ctl_id = I_ADC_RCA1VOL
                case I_ADC_RCA2VOL:
                    menu.mode = CTL_MODE
                    menu.ctl_id = I_ADC_RCA2VOL
                case I_ADC_SPDIFVOL:
                    menu.mode = CTL_MODE
                    menu.ctl_id = I_ADC_SPDIFVOL

                // CDSP
                case I_CDSP_CFILE:
                    menu.mode = CTL_MODE
                    menu.ctl_id = I_CDSP_CFILE
                    menu.cursor_pos = 0
                    cdsp_configs, _ = cdspGetConfigs()
                case I_CDSP_CFG_SAVE:
                    if conf.OUT_active == I_CDSP_OUT {
                        cfgname,err := cdspSaveConfig()
                        if err {
                            msgDisplay(10, 37, "CDSP CFG SAVE ERR", 2)
                            ackMsg = true
                        } else {
                            conf.CDSP_Config = cfgname
                            msgDisplay(5, 37, cfgname + " saved OK", 1)
                            ackMsg = true
                        }
                    }
                case I_CDSP_SAVE:
                    if conf.CDSP_Save {
                        conf.CDSP_Save = false
                    } else {
                        conf.CDSP_Save = true
                    }
                case I_CDSP_SIGHUP:
                    if conf.OUT_active == I_CDSP_OUT {
                        err := cdspReload()
                        if err {
                            msgDisplay(10, 37, "CDSP RELOAD ERR", 2)
                            ackMsg = true
                        } else {
                            msgDisplay(10, 37, "CDSP RELOAD OK", 1)
                            ackMsg = true
                        }
                    }
                case I_CDSP_STAT:
                    if conf.OUT_active == I_CDSP_OUT {
                        clsmp, err := cdspGetClippedSamples()
                        if err {
                            msgDisplay(10, 37, "CDSP CONNECT ERR", 2)
                            ackMsg = true
                            return
                        }
                        disp.oledClear(BLACK)
                        cls := fmt.Sprintf("SClipped: %d", clsmp)
                        disp.oledDraw_string(1, 1, cls, NORMAL_SIZE, WHITE)

                        pload,err := cdspGetPLoad()
                        if err {
                            pload = 0.0
                        }
                        pload_s := fmt.Sprintf("DSPLoad: %.1f", pload)
                        disp.oledDraw_string(1, 10, pload_s, NORMAL_SIZE, WHITE)

                        cdspste,err := cdspGetState()
                        if err {
                            cdspste = "--"
                        }
                        dspst_s := fmt.Sprintf("State: %s", cdspste)
                        disp.oledDraw_string(1, 21, dspst_s, NORMAL_SIZE, WHITE)

                        cdspsr,err := cdspGetStopReason()
                        if err {
                            cdspsr = "--"
                        }
                        dspsr_s := fmt.Sprintf("SReason: %s", cdspsr)
                        disp.oledDraw_string(1, 31, dspsr_s, NORMAL_SIZE, WHITE)

                        disp.oledDisplay()
                        ackMsg = true
                    }

                // IN sel
                case I_RCA1_IN:
                    adcStop(conf.IN_active)
                    conf.IN_active = I_RCA1_IN
                    adcStart(conf.IN_active)
                case I_RCA2_IN:
                    adcStop(conf.IN_active)
                    conf.IN_active = I_RCA2_IN
                    adcStart(conf.IN_active)
                case I_SPDIF_IN:
                    adcStop(conf.IN_active)
                    conf.IN_active = I_SPDIF_IN
                    adcStart(conf.IN_active)
                case I_STREAM_IN:
                    adcStop(conf.IN_active)
                    conf.IN_active = I_STREAM_IN
                case I_UAC_IN:
                    adcStop(conf.IN_active)
                    conf.IN_active = I_UAC_IN
                    adcStart(conf.IN_active)

                // OUT sel
                case I_LINE_OUT:
                    conf.OUT_active = I_LINE_OUT
                    saveConfig("config.json")
                    runCmd("./sysconfig.sh alsa_update_cfg");
                case I_CDSP_OUT:
                    conf.OUT_active = I_CDSP_OUT
                    saveConfig("config.json")
                    runCmd("./sysconfig.sh alsa_update_cfg");
                    time.Sleep(2 * time.Second)
                    cfgLoadCDSP()
                case I_LDSP21_OUT:
                    conf.OUT_active = I_LDSP21_OUT
                    saveConfig("config.json")
                    runCmd("./sysconfig.sh alsa_update_cfg");
                case I_LDSP31_OUT:
                    conf.OUT_active = I_LDSP31_OUT
                    saveConfig("config.json")
                    runCmd("./sysconfig.sh alsa_update_cfg");
                case I_LDSP41_OUT:
                    conf.OUT_active = I_LDSP41_OUT
                    saveConfig("config.json")
                    runCmd("./sysconfig.sh alsa_update_cfg");

                case I_AMP_MUTE:
                    conf.AMP_Mute_Timeout = menu.screens[menu.active_screen].items[menu.cursor_pos].iarg1

                case I_WIFI_AP:
                    if conf.WIFI_AP_Mode {
                        conf.WIFI_AP_Mode = false
                        runCmd("./sysconfig.sh wifiap_stop");
                    } else {
                        conf.WIFI_AP_Mode = true
                        runCmd("./sysconfig.sh wifiap_start");
                    }
                case I_SIGNAL_METR:
                    if conf.Signal_Meter {
                        conf.Signal_Meter = false
                        sMeter = false
                    } else {
                        conf.Signal_Meter = true
                    }
                    saveConfig("config.json")
                case I_USB_ETH:
                    if USB_ETH_Mode {
                        USB_ETH_Mode = false
                        runCmd("./sysconfig.sh usbeth_stop");
                    } else {
                        USB_ETH_Mode = true
                        runCmd("./sysconfig.sh usbeth_start");
                    }

                case I_REBOOT:
                    ctl_ch <- REBOOT
                case I_HALT:
                    ctl_ch <- HALT

            }
        case CTL_MODE:
            switch menu.ctl_id {
                case I_CDSP_CFILE:
                    conf.CDSP_Config = cdsp_configs[menu.cursor_pos]
                    if conf.OUT_active == I_CDSP_OUT {
                        cfgLoadCDSP()
                    }
            }
            menu.mode = NAV_MODE
    }
}

func cfgLoadCDSP() {
    log.Println("Loading CDSP config:", conf.CDSP_Config)
    err := cdspLoadConfig(conf.CDSP_Config)
    if err {
        disp.oledClear(BLACK)
        msgDisplay(0, 35, "CDSP CFG LOAD ERROR", 2)
        log.Println("CDSP config load error")
    }
}

// Left/Right values extract from tinymix output
func lrExtract(s string) (int, int) {
    tmp := strings.Split(s, " ")
    ls := tmp[0][:len(tmp[0])-1]
    rs := tmp[1]
    l, err := strconv.Atoi(ls)
    if err != nil {
        log.Println("tinymix parse value error")
    }
    r, err := strconv.Atoi(rs)
    if err != nil {
        log.Println("tinymix parse value error")
    }
    return l,r
}

func tinymixSet2(id int, v1 int, v2 int) {
    cmd := fmt.Sprintf("tinymix -D%d set %d %d %d", ALSA_CARD, id, v1, v2)
    err := runCmd(cmd)
    if err != 0 {
        log.Println("Unable to set tinymix value")
    }
}

func tinymixGet2(id int) (int, int) {
    cmd := fmt.Sprintf("tinymix -D%d get %d", ALSA_CARD, id)
    stdout := runCmdOut(cmd)
    if len(stdout) < 3 {
        log.Println("Unable to get tinymix value")
        return 0,0
    }
    return lrExtract(stdout)
}

func dac1VolGet() {
    l,r := tinymixGet2(DAC1_VOL_CTL)
    conf.FR_volume = r
    conf.FL_volume = l
    disp.oledDraw_string(35, 15, strconv.Itoa(r), DOUBLE_SIZE ,WHITE)
    disp.oledDraw_string(30, 35, "FR/FL Volume", NORMAL_SIZE ,WHITE)
}

func dac2VolGet() {
    l,r := tinymixGet2(DAC2_VOL_CTL)
    conf.SR_volume = r
    conf.SL_volume = l
    disp.oledDraw_string(35, 15, strconv.Itoa(r), DOUBLE_SIZE ,WHITE)
    disp.oledDraw_string(30, 35, "SR/SL Volume", NORMAL_SIZE ,WHITE)
}

func dac3VolGet(ch int) {
    l,r := tinymixGet2(DAC3_VOL_CTL)
    conf.Sw_volume = r
    conf.Center_volume = l
    v:=r
    if ch == 0 {
        v=l
    }
    disp.oledDraw_string(35, 15, strconv.Itoa(v), DOUBLE_SIZE ,WHITE)
    if ch == 0 {
        disp.oledDraw_string(30, 35, "CE Volume", NORMAL_SIZE ,WHITE)
    } else {
        disp.oledDraw_string(30, 35, "SW Volume", NORMAL_SIZE ,WHITE)
    }
}

func dac4VolGet() {
    l,r := tinymixGet2(DAC4_VOL_CTL)
    conf.DR_volume = r
    conf.DL_volume = l
    disp.oledDraw_string(35, 15, strconv.Itoa(r), DOUBLE_SIZE ,WHITE)
    disp.oledDraw_string(30, 35, "DAC4 Volume", NORMAL_SIZE ,WHITE)
}

func dac1VolSet(dir bool) {
    vstep := conf.Vol_step
    diff := time.Now().Sub(evtTime)
    if diff.Milliseconds() < 50 { vstep += 15 }

    if dir {
        conf.FR_volume += vstep
        conf.FL_volume += vstep
        if conf.FR_volume > 255 { conf.FR_volume = 255 }
        if conf.FL_volume > 255 { conf.FL_volume = 255 }
    } else {
        conf.FR_volume -= vstep
        conf.FL_volume -= vstep
        if conf.FR_volume < 0 { conf.FR_volume = 0 }
        if conf.FL_volume < 0 { conf.FL_volume = 0 }
    }
    tinymixSet2(DAC1_VOL_CTL, conf.FL_volume, conf.FR_volume)
    evtTime = time.Now()
}

func dac2VolSet(dir bool) {
    vstep := conf.Vol_step
    diff := time.Now().Sub(evtTime)
    if diff.Milliseconds() < 50 { vstep += 15 }

    if dir {
        conf.SR_volume += vstep
        conf.SL_volume += vstep
        if conf.SR_volume > 255 { conf.SR_volume = 255 }
        if conf.SL_volume > 255 { conf.SL_volume = 255 }
    } else {
        conf.SR_volume -= vstep
        conf.SL_volume -= vstep
        if conf.SR_volume < 0 { conf.SR_volume = 0 }
        if conf.SL_volume < 0 { conf.SL_volume = 0 }
    }
    tinymixSet2(DAC2_VOL_CTL, conf.SL_volume, conf.SR_volume)
    evtTime = time.Now()
}

func dac3VolSet(dir bool, ch int) {
    vstep := conf.Vol_step
    diff := time.Now().Sub(evtTime)
    if diff.Milliseconds() < 50 { vstep += 15 }

    if dir {
        if ch == 0 {
            conf.Center_volume += vstep
            if conf.Center_volume > 255 { conf.Center_volume = 255 }
        } else {
            conf.Sw_volume += vstep
            if conf.Sw_volume > 255 { conf.Sw_volume = 255 }
        }
    } else {
        if ch == 0 {
            conf.Center_volume -= vstep
            if conf.Center_volume < 0 { conf.Center_volume = 0 }
        } else {
            conf.Sw_volume -= vstep
            if conf.Sw_volume < 0 { conf.Sw_volume = 0 }
        }
    }
    tinymixSet2(DAC3_VOL_CTL, conf.Center_volume, conf.Sw_volume)
    evtTime = time.Now()
}

func dac4VolSet(dir bool) {
    vstep := conf.Vol_step
    diff := time.Now().Sub(evtTime)
    if diff.Milliseconds() < 50 { vstep += 15 }

    if dir {
        conf.DR_volume += vstep
        conf.DL_volume += vstep
        if conf.DR_volume > 255 { conf.DR_volume = 255 }
        if conf.DL_volume > 255 { conf.DL_volume = 255 }
    } else {
        conf.DR_volume -= vstep
        conf.DL_volume -= vstep
        if conf.DR_volume < 0 { conf.DR_volume = 0 }
        if conf.DL_volume < 0 { conf.DL_volume = 0 }
    }
    tinymixSet2(DAC4_VOL_CTL, conf.DL_volume, conf.DR_volume)
    evtTime = time.Now()
}

//ADC vol
func adc1VolGet() {
    l,r := tinymixGet2(ADC1_VOL_CTL)
    conf.ADC1_R_volume = r
    conf.ADC1_L_volume = l
    disp.oledDraw_string(35, 15, strconv.Itoa(r), DOUBLE_SIZE ,WHITE)
    disp.oledDraw_string(12, 35, "RCA1(ADC1) Volume", NORMAL_SIZE ,WHITE)
}

func adc1VolSet(dir bool) {
    vstep := conf.Vol_step
    diff := time.Now().Sub(evtTime)
    if diff.Milliseconds() < 50 { vstep += 15 }

    if dir {
        conf.ADC1_R_volume += vstep
        conf.ADC1_L_volume += vstep
        if conf.ADC1_R_volume > 176 { conf.ADC1_R_volume = 176 }
        if conf.ADC1_L_volume > 176 { conf.ADC1_L_volume = 176 }
    } else {
        conf.ADC1_R_volume -= vstep
        conf.ADC1_L_volume -= vstep
        if conf.ADC1_R_volume < 0 { conf.ADC1_R_volume = 0 }
        if conf.ADC1_L_volume < 0 { conf.ADC1_L_volume = 0 }
    }
    tinymixSet2(ADC1_VOL_CTL, conf.ADC1_L_volume, conf.ADC1_R_volume)
    evtTime = time.Now()
}

func adc2VolGet() {
    l,r := tinymixGet2(ADC2_VOL_CTL)
    conf.ADC2_R_volume = r
    conf.ADC2_L_volume = l
    disp.oledDraw_string(35, 15, strconv.Itoa(r), DOUBLE_SIZE ,WHITE)
    disp.oledDraw_string(12, 35, "RCA2(ADC2) Volume", NORMAL_SIZE ,WHITE)
}

func adc2VolSet(dir bool) {
    vstep := conf.Vol_step
    diff := time.Now().Sub(evtTime)
    if diff.Milliseconds() < 50 { vstep += 15 }

    if dir {
        conf.ADC2_R_volume += vstep
        conf.ADC2_L_volume += vstep
        if conf.ADC2_R_volume > 176 { conf.ADC2_R_volume = 176 }
        if conf.ADC2_L_volume > 176 { conf.ADC2_L_volume = 176 }
    } else {
        conf.ADC2_R_volume -= vstep
        conf.ADC2_L_volume -= vstep
        if conf.ADC2_R_volume < 0 { conf.ADC2_R_volume = 0 }
        if conf.ADC2_L_volume < 0 { conf.ADC2_L_volume = 0 }
    }
    tinymixSet2(ADC2_VOL_CTL, conf.ADC2_L_volume, conf.ADC2_R_volume)
    evtTime = time.Now()
}

func adc3VolGet() {
    l,r := tinymixGet2(ADC3_VOL_CTL)
    conf.ADC3_R_volume = r
    conf.ADC3_L_volume = l
    disp.oledDraw_string(35, 15, strconv.Itoa(r), DOUBLE_SIZE ,WHITE)
    disp.oledDraw_string(10, 35, "SPDIF(ADC3) Volume", NORMAL_SIZE ,WHITE)
}

func adc3VolSet(dir bool) {
    vstep := conf.Vol_step
    diff := time.Now().Sub(evtTime)
    if diff.Milliseconds() < 50 { vstep += 15 }

    if dir {
        conf.ADC3_R_volume += vstep
        conf.ADC3_L_volume += vstep
        if conf.ADC3_R_volume > 176 { conf.ADC3_R_volume = 176 }
        if conf.ADC3_L_volume > 176 { conf.ADC3_L_volume = 176 }
    } else {
        conf.ADC3_R_volume -= vstep
        conf.ADC3_L_volume -= vstep
        if conf.ADC3_R_volume < 0 { conf.ADC3_R_volume = 0 }
        if conf.ADC3_L_volume < 0 { conf.ADC3_L_volume = 0 }
    }
    tinymixSet2(ADC3_VOL_CTL, conf.ADC3_L_volume, conf.ADC3_R_volume)
    evtTime = time.Now()
}

func adcStart(src int) {
    var cmd string
    switch src {
        case I_RCA1_IN:
            cmd = "./rca1_src.sh listen"
        case I_RCA2_IN:
            cmd = "./rca2_src.sh listen"
        case I_SPDIF_IN:
            cmd = "./spdif_src.sh listen"
        case I_UAC_IN:
            cmd = "./uac_src.sh listen"
        case I_STREAM_IN:
            cmd = "./stop_src.sh"
    }
    err := runCmd(cmd)
    if err != 0 {
        log.Println("Unable to start IN src:", src)
    }
}

func adcStop(src int) {
    var cmd string
    switch src {
        case I_RCA1_IN:
            cmd = "./rca1_src.sh stop"
        case I_RCA2_IN:
            cmd = "./rca2_src.sh stop"
        case I_SPDIF_IN:
            cmd = "./spdif_src.sh stop"
        case I_UAC_IN:
            cmd = "./uac_src.sh stop"
        case I_STREAM_IN:
            cmd = "./stop_src.sh"
    }
    err := runCmd(cmd)
    if err != 0 {
        log.Println("Unable to stop IN src:", src)
    }
}

func showCDSPConfigs() {
    var offset int = 0

    n_cfg := len(cdsp_configs)
    if menu.cursor_pos > 5 {
            offset = 6 * (menu.cursor_pos/6)
    }
    its,y := 0,0
    for i := offset; i < n_cfg; i++ {
        if its <= 5 {  //max onscreen items
            disp.oledDraw_string(0, int16(y)*8, cdsp_configs[i], NORMAL_SIZE ,WHITE)
            if i == menu.cursor_pos {
                    disp.oledRect_invert(0, int16(y)*8, disp.width, 8)
            }
            y++
            its++
        }
    }
}
