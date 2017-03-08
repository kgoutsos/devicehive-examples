# DeviceHive Example Device
------------
A simple example of a device communicating with a DeviceHive server and sending notifications to it. Created mainly as a quick-and-easy template to be used while testing different aspects of the DeviceHive platform.

## How to run

1. Copy or rename the example configuration file to ``deviceconf.yml`` and edit it according to your setup details. You can use a different name for the file as long as you update ``dh-cloud-run.sh`` accordingly.
2. Start the DeviceHive cloud client by running ``dh-cloud-run.sh``. The script will check for the required Go packages and install them as needed. If your DBus installation is not configured correctly, you might get a permission error. In that case see the DBus configuration instructions below.
3. Start the device by running go run ``dh-example-device.go``

## DBus configuration

To enable local DBus method calls you should replace the following section in ``/etc/dbus-1/system.conf``
```
<deny own="*"/>
<deny send_type="method_call"/>
```

with

```
<allow own="*"/>
<allow send_type="method_call"/>
```

and then restart the DBus service with ``service dbus restart`` or reboot your system.

It has been reported that in newer Ubuntu distributions the configuration of the DBus permissions has been overhauled and the previous step will not work. However, on Fedora up to version 23 and Ubuntu up to 14.04 (at least) the above instructions still apply.

## More information

The author of this project is in no way affiliated with the development team of DeviceHive apart from merely contributing from time to time via Github.

For more information about [DeviceHive](http://www.devicehive.com) please visit the platform's website as well as its [organisation on Github](http://github.com/devicehive).
