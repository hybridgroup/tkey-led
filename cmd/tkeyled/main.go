package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/pflag"
	"github.com/tillitis/tkeyclient"
)

var (
	cmdSetLED    = appCmd{0x01, "cmdSetLED", tkeyclient.CmdLen32}
	rspSetLED    = appCmd{0x02, "rspSetLED", tkeyclient.CmdLen4}
	cmdSetTiming = appCmd{0x03, "cmdSetTiming", tkeyclient.CmdLen32}
	rspSetTiming = appCmd{0x04, "rspSetTiming", tkeyclient.CmdLen4}
	cmdBlinking  = appCmd{0x05, "cmdBlinking", tkeyclient.CmdLen4}
	rspBlinking  = appCmd{0x06, "rspBlinking", tkeyclient.CmdLen4}
)

type appCmd struct {
	code   byte
	name   string
	cmdLen tkeyclient.CmdLen
}

func (c appCmd) Code() byte {
	return c.code
}

func (c appCmd) CmdLen() tkeyclient.CmdLen {
	return c.cmdLen
}

func (c appCmd) Endpoint() tkeyclient.Endpoint {
	return tkeyclient.DestApp
}

func (c appCmd) String() string {
	return c.name
}

type Blinker struct {
	tk *tkeyclient.TillitisKey // A connection to a TKey
}

// New allocates a struct for communicating with the timer app running
// on the TKey. You're expected to pass an existing connection to it,
// so use it like this:
//
//	tk := tkeyclient.New()
//	err := tk.Connect(port)
//	blinker := NewBlinker(tk)
func NewBlinker(tk *tkeyclient.TillitisKey) Blinker {
	var blinker Blinker

	blinker.tk = tk

	return blinker
}

func (b Blinker) setBool(sendCmd appCmd, expectedReceiveCmd appCmd, on bool) error {
	id := 2
	tx, err := tkeyclient.NewFrameBuf(sendCmd, id)
	if err != nil {
		return fmt.Errorf("NewFrameBuf: %w", err)
	}

	// The boolean
	if on {
		tx[2] = 1
	} else {
		tx[2] = 0
	}
	tkeyclient.Dump("tx", tx)
	if err = b.tk.Write(tx); err != nil {
		return fmt.Errorf("Write: %w", err)
	}

	rx, _, err := b.tk.ReadFrame(expectedReceiveCmd, id)
	tkeyclient.Dump("rx", rx)
	if err != nil {
		return fmt.Errorf("ReadFrame: %w", err)
	}

	if rx[2] != tkeyclient.StatusOK {
		return fmt.Errorf("Command BAD")
	}

	return nil
}

// setInt sets an int with the command cmd
func (b Blinker) setInt(sendCmd appCmd, expectedReceiveCmd appCmd, i int) error {
	id := 2
	tx, err := tkeyclient.NewFrameBuf(sendCmd, id)
	if err != nil {
		return fmt.Errorf("NewFrameBuf: %w", err)
	}

	// The integer
	tx[2] = byte(i)
	tx[3] = byte(i >> 8)
	tx[4] = byte(i >> 16)
	tx[5] = byte(i >> 24)
	tkeyclient.Dump("tx", tx)
	if err = b.tk.Write(tx); err != nil {
		return fmt.Errorf("Write: %w", err)
	}

	rx, _, err := b.tk.ReadFrame(expectedReceiveCmd, id)
	tkeyclient.Dump("rx", rx)
	if err != nil {
		return fmt.Errorf("ReadFrame: %w", err)
	}

	if rx[2] != tkeyclient.StatusOK {
		return fmt.Errorf("Command BAD")
	}

	return nil
}

func (b Blinker) SetLED(led int) error {
	return b.setInt(cmdSetLED, rspSetLED, led)
}

func (b Blinker) SetTiming(ms int) error {
	return b.setInt(cmdSetTiming, rspSetTiming, ms)
}

func (b Blinker) Blinking(on bool) error {
	return b.setBool(cmdBlinking, rspBlinking, on)
}

func main() {
	var devPath string
	var led, timing, speed int
	var blinking, verbose, helpOnly bool
	pflag.CommandLine.SortFlags = false
	pflag.StringVar(&devPath, "port", "",
		"Set serial port device `PATH`. If this is not passed, auto-detection will be attempted.")
	pflag.IntVar(&speed, "speed", tkeyclient.SerialSpeed,
		"Set serial port speed in `BPS` (bits per second).")
	pflag.BoolVar(&verbose, "verbose", false,
		"Enable verbose output.")
	pflag.IntVar(&led, "led", 0,
		"Set LED")
	pflag.IntVar(&timing, "timing", 500,
		"Set blink timing (in ms).")
	pflag.BoolVar(&blinking, "blinking", true,
		"Set blinking on or off.")
	pflag.BoolVar(&helpOnly, "help", false, "Output this help.")
	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n%s", os.Args[0],
			pflag.CommandLine.FlagUsagesWrapped(80))
	}
	pflag.Parse()

	if helpOnly {
		pflag.Usage()
		os.Exit(0)
	}

	if !verbose {
		tkeyclient.SilenceLogging()
	}

	if devPath == "" {
		var err error
		devPath, err = tkeyclient.DetectSerialPort(true)
		if err != nil {
			os.Exit(1)
		}
	}

	tk := tkeyclient.New()
	fmt.Printf("Connecting to device on serial port %s ...\n", devPath)
	if err := tk.Connect(devPath, tkeyclient.WithSpeed(speed)); err != nil {
		fmt.Printf("Could not open %s: %v\n", devPath, err)
		os.Exit(1)
	}
	exit := func(code int) {
		if err := tk.Close(); err != nil {
			fmt.Printf("tk.Close: %v\n", err)
		}
		os.Exit(code)
	}
	handleSignals(func() { exit(1) }, os.Interrupt, syscall.SIGTERM)

	bl := NewBlinker(tk)

	if err := bl.SetLED(led); err != nil {
		fmt.Printf("SetLED: %v\n", err)
		exit(1)
	}

	// if err := bl.SetTiming(timing); err != nil {
	// 	fmt.Printf("SetTiming: %v\n", err)
	// 	exit(1)
	// }

	// if err := bl.Blinking(blinking); err != nil {
	// 	fmt.Printf("Blinking: %v\n", err)
	// 	exit(1)
	// }

	exit(0)
}

func handleSignals(action func(), sig ...os.Signal) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, sig...)
	go func() {
		for {
			<-ch
			action()
		}
	}()
}
