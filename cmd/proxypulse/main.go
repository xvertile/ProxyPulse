package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ProxyPulse/internal/cpu"
	"ProxyPulse/internal/filedescriptors"
	"ProxyPulse/internal/memory"
	"ProxyPulse/internal/network"
	"ProxyPulse/internal/sockets"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var process = "kraken"
var location = "default"

const maxHistoryLength = 100

type Metrics struct {
	Location             string
	TransferRate         int
	TotalSocketsOpen     int32
	TotalFileDescriptors int32
	TotalCPUUsage        float32
	MemoryUsage          float32
	Timestamp            time.Time
}

type MetricsHistory struct {
	TransferRates        []float64
	TotalSocketsOpens    []float64
	TotalFileDescriptors []float64
	TotalCPUUsages       []float64
	MemoryUsages         []float64
}

func collect() (Metrics, error) {
	var metrics Metrics
	metrics.Location = location
	_, txRate, err := network.CalculateTransferRate(1 * time.Second)
	if err != nil {
		return metrics, err
	}
	metrics.TransferRate = int(txRate)
	open, err := sockets.GetTotalOpenSockets()
	if err != nil {
		return metrics, err
	}
	metrics.TotalSocketsOpen = open
	descriptors, err := filedescriptors.GetTotalFileDescriptors(process)
	if err != nil {
		return metrics, err
	}
	metrics.TotalFileDescriptors = descriptors
	usage, err := cpu.GetTotalCPUUsage(process)
	if err != nil {
		return metrics, err
	}
	metrics.TotalCPUUsage = float32(usage)
	totalmemory, err := memory.GetProcessMemoryUsage(process)
	if err != nil {
		return metrics, err
	}
	metrics.MemoryUsage = totalmemory
	metrics.Timestamp = time.Now()
	return metrics, nil
}

