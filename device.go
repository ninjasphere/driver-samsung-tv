package main

import (
	"time"

	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/config"
	"github.com/ninjasphere/go-ninja/devices"
	"github.com/ninjasphere/go-ninja/model"
	"github.com/ninjasphere/go-samsung-tv"
)

type Device struct {
	devices.MediaPlayerDevice
	tv *samsung.TV
}

func (d *Device) updateHost(host string) {
	d.tv.Host = host
}

func newDevice(driver ninja.Driver, conn *ninja.Connection, cfg *TVConfig) (*Device, error) {

	player, err := devices.CreateMediaPlayerDevice(driver, &model.Device{
		NaturalID:     cfg.ID,
		NaturalIDType: "samsung-tv",
		Name:          &cfg.Name,
		Signatures: &map[string]string{
			"ninja:manufacturer": "Samsung",
			"ninja:productName":  "Smart TV",
			"ninja:thingType":    "mediaplayer",
			"ip:mac":             cfg.ID,
		},
	}, conn)

	if err != nil {
		return nil, err
	}

	samsung.EnableLogging = true

	tv := samsung.TV{
		Host:            cfg.Host,
		ApplicationID:   config.MustString("userId"),
		ApplicationName: "Ninja Sphere         ",
	}

	// Volume Channel
	player.ApplyVolumeUp = func() error {
		return tv.SendCommand("KEY_VOLUP")
	}

	player.ApplyVolumeDown = func() error {
		return tv.SendCommand("KEY_VOLDOWN")
	}

	player.ApplyToggleMuted = func() error {
		return tv.SendCommand("KEY_MUTE")
	}

	if err := player.EnableVolumeChannel(false); err != nil {
		player.Log().Fatalf("Failed to enable volume channel: %s", err)
	}

	// Media Control Channel
	player.ApplyPlayPause = func(play bool) error {
		if play {
			return tv.SendCommand("KEY_PLAY")
		}

		return tv.SendCommand("KEY_PAUSE")
	}

	if err := player.EnableControlChannel([]string{}); err != nil {
		player.Log().Fatalf("Failed to enable control channel: %s", err)
	}

	// On-off Channel
	player.ApplyOff = func() error {
		return tv.SendCommand("KEY_POWEROFF")
	}

	if err := player.EnableOnOffChannel("state"); err != nil {
		player.Log().Fatalf("Failed to enable control channel: %s", err)
	}

	go func() {

		// Continuous updates as TV goes online and offline
		for online := range tv.OnlineState(time.Second * 5) {
			player.UpdateOnOffState(online)
		}
	}()

	return &Device{*player, &tv}, nil
}
