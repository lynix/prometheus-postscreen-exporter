# prometheus-postscreen-exporter

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![CI Build](https://github.com/lynix/prometheus-postscreen-exporter/workflows/CI%20Build/badge.svg?branch=master)](https://github.com/lynix/prometheus-postscreen-exporter/actions)

![Grafana Widget](https://github.com/lynix/prometheus-postscreen-exporter/blob/master/grafana.png)

## Summary

*prometheus-postscreen-exporter* exposes check results from Postfix'
[postscreen](http://www.postfix.org/postscreen.8.html) as
[Prometheus](https://prometheus.io) metrics by reading the systemd journal.

**Disclaimer:** This is the result of a hack-a-day and not considered ready for
production. The author had never written a single line of Go before starting
this project.


## Usage

*prometheus-postscreen-exporter* is best run via systemd:

```INI
# prometheus-postscreen-exporter.service

[Unit]
Description=Prometheus exporter for Postfix postscreen

[Service]
ExecStart=/usr/bin/prometheus-postscreen-exporter
User=prometheus
Group=systemd-journal
Restart=always
ProtectSystem=full

[Install]
WantedBy=multi-user.target
```

Options can be set via arguments:

| Flag                  | Description                        | Default    |
|:----------------------|:-----------------------------------|:-----------|
| `-web.listen-address` | Address to listen on for telemetry | `:9101`    |
| `-web.telemetry-path` | Path under which to expose metrics | `/metrics` |


## Bugs / Features

Pull requests are always welcome. Feel free to report bugs or post questions
using the *Issues* function on GitHub.


## License

This project is published under the terms of the *MIT License*. See the file
`LICENSE` for more information.
