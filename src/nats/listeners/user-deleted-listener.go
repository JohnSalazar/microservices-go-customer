package listeners

import (
	"context"
	"customer/src/models"
	natsMetrics "customer/src/nats/interfaces"
	"customer/src/services"
	"encoding/json"
	"fmt"
	"log"

	common_nats "github.com/JohnSalazar/microservices-go-common/nats"
	common_service "github.com/JohnSalazar/microservices-go-common/services"
	"github.com/nats-io/nats.go"

	trace "github.com/JohnSalazar/microservices-go-common/trace/otel"
)

type UserDeletedListener struct {
	service     *services.CustomerService
	email       common_service.EmailService
	publisher   common_nats.Publisher
	natsMetrics natsMetrics.NatsMetric
}

func NewUserDeletedListener(
	service *services.CustomerService,
	email common_service.EmailService,
	publisher common_nats.Publisher,
	natsMetrics natsMetrics.NatsMetric,
) *UserDeletedListener {
	return &UserDeletedListener{
		service:     service,
		email:       email,
		publisher:   publisher,
		natsMetrics: natsMetrics,
	}
}

func (l *UserDeletedListener) ProcessUserDeleted() nats.MsgHandler {
	return func(msg *nats.Msg) {
		ctx := context.Background()
		_, span := trace.NewSpan(ctx, fmt.Sprintf("publish.%s\n", msg.Subject))
		defer span.End()

		customer := &models.Customer{}
		err := json.Unmarshal(msg.Data, customer)

		err = l.service.Delete(context.Background(), customer.ID)
		if err != nil {
			trace.FailSpan(span, fmt.Sprintf("error customer delete: %v", err))
			go l.email.SendSupportMessage(err.Error())
			return
		}

		data, err := json.Marshal(customer)
		if err != nil {
			trace.FailSpan(span, "error json parse")
			log.Printf("error json parse: %v", err)
			return
		}

		err = l.publisher.Publish(string(common_nats.CustomerDeleted), data)
		if err != nil {
			l.natsMetrics.ErrorPublish()
			trace.FailSpan(span, fmt.Sprintf("error publisher: %v", err))
			log.Printf("error publisher: %v", err)
			return
		}

		fmt.Println(fmt.Sprintf("%s processed!!!\n", msg.Subject))

		err = msg.Ack()
		if err != nil {
			log.Printf("stan msg.Ack error: %v", err)
		}
	}
}
