package main

import (
	"net/http"
	"html/template"
	"regexp"
	"io/ioutil"
)



//post file ending
var post_txt = regexp.MustCompile(`\.txt$`)

type Posts struct {
	Title string
	Body string
}

type Content string

func (c *Content) Write(p []byte) (n int, err error) {
	*c = *c+Content(p)
	return len(p), nil
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	var c Content
	
	t_posts, err := template.ParseFiles("content/html_templates/posts.html")
	if err != nil {
		panic(err)
	}
	t_index, err := template.ParseFiles("webbstuff/GoCrazyIndex.html")
	if err != nil {
		panic(err)
	}
	
	//Skulle vilja ha en bättre logik över vars saker ligger. 
	dir, err := ioutil.ReadDir(`./content/posts/GoCrazyIndex/`)
	if err != nil {
		panic(err)
	}
	p := make([]Posts, len(dir))
	
	//Finding the post files
	j := 0
	for i := range dir {
		file := dir[i].Name()
		if post_txt.MatchString(file) {
			s := post_txt.FindStringIndex(file)
			p[j].Title = string(file[:s[0]])
			j++
		}
	}
	
	//Reading in post files content and parse them into template
	for i:=0; i<j;i++ {
		holder, _ := ioutil.ReadFile("./content/posts/GoCrazyIndex/"+p[i].Title + ".txt")
		p[i].Body = string(holder)
		t_posts.Execute(&c, p[i])
	}
	
	t_index.Execute(w, template.HTML(c))
}
