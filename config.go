package main

import ("log"
        "os"
        "io"
        "encoding/json"
    )

type Config struct {
    Sw_volume     int
    Center_volume int
    FR_volume     int
    FL_volume     int
    SR_volume     int
    SL_volume     int
    //DAC4 direct output
    DR_volume     int
    DL_volume     int

    //ADC volumes
    ADC1_L_volume  int
    ADC1_R_volume  int
    ADC2_L_volume  int
    ADC2_R_volume  int
    ADC3_L_volume  int
    ADC3_R_volume  int

    //Active input
    IN_active     int
    OUT_active    int

    //Rotary switch step for vol reg
    Vol_step      int

    OLED_Contrast uint8
    // OLED timeout in sec
    OLED_Timeout  int

    // Put amp to idle mode after n-sec of inactivity
    AMP_Mute_Timeout int

    CDSP_Config string
    CDSP_Save bool

    WIFI_AP_Mode bool
}

var conf Config

func loadConfig(name string) {
    fn, err := os.Open(name)
    if err != nil {
        log.Fatal("Unable to open config: ", err)
    }
    defer fn.Close()
    data, err := io.ReadAll(fn)
    if err != nil {
        log.Fatal("Unable to load config: ", err)
    }
    err = json.Unmarshal([]byte(data), &conf)
    if err != nil {
        log.Fatal("Unable to parse config: ", err)
    }

}

func saveConfig(name string) {
    fn, err := os.Create(name)
    if err != nil {
        log.Fatal("Unable to save config: ", err)
    }
    defer fn.Close()
    b, err := json.MarshalIndent(conf, "", "\t")
    if err != nil {
        log.Fatal("Unable to encode config: ", err)
    }
    fn.Write(b)
}

func applyConfig() {
    tinymixSet2(DAC1_VOL_CTL, conf.FL_volume, conf.FR_volume)
    tinymixSet2(DAC2_VOL_CTL, conf.SL_volume, conf.SR_volume)
    tinymixSet2(DAC3_VOL_CTL, conf.Center_volume, conf.Sw_volume)
    tinymixSet2(DAC4_VOL_CTL, conf.DL_volume, conf.DR_volume)
    tinymixSet2(ADC1_VOL_CTL, conf.ADC1_L_volume, conf.ADC1_R_volume)
    tinymixSet2(ADC2_VOL_CTL, conf.ADC2_L_volume, conf.ADC2_R_volume)
    tinymixSet2(ADC3_VOL_CTL, conf.ADC3_L_volume, conf.ADC3_R_volume)
    adcStop(I_STREAM_IN) //kill every input listener if any
    adcStart(conf.IN_active)
}

