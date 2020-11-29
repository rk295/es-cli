package cmd

import "github.com/jedib0t/go-pretty/table"

var (
	// Defaults
	enableColourDefault = true
	esURLDefault        = "http://localhost:9200/"

	// Flags
	byteFlag   = "byte"
	colourFlag = "colour"
	esURLFlag  = "es-url"
	sortFlag   = "sort"

	// Output
	tableStyle = table.StyleRounded
)
