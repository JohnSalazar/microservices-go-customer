package controllers

import (
	"customer/src/dtos"
	"customer/src/models"
	natsMetrics "customer/src/nats/interfaces"
	"customer/src/services"
	"customer/src/validators"
	"net/http"

	"github.com/JohnSalazar/microservices-go-common/config"
	"github.com/JohnSalazar/microservices-go-common/helpers"
	httputil "github.com/JohnSalazar/microservices-go-common/httputil"
	common_nats "github.com/JohnSalazar/microservices-go-common/nats"
	common_security "github.com/JohnSalazar/microservices-go-common/security"
	trace "github.com/JohnSalazar/microservices-go-common/trace/otel"
	"github.com/gin-gonic/gin"
)

type CustomerController struct {
	config              *config.Config
	publisher           common_nats.Publisher
	natsMetrics         natsMetrics.NatsMetric
	managerTokensCommon *common_security.ManagerTokens
	serviceCustomer     *services.CustomerService
	serviceAddress      *services.AddressService
}

func NewCustomerController(
	config *config.Config,
	publisher common_nats.Publisher,
	natsMetrics natsMetrics.NatsMetric,
	managerTokensCommon *common_security.ManagerTokens,
	serviceCustomer *services.CustomerService,
	serviceAddress *services.AddressService,
) *CustomerController {
	return &CustomerController{
		config:              config,
		publisher:           publisher,
		natsMetrics:         natsMetrics,
		managerTokensCommon: managerTokensCommon,
		serviceCustomer:     serviceCustomer,
		serviceAddress:      serviceAddress,
	}
}

func (customer *CustomerController) Profile(c *gin.Context) {
	ID, customerIDOk := c.Get("user")
	if !customerIDOk {
		httputil.NewResponseError(c, http.StatusForbidden, "invalid user")
		return
	}

	isID := helpers.IsValidID(ID.(string))
	if !isID {
		httputil.NewResponseError(c, http.StatusBadRequest, "invalid customerId")
		return
	}

	_customer, err := customer.serviceCustomer.FindByID(c.Request.Context(), helpers.StringToID(ID.(string)))
	if err != nil {
		httputil.NewResponseError(c, http.StatusBadRequest, "customer not found")
		return
	}

	customerDTO := mapToCustomerDTO(_customer)

	c.JSON(http.StatusOK, customerDTO)
}

func (customer *CustomerController) GetAddressesCustomer(c *gin.Context) {
	_, span := trace.NewSpan(c.Request.Context(), "CustomerController.GetAddressesCustomer")
	defer span.End()

	ID, customerIDOk := c.Get("user")
	if !customerIDOk {
		httputil.NewResponseError(c, http.StatusForbidden, "invalid customer")
		return
	}

	isID := helpers.IsValidID(ID.(string))
	if !isID {
		httputil.NewResponseError(c, http.StatusBadRequest, "invalid customerId")
		return
	}

	customerID := helpers.StringToID(ID.(string))

	addresses, err := customer.serviceAddress.GetByCustomerID(c.Request.Context(), customerID)
	if err != nil {
		httputil.NewResponseError(c, http.StatusBadRequest, "addresses not found")
		return
	}

	var addressesDTO []*dtos.Address

	for _, address := range addresses {
		addressDTO := mapToAddressDTO(address)
		addressesDTO = append(addressesDTO, addressDTO)
	}

	c.JSON(http.StatusOK, addressesDTO)
}

func (customer *CustomerController) GetAddress(c *gin.Context) {
	_, span := trace.NewSpan(c.Request.Context(), "CustomerController.GetAddress")
	defer span.End()

	isID := helpers.IsValidID(c.Param("id"))
	if !isID {
		httputil.NewResponseError(c, http.StatusBadRequest, "invalid id")
		return
	}

	ID := helpers.StringToID(c.Param("id"))

	address, err := customer.serviceAddress.FindByID(c.Request.Context(), ID)
	if err != nil {
		httputil.NewResponseError(c, http.StatusBadRequest, "address not found")
		return
	}

	addressDTO := mapToAddressDTO(address)

	c.JSON(http.StatusOK, addressDTO)
}

