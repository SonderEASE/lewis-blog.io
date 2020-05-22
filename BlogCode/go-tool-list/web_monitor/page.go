package web_monitor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"wx-gitlab.xunlei.cn/scdn/x/web_monitor/dom"
)

// Page 页(一个html页面)
type Page struct {
	title string

	jsonDom   document
	documents []document
	scripts   string

	style string

	// id分配器
	idDistributor int
}

func (p *Page) byte() []byte {
	if p.jsonDom != nil {
		return []byte(p.jsonDom.Html())
	}
	return []byte(fmt.Sprintf(`
<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>%s</title>

  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap@4.4.1/dist/css/bootstrap.min.css" integrity="sha384-Vkoo8x4CGsO3+Hhxv8T/Q5PaXtkKtu6ug5TOeNV6gBiFeWPGFN9MuhOf23Q9Ifjh" crossorigin="anonymous">

  <!--[if lt IE 9]>
    <script src="https://cdn.jsdelivr.net/npm/html5shiv@3.7.3/dist/html5shiv.min.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/respond.js@1.4.2/dest/respond.min.js"></script>
  <![endif]-->

  <script src="https://cdn.jsdelivr.net/npm/jquery@3.4.1/dist/jquery.slim.min.js" integrity="sha384-J6qa4849blE2+poT4WnyKhv5vZF5SrPo0iEjwBvKU7imGFAV0wwj1yYfoRSJoZ+n" crossorigin="anonymous"></script>
  <script src="https://cdn.jsdelivr.net/npm/popper.js@1.16.0/dist/umd/popper.min.js" integrity="sha384-Q6E9RHvbIyZFJoft+2mJbHaEWldlvI9IOYy5n3zV9zzTtmI3UksdQRVvoxMfooAo" crossorigin="anonymous"></script>
  <script src="https://cdn.jsdelivr.net/npm/bootstrap@4.4.1/dist/js/bootstrap.min.js" integrity="sha384-wfSDF2E50Y2D1uUdj0O3uMBJnjuUD4Ih7YwaYd1iqfktj0Uod8GCExl3Og8ifwB6" crossorigin="anonymous"></script>

  <script src="https://cdn.bootcss.com/echarts/4.2.1/echarts.min.js"></script>
</head>
<style>
  %s
</style>
<body style="margin: 20px;font-size: 14px">
  %s
  <script>
    %s
  </script>
</body>
</html>
`, p.title, p.styleHTML(), p.body(), p.script()))
}

func (p *Page) styleHTML() string {
	return `table{
	border: 1px solid silver;
    border-collapse: collapse;
    word-break: break-word;
}
td {
	border: 1px solid silver;
	border-collapse: collapse;
	font-size: 13px;
	height: 15px;
	padding-left: 5px;
	padding-right: 5px;
	text-align: right;
}

th {
	border: 1px solid silver;
	border-collapse: collapse;
	height: 20px;
	padding-left: 5px;
	padding-right: 5px;
}
tr:hover
{
	background-color: #e6f3ff;
}
r{
	color: red;
}
y{
	color:#e2d81f;
}
g{
	color:green;
}
table {
  width:auto !important;
}
.panel {
  display: inline-block !important;
}
th {
  text-align: center;
}
.bg-primary {
    background-color: #337ab7!important;
}
.btn-primary {
    color: #fff;
    background-color: #337ab7!important;
    border-color: #337ab7!important;
}
a {
    color: #337ab7!important;
    text-decoration: none;
    background-color: transparent;
}
` + p.style
}

func (p *Page) AddStyle(s string) {
	p.style += s
}

func (p *Page) body() string {
	html := ""
	for _, d := range p.documents {
		html += d.Html()
	}
	return html
}

