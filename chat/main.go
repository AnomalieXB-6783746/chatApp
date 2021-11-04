package main

import (
	"flag"
	"github.com/AnomalieXB-6783746/chatApp/trace"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"
)

// templ represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTP handles the HTTP request.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates",
			t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse() // parse the flags
	// setup gomniauth
	gomniauth.SetSecurityKey("DLw6wKwh25Go8y587RjR")
	gomniauth.WithProviders(
		/*		facebook.New("key", "secret",
				"http://localhost:8080/auth/callback/facebook"),*/
		github.New("582e9464e51bb3301734", "f3cff97c344a76a05fdb228b1384c08e86956e2e",
			"http://localhost:8080/auth/callback/github"),
		google.New("354093525754-nre11d6a9aicq2u4p7oj27lb7d6ug1j1.apps.googleusercontent.com",
			"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",
			"http://localhost:8080/auth/callback/google"),
	)
	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	// listen on the root path
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	go r.run()
	// start the web server
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
