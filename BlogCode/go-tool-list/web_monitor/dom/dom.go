package dom

import "strings"

type Base struct {
	style []string
	class []string
	id    string
}

func (d *Base) AddStyle(style ...string) *Base {
	if d == nil {
		return nil
	}
	d.style = append(d.style, style...)
	return d
}

func (d *Base) AddClass(class ...string) *Base {
	if d == nil {
		return nil
	}
	d.class = append(d.class, class...)
	return d
}

func (d *Base) SetID(id string) *Base {
	if d == nil {
		return nil
	}
	d.id = id
	return d
}

func (d *Base) StyleHtml() string {
	if d == nil || len(d.style) == 0 {
		return ""
	}

	return `style="` + strings.Join(d.style, ";") + `"`
}

func (d *Base) IdHtml() string {
	if d == nil || d.id == "" {
		return ""
	}
	return `id="` + d.id + `"`
}

type Interface interface {
	Html() string
}

type Html struct {
	inner string
}

func (h *Html) Html() string {
	if h == nil {
		return ""
	}
	return h.inner
}
