package main

import (
    "os"
    "log"
    "strconv"
    "github.com/gorilla/websocket"
    "github.com/buger/jsonparser"
)

const CDSP_ADDR = "ws://127.0.0.1:5050"
const CDSP_CFG_DIR = "/opt/dsp/configs"

func GetFloatArray(data []byte, keys ...string) ([]float64, error) {
    var floats []float64
    _, err := jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
        if err != nil {
            return
        }
        num, err := jsonparser.ParseFloat(value)
        if err != nil {
            return
        }
        floats = append(floats, num)
    }, keys...)
    if err != nil {
        return nil, err
    }
    return floats, nil
}

func cdspGetPlaybackSignalRms() (reason []float64, error bool) {
    ws, _, err := websocket.DefaultDialer.Dial(CDSP_ADDR, nil)
    if err != nil {
        return nil, true
    }
    defer ws.Close()

    msgs := "\"GetPlaybackSignalRms\""
    if err := ws.WriteMessage(websocket.TextMessage, []byte(msgs)); err != nil {
        return nil, true
    }
    _, msg, err := ws.ReadMessage()
    if err != nil {
         return nil, true
    }

    res, _ := jsonparser.GetString(msg, "GetPlaybackSignalRms", "result")
    if res != "Ok" {
        return nil, true
    }

    resn, err := GetFloatArray(msg, "GetPlaybackSignalRms", "value")
    if err != nil {
         return nil, true
    }

    return resn, false
}

func cdspGetVolume() (volume float64, error bool) {
    ws, _, err := websocket.DefaultDialer.Dial(CDSP_ADDR, nil)
    if err != nil {
        return 0, true
    }
    defer ws.Close()

    msgs := "\"GetVolume\""
    if err := ws.WriteMessage(websocket.TextMessage, []byte(msgs)); err != nil {
        return 1, true
    }
    _, msg, err := ws.ReadMessage()
    if err != nil {
         return 1, true
    }

    res, _ := jsonparser.GetString(msg, "GetVolume", "result")
    if res != "Ok" {
        return 2, true
    }
    vol, _ := jsonparser.GetFloat(msg, "GetVolume", "value")

    return vol, false
}

func cdspSetVolume(volume float64) (error bool) {
    ws, _, err := websocket.DefaultDialer.Dial(CDSP_ADDR, nil)
    if err != nil {
        return true
    }
    defer ws.Close()

    msgs := "{\"SetVolume\":" + strconv.FormatFloat(volume, 'f', 2, 64) + "}"
    if err := ws.WriteMessage(websocket.TextMessage, []byte(msgs)); err != nil {
        return true
    }
    _, msg, err := ws.ReadMessage()
    if err != nil {
         return true
    }
    res, _ := jsonparser.GetString(msg, "SetVolume", "result")
    if res != "Ok" {
        return true
    }

    return false
}

func cdspGetConfig() (cfg string, error bool) {
    ws, _, err := websocket.DefaultDialer.Dial(CDSP_ADDR, nil)
    if err != nil {
        return "", true
    }
    defer ws.Close()

    msgs := "\"GetConfig\""
    if err := ws.WriteMessage(websocket.TextMessage, []byte(msgs)); err != nil {
        return "", true
    }
    _, msg, err := ws.ReadMessage()
    if err != nil {
         return "", true
    }

    res, _ := jsonparser.GetString(msg, "GetConfig", "result")
    if res != "Ok" {
        return "", true
    }

    cfgy, _ := jsonparser.GetString(msg, "GetConfig", "value")
    return cfgy, false
}

func cdspGetConfigPath() (cpath string, error bool) {
    ws, _, err := websocket.DefaultDialer.Dial(CDSP_ADDR, nil)
    if err != nil {
        return "", true
    }
    defer ws.Close()

    msgs := "\"GetConfigFilePath\""
    if err := ws.WriteMessage(websocket.TextMessage, []byte(msgs)); err != nil {
        return "", true
    }
    _, msg, err := ws.ReadMessage()
    if err != nil {
         return "", true
    }

    res, _ := jsonparser.GetString(msg, "GetConfigFilePath", "result")
    if res != "Ok" {
        return "", true
    }

    cp, _ := jsonparser.GetString(msg, "GetConfigFilePath", "value")
    return cp, false
}

func cdspGetConfigTitle() (title string, error bool) {
    ws, _, err := websocket.DefaultDialer.Dial(CDSP_ADDR, nil)
    if err != nil {
        return "", true
    }
    defer ws.Close()

    msgs := "\"GetConfigTitle\""
    if err := ws.WriteMessage(websocket.TextMessage, []byte(msgs)); err != nil {
        return "", true
    }
    _, msg, err := ws.ReadMessage()
    if err != nil {
         return "", true
    }

    res, _ := jsonparser.GetString(msg, "GetConfigTitle", "result")
    if res != "Ok" {
        return "", true
    }

    ct, _ := jsonparser.GetString(msg, "GetConfigTitle", "value")
    return ct, false
}

func cdspGetClippedSamples() (clsm int, error bool) {
    ws, _, err := websocket.DefaultDialer.Dial(CDSP_ADDR, nil)
    if err != nil {
        return 0, true
    }
    defer ws.Close()

    msgs := "\"GetClippedSamples\""
    if err := ws.WriteMessage(websocket.TextMessage, []byte(msgs)); err != nil {
        return 0, true
    }
    _, msg, err := ws.ReadMessage()
    if err != nil {
         return 0, true
    }

    res, _ := jsonparser.GetString(msg, "GetClippedSamples", "result")
    if res != "Ok" {
        return 0, true
    }

    cs, _ := jsonparser.GetInt(msg, "GetClippedSamples", "value")
    if err != nil { return 0, true }
    return int(cs), false
}

