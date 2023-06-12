package nats

import (
	"customer/src/nats/interfaces"
	"customer/src/nats/listeners"
	"customer/src/services"
	"log"

	common_nats "github.com/JohnSalazar/microservices-go-common/nats"
	common_service "github.com/JohnSalazar/microservices-go-common/services"
	"github.com/nats-io/nats.go"
)

type listen struct {
	js nats.JetStreamContext
}

const queueGroupName string = "customers-service"

var subscribe common_nats.Listener
var userDelete *listeners.UserDeletedListener

func NewListen(
	js nats.JetStreamContext,
	service *services.CustomerService,
	email common_service.EmailService,
	publisher common_nats.Publisher,
	natsMetrics interfaces.NatsMetric,
) *listen {
	subscribe = common_nats.NewListener(js)
	userDelete = listeners.NewUserDeletedListener(service, email, publisher, natsMetrics)
	return &listen{
		js: js,
	}
}

func (l *listen) Listen() {
	go subscribe.Listener(string(common_nats.UserDeleted), queueGroupName, queueGroupName+"_0", userDelete.ProcessUserDeleted())

	log.Printf("Listener on!!!")
}
