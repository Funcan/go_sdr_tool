package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/funcan/soapyradiotool/mainwindow"
	"github.com/funcan/soapyradiotool/mathtools"
)

type U8Data struct {
	Len        int64
	SampleRate int64
	I          []uint8
	J          []uint8
	Real       []float64
}

func ReadU8(filename string, samplerate int64) (*U8Data, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	stats, err := f.Stat()
	if err != nil {
		return nil, err
	}

	size := stats.Size()
	bytes := make([]byte, size)
	reader := bufio.NewReader(f)
	_, err = reader.Read(bytes)

	ret := U8Data{
		Len:         size / 2,
		SampleRate:  samplerate,
		I:           make([]uint8, size/2),
		J:           make([]uint8, size/2),
		Real:        make([]float64, size/2),
	}

	for i := int64(0); i < size / 2; i++ {
		ret.I[i] = uint8FromBuf(bytes, i * 2)
		ret.J[i] = uint8FromBuf(bytes, i * 2 + 1)

		ret.Real[i] = math.Sqrt(float64(ret.I[i]) * float64(ret.J[i]))
	}

	return &ret, nil
}

func uint8FromBuf(buf []byte, offset int64) uint8 {
	var ret uint8

	tmp := bytes.NewBuffer(buf[offset : offset+1])

	err := binary.Read(tmp, binary.LittleEndian, &ret)
	// FIXME: Pass up err rather than panic
	if err != nil {
		panic(err)
	}

	return ret
}

func readFile(arg interface{}) {
	filename, ok := arg.(string)
	if !ok {
		log.Printf("readFile invalid parameter %T", filename)
		return
	}

	data, err := ReadU8(filename, 2000000)
	if err != nil {
		log.Printf("readFile error: %s", err)
		return
	}

	f := registerSource("show data")
	f(data.Real)
}

func oldmain() {
	filename := "/home/duncan/src/rtl_433_tests/tests/generic_remote/01/gfile001.cu8"

	data, err := ReadU8(filename, 2000000)
	if err != nil {
		panic(err)
	}

	mean := mathtools.Mean(data.Real)
	stdDev := mathtools.StdDev(data.Real)


	squelched := mathtools.Squelch(data.Real, mean + stdDev)
	//avg := mathtools.RollingAverage(squelched, 5)

//	for _, d := range(squelched) {
//		fmt.Printf("%f\n", d)
//	}

	hist := mathtools.Histogram(squelched, 100)

	for _, d := range(hist) {
		fmt.Printf("%d\n", d)
	}
}

var listeners map[string][]func(e interface{})

func registerListener(event string, f func(e interface{})) {
	if listeners == nil {
		listeners = make(map[string][]func(e interface{}))
	}

	val, ok := listeners[event]

	if !ok {
		listeners[event] = make([]func(e interface{}), 0)
		val = listeners[event]
	}

	listeners[event] = append(val, f)
}

func registerSource(event string) func(e interface{}) {
	return func(e interface{}) {
		log.Printf("event '%s' -> %T", event, e)
		for s, funcs := range(listeners) {
			if s == event {
				for _, f := range(funcs) {
					f(e)
				}
			}
		}
	}
}

func main() {
	registerListener("load file", readFile)
	mainwindow.Setup(registerListener, registerSource)
}
