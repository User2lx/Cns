//go:generate picopacker index.html index.go index
//go:generate picopacker add.html add.go add
//go:generate picopacker opensearch.template.xml opensearch.go opensearch
package main

import (
	"bufio"
	"flag"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	ttemplate "text/template"
	"time"

	"github.com/ajanicij/goduckgo/goduckgo"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	port := flag.Int("port", 8080, "http server listen port")
	u := flag.String("url", "https://example.com", "http server external url")
	bangFile := flag.String("file", "bangs.txt", "bang file name")
	flag.Parse()

	bangs := make(map[string]string, 128)

	f, err := os.OpenFile(*bangFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		panic(err)
	}
	s := bufio.NewScanner(f)
	for s.Scan() {
		l := s.Text()
		i := strings.IndexRune(l, ' ')
		if i <= 0 || i+1 == len(l) {
			continue
		}
		bangs[l[:i]] = l[i+1:]
	}
	fw := bufio.NewWriter(f)
	defer f.Close()

	bangRegex := regexp.MustCompile("!([A-Za-z0-9]*)")
	opensearchTemplate := ttemplate.Must(ttemplate.New("opensearch").Parse(string(opensearch)))
	addTemplate := template.Must(template.New("add").Parse(string(add)))

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			addTemplate.Execute(w, struct {
				Bang string `html:"bang"`
			}{
				Bang: "",
			})
			return
		}
		if r.Method != "POST" {
			w.WriteHeader(404)
			return
		}

		if err := r.ParseForm(); err != nil {
			w.WriteHeader(400)
			return
		}

		f := r.PostForm
		key := f.Get("key")
		u := f.Get("url")

		if len(key) == 0 || len(u) == 0 {
			w.WriteHeader(400)
			return
		}

		i := strings.Index(u, ":/")
		if i >= 0 && len(u) > i+2 && u[i+2] != '/' {
			u = u[:i] + "://" + u[i+2:]
		}

		bangs[key] = u
		fw.WriteString(key + " " + u + "\n")
		fw.Flush()

		addTemplate.Execute(w, struct {
			Bang string `html:"bang"`
		}{
			Bang: key,
		})
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(404)
			return
		}

		query := r.URL.Query().Get("q")

		matches := bangRegex.FindStringSubmatchIndex(query)
		if matches == nil {
			w.Header().Set("Location", "https://google.com/search?hl=en&q="+url.QueryEscape(query))
			w.WriteHeader(302)
			return
		}

		u, ok := bangs[query[matches[2]:matches[3]]]
		if !ok {
			m, err := goduckgo.Query(query)
			if err != nil || m.Redirect == "" {
				w.Header().Set("Location", "https://google.com/search?hl=en&q="+url.QueryEscape(query))
				w.WriteHeader(302)
				return
			}
			w.Header().Set("Location", m.Redirect)
			w.WriteHeader(302)
			return
		}

		query = query[:matches[0]] + query[matches[1]:]
		u = strings.ReplaceAll(u, "{{{q}}}", query)

		if !strings.Contains(u, "://") {
			u = "http://" + u
		}

		w.Header().Set("Location", u)
		w.WriteHeader(302)
	})

	http.HandleFunc("/opensearch.xml", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(404)
			return
		}
		opensearchTemplate.Execute(w, struct {
			Url string `html:"url"`
		}{
			Url: *u,
		})
	})
	http.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(404)
			return
		}
		w.Write([]byte("User-agent: *\r\nDisallow: /\r\n"))
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(404)
			return
		}
		w.Write(index)
	})
	fmt.Println("listening on localhost:" + strconv.Itoa(*port))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}
