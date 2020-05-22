package web_monitor

import (
	"fmt"
	"strings"
	"time"

	"wx-gitlab.xunlei.cn/scdn/x/conv"
)

type DomLine struct {
	id     string
	title  string
	width  string
	height string
	x      interface{}
	yAxis  []*LineOptionYAxis
	series []*LineOptionSeries
}

func (d *DomLine) SetX(x interface{}) *DomLine {
	d.x = x
	return d
}
func (d *DomLine) SetYAxis(yAxis []*LineOptionYAxis) *DomLine {
	d.yAxis = yAxis
	return d
}
func (d *DomLine) SetSeries(series ...*LineOptionSeries) *DomLine {
	d.series = series
	return d
}

func (d *DomLine) Html() string {
	if d == nil {
		return ""
	}
	return fmt.Sprintf(
		`<div id="%s" style="width:%s;height:%s;border: 1px solid #e2e2e2;"></div>`,
		d.id, d.width, d.height,
	)

}

func (d *DomLine) Script() string {
	if d == nil {
		return ""
	}
	option := lineOption{
		Title:  lineOptionTitle{Text: d.title},
		Legend: lineOptionLegend{Data: nil},
		Toolbox: lineOptionToolbox{Feature: lineOptionToolboxFeature{
			DataZoom:    lineOptionToolboxFeatureDataZoom{YAxisIndex: "none"},
			SaveAsImage: lineOptionToolboxFeatureSaveAsImage{Name: d.title + time.Now().Format("_20060102_150405")},
			MagicType:   lineOptionToolboxFeatureMagicType{Type: []string{"line", "stack", "bar", "tiled"}},
			Restore:     lineOptionToolboxFeatureRestore{Show: true},
		}},
		Tooltip:  lineOptionTooltip{Trigger: "axis", AxisPointer: lineOptionTooltipAxisPointer{Type: "cross"}, Formatter: "${Tooltip.Formatter}"},
		XAxis:    lineOptionXAxis{Data: d.x, SplitLine: &LineOptionSplitLine{Show: true}}, // 行
		YAxis:    d.yAxis,
		Series:   d.series,
		DataZoom: []*lineOptionDataZoom{{Start: 0, End: 100}},
	}

	for _, v := range d.yAxis {
		if v.Type == "" {
			v.Type = "value"
		}
		if v.AxisLabel.Formatter == "" {
			v.AxisLabel.Formatter = "{value}"
		}
	}

	for _, v := range d.series {
		if v.Type == "" {
			v.Type = "line"
		}
	}

	b := conv.ToJsonString(option)
	b = strings.ReplaceAll(b, `"${Tooltip.Formatter}"`, option.MakeTooltipFormatter())
	return fmt.Sprintf(`echarts.init(document.getElementById('%s'), 'vintage').setOption(%s);`, d.id, b)
}

type lineOptionTitle struct {
	Text string `json:"text"`
}

type lineOptionToolboxFeatureSaveAsImage struct {
	Name string `json:"name"`
}

type lineOptionToolboxFeatureDataZoom struct {
	YAxisIndex string `json:"yAxisIndex"`
}

type lineOptionToolboxFeatureMagicType struct {
	Type []string `json:"type"`
}

type lineOptionToolboxFeatureRestore struct {
	Show bool `json:"show"`
}

type lineOptionToolboxFeature struct {
	SaveAsImage lineOptionToolboxFeatureSaveAsImage `json:"saveAsImage"`
	DataZoom    lineOptionToolboxFeatureDataZoom    `json:"dataZoom"`
	MagicType   lineOptionToolboxFeatureMagicType   `json:"magicType"`
	Restore     lineOptionToolboxFeatureRestore     `json:"restore"`
}

type lineOptionToolbox struct {
	Feature lineOptionToolboxFeature `json:"feature"`
}

type lineOptionTooltipAxisPointer struct {
	Type string `json:"type"`
}

type lineOptionTooltip struct {
	Trigger     string                       `json:"trigger"`
	AxisPointer lineOptionTooltipAxisPointer `json:"axisPointer"`
	Formatter   string                       `json:"formatter"` // '{a0}:{c0}万'
}