func (p *Page) script() string {
	script := `
$(function () {
  $('[data-toggle="tooltip"]').tooltip()
})
$(function () {
  $('[data-toggle="popover"]').popover()
})
function SortTable(obj){
  const index = obj.cellIndex;
  var desc = false;
  if (obj.className == "asc"){
    obj.className = "desc";
    desc = true;
  }else{
    obj.className = "asc";
  }

  const table = obj.parentElement.parentElement.parentElement.children[1];
  let trs = table.children;
  let SortSlice = [];
  for(let k = 0; k < trs.length; k++){
    SortSlice.push(k);
  }
  
  SortSlice.sort(compare(table,index,desc))

  for (let i=0; i<SortSlice.length; i++){
    var tem1 = table.children[SortSlice[i]].cloneNode(true);
    var tem2 = table.children[i].cloneNode(true);
    table.replaceChild(tem2, table.children[SortSlice[i]]);
    table.replaceChild(tem1, table.children[i]);
    for (let j=i+1; j<SortSlice.length; j++){
      if (SortSlice[j] == (i) ){
        SortSlice[j]=SortSlice[i];
        break;
      }
    }
  }
}
function compare(table,index,desc){ //这是比较函数
	if(desc){
		return function(i,j){
        var a = table.children[i].children[index].innerText;
        var b = table.children[j].children[index].innerText;
		if(a===""&&b===""){return 0;}
		if(a===""){return -1;}
		if(b===""){return 1;}
		aa=convtonumber(a);
		bb=convtonumber(b);
        if(aa!=="compare string"&&bb!=="compare string"){return aa - bb;}//升序
        if(a>b){return 1;}
        if(a<b){return -1;}
        return 0;
        }
	}else{
	return function(j,i){
        var a = table.children[i].children[index].innerText;
        var b = table.children[j].children[index].innerText;
		if(a===""&&b===""){return 0;}
		if(a===""){return -1;}
		if(b===""){return 1;}
		aa=convtonumber(a);
		bb=convtonumber(b);
        if(aa!=="compare string"&&bb!=="compare string"){return aa - bb;}//降序
        if(a>b){return 1;}
        if(a<b){return -1;}
        return 0;
        }
	}   
}
	function convtonumber(x){
	if(x===""){return 0;}
		if (convtimetonumber(x)==="compare string"){
		return convunittonumber(x)
		}
		return convtimetonumber(x)
	}
	function convunittonumber(x){
	if(x===""){return 0;}
	if(x[0]==="."){return "compare string";}
		var index=0,dotnum=0;
	for (let i=0;i<x.length;i++){
	if(x[i]==="."){
	    if(dotnum===0){index++;dotnum++;continue;}
		else{break;}	
	}else if(isnumber(x[i])){index++;continue;}
	else{break;}	
	}
    if(index===0){return "compare string"}
	var number=Number(x.substring(0,index));	
	var unit=x.substring(index)
		switch(unit.toUpperCase()){
        case "":
            return number;
        case "/S":
            return number; 
        case "%":
            return number;
		case "B":
            return number;
        case "KB":
            return number*1024;
        case "MB":
            return number*Math.pow(1024,2);
        case "GB":
            return number*Math.pow(1024,3);
        case "TB":
            return number*Math.pow(1024,4);
        case "PB":
            return number*Math.pow(1024,5);
        default:
            return "compare string";
    }		
	}
    function convtimetonumber(x){
	if(x===""){return 0;}
	if(!isnumber(x[0])||isnumber(x[x.length-1])){return "compare string";}
		var index=0,unitleft=0;
        var flagindex=false,flagunit=false;
        var res=0
	for (let i=0;i<x.length;i++){
	if(isnumber(x[i])){
	    if(flagindex===false){
			if (flagunit===true){
				var unit=convtimeunittonumber(x.substring(unitleft,i))
				if (unit==="compare string"){
					return "compare string"
				}
				var number=Number(x.substring(index,unitleft))
				res+=number*unit
			}
		index=i;flagindex=true;flagunit=false;continue;
		}
		else{continue;}	
	}else if(!isnumber(x[i])){
     if(flagunit===false){unitleft=i;flagindex=false;flagunit=true;continue;}
		else{
		continue;}	
    }
	}
	var unit=convtimeunittonumber(x.substring(unitleft))
				if (unit==="compare string"){
					return "compare string"
				}
				var number=Number(x.substring(index,unitleft))
				res+=number*unit
	return res
	}
	
	
	function convtimeunittonumber(unit){
	switch(unit.toUpperCase()){
        case "D":
            return 60*60*24;
		case "H":
            return 60*60;
        case "M":
            return 60;
        case "S":
            return 1;
        case "MS":
            return Math.pow(0.1,3);
        case "μs":
            return Math.pow(0.1,6);
        case "NS":
            return Math.pow(0.1,9);
        default:
            return "compare string";
    }		
    }
	function isnumber(x){
	if(x==="0"||x==="1"||x==="2"||x==="3"||x==="4"||x==="5"||x==="6"||x==="7"||x==="8"||x==="9"){return true;}
	else{return false;}
	}
	function prePage() {
		var search = window.location.search;
		var page = getQueryString("page");
		page = Number(page);
		var pageNext;
		if (page <= 1){
		  pageNext = 1
		}else{
		  pageNext = page-1;
		}
		if (window.location.search.search("page="+page) >= 0){
		  window.location.search = window.location.search.replace("page="+page, "page="+pageNext);
		}else{
		  if (window.location.search == "") {
			window.location.search = "?page="+pageNext
		  }else{
			window.location.search += "&page="+pageNext
		  }
		}
	  }
	  function newPage() {
		var search = window.location.search;
		var page = getQueryString("page");
		page = Number(page);
		if (page == 0){
		  page = 1
		}
		var pageNext = page+1;
		if (window.location.search.search("page="+page) >= 0){
		  window.location.search = window.location.search.replace("page="+page, "page="+pageNext);
		}else{
		  if (window.location.search == "") {
			window.location.search = "?page="+pageNext
		  }else{
			window.location.search += "&page="+pageNext
		  }
		}
	  }
	function getQueryString(name) {
		var reg = new RegExp('(^|&)' + name + '=([^&]*)(&|$)', 'i');
		var r = window.location.search.substr(1).match(reg);
		if (r != null) {
		  return unescape(r[2]);
		}
		return '';
	}
	Date.prototype.Format = function (fmt) { //author: meizz
	  var o = {
		"M+": this.getMonth() + 1, //月份
		"d+": this.getDate(), //日
		"h+": this.getHours(), //小时
		"m+": this.getMinutes(), //分
		"s+": this.getSeconds(), //秒
		"q+": Math.floor((this.getMonth() + 3) / 3), //季度
		"S": this.getMilliseconds() //毫秒
	  };
	  if (/(y+)/.test(fmt)) fmt = fmt.replace(RegExp.$1, (this.getFullYear() + "").substr(4 - RegExp.$1.length));
	  for (var k in o)
	  if (new RegExp("(" + k + ")").test(fmt)) fmt = fmt.replace(RegExp.$1, (RegExp.$1.length == 1) ? (o[k]) : (("00" + o[k]).substr(("" + o[k]).length)));
	  return fmt;
	}
	` + p.scripts
	for _, v := range p.documents {
		script += v.Script()
	}
	return script
}

