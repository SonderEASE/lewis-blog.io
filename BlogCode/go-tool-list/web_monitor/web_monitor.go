package web_monitor

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/shirou/gopsutil/load"
	"golang.org/x/net/websocket"
	"wx-gitlab.xunlei.cn/scdn/x/conv"
	"wx-gitlab.xunlei.cn/scdn/x/exit"
	"wx-gitlab.xunlei.cn/scdn/x/logger"
	"wx-gitlab.xunlei.cn/scdn/x/sysx"
	"wx-gitlab.xunlei.cn/scdn/x/timex"
)

var (
	opt           *webMonitorOptions
	handles       []*handle
	noAccountPath = map[string]bool{}
	_hostname     string
)

// HandleFunc 打开页面时触发回调函数
type HandleFunc func(*Page, *Params)

type handle struct {
	path   string
	title  string
	isShow bool

	hf HandleFunc
}

func (h *handle) NewPage() *Page {
	return &Page{
		title: h.title,
	}
}

// Init 初始化WebMonitor
// 参数1：监听端口
// 参数2...：参考web_monitor.WithXXX()
func Init(port int, opts ...WebMonitorOption) {
	opt = newDefaultWebMonitorOptions()
	opt.Init(opts)
	opt.port = port

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleBase)
	mux.HandleFunc("/user/login", handleLogin)
	mux.HandleFunc("/file/upload", handleUploadFile)
	mux.HandleFunc("/log/dir", handleLogDir)
	mux.HandleFunc("/log/file", handleLogFile)
	mux.Handle("/log/tail/ws", websocket.Handler(handleLogTailWs))
	mux.HandleFunc("/log/tail", handleLogTail)
	mux.HandleFunc("/sys/cpu", handleSysCPU)
	mux.HandleFunc("/sys/mem", handleSysMem)
	mux.HandleFunc("/sys/load", handleSysLoad)
	if opt.pProf {
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	go func() {
		var err error
		for i := 0; i < 3; i++ {
			err = http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
			time.Sleep(time.Millisecond * 100)
		}
		if err != nil {
			exit.ExitError(fmt.Errorf("[web_monitor] %w", err))
		}
		logger.Error(err)
	}()

	initSYS()
	err := sysx.InitNet()
	if err != nil {
		logger.Errorf("[web_monitor] %v", err)
	}

	_hostname, err = os.Hostname()
	if err != nil {
		_hostname = err.Error()
	}

	logger.Infof("[web_monitor] init, listen at %d", port)
}

// AddPage 添加一个页面，显示在首页跳转列表里
// 例：web_monitor.AddPage("/monitor/domain/info", "domain", func(page *web_monitor.Page, params *web_monitor.Params) {})
func AddPage(path, title string, hf HandleFunc) {
	handles = append(handles, &handle{
		hf:     hf,
		title:  title,
		path:   path,
		isShow: true,
	})
}

// AddHidePage 添加一个页面，不显示在首页跳转列表里(隐藏按钮可以显示)
func AddHidePage(path, title string, hf HandleFunc) {
	handles = append(handles, &handle{
		hf:     hf,
		title:  title,
		path:   path,
		isShow: false,
	})
}

// AddHidePageWithNoAccount 添加无需登录的隐藏页面
func AddHidePageWithNoAccount(path, title string, hf HandleFunc) {
	noAccountPath[path] = true
	AddHidePage(path, title, hf)
}

func handleBase(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/favicon.ico" {
		return
	}

	logger.Infof("[web_monitor] handle path %s, from %s", path, r.RemoteAddr)

	if path == "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if !noAccountPath[path] && !checkAccount(r) {
		http.Redirect(w, r, "/user/login?r="+base64.URLEncoding.EncodeToString([]byte(r.URL.String())), http.StatusFound)
		return
	}

	if path == "/monitor" || path == "/monitor/" {
		handleHomePage(w, r)
		return
	}

	var handle *handle
	for _, h := range handles {
		if h.path == path {
			handle = h
			break
		}
	}
	if handle == nil {
		write(w, []byte("unknown path "+path))
		return
	}

	page := handle.NewPage()
	handle.hf(page, newParams(w, r))
	page.output(w)
}

func write(w http.ResponseWriter, data []byte) {
	_, err := w.Write(data)
	if err != nil {
		logger.Errorf("[web_monitor] %v", err)
	}
}

