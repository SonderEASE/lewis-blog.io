package dom

import (
	"fmt"
	"sort"
	"strconv"

	"wx-gitlab.xunlei.cn/scdn/x/conv"
)

// Table 表格dom
type Table struct {
	Name      string
	Header    []interface{}
	Rows      [][]interface{}
	No        bool // 是否显示编号
	maxHeight int  // 表格最大高度，默认没有，超过最大高度时显示滚动条。

	tdStyle map[int]string // col => string
}

func NewTable(name string) *Table {
	return &Table{
		Name:    name,
		tdStyle: map[int]string{},
	}
}

// AddHeader 添加表头
func (t *Table) AddHeader(row ...interface{}) {
	t.Header = row
}

// AddRow 添加表行
func (t *Table) AddRow(row ...interface{}) {
	t.Rows = append(t.Rows, row)
}

// Sort 排序
// 参数1：列index，从0开始
// 参数2：是否从大到小排序
func (t *Table) Sort(colIndex int, desc bool) {
	if len(t.Rows) == 0 || colIndex >= len(t.Rows[0]) {
		return
	}

	switch t.Rows[0][colIndex].(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		sort.Slice(t.Rows, func(i, j int) bool {
			if len(t.Rows[i]) <= colIndex {
				return desc
			}
			if len(t.Rows[j]) <= colIndex {
				return desc
			}
			return desc != (conv.ToFloat64(t.Rows[i][colIndex]) < conv.ToFloat64(t.Rows[j][colIndex]))
		})
	default:
		sort.Slice(t.Rows, func(i, j int) bool {
			if len(t.Rows[i]) <= colIndex {
				return desc
			}
			if len(t.Rows[j]) <= colIndex {
				return desc
			}
			return desc != (fmt.Sprint(t.Rows[i][colIndex]) < fmt.Sprint(t.Rows[j][colIndex]))
		})
	}
}

// MultSort 多重排序
// 参数 排序index和是否逆序(index从0开始)
// 例： MultSort(1, true, 5, false) 按第2列逆序，第一列相同的，按第6列顺序排序
func (t *Table) MultSort(ids ...interface{}) {
	if t == nil || len(t.Rows) == 0 || len(ids)%2 != 0 {
		return
	}

	sort.Slice(t.Rows, func(i, j int) bool {
		a := t.Rows[i]
		b := t.Rows[j]

		for m := 0; m < len(ids); m += 2 {
			index := ids[m].(int)
			desc := ids[m+1].(bool)
			if len(a) <= index || len(b) <= index {
				return false
			}

			rolA, okA := t.toFloat64(a[index])
			rolB, okB := t.toFloat64(b[index])

			// 两个都是数字
			if okA && okB {
				if rolA == rolB {
					// 值相同，按下一个条件过滤
					continue
				}
				return desc != (rolA < rolB)
			}

			// 按字符串比较
			colAstr := fmt.Sprint(a[index])
			colBstr := fmt.Sprint(b[index])
			if colAstr != colBstr {
				return desc != (colAstr < colBstr)
			}
		}

		return false
	})

}

func (t *Table) toFloat64(v interface{}) (float64, bool) {
	if v == nil {
		return 0, false
	}
	str := ""
	switch tv := v.(type) {
	case []byte:
		str = string(tv)
	case int:
		return float64(tv), true
	case int32:
		return float64(tv), true
	case int64:
		return float64(tv), true
	case float32:
		return float64(tv), true
	case float64:
		return tv, true
	default:
		str = fmt.Sprint(v)
	}
	ret, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, false
	}
	return ret, true
}

// AddNo 添加编号列
// 编号会添加在第一列，与sort的列index不冲突。
func (t *Table) AddNo() {
	t.No = true
}

func (t *Table) rowHTML(index int, row []interface{}, ti string) string {
	if len(row) == 0 {
		return ""
	}
	html := "<tr>"
	if t.No {
		if ti == "th" {
			html += fmt.Sprintf("<%s onclick=\"SortTable(this)\" style=\"text-align:center\">#</%s>", ti, ti)
		} else {
			html += fmt.Sprintf("<%s>%d</%s>", ti, index+1, ti)
		}
	}
	for i, item := range row {
		s := ""
		switch item.(type) {
		case int, int32, int64:
			s = fmt.Sprintf("%d", item)
		case float32, float64:
			s = fmt.Sprintf("%.2f", item)
		default:
			s = fmt.Sprintf("%v", item)
		}
		if ti == "th" {
			html += "<" + ti + " onclick=\"SortTable(this)\" style=\"text-align:center\">" + fmt.Sprint(s) + "</" + ti + ">"
		} else {
			html += fmt.Sprintf(`<td style="%s">%s</td>`, t.tdStyle[i], fmt.Sprint(s))
		}
	}
	html += "</tr>"
	return html
}

// TextAlignLeft 设置文字对齐方式：左对齐
// 默认右对齐
func (t *Table) TextAlignLeft(cols ...int) {
	for _, col := range cols {
		t.tdStyle[col] = t.tdStyle[col] + "text-align:left;"
	}
}

// TextAlignLeft 设置文字对齐方式：居中
// 默认右对齐
func (t *Table) TextAlignCenter(cols ...int) {
	for _, col := range cols {
		t.tdStyle[col] = t.tdStyle[col] + "text-align:center;"
	}
}

// Width 设置列宽度
// 支持多组设置，列宽成对出现。
// 例 table.Width(0, 300, 2, 100) 第0列宽300px 第2列宽100px 其他自动
func (t *Table) Width(colsAndWidth ...int) {
	for i := 0; i < len(colsAndWidth)-1; i += 2 {
		t.tdStyle[colsAndWidth[i]] = t.tdStyle[colsAndWidth[i]] + fmt.Sprintf("width:%dpx;", colsAndWidth[i+1])
	}
}

func (t *Table) MaxWidth(colsAndMaxWidth ...int) {
	for i := 0; i < len(colsAndMaxWidth)-1; i += 2 {
		t.tdStyle[colsAndMaxWidth[i]] = t.tdStyle[colsAndMaxWidth[i]] + fmt.Sprintf("max-width:%dpx;", colsAndMaxWidth[i+1])
	}
}

// maxHeight 表格最大高度，默认没有，超过最大高度时显示滚动条。
func (t *Table) MaxHeight(height int) {
	t.maxHeight = height
}

func (t *Table) Html() string {
	rowHTML := `<thead>` + t.rowHTML(0, t.Header, "th") + "</thead><tbody>"
	for index, row := range t.Rows {
		rowHTML += t.rowHTML(index, row, "td")
	}
	rowHTML += "</tbody>"

	tableHTML := fmt.Sprintf(`
<table class="table table-striped table-bordered table-hover table-condensed table-sm" style="margin-bottom:0px">
	%s
</table>
`, rowHTML)

	var maxHeightCSS string
	if t.maxHeight > 0 {
		maxHeightCSS = fmt.Sprintf(`max-height: %dpx;overflow-y: auto;`, t.maxHeight)
	}

	return fmt.Sprintf(`
<div class="row" style="padding-left:10px;padding-bottom: 10px;">
    <div class="card bg-primary text-white">
      <div class="card-header" style="padding: 5px 10px;font-size: 16px;">%s<span style="float:right">[%d]</span></div>
      <div class="card-body" style="padding: 0px; background-color: white;">
        <div style="%s">
            %s
          </div>
      </div>
    </div>
</div>
`, t.Name, len(t.Rows), maxHeightCSS, tableHTML)
}

func (t *Table) Script() string {
	return ""
}
