package web_monitor

// DomHTML ...
type DomHTML struct {
	s string
}

func (d *DomHTML) Html() string {
	if d == nil {
		return ""
	}

	return d.s
}

func (d *DomHTML) Script() string {
	return ""
}
