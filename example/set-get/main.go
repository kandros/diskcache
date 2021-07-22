package main

import (
	"github.com/kandros/diskcache"
	"github.com/sanity-io/litter"
	"github.com/spf13/afero"
)

func main() {
	fs := afero.NewOsFs()
	dc := diskcache.New(fs, "mytest")

	type a struct {
		Name string
	}

	aa := a{Name: "jaga"}
	dc.Set("mydata", aa)
	bb := &a{} /* Name: "" */
	litter.Dump(bb)
	err := dc.Get("mydata", bb)
	if err != nil {
		litter.Dump(err.Error())
		return
	}
	litter.Dump(bb) /* Name: "jaga" */
}
