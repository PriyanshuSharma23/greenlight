package main

import (
	"flag"
	"log"
	"net/http"
)

var html = `
<!doctype html>
<html>
  <head>
    <title>Simple Cors</title>
  </head>
  <body>
    <h1>Simple CORS</h1>
    <pre id="output"></pre>
    <script>
        window.addEventListener("DOMContentLoaded", () => {
          fetch("http://localhost:4000/v1/healthcheck")
          .then(res => res.text())
          .then(data => document.getElementById("output").innerHTML=data)
          .catch(err => document.getElementById("output").innerHTML=err)
        })
    </script>
  </body>
</html>
`

func main() {
	addr := flag.String("addr", ":9000", "Server address")
	flag.Parse()
	log.Printf("starting server on %s", *addr)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(html))
	})

	log.Fatal(http.ListenAndServe(*addr, mux))
}
