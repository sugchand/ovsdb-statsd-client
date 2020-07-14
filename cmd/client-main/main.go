//
//  Copyright (c) 2020 Sugesh Chandran
//
package main

import (
	"fmt"
	"os"

	// "github.com/cactus/go-statsd-client/statsd"
	flags "github.com/jessevdk/go-flags"
	"ovsdb-statsd-client/config"
	"ovsdb-statsd-client/internal/ovsdb-reader"
	"ovsdb-statsd-client/internal/statsD-writer"
	"ovsdb-statsd-client/pkg/errors"
)


// func main() {
// 	ExampleClient_listDatabases()
// }

// func main() {

// 	// command line flags
// 	var opts struct {
// 		HostPort  string        `long:"host" default:"127.0.0.1:8125" description:"host:port of statsd server"`
// 		Prefix    string        `long:"prefix" default:"test-client" description:"Statsd prefix"`
// 		StatType  string        `long:"type" default:"count" description:"stat type to send. Can be one of: timing, count, gauge"`
// 		StatValue int64         `long:"value" default:"1" description:"Value to send"`
// 		Name      string        `short:"n" long:"name" default:"counter" description:"stat name"`
// 		Rate      float32       `short:"r" long:"rate" default:"1.0" description:"sample rate"`
// 		Volume    int           `short:"c" long:"count" default:"1000" description:"Number of stats to send. Volume."`
// 		Nil       bool          `long:"nil" description:"Use nil client"`
// 		Buffered  bool          `long:"buffered" description:"Use a buffered client"`
// 		Duration  time.Duration `short:"d" long:"duration" default:"10s" description:"How long to spread the volume across. For each second of duration, volume/seconds events will be sent."`
// 	}

// 	// parse said flags
// 	_, err := flags.Parse(&opts)
// 	if err != nil {
// 		if e, ok := err.(*flags.Error); ok {
// 			if e.Type == flags.ErrHelp {
// 				os.Exit(0)
// 			}
// 		}
// 		fmt.Printf("Error: %+v\n", err)
// 		os.Exit(1)
// 	}

// 	if opts.Nil && opts.Buffered {
// 		fmt.Printf("Specifying both nil and buffered together is invalid\n")
// 		os.Exit(1)
// 	}

// 	if opts.Name == "" || statsd.CheckName(opts.Name) != nil {
// 		fmt.Printf("Stat name contains invalid characters\n")
// 		os.Exit(1)
// 	}

// 	if statsd.CheckName(opts.Prefix) != nil {
// 		fmt.Printf("Stat prefix contains invalid characters\n")
// 		os.Exit(1)
// 	}

// 	config := &statsd.ClientConfig{
// 		Address:     opts.HostPort,
// 		Prefix:      opts.Prefix,
// 		ResInterval: time.Duration(0),
// 	}

// 	var client statsd.Statter
// 	if !opts.Nil {
// 		if opts.Buffered {
// 			config.UseBuffered = true
// 			config.FlushInterval = opts.Duration / time.Duration(4)
// 			config.FlushBytes = 0
// 		}

// 		client, err = statsd.NewClientWithConfig(config)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		defer client.Close()
// 	}

// 	var stat func(stat string, value int64, rate float32) error
// 	switch opts.StatType {
// 	case "count":
// 		stat = func(stat string, value int64, rate float32) error {
// 			return client.Inc(stat, value, rate)
// 		}
// 	case "gauge":
// 		stat = func(stat string, value int64, rate float32) error {
// 			return client.Gauge(stat, value, rate)
// 		}
// 	case "timing":
// 		stat = func(stat string, value int64, rate float32) error {
// 			return client.Timing(stat, value, rate)
// 		}
// 	default:
// 		log.Fatal("Unsupported state type")
// 	}

// 	pertick := opts.Volume / int(opts.Duration.Seconds()) / 10
// 	// add some extra time, because the first tick takes a while
// 	ender := time.After(opts.Duration + 100*time.Millisecond)
// 	c := time.Tick(time.Second / 10)
// 	count := 0
// 	for {
// 		select {
// 		case <-c:
// 			for x := 0; x < pertick; x++ {
// 				err := stat(opts.Name, opts.StatValue, opts.Rate)
// 				if err != nil {
// 					log.Printf("Got Error: %+v\n", err)
// 					break
// 				}
// 				count++
// 			}
// 		case <-ender:
// 			log.Printf("%d events called\n", count)
// 			return
// 		}
// 	}
// }

func ovsDBReader() {
	reader := ovsdbreader.CreateNewOVSDBReader(&config.InitConfig.OvsDBConf)
	reader.ConnectDB()
	res := reader.ReadOVSDB()
	for i, elem := range res.Rows {
		fmt.Printf("\n Row : %d ...", i)
		for _, data := range elem.DataSet {
			fmt.Printf("\n %+v", data)
		}
	}
	reader.CloseDBConn()

}

func statsDWriter() {
	writer:= statsdwriter.CreateStatsDWriter(&config.InitConfig.StatsDConf)
	writer.Connect()
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
	ovsDBReader()
	statsDWriter()
 }

