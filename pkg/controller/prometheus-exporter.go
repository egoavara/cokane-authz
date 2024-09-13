package controller

import (
	"fmt"

	"egoavara.net/authz/pkg/util"
	"github.com/Depado/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
)

type PrometheusExporter struct {
	Prometheus *ginprom.Prometheus
	Paths      []string
}

func NewPrometheusExporter(namespace string, subsystem string, paths []string) *PrometheusExporter {
	Prometheus := ginprom.New(
		ginprom.Namespace(namespace),
		ginprom.Subsystem(subsystem),
		ginprom.BucketSize([]float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}),
	)
	return &PrometheusExporter{
		Prometheus: Prometheus,
		Paths:      paths,
	}
}

func (p *PrometheusExporter) Use(engine *gin.Engine) {
	// CPU 정보
	p.Prometheus.AddCustomGauge("cpu_usage_ratio", "CPU Usage", []string{"cpu"})
	// Memory 정보 (가상 메모리)
	p.Prometheus.AddCustomGauge("memory_total_bytes", "Memory Total", []string{})
	p.Prometheus.AddCustomGauge("memory_used_bytes", "Memory Used", []string{})
	p.Prometheus.AddCustomGauge("memory_available_bytes", "Memory Available", []string{})
	p.Prometheus.AddCustomGauge("memory_free_bytes", "Memory Free", []string{})
	// Swap Memory 정보
	p.Prometheus.AddCustomGauge("swap_memory_total_bytes", "Swap Memory Total", []string{})
	p.Prometheus.AddCustomGauge("swap_memory_used_bytes", "Swap Memory Used", []string{})
	p.Prometheus.AddCustomGauge("swap_memory_free_bytes", "Swap Memory Free", []string{})
	p.Prometheus.AddCustomGauge("swap_memory_sin_bytes", "Swap Memory sin", []string{})
	p.Prometheus.AddCustomGauge("swap_memory_sout_bytes", "Swap Memory sout", []string{})
	p.Prometheus.AddCustomGauge("swap_memory_pgin_bytes", "Swap Memory pgin", []string{})
	p.Prometheus.AddCustomGauge("swap_memory_pgout_bytes", "Swap Memory pgout", []string{})
	p.Prometheus.AddCustomGauge("swap_memory_pgfault_bytes", "Swap Memory pgfault", []string{})
	// 디스크 정보
	p.Prometheus.AddCustomGauge("disk_total_bytes", "Disk Total", []string{"device", "path", "fstype"})
	p.Prometheus.AddCustomGauge("disk_free_bytes", "Disk Free", []string{"device", "path", "fstype"})
	p.Prometheus.AddCustomGauge("disk_used_bytes", "Disk Used", []string{"device", "path", "fstype"})
	// 네트워크 정보
	p.Prometheus.AddCustomGauge("network_send_bytes", "Network Send Bytes", []string{"nic"})
	p.Prometheus.AddCustomGauge("network_recv_bytes", "Network Receive Bytes", []string{"nic"})
	p.Prometheus.AddCustomGauge("network_send_packet_counts", "Network Send Packet Count", []string{"nic"})
	p.Prometheus.AddCustomGauge("network_recv_packet_counts", "Network Receive Packet Count", []string{"nic"})
	// 정보를 업데이트하는 미들웨어
	engine.Use(
		func(ctx *gin.Context) {
			// CPU 정보
			util.Result(cpu.Percent(0, false)).RunOk(func(cpuPercentages []float64) {
				p.Prometheus.SetGaugeValue("cpu_usage_ratio", []string{"total"}, cpuPercentages[0])
			})
			util.Result(cpu.Percent(0, true)).RunOk(func(cpuPercentages []float64) {
				for cpuIndex, cpuUsage := range cpuPercentages {
					p.Prometheus.SetGaugeValue("cpu_usage_ratio", []string{fmt.Sprintf("cpu_%d", cpuIndex)}, cpuUsage)
				}
			})
			// Memory 정보 (가상 메모리)
			util.Result(mem.VirtualMemory()).RunOk(func(virtualMemory *mem.VirtualMemoryStat) {
				p.Prometheus.SetGaugeValue("memory_total_bytes", []string{}, float64(virtualMemory.Total))
				p.Prometheus.SetGaugeValue("memory_used_bytes", []string{}, float64(virtualMemory.Used))
				p.Prometheus.SetGaugeValue("memory_available_bytes", []string{}, float64(virtualMemory.Available))
				p.Prometheus.SetGaugeValue("memory_free_bytes", []string{}, float64(virtualMemory.Free))
			})
			// Swap Memory 정보
			util.Result(mem.SwapMemory()).RunOk(func(swapMemory *mem.SwapMemoryStat) {
				p.Prometheus.SetGaugeValue("swap_memory_total_bytes", []string{}, float64(swapMemory.Total))
				p.Prometheus.SetGaugeValue("swap_memory_used_bytes", []string{}, float64(swapMemory.Used))
				p.Prometheus.SetGaugeValue("swap_memory_free_bytes", []string{}, float64(swapMemory.Free))
				p.Prometheus.SetGaugeValue("swap_memory_sin_bytes", []string{}, float64(swapMemory.Sin))
				p.Prometheus.SetGaugeValue("swap_memory_sout_bytes", []string{}, float64(swapMemory.Sout))
				p.Prometheus.SetGaugeValue("swap_memory_pgin_bytes", []string{}, float64(swapMemory.PgIn))
				p.Prometheus.SetGaugeValue("swap_memory_pgout_bytes", []string{}, float64(swapMemory.PgOut))
				p.Prometheus.SetGaugeValue("swap_memory_pgfault_bytes", []string{}, float64(swapMemory.PgFault))
			})
			// 디스크 정보
			util.Result(disk.Partitions(false)).Run(func(partitions []disk.PartitionStat) error {
				for _, partition := range partitions {
					usage, err := disk.Usage(partition.Mountpoint)
					if err != nil {
						return err
					}
					p.Prometheus.SetGaugeValue("disk_total_bytes", []string{partition.Device, partition.Mountpoint, partition.Fstype}, float64(usage.Total))
					p.Prometheus.SetGaugeValue("disk_used_bytes", []string{partition.Device, partition.Mountpoint, partition.Fstype}, float64(usage.Used))
					p.Prometheus.SetGaugeValue("disk_free_bytes", []string{partition.Device, partition.Mountpoint, partition.Fstype}, float64(usage.Free))
				}
				return nil
			})
			util.Result(net.IOCounters(true)).RunOk(func(ioCounters []net.IOCountersStat) {
				for _, ioCounter := range ioCounters {
					p.Prometheus.SetGaugeValue("network_send_bytes", []string{ioCounter.Name}, float64(ioCounter.BytesSent))
					p.Prometheus.SetGaugeValue("network_recv_bytes", []string{ioCounter.Name}, float64(ioCounter.BytesRecv))
					p.Prometheus.SetGaugeValue("network_send_packet_counts", []string{ioCounter.Name}, float64(ioCounter.PacketsSent))
					p.Prometheus.SetGaugeValue("network_recv_packet_counts", []string{ioCounter.Name}, float64(ioCounter.PacketsRecv))
				}
			})

			ctx.Next()
		},
		p.Prometheus.Instrument(),
	)
	for _, path := range p.Paths {
		p.Prometheus.MetricsPath = path
		p.Prometheus.Use(engine)
	}
}
