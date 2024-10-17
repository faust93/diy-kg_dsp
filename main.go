package main

import (
    "fmt"
    "log"
    "time"
    "os"
    "os/exec"
    "os/signal"
    "errors"
    "syscall"
    "sync"
    "strings"
    "strconv"
    "github.com/warthog618/go-gpiocdev"
    )

// HW contstants
const (
    // OLED I2c
    OLED_I2C_CH = 0x0
    OLED_I2C_ID = 0x3c

    // RotarySwitch GPIOs
    RSW_PIN_CLK_CH = "gpiochip2"
    RSW_PIN_CLK = 0
    RSW_PIN_DT_CH = "gpiochip2"
    RSW_PIN_DT  = 1
    RSW_PIN_SW_CH = "gpiochip1"
    RSW_PIN_SW  = 22

    // AMP mute
    MUTE_PIN_CH = "gpiochip0"
    MUTE_PIN = 15
    MUTE_TH = 60 // AMP mute threshold sec

    DAC_STATE_F = "/proc/asound/cirruscs42448/pcm0p/sub0/status"
    DAC_P_PARAMS_F = "/proc/asound/cirruscs42448/pcm0p/sub0/hw_params"
    DAC_C_PARAMS_F = "/proc/asound/cirruscs42448/pcm0c/sub0/hw_params"
)

var mutex = &sync.Mutex{}

var disp OLED
var rsw RotarySwitch

var muteGPIO *gpiocdev.Line

var unmuteTime time.Time
var AMP_Mute int

var OLED_off bool
var OLED_time int
var OLED_LowContrast bool = false

// WIFI connection state
var WifiState bool = false

// DAC playback device state
var DacPRate, DacPCh, DacPFmt string

// USB Ethernet gadget (for debug/config purposes)
var USB_ETH_Mode bool = false

// control channel
var ctl_ch = make(chan int, 16)
var StopReason int = STERM
const (
    DISP_UPDATE = 1
    HALT = 0xfd
    REBOOT = 0xfe
    STERM = 0xff
)

func pbHandler(fn int) {
    if fn == 0 {
        MenuClick()
    } else if fn == 1 {
        fmt.Println("Long press")
    } else {
        fmt.Println("Very Long press")
    }
    ctl_ch <- DISP_UPDATE
}

func rHandler(dir bool) {
    if dir {
        MenuSwitch(true)
    } else {
        MenuSwitch(false)
    }
    ctl_ch <- DISP_UPDATE
}

func signalHandler() {
    signChan := make(chan os.Signal, 1)
    signal.Notify(signChan, syscall.SIGINT, syscall.SIGTERM)
    <-signChan
    if StopReason != STERM {
        return
    } else {
        ctl_ch <- STERM
    }
}

func runCmd(cmdStr string) int {
        cmd := exec.Command("sh", "-c", cmdStr)
        err := cmd.Run()
        var exitErr *exec.ExitError
        if errors.As(err, &exitErr) {
            log.Println("Error launching command:", err)
            return exitErr.ExitCode()
        }
        return 0
}

func runCmdOut(cmdStr string) string {
        out, err := exec.Command("sh", "-c", cmdStr).CombinedOutput()
        if err != nil {
            log.Println("Error launching command:", err)
            return ""
        }
        return string(out)
}

// Mute/Unmute amp control thread
func ampControl() {
    f, err := os.Open(DAC_STATE_F)
    if err != nil {
        msgDisplay(1, 30, "HW ERROR: no DAC found", 2)
        log.Println("Sound system error, no DAC found")
        time.Sleep(20 * time.Second)
    }
    defer f.Close()

    buf := make([]byte, 128)
    for {
        for i := range buf { buf[i] = 0 }
        _, err := f.Seek(0, 0)
        _, err = f.Read(buf)
        if err != nil {
            log.Println("Sound system error reading DAC state")
        }

        astate := strings.Contains(string(buf), "RUNNING")
        astate1 := strings.Contains(string(buf), "SETUP")
        astate2 := strings.Contains(string(buf), "closed")

        if astate && AMP_Mute == 0 {
            log.Println("unmute AMP")
            muteAMP(1)  //unmute
            unmuteTime = time.Now()
            ctl_ch <- DISP_UPDATE
        } else if (astate1 || astate2) && AMP_Mute == 1 {

            diff := time.Now().Sub(unmuteTime)
            if diff.Seconds() > float64(MUTE_TH) {
                log.Println("mute AMP")
                muteAMP(0) //mute
                ctl_ch <- DISP_UPDATE
            }
        }
        time.Sleep(1 * time.Second)
    }
}

