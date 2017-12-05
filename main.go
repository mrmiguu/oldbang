package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

func main() {
	nameExp := `[A-Za-z_][0-9A-Za-z_]*`
	Name := regexp.MustCompile(nameExp)

	Package := regexp.MustCompile(`^!` + nameExp)
	Namespace := regexp.MustCompile(`\n` + nameExp)
	Var := regexp.MustCompile(nameExp + ` : `)
	Append := regexp.MustCompile(nameExp + ` \[\]=`)

	src, err := ioutil.ReadFile("test.!")
	must(err)

	//
	// parser starts here
	//

	// package name
	pkg := Package.Find(src)
	src = Package.ReplaceAll(src, append([]byte("package "), pkg[1:]...))

	// range
	src = bytes.Replace(src, []byte(":=<"), []byte(":=range"), -1)

	// var
	src = Var.ReplaceAllFunc(src, func(match []byte) []byte {
		return append([]byte("var "), match[:len(match)-2]...)
	})

	// append
	println("[[append]]")
	src = Append.ReplaceAllFunc(src, func(match []byte) []byte {
		nm := Name.Find(match)
		return append(append([]byte{}, nm...), append(append([]byte(`=append(`), nm...), ',')...)
	})
	// println(string(src))

	// namespace
	var matches []string
	bb := Namespace.FindAll(src, -1)
	for _, b := range bb {
		matches = append(matches, string(bytes.Trim(b, "\r\n")))
	}
	fmt.Println("namespace", matches)

	dst, err := os.Create("test.go")
	must(err)
	_, err = dst.Write(src)
	must(err)
}

func must(err error) {
	if err == nil {
		return
	}
	panic(err)
}