// AddTable 添加表格
func (p *Page) AddTable(name string) *DomTable {
	table := newDomTable(name)
	p.documents = append(p.documents, table)
	return table
}

// AddPanel 添加空白区域
func (p *Page) AddPanel() *DomPanel {
	panel := &DomPanel{page: p}
	p.documents = append(p.documents, panel)
	return panel
}

// AddRadio 添加单选按钮组
//  例子:
//    tp := params.GetStringDefault("type", "302")
//    radio := page.AddRadio("接口类型", "type")
//    radio.AddItem("302模式", "302", false)
//    radio.AddItem("api模式", "api", false)
//    radio.SetDefault(tp)
//    page.AddPanel().AddSubmitButton("提交")
func (p *Page) AddRadio(name, key string) *DomRadio {
	dom := newRadio(p)
	dom.name = name
	dom.key = key

	p.documents = append(p.documents, dom)
	return dom
}

func (p *Page) AddContainer() *dom.Container {
	c := dom.NewContainer()
	p.documents = append(p.documents, c)
	return c
}

// AddLine 添加曲线
func (p *Page) AddLine(title string, x interface{}, width, height string, data ...interface{}) {
	panel := &DomPanel{}
	p.documents = append(p.documents, panel)
	id := fmt.Sprintf("line-%d", len(p.documents))
	panel.AddHTML(fmt.Sprintf(`<div id="%s" style="width:%s;height:%s;border: 1px solid #e2e2e2;"></div>`, id, width, height))

	option := lineOption{
		Title:  lineOptionTitle{Text: title},
		Legend: lineOptionLegend{Data: nil},
		Toolbox: lineOptionToolbox{Feature: lineOptionToolboxFeature{
			DataZoom:    lineOptionToolboxFeatureDataZoom{YAxisIndex: "none"},
			SaveAsImage: lineOptionToolboxFeatureSaveAsImage{Name: title + time.Now().Format("_20060102_150405")},
			MagicType:   lineOptionToolboxFeatureMagicType{Type: []string{"line", "stack", "bar", "tiled"}},
			Restore:     lineOptionToolboxFeatureRestore{Show: true},
		}},
		Tooltip:  lineOptionTooltip{Trigger: "axis", AxisPointer: lineOptionTooltipAxisPointer{Type: "cross"}, Formatter: "{a0}:{c0}"},
		XAxis:    lineOptionXAxis{Data: x}, // 行
		YAxis:    []*LineOptionYAxis{{}},
		Series:   []*LineOptionSeries{},
		DataZoom: []*lineOptionDataZoom{{Start: 0, End: 100}},
	}

	name := ""
	for i := 0; i < len(data); i++ {
		if i%2 == 0 {
			name = data[i].(string)
			option.Legend.Data = append(option.Legend.Data, name)
		} else {
			option.Series = append(option.Series, &LineOptionSeries{
				Name:   name,
				Type:   "line",
				Data:   data[i],
				Smooth: true,
			})
		}
	}

	b, _ := json.Marshal(option)
	p.AddScript(fmt.Sprintf(`echarts.init(document.getElementById('%s'), 'vintage').setOption(%s);`, id, b))
}

