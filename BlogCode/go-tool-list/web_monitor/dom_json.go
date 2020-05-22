package web_monitor

// DomJson json模式的dom
type DomJson struct {
	data []byte
}

func (d *DomJson) Html() string {
	if d == nil {
		return ""
	}

	return string(d.data)
}

func (d *DomJson) Script() string {
	return ""
}
