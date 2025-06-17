package main

import (
	"fmt"

	nut "github.com/fortuneehis/gonut"
)

func main() {
	template := `# declaring a html tag
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
