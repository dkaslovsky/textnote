package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dkaslovsky/TextNote/pkg/template"
)

func main() {
	date := time.Now()
	b := template.NewBody(date,
		template.NewSection("TODO",
			"- nothing",
			"- things\n  - all of them\n\n  - some of them",
			"- others\n\n\n\n",
			"- still more",
		),
		template.NewSection("DONE"),
		template.NewSection("NOTES", "- foo\n", "- bar"),
	)

	fmt.Println(b.GetFileName())
	fmt.Println("-------------------")

	err := b.Write(os.Stdout)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("-------------------")

	file := strings.NewReader(b.String())

	b2, err := template.Read(file)
	if err != nil {
		log.Fatal(err)
	}

	err = b2.Write(os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
