package main

import (
	"log"

	"github.com/dkaslovsky/TextNote/cmd"
)

// func main() {
// 	date := time.Now()
// 	b := template.NewBody(date,
// 		template.NewSection("TODO",
// 			"- nothing",
// 			"- things\n  - all of them\n\n  - some of them",
// 			"- others\n\n\n\n",
// 			"- still more",
// 		),
// 		template.NewSection("DONE"),
// 		template.NewSection("NOTES", "- foo\n", "- bar"),
// 	)

// 	fmt.Println(b.GetFileName())
// 	fmt.Println("-------------------")

// 	err := b.Write(os.Stdout)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("-------------------")

// 	file := strings.NewReader(b.String())

// 	b2, err := template.Read(file)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	err = b2.Write(os.Stdout)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

func main() {
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
