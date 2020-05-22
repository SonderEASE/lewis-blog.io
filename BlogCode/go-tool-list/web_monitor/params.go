package web_monitor

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"wx-gitlab.xunlei.cn/scdn/x/logger"
)

// Params 请求参数
type Params struct {
	w http.ResponseWriter
	r *http.Request
}

func newParams(w http.ResponseWriter, r *http.Request) *Params {
	err := r.ParseForm()
	if err != nil {
		logger.Errorf("[web_monitor] %v", err)
	}
	return &Params{
		w: w,
		r: r,
	}
}

// GetString 获取string类型的参数
func (p *Params) GetString(key string) (string, bool) {
	if p == nil || p.r == nil {
		return "", false
	}

	v, ok := p.r.Form[key]
	if !ok || len(v) == 0 {
		return "", false
	}
	return strings.TrimSpace(v[0]), ok
}

// GetStringDefault 获取string类型的参数
func (p *Params) GetStringDefault(key, def string) string {
	v, ok := p.GetString(key)
	if !ok {
		return def
	}
	return v
}

func (p *Params) GetRequest() *http.Request {
	return p.r
}

func (p *Params) GetW() http.ResponseWriter {
	return p.w
}

// GetInt 获取int类型的参数
func (p *Params) GetInt(key string) (int, bool) {
	v, ok := p.GetString(key)
	if !ok {
		return 0, false
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return 0, false
	}
	return i, true
}

// GetIntDefault 获取int类型的参数
func (p *Params) GetIntDefault(key string, def int) int {
	v, ok := p.GetInt(key)
	if !ok {
		return def
	}
	return v
}

// GetTime 获取时间参数
// 支持格式1 2019-12-05 16:53:22
// 支持格式2 2019-12-05
// 支持格式3 1575536030 秒时间戳
// 支持格式4 1575536030000 毫秒时间戳
func (p *Params) GetTime(key string) (time.Time, bool) {
	if p == nil || p.r == nil {
		return time.Time{}, false
	}

	v, ok := p.r.Form[key]
	if !ok || len(v) == 0 {
		return time.Time{}, false
	}

	t, err := time.ParseInLocation("2006-01-02 15:04:05", v[0], time.Local)
	if err == nil {
		return t, true
	}

	t, err = time.ParseInLocation("2006-01-02", v[0], time.Local)
	if err == nil {
		return t, true
	}

	i, err := strconv.ParseFloat(v[0], 64)
	if err == nil {
		if i > 1000000000000 {
			return time.Unix(0, int64(i)), true
		}
		if i > 1000000000 {
			return time.Unix(int64(i), 0), true
		}
	}

	return time.Time{}, false
}

// GetBETime 获取起止时间参数
func (p *Params) GetBETime(defaultBegin, defaultEnd time.Time) (begin, end time.Time) {
	var ok bool
	begin, ok = p.GetTime("begin")
	if !ok {
		begin = defaultBegin
	}
	end, ok = p.GetTime("end")
	if !ok {
		end = defaultEnd
	}

	return
}
