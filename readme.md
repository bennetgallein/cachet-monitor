cachet-monitor

a fork of the cachet-monitor project originally developed by CastawayLabs but abandoned. Extended with new features

## Features

:heavy_check_mark: Interval based checking of predefined Resources  
:heavy_check_mark: Posts monitor lag to cachet graphs  
:heavy_check_mark: Creates & Resolves Incidents  
:heavy_check_mark: Updates Component to Partial Outage  
:heavy_check_mark: Updates Component to Major Outage if already in Partial Outage (works with distributed monitors)  
:heavy_check_mark: Can be run on multiple servers and geo regions  
:heavy_check_mark: HTTP Checks (body/status code)  
:heavy_check_mark: DNS Checks  
:heavy_check_mark: TCP Checks  
:heavy_check_mark: ICMP Checks

## Quick Start

Configuration can be done in either yaml or json format. An example JSON-File would look something like this:

```json
{
  "api": {
    "url": "https://demo.cachethq.io/api/v1",
    "token": "9yMHsdioQosnyVK4iCVR",
    "insecure": false
  },
  "date_format": "02/01/2006 15:04:05 MST",
  "monitors": [
    {
      "active": false,
      "name": "google",
      "target": "https://google.com",
      "strict": true,
      "method": "POST",
      "component_id": 1,
      "metric_id": 4,
      "template": {
        "investigating": {
          "subject": "{{ .Monitor.Name }} - {{ .SystemName }}",
          "message": "{{ .Monitor.Name }} check **failed** (server time: {{ .now }})\n\n{{ .FailReason }}"
        },
        "fixed": {
          "subject": "I HAVE BEEN FIXED"
        }
      },
      "interval": 1,
      "timeout": 1,
      "threshold": 80,
      "headers": {
        "Authorization": "Basic <hash>"
      },
      "expected_status_code": 200,
      "expected_body": "P.*NG"
    }
  ]
}
```

## Installation

