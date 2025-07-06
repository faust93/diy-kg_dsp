// SH1106 OLED control for Go
// by faust93
// Based on https://github.com/jburrell7/RPi-Pico-OLED-DRIVER
package main

import (
        "fmt"
        "os"
        "unsafe"
        "syscall"
        "time"
        _ "reflect"
       )

var DBG bool = false

// supported OLED types
const (
    SH1106  = 0
    SSD1306 = 1
)
const (
    W_96  = 96
    W_128 = 128
    W_132 = 132
    H_16  = 16
    H_32  = 32
    H_64  = 64
)

const (
    BLACK = 0
    WHITE = 1
)

const (
    NORMAL_SIZE = 0
    DOUBLE_SIZE = 1

    OLED_FONT_HEIGHT = 8
    OLED_FONT_WIDTH = 6
)

const (
    NO_SCROLLING = 0
    HORIZONTAL_RIGHT = 0x26
    HORIZONTAL_LEFT = 0x27
    DIAGONAL_RIGHT = 0x29
    DIAGONAL_LEFT = 0x2A
)

const (
    HOLLOW = 0
    SOLID = 1
)

const (
    OLED_DATA     = 0x40
    OLED_CONTRAST = 0x81
    OLED_INVERT   = 0xa7
    OLED_OFF      = 0xae
    OLED_ON       = 0xaf
    )

const (
    I2C_SLAVE = 0x0703
    I2C_M_RD  = 0x0001
    I2C_RDWR  = 0x0707
)

type OLED struct {
    i2c I2C

    oled_model uint8
    height, width int16
    usingOffset bool

    fontInverted bool
    color uint8
    scaling uint8
    scroll_type uint8

    x,y int16

    pages int16
    bufsize int16
    buffer [1024]byte
}

// I2C stuff
type I2C struct {
    bus int
    device uint8
    fh *os.File
    led int
    in, out []byte
}

type i2c_msg struct {
    addr      uint16
    flags     uint16
    len       uint16
    __padding uint16
    buf       uintptr
}

type i2c_rdwr_ioctl_data struct {
    msgs  uintptr
    nmsgs uint32
}

func (p *I2C) Init() {
    var err error
    p.fh, err = os.OpenFile(fmt.Sprintf("/dev/i2c-%d", p.bus), os.O_RDWR, 0600)
    if err != nil {
        panic(err)
    }
    if _, _, errno := syscall.Syscall6(syscall.SYS_IOCTL, p.fh.Fd(), I2C_SLAVE, uintptr(p.device), 0, 0, 0); errno != 0 {
        fmt.Printf("Error I2C device init: %x Device: %x\n", p.bus, p.device)
        panic(err)
    }

    p.in = make([]byte, 150)
    p.out = make([]byte, 150)

    if DBG {
        fmt.Printf("I2C Bus: %x Device: %x\n", p.bus, p.device)
    }
}

func (p *I2C) Write(buf *[]byte) int {
    n, err := p.fh.Write(*buf)
    if err != nil {
        fmt.Println("I2C write error")
        n = 0
    }
    return n
}

func (p *I2C) WriteN(buf *[]byte, len int) int {
    n, err := p.fh.Write((*buf)[:len])
    if err != nil {
        fmt.Println("I2C write error")
        n = 0
    }
    return n
}

func (p *I2C) Read(buf []byte) int {
    n, err := p.fh.Read(buf)
    if err != nil {
        fmt.Println("I2C read error")
        n = 0
    }
    return n
}

func (p *I2C) Clean_Up() {
    for i := range(p.in){ p.in[i] = 0; p.out[i] = 0 }
}

func (p *I2C) Close() {
    p.fh.Close()
}

