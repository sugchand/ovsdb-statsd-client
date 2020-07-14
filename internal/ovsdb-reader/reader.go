// Package ovsdb-reader : Client to interact with ovsdb
//
//  Copyright (c) 2020 Sugesh Chandran
//
package ovsdbreader

import (
	"time"
	"context"
	"log"
	"ovsdb-statsd-client/config"
	"ovsdb-statsd-client/pkg/errors"
	"github.com/digitalocean/go-openvswitch/ovsdb"
)

const (
	DEFAULT_OVSDB_OP_TIMEOUT = 3*time.Second
)

type OVSDBReader struct {
	networkType string
	address string
	dbName string
	tableName string
	conn *ovsdb.Client
	cols []*DBColConf
	report *DBReport
	
}

type DBColConf struct {
	ColName string
	reportType config.ReportValueType
}

type DBReport struct {
	Rows []*ReportRow
}

type ReportRow struct {
	DataSet []*ReportCol
}

type ReportCol struct {
	ColName string
	Data interface{}
	ReportType config.ReportValueType
}

func (reader *OVSDBReader)ConnectDB() error {
	// Dial an OVSDB connection and create a *ovsdb.Client.
	var err error
	reader.conn, err = ovsdb.Dial("unix", "/usr/local/var/run/openvswitch/db.sock")
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
		return err
	}
	return errors.ErrNil
}

func (reader *OVSDBReader)CloseDBConn() {
	reader.conn.Close()
}

func (reader *OVSDBReader)ReadOVSDB()*DBReport {

	var report DBReport
	selectTransact := ovsdb.Select{
		Table : reader.tableName,
	}

	txops := []ovsdb.TransactOp {
		selectTransact,
	}
	// read table from ovsdb
	tblCtx, tblCancel := context.WithTimeout(context.Background(),
							DEFAULT_OVSDB_OP_TIMEOUT)
	defer tblCancel()
	rows, err := reader.conn.Transact(tblCtx, reader.dbName, txops)
	if err != nil {
		log.Fatalf("Failed to read from db. err : %s", err)
		return nil
	}
	if len(rows) == 0 {
		return nil
	}
	report.Rows = make([]*ReportRow, 0)
	for _, row := range rows {
		reportRow := new(ReportRow)
		reportRow.DataSet = make([]*ReportCol, 0)
		for _, col := range reader.cols {
			if colVal, ok := row[col.ColName];ok {
				reportCol := &ReportCol {
					ColName : col.ColName,
					Data : colVal,
					ReportType : col.reportType,
				}
				reportRow.DataSet = append(reportRow.DataSet, reportCol)
			}
		}
		report.Rows = append(report.Rows, reportRow)
	}
	reader.report = &report
	return reader.report
}

func CreateNewOVSDBReader(conf *config.OVSDBConfig)*OVSDBReader {
	reader := new(OVSDBReader)
	reader.networkType = conf.Network
	reader.address = conf.Address
	reader.dbName = conf.DB
	reader.tableName = conf.Table
	reader.cols = make([]*DBColConf, 0)
	for _, element := range conf.Cols {
		col := &DBColConf {
			ColName : element.Name,
			reportType : element.Type,
		}
		reader.cols = append(reader.cols, col)
	}
	log.Printf("Initialized a new OVSDB reader with parameters %+v", *reader)
	return reader
}