package main

import (
	"net/http"
	"html/template"
	"regexp"
	"io/ioutil"
	"os"
	"io"
	"bufio"
	"time"
)


var post_txt = regexp.MustCompile(`\.txt$`)	//post file ending
var files = "content/root/"					//Files related to this page

type Posts struct {
	Title string
	Body string
	CreatTime time.Time
}

type Content string

func (c *Content) Write(p []byte) (n int, err error) {
	*c = *c+Content(p)
	return len(p), nil
}

//Time stamp related =======#
func firstNil(b []byte) int {
	var i int
	for i = len(b)-1; ; i-- {
		if b[i] != 0 {
			break;
		}
	}
	return i+1
}

func stamp(f *os.File) {
	f.Seek(0,0)			// Making sure we start correctly. 
	f_r := bufio.NewReader(f)
	f_w := bufio.NewWriter(f)

	buf := make([]byte, 1024)
	for {
		n, err := f_r.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break;
		}
	}
	
	// Stamp the buffer
	t := time.Now().Format(time.UnixDate)
	t += "\n"
	buf = append([]byte(t), buf...)
	
	// Write the buffer
	f.Seek(0,0) //Reset the seek after all the reads.
	n := firstNil(buf)	//Making sure not to write nil bytes. 
	_, err := f_w.Write(buf[:n])
	if err != nil {
		panic(err)
	}
	err = f_w.Flush()
	if err != nil {
		panic(err)
	}
	
}

func isStamped(f *os.File) bool {
	f_r := bufio.NewReader(f)
	line, err := f_r.ReadSlice('\n')
	if err != nil && err != io.EOF {
		panic(err)
	}
	_, err = time.Parse(time.UnixDate, string(line[:len(line)-1])) // removing '\n'
	if err != nil {
		return false
	}
	return true

}


func postCreateTime(f *os.File) time.Time {
	f.Seek(0,0) // Make sure we're at the start
	f_r := bufio.NewReader(f)
	line, err := f_r.ReadSlice('\n')
	if err != nil && err != io.EOF {
		panic(err)
	}
	
	t, err := time.Parse(time.UnixDate, string(line[:len(line)-1])) // removing '\n'
	if err != nil {
		panic(err)
	}
	return t
}
//==================#

func rootHandler(w http.ResponseWriter, r *http.Request) {
	var c Content
	
	t_posts, err := template.ParseFiles(files+"posts.html")
	if err != nil {
		panic(err)
	}
	t_index, err := template.ParseFiles("main_web/main.html")
	if err != nil {
		panic(err)
	}
	
	dir, err := ioutil.ReadDir(files+"posts/")
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
		f_name, err := (files+"posts/"+p[i].Title + ".txt")
		if err != nil {
			panic(err)
		}
		f, err := os.Open(files+"posts/"+p[i].Title + ".txt", os.O_RDWR, 0644)
		if !isStamped(f){
			stamp(f)
		}
		defer f.Close()
		p[i].CreateTime = postCreateTime(f)
	}
	
	for i:=0;i<j;i++{
		//Implement quicksort ='(
	}
	
	t_index.Execute(w, template.HTML(c))
}
