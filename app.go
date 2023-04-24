package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"runtime"
	"time"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type AppConfig struct {
	Subscriptions []Subscription `json:"subscriptions"`
}

type Subscription struct {
	LastUpdate  int64       `json:"last_update"` // time.Now().UnixMilli()
	DisplayName string      `json:"display_name"`
	URL         string      `json:"url"`
	Config      interface{} `json:"config"`
}

// App struct
type App struct {
	ctx           context.Context
	ctxCancelFunc *context.CancelFunc
	dataPath      string
	appConfigFile *os.File
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// ChanneledWriter Receives data from Write([]byte) function and sends it to its channel
type ChanneledIO struct {
	Data chan []byte
}

func NewChanneledIO(capacity int) (cio *ChanneledIO) {
	cio = &ChanneledIO{}
	if capacity > 0 {
		cio.Data = make(chan []byte, capacity)
	} else {
		cio.Data = make(chan []byte)
	}
	return
}

func (cio *ChanneledIO) Close() {
	close(cio.Data)
}

func (cio *ChanneledIO) Read(res []byte) (n int, err error) {
	data := <-cio.Data
	n = len(data)
	err = nil
	if n > len(res) {
		n = len(res)
		err = fmt.Errorf("not enough size to write all of the data: %d bytes needed, %d provided", len(data), n)
	}
	for i := 0; i < n; i++ {
		res[i] = data[i]
	}
	return
}

func (cio *ChanneledIO) Write(inp []byte) (n int, err error) {
	cio.Data <- inp
	n = len(inp)
	err = nil
	return
}

func (cio *ChanneledIO) GetWholeArray() []byte {
	return <-cio.Data
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	var err error
	a.ctx = ctx
	a.dataPath, err = os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	a.dataPath = path.Join(a.dataPath, "ecorp-vpn")
	os.MkdirAll(a.dataPath, os.ModePerm)
	a.appConfigFile, err = os.OpenFile(path.Join(a.dataPath, "appConfig.json"), os.O_CREATE|os.O_RDWR, 0660)
	if err != nil {
		panic(err)
	}
	stat, err := a.appConfigFile.Stat()
	if err != nil {
		panic(err)
	}
	if stat.Size() == 0 {
		a.appConfigFile.Seek(0, 0)
		a.appConfigFile.WriteString("{[]}")
	}
}

func (a *App) shutdown(ctx context.Context) {
	if a.ctxCancelFunc != nil {
		(*a.ctxCancelFunc)()
	}
	fmt.Println("Goodbye!")
}

// ParseConfig reads configurations from a.appConfigFile
func (a *App) ParseConfig() AppConfig {
	var result AppConfig
	a.appConfigFile.Seek(0, 0)
	rawData, err := ioutil.ReadAll(a.appConfigFile)
	if err != nil {
		return result
	}
	_ = json.Unmarshal(rawData, &result)
	return result
}

// WriteConfig returns true on success
func (a *App) WriteConfig(cfg AppConfig) bool {
	rawData, err := json.Marshal(cfg)
	if err != nil {
		return false
	}
	a.appConfigFile.Seek(0, 0)
	n, err := a.appConfigFile.Write(rawData)
	if err != nil || n != len(rawData) {
		return false
	}
	err = a.appConfigFile.Truncate(int64(n))
	return err == nil
}

// Emits its events to "VPNLog" event name
func RunVPN(appCtx context.Context, ctx context.Context, cfg Subscription) (err error) {
	cfgBytes, err := json.Marshal(cfg.Config)
	if err != nil {
		return
	}
	executable_name := "./sing-box"
	if runtime.GOOS == "windows" {
		executable_name += ".exe"
	}
	cmd := exec.CommandContext(ctx, executable_name, "run", "-c", "stdin")
	stdin_reader, stdin_writer := io.Pipe()
	output_io := NewChanneledIO(0)
	cmd.Stdin = stdin_reader
	cmd.Stdout = output_io
	cmd.Stderr = output_io

	go func() {
		stdin_writer.Write(cfgBytes)
		stdin_writer.Close()
		stdin_reader.Close()

		defer cmd.Cancel()
		for {
			select {
			case data := <-output_io.Data:
				strData := string(data)
				wailsRuntime.EventsEmit(appCtx, "VPNLog", strData)
				fmt.Println(strData)
			case <-ctx.Done():
				fmt.Println("Shutting down sing-box")
				output_io.Close()
				cmd.Cancel()
				cmd.Wait()
				return
			}
		}
	}()
	return cmd.Start()
}

func (a *App) StartVPN(cfg Subscription) (err error) {
	if a.ctxCancelFunc != nil {
		(*a.ctxCancelFunc)()
	}
	vpnCtx := context.Background()
	vpnCtx, cancelFunc := context.WithCancel(vpnCtx)
	a.ctxCancelFunc = &cancelFunc

	wailsRuntime.EventsOff(a.ctx, "StopVPN")

	time.Sleep(50 * time.Millisecond)

	wailsRuntime.EventsOnce(a.ctx, "StopVPN", func(optionalData ...interface{}) {
		if a.ctxCancelFunc != nil {
			(*a.ctxCancelFunc)()
		}
		a.ctxCancelFunc = nil
	})

	return RunVPN(a.ctx, vpnCtx, cfg)
}
