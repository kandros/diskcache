# Diskcache

An easy way to write and read golang data structures to disk to be used as a cache.

## Why i built it

To allow my CLIs that load data from the network to conserve them for some times for consecutive usage.
eg. fetching a list of Namespaces from a Kubernetes cluster, this is an information that don't change very often and i can reuse the data between different stateless command usage.

## How can i prevent data that is too old?

using the method `GetIfMaxAge(string, time.Duration)` will not return the data if older than `time.Duration`

## How to use?

Check the [examples](./example/main.go) .

## Using

- https://github.com/spf13/afero to have a testable filesystem in memory
