package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

func main() {
	err := filepath.Walk(os.Args[1], func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != ".proto" {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}

		goPkg := fmt.Sprintf("option go_package = \"github.com/lyft/clutch/backend/%s;%sv1\";", filepath.Dir(path), filepath.Base(filepath.Dir(filepath.Dir(path))))

		var b bytes.Buffer
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			l := scanner.Bytes()
			re := regexp.MustCompile(`option go_package`)
			if len(re.Find(l)) != 0 {
				b.Write([]byte(goPkg))
			} else {
				b.Write(scanner.Bytes())
			}
			b.WriteString("\n")
		}
		f.Close()

		f, err = os.Create(path)
		if err != nil {
			return err
		}
		if _, err := b.WriteTo(f); err != nil {
			return err
		}



		//panic(err)
		//
		//
		//b, err := ioutil.ReadAll(f)
		//if err := f.Close(); err != nil {
		//	return err
		//}
		//if err != nil {
		//	return err
		//}
		//packageName := fmt.Sprintf("package clutch.%s;", strings.TrimPrefix(
		//	strings.ReplaceAll(filepath.Dir(path), "/", "."),
		//	"api."))
		//_ = packageName
		//
		//re := regexp.MustCompile(`go_package=[\w.]+;`)
		//o := re.Find(b)
		//
		//if packageName != string(o) {
		//	println(path, string(o), packageName)
		//	panic(path)
		//}
		//
		//fmt.Println("option go_package=github.com/lyft/clutch/backend/api;")
		//
		return nil
	})
	if err != nil {
		panic(err)
	}
}
