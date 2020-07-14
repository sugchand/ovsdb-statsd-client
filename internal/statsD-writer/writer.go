// Package statsD-writer : Client to interact with ovsdb
//
//  Copyright (c) 2020 Sugesh Chandran
//
package statsdwriter

import (
	"time"
	"context"
	"log"
	"ovsdb-statsd-client/config"
	"ovsdb-statsd-client/pkg/errors"
	"github.com/cactus/go-statsd-client/statsd"
)

type statsDWriter struct {
	address string // network address in host:port format
	pollInterval uint16
	reportPrefix string /// prefix used when reporting data
	


}


func CreateStatsDWriter(conf *config.StatsDConf) *statsDWriter {
	writer = new(statsDWriter)

}