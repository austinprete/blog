package blog

import (
	"fmt"
	"html/template"
	"net/http"
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

// PostData is the view representation of a post.
type PostData struct {
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
	http.HandleFunc("/admin/edit", adminEditPostHandler)
	http.HandleFunc("/admin/edit/submit", submitEditHandler)
	http.HandleFunc("/post", postViewHandler)
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/index.html", "templates/base.html"))
	c := appengine.NewContext(r)
	qry := datastore.NewQuery("Post").Order("-Date")
	var postAry []Post
	keys, _ := qry.GetAll(c, &postAry)
	postDataAry := make([]PostData, len(postAry))
	for i := range postAry {
		postDataAry[i].Content = postAry[i].Content
		if len(postAry[i].Content) > 60 {
			postDataAry[i].Content = postAry[i].Content[0:60] + " [...]"
		}
		postDataAry[i].ID = keys[i].IntID()
		postDataAry[i].Date = postAry[i].Date.Format("Jan 2, 2006")
		postDataAry[i].Author = postAry[i].Author
		postDataAry[i].Title = postAry[i].Title
	}
	t.ExecuteTemplate(w, "base", postDataAry)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/about.html", "templates/base.html"))
	t.ExecuteTemplate(w, "base", nil)
}

func portfolioHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/portfolio.html", "templates/base.html"))
	err := t.ExecuteTemplate(w, "base", nil)
	if err != nil {
		fmt.Fprint(w, err)
	}
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
	//u := user.Current(c)
	post := &Post{
		Title:   r.FormValue("title"),
		Content: r.FormValue("content"),
		Date:    time.Now(),
		Author:  "Austin Prete", // u.String() if you'd like to access from logged in user
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
	t := template.Must(template.ParseFiles("templates/post.html", "templates/base.html"))
	t.ExecuteTemplate(w, "base", post)
}