type lineOptionXAxis struct {
	Data      interface{}          `json:"data"`
	SplitLine *LineOptionSplitLine `json:"splitLine,omitempty"` // 坐标轴在 grid 区域中的分隔线
}

type LineOptionAxisLabel struct {
	Formatter string `json:"formatter"` // {value} MB/s
}

// 坐标轴在 grid 区域中的分隔线
type LineOptionSplitLine struct {
	Show bool `json:"show"` // 默认true
}

type LineOptionYAxis struct {
	Type      string               `json:"type,omitempty"` //
	Name      string               `json:"name,omitempty"`
	Position  string               `json:"position,omitempty"` // left | right
	Offset    int                  `json:"offset"`             // x轴偏移量
	Min       float64              `json:"min,omitempty"`
	Max       float64              `json:"max,omitempty"`
	AxisLabel *LineOptionAxisLabel `json:"axisLabel,omitempty"`
	SplitLine *LineOptionSplitLine `json:"splitLine,omitempty"` // 坐标轴在 grid 区域中的分隔线
}

type LineOptionSeriesMarkPointData struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type LineOptionSeriesMarkPoint struct {
	Data []LineOptionSeriesMarkPointData `json:"data"`
}

type LineOptionSeries struct {
	Name       string                     `json:"name"`
	Type       string                     `json:"type,omitempty"`
	Data       interface{}                `json:"data"`
	MarkPoint  *LineOptionSeriesMarkPoint `json:"markPoint,omitempty"`
	YAxisIndex int                        `json:"yAxisIndex,omitempty"` // y轴数组下标
	Smooth     bool                       `json:"smooth,omitempty"`     // 是否平滑曲线
}

func NewLineOptionSeries(name string, data interface{}) *LineOptionSeries {
	return &LineOptionSeries{Name: name, Data: data}
}

// WithMaxMinPoint 设置最大最小值的锚点
func (l *LineOptionSeries) WithMaxMinPoint() *LineOptionSeries {
	l.MarkPoint = &LineOptionSeriesMarkPoint{
		Data: []LineOptionSeriesMarkPointData{
			{Type: "max", Name: "最大值"},
			{Type: "min", Name: "最小值"},
		},
	}
	return l
}

// WithYAxisIndex 设置Y周下标(对多个Y周有效，从0开始)
func (l *LineOptionSeries) WithYAxisIndex(index int) *LineOptionSeries {
	l.YAxisIndex = index
	return l
}

type lineOptionLegend struct {
	Data []string `json:"data"`
}

type lineOptionDataZoom struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

type lineOption struct {
	Title     lineOptionTitle       `json:"title"`
	Legend    lineOptionLegend      `json:"legend"`
	Toolbox   lineOptionToolbox     `json:"toolbox"`
	Tooltip   lineOptionTooltip     `json:"tooltip"`
	XAxis     lineOptionXAxis       `json:"xAxis"`
	YAxis     []*LineOptionYAxis    `json:"yAxis"`
	Series    []*LineOptionSeries   `json:"series"`
	DataZoom  []*lineOptionDataZoom `json:"dataZoom"`
	Animation bool                  `json:"animation"` // 是否开启动画 默认不开启
}

// GetLineUnit 获取曲线单位
func (l *lineOption) GetLineUnit(xIndex int) string {
	if xIndex < 0 || xIndex >= len(l.Series) {
		return ""
	}
	yIndex := l.Series[xIndex].YAxisIndex
	if yIndex < 0 || yIndex >= len(l.YAxis) {
		return ""
	}
	axisLabel := l.YAxis[yIndex].AxisLabel
	if axisLabel == nil {
		return ""
	}

	return strings.Split(axisLabel.Formatter, " ")[1]
}

func (l *lineOption) MakeTooltipFormatter() string {
	m := map[string]string{}
	for i, series := range l.Series {
		m[series.Name] = l.GetLineUnit(i)
	}
	return fmt.Sprintf(`
function(params) {
var result = params[0].name+"<br>";
const unit=%s;
params.forEach(function (item) {
    result+=item.marker+" "+item.seriesName+" "+item.value+unit[item.seriesName]+"</br>";
});
return result;}`, conv.ToJsonString(m))
}
