package routes

import (
	//"context"
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/syedomair/ex-paygate-capture/routes/capture"
	"github.com/syedomair/ex-paygate-lib/lib/container"
	log "github.com/syedomair/ex-paygate-lib/lib/tools/logger"
	"github.com/syedomair/ex-paygate-lib/lib/tools/request"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

// NewRouter comment
func NewRouter(c container.Container) *chi.Mux {
	router := chi.NewRouter()
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST, GET"},
		AllowedHeaders:   []string{"ApiKey", "RefreshToken", "Token", "FrontendURL", "Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Access-Control-Allow-Origin"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	router.Use(cors.Handler)

	router.Route("/v1", func(r chi.Router) {
		r.Mount("/", routerSetup(
			c.Db(),
			c.Logger(),
			c.SigningKey(),
		))
	})

	return router
}

// Route Public
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
	Access      string
}

// Routes List Public
type Routes []Route

func routerSetup(db *gorm.DB, logger log.Logger, signingKey string) *chi.Mux {

	payCapture := capture.NewPaymentService(logger)
	repoCapture := capture.NewPostgresRepository(db, logger, payCapture)

	router := chi.NewRouter()
	captureController := &capture.Controller{
		Logger: logger,
		Repo:   repoCapture,
		Pay:    payCapture,
	}

	var routes = Routes{
		Route{
			"Ping",
			"GET",
			"/ping",
			captureController.Ping,
			"public",
		},
		Route{
			"Capture",
			"POST",
			"/capture",
			captureController.CaptureAction,
			"public",
		},
	}

	for _, route := range routes {
		var handler http.Handler
		strKey := uuid.NewString()
		handler = route.HandlerFunc
		handler = loggerMiddleware(strKey, logger, handler, route.Method, route.Pattern)
		router.Method(route.Method, route.Pattern, handler)
	}
	return router
}

func loggerMiddleware(strKey string, logger log.Logger, handle http.Handler, method string, pattern string) http.Handler {
	methodName := "loggerMiddleware"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testNum := r.Header.Get("Testnum")
		ctx := r.Context()
		ctx = context.WithValue(ctx, request.ContextKeyRequestID, strKey)
		r = r.WithContext(ctx)
		logger.Debug(strKey, "******* %v ****%v:%v**********", testNum, method, pattern)
		logger.Debug(strKey, "M:%v ", methodName)
		handle.ServeHTTP(w, r)
	})
}
