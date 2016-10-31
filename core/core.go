package core

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/lars"
)

type ApplicationGlobals struct {
	Log       *log.Logger
	Templates *template.Template
}

// Reset gets called just before a new HTTP request starts calling
// middleware + handlers
func (g *ApplicationGlobals) Reset() {
	// DB = new database connection or reset....
	// We don't touch translator + log as they don't change per request
}

// Done gets called after the HTTP request has completed right before
// Context gets put back into the pool
func (g *ApplicationGlobals) Done() {
	// DB.Close()
}

func newGlobals() *ApplicationGlobals {

	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	templates := template.Must(template.ParseGlob("templates/**"))

	return &ApplicationGlobals{
		Log:       logger,
		Templates: templates,
	}
}

// Context is a custom context
type Context struct {
	*lars.Ctx  // a little dash of Duck Typing....
	AppContext *ApplicationGlobals
}

type Potion struct {
	lars.LARS
}

// RequestStart overriding
func (mc *Context) RequestStart(w http.ResponseWriter, r *http.Request) {
	// call lars context reset, must be done
	mc.Ctx.RequestStart(w, r) // MUST be called!
	mc.AppContext.Reset()
}

// RequestEnd overriding
func (mc *Context) RequestEnd() {
	mc.AppContext.Done()
	mc.Ctx.RequestEnd() // MUST be called!
}

func New() *Potion {
	l := lars.New()

	l.RegisterContext(func(l *lars.LARS) lars.Context {
		return &Context{
			Ctx:        lars.NewContext(l),
			AppContext: newGlobals(),
		}
	})

	l.RegisterCustomHandler(func(*Context) {}, func(c lars.Context, handler lars.Handler) {
		h := handler.(func(*Context))
		ctx := c.(*Context)
		h(ctx)
	})

	//MiddleWare is LI-FO
	l.Use(Logger)
	return &Potion{*l}
}

func (p *Potion) Pubilc(path string) {
	if path[len(path)-1] != []byte("/")[0] {
		path = path + "/"
	}

	files, _ := ioutil.ReadDir("./" + path)
	for _, f := range files {
		if f.IsDir() {
			fs := http.FileServer(http.Dir(path + f.Name()))
			p.Get("/"+f.Name()+"/*", http.StripPrefix("/"+f.Name(), fs))
		} else {
			p.Handle("GET", "/"+f.Name(), func(w http.ResponseWriter, r *http.Request) {
				http.ServeFile(w, r, path+f.Name())
			})
		}
	}
}

// Logger
func Logger(c lars.Context) {
	start := time.Now()
	c.Next()
	stop := time.Now()
	path := c.Request().URL.Path
	if path == "" {
		path = "/"
	}
	log.Printf("%s %d %s %s", c.Request().Method, c.Response().Status(), path, stop.Sub(start))
}
