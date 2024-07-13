package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"sync"
	"time"
)

type LastSeenMonitor struct {
	Timestamp time.Time
	RespTime  int64
	StatusUp  bool
}

var LastSeenMonitors = make(map[string]LastSeenMonitor)
var lastSeenMonitorsMutex = sync.Mutex{}

type DhmbCollector struct {
	resptimeMetric *prometheus.Desc
}

func NewDHMBbCollector() *DhmbCollector {
	return &DhmbCollector{resptimeMetric: prometheus.NewDesc("dhmb_resptime", "Response time in ms", []string{"name", "status"}, nil)}
}

// Describe - Each and every collector must implement the Describe function. It essentially writes all descriptors to the exporter desc channel.
func (collector *DhmbCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.resptimeMetric
}

func (collector *DhmbCollector) Collect(ch chan<- prometheus.Metric) {
	var status string
	lastSeenMonitorsMutex.Lock()
	defer lastSeenMonitorsMutex.Unlock()
	for name, lastSeenMon := range LastSeenMonitors {
		if lastSeenMon.Timestamp.After(time.Now().Add(-5 * time.Minute)) {
			status = "DOWN"
			if lastSeenMon.StatusUp {
				status = "UP"
			}
			ch <- prometheus.MustNewConstMetric(collector.resptimeMetric, prometheus.GaugeValue, float64(lastSeenMon.RespTime), name, status)
		}
	}
}
