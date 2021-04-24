# sysaudit

Journal System Audits

## Info

This package implements the logic for capturing systemd events and sending them to slack.

Currently the audit is implemented for the sshd daemon only but the logic to extend it for 
additional units or events is the same.

## Installation

First in `main.go` at line 19 update the slack webhook url to your webhook.

Then build

```bash

    go build -o sys-audit main.go

```

And install under `/usr/local/bin`

```bash

    sudo install -m 0755 sys-audit /usr/local/bin/sys-audit

```

Copy the `sys-audit.service` under `/etc/systemd/system`

```bash

    cp sys-audit.service /etc/systemd/system/sys-audit.service

```

Do a daemon reload

```bash

    systemct daemon-reload

```

Enable the systemd unit

```bash

    systemctl enable --now sys-audit.service

```

That's it, now anytime a sshd event is logged into your system it will be also forwarded 
to your slack channel.

