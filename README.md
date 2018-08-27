[![Build Status](https://travis-ci.org/uudashr/go-module.svg?branch=master)](https://travis-ci.org/uudashr/go-module) [![GoDoc](https://godoc.org/github.com/uudashr/go-module?status.svg)](https://godoc.org/github.com/uudashr/go-module)
# go-module

This are project for fun. Parse and read go module. Read [go module](https://research.swtch.com/vgo-module)



## Usage

Install `go get -u github.com/uudashr/go-module`



## Example

Given `go.mod` file

```
module my/thing
require other/thing v1.0.2
require new/thing v2.3.4
exclude old/thing v1.2.3
replace bad/thing v1.4.5 => good/thing v1.4.5
```



```go
package main

import module "github.com/uudashr/go-module"

func main() {
    b, err := ioutil.ReadFile("go.mod")
    if err != nil {
        panic(err)
    }
    
    mod, err := module.Parse(string(b))
    if err != nil {
        panic(err)
    }
    
    _ = mod // now we have mod file
    // mod.Requires
    // mod.Excludes
    // mod.Repalces
}
```