func (p *I2C) Rdwr(rw bool, size int) int {
    var error int = 0

    msg := []i2c_msg{
        {
            addr: uint16(p.device),
            flags: 0x4000,
            buf: uintptr(unsafe.Pointer(&p.out[0])),
            len: uint16(size),
        },
//        {
//            addr: uint16(p.device),
//            flags: uint16(I2C_M_RD),
//            buf: uintptr(unsafe.Pointer(&p.in[0])),
//            len: uint16(size),
//        },
    }

    if DBG {
        fmt.Println(p.out[:size])
    }

    err := i2c_transfer(p.fh, &msg[0], 1)
    if err != nil {
        if DBG {
            fmt.Println("I2C Write error")
        }
        for i:=0; i < 5; i++ {
            time.Sleep(80 * time.Millisecond)
            err := i2c_transfer(p.fh, &msg[0], 1)
            error++
            if err == nil { error = 0; break }
        }
    }

    if !rw {
        return error
    }
/*
    time.Sleep(80 * time.Millisecond)
    err = i2c_transfer(p.fh, &msg[1], 1)

    if DBG {
        fmt.Println(p.in[:size])
    }


    if (p.in[0] != 22 && p.in[1] != 22) {
        for i:=0; i < 5; i++ {
            if DBG {
                fmt.Println("I2C Read error")
            }
            _ = i2c_transfer(p.fh, &msg[0], 1)
            time.Sleep(80 * time.Millisecond)
            _ = i2c_transfer(p.fh, &msg[1], 1)
            error++
            if (p.in[0] == 22 && p.in[1] == 22) { error = 0; break }
        }
    }
*/
    return error
}

func i2c_transfer(f *os.File, msgs *i2c_msg, n int) (err error) {
    data := i2c_rdwr_ioctl_data{
        msgs:  uintptr(unsafe.Pointer(msgs)),
        nmsgs: uint32(n),
    }
    err = nil
    _, _, errno := syscall.Syscall(
        syscall.SYS_IOCTL,
        uintptr(f.Fd()),
        uintptr(I2C_RDWR),
        uintptr(unsafe.Pointer(&data)),
    )
    if (errno != 0) {
        err = errno
    }
    return err
}


//  OLED stuff
func (o *OLED) Close() {
    o.i2c.Close()
}

func (o *OLED) oledInit(i2c_bus int, i2c_addr uint8) int {

    var comPin0, comPin1 byte

    o.i2c = I2C{bus: i2c_bus, device: i2c_addr}
    o.i2c.Init()

    o.fontInverted  = false
    o.color         = WHITE
    o.scaling       = NORMAL_SIZE
    o.x             = 0
    o.y             = 0
    o.pages         = ((o.height + 7) / 8)
    o.scroll_type   = NO_SCROLLING

    o.usingOffset = false
    if o.width == W_132 {
        o.width = W_128
        o.usingOffset = true
    }
    o.bufsize = (int16)(o.pages * o.width)

    comPin0 = 0xDA  // com pins hardware configuration
    if o.width == W_128 && o.height == H_32 {
        comPin1 = 0x02
    } else if o.width == W_128 && o.height == H_64 {
        comPin1 = 0x12
    } else if o.width == W_96 && o.height == H_16 {
        comPin1 = 0x02
    }

    var params = []byte{
    0x00,                           // command
    0xAE,                           // display off
    0xD5,                           // clock divider
    0x80,
    0xA8,                           // multiplex ratio
    byte(o.height - 1),
    0xD3,                           // no display offset
    0x00,
    0x40,                           // start line address=0
    0x8D,                           // enable charge pump
    0x14,
    0x20,                           // memory adressing mode=horizontal
    0x00,
    0xA1,                           // segment remapping mode
    0xC8,                           // COM output scan direction
    comPin0,
    comPin1,
    0x81,                           // contrast control
    0x80,
    0xD9,                           // pre-charge period
    0x22,
    0xDB,                           // set vcomh deselect level
    0x20,
    0xA4,                           // output RAM to display
    0xA6,                           // display mode A6=normal, A7=inverse
    0x2E}                           // stop scrolling

    wb := o.i2c.Write(&params)
    if wb == 0 {
        panic("Unable to initialize OLED device")
    }

    o.oledSet_power(true)
    o.oledClear(BLACK)
    o.oledDisplay()

    return 0
}

