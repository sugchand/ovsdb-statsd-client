// Package config : Initial configuration for the ovsdb-statsd client
// refer config.yaml for more details.
//
//  Copyright (c) 2020 Sugesh Chandran
//
package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"gopkg.in/yaml.v2"
	"ovsdb-statsd-client/pkg/errors"
)

type ReportValueType uint16
const (
	// 0 : undefined
    // 1 : Counter
    // 2 : Gauge
    // 3 : Timer
	Undefined ReportValueType = iota
	Counter
	Gauge
	Timer
)

type StatsDConfig struct {
	Host string `yaml:"Host"`
	Port uint16 `yaml:"Port"`
	PollInterval uint16 `yaml:"PollInterval"`
	Prefix string `yaml:2Prefix"`
	SampleRate uint16 `yaml:"SampleRate"`
}

type DBColumns struct {
	Name string `yaml:"Name"`
	Type ReportValueType `yaml:"Type"`
}

type OVSDBConfig struct {
	Network string `yaml:"Network"`
	Address string `yaml:"Address"`
	DB string `yaml:"DB"`
	Table string `yaml:"Table"`
	Cols []DBColumns `yaml:"Columns"`
}

// StartupConfig : YAML configuration struct for the agent
type StartupConfig struct {
	OvsDBConf OVSDBConfig `yaml:"OvsDBConfig"`
	StatsDConf StatsDConfig `yaml:"StatsDConfig"`
}

// InitConfig : Global variable to hold the startup configuration.
var InitConfig StartupConfig

// ParseYAMLConfig : Must be called at the startup before any other agent operation.
func ParseYAMLConfig(yamlFilePath string) error {
	var configFile string
	var err error
	if yamlFilePath == "" {
		configFile, err = os.Executable()
		if err != nil {
			panic("Failed to get YAML configuration filepath err : " + err.Error())
		}
		configFile, err = filepath.Abs(filepath.Dir(filepath.Dir(configFile)))
		if err != nil {
			panic("Failed to get YAML file absolute directory path")
		}
		configFile = configFile + "/config/config.yaml"
	} else {
		configFile = yamlFilePath
	}
	yamlConfig, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Error in reading YAML file %s err: %s", configFile, err)
		return err
	}
	err = yaml.Unmarshal(yamlConfig, &InitConfig)
	if err != nil {
		fmt.Printf("Error in unmarshal YAML config file %s, err : %s ",
			yamlConfig, err)
		return err
	}
	return errors.ErrNil
}