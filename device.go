package main

import (
	"fmt"

	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/channels"
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

	player.ApplyVolume = device.applyVolume // TODO: Remove me!

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

// TODO: Remove me! Just needed till led controller is better
func (d *MediaPlayer) applyVolume(state *channels.VolumeState) error {
	d.player.Log().Infof("applyVolume called, volume %v", state)

	if *state.Level > 0 && d.lastVolume == 0 {
		return sendCommand(d.ip, "KEY_MUTE")
	}

	if *state.Level == 0 && d.lastVolume > 0 {
		return sendCommand(d.ip, "KEY_MUTE")
	}

	// Do nothing for now
	return nil
}

func (d *MediaPlayer) applyPlayPause(play bool) error {
	if play {
		return sendCommand(d.ip, "KEY_PLAY")
	}
	return sendCommand(d.ip, "KEY_PAUSE")
}

func (d *MediaPlayer) applyToggleMuted() error {
	d.player.Log().Infof("applyToggleMuted called")

	return sendCommand(d.ip, "KEY_MUTE")
}

func (d *MediaPlayer) applyVolumeUp() error {
	d.player.Log().Infof("applyVolumeUp called")

	return sendCommand(d.ip, "KEY_VOLUP")
}

func (d *MediaPlayer) applyVolumeDown() error {
	d.player.Log().Infof("applyVolumeDown called")

	return sendCommand(d.ip, "KEY_VOLDOWN")
}
