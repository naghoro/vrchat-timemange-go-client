# vrchat-timemange-go-client
time manage client for VRChat, send OSC once an hour

## Overview

When client start, client send OSC to 127.0.0.1:9000.

It send OSC once a hour.
When client start, it send firstly and then send second time after one hour.

If the current time is 15:15, message will be sent as follows:

```
/time/hour/14 0
/time/hour/15 1
```

It sends twice. The first is off for one hour ago, the second is on for current hour.

## Prerequsites

golang version 1.3 or higher.

## Usage

### build

exec command line

```
vrchat-timemange-go-client$ make
```

### start client

send OSC once an hour

```
./appmanager
```

You can change second of hour for testing.
Below, send OSC once a second.

```
./appmanager -sec=1s
```

