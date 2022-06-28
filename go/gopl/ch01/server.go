// A minimal "echo" and counter server
package main

import (
    "fmt"
    "log"
    "net/http"
    "sync"
    "image"
    "image/color"
    "image/gif"
    "io"
    "math"
    "math/rand"
)

var mu sync.Mutex
var count int

func main() {
    http.HandleFunc("/", handler)
    http.HandleFunc("/count", counter)

    lissajousHandler := func(w http.ResponseWriter, r *http.Request) {
        lissajous(w)
    }
    http.HandleFunc("/lissajous", lissajousHandler)

    log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

// handler echoes the HTTP request
func handler(w http.ResponseWriter, r *http.Request) {
    mu.Lock()
    count++
    mu.Unlock()

    fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto)
    for k, v := range r.Header {
        fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
    }
    fmt.Fprintf(w, "Host = %q\n", r.Host)
    fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)
    if err := r.ParseForm(); err != nil {
        log.Print(err)
    }
    for k, v := range r.Form {
        fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
    }
}

//counter echoes the number of calls so far
func counter(w http.ResponseWriter, r *http.Request) {
    mu.Lock()
    fmt.Fprintf(w, "Count %d\n", count)
    mu.Unlock()
}

// setup for the lissajous function
var green1 = color.RGBA{0x00, 0x33, 0x00, 0xff}
var green2 = color.RGBA{0x00, 0x99, 0x00, 0xff}
var green3 = color.RGBA{0x00, 0xcc, 0x00, 0xff}
var palette = []color.Color{color.White, green1, green2, green3}

const (
    whiteIndex = 0 // first color in palette
    blackIndex = 1 // next color in palette
    greenIndex = 2 // final color in palette
)

func lissajous(out io.Writer) {
    const (
        cycles  = 5     // number of complete x oscillator revolutions
        res     = 0.001 // angular resolution
        size    = 100   // image canvas covers [-size..+size]
        nframes = 64    // number of animation frames
        delay   = 8     // delay between frames in 10ms units
    )

    freq := rand.Float64() * 3.0 // relative freq of y oscillator
    anim := gif.GIF{LoopCount: nframes}
    phase := 0.0 // phase difference

    for i := 0; i < nframes; i++ {
        rect := image.Rect(0, 0, 2*size+1, 2*size+1)
        img := image.NewPaletted(rect, palette)
        
        for t := 0.0; t < cycles*2*math.Pi; t += res {
            x := math.Sin(t)
            y := math.Sin(t*freq + phase)
            img.SetColorIndex(
                size+int(x*size+0.5),
                size+int(y*size+0.5),
                uint8(i%3+1))
        }
        
        phase += 0.1
        anim.Delay = append(anim.Delay, delay)
        anim.Image = append(anim.Image, img)
    }
    gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}

