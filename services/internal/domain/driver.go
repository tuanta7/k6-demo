package domain

type Driver struct {
	ID     string
	Name   string
	Phone  string
	Email  string
	Rating float64
}

type Behavior struct {
	PreferredSpeedBand []int
	MicroRouteBiases   map[string]any
	StopPattern        map[string]any
}
