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
	"reflect"
	"ovsdb-statsd-client/config"
	"ovsdb-statsd-client/pkg/errors"
	"ovsdb-statsd-client/internal/ovsdb-reader"
	"github.com/cactus/go-statsd-client/statsd"
)

type SWriter struct {
	address string // network address in host:port format
	flushInterval uint16
	reportPrefix string /// prefix used when reporting data
	sampleRate float32
	prefix string
	conn statsd.Statter

}

func (writer *SWriter)Connect() error {
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

func (writer *SWriter)Disconnect() {
	writer.conn.Close()
}


func (writer *SWriter)getRowName(row *ovsdbreader.ReportRow) string {
	rowName := ""
	for _, col := range row.DataSet {
		if col.ReportType == config.TagName {
			name, ok := writer.stringValue(col.Data)
			if !ok {
				continue
			}
			rowName = rowName + name
		}
	}
	return rowName
}

func (writer *SWriter)isValidColData(data interface{}) bool{
	switch v:= reflect.ValueOf(data); v.Kind() {
		case reflect.Map, reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, 
		   reflect.Uint,reflect.Uint8,reflect.Uint16,reflect.Uint32,reflect.Uint64, reflect.Int64:
			return true
		default:
			log.Printf("Invalid column type %+v", v.Kind())
			return false
	}
}

func (writer *SWriter)stringValue(data interface{}) (string, bool) {
	switch v:= reflect.ValueOf(data); v.Kind() {
	case reflect.String:
		res, ok := data.(string)
		if !ok {
			return "", false
		}
		return res, true
	default:
		return "", false
	}
}

func (writer *SWriter)numericValue(data interface{}) (int64, bool) {
	switch v:= reflect.ValueOf(data); v.Kind() {
	case reflect.Int:
		res, ok := data.(int)
		if !ok {
			return 0 , false
		}
		return int64(res), true
	case reflect.Int8:
		res, ok := data.(int8)
		if !ok {
			return 0 , false
		}
		return int64(res), true
	case reflect.Int16:
		res, ok := data.(int16)
		if !ok {
			return 0 , false
		}
		return int64(res), true
	case reflect.Int32:
		res, ok := data.(int32)
		if !ok {
			return 0 , false
		}
		return int64(res), true
	case reflect.Uint:
		res, ok := data.(int64)
		if !ok {
			return 0 , false
		}
		return int64(res), true
	case reflect.Uint8:
		res, ok := data.(uint8)
		if !ok {
			return 0 , false
		}
		return int64(res), true
	case reflect.Uint16:
		res, ok := data.(uint16)
		if !ok {
			return 0 , false
		}
		return int64(res), true
	case reflect.Uint32:
		res, ok := data.(uint32)
		if !ok {
			return 0 , false
		}
		return int64(res), true
	case reflect.Uint64:
		res, ok := data.(uint64)
		if !ok {
			return 0 , false
		}
		return int64(res), true
	case reflect.Int64:
		res, ok := data.(int64)
		if !ok {
			return 0 , false
		}
		return int64(res), true
	default:
		return 0, false
	}
}

func (writer *SWriter)ColDataIsMap(data interface{}) bool{
	switch v:= reflect.ValueOf(data); v.Kind() {
	case reflect.Map:
		return true
	default:
		return false
	}
}

func (writer *SWriter)writeCol(name string, value int64, rType config.ReportValueType) {
	switch rType {
	case config.Counter:
		writer.conn.Inc(name, value, writer.sampleRate)
	case config.Gauge:
		writer.conn.Gauge(name, value, writer.sampleRate)
	case config.Timer:
		writer.conn.Timing(name, value, writer.sampleRate)
	default:
		log.Printf("Invalid report type, %d, Exiting..", rType)
	}

}
func (writer *SWriter)processCol(name string, col *ovsdbreader.ReportCol) {
	statName := name + col.ColName
	if writer.isValidColData(col.Data) == false {
		return
	}
	if writer.ColDataIsMap(col.Data) == true {
		// process a map
		colMap, _ := col.Data.(map[interface{}]interface{})
		for key, val := range colMap {
			name, ok := writer.stringValue(key)
			if !ok {
				continue
			}
			name = statName + name
			dataVal, ok := writer.numericValue(val)
			if !ok {
				continue
			}
			writer.writeCol(name, dataVal, col.ReportType)

		}

	}
	dataVal, ok := writer.numericValue(col.Data)
	if !ok {
		return
	}
	writer.writeCol(statName, dataVal, col.ReportType)
}

func (writer *SWriter)processRow(row *ovsdbreader.ReportRow) {
	rowName := writer.getRowName(row)
	for _, col := range row.DataSet {
		writer.processCol(rowName, col)
	}
}

func (writer *SWriter)Write(report *ovsdbreader.DBReport) {
	for _, row := range report.Rows {
		// process each row in the report.
		writer.processRow(row)
	}
}

func CreateSWriter(conf *config.StatsDConfig) *SWriter {
	writer := new(SWriter)
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
