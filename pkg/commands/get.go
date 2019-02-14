package commands

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-getter"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
)

type Mode int

const _Mode_name = "ANYFILEDIR"

const (
	ANY Mode = iota
	FILE
	DIR
)

var _Mode_index = [...]uint8{0, 3, 7, 10}

func load(src, dst string) {
	var moder = getter.ClientModeAny

	// Get the pwd
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting wd: %v", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	// Build the client
	client := &getter.Client{
		Ctx:  ctx,
		Src:  src,
		Dst:  dst,
		Pwd:  pwd,
		Mode: moder,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	errChan := make(chan error, 2)
	go func() {
		defer wg.Done()
		defer cancel()
		if err := client.Get(); err != nil {
			errChan <- err
		}
	}()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	select {
	case sig := <-c:
		signal.Reset(os.Interrupt)
		cancel()
		wg.Wait()
		log.Printf("signal %v", sig)
	case <-ctx.Done():
		wg.Wait()
		log.Printf("success!")
	case err := <-errChan:
		wg.Wait()
		log.Fatalf("Error downloading: %s", err)
	}
}

func (i Mode) String() string {
	if i < 0 || i >= Mode(len(_Mode_index)-1) {
		return "Mode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Mode_name[_Mode_index[i]:_Mode_index[i+1]]
}
