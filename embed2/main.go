package main

import (
	"bytes"
	"embed"
	"html/template"
	"log"
)

//go:embed template/*.tpl
var mytpl embed.FS

func main() {
	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}
	list, err := mytpl.ReadDir("template")
	check(err)
	for _, item := range list {
		log.Println("template/" + item.Name())
	}

	data, err := mytpl.ReadFile("template/hello.tpl")
	log.Println("data:", string(data))
	check(err)
	t, err := template.New("hello").Parse(string(data))

	content := bytes.NewBufferString("")
	err = t.Execute(content, map[string]string{"name": "张三"})
	check(err)
	log.Println(content.String())
}
