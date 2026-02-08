// Rotary Switch
package main

import (
    "log"
    "time"
    "github.com/warthog618/go-gpiocdev"
    )

type RotarySwitch struct {
    gpio_clk_chip string
    gpio_clk int
    gpio_dt_chip string
    gpio_dt int
    gpio_sw_chip string
    gpio_sw int
    outA *gpiocdev.Line  //CLK
    outB *gpiocdev.Line  //DT
    outSW *gpiocdev.Line //SS
    evtTime time.Time
    pb_evtTime time.Time
    pb_state bool
    longPthr int64
    long2Pthr int64
    swLongPress func(fn int)
    swLong2Press func(fn int)
    swShortPress func(fn int)
    rotHandler func(dir bool) // rotation handler, true = cw, false = ccw
}

func (sw *RotarySwitch) RotarySwitchInit() {
    var err error

    period := 10 * time.Millisecond
    sw.outA, err = gpiocdev.RequestLine(sw.gpio_clk_chip, sw.gpio_clk, gpiocdev.WithEventHandler(sw.rotaryHandler), gpiocdev.WithBothEdges, gpiocdev.WithDebounce(period))
    if err != nil {
            log.Fatal("RequestLine returned error: %w", err)
    }
    sw.outA.Reconfigure(gpiocdev.WithPullUp)

    sw.outB, err = gpiocdev.RequestLine(sw.gpio_dt_chip, sw.gpio_dt, gpiocdev.WithEventHandler(sw.rotaryHandler), gpiocdev.WithBothEdges, gpiocdev.WithDebounce(period))
    if err != nil {
            log.Fatal("RequestLine returned error: %w", err)
    }
    sw.outB.Reconfigure(gpiocdev.WithPullUp)

    sw.outSW, err = gpiocdev.RequestLine(sw.gpio_sw_chip, sw.gpio_sw, gpiocdev.WithEventHandler(sw.rotaryHandler), gpiocdev.WithBothEdges)
    if err != nil {
            log.Fatal("RequestLine returned error: %w", err)
    }
    sw.outSW.Reconfigure(gpiocdev.WithPullUp)
    sw.pb_state = false
}

func (sw *RotarySwitch) Close() error {
    if sw.outA != nil {
        log.Println("Closing line")
        _ = sw.outA.Close()
        _ = sw.outB.Close()
        _ = sw.outSW.Close()
    }
    return nil
}

func (sw *RotarySwitch) rotaryHandler(evt gpiocdev.LineEvent) {

      if evt.Offset == sw.gpio_sw {
        if evt.Type == gpiocdev.LineEventRisingEdge {
            diff := time.Now().Sub(sw.pb_evtTime)
            if diff.Milliseconds() < 5 {
                return
            }
            sw.pb_evtTime = time.Now()
            if sw.pb_state == false {
                sw.pb_state = true

                diff := time.Now().Sub(sw.pb_evtTime)
                if diff.Milliseconds() > sw.long2Pthr {
                    sw.swLong2Press(2)
                } else if diff.Milliseconds() > sw.longPthr {
                    sw.swLongPress(1)
                } else {
                    sw.swShortPress(0)
                }
            }
        } else if evt.Type == gpiocdev.LineEventFallingEdge {
            sw.pb_state = false
        }
      return
      }

      var rd, upd bool

      diff := time.Now().Sub(sw.evtTime)
      if diff.Milliseconds() < 10 { return }

      Clk, _ := sw.outA.Value()
      Dt, _ := sw.outB.Value()
      if evt.Type == gpiocdev.LineEventFallingEdge {
        if Clk == 1 && Dt == 0 {
            rd = true
            upd = true
        } else if Clk == 0 && Dt == 1 {
            rd = false
            upd = true
        }
      } else  if evt.Type == gpiocdev.LineEventRisingEdge {
        if Clk == 0 && Dt == 1 {
            rd = true
            upd = true
        } else if Clk == 1 && Dt == 0 {
            rd = false
            upd = true
        }
      }

    if upd {
     if rd {
        sw.rotHandler(true)
     } else {
        sw.rotHandler(false)
     }
     upd = false
     sw.evtTime = time.Now()
    }
}
