package blog

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

// Post is the datastore structure for a blog post
type Post struct {
	Title   string
	Content string
	Date    time.Time
	Author  string
	ID      int64
}

type PostContext struct {
	Title   string
	Content string
	Date    string
	Author  string
	ID      int64
}

type blog struct {
	PostInfo Post
	ID       string
}

func init() {
	http.HandleFunc("/", blogHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/portfolio", portfolioHandler)
	http.HandleFunc("/admin", adminHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/admin/post", postHandler)
	http.HandleFunc("/post", postViewHandler)
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
	postURL, _ := url.Parse(r.URL.String())
	postQuery := postURL.Query()
	postnumber := postQuery.Get("post")
	if postnumber == "" {
		t, _ := template.ParseFiles("index.html")
		c := appengine.NewContext(r)
		qry := datastore.NewQuery("Post").Order("Date")
		var postAry []Post
		keys, _ := qry.GetAll(c, &postAry)
		postContextAry := make([]PostContext, len(postAry))
		for i := range postAry {
			postContextAry[i].Content = postAry[i].Content
			if len(postAry[i].Content) > 60 {
				postContextAry[i].Content = postAry[i].Content[0:60] + " [...]"
			}
			postContextAry[i].ID = keys[i].IntID()
			postContextAry[i].Date = postAry[i].Date.Format("Jan 2, 2006")
			postContextAry[i].Author = postAry[i].Author
			postContextAry[i].Title = postAry[i].Title
		}
		t.Execute(w, postContextAry)
	} else {
		t, _ := template.ParseFiles("post.html")
		t.Execute(w, nil)
	}
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	file, _ := ioutil.ReadFile("about.html")
	fmt.Fprint(w, string(file))
}

func portfolioHandler(w http.ResponseWriter, r *http.Request) {
	file, _ := ioutil.ReadFile("portfolio.html")
	fmt.Fprint(w, string(file))
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		url, _ := user.LoginURL(c, r.URL.String())
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return
	}
	t, _ := template.ParseFiles("admin.html")
	t.Execute(w, nil)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusFound)
		return
	}
	url, _ := user.LogoutURL(c, "/")
	w.Header().Set("location", url)
	w.WriteHeader(http.StatusFound)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	post := &Post{
		Title:   r.FormValue("title"),
		Content: r.FormValue("content"),
		Date:    time.Now(),
		Author:  u.String(),
	}
	key := datastore.NewIncompleteKey(c, "Post", nil)
	datastore.Put(c, key, post)
	http.Redirect(w, r, "/admin", http.StatusFound)
}

func postViewHandler(w http.ResponseWriter, r *http.Request) {
	postURL := r.URL
	postString := postURL.Query().Get("p")
	c := appengine.NewContext(r)
	postNum, _ := strconv.ParseInt(postString, 10, 64)
	k := datastore.NewKey(c, "Post", "", postNum, nil)
	qry := datastore.NewQuery("Post").Filter("__key__ =", k).Limit(1)
	var posts []Post
	key, _ := qry.GetAll(c, &posts)
	post := posts[0]
	post.ID = key[0].IntID()
	t, _ := template.ParseFiles("post.html")
	t.Execute(w, post)
}
