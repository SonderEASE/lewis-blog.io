package web_monitor

import (
	"fmt"
	"net/http"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"wx-gitlab.xunlei.cn/scdn/x/conv"
)

type sysCPU struct {
	Info           []cpu.InfoStat `json:"info"`
	Percent        []float64      `json:"percent"`
	IntervalSecond int            `json:"interval_second"`
}

type sysMem struct {
	VirtualMemory *mem.VirtualMemoryStat
}

var (
	gSysCPU sysCPU
	gSysMem sysMem
)

func (s *sysCPU) Format() string {
	if s == nil {
		return "unknown"
	}
	p := s.Percent
	var percent float64
	for _, v := range p {
		percent += v
	}

	return fmt.Sprintf("%.2f%%[%d]", percent, len(p))
}

func (s *sysMem) Format() string {
	if s == nil {
		return "unknown"
	}
	v := s.VirtualMemory
	if v == nil {
		return "unknown"
	}

	return fmt.Sprintf("%s/%s(%.2f%%)",
		conv.ByteFormat(int64(s.VirtualMemory.Used)),
		conv.ByteFormat(int64(s.VirtualMemory.Total)),
		s.VirtualMemory.UsedPercent)
}

func initSYS() {
	gSysCPU.IntervalSecond = 5
	gSysCPU.Info, _ = cpu.Info()
	go func() {
		for {
			gSysCPU.Percent, _ = cpu.Percent(time.Second*time.Duration(gSysCPU.IntervalSecond), true)
		}
	}()

	go func() {
		for {
			gSysMem.VirtualMemory, _ = mem.VirtualMemory()
			time.Sleep(time.Second * 5)
		}
	}()
}

func handleSysCPU(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	if r.Form.Get("json") != "" {
		_, _ = w.Write(conv.ToJsonByte(gSysCPU))
		return
	}

	page := Page{
		title: opt.serviceName + " CPU",
	}

	table := page.AddTable("Info")
	table.AddHeader("CPU", "VendorID", "Family", "Model", "Stepping", "PhysicalID", "CoreID", "Cores", "ModelName", "Mhz", "CacheSize")
	for _, c := range gSysCPU.Info {
		table.AddRow(c.CPU, c.VendorID, c.Family, c.Model, c.Stepping, c.PhysicalID, c.CoreID, c.Cores, c.ModelName, c.Mhz, c.CacheSize)
	}

	table = page.AddTable("Percent")
	table.AddHeader("使用率百分比")
	for _, p := range gSysCPU.Percent {
		table.AddRow(fmt.Sprintf("%.2f%%", p))
	}
	table.AddNo()

	page.output(w)
}

