package degos_helpers

import (
	"flag"
	"fmt"
)

func PrintUsage() {
	fmt.Println("Usage: degos [flags] <PDF source>")
	fmt.Println("Flags:")
	flag.PrintDefaults()
}

func PrintHelp() {
	fmt.Println("Bulk Convert PDFs to Images")
	fmt.Println("Reads all PDFs in a directory 'A/B', moves each one to new directories 'A/<PDF_name>' and and parses all pages in the PDF to images which are stored in 'A/<PDF_name>/\n")
	PrintUsage()
}