func (o *OLED) oledUseOffset(enabled bool) {
    if o.oled_model == SH1106 {
        o.usingOffset = enabled
    }
}

func (o *OLED) oledSet_power(enable bool) {
    if enable {
        o.i2c.out[0] = 0x00
        o.i2c.out[1] = 0x8D
        o.i2c.out[2] = 0x14
        o.i2c.out[3] = 0xAF
        o.i2c.Rdwr(false, 4)
    } else {
        o.i2c.out[0] = 0x00
        o.i2c.out[1] = 0xAE
        o.i2c.out[2] = 0x8D
        o.i2c.out[3] = 0x10
        o.i2c.Rdwr(false, 4)
    }
}

func (o *OLED) oledClear(color uint8) {
    var c uint8

    if color == WHITE {
        c = 0xFF
    } else {
        c = 0x00
    }
    var i int16
    for i = 0; i < o.bufsize; i++ { o.buffer[i] = c }
    o.x = 0
    o.y = 0
}

func (o *OLED) oledSetCursor(x int16, y int16) {
    o.x = x
    o.y = y
}

func (o *OLED) oledSet_font_inverted(enabled bool) {
    o.fontInverted = enabled
}

func (o *OLED) oledSet_invert(enable bool) {
    o.i2c.out[0] = 0x00
    if enable {
        o.i2c.out[1] = 0xA7
    }  else {
        o.i2c.out[1] = 0xA6
    }
    o.i2c.Rdwr(false, 2)
}

func (o *OLED) oledSet_contrast(contrast uint8) {
    o.i2c.out[0] = 0x00
    o.i2c.out[1] = 0x81
    o.i2c.out[2] = contrast
    o.i2c.Rdwr(false, 3)
}

func (o *OLED) oledSet_scrolling(scroll_type uint8, first_page uint8, last_page uint8) {

    if scroll_type == NO_SCROLLING {
        o.i2c.out[0] = 0x00
        o.i2c.out[1] = 0x2E
        o.i2c.Rdwr(false, 2)
        o.scroll_type = scroll_type
        return
    }

    if scroll_type == DIAGONAL_LEFT || scroll_type == DIAGONAL_RIGHT {
        o.i2c.out[0] = 0x00
        o.i2c.out[1] = 0x2E
        o.i2c.out[2] = 0xA3
        o.i2c.out[3] = 0x00
        o.i2c.out[4] = byte((o.height - 1))
        o.i2c.out[5] = scroll_type
        o.i2c.out[6] = 0x00
        o.i2c.out[7] = first_page
        o.i2c.out[8] = 0x00
        o.i2c.out[9] = last_page
        o.i2c.out[10] = 0x01
        o.i2c.out[11] = 0x2F
        o.i2c.Rdwr(false, 12)

        o.scroll_type = scroll_type
        return
    }

    if scroll_type == HORIZONTAL_RIGHT || scroll_type == HORIZONTAL_LEFT {
        o.i2c.out[0] = 0x00
        o.i2c.out[1] = 0x2E
        o.i2c.out[2] = scroll_type
        o.i2c.out[3] = 0x00
        o.i2c.out[4] = first_page
        o.i2c.out[5] = 0x00
        o.i2c.out[6] = last_page
        o.i2c.out[7] = 0x00
        o.i2c.out[8] = 0xFF
        o.i2c.out[9] = 0x2F
        o.i2c.Rdwr(false, 10)
        o.scroll_type = scroll_type
        return
    }
}