func muteAMP(state int) {
    AMP_Mute = state
    muteGPIO.SetValue(state)
}

func wifiControl() {
    f, err := os.Open("/proc/net/arp")
    if err != nil {
        log.Fatal("/proc/net/arp not found")
    }
    defer f.Close()

    buf := make([]byte, 256)
    for {
        for i := range buf { buf[i] = 0 }
        _, err := f.Seek(0, 0)
        _, err = f.Read(buf)
        if err != nil {
            log.Println("/proc/net/arp read error")
        }

        arpRec := strings.Contains(string(buf), "wlan0")

        if arpRec && !WifiState {
            log.Println("WIFI connected")
            WifiState = true
            ctl_ch <- DISP_UPDATE
        } else if !arpRec && WifiState {
            log.Println("WIFI disconnected")
            WifiState = false
            ctl_ch <- DISP_UPDATE
        }
        time.Sleep(5 * time.Second)
    }
}

func oledControl() {
    for {
        if OLED_off && OLED_time > 0 {
            disp.oledSet_power(true)
            disp.oledSet_contrast(conf.OLED_Contrast)
            OLED_off = false
            OLED_LowContrast = false
        } else if !OLED_off && OLED_time <= 0 {
            disp.oledSet_power(false)
            OLED_off = true
        }
        if OLED_time > 0 {
            mutex.Lock()
            OLED_time--
            mutex.Unlock()
            if OLED_time < (conf.OLED_Timeout / 2) && !OLED_LowContrast {
                OLED_LowContrast = true
                disp.oledSet_contrast(conf.OLED_Contrast / 3)
            }
        }
        time.Sleep(1 * time.Second)
    }
}

func dacParamsFetch() {
    dacP, err := os.Open(DAC_P_PARAMS_F)
    if err != nil {
        log.Fatal(DAC_P_PARAMS_F, "not found")
    }
    defer dacP.Close()

    buf := make([]byte, 130)
    for {
        for i := range buf { buf[i] = 0 }
        _, err := dacP.Seek(0, 0)
        _, err = dacP.Read(buf)
        if err != nil {
            log.Println(DAC_P_PARAMS_F, "read error")
        }
        var stUpd bool = false

        state := strings.Contains(string(buf), "closed")
        if !state {
            if dfmt, err := getStrKv(string(buf),"format:"); !err {
                if dfmt != DacPFmt { stUpd = true }
                DacPFmt = dfmt
            }
            if dch, err := getStrKv(string(buf),"channels:"); !err {
                if dch != DacPCh { stUpd = true }
                DacPCh = dch
            }
            if drate, err := getStrKv(string(buf),"rate:"); !err {
                if drate != DacPRate { stUpd = true }
                DacPRate = drate
            }
            if stUpd || !OLED_off {
                ctl_ch <- DISP_UPDATE
                stUpd = false
            }
        }

        time.Sleep(5 * time.Second)
    }
}

func getStrKv(str string, key string) (val string, err bool) {
    if ssidx := strings.Index(str, key); ssidx != -1 {
            ssend := strings.IndexByte(str[ssidx:], 0x0a)
            return string(str[ssidx+len(key)+1:ssidx+ssend]), false
    }
    return "", true
}


func msgDisplay(x int, y int, msg string, mtype int) {
    // mtype: 0 plain message, 1 warning, 2 error
    switch mtype {
        case 1:
            disp.oledClear(BLACK)
            disp.oledDraw_bitmap(45, 1, 24, 24, &warning_img, WHITE)
        case 2:
            disp.oledClear(BLACK)
            disp.oledDraw_bitmap(45, 1, 24, 24, &error_img, WHITE)
    }
    disp.oledDraw_string(int16(x), int16(y), msg, NORMAL_SIZE ,WHITE)
    disp.oledDisplay()
}

