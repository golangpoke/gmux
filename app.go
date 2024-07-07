package gmux

import (
	"fmt"
	"log/slog"
	"net/http"
)

type level int

const (
	LevelRelease level = iota
	LevelDebug
)

type HandleFunc func(c *Ctx) Api
type Middleware func(next HandleFunc) HandleFunc

type App struct {
	logLevel level
	mux      *http.ServeMux
	route    *route
}

type route struct {
	path        string
	hashMap     map[string]HandleFunc
	Middlewares []Middleware
}

func New() *App {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	return &App{
		logLevel: LevelDebug,
		mux:      http.NewServeMux(),
		route: &route{
			hashMap: make(map[string]HandleFunc),
		},
	}
}

func (a *App) SetLogLevel(level level) {
	a.logLevel = level
	if level == LevelRelease {
		slog.SetLogLoggerLevel(slog.LevelInfo)
	}
}

func (a *App) Group(path string) *App {
	newApp := &App{
		logLevel: a.logLevel,
		mux:      a.mux,
		route: &route{
			path:        a.route.path + path,
			hashMap:     a.route.hashMap,
			Middlewares: a.route.Middlewares,
		},
	}
	return newApp
}
func (a *App) Use(middlewares ...Middleware) {
	a.route.Middlewares = append(a.route.Middlewares, middlewares...)
}
func (a *App) RouteMap() map[string]HandleFunc {
	return a.route.hashMap
}
func (a *App) GET(path string, handle HandleFunc) {
	a.addRoute(http.MethodGet, path, handle)
}
func (a *App) POST(path string, handle HandleFunc) {
	a.addRoute(http.MethodPost, path, handle)
}
func (a *App) PUT(path string, handle HandleFunc) {
	a.addRoute(http.MethodPut, path, handle)
}
func (a *App) DELETE(path string, handle HandleFunc) {
	a.addRoute(http.MethodDelete, path, handle)
}
func (a *App) addRoute(method string, path string, handleFunc HandleFunc) {
	path = a.route.path + path
	pattern := fmt.Sprintf("%s %s", method, path)
	a.route.hashMap[pattern] = handleFunc
	a.mux.Handle(pattern, a.handle(handleFunc))
}

func (a *App) handle(handleFunc HandleFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := newContext(w, r)
		//handle middleware
		for _, middleware := range a.route.Middlewares {
			mh := middleware(handleFunc)
			mh(c)
		}
		//handle context
		api := handleFunc(c).(Api)
		//handle api
		a.ApiHandleFunc(c, api)
	})
}

func (a *App) ApiHandleFunc(c *Ctx, api Api) {
	if r, ok := api.(*R); ok {
		m := Map{}
		m["code"] = r.Code
		m["message"] = resultMaps[r.Code]
		m["data"] = r.Data
		if a.logLevel == LevelDebug {
			if r.Error != nil {
				m["debug"] = r.Error.Error()
				slog.Debug(fmt.Sprintf("Api Has Error Error: %s", r.Error.Error()))
			}
			slog.Debug(fmt.Sprintf("Request Api Data: %v", m))

		}
		c.JSON(http.StatusOK, m)
		return
	}
	c.String(http.StatusInternalServerError, "Error Api Type")
}

func (a *App) Run(port int) error {
	return http.ListenAndServe(fmt.Sprintf(":%d", port), a.mux)
}