func (customer *CustomerController) AddCustomer(c *gin.Context) {
	_, span := trace.NewSpan(c.Request.Context(), "CustomerController.AddCustomer")
	defer span.End()

	ID, customerIDOk := c.Get("user")
	if !customerIDOk {
		httputil.NewResponseError(c, http.StatusForbidden, "invalid user")
		return
	}

	isID := helpers.IsValidID(ID.(string))
	if !isID {
		httputil.NewResponseError(c, http.StatusBadRequest, "invalid customerId")
		return
	}

	addCustomerDTO := &dtos.AddCustomer{}
	err := c.BindJSON(addCustomerDTO)
	if err != nil {
		trace.FailSpan(span, "Error json parse")
		httputil.NewResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	customerID := helpers.StringToID(ID.(string))

	result := validators.ValidateAddCustomer(customerID, addCustomerDTO)
	if result != nil {
		httputil.NewResponseError(c, http.StatusBadRequest, result)
		return
	}

	customerExists, _ := customer.serviceCustomer.FindByID(c.Request.Context(), customerID)
	if customerExists != nil {
		httputil.NewResponseError(c, http.StatusBadRequest, "customer already exists")
		return
	}

	customerMapped := &models.Customer{
		ID:        customerID,
		Email:     addCustomerDTO.Email,
		FirstName: addCustomerDTO.FirstName,
		LastName:  addCustomerDTO.LastName,
		Avatar:    addCustomerDTO.Avatar,
		Phone:     addCustomerDTO.Phone,
	}

	_customer, err := customer.serviceCustomer.Create(c.Request.Context(), customerMapped)
	if err != nil {
		trace.AddSpanError(span, err)
		httputil.NewResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	customer.natsMetrics.SuccessPublishCustomerCreated()

	c.JSON(http.StatusCreated, _customer)
}

func (customer *CustomerController) UpdateCustomer(c *gin.Context) {
	_, span := trace.NewSpan(c.Request.Context(), "CustomerController.UpdateCustomer")
	defer span.End()

	ID, customerIDOk := c.Get("user")
	if !customerIDOk {
		httputil.NewResponseError(c, http.StatusForbidden, "invalid user")
		return
	}

	isID := helpers.IsValidID(ID.(string))
	if !isID {
		httputil.NewResponseError(c, http.StatusBadRequest, "invalid customerId")
		return
	}

	customerID := helpers.StringToID(ID.(string))

	updateCustomerDTO := &dtos.UpdateCustomer{}
	err := c.BindJSON(updateCustomerDTO)
	if err != nil {
		httputil.NewResponseError(c, http.StatusForbidden, err.Error())
		return
	}

	result := validators.ValidateUpdateCustomer(customerID, updateCustomerDTO)
	if result != nil {
		httputil.NewResponseError(c, http.StatusBadRequest, result)
		return
	}

	_, err = customer.serviceCustomer.FindByID(c.Request.Context(), customerID)
	if err != nil {
		httputil.NewResponseError(c, http.StatusBadRequest, "customer not found")
		return
	}

	customerMapped := &models.Customer{
		ID:        customerID,
		Email:     updateCustomerDTO.Email,
		FirstName: updateCustomerDTO.FirstName,
		LastName:  updateCustomerDTO.LastName,
		Avatar:    updateCustomerDTO.Avatar,
		Phone:     updateCustomerDTO.Phone,
		Version:   updateCustomerDTO.Version,
	}

	_customer, err := customer.serviceCustomer.Update(c.Request.Context(), customerMapped)
	if err != nil {
		httputil.NewResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	customer.natsMetrics.SuccessPublishCustomerUpdated()

	c.JSON(http.StatusOK, _customer)
}

func (customer *CustomerController) AddAddress(c *gin.Context) {
	_, span := trace.NewSpan(c.Request.Context(), "CustomerController.AddAddress")
	defer span.End()

	ID, customerIDOk := c.Get("user")
	if !customerIDOk {
		httputil.NewResponseError(c, http.StatusForbidden, "invalid user")
		return
	}

	isID := helpers.IsValidID(ID.(string))
	if !isID {
		httputil.NewResponseError(c, http.StatusBadRequest, "invalid customerId")
		return
	}

	addAddressDTO := &dtos.AddAddress{}
	err := c.BindJSON(addAddressDTO)
	if err != nil {
		trace.FailSpan(span, "Error json parse")
		httputil.NewResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	customerID := helpers.StringToID(ID.(string))

	result := validators.ValidateAddAddress(addAddressDTO)
	if result != nil {
		httputil.NewResponseError(c, http.StatusBadRequest, result)
		return
	}

	addressMapped := &models.Address{
		CustomerID: customerID,
		Street:     addAddressDTO.Street,
		City:       addAddressDTO.City,
		Province:   addAddressDTO.Province,
		Code:       addAddressDTO.Code,
		Type:       addAddressDTO.Type,
	}

	addressExists := customer.serviceAddress.FindAddressExists(c.Request.Context(), addressMapped)
	if addressExists {
		httputil.NewResponseError(c, http.StatusBadRequest, "address already exists")
		return
	}

	address, err := customer.serviceAddress.Create(c.Request.Context(), addressMapped)
	if err != nil {
		trace.AddSpanError(span, err)
		httputil.NewResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	customer.natsMetrics.SuccessPublishAddressCreated()

	c.JSON(http.StatusCreated, address)
}

func (customer *CustomerController) UpdateAddress(c *gin.Context) {
	_, span := trace.NewSpan(c.Request.Context(), "CustomerController.UpdateAddress")
	defer span.End()

	isID := helpers.IsValidID(c.Param("id"))
	if !isID {
		httputil.NewResponseError(c, http.StatusBadRequest, "invalid Id")
		return
	}

	ID := helpers.StringToID(c.Param("id"))

	updateAddressDTO := &dtos.UpdateAddress{}
	err := c.BindJSON(updateAddressDTO)
	if err != nil {
		httputil.NewResponseError(c, http.StatusForbidden, err.Error())
		return
	}

	result := validators.ValidateUpdateAddress(updateAddressDTO)
	if result != nil {
		httputil.NewResponseError(c, http.StatusBadRequest, result)
		return
	}

	addressExists, err := customer.serviceAddress.FindByID(c.Request.Context(), ID)
	if err != nil {
		httputil.NewResponseError(c, http.StatusBadRequest, "customer not found")
		return
	}

	addressMapped := &models.Address{
		ID:       addressExists.ID,
		Street:   updateAddressDTO.Street,
		City:     updateAddressDTO.City,
		Province: updateAddressDTO.Province,
		Code:     updateAddressDTO.Code,
		Type:     updateAddressDTO.Type,
		Version:  updateAddressDTO.Version,
	}

	address, err := customer.serviceAddress.Update(c.Request.Context(), addressMapped)
	if err != nil {
		httputil.NewResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	customer.natsMetrics.SuccessPublishAddressUpdated()

	addressDTO := mapToAddressDTO(address)

	c.JSON(http.StatusOK, addressDTO)
}

func (customer *CustomerController) DeleteAddress(c *gin.Context) {
	_, span := trace.NewSpan(c.Request.Context(), "CustomerController.DeleteAddress")
	defer span.End()

	isID := helpers.IsValidID(c.Param("id"))
	if !isID {
		httputil.NewResponseError(c, http.StatusBadRequest, "invalid Id")
		return
	}

	ID := helpers.StringToID(c.Param("id"))

	err := customer.serviceAddress.Delete(c.Request.Context(), ID)
	if err != nil {
		httputil.NewResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	customer.natsMetrics.SuccessPublishAddressDeleted()

	httputil.NewResponseSuccess(c, http.StatusOK, "deleted address")
}

func mapToCustomerDTO(user *models.Customer) *dtos.Customer {
	return &dtos.Customer{
		ID:        user.ID.Hex(),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Version:   user.Version,
	}
}

func mapToAddressDTO(address *models.Address) *dtos.Address {
	return &dtos.Address{
		ID:        address.ID.Hex(),
		Street:    address.Street,
		City:      address.City,
		Province:  address.Province,
		Code:      address.Code,
		Type:      address.Type,
		CreatedAt: address.CreatedAt,
		UpdatedAt: address.UpdatedAt,
		Version:   address.Version,
	}
}
