package main

import (
	"fmt"

	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/devices"
	"github.com/ninjasphere/go-ninja/model"
)

type MediaPlayer struct {
	player     *devices.MediaPlayerDevice
	ip         string
	lastVolume float64
}

func NewMediaPlayer(driver ninja.Driver, conn *ninja.Connection, ip string) (*MediaPlayer, error) {
	name := fmt.Sprintf("Samsung TV: %s", ip)

	player, err := devices.CreateMediaPlayerDevice(driver, &model.Device{
		NaturalID:     ip,
		NaturalIDType: "samsung-tv",
		Name:          &name,
		Signatures: &map[string]string{
			"ninja:manufacturer": "Samsung",
			"ninja:productName":  "Smart TV",
			"ninja:thingType":    "mediaplayer",
		},
	}, conn)

	if err != nil {
		return nil, err
	}

	device := &MediaPlayer{
		player: player,
		ip:     ip,
	}

	player.ApplyVolumeUp = device.applyVolumeUp
	player.ApplyVolumeDown = device.applyVolumeDown
	player.ApplyToggleMuted = device.applyToggleMuted
	if err := player.EnableVolumeChannel(false); err != nil {
		player.Log().Fatalf("Failed to enable volume channel: %s", err)
	}

	player.ApplyPlayPause = device.applyPlayPause
	if err := player.EnableControlChannel([]string{}); err != nil {
		player.Log().Fatalf("Failed to enable control channel: %s", err)
	}

	return device, nil
}

func (d *MediaPlayer) applyPlayPause(play bool) error {
	if play {
		go sendCommand(d.ip, "KEY_PLAY")
	} else {
		go sendCommand(d.ip, "KEY_PAUSE")
	}
	return nil
}

func (d *MediaPlayer) applyToggleMuted() error {
	d.player.Log().Infof("applyToggleMuted called")

	go sendCommand(d.ip, "KEY_MUTE")
	return nil
}

func (d *MediaPlayer) applyVolumeUp() error {
	d.player.Log().Infof("applyVolumeUp called")

	go sendCommand(d.ip, "KEY_VOLUP")
	return nil
}

func (d *MediaPlayer) applyVolumeDown() error {
	d.player.Log().Infof("applyVolumeDown called")

	go sendCommand(d.ip, "KEY_VOLDOWN")
	return nil
}
