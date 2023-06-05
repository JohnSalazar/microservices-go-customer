package interfaces

type NatsMetric interface {
	SuccessPublishCustomerCreated()
	SuccessPublishCustomerUpdated()
	SuccessPublishCustomerDeleted()
	SuccessPublishAddressCreated()
	SuccessPublishAddressUpdated()
	SuccessPublishAddressDeleted()
	ErrorPublish()
}