func (o *OLED) oledDisplay() {
    var bufIndex uint16
    var index uint16 = 0
    var bytesToSend uint16 = 0

    var page int16
    for page = 0; page < o.pages; page++ {
        if o.oled_model == SH1106 {
            o.i2c.out[0] = 0x00
            o.i2c.out[1] = 0xB0 + byte(page)
            o.i2c.out[2] = 0x00
            o.i2c.out[3] = 0x10
            o.i2c.Rdwr(false, 4)
        } else {
            o.i2c.out[0] = 0x00
            o.i2c.out[1] = 0xB0 + byte(page)
            o.i2c.out[2] = 0x21
            o.i2c.out[3] = 0x00
            o.i2c.out[4] = byte(o.width - 1)
            o.i2c.Rdwr(false, 5)
        }
        o.i2c.out[0] = 0x40
        bufIndex    = 1
        bytesToSend = (uint16)(o.width + 1)

        if o.usingOffset {
        // send two dummy bytes if the width of the display
        //  is > 128 pixels
            o.i2c.out[1] = 0x00
            o.i2c.out[2] = 0x00
            bufIndex     = 3
            bytesToSend += 2
        }

        copy(o.i2c.out[bufIndex:], o.buffer[index:index+uint16(o.width)])
        o.i2c.Rdwr(false, int(bytesToSend))
        //o.i2c.WriteN(&o.i2c.out, int(bytesToSend))
        index += uint16(o.width)
    }
}

func (o *OLED) oledDraw_pixel(x int16, y int16, color uint8) {
    if x >= o.width || y >= o.height { return }
    if color == WHITE {
        o.buffer[x + (y / 8) * o.width] |= (1 << (y & 7)) // set bit
    } else {
        o.buffer[x + (y / 8) * o.width] &= ^(1 << (y & 7)) // clear bit
    }
}

func abs(x int16) int16 {
   if x < 0 {
      return -x
   }
   return x
}

func (o *OLED) oledDraw_line(x0 int16, y0 int16, x1 int16, y1 int16, color uint8) {
    // Algorithm copied from Wikipedia
    var dx int16 = abs(x1 - x0)
    var sx int16
    if x0 < x1 {
        sx = 1
    } else {
        sx = -1
    }
    var dy int16 = -abs(y1 - y0)
    var sy int16
    if y0 < y1 {
        sy = 1
    } else {
        sy = -1
    }
    var err int16 = dx + dy
    var e2 int16

    for {
        o.oledDraw_pixel(x0, y0, color)
        if x0 == x1 && y0 == y1 { break }
        e2 = 2 * err
        if e2 > dy {
            err += dy
            x0 += sx
        }
        if e2 < dx {
            err += dx
            y0 += sy
        }
    }
}

func (o *OLED) oledDraw_circle(x0 int16, y0 int16, radius int16, fillMode uint8, color uint8) {
    // Algorithm copied from Wikipedia
    var f int16 = 1 - radius
    var ddF_x int16 = 0
    var ddF_y int16 = -2 * radius
    var x int16 = 0
    var y int16 = radius

    if fillMode == SOLID {
        o.oledDraw_pixel(x0, y0 + radius, color)
        o.oledDraw_pixel(x0, y0 - radius, color)
        o.oledDraw_line(x0 - radius, y0, x0 + radius, y0, color)
    } else {
        o.oledDraw_pixel(x0, y0 + radius, color)
        o.oledDraw_pixel(x0, y0 - radius, color)
        o.oledDraw_pixel(x0 + radius, y0, color)
        o.oledDraw_pixel(x0 - radius, y0, color)
    }

    for x < y {
        if f >= 0 {
            y--
            ddF_y += 2
            f += ddF_y
        }
        x++
        ddF_x += 2
        f += ddF_x + 1

        if fillMode == SOLID {
            o.oledDraw_line(x0 - x, y0 + y, x0 + x, y0 + y, color)
            o.oledDraw_line(x0 - x, y0 - y, x0 + x, y0 - y, color)
            o.oledDraw_line(x0 - y, y0 + x, x0 + y, y0 + x, color)
            o.oledDraw_line(x0 - y, y0 - x, x0 + y, y0 - x, color)
        } else {
            o.oledDraw_pixel(x0 + x, y0 + y, color)
            o.oledDraw_pixel(x0 - x, y0 + y, color)
            o.oledDraw_pixel(x0 + x, y0 - y, color)
            o.oledDraw_pixel(x0 - x, y0 - y, color)
            o.oledDraw_pixel(x0 + y, y0 + x, color)
            o.oledDraw_pixel(x0 - y, y0 + x, color)
            o.oledDraw_pixel(x0 + y, y0 - x, color)
            o.oledDraw_pixel(x0 - y, y0 - x, color)
        }
    }
}