func handleHomePage(w http.ResponseWriter, r *http.Request) {
	eUp, eDown := sysx.GetExternalBandWidth()
	iUp, iDown := sysx.GetInternalBandWidth()
	avgStat, _ := load.Avg()
	if avgStat == nil {
		avgStat = &load.AvgStat{}
	}

	page := Page{
		title: opt.serviceName,
	}
	container := page.AddContainer()
	container.AddStyle("background-color: #f9f9f9;")

	cal := container.AddRow().AddCol()
	inner := fmt.Sprintf(`<div style="font-size:60px;font-weight: 700;text-decoration: underline;font-family: Helvetica Neue;color: #333;"><strong>%s</strong></div>`, opt.serviceName) +
		`<div>Hostname: <strong>` + _hostname + `</strong></div>` +
		fmt.Sprintf(`<div>版本: <strong>%s</strong> 启动时间: <strong>%s</strong> 运行时长: <strong>%s</strong></div>`,
			opt.version,
			opt.startTime.Format("2006-01-02 15:04:05"),
			timex.DurationFormat(time.Since(opt.startTime)),
		) +
		fmt.Sprintf("<div>CPU: <strong>%s</strong> MEM: <strong>%s</strong> Load Avg: <strong>%.2f, %.2f, %.2f</strong></div>",
			LinkA("/sys/cpu", gSysCPU.Format()),
			LinkA("/sys/mem", gSysMem.Format()),
			avgStat.Load1, avgStat.Load5, avgStat.Load15,
		) +
		fmt.Sprintf("<div>带宽 内网: <strong>%s|%s</strong> 外网: <strong>%s|%s</strong></div>",
			conv.ByteFormat(int64(iUp)), conv.ByteFormat(int64(iDown)),
			conv.ByteFormat(int64(eUp)), conv.ByteFormat(int64(eDown)),
		)
	if opt.pProf {
		ip := strings.Split(r.Host, ":")[0]
		inner += fmt.Sprintf(`
<div style="overflow-x: scroll">
	<table class="table table-striped table-bordered table-hover table-condensed table-sm" style="margin-bottom:0px;background-color: white;margin-top: 10px;min-width:590px">
		<tbody>
			<tr><td>MEM</td><td>go tool pprof -http : http://%s:%d/debug/pprof/heap</td></tr>
			<tr><td>CPU</td><td>go tool pprof -http : -seconds 10 http://%s:%d/debug/pprof/profile</td></tr>
		</tbody>
	</table>
</div>
`, ip, opt.port, ip, opt.port)
	}
	cal.SetInner(inner)

	sort.Slice(handles, func(i, j int) bool {
		return handles[i].title > handles[j].title
	})

	var linkShow [][]string
	var linkHide [][]string
	for _, handle := range handles {
		if !handle.isShow {
			if len(linkHide) == 0 || len(linkHide[len(linkHide)-1]) == 5 {
				linkHide = append(linkHide, []string{})
			}
			linkHide[len(linkHide)-1] = append(linkHide[len(linkHide)-1], `<li class="list-group-item" style="padding: 0.2rem 1.25rem;">`+LinkA(handle.path, handle.title)+`</li>`)
		} else {
			if len(linkShow) == 0 || len(linkShow[len(linkShow)-1]) == 10 {
				linkShow = append(linkShow, []string{})
			}
			linkShow[len(linkShow)-1] = append(linkShow[len(linkShow)-1], `<li class="list-group-item" style="padding: 0.2rem 1.25rem;">`+LinkA(handle.path, handle.title)+`</li>`)
		}
	}

	row := container.AddRow()
	row.AddStyle(`margin-top: 10px;`)
	for _, link := range linkShow {
		row.AddCol().SetInner(strings.Join(link, "")).AddStyle("min-width: 320px", "padding-bottom: 10px")
	}

	row = container.AddRow()
	row.AddStyle(`margin-top: 10px;`)
	row.AddCol().SetInner(`<button type="button" class="btn btn-secondary btn-sm" onclick="showHideMenu()" id="id-btn-show-hide" style="font-size: 11px;">显示隐藏连接</button>`)

	row = container.AddRow()
	row.SetID(`id-hide-menu`)
	row.AddStyle(`display: none`)
	row.AddStyle(`margin-top: 10px;`)
	for _, link := range linkHide {
		row.AddCol().SetInner(strings.Join(link, "")).AddStyle("min-width: 320px", "padding-bottom: 10px")
	}

	page.AddScript(`
function showHideMenu(){
	var dom = document.getElementById("id-hide-menu");
	if (dom.style.display === "none") {
		dom.style.display = "";
		document.getElementById("id-btn-show-hide").innerText = "收起隐藏连接";
	} else {
		dom.style.display = "none";
		document.getElementById("id-btn-show-hide").innerText = "显示隐藏连接";
	}
}`)
	page.output(w)
}

func GetEncryptionMethod() Encryption {
	if opt != nil {
		return opt.encryption
	}
	return nil
}

func GetHostname() string {
	return _hostname
}
