package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mostlygeek/arp"
	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/logger"
	"github.com/ninjasphere/go-ninja/model"
	"github.com/ninjasphere/go-ninja/support"
	"github.com/ninjasphere/go-samsung-tv"
)

var info = ninja.LoadModuleInfo("./package.json")
var log = logger.GetLogger(info.Name)

type Driver struct {
	support.DriverSupport
	config  Config
	devices map[string]*Device
}

type Config struct {
	TVs map[string]*TVConfig
}

func (c *Config) get(id string) *TVConfig {
	for _, tv := range c.TVs {
		if tv.ID == id {
			return tv
		}
	}
	return nil
}

type TVConfig struct {
	ID   string
	Name string
	Host string
}

func NewDriver() (*Driver, error) {

	driver := &Driver{
		devices: make(map[string]*Device),
	}

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

func (d *Driver) deleteTV(id string) error {
	delete(d.config.TVs, id)

	err := d.SendEvent("config", &d.config)

	// TODO: Can't unexport devices at the moment, so we should restart the driver...
	go func() {
		time.Sleep(time.Second * 2)
		os.Exit(0)
	}()

	return err
}

func (d *Driver) saveTV(tv TVConfig) error {

	if !(&samsung.TV{Host: tv.Host}).Online(time.Second * 5) {
		return fmt.Errorf("Could not connect to TV. Is it online?")
	}

	mac, err := getMACAddress(tv.Host, time.Second*10)

	if err != nil {
		return fmt.Errorf("Failed to get mac address for TV. Is it online?")
	}

	existing := d.config.get(mac)

	if existing != nil {
		existing.Host = tv.Host
		existing.Name = tv.Name
		device, ok := d.devices[mac]
		if ok {
			device.tv.Host = tv.Host
		}
	} else {
		tv.ID = mac
		d.config.TVs[mac] = &tv

		go d.createTVDevice(&tv)
	}

	return d.SendEvent("config", d.config)
}

func (d *Driver) Start(config *Config) error {
	log.Infof("Driver Starting with config %+v", config)

	if config.TVs == nil {
		config.TVs = make(map[string]*TVConfig)
	}

	d.config = *config

	for _, cfg := range config.TVs {
		d.createTVDevice(cfg)
	}

	d.Conn.MustExportService(&configService{d}, "$driver/"+info.ID+"/configure", &model.ServiceAnnouncement{
		Schema: "/protocol/configuration",
	})

	return nil
}

func (d *Driver) createTVDevice(cfg *TVConfig) {

	device, err := newDevice(d, d.Conn, cfg)

	if err != nil {
		log.Fatalf("Failed to create new Samsung TV device host:%s id:%s name:%s : %s", cfg.Host, cfg.ID, cfg.Name, err)
	}

	d.devices[cfg.ID] = device
}

func getMACAddress(host string, timeout time.Duration) (string, error) {

	timedOut := false
	success := make(chan string, 1)

	go func() {
		for {
			if timedOut {
				break
			}
			id := arp.Search(host)
			if id != "" && id != "(incomplete)" {
				success <- id
			}

			time.Sleep(time.Millisecond * 500)
		}
	}()

	select {
	case mac := <-success:
		return mac, nil
	case <-time.After(timeout):
		timedOut = true
		return "", fmt.Errorf("Timed out searching for MAC address")
	}
}
