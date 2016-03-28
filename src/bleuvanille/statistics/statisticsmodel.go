package statistics

// Counter holds a statistics
type Counter struct {
	CounterName string `json:"counter"`
	Count       uint32 `json:"count"`
}