func main() {
	processName := flag.String("p", "", "Process name to monitor")
	flag.Parse()
	if *processName == "" {
		log.Fatalf("Please provide a process name to monitor using -p flag")
	}
	process = *processName
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	locationName := flag.String("l", "default", "Location name to monitor")
	flag.Parse()
	location = *locationName
	defer ui.Close()

	title := widgets.NewParagraph()
	title.Text = `
    ____                        ____        __        
   / __ \_________  _  ____  __/ __ \__  __/ /_______ 
  / /_/ / ___/ __ \| |/_/ / / / /_/ / / / / / ___/ _ \
 / ____/ /  / /_/ />  </ /_/ / ____/ /_/ / (__  )  __/
/_/   /_/   \____/_/|_|\__, /_/    \__,_/_/____/\___/ 
                      /____/

A ByteZero project | Process: ` + process + `
Press q or <C-c> to quit
`
	title.SetRect(0, 0, 70, 12)
	title.TextStyle.Fg = ui.ColorCyan

	metricsText := widgets.NewParagraph()
	metricsText.SetRect(0, 12, 70, 20)
	plotHeight := 14
	plotWidth := 70

	cpuUsagePlot := widgets.NewPlot()
	cpuUsagePlot.Title = "CPU Usage (%)"
	cpuUsagePlot.SetRect(70, 0, 70+plotWidth, 20)
	cpuUsagePlot.AxesColor = ui.ColorWhite
	cpuUsagePlot.LineColors[0] = ui.ColorCyan
	cpuUsagePlot.MaxVal = 100

	transferRatePlot := widgets.NewPlot()
	transferRatePlot.Title = "Transfer Rate (MB/s)"
	transferRatePlot.SetRect(0, 20, plotWidth, 20+plotHeight)
	transferRatePlot.AxesColor = ui.ColorWhite
	transferRatePlot.LineColors[0] = ui.ColorCyan

	socketsOpenPlot := widgets.NewPlot()
	socketsOpenPlot.Title = "Total Sockets Open"
	socketsOpenPlot.SetRect(0, 20+plotHeight, plotWidth, 20+2*plotHeight)
	socketsOpenPlot.AxesColor = ui.ColorWhite
	socketsOpenPlot.LineColors[0] = ui.ColorCyan

	fileDescriptorsPlot := widgets.NewPlot()
	fileDescriptorsPlot.Title = "Total File Descriptors"
	fileDescriptorsPlot.SetRect(plotWidth, 20, 2*plotWidth, 20+plotHeight)
	fileDescriptorsPlot.AxesColor = ui.ColorWhite
	fileDescriptorsPlot.LineColors[0] = ui.ColorCyan

	memoryUsagePlot := widgets.NewPlot()
	memoryUsagePlot.Title = "Memory Usage (GB)"
	memoryUsagePlot.SetRect(plotWidth, 20+plotHeight, 2*plotWidth, 20+2*plotHeight)
	memoryUsagePlot.AxesColor = ui.ColorWhite
	memoryUsagePlot.LineColors[0] = ui.ColorCyan

	var history MetricsHistory

	updateMetrics := func() {
		metrics, err := collect()
		if err != nil {
			metricsText.Text = fmt.Sprintf("Error collecting metrics: %v", err)
			return
		}
		metricsText.Text = fmt.Sprintf(
			"Location: %s\nTransfer Rate: %d MB/s\nTotal Sockets Open: %d\nTotal File Descriptors: %d\nTotal CPU Usage: %.2f%%\nMemory Usage: %.2f GB\nTimestamp: %s\n",
			metrics.Location,
			metrics.TransferRate,
			metrics.TotalSocketsOpen,
			metrics.TotalFileDescriptors,
			metrics.TotalCPUUsage,
			metrics.MemoryUsage,
			metrics.Timestamp.Format(time.RFC3339),
		)

		history.TransferRates = append(history.TransferRates, float64(metrics.TransferRate))
		history.TotalSocketsOpens = append(history.TotalSocketsOpens, float64(metrics.TotalSocketsOpen))
		history.TotalFileDescriptors = append(history.TotalFileDescriptors, float64(metrics.TotalFileDescriptors))
		history.TotalCPUUsages = append(history.TotalCPUUsages, float64(metrics.TotalCPUUsage))
		history.MemoryUsages = append(history.MemoryUsages, float64(metrics.MemoryUsage))

		if len(history.TransferRates) > maxHistoryLength {
			history.TransferRates = history.TransferRates[1:]
			history.TotalSocketsOpens = history.TotalSocketsOpens[1:]
			history.TotalFileDescriptors = history.TotalFileDescriptors[1:]
			history.TotalCPUUsages = history.TotalCPUUsages[1:]
			history.MemoryUsages = history.MemoryUsages[1:]
		}

		if len(history.TransferRates) > 5 {
			transferRatePlot.Data = [][]float64{history.TransferRates}
			cpuUsagePlot.Data = [][]float64{history.TotalCPUUsages}
			socketsOpenPlot.Data = [][]float64{history.TotalSocketsOpens}
			fileDescriptorsPlot.Data = [][]float64{history.TotalFileDescriptors}
			memoryUsagePlot.Data = [][]float64{history.MemoryUsages}

			ui.Render(title, metricsText, cpuUsagePlot, transferRatePlot, socketsOpenPlot, fileDescriptorsPlot, memoryUsagePlot)
		} else {
			ui.Render(title, metricsText, cpuUsagePlot)
		}
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	uiEvents := ui.PollEvents()
	sigEvents := make(chan os.Signal, 1)
	signal.Notify(sigEvents, os.Interrupt, syscall.SIGTERM)

	exitChan := make(chan bool)

	go func() {
		for {
			select {
			case e := <-uiEvents:
				if e.Type == ui.KeyboardEvent {
					switch e.ID {
					case "q", "<C-c>":
						exitChan <- true
						return
					}
				}
			case <-sigEvents:
				exitChan <- true
				return
			}
		}
	}()

	for {
		select {
		case <-ticker.C:
			updateMetrics()
		case <-exitChan:
			return
		}
	}
}
