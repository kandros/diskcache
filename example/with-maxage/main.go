package main

import (
	"time"

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
	dc.Set("my_data_with_max_age", aa)
	bb := &a{}
	litter.Dump(bb)
	time.Sleep(3)
	err := dc.GetIfMaxAge("mydata", bb, 1*time.Second)
	litter.Dump(err.Error()) /* expired */
}
