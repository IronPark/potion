package core

func (c *Potion) HTML(code int, name string, data interface{}) {
	resp := c.Response()
	resp.WriteHeader(code)
	header := resp.Header()
	header.Set("Content-Type", "text/html; charset=utf-8")
	c.AppContext.Templates.ExecuteTemplate(resp, name, data)
}

func (c *Potion) NotFound(name string) {
	resp := c.Response()
	resp.WriteHeader(404)
	header := resp.Header()
	header.Set("Content-Type", "text/html; charset=utf-8")
	c.AppContext.Templates.ExecuteTemplate(resp, name, nil)
}
