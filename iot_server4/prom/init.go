package prom

import (
	"iot_server4/config"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	onlineNum *prometheus.GaugeVec
)

func SetUp() {
	onlineNum = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "online_client_num",
			Help: "online client num",
		},
		[]string{"target"},
	)
	prometheus.MustRegister(onlineNum)
	go func() {
		http.Handle("/metrics", promhttp.HandlerFor(
			prometheus.DefaultGatherer,
			promhttp.HandlerOpts{
				EnableOpenMetrics: true,
			},
		))
		http.ListenAndServe(":"+strconv.Itoa(config.Conf.PromPort), nil)
	}()
}

func ClientOnline() {
	onlineNum.WithLabelValues("client").Add(1)
}

func ClientOffline() {
	onlineNum.WithLabelValues("client").Sub(1)
}

func DeviceOnline() {
	onlineNum.WithLabelValues("device").Add(1)
}

func DeviceOffline() {
	onlineNum.WithLabelValues("device").Sub(1)
}
