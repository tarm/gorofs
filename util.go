package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func testFile(name string) (err os.Error) {
	f, err := os.Open(name)
	if err != nil {
		return
	}
	fi, err := f.Stat()
	if err != nil {
		return
	}
	r, err := zip.NewReader(f, fi.Size)
	fmt.Println(fi.Size)
	for _, f := range r.File {
		fmt.Println(f.FileHeader.Name)
	}
	return
}

func main() {
	pkgName := flag.String("pkg", "main", "package name")
	varName := flag.String("var", "rofs", "variable name")
	outName := flag.String("out", "rofs.go", "file name")
	srcName := flag.String("src", "src.zip", "Src zip file")
	flag.Parse()
	err := testFile(*srcName)
	if err != nil {
		log.Fatal(err)
	}
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
	fmt.Fprintf(out, `package %s

var %s = [...]byte{`, *pkgName, *varName)
	for i, b := range buf {
		if i%8 == 0 {
			if i == 0 {
				fmt.Fprintf(out, "0x%02x", b)
			} else {
				fmt.Fprintf(out, ",\n\t0x%02x", b)
			}
		} else {
			fmt.Fprintf(out, ", 0x%02x", b)
		}
	}
	fmt.Fprintf(out, "}\n")
	if err != nil {
		log.Fatal(err)
	}
	out.Close()
}
