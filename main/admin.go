// All handlers and functions related to the admin page

package blog

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

func adminHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/admin.html", "templates/base.html"))
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
		if len(postAry[i].Title) > 20 {
			postDataAry[i].Title = postAry[i].Title[0:20] + " [...]"
		}
	}
	t.ExecuteTemplate(w, "base", postDataAry)
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
	time.Sleep(100 * time.Millisecond)
	http.Redirect(w, r, "/admin", http.StatusFound)
}

func adminEditPostHandler(w http.ResponseWriter, r *http.Request) {
	postURL := r.URL
	postString := postURL.Query().Get("p")
	c := appengine.NewContext(r)
	postNum, _ := strconv.ParseInt(postString, 10, 64)
	key := datastore.NewKey(c, "Post", "", postNum, nil)
	qry := datastore.NewQuery("Post").Filter("__key__ =", key).Limit(1)
	var posts []Post
	postKeys, _ := qry.GetAll(c, &posts)
	post := posts[0]
	post.ID = postKeys[0].IntID()
	t := template.Must(template.ParseFiles("templates/edit.html", "templates/base.html"))
	err := t.ExecuteTemplate(w, "base", post)
	if err != nil {
		fmt.Fprint(w, err)
	}
}

func submitEditHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	postURL := r.URL
	postString := postURL.Query().Get("p")
	postNum, _ := strconv.ParseInt(postString, 10, 64)
	key := datastore.NewKey(c, "Post", "", postNum, nil)
	qry := datastore.NewQuery("Post").Filter("__key__ =", key).Limit(1)
	var posts []Post
	qry.GetAll(c, &posts)
	post := posts[0]
	post.Title = r.FormValue("title")
	post.Content = r.FormValue("content")
	datastore.Put(c, key, &post)
	time.Sleep(100 * time.Millisecond)
	http.Redirect(w, r, "/admin", http.StatusFound)
}

func removePostHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	postURL := r.URL
	postString := postURL.Query().Get("p")
	postNum, _ := strconv.ParseInt(postString, 10, 64)
	key := datastore.NewKey(c, "Post", "", postNum, nil)
	datastore.Delete(c, key)
	time.Sleep(100 * time.Millisecond)
	http.Redirect(w, r, "/admin", http.StatusFound)
}
