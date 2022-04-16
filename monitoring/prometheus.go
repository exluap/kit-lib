package monitoring

import (
	"fmt"
	"github.com/exluap/kit/log"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"regexp"
)

type prometheusMetricsSrv struct {
	logger     log.CLoggerFunc
	registerer *prometheus.Registry
	router     *mux.Router
	httpSrv    *http.Server
}

func NewMetricsServer(logger log.CLoggerFunc) MetricsServer {
	srv := &prometheusMetricsSrv{
		logger: logger,
	}
	return srv
}

func (s *prometheusMetricsSrv) Init(config *Config, metricProviders ...MetricsProvider) error {

	if match, _ := regexp.MatchString("^\\d{1,6}$", config.Port); !match {
		return ErrPrometheusInvalidPort(config.Port)
	}

	url := config.UrlPath
	if url == "" {
		url = "/metrics"
	}

	s.registerer = prometheus.NewRegistry()

	for _, pr := range metricProviders {
		for _, m := range pr.GetCollector()() {
			if err := s.registerer.Register(m); err != nil {
				return ErrPrometheusRegisterAppMetrics(err)
			}
		}
	}

	if config.GoMetrics {
		if err := s.registerer.Register(prometheus.NewGoCollector()); err != nil {
			return ErrPrometheusRegisterGoMetrics(err)
		}
		if err := s.registerer.Register(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{})); err != nil {
			return ErrPrometheusRegisterProcessMetrics(err)
		}
	}

	s.router = mux.NewRouter()
	s.router.Path(url).Handler(promhttp.HandlerFor(s.registerer, promhttp.HandlerOpts{}))

	s.httpSrv = &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Port),
		Handler: s.router,
	}

	return nil

}

func (s *prometheusMetricsSrv) Listen() {
	go func() {
		l := s.logger().Pr("http").Cmp("prometheus").
			Mth("listen").
			F(log.FF{"url": s.httpSrv.Addr}).
			Inf("start listening")

		if err := s.httpSrv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				l.E(ErrPrometheusHttpServer(err)).St().Err()
			} else {
				l.Dbg("server closed")
			}
		}
	}()
}

func (s *prometheusMetricsSrv) Close() {
	_ = s.httpSrv.Close()
	s.logger().Pr("http").Cmp("prometheus").Inf("closed")
}
