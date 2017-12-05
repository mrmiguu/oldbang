package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/WedgeNix/warn"
)

func main() {
	exName := `[A-Za-z_][0-9A-Za-z_]*`

	src, err := ioutil.ReadFile("test.!")
	must(err)

	//
	// parser starts here
	//

	println(string(src))
	warn.Do("wipe comments")

	// comments
	Comment := regexp.MustCompile(`//.*\n`)
	src = Comment.ReplaceAll(src, []byte{'\n'})

	println(string(src))
	warn.Do("clear empty lines")

	// empty lines
	Empty := regexp.MustCompile(`[\r\n]{2,}`)
	src = Empty.ReplaceAll(src, []byte{'\n'})

	println(string(src))
	warn.Do("change bangs to returns")

	// return statements
	Bang := regexp.MustCompile(`[^\s]*! ?\n`)
	src = Bang.ReplaceAllFunc(src, func(match []byte) []byte {
		i := bytes.LastIndex(match, []byte{'!'})
		return append([]byte("return "), append(append([]byte{}, match[:i]...), match[i+1:]...)...)
	})

	println(string(src))
	warn.Do("close functions")

	// close namespace functions
	Close := regexp.MustCompile(`\n\s.*(\n\s.*)*`)
	src = Close.ReplaceAllFunc(src, func(match []byte) []byte {
		return append(append([]byte{'{'}, match...), '}')
	})

	println(string(src))
	warn.Do("fix function signatures")

	// function signatures
	Func := regexp.MustCompile(exName + ` ?::.*\n`)
	src = Func.ReplaceAllFunc(src, func(match []byte) []byte {
		match = bytes.Replace(match, []byte("::"), []byte{'('}, 1)
		match = append(match[:len(match)-2], append([]byte{')'}, match[len(match)-2:]...)...)
		match = bytes.Replace(match, []byte{';'}, []byte(")("), 1)
		return append([]byte("func "), match...)
	})

	println(string(src))
	warn.Do("fix package name")

	// package name
	Package := regexp.MustCompile(`^!` + exName)
	pkg := Package.Find(src)
	src = Package.ReplaceAll(src, append([]byte("package "), pkg[1:]...))

	println(string(src))
	warn.Do("convert for-range(s)")

	// range
	src = bytes.Replace(src, []byte(":=@"), []byte(":=range"), -1)

	println(string(src))
	warn.Do("convert var(s)")

	// var
	Var := regexp.MustCompile(exName + ` :`)
	src = Var.ReplaceAllFunc(src, func(match []byte) []byte {
		return append([]byte("var "), match[:len(match)-1]...)
	})

	println(string(src))
	warn.Do("fix append shorthands")

	// append
	Append := regexp.MustCompile(exName + ` \[\]=`)
	Name := regexp.MustCompile(exName)
	src = Append.ReplaceAllFunc(src, func(match []byte) []byte {
		nm := Name.Find(match)
		return append(append([]byte{}, nm...), append(append([]byte(` = append `), nm...), []byte(", ")...)...)
	})

	println(string(src))
	warn.Do("compile bang=>go")

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
