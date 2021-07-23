# Diskcache

An easy way to write and read golang data structures to disk to be used as a cache.

## Why i built it

To allow my CLIs that load data from the network to conserve them for some times for consecutive usage.
eg. fetching a list of Namespaces from a Kubernetes cluster, this is an information that don't change very often and i can reuse the data between different stateless command usage.

## How can i prevent data that is too old?

using the method `GetIfMaxAge(string, time.Duration)` will not return the data if older than `time.Duration`

## How to use?

Check the [examples](./example/main.go) .

### Basic usage

```

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
```

#### Refuse value if stale

```
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
```

## Using

- https://github.com/spf13/afero to have a testable filesystem in memory
