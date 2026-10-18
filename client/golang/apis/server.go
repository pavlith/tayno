package apis

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"

	"myst/client/golang/apis/generated"
	"myst/client/golang/core"
	"myst/pkgs/logger"
	"myst/pkgs/server"
)

//go:generate oapi-codegen --config oapi-codegen-models.yaml openapi.json
//go:generate oapi-codegen --config oapi-codegen-client.yaml openapi.json
//go:generate oapi-codegen --config oapi-codegen-models.yaml openapi.json
//go:generate npx openapi-typescript-codegen --input openapi.json --output ../../web/src/api/generated --client fetch --useOptions --useUnionTypes

var log = logger.New("router", logger.Cyan)

type Server struct {
	*gin.Engine
	app    core.Application
	server *server.Server
}

func NewServer(app core.Application, ui fs.FS) *Server {
	s := new(Server)

	s.app = app

	// Set gin mode
	gin.SetMode(gin.DebugMode)

	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard

	r := gin.New()
	s.Engine = r

	// Do not redirect folders to trailing slash
	r.RedirectTrailingSlash = true
	r.RedirectFixedPath = true

	// Custom PrintRouteFunc
	gin.DebugPrintRouteFunc = printRoutes

	r.Use(recoveryMiddleware)
	r.Use(loggerMiddleware)
	r.NoRoute(noRouteMiddleware("/", embedFolder(ui, "static")))
	// Cross-origin access is restricted to explicitly trusted web origins.
	// Vite development uses a same-origin proxy by default.
	origins := []string{"http://localhost:8082", "http://127.0.0.1:8082"}
	if configured := strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGINS")); configured != "" {
		origins = nil
		for _, origin := range strings.Split(configured, ",") {
			if origin = strings.TrimSpace(origin); origin != "" {
				origins = append(origins, origin)
			}
		}
	}
	r.Use(cors.New(cors.Options{
		AllowedOrigins: origins,
		AllowedHeaders: []string{"Accept", "Content-Type", "Authorization", "X-Myst-Account"},
		AllowedMethods: []string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		MaxAge: 600,
	}))

	api := r.Group("api")

	api.POST("/auth/login", s.LoginAccount)
	api.POST("/auth/signup", s.SignupAccount)
	api.GET("/account", s.AccountStatus)
	api.POST("/authenticate", s.Authenticate)
	api.POST("/enclave", s.CreateEnclave)

	sec := api.Group("/")
	sec.Use(s.authenticationMiddleware())

	sec.POST("/logout", s.Logout)
	sec.GET("/enclave", s.Enclave)
	sec.GET("/health", s.HealthCheck)
	sec.POST("/auth/register", s.Register)
	sec.GET("/user", s.CurrentUser)
	sec.Use(s.registrationMiddleware())

	sec.POST("/keystores", s.CreateKeystore)
	sec.GET("/keystores", s.Keystores)
	sec.GET("/keystore/:keystoreId", s.Keystore)
	sec.DELETE("/keystore/:keystoreId", s.DeleteKeystore)
	sec.POST("/keystore/:keystoreId/entries", s.CreateEntry)
	sec.PATCH("/keystore/:keystoreId/entry/:entryId", s.UpdateEntry)
	sec.DELETE("/keystore/:keystoreId/entry/:entryId", s.DeleteEntry)
	sec.GET("/invitations", s.GetInvitations)
	sec.GET("/invitation/:invitationId", s.GetInvitation)
	sec.POST("/keystore/:keystoreId/invitations", s.CreateInvitation)
	sec.PATCH("/invitation/:invitationId", s.AcceptInvitation)
	sec.DELETE("/invitation/:invitationId", s.DeclineOrCancelInvitation)
	sec.POST("/invitation/:invitationId", s.FinalizeInvitation)

	return s
}

func (s *Server) Start(addr string) error {
	log.Println("app started on", addr)

	httpServer, err := server.New(addr, s.Engine)
	if err != nil {
		return err
	}

	s.server = httpServer
	return nil
}

func (s *Server) Stop() error {
	return s.server.Stop()
}

func (s *Server) HealthCheck(c *gin.Context) {
	sid := sessionId(c)

	err := s.app.HealthCheck(sid)
	if err != nil {
		log.Error(err)
		Error(c, http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

type embedFileSystem struct {
	http.FileSystem
}

func (e embedFileSystem) Exists(_ string, path string) bool {
	_, err := e.Open(path)
	return err == nil
}

func embedFolder(fsEmbed fs.FS, targetPath string) static.ServeFileSystem {
	fsys, err := fs.Sub(fsEmbed, targetPath)
	if err != nil {
		panic(err)
	}
	return embedFileSystem{
		FileSystem: http.FS(fsys),
	}
}

func Error(c *gin.Context, statusCode int, errorCodeAndOptionalMessage ...string) {
	code := ""
	msg := ""

	if len(errorCodeAndOptionalMessage) > 0 {
		code = errorCodeAndOptionalMessage[0]
	}

	if len(errorCodeAndOptionalMessage) > 1 {
		msg = errorCodeAndOptionalMessage[1]
	}

	if code == "" {
		code = fmt.Sprintf("%d", statusCode)
		msg = http.StatusText(statusCode)
	}

	c.JSON(statusCode, generated.Error{
		Code:    code,
		Message: msg,
	})
}
