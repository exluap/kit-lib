package monitoring

import "github.com/exluap/kit/er"

const (
	ErrCodePrometheusRegisterGoMetrics      = "MON-001"
	ErrCodePrometheusRegisterProcessMetrics = "MON-002"
	ErrCodePrometheusHttpServer             = "MON-003"
	ErrCodePrometheusInvalidPort            = "MON-004"
	ErrCodePrometheusRegisterAppMetrics     = "MON-005"
)

var (
	ErrPrometheusRegisterGoMetrics = func(cause error) error {
		return er.WrapWithBuilder(cause, ErrCodePrometheusRegisterGoMetrics, "").Err()
	}
	ErrPrometheusRegisterProcessMetrics = func(cause error) error {
		return er.WrapWithBuilder(cause, ErrCodePrometheusRegisterProcessMetrics, "").Err()
	}
	ErrPrometheusHttpServer  = func(cause error) error { return er.WrapWithBuilder(cause, ErrCodePrometheusHttpServer, "").Err() }
	ErrPrometheusInvalidPort = func(port string) error {
		return er.WithBuilder(ErrCodePrometheusInvalidPort, "invalid port").F(er.FF{"port": port}).Err()
	}
	ErrPrometheusRegisterAppMetrics = func(cause error) error {
		return er.WrapWithBuilder(cause, ErrCodePrometheusRegisterAppMetrics, "").Err()
	}
)
