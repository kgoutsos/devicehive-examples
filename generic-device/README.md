# DeviceHive Example Device

A simple example of a device communicating with a DeviceHive server and sending notifications to it. Created mainly as a quick-and-easy template to be used while testing different aspects of the DeviceHive platform.

## How to run

1. Edit the configuration file `cloud/device-conf.yml` to match your setup.
2. Start the DeviceHive cloud client by running `cloud/start.sh`. See cloud/README.md for details and troubleshooting.
3. Start the device with `go run generic-device.go`
