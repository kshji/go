# Using  commandline options and arguments as \*nix persons calls or flags as Go calls #

Basic using,
this version function call need update always when you add new flag.
```go
// main.flag.ver0.go
package main

// go run main.flag.ver0.go -d ";" -f xxxx

import (
    "flag"
    "fmt"
)



func main() {
	// variables in main
        var (
		// option -f, default empty string,  and help text
                xlsxPath = flag.String("f", "", "Path to an XLSX file")
                delimiter = flag.String("d", ";", "Delimiter to use between fields")
                )
        flag.Parse()
        testflags(*xlsxPath, *delimiter)
}

func testflags(str1 string, str2 string) {
    fmt.Println("str1",str1)
    fmt.Println("str",str2)
}

```


## Using global variables in package ##
This version function call is stable, but need to use global variables.
```go
// main.flag.ver1.go
package main

// go run main.flag.ver1.go -d ";" -s sheet -f xxxx -i 1

import (
    "flag"
    "fmt"
)


// global variables in package, pointers to the data
var (
        xlsxPath = flag.String("f", "", "Path to an XLSX file")
        sheetIndex = flag.Int("i", 0, "Index of sheet to convert, zero based")
        sheetName = flag.String("s", "", "Name of sheet to convert")
        delimiter = flag.String("d", ";", "Delimiter to use between fields")
        )

func main() {
        flag.Parse()
        testflags()  // using global variables
}

func testflags() {
    // using global variables
    fmt.Println("xlsxPath",*xlsxPath)
    fmt.Println("sheetIndex",*sheetIndex)
    fmt.Println("sheetName",*sheetName)
    fmt.Println("delimiter",*delimiter)
}

```

## Using structured set of flags ##
This version function call is always stable, using struct.

```go
// main.flag.ver2.go
package main

// go run main.flag.ver2.go -d ";" -s sheet -f xxxx -i 1

import (
    "flag"
    "fmt"
)

// define struct
type MyConfig struct {
        xlsxPath string
        sheetIndex int
        sheetName string
        delimiter string
}

func main() {
        myParam := new(MyConfig)  // create object using struct
	// flag support also to set flags to the varables, not using pointers = make copy from flags data
        flag.StringVar(&myParam.xlsxPath,"f", "", "Path to an XLSX file")
        flag.IntVar(&myParam.sheetIndex,"i", 0, "Index of sheet to convert, zero based")
        flag.StringVar(&myParam.sheetName,"s", "", "Name of sheet to convert")
        flag.StringVar(&myParam.delimiter,"d", ";", "Delimiter to use between fields")
        flag.Parse()
        testflags(*myParam)  // using struct
}

func testflags(Par MyConfig) {
    fmt.Println(Par)
    fmt.Println("xlsxPath",Par.xlsxPath)
    fmt.Println("sheetIndex",Par.sheetIndex)
    fmt.Println("sheetName",Par.sheetName)
    fmt.Println("delimiter",Par.delimiter)
}
```


## Using structured set of flags ##
Include also usage example and command line include options and argument(s).

```go
// main.flag.ver3.go
package main

// go run main.flag.ver3.go -d ";" -s sheet -f xxxx -i 1
// go run main.flag.ver3.go -d ";" -s sheet -f xxxx -i 1 somefile                                                                                                                 
import (
    "flag"
    "fmt"
    "os"
)

type MyConfig struct {
        xlsxPath string
        sheetIndex int
        sheetName string
        delimiter string
}

func main() {
        myParam := new(MyConfig)
        // - options
        flag.StringVar(&myParam.xlsxPath,"f", "", "Path to an XLSX file")
        flag.IntVar(&myParam.sheetIndex,"i", 0, "Index of sheet to convert, zero based, default 0")
        flag.StringVar(&myParam.sheetName,"s", "", "Name of sheet to convert")
        flag.StringVar(&myParam.delimiter,"d", ";", "Delimiter to use between fields")

        // setup usage
        flag.Usage = func() {
                fmt.Fprintf(os.Stderr, `
%s do something ...
Usage: %s [flags] somefile

`, os.Args[0], os.Args[0])
                flag.PrintDefaults()
        }

        flag.Parse()

	// test is there argument(s), in this example need something ...
        if flag.NArg() < 1 {  // somefile is not set
                flag.Usage()
                os.Exit(1)
        }

        fmt.Println("arguments:", flag.Args())
        fmt.Println("arg 1:", flag.Arg(0))

        testflags(*myParam)  // using struct
}

func testflags(Par MyConfig) {
    fmt.Println(Par)
    fmt.Println("xlsxPath",Par.xlsxPath)
    fmt.Println("sheetIndex",Par.sheetIndex)
    fmt.Println("sheetName",Par.sheetName)
    fmt.Println("delimiter",Par.delimiter)
}
```
