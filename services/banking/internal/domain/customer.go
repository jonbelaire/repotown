package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Common errors
var (
	ErrCustomerNotFound = errors.New("customer not found")
	ErrCustomerExists   = errors.New("customer with email already exists")
)

// CustomerStatus represents the status of a customer
type CustomerStatus string

// Customer statuses
const (
	CustomerStatusActive   CustomerStatus = "active"
	CustomerStatusInactive CustomerStatus = "inactive"
	CustomerStatusBlocked  CustomerStatus = "blocked"
)

// Customer represents a bank customer
type Customer struct {
	ID          string         `json:"id"`
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	Email       string         `json:"email"`
	PhoneNumber string         `json:"phone_number"`
	Address     Address        `json:"address"`
	Status      CustomerStatus `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// Address represents a physical address
type Address struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	PostalCode string `json:"postal_code"`
}

// NewCustomer creates a new customer
func NewCustomer(firstName, lastName, email, phoneNumber string, address Address) *Customer {
	now := time.Now()
	return &Customer{
		ID:          uuid.New().String(),
		FirstName:   firstName,
		LastName:    lastName,
		Email:       email,
		PhoneNumber: phoneNumber,
		Address:     address,
		Status:      CustomerStatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// FullName returns the customer's full name
func (c *Customer) FullName() string {
	return c.FirstName + " " + c.LastName
}

// IsActive checks if the customer is active
func (c *Customer) IsActive() bool {
	return c.Status == CustomerStatusActive
}

// UpdateStatus updates the customer status
func (c *Customer) UpdateStatus(status CustomerStatus) {
	c.Status = status
	c.UpdatedAt = time.Now()
}

// UpdateAddress updates the customer's address
func (c *Customer) UpdateAddress(address Address) {
	c.Address = address
	c.UpdatedAt = time.Now()
}

// UpdateContact updates the customer's contact information
func (c *Customer) UpdateContact(email, phoneNumber string) {
	c.Email = email
	c.PhoneNumber = phoneNumber
	c.UpdatedAt = time.Now()
}