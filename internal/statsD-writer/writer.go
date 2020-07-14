// Package statsD-writer : Client to interact with ovsdb
//
//  Copyright (c) 2020 Sugesh Chandran
//
package statsdwriter

import (
	"time"
	//"context"
	"log"
	"strconv"
	"ovsdb-statsd-client/config"
	"ovsdb-statsd-client/pkg/errors"
	"github.com/cactus/go-statsd-client/statsd"
)

type statsDWriter struct {
	address string // network address in host:port format
	flushInterval uint16
	reportPrefix string /// prefix used when reporting data
	sampleRate uint16
	prefix string
	conn statsd.Statter

}

func (writer *statsDWriter)Connect() error {
	var err error

	flush_interval := time.Duration(writer.flushInterval)* time.Nanosecond
	config := &statsd.ClientConfig{
	Address: writer.address,
	Prefix: writer.reportPrefix,
	ResInterval: time.Duration(0),
	UseBuffered : true,
	FlushInterval : flush_interval,
	FlushBytes : 0,
	}
	writer.conn,err = statsd.NewClientWithConfig(config)
	if err != nil {
		log.Fatalf("Failed to connect to statsd server %s", err)
		return err
	}
	return errors.ErrNil
}

func (writer *statsDWriter)Disconnect() {
	writer.conn.Close()
}

func CreateStatsDWriter(conf *config.StatsDConfig) *statsDWriter {
	writer := new(statsDWriter)
	var host, port string
	if conf.Host == "" {
		host = config.DEFAULT_STATSD_SERVER_IP
	} else {
		host = conf.Host
	}
	if conf.Port == 0 {
		port = strconv.Itoa(config.DEFAULT_STATSD_SERVER_PORT)
	} else {
		port = strconv.Itoa(int(conf.Port))
	}
	writer.address = host + ":" + port
	writer.flushInterval = config.DEFAULT_STATSD_FLUSH_INTERVAL
	if conf.FlushInterval != 0 {
		writer.flushInterval = conf.FlushInterval
	}
	writer.sampleRate = config.DEFAULT_STATSD_SAMPLE_RATE
	if conf.SampleRate != 0 {
		writer.sampleRate = conf.SampleRate
	}
	writer.prefix = config.DEFAULT_STATSD_PREFIX
	if conf.Prefix != "" {
		writer.prefix = conf.Prefix
	}
	log.Printf("Initialized a new write with parameters %+v", *writer)
	return writer
}
