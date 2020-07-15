//
//  Copyright (c) 2020 Sugesh Chandran
//
package main

import (
	"fmt"
	"os"
	"time"
	"syscall"
	"os/signal"
	flags "github.com/jessevdk/go-flags"
	"ovsdb-statsd-client/config"
	"ovsdb-statsd-client/internal/ovsdb-reader"
	"ovsdb-statsd-client/internal/statsD-writer"
	"ovsdb-statsd-client/pkg/errors"
)

func run(reader *ovsdbreader.OVSDBReader, writer *statsdwriter.SWriter) {
	report := reader.ReadOVSDB()
	//reader.DisplayReport()
	writer.Write(report)
}

func statsDLoop() {
	reader := ovsdbreader.CreateNewOVSDBReader(&config.InitConfig.OvsDBConf)
	err := reader.ConnectDB()
	if err != errors.ErrNil {
		return
	}
	writer:= statsdwriter.CreateSWriter(&config.InitConfig.StatsDConf)
	err = writer.Connect()
	if err != errors.ErrNil {
		return
	}
	exitsignal := make(chan os.Signal, 1)
    signal.Notify(exitsignal, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-exitsignal:
			reader.CloseDBConn()
			writer.Disconnect()
		default:
			start := time.Now()
			run(reader, writer)
			duration := time.Now().Sub(start)
			duration = time.Second -duration
			time.Sleep(duration)
		}

	}

}
func main() {

 	// command line flags
 	var opts struct {
		ConfFile  string  `long:"config" short:"c" default:"./config/config.yaml" description:"config : Input configuration file"`
	 }

	// parse said flags
	_, err := flags.Parse(&opts)
	if err != nil {
		if e, ok := err.(*flags.Error); ok {
			if e.Type == flags.ErrHelp {
				os.Exit(0)
			}
		}
		fmt.Printf("Error: %+v\n", err)
		os.Exit(1)
	}
	if opts.ConfFile == "" {
		err = config.ParseYAMLConfig("")
	} else {
		err = config.ParseYAMLConfig(opts.ConfFile)
	}
	if err != errors.ErrNil {
		return
	}
	statsDLoop()
 }

