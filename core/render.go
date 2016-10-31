package core

import "github.com/russross/blackfriday"
import "io/ioutil"
import "log"

func (c *Context) markDown(path string) ([]byte, error) {
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return blackfriday.MarkdownBasic(fileBytes), nil
}

func (c *Context) MD(code int, path string) {
	mark, err := c.markDown(path)
	if err != nil {
		log.Println("Fail to Parsing Markdown File", path, err)
		return
	}
	resp := c.Response()
	resp.WriteHeader(code)
	header := resp.Header()
	header.Set("Content-Type", "text/html; charset=utf-8")
	resp.ResponseWriter.Write(mark)
}

func (c *Context) HTML(code int, name string, data interface{}) {
	resp := c.Response()
	resp.WriteHeader(code)
	header := resp.Header()
	header.Set("Content-Type", "text/html; charset=utf-8")
	c.AppContext.Templates.ExecuteTemplate(resp, name, data)
}

func (c *Context) NotFound(name string) {
	resp := c.Response()
	resp.WriteHeader(404)
	header := resp.Header()
	header.Set("Content-Type", "text/html; charset=utf-8")
	c.AppContext.Templates.ExecuteTemplate(resp, name, nil)
}
