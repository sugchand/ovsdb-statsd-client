// Package ovsdb-reader : Client to interact with ovsdb
//
//  Copyright (c) 2020 Sugesh Chandran
//
package ovsdbreader

import (
	"time"
	"fmt"
	"context"
	"log"
	"ovsdb-statsd-client/config"
	"ovsdb-statsd-client/pkg/errors"
	"github.com/digitalocean/go-openvswitch/ovsdb"
)

const (
	DEFAULT_OVSDB_OP_TIMEOUT = 3*time.Second
)

type tableConf struct {
	name string
	cols []*DBColConf
}

type OVSDBReader struct {
	networkType string
	address string
	dbName string
	tables []*tableConf
	conn *ovsdb.Client
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
	reader.conn, err = ovsdb.Dial(
		reader.networkType, reader.address)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
		return err
	}
	return errors.ErrNil
}

func (reader *OVSDBReader)CloseDBConn() {
	reader.conn.Close()
}

func (reader *OVSDBReader)getNextTable(tableconf **tableConf, tableIdx *int) bool {
	if len(reader.tables) > *tableIdx {
		*tableconf = reader.tables[*tableIdx]
		*tableIdx = *tableIdx + 1
		return true
	}
	return false
}

func (reader *OVSDBReader)rowIsInSameTable(prevRow ovsdb.Row, currRow ovsdb.Row) bool {
	if len(prevRow) != len(currRow) {
		return false
	}
	for key, _ := range prevRow {
		if _, found := currRow[key]; found == false {
			// key in both row are mismatch
			return false
		}
	}
	return true
}

func (reader *OVSDBReader)ReadOVSDB()*DBReport {

	var report DBReport
	// create list of transaction for each table
	txops := make([]ovsdb.TransactOp, 0)
	for _, tabl := range reader.tables {
		transact := &ovsdb.Select{
			Table : tabl.name,
		}
		txops = append(txops, *transact)
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
	var currTableConf *tableConf
	var currTableIdx int
	var prevRow ovsdb.Row
loop:
	for _, row := range rows {
		// in this list of rows, data is from all
		// the different tables.Now lets find out
		// which table a row belongs to.
		// We follow a bit naive approach here by just
		// looking at the order of table input vs the order
		// of rows returned, assuming ovsdb client return the results
		// in the order of input query.
		if reader.rowIsInSameTable(prevRow, row) == false {
			if reader.getNextTable(&currTableConf, &currTableIdx) == false {
				// failed to get next table, cannot process further.
				break loop
			}
		}

		reportRow := new(ReportRow)
		reportRow.DataSet = make([]*ReportCol, 0)
		for _, col := range currTableConf.cols {
			if colVal, ok := row[col.ColName];ok {
				reportCol := &ReportCol {
					ColName : currTableConf.name + "-" + col.ColName,
					Data : colVal,
					ReportType : col.reportType,
				}
				reportRow.DataSet = append(reportRow.DataSet, reportCol)
			}
		}
		report.Rows = append(report.Rows, reportRow)
		prevRow = row
	}
	reader.report = &report
	return reader.report
}

func (reader *OVSDBReader)DisplayReport() {
	for i, row := range reader.report.Rows {
		fmt.Printf("\n %dth row :", i)
		colDataStr := "\tData: "
		for _, data := range row.DataSet {
			colDataStr = fmt.Sprintf("%s %s(%d) : %+v", colDataStr, data.ColName, data.ReportType,
		                                data.Data)
		}
		fmt.Printf("\n %s", colDataStr)
	}
}

func CreateNewOVSDBReader(conf *config.OVSDBConfig)*OVSDBReader {
	reader := new(OVSDBReader)
	reader.networkType = conf.Network
	reader.address = conf.Address
	reader.dbName = conf.DB
	// set up table conf
	reader.tables = make([]*tableConf, 0)
	for _, table := range conf.Tables {
		tabl := &tableConf {
			name : table.Name,
		}
		for _, element := range table.Cols {
			col := &DBColConf {
			ColName : element.Name,
			reportType : element.Type,
			}
			tabl.cols = append(tabl.cols, col)
		}
		reader.tables = append(reader.tables, tabl)
	}
	log.Printf("Initialized a new OVSDB reader with parameters %+v", *reader)
	return reader
}
