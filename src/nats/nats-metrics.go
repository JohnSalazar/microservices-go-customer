package nats

import (
	"github.com/JohnSalazar/microservices-go-common/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type natsMetric struct {
	config *config.Config
}

var successCustomerCreated prometheus.Counter
var successCustomerUpdated prometheus.Counter
var successCustomerDeleted prometheus.Counter
var successAddressCreated prometheus.Counter
var successAddressUpdated prometheus.Counter
var successAddressDeleted prometheus.Counter
var errorPublish prometheus.Counter

func NewNatsMetric(
	config *config.Config,
) *natsMetric {
	successCustomerCreated = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: config.AppName + "_nats_success_customer_created_total",
			Help: "The total number of success customer created NATS messages",
		},
	)

	successCustomerUpdated = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: config.AppName + "_nats_success_customer_updated_total",
			Help: "The total number of success customer updated NATS messages",
		},
	)

	successCustomerDeleted = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: config.AppName + "_nats_success_customer_deleted_total",
			Help: "The total number of success customer deleted NATS messages",
		},
	)

	successAddressCreated = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: config.AppName + "_nats_success_address_created_total",
			Help: "The total number of success address created NATS messages",
		},
	)

	successAddressUpdated = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: config.AppName + "_nats_success_address_updated_total",
			Help: "The total number of success address updated NATS messages",
		},
	)

	successAddressDeleted = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: config.AppName + "_nats_success_address_deleted_total",
			Help: "The total number of success address deleted NATS messages",
		},
	)

	errorPublish = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: config.AppName + "_nats_error_publish_message_total",
			Help: "The total number of error NATS publish message",
		},
	)

	return &natsMetric{
		config: config,
	}
}

func (nats *natsMetric) SuccessPublishCustomerCreated() {
	successCustomerCreated.Inc()
}

func (nats *natsMetric) SuccessPublishCustomerUpdated() {
	successCustomerUpdated.Inc()
}

func (nats *natsMetric) SuccessPublishCustomerDeleted() {
	successCustomerDeleted.Inc()
}

func (nats *natsMetric) SuccessPublishAddressCreated() {
	successAddressCreated.Inc()
}

func (nats *natsMetric) SuccessPublishAddressUpdated() {
	successAddressUpdated.Inc()
}

func (nats *natsMetric) SuccessPublishAddressDeleted() {
	successAddressDeleted.Inc()
}

func (nats *natsMetric) ErrorPublish() {
	errorPublish.Inc()
}
