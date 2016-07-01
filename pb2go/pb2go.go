package main

import (
	"flag"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func CopyFile(src, des string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	desFile, err := os.Create(des)
	if err != nil {
		return err
	}
	defer desFile.Close()

	_, err = io.Copy(desFile, srcFile)
	return err
}

func main() {
	var inPath string
	var outPath string

	flag.StringVar(&inPath, "in", ".", "input dir")
	flag.StringVar(&outPath, "out", ".", "input dir")
	flag.Parse()

	log.Println("start...")

	dir, err := os.Open(inPath)
	if err != nil {
		log.Fatalln("open path", inPath, "failed:", err)
	}
	defer dir.Close()

	names, err := dir.Readdirnames(0)
	if err != nil {
		log.Fatalln("read files name failed", err)
	}

	for _, name := range names {
		if !strings.HasSuffix(name, ".proto") {
			continue
		}

		goname := strings.TrimSuffix(name, ".proto") + ".pb.go"

		log.Println("work on:", name)

		pbpath := inPath + "/" + name
		tmppath := os.TempDir() + "/" + goname

		out, err := exec.Command("protoc", "--proto_path="+inPath, "--go_out="+os.TempDir(), pbpath).Output()
		if err != nil {
			log.Println("exec protoc failed:", err)
			log.Fatalln(out)
		}

		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, tmppath, nil, parser.ImportsOnly)
		if err != nil {
			log.Fatalln("go file parse failed:", err)
		}

		pkgname := f.Name.Name

		log.Println("packagename=", pkgname)

		os.Mkdir(outPath+"/"+pkgname, os.ModePerm)

		err = CopyFile(tmppath, outPath+"/"+pkgname+"/"+goname)
		if err != nil {
			log.Fatalln("copy failed:", err)
		}
	}

	log.Println("end")
}