func preExit() {
    if OLED_off {
        disp.oledSet_power(true)
    }
    disp.oledClear(BLACK)

    if conf.OUT_active == I_CDSP_OUT && conf.CDSP_Save {
        cfgname,err := cdspSaveConfig()
        if err {
            msgDisplay(0, 37, "CDSP CFG SAVE ERR", 2)
            time.Sleep(30 * time.Second)
        } else {
            conf.CDSP_Config = cfgname
        }
    }

    saveConfig("config.json")
}

func mutePipe() {
    msg := make([]byte, 128)
    mute_pipe := "/tmp/mute"
    os.Remove(mute_pipe)
    err := syscall.Mkfifo(mute_pipe, 0666)
    if err != nil {
        log.Fatal("Make mute pipe file error:", err)
    }

    file, err := os.OpenFile(mute_pipe, os.O_CREATE|os.O_RDWR, os.ModeNamedPipe)
    if err != nil {
        log.Fatal("Open mute pipe file error:", err)
    }
    defer file.Close()

    for {
        nb, err := file.Read(msg)
        if err == nil {
            mstat, err := strconv.Atoi(strings.TrimSpace(string(msg[:nb])))
            if err != nil {
                log.Println("Error parsing mute value")
            } else {
                log.Println("Mute pipe received:",mstat)
                muteGPIO.SetValue(mstat)
            }
        } else {
            log.Println("Error reading mute pipe")
        }
    }
}

func main() {
    log.Println("Kirogaz DSP control program v0.1")

    loadConfig("config.json")

    OLED_off = false
    OLED_time = conf.OLED_Timeout
    disp = OLED{oled_model: SH1106, width: W_132, height: H_64}
    disp.oledInit(OLED_I2C_CH, OLED_I2C_ID)
    disp.oledSet_contrast(conf.OLED_Contrast)
    defer disp.Close()

    var err error
    muteGPIO, err = gpiocdev.RequestLine(MUTE_PIN_CH, MUTE_PIN, gpiocdev.AsOutput(0))
    if err != nil {
            msgDisplay(5, 30, "HW ERROR: Mute PIN", 2)
            log.Fatal("AMP GPIO RequestLine returned error: %w", err)
    }
    defer muteGPIO.Close()

    muteAMP(0) //0 mute, 1 unmute

    rsw = RotarySwitch{
        gpio_clk_chip: RSW_PIN_CLK_CH,
        gpio_clk: RSW_PIN_CLK,
        gpio_dt_chip: RSW_PIN_DT_CH,
        gpio_dt: RSW_PIN_DT,
        gpio_sw_chip: RSW_PIN_SW_CH,
        gpio_sw: RSW_PIN_SW,
        longPthr: 500,
        long2Pthr: 2000,
        swLongPress: pbHandler,
        swLong2Press: pbHandler,
        swShortPress: pbHandler,
        rotHandler: rHandler,
    }
    rsw.RotarySwitchInit()
    defer rsw.Close()

    go signalHandler()
    go ampControl()
    go oledControl()
    go wifiControl()
    go dacParamsFetch()
    go mutePipe()

    log.Println("Starting services")
    if err := runCmd("./sysconfig.sh start"); err != 0 {
        log.Println("Some services failed to start")
    }

    applyConfig()

    if conf.OUT_active == I_CDSP_OUT {
        time.Sleep(2 * time.Second)
        cfgLoadCDSP()
    }

    MenuInit()
    MenuDisplay()

    ctl_ch = make(chan int)
    for {
        select {
            case ctl := <-ctl_ch:
                switch ctl {
                    case DISP_UPDATE:
                        if StopReason == STERM {
                            MenuDisplay()
                        }
                    case REBOOT:
                        StopReason = REBOOT
                        preExit()
                        msgDisplay(15, 25, "Rebooting..", 0)
                        runCmd("reboot")
                    case HALT:
                        StopReason = HALT
                        preExit()
                        msgDisplay(15, 25, "Power OFF..", 0)
                        runCmd("halt")
                    case STERM:
                        preExit()
                        log.Println("Stopping services")
                        runCmd("./sysconfig.sh terminate")
                        os.Exit(1)
                }
        }
    }
}