func (o *OLED) oledDraw_rectangle(x0 int16, y0 int16, x1 int16, y1 int16, fillMode uint8, color uint8) {
    // Swap x0 and x1 if in wrong order
    if x0 > x1 {
        tmp := x0
        x0 = x1
        x1 = tmp
    }
    // Swap y0 and y1 if in wrong order
    if y0 > y1 {
        tmp := y0
        y0 = y1
        y1 = tmp
    }
    if fillMode == SOLID {
        for  y := y0; y <= y1; y++ {
            o.oledDraw_line(x0, y, x1, y, color)
        }
    } else {
        o.oledDraw_line(x0, y0, x1, y0, color)
        o.oledDraw_line(x0, y1, x1, y1, color)
        o.oledDraw_line(x0, y0, x0, y1, color)
        o.oledDraw_line(x1, y0, x1, y1, color)
    }
}

func (o *OLED) oledDraw_hbar(x0 int16, y0 int16, height int16, width int16) {
    var i int16
    for i = 0; i < width; i+=2 {
        o.oledDraw_line(x0 + i, y0, x0 + i, y0 + height, WHITE)
    }
}

func (o *OLED) oledDraw_smBar(x0 int16, height int16, width int16) {
    var h int16

    h = o.height - height
    for  y := o.height; y >= h; y-=2 {
        o.oledDraw_line(x0, y, x0 + width, y, WHITE)
    }
}

func (o *OLED) oledDraw_string(x int16, y int16, s string, scaling uint8, color uint8) {
    for i:=0; i < len(s); i++ {
        o.oledDraw_character(x, y, s[i], scaling, color)
        if scaling == DOUBLE_SIZE {
            x += 12
        } else { // NORMAL_SIZE
            x += 6
        }
    }
}

func (o *OLED) oledWriteChar(c byte) int {
    n := 1
    n = o.oledDraw_character(o.x, o.y, c, o.scaling, o.color)
    o.x += OLED_FONT_WIDTH
    return n
}

func (o *OLED) oledDraw_character(x int16, y int16, c byte, scaling uint8, color uint8) int {
    // Invalid position
    if x >= o.width || y >= o.height || c < 32 { return 0 }

    // Remap extended Latin1 character codes

    switch c {
        case 252: /* u umlaut */
            c = 127
        case 220: /* U umlaut */
            c = 128
        case 228: /* a umlaut */
            c = 129
        case 196: /* A umlaut */
            c = 130
        case 246: /* o umlaut */
            c = 131
        case 214: /* O umlaut */
            c = 132
        case 176: /* degree   */
            c = 133
        case 223: /* szlig    */
            c = 134
    }

    font_index := (uint16(c)*6)// - 32)*6

    // Invalid character code/font index
    if font_index >= uint16(len(oled_font5x7)) { return 0 }

    tmp := make([]byte, 6)
    copy(tmp, oled_font5x7[font_index:])
    o.oledDraw_bytes(x, y, &tmp, 6, scaling, color)
    return 1
}