func handleSysMem(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	if r.Form.Get("json") != "" {
		_, _ = w.Write(conv.ToJsonByte(gSysMem))
		return
	}

	virtualMemory := gSysMem.VirtualMemory
	if virtualMemory == nil {
		_, _ = w.Write([]byte("empty"))
		return
	}

	page := Page{
		title: opt.serviceName + " MEM",
	}

	table := page.AddTable("VirtualMemory")

	table.AddRow("Total", conv.ByteFormat(int64(virtualMemory.Total)))
	table.AddRow("Available", conv.ByteFormat(int64(virtualMemory.Available)))
	table.AddRow("Used", conv.ByteFormat(int64(virtualMemory.Used)))
	table.AddRow("UsedPercent", fmt.Sprintf("%.2f%%", virtualMemory.UsedPercent))
	table.AddRow("Free", conv.ByteFormat(int64(virtualMemory.Free)))
	table.AddRow("Active", conv.ByteFormat(int64(virtualMemory.Active)))
	table.AddRow("Inactive", conv.ByteFormat(int64(virtualMemory.Inactive)))
	table.AddRow("Wired", conv.ByteFormat(int64(virtualMemory.Wired)))
	table.AddRow("Laundry", conv.ByteFormat(int64(virtualMemory.Laundry)))
	table.AddRow("Buffers", conv.ByteFormat(int64(virtualMemory.Buffers)))
	table.AddRow("Cached", conv.ByteFormat(int64(virtualMemory.Cached)))
	table.AddRow("Writeback", conv.ByteFormat(int64(virtualMemory.Writeback)))
	table.AddRow("Dirty", conv.ByteFormat(int64(virtualMemory.Dirty)))
	table.AddRow("WritebackTmp", conv.ByteFormat(int64(virtualMemory.WritebackTmp)))
	table.AddRow("Shared", conv.ByteFormat(int64(virtualMemory.Shared)))
	table.AddRow("Slab", conv.ByteFormat(int64(virtualMemory.Slab)))
	table.AddRow("SReclaimable", conv.ByteFormat(int64(virtualMemory.SReclaimable)))
	table.AddRow("SUnreclaim", conv.ByteFormat(int64(virtualMemory.SUnreclaim)))
	table.AddRow("PageTables", conv.ByteFormat(int64(virtualMemory.PageTables)))
	table.AddRow("SwapCached", conv.ByteFormat(int64(virtualMemory.SwapCached)))
	table.AddRow("CommitLimit", conv.ByteFormat(int64(virtualMemory.CommitLimit)))
	table.AddRow("CommittedAS", conv.ByteFormat(int64(virtualMemory.CommittedAS)))
	table.AddRow("HighTotal", conv.ByteFormat(int64(virtualMemory.HighTotal)))
	table.AddRow("HighFree", conv.ByteFormat(int64(virtualMemory.HighFree)))
	table.AddRow("LowTotal", conv.ByteFormat(int64(virtualMemory.LowTotal)))
	table.AddRow("LowFree", conv.ByteFormat(int64(virtualMemory.LowFree)))
	table.AddRow("SwapTotal", conv.ByteFormat(int64(virtualMemory.SwapTotal)))
	table.AddRow("SwapFree", conv.ByteFormat(int64(virtualMemory.SwapFree)))
	table.AddRow("Mapped", conv.ByteFormat(int64(virtualMemory.Mapped)))
	table.AddRow("VMallocTotal", conv.ByteFormat(int64(virtualMemory.VMallocTotal)))
	table.AddRow("VMallocUsed", conv.ByteFormat(int64(virtualMemory.VMallocUsed)))
	table.AddRow("VMallocChunk", conv.ByteFormat(int64(virtualMemory.VMallocChunk)))
	table.AddRow("HugePagesTotal", conv.ByteFormat(int64(virtualMemory.HugePagesTotal)))
	table.AddRow("HugePagesFree", conv.ByteFormat(int64(virtualMemory.HugePagesFree)))
	table.AddRow("HugePageSize", conv.ByteFormat(int64(virtualMemory.HugePageSize)))

	swapMemoryStat, err := mem.SwapMemory()
	if err == nil {
		table = page.AddTable("SwapMemory")
		table.AddRow("Total", conv.ByteFormat(int64(swapMemoryStat.Total)))
		table.AddRow("Used", conv.ByteFormat(int64(swapMemoryStat.Used)))
		table.AddRow("Free", conv.ByteFormat(int64(swapMemoryStat.Free)))
		table.AddRow("UsedPercent", fmt.Sprintf("%.2f%%", swapMemoryStat.UsedPercent))
		table.AddRow("Sin", conv.ByteFormat(int64(swapMemoryStat.Sin)))
		table.AddRow("Sout", conv.ByteFormat(int64(swapMemoryStat.Sout)))
		table.AddRow("PgIn", conv.ByteFormat(int64(swapMemoryStat.PgIn)))
		table.AddRow("PgOut", conv.ByteFormat(int64(swapMemoryStat.PgOut)))
		table.AddRow("PgFault", conv.ByteFormat(int64(swapMemoryStat.PgFault)))
	}

	page.output(w)
}

func handleSysLoad(w http.ResponseWriter, r *http.Request) {
	avgStat, _ := load.Avg()
	if avgStat == nil {
		avgStat = &load.AvgStat{}
	}

	_, _ = w.Write(conv.ToJsonByte(avgStat))
}
