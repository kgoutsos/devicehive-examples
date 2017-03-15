# DeviceHive Cloud Client

## How to run

1. Edit the configuration file `cloud/device-conf.yml` to match your setup.
2. Start the DeviceHive cloud client by running `cloud/start.sh`.

The script will check for the required Go packages and install them as needed.

In case you get a permission error, make sure DBus is configured as described below.

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
