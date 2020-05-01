package apps

import (
	"codeunic.com/DocumentationApp/models"
	"codeunic.com/DocumentationApp/routes"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo-contrib/prometheus"
	"golang.org/x/crypto/acme/autocert"

	"html/template"
	"io/ioutil"
	"os"
)

type Main struct {
	init         *echo.Echo
	dirTemplates string
	port         string
	signingKey   string
}

func (then *Main) Initialize() {
	then.dirTemplates = "templates/**/*.html"
	then.port = ":1323"
	then.signingKey = "1234asjosadcpe02@asjas#"

	then.init = echo.New()
	then.static()
	then.templates()
	then.routes()
	then.saveRoutes()
}

func (then *Main) templates() {
	renderer := &models.TemplateRenderer{
		Templates: template.Must(template.ParseGlob(then.dirTemplates)),
	}
	then.init.Renderer = renderer
}

func (then *Main) static() {
	then.init.Static("/static", "assets")
}

func (then *Main) routes() {
	then.init.GET("/", routes.HomeRoute).Name = "main"
}

func (then *Main) saveRoutes() {
	data, err := json.MarshalIndent(then.init.Routes(), "", "  ")
	if err != nil {
		return
	}
	ioutil.WriteFile("routes.json", data, 0644)
}

func (then *Main) MidGzip() {
	then.init.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
}

func (then *Main) MidCors() {
	then.init.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://localhost", "http://localhost"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
}
func (then *Main) MidJWT() {
	then.init.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  []byte(then.signingKey),
		ContextKey:  "user",
		TokenLookup: "header:" + echo.HeaderAuthorization,
		AuthScheme:  "Bearer",
		Claims:      jwt.MapClaims{},
	}))
}
func (then *Main) MidLogger(basic bool) {
	then.init.Use(middleware.Logger())
	var log middleware.LoggerConfig = middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}
	if !basic {
		log = middleware.LoggerConfig{
			Format: `{"time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}","host":"${host}",` +
				`"method":"${method}","uri":"${uri}","status":${status},"error":"${error}","latency":${latency},` +
				`"latency_human":"${latency_human}","bytes_in":${bytes_in},` +
				`"bytes_out":${bytes_out}}` + "\n",
			Output: os.Stdout,
		}
	}
	then.init.Use(middleware.LoggerWithConfig(log))
}
func (then *Main) Prometheus() {
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(then.init)
}
func (then *Main) Recover() {
	then.init.Use(middleware.Recover())
	then.init.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         4 << 10, // 4 KB
		DisableStackAll:   false,
		DisablePrintStack: false,
	}))
}
func (then *Main) Start() {
	then.init.Logger.Fatal(then.init.Start(then.port))
}
func (then *Main) MidTrailingSlash(b bool) {
	then.init.Pre(middleware.AddTrailingSlash())
}

func (then *Main) Cache(b bool) {
	then.init.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")
}
