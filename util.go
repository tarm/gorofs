package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	pkgName := flag.String("pkg", "main", "package name")
	varName := flag.String("var", "rofs", "variable name")
	outName := flag.String("out", "rofs.go", "file name")
	srcName := flag.String("src", "src.zip", "Src zip file")
	flag.Parse()
	src, err := os.Open(*srcName)
	if err != nil {
		log.Fatal(err)
	}
	buf, err := ioutil.ReadAll(src)
	if err != nil {
		log.Fatal(err)
	}
	src.Close()
	out, err := os.Create(*outName)
	if err != nil {
		log.Fatal(err)
	}
	_, err = fmt.Fprintf(out, `package %s

const %s = %q
`, *pkgName, *varName, buf)
	if err != nil {
		log.Fatal(err)
	}
	out.Close()
}
