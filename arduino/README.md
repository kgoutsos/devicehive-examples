# DeviceHive Example Device

A simple example of the DeviceHive server controlling an Arduino through an intermediate router running a bridge application.

## How to run

1. Start the DeviceHive cloud client as described in [cloud-client/README.md](../cloud-client/README.md)
2. Load [serial-led-control.ino](serial-led-control.ino) to your Arduino using your preferred method.
3. Start the device with `go run bridge.go`

The [serial-led-control.ino](serial-led-control.ino) sketch on Arduino will keep reading bytes from the serial port until it gets an "1" or a "0". When that happens, the LED (or anything else connected to pin 13) is toggled accordingly and a response is written to the serial port.

The Go code in [bridge.go](bridge.go) listens to the DeviceHive server for On/Off commands and sends the appropriate bytes to the Arduino through the serial port. Depending on your setup you may have to change the serial port to match the one to which your Arduino is connected.