func (o *OLED) oledDraw_bytes(x int16, y int16, data *[]byte, size uint8, scaling uint8, color uint8) {
    var column uint8
    for column = 0; column < size; column++ {
        b := (*data)[column]

        if scaling == DOUBLE_SIZE {
            // Stretch vertically
            var w uint16 = 0
            var bit uint8
            for bit = 0; bit < 7; bit++ {
                if b & (1 << bit) != 0 {
                    pos := bit << 1
                    w |= ((1 << pos) | (1 << (pos + 1)))
                }
            }

            // Output 2 times to strech hozizontally
            o.oledDraw_byte(x, y, byte(w & 0xFF), color)
            o.oledDraw_byte(x, y + 8, byte((w >> 8)), color)
            x++
            o.oledDraw_byte(x, y, byte(w & 0xFF), color)
            o.oledDraw_byte(x, y + 8, byte((w >> 8)), color)
            x++
        } else { // NORMAL_SIZE
            x++
            o.oledDraw_byte(x, y, b, color)
        }
    }
}

func (o *OLED) oledDraw_byte(x int16, y int16, b byte, color uint8) {
    // Invalid position
    if x >= o.width || y >= o.height { return }

    buffer_index := y / 8 * o.width + x

    if o.fontInverted {
        b^=255
    }

    if color == WHITE {
        // If the y position matches a page, then it goes quicker
        if y % 8 == 0 {
            if buffer_index < o.bufsize {
                o.buffer[buffer_index] |= b
            }
        } else {
            w := uint16(b) << (y % 8)
            if buffer_index < o.bufsize {
                o.buffer[buffer_index] |= byte((w & 0xFF))
            }
            buffer_index2 := buffer_index + o.width
            if buffer_index2 < o.bufsize {
                o.buffer[buffer_index2] |= byte(w >> 8)
            }
        }
    } else {
        // If the y position matches a page, then it goes quicker
        if y % 8 == 0 {
            if buffer_index < o.bufsize {
                o.buffer[buffer_index] &= ^b
            }
        } else {
            w := uint16(b) << (y % 8)
            if buffer_index < o.bufsize {
                o.buffer[buffer_index] &= ^byte(w & 0xFF)
            }
            buffer_index2 := buffer_index + o.width
            if buffer_index2 < o.bufsize {
                o.buffer[buffer_index2] &= ^byte(w >> 8)
            }
        }
    }
}

func (o *OLED) oledRect_invert(x int16, y int16, w int16, h int16) {
    // Invalid position
    if x >= o.width || y >= o.height { return }

    var i,v int16
    for v = 0; v<h; v++ {
        x2 := x
        for i = 0; i<w; i++ {
            o.buffer[x2 + (y / 8) * o.width] ^= (1 << (y & 7)) // clear bit
            x2++
        }
    y++
    }
}

func (o *OLED) oledDraw_bitmap(x int16, y int16, bitmap_width uint16, bitmap_height uint16, data *[]byte, color uint8) {
    num_pages := (bitmap_height + 7) / 8
    tmp := make([]byte, bitmap_width)
    var i,page uint16
    for page = 0; page < num_pages; page++ {
        copy(tmp, (*data)[i:])
        o.oledDraw_bytes(x, y, &tmp, uint8(bitmap_width), NORMAL_SIZE, color)
        i += bitmap_width
        y += 8
    }
}

func (o *OLED) oledScroll_up(num_lines int16) {
    // Scroll full pages, fast algorithm
    scroll_pages := num_lines / 8
    var i,x int16
    for i = 0; i < o.pages; i++ {
        for x = 0; x < o.width; x++ {
            index := i * o.width + x
            index2 := (i + scroll_pages) * o.width + x
            if index2 < o.bufsize {
                o.buffer[index] = o.buffer[index2]
            } else  {
                o.buffer[index] = 0
            }
        }
    }
    num_lines -= scroll_pages * 8
}

func oledToCol(x int) int {
    return x/OLED_FONT_WIDTH
}

func oledToRow(y int) int {
    return y/OLED_FONT_HEIGHT
}

func oledToX(col int) int {
    return col*OLED_FONT_WIDTH
}

func oledToY(row int) int {
    return row*OLED_FONT_HEIGHT
}