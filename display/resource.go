package display

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/olekukonko/tablewriter"
	"github.com/wallix/awless/cloud/aws"
)

type ResourceDisplayer interface {
	Print(io.Writer) error
	SetResource(*aws.Resource)
}

func BuildResourceDisplayer(headers []ColumnDefinition, format string) ResourceDisplayer {
	switch format {
	case "table":
		return &tableResourceDisplayer{headers: headers}
	default:
		fmt.Fprintf(os.Stderr, "unknown displayer for %s", format)
	}
	return &tableResourceDisplayer{headers: headers}
}

type tableResourceDisplayer struct {
	r       *aws.Resource
	headers []ColumnDefinition
}

func (d *tableResourceDisplayer) Print(w io.Writer) error {
	values := make(table, len(d.r.Properties()))

	i := 0
	for prop, val := range d.r.Properties() {
		var header ColumnDefinition
		for _, h := range d.headers {
			if h.propKey() == prop {
				header = h
			}
		}
		if header == nil {
			header = StringColumnDefinition{Prop: prop}
		}

		if v := values[i]; v == nil {
			values[i] = make([]interface{}, 2)
		}
		values[i][0] = header.title(false)
		values[i][1] = header.format(val)
		i++
	}

	sort.Sort(byCols{table: values, sortBy: []int{0}})

	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"Property" + ascSymbol, "Value"})

	for i := range values {
		table.Append([]string{fmt.Sprint(values[i][0]), fmt.Sprint(values[i][1])})
	}

	table.Render()

	return nil
}

func (d *tableResourceDisplayer) SetResource(r *aws.Resource) {
	d.r = r
}
