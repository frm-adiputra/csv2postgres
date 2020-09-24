# csv2postgres

Golang code generator for importing CSV to PostgreSQL with fun features

## Features

- [x] Importing data from CSV to PostgreSQL
- [x] Create tables based on CSV data
- [x] Create additional views
- [x] Export from PostgreSQL to CSV
- [x] Manage dependencies between tables and views
- [ ] Use dependency graph from database for rolling back migration, and from code for migrate

## Project structure

### Tables

Table specifications must be put inside `tables` directory.
The specification file must be named using format `schema_name.table_name.yaml`.
Schema name can be omitted, and if omitted, csv2postgres will use the value of
`DefaultSchema` field of `Generator` struct.

### Views

View specifications must be put inside `views` directory.
The specification file must be named using format `schema_name.view_name.yaml`.
Schema name can be omitted, and if omitted, csv2postgres will use the value of
`DefaultSchema` field of `Generator` struct.

### Setup

In the project root directory there are two files that must be created.
The first is `main.go` with content like this.

```go
package main

//go:generate go run gen.go
```

The `main.go` file will be used to setup the go generate command.

The other file is `gen.go`.
The `gen.go` contains the setup for the generator.

```go
// +build ignore
// This program generates codes for data processing.
// It must be invoked by running go generate

package main

import (
	"fmt"
	"os"

	"github.com/frm-adiputra/csv2postgres"
)

func main() {
	g := csv2postgres.Generator{
		BaseImportPath: "github.com/this/project/import/path",
		RootDir:        "generated", // default to "."
		DefaultSchema:  "my_schema", // default to "public"
	}
	if err := g.Generate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
```

## Running the project

1. Run `go generate` on your project root directory
2. Run `go run . help` on your project root to show available targets
