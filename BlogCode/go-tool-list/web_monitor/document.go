package web_monitor

import (
	"fmt"
	"strconv"
	"strings"

	"wx-gitlab.xunlei.cn/scdn/x/locate"
)

type document interface {
	Html() string
	Script() string
}

// RedSpan 生成红色字体的span
func RedSpan(inner interface{}) string {
	return "<span style=\"color:red\">" + fmt.Sprint(inner) + "</span>"
}

// ColorBool 根据bool生成红绿span
func ColorBool(ok bool) string {
	if ok {
		return GreenSpan(true)
	}
	return RedSpan(false)
}

// RedIfSpan 如果is==true输出红色，否则原色。
func RedIfSpan(inner interface{}, is bool) string {
	if is {
		return RedSpan(inner)
	}
	return fmt.Sprint(inner)
}

// GreenSpan 生成绿色字体的span
func GreenSpan(inner interface{}) string {
	return "<span style=\"color:green\">" + fmt.Sprint(inner) + "</span>"
}

// TipSpan 带tip的span
func TipSpan(html, tip string) string {
	return fmt.Sprintf(`<span data-toggle="tooltip" data-placement="top" title="%s">%s</span>`,
		tip, html,
	)
}

// PopoverSpanLimitLength 弹出框
// 文本显示不超过showLengthMax个字符，点击后显示完整内容
func PopoverSpanLimitLength(text string, showLengthMax int) string {
	var show string
	if len(text) > showLengthMax {
		show = text[:showLengthMax] + "..."
		return fmt.Sprintf(`<a type="button" data-placement="top" data-toggle="popover" data-content="%s">%s</a>`,
			text, show,
		)
	}
	return text
}

// CollapseSpanExcerptTitle 下拉文本(节选部分内容做标题)
func CollapseSpanExcerptTitle(body string, length int) string {
	if len(body) <= length {
		return body
	}
	title := body[:length] + "..."
	return CollapseSpan(title, body)
}

var _collapseID int32

// CollapseSpan 下拉文本
func CollapseSpan(title, body string) string {
	_collapseID++
	return fmt.Sprintf(`
  <div>
    <a class="collapsed" data-toggle="collapse" data-parent="#accordion" href="#collapse-id-%d" aria-expanded="false" aria-controls="collapseOne">
      %s
    </a>
    <div id="collapse-id-%d" class="panel-collapse collapse">
      <div style="margin-top: 5px;border-top: 1px dashed #b1b1b1">%s</div>
    </div>
  </div>`, _collapseID, title, _collapseID, body)
}

// IPList 将ip列表格式化展示
// 电信|联通|移动|内网|其他 分颜色，带tip
func IPList(ipList []string, moreJumpLink ...string) string {
	var str []string
	var more []string
	//var innerIPList []string
	for _, addr := range ipList {
		ip := strings.Split(addr, ":")[0]
		loc, _ := locate.IPLocate(ip)
		if loc != nil && (strings.HasPrefix(ip, "10.") || strings.HasPrefix(ip, "192.")) {
			loc.Isp = 123
		}
		switch loc.GetIsp() {
		case 0:
			ip = GraySpan(ip)
		case 1:
			ip = GreenSpan(ip)
		case 2:
			ip = RedSpan(ip)
		case 5:
			ip = YellowSpan(ip)
		case 123:
			more = append(more, fmt.Sprintf("%s(%s)", ip, loc.GetIspFormat()))
			//innerIPList = append(innerIPList, ip)
			continue
		}
		if len(str) < 3 {
			str = append(str, TipSpan(ip, loc.GetIspFormat()))
			continue
		}
		more = append(more, fmt.Sprintf("%s(%s)", strings.Split(addr, ":")[0], loc.GetIspFormat()))
	}

	if len(more) > 0 {
		html := GraySpan("更多[" + strconv.Itoa(len(more)) + "]")
		if len(moreJumpLink) > 0 {
			html = LinkA(moreJumpLink[0], html)
		}
		str = append(str, TipSpan(html, strings.Join(more, "\n")))
	}
	//str = append(str, TipSpan(GraySpan("内网["+strconv.Itoa(len(innerIPList))+"]"), strings.Join(innerIPList, "\n")))

	return strings.Join(str, " | ")
}

func IPListTable(title string, ipList []string) *DomTable {
	table := newDomTable(title)
	table.AddHeader("IP", "区域", "省份", "城市", "运营商")

	for _, ip := range ipList {
		lo, err := locate.IPLocate(ip)
		if err != nil {
			lo = &locate.IPLocation{}
		}
		table.AddRow(ip, lo.GetZoneFormat(), lo.GetProvFormat(), lo.GetCityFormat(), lo.GetIspFormat())
	}
	table.AddNo()
	table.Sort(0, false)

	return table
}

// GreenIfSpan 如果is==true输出绿色，否则原色。
func GreenIfSpan(inner interface{}, is bool) string {
	if is {
		return GreenSpan(inner)
	}
	return fmt.Sprint(inner)
}

// YellowSpan 生成黄色字体的span
func YellowSpan(inner interface{}) string {
	return "<span style=\"color:#e2942b\">" + fmt.Sprint(inner) + "</span>"
}

// YellowIfSpan 如果is==true输出黄色，否则原色。
func YellowIfSpan(inner interface{}, is bool) string {
	if is {
		return YellowSpan(inner)
	}
	return fmt.Sprint(inner)
}

// GraySpan 输出灰色
func GraySpan(inner interface{}) string {
	return "<span style=\"color:gray\">" + fmt.Sprint(inner) + "</span>"
}

// GrayIfSpan 如果is==true输出灰色，否则原色。
func GrayIfSpan(inner interface{}, is bool) string {
	if is {
		return GraySpan(inner)
	}
	return fmt.Sprint(inner)
}

var gAutoColor = []string{"black", "green", "#e2942b", "red"}

// AutoColorSpan 自动颜色span
// lv 取值范围 [0, 100]，值越大，颜色越鲜艳
func AutoColorSpan(inner interface{}, lv int) string {
	if lv == 0 {
		return fmt.Sprintf(`<span style="color:gray">%v</span>`, inner)
	}
	if lv >= 100 {
		lv = 99
	}
	if lv < 0 {
		lv = 0
	}
	return fmt.Sprintf(`<span style="color:%s">%v</span>`, gAutoColor[int(float64(lv)/100*float64(len(gAutoColor)))], inner)
}

// GreenBackSpan 生成绿色背景的span
func GreenBackSpan(inner interface{}) string {
	return "<span style=\"background-color:#9aff8e\">" + fmt.Sprint(inner) + "</span>"
}

// LinkA 生成一个超链接
func LinkA(url, html string) string {
	return fmt.Sprintf(`<a target="_blank" href="%s">%s</a>`, url, html)
}

// LinkAWithParams 生成一个带参数的超链接
// 例：LinkAWithParams("/user", "用户", "name", "张三", "age", 13) => <a target="_blank" href="/user?name=张三&age=13">用户</a>
func LinkAWithParams(url, html string, params ...interface{}) string {
	var s []string
	for i := 1; i < len(params); i += 2 {
		s = append(s, fmt.Sprintf("%v=%v", params[i-1], params[i]))
	}
	return fmt.Sprintf(`<a target="_blank" href="%s?%s">%s</a>`, url, strings.Join(s, "&"), html)
}

// PageButton 翻页按钮
func PageButton() string {
	return `
<div class="row" style="margin-bottom: 10px;">
  	<div class="col-md-12">
		<div>
		  <button onclick="prePage()">上一页</button>
		  <button onclick="newPage()">下一页</button>
		</div>
	</div>
</div>
`
}
