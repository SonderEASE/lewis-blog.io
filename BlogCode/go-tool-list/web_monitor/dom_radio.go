package web_monitor

import (
	"fmt"
	"strconv"
	"strings"
)

// DomRadio 单选框
type DomRadio struct {
	page *Page
	id   string

	name  string
	key   string
	items []*stRadioItem
}

func (p *DomRadio) Script() string {
	return ""
}

type stRadioItem struct {
	html      string
	value     string
	isDefault bool
}

func newRadio(page *Page) *DomRadio {
	p := &DomRadio{
		page: page,
		id:   "id-radio-" + page.newID(),
	}

	return p
}

func (p *DomRadio) AddItem(html, value string, isDefault bool) {
	p.items = append(p.items, &stRadioItem{html: html, value: value, isDefault: isDefault})
}

func (p *DomRadio) SetDefault(value string) {
	for _, v := range p.items {
		if v.value == value {
			v.isDefault = true
			return
		}
	}
}

func (p *DomRadio) Html() string {
	var ss []string
	for i, item := range p.items {
		id := p.id + "-" + strconv.Itoa(i)
		checked := ""
		if item.isDefault {
			checked = "checked"
		}
		ss = append(ss, fmt.Sprintf(`
<div class="custom-control custom-radio custom-control-inline">
  <input class="custom-control-input" type="radio" name="`+p.id+`" id="`+id+`" value="`+item.value+`" `+checked+`>
  <label class="custom-control-label" for="`+id+`">
    `+item.html+`
  </label>
</div>
`))
	}

	panel := DomPanel{page: p.page}
	panel.SetHtml(`
<form class="form-inline">
  <div class="form-row align-items-center c-radio" id="` + p.id + `" name="` + p.key + `">
    <label>` + p.name + `：</label>
      ` + strings.Join(ss, "\n") + `
  </div>
</form>
`)
	return panel.Html()
}
