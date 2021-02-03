package got

import (
	"fmt"
	"log"
	"net/http"
)

var formStr = `<form action="/add" method="post">
URL: <input type="text" name="url">
<input type="submit" value="提交">
</form>`

var store *Store

func init() {
	store = NewStore()
}

func Add(w http.ResponseWriter, r *http.Request) {
	url := r.PostFormValue("url")
	if url == "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, formStr)
	} else {
		short, err := store.Add(url)
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}
		fmt.Fprint(w, short)
	}
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	path, err := store.Get(r.URL.Path[1:])
	log.Println("path", path)
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	http.Redirect(w, r, path, http.StatusFound)
}
