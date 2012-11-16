package main

import (
	"bufio"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sort"
	"time"
)

var post_txt = regexp.MustCompile(`\.txt$`) //post file ending
var files = "content/root/"                 //Files related to this page

type Post struct {
	Title     string
	Body      string
	CreatTime time.Time
}

type PostSlice []Post

//Sort methods for PostSlice
func (p PostSlice) Len() int           { return len(p) }
func (p PostSlice) Less(i, j int) bool { return p[i].CreatTime.Unix() > p[j].CreatTime.Unix() }
func (p PostSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type Content string

func (c *Content) Write(p []byte) (n int, err error) {
	*c = *c + Content(p)
	return len(p), nil
}

//Time stamp related =======#
func firstNil(b []byte) int {
	var i int
	for i = len(b) - 1; ; i-- {
		if b[i] != 0 {
			break
		}
	}
	return i + 1
}

func stamp(f *os.File) {
	f.Seek(0, 0) // Making sure we start correctly. 
	f_r := bufio.NewReader(f)
	f_w := bufio.NewWriter(f)

	buf := make([]byte, 1024)
	for {
		n, err := f_r.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
	}

	// Stamp the buffer
	t := time.Now().Format(time.UnixDate)
	t += "\n"
	buf = append([]byte(t), buf...)

	// Write the buffer
	f.Seek(0, 0)       // Making sure we start correctly.
	n := firstNil(buf) //Making sure not to write nil bytes. 
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
	f.Seek(0, 0) //Reset the seek after all the reads.
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
	f.Seek(0, 0) // Make sure we're at the start
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

func readPost(f *os.File) string {
	f.Seek(0, 0) //Reset the seek after all the reads.
	f_r := bufio.NewReader(f)
	_, err := f_r.ReadSlice('\n') // Do not want the timestamp in the post
	if err != nil && err != io.EOF {
		panic(err)
	}

	//Some buffer optimisation, see golang ioutil ReadFile source
	var b_size int64
	if fi, err := f.Stat(); err == nil {
		if f_size := fi.Size(); f_size < 1e9 {
			b_size = f_size
		}
	}
	if b_size == 0 {
		b_size = 10
	}

	postBody := make([]byte, b_size)
	for {
		n, err := f_r.Read(postBody)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
	}
	m := firstNil(postBody) // just incase
	return string(postBody[:m])
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	var c Content

	t_posts, err := template.ParseFiles(files + "posts.html")
	if err != nil {
		panic(err)
	}
	t_index, err := template.ParseFiles("main_web/main.html")
	if err != nil {
		panic(err)
	}

	dir, err := ioutil.ReadDir(files + "posts/")
	if err != nil {
		panic(err)
	}
	var p PostSlice
	p = make([]Post, len(dir))

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

	//Reading in posts and get the timestamp/timestamp them
	for i := 0; i < len(p); i++ {
		f, err := os.OpenFile(files+"posts/"+p[i].Title+".txt", os.O_RDWR, 0644)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if !isStamped(f) {
			stamp(f)
		}
		defer f.Close()
		p[i].Body = readPost(f)
		p[i].CreatTime = postCreateTime(f)
	}
	//sort Posts after create time
	sort.Sort(p)
	//Parse the posts
	for i := 0; i < len(p); i++ {
		t_posts.Execute(&c, p[i])
	}

	t_index.Execute(w, template.HTML(c))
}
