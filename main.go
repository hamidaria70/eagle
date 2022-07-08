package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"strconv"
	"time"

	"code.cloudfoundry.org/bytefmt"
	"github.com/fbiville/markdown-table-formatter/pkg/markdown"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/loadavg"
	"github.com/mackerelio/go-osstat/memory"
	"github.com/mackerelio/go-osstat/uptime"
	"github.com/ricochet2200/go-disk-usage/du"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func memoryUsage() (string, string) {
	memory, err := memory.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}

	totalMemory := bytefmt.ByteSize(memory.Total)
	percentageMemory := strconv.FormatUint(memory.Used*100/memory.Total, 10)

	return totalMemory, percentageMemory
}

func cpuUsage() (string, string) {
	before, err := cpu.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
	time.Sleep(time.Duration(1) * time.Second)
	after, err := cpu.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
	total := float64(after.Total - before.Total)

	percentageCpu := fmt.Sprintf("%.2f", 100-float64(after.Idle-before.Idle)/total*100)
	cpuCores := fmt.Sprintf("%v", runtime.NumCPU())

	return percentageCpu, cpuCores
}

func loadAvarage() string {
	load, err := loadavg.Get()
	checkError(err)

	loadAverage := fmt.Sprintf("%.2f ,%.2f ,%.2f", float64(load.Loadavg1), float64(load.Loadavg5), float64(load.Loadavg15))

	return loadAverage
}

func upTime() time.Duration {
	uptime, err := uptime.Get()
	checkError(err)

	return uptime
}

func diskUsage() (string, string) {
	usage := du.NewDiskUsage(".")
	diskSize := bytefmt.ByteSize(usage.Size())
	percentageDisk := fmt.Sprintf("%.2f", usage.Usage()*100)

	return diskSize, percentageDisk
}

func markdownGenerator(hostname string, ip, uptime string, percentagecpu string, cpuCores string,
	percentagedisk string, disksize string, percentagememory string, totalmemory string,
	loadavarage string) {
	percentagememory = fmt.Sprintf("%v %% out of %vB", percentagememory, totalmemory)
	percentagedisk = fmt.Sprintf("%v %% out of %vB", percentagedisk, disksize)
	percentagecpu = fmt.Sprintf("%v %% out of %v cores", percentagecpu, cpuCores)
	basicTable, _ := markdown.NewTableFormatterBuilder().
		WithPrettyPrint().
		Build("Hostname", "IP Address", "Up Time", "CPU Usage", "Disk Usage",
			"Memory Usage", "Load Average").
		Format([][]string{
			{hostname, ip, uptime, percentagecpu, percentagedisk, percentagememory, loadavarage},
		})

	f, err := os.Create("data.md")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(basicTable)

	if err2 != nil {
		log.Fatal(err2)
	}

}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func main() {
	totalmemory, percentagememory := memoryUsage()
	percentagecpu, cpuCores := cpuUsage()
	loadaverage := loadAvarage()
	uptime := upTime().String()
	disksize, percentagedisk := diskUsage()
	hostname, _ := os.Hostname()
	ip := GetOutboundIP().String()

	markdownGenerator(hostname, ip, uptime, percentagecpu, cpuCores, percentagedisk, disksize,
		percentagememory, totalmemory, loadaverage)
}
