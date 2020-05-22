package web_monitor

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hpcloud/tail"
	"golang.org/x/net/websocket"
	"wx-gitlab.xunlei.cn/scdn/x/conv"
	"wx-gitlab.xunlei.cn/scdn/x/logger"
)

func handleLogDir(w http.ResponseWriter, r *http.Request) {
	page := Page{
		title: opt.serviceName + " 日志目录",
	}

	fs, err := ioutil.ReadDir("../log")
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	table := page.AddTable("../log")
	table.AddHeader("name", "size", "time", "实时日志")

	ts := time.Now().Unix()
	for _, f := range fs {
		if f.IsDir() {
			continue
		}
		name := f.Name()

		// 只给修改时间在1分钟内的返回tail连接
		var tailLink string
		if ts-f.ModTime().Unix() < 60 {
			tailLink = LinkAWithParams("/log/tail", "打开", "file_name", name)
		}

		table.AddRow(
			LinkAWithParams("/log/file", name, "file_name", name),
			conv.ByteFormat(f.Size()),
			f.ModTime().Format("2006-01-02 15:04:05"),
			tailLink,
		)
	}
	table.MultSort(3, true, 0, false)
	table.AddNo()

	page.output(w)
}

func handleLogFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	fileName := r.Form.Get("file_name")
	if fileName == "" {
		_, _ = w.Write([]byte("缺少参数 file_name"))
		return
	}

	f, err := os.Open("../log/" + fileName)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	defer func() { _ = f.Close() }()

	page := conv.ToInt(r.Form.Get("page"))
	pageSize := conv.ToInt(r.Form.Get("page_size"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 40
	}
	if pageSize > 3000 {
		pageSize = 3000
	}

	rd := bufio.NewReader(f)

	p := Page{
		title: opt.serviceName + " 日志",
	}
	p.AddPanel().AddHTML(PageButton())
	table := p.AddTable(fileName)

	for i := 0; i < page*pageSize; i++ {
		line, err := rd.ReadString('\n')
		if err != nil {
			break
		}
		if i < (page-1)*pageSize {
			continue
		}

		line = strings.Replace(line, "[EROR]", RedSpan("[EROR]"), 1)
		line = strings.Replace(line, "[WARN]", YellowSpan("[WARN]"), 1)
		table.AddRow(i+1, line)
	}
	table.TextAlignLeft(1)

	p.output(w)
}

func handleLogTail(w http.ResponseWriter, r *http.Request) {
	logger.Infof("[web_monitor] /log/tail from %s", r.RemoteAddr)
	err := r.ParseForm()
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	fileName := r.Form.Get("file_name")
	if fileName == "" {
		_, _ = w.Write([]byte("缺少参数 file_name"))
		return
	}

	html := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="zh-C">
<head>
    <meta charset="UTF-8">
    <title>Tools</title>
</head>
<body>
<div id="id_log" style="height:100%%">
<!-- 	<div class="terminal">
		<div class="log" id="id_log">
		</div>
	</div> -->
</div>
<style>
.terminal {
	margin-top: 20px;
	height: 100%%;
}
.terminal .log {
	background-color: black;
    font-family: monospace;
    color: #22da26;
    
	padding: 10px;
	overflow:auto;
	height: 85%%;
	word-break:break-all
}
.terminal .title {
	font-size: 20px;
}
html,body{
         height:100%%;
        }
</style>
<script type="text/javascript">
	function getQueryString(name) {
		var reg = new RegExp('(^|&)' + name + '=([^&]*)(&|$)', 'i');
		var r = window.location.search.substr(1).match(reg);
		if (r != null) {
		  return unescape(r[2]);
		}
		return '';
	}

    var ws = new WebSocket("ws://%s/log/tail/ws");

	ws.onopen = function(){
		console.log("open");
		var fileName = getQueryString("file_name");
		if (fileName == "") {
			console.log("缺少参数 file_name");
			ws.close();
			return;
        }
		console.log(fileName);
		ws.send(fileName);
	};
	
	ws.onmessage = function(evt){
		var msg = evt.data;
		var id = 'id-log-xxx';
		var log = document.getElementById(id);
		if (log == null){
			var d = document.getElementById('id_log');
			var terminal = document.createElement("div");
			terminal.className = "terminal";

			var title = document.createElement("div");
			title.className = "title";
			title.innerHTML = "%s";
			terminal.appendChild(title);

			log = document.createElement("div");
			log.className = "log";
			log.id = id;
			terminal.appendChild(log);

            d.appendChild(terminal);
		}
		if(msg.indexOf('EROR') !== -1){
			log.innerHTML+= "<br><span style=\"color:#ff5151;font-size: 14px;\">"+msg+"</span>";
		} else if(msg.indexOf('WARN') !== -1){
			log.innerHTML+= "<br><span style=\"color:#c3b900;font-size: 13px;\">"+msg+"</span>";
		} else if(msg.indexOf('INFO') !== -1){
			log.innerHTML+= "<br><span style=\"color:#00b34c;font-size: 12px;\">"+msg+"</span>";
		} else {
			log.innerHTML+= "<br><span style=\"color:#a2a2a2;font-size: 11px;\">"+msg+"</span>";
		}
		if(log.innerHTML.length > 32*1024){
			log.innerHTML = log.innerHTML.substring(log.innerHTML.length-32*1024)
		}
		
        log.scrollTop = log.scrollHeight;
	};
	
	ws.onclose = function(evt){
	  console.log("WebSocketClosed!");
	};
	
	ws.onerror = function(evt){
	  console.log("WebSocketError!");
	};
</script>
</body>
</html>
`, r.Host, fileName)

	_, _ = w.Write([]byte(html))
}

func handleLogTailWs(ws *websocket.Conn) {
	var fileName string
	err := websocket.Message.Receive(ws, &fileName)
	if err != nil {
		_, _ = ws.Write([]byte(err.Error()))
		return
	}
	logger.Infof("[web_monitor] handle log file %s", fileName)
	t, err := tail.TailFile("../log/"+fileName, tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	})
	if err != nil {
		_, _ = ws.Write([]byte(err.Error()))
		return
	}
	defer func() {
		logger.Infof("[web_monitor] tail log file %s closed", fileName)
		_ = t.Stop()
		t.Cleanup()
	}()

	for {
		msg, ok := <-t.Lines
		if !ok {
			_, _ = ws.Write([]byte("close"))
			return
		}
		if msg.Err != nil {
			_, _ = ws.Write([]byte(err.Error()))
			return
		}
		_, err = ws.Write([]byte(msg.Text))
		if err != nil {
			_, _ = ws.Write([]byte(err.Error()))
			return
		}
	}
}
