package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MittMetrics struct {
	TotalMitts prometheus.Gauge
	TotalLikes prometheus.Gauge
}

func NewMittMetrics() *MittMetrics {
	return &MittMetrics{
		TotalMitts: promauto.NewGauge(prometheus.GaugeOpts{Name: "mitter_mitts_total"}),
		TotalLikes: promauto.NewGauge(prometheus.GaugeOpts{Name: "mitter_mitts_likes_total"}),
	}
}

func (m *MittMetrics) AddMitt() {
	m.TotalMitts.Inc()
}

func (m *MittMetrics) DeleteMitt() {
	m.TotalMitts.Dec()
}

func (m *MittMetrics) AddLike() {
	m.TotalLikes.Inc()
}

func (m *MittMetrics) DeleteLike() {
	m.TotalLikes.Dec()
}