func cdspGetPLoad() (load float64, error bool) {
    ws, _, err := websocket.DefaultDialer.Dial(CDSP_ADDR, nil)
    if err != nil {
        return 0.0, true
    }
    defer ws.Close()

    msgs := "\"GetProcessingLoad\""
    if err := ws.WriteMessage(websocket.TextMessage, []byte(msgs)); err != nil {
        return 0.0, true
    }
    _, msg, err := ws.ReadMessage()
    if err != nil {
         return 0.0, true
    }

    res, _ := jsonparser.GetString(msg, "GetProcessingLoad", "result")
    if res != "Ok" {
        return 0.0, true
    }

    pl, _ := jsonparser.GetFloat(msg, "GetProcessingLoad", "value")
    return pl, false
}

func cdspGetState() (state string, error bool) {
    ws, _, err := websocket.DefaultDialer.Dial(CDSP_ADDR, nil)
    if err != nil {
        return "", true
    }
    defer ws.Close()

    msgs := "\"GetState\""
    if err := ws.WriteMessage(websocket.TextMessage, []byte(msgs)); err != nil {
        return "", true
    }
    _, msg, err := ws.ReadMessage()
    if err != nil {
         return "", true
    }

    res, _ := jsonparser.GetString(msg, "GetState", "result")
    if res != "Ok" {
        return "", true
    }

    st, _ := jsonparser.GetString(msg, "GetState", "value")
    return st, false
}

func cdspGetStopReason() (reason string, error bool) {
    ws, _, err := websocket.DefaultDialer.Dial(CDSP_ADDR, nil)
    if err != nil {
        return "", true
    }
    defer ws.Close()

    msgs := "\"GetStopReason\""
    if err := ws.WriteMessage(websocket.TextMessage, []byte(msgs)); err != nil {
        return "", true
    }
    _, msg, err := ws.ReadMessage()
    if err != nil {
         return "", true
    }

    res, _ := jsonparser.GetString(msg, "GetStopReason", "result")
    if res != "Ok" {
        return "", true
    }

    resn, _ := jsonparser.GetString(msg, "GetStopReason", "value")
    return resn, false
}

func cdspReload() (error bool) {
    ws, _, err := websocket.DefaultDialer.Dial(CDSP_ADDR, nil)
    if err != nil {
        return true
    }
    defer ws.Close()

    msgs := "\"Reload\""
    if err := ws.WriteMessage(websocket.TextMessage, []byte(msgs)); err != nil {
        return true
    }
    _, msg, err := ws.ReadMessage()
    if err != nil {
         return true
    }

    res, _ := jsonparser.GetString(msg, "Reload", "result")
    if res != "Ok" {
        return true
    }

    return false
}

func cdspSaveConfig() (conf string, error bool) {
    cfg, err := cdspGetConfig()
    if err { return "",true }

    name, err := cdspGetConfigTitle()
    if err { return "",true }

    cfg_path := CDSP_CFG_DIR + "/" + name + ".yml"
    log.Println("Saving CDSP config:", cfg_path)

    fn, errio := os.Create(cfg_path)
    if errio != nil {
        log.Println("Unable to save CDSP config: ", errio)
    }
    defer fn.Close()

    _,errio = fn.Write([]byte(cfg))
    if errio != nil {
        log.Println("Unable to save CDSP config: ", errio)
        return "",true
    }
    return name + ".yml",false
}

func cdspLoadConfig(conf string) (error bool) {
    ws, _, err := websocket.DefaultDialer.Dial(CDSP_ADDR, nil)
    if err != nil {
        log.Println("CDSP WS Error")
        return true
    }
    defer ws.Close()

    f, err := os.Open(CDSP_CFG_DIR + "/" + conf)
    if err != nil {
        log.Println("Unable to open CDSP config:", conf)
        return true
    }
    defer f.Close()

    cfg := make([]byte, 32000)
    nb, err := f.Read(cfg)
    if err != nil {
        log.Println("Unable to read CDSP config:", conf)
        return true
    }

    cfg_str := string(cfg[:nb])
    msgs := `{"SetConfig":` + strconv.Quote(cfg_str) + `}`
    if err := ws.WriteMessage(websocket.TextMessage, []byte(msgs)); err != nil {
        return true
    }
    _, msg, err := ws.ReadMessage()
    if err != nil {
         return true
    }
    res, _ := jsonparser.GetString(msg, "SetConfig", "result")
    if res != "Ok" {
        log.Println("CDSP config apply failed")
        return true
    }
    log.Println("CDSP config apply OK")
    return false
}

func cdspGetConfigs() ( configs []string, error bool) {
    var cfgs []string

    f, err := os.Open(CDSP_CFG_DIR)
    if err != nil {
        log.Println("Unable to open CDSP configs dir:", err)
        return cfgs,true
    }
    defer f.Close()

    files, err := f.Readdir(0)
    if err != nil {
        log.Println("Unable to read CDSP configs dir:", err)
        return cfgs,true
    }

    for _, v := range files {
        if !v.IsDir() {
            cfgs = append(cfgs,v.Name())
        }
    }
    return cfgs, false
}
