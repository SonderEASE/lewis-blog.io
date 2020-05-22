package dom

import (
	"fmt"
	"strings"
)

type Container struct {
	Base
	rows []*Row
}

func NewContainer() *Container {
	c := &Container{}

	return c
}

func (c *Container) AddRow() *Row {
	r := &Row{
		container: c,
	}

	c.rows = append(c.rows, r)
	return r
}

func (c *Container) Html() string {
	html := fmt.Sprintf(`<div class="container" %s>`, c.StyleHtml())
	for _, row := range c.rows {
		html += row.Html()
	}
	html += `</div>`

	return html
}

func (c *Container) Script() string {
	return ""
}

type Row struct {
	Base
	container *Container

	cols []*Col
}

func (r *Row) AddCol() *Col {
	col := &Col{
		row: r,
	}

	r.cols = append(r.cols, col)
	return col
}

func (r *Row) Html() string {
	var id string
	if r.id != "" {
		id = fmt.Sprintf(`id="%s"`, r.id)
	}
	html := fmt.Sprintf(`<div class="%s" %s style="%s">`, strings.Join(append(r.class, "row"), " "), id, strings.Join(r.style, ";"))
	for _, col := range r.cols {
		html += col.Html()
	}
	html += `</div>`

	return html
}

type Col struct {
	Base
	row *Row

	suffix []interface{}
	inner  []Interface
}

// SetColClassSuffix 设置 col-* 后缀
// .SetSuffix("sm", "lg-2", 5, "") => class="col-sm col-lg-2 col-5 col"
// 默认 class="col"
func (c *Col) SetColClassSuffix(suffix ...interface{}) *Col {
	if c == nil {
		return nil
	}
	c.suffix = suffix
	return c
}

func (c *Col) SetInner(inner string) *Col {
	if c == nil {
		return nil
	}
	c.inner = []Interface{&Html{inner: inner}}
	return c
}

func (c *Col) SetInnerH5(inner string) *Col {
	return c.SetInner(fmt.Sprintf(`<h5>%s</h5>`, inner))
}

func (c *Col) Class() []string {
	ss := []string{"col"}
	if len(c.suffix) > 0 {
		ss = []string{}
		for _, s := range c.suffix {
			if s == "" {
				ss = append(ss, "col")
				continue
			}
			ss = append(ss, fmt.Sprintf("col-%v", s))
		}
	}

	return ss
}

func (c *Col) Html() string {
	if c == nil {
		return ""
	}
	var ss []string
	for _, in := range c.inner {
		ss = append(ss, in.Html())
	}
	return fmt.Sprintf(`<div class="%s" %s %s>%s</div>`,
		strings.Join(append(c.Class(), c.class...), " "),
		c.IdHtml(),
		c.StyleHtml(),
		strings.Join(ss, ""),
	)
}

func (c *Col) AddTable(name string) *Table {
	table := NewTable(name)
	c.inner = append(c.inner, table)
	return table
}
