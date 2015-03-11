# Ninja Sphere - Samsung TV Driver


[![godoc](http://img.shields.io/badge/godoc-Reference-blue.svg)](https://godoc.org/github.com/ninjasphere/driver-samsung-tv)
[![MIT License](https://img.shields.io/badge/license-MIT-yellow.svg)](LICENSE)
[![Ninja Sphere](https://img.shields.io/badge/built%20by-ninja%20blocks-lightgrey.svg)](http://ninjablocks.com)
[![Ninja Sphere](https://img.shields.io/badge/works%20with-ninja%20sphere-8f72e3.svg)](http://ninjablocks.com)

---


### Introduction
This is a driver for Samsung Smart TVs, allowing them to be used as part of Ninja Sphere.

It is basically a port of https://github.com/natalan/samsung-remote

** VERY EARLY - Set the TV ip in tv.txt **

### Supported Sphere Protocols

| Name | URI | Supported Events | Supported Methods |
| ------ | ------------- | ---- | ----------- |
| volume | [http://schema.ninjablocks.com/protocol/volume](https://github.com/ninjasphere/schemas/blob/master/protocol/volume.json) | | set, volumeUp, volumeDown, mute, unmute, toggleMute |
| media-control | [http://schema.ninjablocks.com/protocol/media-control](https://github.com/ninjasphere/schemas/blob/master/protocol/media-control.json) | play, pause  | |

#### Can't Do
* There is currently no way to get state back from the television.

### Requirements

* Go 1.3

### Building

This project can be built with `go build`, but a makefile is also provided.

### Running

`DEBUG=* ./driver-samsung-tv`

### Installing the driver

This driver isn't currently officially supported by Ninja Blocks and so is not installed by default.

To install the driver in a manner that survives reboots of the Ninja Sphere, connect to the sphere with SSH or a mini-USB cable and edit and then paste the following command into a terminal:

	sudo apt-get update -y &&
	sudo apt-get install -y ninja-samsung-tv &&
	export SAMSUNG_IP={edit-this-to-add-your-samsung-tv-ip} # edit this with the IP address of your Samsung Smart TV

Then, paste the following command into the terminal:

	sudo mkdir -p /data/sphere/user-autostart/drivers &&
	sudo ln -sf /opt/ninjablocks/drivers/driver-samsung-tv /data/sphere/user-autostart/drivers &&
	echo $SAMSUNG_IP | sudo with-rw tee /opt/ninjablocks/drivers/driver-samsung-tv/tv.txt

Then, run the following command to reboot the sphere.

	sudo reboot # to reboot the sphere

Then, open the Ninja Sphere phone app and go to the "Add Things" page to add your TV to a room.

### Disabling the driver

If you would like to disable the driver running, connect to the sphere with SSH or a mini-USB cable and run:

	nservice -q driver-samsung-tv stop &&
	sudo rm /data/sphere/user-autostart/drivers/driver-samsung-tv

### Uninstalling the driver

	sudo apt-get remove ninja-samsung-tv

### Options

The usual Ninja Sphere configuration and parameters apply, but these are the most useful during development.

* `--autostart` - Doesn't wait to be started by Ninja Sphere
* `--mqtt.host=HOST` - Override default mqtt host
* `--mqtt.port=PORT` - Override default mqtt host

### More Information

More information can be found on the [project site](http://github.com/ninjasphere/driver-samsung-tv) or by visiting the Ninja Blocks [forums](https://discuss.ninjablocks.com).

### Contributing Changes

To contribute code changes to the project, please clone the repository and submit a pull-request ([What does that mean?](https://help.github.com/articles/using-pull-requests/)).

### License
This project is licensed under the MIT license, a copy of which can be found in the [LICENSE](LICENSE) file.

### Copyright
This work is Copyright (c) 2014-2015 - Ninja Blocks Inc.