// AddCustomLine 添加自定义曲线
func (p *Page) AddCustomerLine(title string, width, height string, x interface{}, yAxis []*LineOptionYAxis, series ...*LineOptionSeries) *DomLine {
	line := &DomLine{
		id:     fmt.Sprintf("line-%d", len(p.documents)+1),
		title:  title,
		width:  width,
		height: height,
		x:      x,
		yAxis:  yAxis,
		series: series,
	}
	p.documents = append(p.documents, line)

	return line
}

// AddScript 添加脚本
func (p *Page) AddScript(s string) {
	p.scripts += s
}

func (p *Page) SetHTML(s string) {
	if p == nil {
		return
	}

	p.jsonDom = &DomHTML{
		s: s,
	}
}

// SetJson 添加json页面
func (p *Page) SetJson(obj interface{}) {
	if p == nil {
		return
	}

	data, err := json.Marshal(obj)
	if err != nil {
		data = []byte(err.Error())
	}

	p.jsonDom = &DomJson{
		data: data,
	}
}

// AddMapTable 添加一个Map索引的表
func (p *Page) AddMapTable(tableName string, m map[string]interface{}) {
	if p == nil {
		return
	}

	t := p.AddTable(tableName)
	for k, v := range m {
		t.AddRow(k, v)
	}
}

// AddError 添加错误信息
// 可以添加多次，以表格的形式顺序展示
func (p *Page) AddError(err error) {
	table := p.AddTable("错误")
	table.AddRow(RedSpan(err.Error()))
}

func (p *Page) output(w http.ResponseWriter) {
	write(w, p.byte())
}

func (p *Page) newID() string {
	p.idDistributor++
	return fmt.Sprint(p.idDistributor)
}
