package routers

import (
	"customer/src/controllers"
	"fmt"

	"github.com/JohnSalazar/microservices-go-common/config"
	"github.com/JohnSalazar/microservices-go-common/middlewares"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	common_service "github.com/JohnSalazar/microservices-go-common/services"
)

type Router struct {
	config             *config.Config
	serviceMetrics     common_service.Metrics
	authentication     *middlewares.Authentication
	customerController *controllers.CustomerController
}

func NewRouter(
	config *config.Config,
	serviceMetrics common_service.Metrics,
	authentication *middlewares.Authentication,
	customerController *controllers.CustomerController,
) *Router {
	return &Router{
		config:             config,
		serviceMetrics:     serviceMetrics,
		authentication:     authentication,
		customerController: customerController,
	}
}

func (r *Router) RouterSetup() *gin.Engine {
	router := r.initRouter()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middlewares.CORS())
	router.Use(location.Default())
	router.Use(otelgin.Middleware(r.config.Jaeger.ServiceName))
	router.Use(middlewares.Metrics(r.serviceMetrics))

	router.GET("/healthy", middlewares.Healthy())
	router.GET("/metrics", middlewares.MetricsHandler())

	v1 := router.Group(fmt.Sprintf("/api/%s", r.config.ApiVersion))
	v1.Use(r.authentication.Verify())

	v1.GET("/profile", r.customerController.Profile)
	v1.POST("/", r.customerController.AddCustomer)
	v1.PUT("/", r.customerController.UpdateCustomer)
	v1.GET("/addresses", r.customerController.GetAddressesCustomer)
	v1.GET("/address/:id", r.customerController.GetAddress)
	v1.POST("/address", r.customerController.AddAddress)
	v1.PUT("/address/:id", r.customerController.UpdateAddress)
	v1.DELETE("/address/:id", r.customerController.DeleteAddress)

	return router
}

func (r *Router) initRouter() *gin.Engine {
	if r.config.Production {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	return gin.New()
}
