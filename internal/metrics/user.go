package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type UserMetrics struct {
	TotalUsers prometheus.Gauge
}

func NewUserMetrics() *UserMetrics {
	return &UserMetrics{TotalUsers: promauto.NewGauge(prometheus.GaugeOpts{
		Name: "mitter_users_total",
	})}
}

func (m *UserMetrics) AddUser() {
	m.TotalUsers.Inc()
}

func (m *UserMetrics) DeleteUser() {
	m.TotalUsers.Dec()
}
