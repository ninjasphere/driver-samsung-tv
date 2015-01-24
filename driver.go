package main

import (
	"io/ioutil"
	"strings"

	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/logger"
	"github.com/ninjasphere/go-ninja/support"
)

var info = ninja.LoadModuleInfo("./package.json")
var log = logger.GetLogger(info.Name)

type Driver struct {
	support.DriverSupport
}

func NewDriver() (*Driver, error) {

	driver := &Driver{}

	err := driver.Init(info)
	if err != nil {
		log.Fatalf("Failed to initialize driver: %s", err)
	}

	err = driver.Export(driver)
	if err != nil {
		log.Fatalf("Failed to export driver: %s", err)
	}

	return driver, nil
}

func (d *Driver) Start(_ interface{}) error {
	log.Infof("Driver Starting")

	ipB, err := ioutil.ReadFile("./tv.txt")
	if err != nil {
		log.Fatalf("Failed to load IP address from tv.txt: %s", err)
	}

	ip := strings.TrimSpace(string(ipB))

	log.Infof("Temporary: Using hardcoded IP address: %s", ip)

	p, err := NewMediaPlayer(d, d.Conn, ip)

	//p.applyPlayPause(false)

	return nil
}
