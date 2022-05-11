package router

import (
	"github.com/DeniesKresna/jhapi2/config"
	"github.com/DeniesKresna/jhapi2/repository/cache/redis"
	userSql "github.com/DeniesKresna/jhapi2/repository/sql/user"
	userService "github.com/DeniesKresna/jhapi2/service/user"
	"github.com/DeniesKresna/jhapi2/utils"
	"github.com/go-playground/validator/v10"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type Object struct {
}

func Provide() {
	utility := utils.Provide()

	// repositories
	cacheRepo := redis.NewRepository(config.GetRedisClient().Cache)
	userRepo := userSql.Provide(config.GetGormClient().DB, &cacheRepo, utility)
	userSvc := userService.Provide(userRepo, utility)

	iCotrl := controller.Provide(validator.New(), userSvc, utility)
}

func GetRouter() (router http.Handler, err error) {

	cfg := config.NewConfig()

	var app *newrelic.Application
	app, err = newrelic.NewApplication(
		newrelic.ConfigAppName(*config.Get().Application.Name),
		newrelic.ConfigLicense(*config.Get().NewRelic.License),
		newrelic.ConfigDistributedTracerEnabled(true),
		func(c *newrelic.Config) {
			enable, _ := strconv.ParseBool(*config.Get().NewRelic.Enable)
			c.Enabled = enable
		},
	)
	if nil != err {
		log.Panic().AnErr("error", err).Msg("Fail to create instance agent new relic")
	}
	subRouter := mux.NewRouter().StrictSlash(true)
	subRouter.Use(nrgorilla.Middleware(app))
	ouputPanic := capture.Slack
	if *config.Get().Application.Environment == "dev" {
		ouputPanic = capture.Stdout
	}
	capture.InitCapture(&capture.Config{
		OutputPanic: ouputPanic,
		ShowTrace:   true,
		SlackClient: capture.SlackClient{
			WebHookUrl:  *config.Get().Slack.WebHook,
			UserName:    "Alert Tataskola API",
			Channel:     *config.Get().Slack.Channel,
			MentionUser: strings.Split(*config.Get().Slack.MentionUser, ","),
			MentionHere: true,
			Environment: *config.Get().Application.Environment,
			TimeOut:     5 * time.Second,
		},
	})

	subRouter.Use(capture.CapturePanic)

	log.Info().Msg("register routes")

	for _, route := range ProvideRoute() {

		//Add api prefix
		//route.Pattern = "/api" + route.Pattern

		//Check if route should go through jwt middleware
		if !route.IsExcluded {
			route.HandlerFunc = MidlwrSvc.ValidateRequest(route.Pattern, route.WhoCanAccess, route.HandlerFunc)
		}

		//Default middleware : should apply to every route
		route.HandlerFunc = cfg.InjectConfig(route.HandlerFunc)
		subRouter.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	//this changes for cors for browser
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://127.0.0.1:8787", "http://localhost:8787", "https://10.184.0.17:8787", "https://10.184.0.18:8787", "https://10.184.0.19:8787", "https://10.11.11.15:8787", "http://10.11.11.15:8787", "http://10.22.0.101:8787", "http://tataskola-dashboard-prod:8787", "http://tataskola-apis:8080", "http://34.93.230.42:30180", "http://13.126.115.212:30790/", "http://34.101.145.198:30888/", "https://web-payments.tataskola.id/", "https://web-payments-stage.tataskola.id", "https://web-payments-pre.tataskola.id", "https://www.tataskola.id", "https://stage-dashboard.tataskola.id", "https://pre-dashboard.tataskola.id", "https://dashboard.tataskola.id", "https://prod-dashboard.tataskola.id", "https://pembayaran.tataskola.id", "https://pembayaran-stage.tataskola.id", "https://pembayaran-pre.tataskola.id", "https://pembayaran-prod.tataskola.id", "http://tataskola-apis-prod:8080", "*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, "Post", http.MethodOptions},
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})

	router = c.Handler(subRouter)

	return router, nil
}
