# Fingy - Device gateway
This is the Fingy component facing the devices (i.e. the WebSocket endpoint that devices implementing the Fingy client contact). It manages the connection between the devices and Fingy, and dispatches received events.

**NOTE: this is still a pre-alpha version. Everything can change and documentation is close to non-existing.**

## How to use?
The gateway is far from being finished. It does not connect to any database, automatically registers incoming devices without any check, and has an hardcoded list of destination services. To use, just run the executable.

## Remaining work
* Package the gateway better (configuration, documentation).
* Save device & service registrations to a database.
* Implement lookup device lookup service (to know to which gateway a specific device is connected to).
* Local buffer of messages on device/service unavailability, with timeout management.
* Registration/unregistration of devices endpoints (should be another service).
* Built-in aliveness check.
