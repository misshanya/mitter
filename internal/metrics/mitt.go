package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MittMetrics struct {
	TotalMitts   prometheus.Gauge
	TotalLikes   prometheus.Gauge
	ViewedInFeed prometheus.Counter
}

func NewMittMetrics() *MittMetrics {
	return &MittMetrics{
		TotalMitts:   promauto.NewGauge(prometheus.GaugeOpts{Name: "mitter_mitts_total"}),
		TotalLikes:   promauto.NewGauge(prometheus.GaugeOpts{Name: "mitter_mitts_likes_total"}),
		ViewedInFeed: promauto.NewCounter(prometheus.CounterOpts{Name: "mitter_mitts_feed_viewed"}),
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

func (m *MittMetrics) ViewInFeed(count float64) {
	m.ViewedInFeed.Add(count)
}
