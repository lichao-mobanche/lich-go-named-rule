# lich-go-named-rule

A common library for named rule use case.

## Install

You can get the library with ``go get``

```
go get -u github.com/lichao-mobanche/lich-go-pattern-rule
```

## Usage

```
package main

import (
	"fmt"
	"github.com/lichao-mobanche/lich-go-named-rule/cabinet"
)

func main() {
	cab := cabinet.NewCabinet()
	gro := cabinet.NewGroup("first")
	cab.LoadGroup(gro)
	fmt.Println(gro.LoadTage("1", 1))
	fmt.Println(gro.LoadTage("2", 2))
	fmt.Println(gro.LoadTage("3", 3))
	fmt.Println(gro.LoadSubTage("1", "->", "2"))
	fmt.Println("1 next tag is : ", gro.CheckNextTage("1", "->"))
	fmt.Println("gro graph is : ", gro.GroupInfo())
	gro.LoadSubTage("3", "->", "2")
	fmt.Println("gro graph is : ", gro.GroupInfo())
```

## License
  MIT licensed.