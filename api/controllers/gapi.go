package controllers

import (
	"encoding/json"
	"runtime"
	"strconv"

	"github.com/Glintt/gAPI/api/cache"
	"github.com/Glintt/gAPI/api/config"
	"github.com/Glintt/gAPI/api/http"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
)

func InvalidateCache(c *routing.Context) error {
	cache.InvalidateCache()
	c.Response.SetBody([]byte(`{"error":false, "msg": "Invalidation finished."}`))
	c.Response.Header.SetContentType("application/json")
	return nil
}

func ReloadServices(c *routing.Context) error {
	// InitServices()

	cache.InvalidateCache()

	c.Response.SetBody([]byte(`{"error":false, "msg": "Reloaded successfully."}`))
	c.Response.Header.SetContentType("application/json")
	return nil
}

func ProfileGApiUsage(c *routing.Context) error {

	var profileStats map[string]interface{}
	profileStats = make(map[string]interface{})

	profileStats["OS"] = runtime.GOOS

	// Go Routines
	profileStats["GoRoutines"] = runtime.NumGoroutine()

	// Memory Stats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	var memStats map[string]interface{}
	memStats = make(map[string]interface{})

	memStats["Alloc"] = strconv.Itoa(int(bToMb(m.Alloc))) + " Mb"
	memStats["TotalAlloc"] = strconv.Itoa(int(bToMb(m.TotalAlloc))) + " Mb"
	memStats["Sys"] = strconv.Itoa(int(bToMb(m.Sys))) + " Mb"
	memStats["NumGC"] = m.NumGC
	profileStats["Memory"] = memStats

	// CPU Stats
	cpuStat, _ := cpu.Info()
	profileStats["CPU"] = cpuStat

	percentage, _ := cpu.Percent(0, true)

	var cpuUsageMap []string
	for _, cpupercent := range percentage {
		cpuUsageMap = append(cpuUsageMap, strconv.FormatFloat(cpupercent, 'f', 2, 64))
	}
	profileStats["CPU_Usage"] = cpuUsageMap

	// Host Stats
	hostStat, _ := host.Info()
	profileStats["HostInfo"] = hostStat

	jsonBytes, _ := json.Marshal(profileStats)

	http.Response(c, string(jsonBytes), 200, config.GAPI_SERVICE_NAME, "application/json")
	return nil
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
