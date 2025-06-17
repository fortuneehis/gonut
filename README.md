# GoNut 🍩 – Lightweight HTML-like Templating Engine in Go

GoNut is a minimal, fast, and readable templating engine for Go, designed to help developers create HTML-like views with clean syntax and zero bloat.

Ideal for internal tools, microservices, or side projects where full templating engines are overkill.

## 🚀 Getting Started

### Installation

````bash
go get github.com/fortuneehis/gonut
````

### Basic Usage

````go
package main

import (
	"fmt"

	nut "github.com/fortuneehis/gonut"
)

func main() {
	template := `
# declaring a html tag
html
#head tag
    head
        style
            [
            .body {
                width: 100%;
                height: 100%;
            }
            ]
        title
            text(value="{name}")
 #This is a comment
 body
        header # header tag
        section
            br
            div(class="container section")
    text(value="")
    script(src="./index.js")
    script
        [
        console.log("yay")
        ]`

	output, err := nut.Run([]byte(template), map[string]string{
		"name": "Welcome to GoNut",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(output)
}
````

---

## 🧰 Use Cases

* Small web apps & dashboards
* Internal admin panels
* Configuration-based views
* Templates for static site generation
* Email templating

---


## 🤝 Contributing

Pull requests and issues are welcome!
If you find a bug or want a feature, feel free to open an issue or submit a PR.

---

## 📄 License

MIT License © [Fortune Ehijianbhulu](https://github.com/fortuneehis)