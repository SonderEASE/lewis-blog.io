package web_monitor

import (
	"fmt"
)

// DomPanel 空dom
type DomPanel struct {
	page     *Page
	inner    string
	rowStyle string
}

func (p *DomPanel) Html() string {
	return fmt.Sprintf(`
<div class="row" style="margin-bottom: 10px;%s">
  	<div class="col-md-12">
		%s
	</div>
</div>
`, p.rowStyle, p.inner)
}

func (p *DomPanel) SetRowStyle(style string) {
	p.rowStyle = style
}

// SetHtml 设置html内容
func (p *DomPanel) SetHtml(html string) {
	p.inner = html
}

// AddH1 添加 <h1> dom
func (p *DomPanel) AddH1(text string) {
	p.inner += "<h1>" + text + "</h1>"
}

// AddH2 添加 <h2> dom
func (p *DomPanel) AddH2(text string) {
	p.inner += "<h2>" + text + "</h2>"
}

// AddH3 添加 <h3> dom
func (p *DomPanel) AddH3(text string) {
	p.inner += "<h3>" + text + "</h3>"
}

// AddH4 添加 <h4> dom
func (p *DomPanel) AddH4(text string) {
	p.inner += "<h4>" + text + "</h4>"
}

// AddH5 添加 <h5> dom
func (p *DomPanel) AddH5(text string) {
	p.inner += "<h5>" + text + "</h5>"
}

// AddP 添加 <p> dom
func (p *DomPanel) AddP(text string) {
	p.inner += "<p>" + text + "</p>"
}

// AddPSmall 添加 <p><small>
func (p *DomPanel) AddPSmall(text string) {
	p.inner += "<p><small>" + text + "</small></p>"
}

// AddLink 添加超链接
func (p *DomPanel) AddLink(url, title string) {
	p.inner += LinkA(url, title)
}

// AddHTML 添加HTML内容
func (p *DomPanel) AddHTML(html string) {
	p.inner += html
}

// AddLinkWithParams 添加带参数的跳转链接
// 例 panel.AddLinkWithParams("/user", "用户", "name", "张三", "age", 13) => <a target="_blank" href="/user?name=张三&age=13">用户</a>
func (p *DomPanel) AddLinkWithParams(url, title string, parms ...interface{}) {
	p.inner += LinkAWithParams(url, title, parms...)
}

// AddInput 添加输入框
func (p *DomPanel) AddInput(name, key, value string) {
	p.inner += fmt.Sprintf(`<span>%s: </span><input style="width:300px" class="id-input" key="%s" value="%s"></input>`, name, key, value)
}

// AddSubmitButton 添加提交按钮
func (p *DomPanel) AddSubmitButton(name string) {
	p.inner += fmt.Sprintf(`<button class="btn btn-primary btn-sm" onclick="window.onsubmit()" style="min-width: 100px;">%s</button><script>
	function onsubmit(){
	  var url = location.href;
	  if (url.indexOf("?") < 0) {
	  	url+="?";
	  }

	  var inputs = document.getElementsByClassName("id-input");
	  for (let i=0; i<inputs.length; i++) {
	  	let input = inputs[i];
	  	console.log(input.getAttribute("key"), input.value);
	  	url = addURLParam(url, input.getAttribute("key"), encodeURIComponent(input.value));
	  }

	  var radios = document.getElementsByClassName("c-radio");
      for(var i=0; i<radios.length; i++){
        var radio = radios[i];
        for(var j=0;;j++){
          var item = document.getElementById(radio.id+"-"+j);
          if (!item){
            break
          }
          if(item.checked){
			url = addURLParam(url, radio.getAttribute("name"), encodeURIComponent(item.value));
            break;
          }
        }
      }
		window.location.href=url;
	}
	function addURLParam(url, key, value){
		var urlSplit = url.split("?");
		if (urlSplit.length > 1) {
			var arrPara = urlSplit[1].split("&")
			var find = false;
			for (var i = 0; i < arrPara.length; i++) {
				if (arrPara[i].split("=")[0] == key) {
					arrPara[i] = key+"="+value;
					find = true;
					break;
				}
			}
			if (!find){
				arrPara.push(key+"="+value);
			}
			return urlSplit[0]+"?"+arrPara.join("&");
		}else{
			return url+"?"+key+"="+value;
		}
	}
	</script>`, name)
}

// AddTimeBeginEndInput 添加起止时间选择控件
// 包涵起止时间两个输入框和多个快捷设置按钮
func (p *DomPanel) AddTimeBeginEndInput(begin, end string) {
	p.inner += fmt.Sprintf(`
<div class="form-inline">
  <div class="form-row align-items-center">
    <label>Begin：</label>
    <input type="text2" class="id-input" key="begin" id="id-input-time-begin" value="%s">
    <label style="padding-left: 10px;">End：</label>
    <input type="text" class="id-input" key="end" id="id-input-time-end" value="%s">
    <div class="col-auto" style="padding-left: 10px;">
      <button class="btn btn-secondary btn-sm" onclick="window.onChoiceTimeButton(60*60)">1小时</button>
      <button class="btn btn-secondary btn-sm" onclick="window.onChoiceTimeButton(60*60*3)">3小时</button>
      <button class="btn btn-secondary btn-sm" onclick="window.onChoiceTimeButton(60*60*6)">6小时</button>
      <button class="btn btn-secondary btn-sm" onclick="window.onChoiceTimeButton(60*60*12)">12小时</button>
      <button class="btn btn-secondary btn-sm" onclick="window.onChoiceTimeButton(60*60*24)">1天</button>
      <button class="btn btn-secondary btn-sm" onclick="window.onChoiceTimeButton(60*60*24*2)">2天</button>
      <button class="btn btn-secondary btn-sm" onclick="window.onChoiceTimeButton(60*60*24*7)">7天</button>
    </div>
  </div>
</div>
<script>
	function onChoiceTimeButton(second){
        var t = (new Date());
		document.getElementById("id-input-time-end").value = t.Format("yyyy-MM-dd hh:mm:ss");
        t.setSeconds(t.getSeconds()-second);
		document.getElementById("id-input-time-begin").value = t.Format("yyyy-MM-dd hh:mm:ss");
    }
</script>
`, begin, end)
}

func (p *DomPanel) Script() string {
	return ""
}
