package main

import (
	"log"
	"net/http"
)

func main() {
	// listen on the root path
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		//respond HTML
		writer.Write([]byte(`
			<html>
				<head>
					<title>Chat</title>
				</head>
				<body>
					Let's chat!
				</body>
			</html>`))
	})
	// start the web server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
