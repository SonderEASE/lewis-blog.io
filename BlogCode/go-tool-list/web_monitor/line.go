package web_monitor

var h = []string{"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}
var hLen = int(len(h))

func BarGetInt64(max, current int64) string {
	return BarGetInt(int(max), int(current))
}

func BarGetInt(max, current int) string {
	if max <= 0 {
		return BarMin()
	}

	return h[(hLen-1)*current/max]
}

func BarMin() string {
	return h[0]
}
