# ProxyPulse
![ProxyPulse](https://i.imgur.com/jNkQPg7.png)

ProxyPulse is a real-time monitoring tool for system metrics, specifically for monitoring the performance of a proxy server. The application displays various metrics, such as transfer rate, CPU usage, total sockets open, total file descriptors, and memory usage, in a user-friendly terminal interface.
## Features

- Real-time monitoring of system metrics
- Visual representation of metrics using line graphs
- Easily extendable to include additional metrics
- User-friendly terminal interface

## Prerequisites

- Ubuntu 18.04 or later
- Requires the process to be registered as a system service (/etc/systemd/system/proxy.service). See the System Process section for more information.

## Install

```bash
You can install ProxyPulse with a single command:

```bash
curl -sL https://github.com/xvertile/ProxyPulse/releases/download/release/proxypulse -o /usr/local/bin/proxypulse && chmod +x /usr/local/bin/proxypulse
```
## Usage
```bash
./proxypulse -p squid -l usa-datacenter
```
### Options
- `-p`: The process name of the proxy server (e.g., squid).
- `-l` (optional): The location of the proxy server (e.g., usa-datacenter, europe-datacenter).

### Example Output

The application displays the following information:
- Transfer Rate (MB/s)
- CPU Usage (%)
- Total Sockets Open
- Total File Descriptors
- Memory Usage (GB)

Each metric is displayed as a line graph in the terminal.

## System Process
I recommend running any proxy server as a system process to ensure that it runs in the background and restarts automatically in case of a failure. This also allows to set the ulimit values for the process, which is essential for a proxy server that requires a high number of file descriptors. This is set via the `LimitNOFILE` parameter in the systemd service file.

You can make a system process with the following steps:
- 1 - Create a new file in `/etc/systemd/system/proxy.service`
- 2 - Add the following content to the file:
```bash
[Unit]
Description=My Awesome Proxy Service
After=network.target

[Service]
Type=simple
ExecStart=/root/proxy/proxy
WorkingDirectory=/root/proxy/proxy
Restart=on-failure
RestartSec=1s
StartLimitBurst=100
StartLimitIntervalSec=3600
[Install]
WantedBy=multi-user.target
```
- 3 - Reload the systemd manager configuration with the following command:
```bash
systemctl daemon-reload
```
- 4 - Start the service with the following command:
```bash
sudo service proxy start
```

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue if you have any suggestions or bug reports.

## License

This project is licensed under the MIT License.