1. Download binary from [release page](https://github.com/bennetgallein/cachet-monitor/releases)
2. Add the binary to an executable path (/usr/bin, etc.)
3. Create a configuration following provided examples
4. `cachet-monitor -c /etc/cachet-monitor.json`

pro tip: run in background using `nohup cachet-monitor 2>&1 > /var/log/cachet-monitor.log &`, or use a tmux/screen session

```
Usage:
  cachet-monitor (-c PATH | --config PATH) [--log=LOGPATH] [--name=NAME] [--immediate]
  cachet-monitor -h | --help | --version

Arguments:
  PATH     path to config.json
  LOGPATH  path to log output (defaults to STDOUT)
  NAME     name of this logger

Examples:
  cachet-monitor -c /root/cachet-monitor.json
  cachet-monitor -c /root/cachet-monitor.json --log=/var/log/cachet-monitor.log --name="development machine"

Options:
  -c PATH.json --config PATH     Path to configuration file
  -h --help                      Show this screen.
  --version                      Show version
  --immediate                    Tick immediately (by default waits for first defined interval)
  --restarted                    Get open incidents before start monitoring (if monitor died or restarted)

Environment varaibles:
  CACHET_API      override API url from configuration
  CACHET_TOKEN    override API token from configuration
  CACHET_DEV      set to enable dev logging
```

## Init script

If your system is running systemd (like Debian, Ubuntu 16.04, Fedora, RHEL7, or Archlinux) you can use the provided example file: [example.cachet-monitor.service](https://github.com/bennetgallein/cachet-monitor/blob/master/example.cachet-monitor.service).

1. Simply put it in the right place with `cp example.cachet-monitor.service /etc/systemd/system/cachet-monitor.service`
2. Then do a `systemctl daemon-reload` in your terminal to update Systemd configuration
3. Finally you can start cachet-monitor on every startup with `systemctl enable cachet-monitor.service`! 👍

## Templates

This package makes use of [`text/template`](https://godoc.org/text/template). [Default HTTP template](https://github.com/CastawayLabs/cachet-monitor/blob/master/http.go#L14)

The following variables are available:

| Root objects  | Description                         |
| ------------- | ----------------------------------- |
| `.SystemName` | system name                         |
| `.API`        | `api` object from configuration     |
| `.Monitor`    | `monitor` object from configuration |
| `.now`        | formatted date string               |

| Monitor variables |
| ----------------- |
| `.Name`           |
| `.Target`         |
| `.Type`           |
| `.Strict`         |
| `.MetricID`       |
| ...               |

All monitor variables are available from `monitor.go`

## Monitor Types

We support a variety of monitor-types. Here are the configuration parameters for each of them

Also, the following parameters are shared for all monitors.

| Key                     | Description                                                                    |
| ----------------------- | ------------------------------------------------------------------------------ |
| name                    | a friendly name for the monitor                                                |
| target                  | target for the check (e.g. a domain or IP)                                     |
| active                  | a boolean wether or not this test is currently active                          |
| type                    | type type of the check, see supported types above or below                     |
| interval                | the interval in seconds in which to check the monitor                          |
| timeout                 | the timeout for each check. Needs to be smaller than the interval              |
| metric_id               | the ID of the metric. Metrics are used for graphing values                     |
| component_id            | the ID of the component inside of Cachet. Used for creating incidents          |
| templates.investigating | template to use as a message for when the check enters the investigating stage |
| templates.fixed         | template to use as a message for when the check enters the fixed stage         |
| threshold               | If % of downtime is over this threshold, open an incident                      |
| threshold_count         | the number of checks that count into the threshold (defaults to 10)            |

### HTTP

Either expected_body or expected_status_code needs to be set.

| Key                  | Description                                                                                 |
| -------------------- | ------------------------------------------------------------------------------------------- |
| method               | the HTTP Request method to use (Defaults to GET)                                            |
| headers              | a key-value array of additional headers to use for the request                              |
| expected_status_code | the expected status-code returned from the request                                          |
| expected_body        | a regex or normal string that will be used to test against the returned body of the request |
| expected_md5sum      | a md5 checksum of the body, which will be checked against                                   |
| expected_length      | the length of the string of the response body                                               |
| data                 | body-data for a post request                                                                |

### DNS

| Key      | Description                                                                                          |
| -------- | ---------------------------------------------------------------------------------------------------- |
| dns      | set a custom DNS-Resolver (IP:Port format)                                                           |
| question | the type of DNS-Request to execute (e.g. A, MX, CNAME...). Can also be a List (['A', 'MX', 'CNAME']) |
| answers  | an array of response objects. see below                                                              |

#### Answer Object

| Key   | Description                                                           |
| ----- | --------------------------------------------------------------------- |
| regex | if you want to use a regexp, use this key and the regexp as the value |
| exact | exact match for the response value                                    |

### TCP

| Key  | Description                        |
| ---- | ---------------------------------- |
| port | the port to do a tcp connection to |

### ICMP

_No special variables needed_

## Vision and goals

We made this tool because we felt the need to have our own monitoring software (leveraging on Cachet).
The idea is a stateless program which collects data and pushes it to a central cachet instance.

This gives us power to have an army of geographically distributed loggers and reveal issues in both latency & downtime on client websites.

## Package usage

When using `cachet-monitor` as a package in another program, you should follow what `cli/main.go` does. It is important to call `Validate` on `CachetMonitor` and all the monitors inside.

[API Documentation](https://godoc.org/github.com/CastawayLabs/cachet-monitor)

# Contributions welcome

We'll happily accept contributions for anything usefull.

# Build on Linux/MacOS

1. Read and install with https://ahmadawais.com/install-go-lang-on-macos-with-homebrew/
2. Test in console with `go get -u` and `go build cli/main.go`
3. Run `./go-executable-build.sh cli/main.go`
4. This will create a `build/cli/main.go-linux-amd64`-file, which is the executable binary
5. `mv build/cli/main.go-linux-amd64 /usr/bin/cachet-monitor`
