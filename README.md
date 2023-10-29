# lvm2

A Go wrapper for working with [LVM2](https://sourceware.org/lvm2).

## Usage

For detailed examples see the [lvm2_test.go](./lvm2_test.go) file.

```go
package main

import (
    "context"
    "log"

    "github.com/dpeckett/lvm2"
)

func main() {
    c := lvm2.NewClient()

    ctx := context.Background()

    // List all volume groups
    vgs, err := c.ListVolumeGroups(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }

    for _, vg := range vgs {
        log.Printf("Found Volume Group: %s", vg.Name)
    }
}
```

## Contributing

Pull requests are more than welcome. 

The LVM2 project has a huge surface area and I'm sure I've missed things. If you find a bug or something that is missing please open an issue or submit a pull request.