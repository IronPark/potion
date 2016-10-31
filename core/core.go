package core

import (
	"html/template"
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

// Potion is a custom context
type Potion struct {
	*lars.Ctx  // a little dash of Duck Typing....
	AppContext *ApplicationGlobals
}

// RequestStart overriding
func (mc *Potion) RequestStart(w http.ResponseWriter, r *http.Request) {
	// call lars context reset, must be done
	mc.Ctx.RequestStart(w, r) // MUST be called!
	mc.AppContext.Reset()
}

// RequestEnd overriding
func (mc *Potion) RequestEnd() {
	mc.AppContext.Done()
	mc.Ctx.RequestEnd() // MUST be called!
}

func newContext(l *lars.LARS) lars.Context {
	return &Potion{
		Ctx:        lars.NewContext(l),
		AppContext: newGlobals(),
	}
}

func castCustomContext(c lars.Context, handler lars.Handler) {
	// could do it in all one statement, but in long form for readability
	h := handler.(func(*Potion))
	ctx := c.(*Potion)
	h(ctx)
}

func New() *lars.LARS {
	l := lars.New()
	l.RegisterContext(newContext) // all gets cached in pools for you
	l.RegisterCustomHandler(func(*Potion) {}, castCustomContext)
	//MiddleWare is LI-FO
	l.Use(Logger)
	return l
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
