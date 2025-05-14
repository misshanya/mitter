package models

type MittMetrics interface {
	AddMitt()
	DeleteMitt()

	AddLike()
	DeleteLike()

	ViewInFeed(count float64)
}
