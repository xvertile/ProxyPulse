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

## Install
You can install ProxyPulse with a single command:

```bash
curl -sL https://github.com/yourusername/ProxyPulse/releases/latest/download/proxypulse -o /usr/local/bin/proxypulse && chmod +x /usr/local/bin/proxypulse
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

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue if you have any suggestions or bug reports.

## License

This project is licensed under the MIT License.