package commonStructs

import (
	"time"

	"github.com/google/uuid"
)

type TicketStatus string

const (
	Open    TicketStatus = "OPEN"
	Ongoing TicketStatus = "ONGOING"
	Booked  TicketStatus = "BOOKED"
)

type UpdateTicketStatusRequest struct {
	InvoiceId string        `json:"invoiceId" form:"invoiceId" validate:"required"`
	Status    PaymentStatus `json:"status" form:"status" validate:"required,is_payment_status"`
	Message   string        `json:"message" form:"message" validate:"required"`
}

type UpdateTicketStatusServicePayload struct {
	InvoiceId string
	Status    PaymentStatus
	UserId    string
	Message   string
	JWTToken  string
}

type CreateTicketServicePayload struct {
	Price   int       `json:"price" form:"price" validate:"required,is_price"`
	EventId uuid.UUID `json:"eventId" form:"eventId" validate:"required"`
	SeatId  string    `json:"seatId" form:"seatId" validate:"required,is_seat_number"`
}

type UpdateTicketServicePayload struct {
	Price   int          `json:"price" form:"price" validate:"is_price"`
	EventId uuid.UUID    `json:"eventId" form:"eventId"`
	SeatId  string       `json:"seatId" form:"seatId" validate:"is_seat_number"`
	Status  TicketStatus `json:"status" form:"status"`
}

type UpdateStatusServicePayload struct {
	InvoiceId string        `json:"invoiceId" form:"invoiceId"`
	Status    PaymentStatus `json:"status" form:"status"`
	UserId    string        `json:"userId" form:"userId"`
}

type GetManyTicketsByIdsPayload struct {
	Ids string `query:"ids"`
}

type TicketWithEvent struct {
	Price     int       `json:"price"`
	EventName string    `json:"eventName"`
	EventTime time.Time `json:"eventTime"`
	Location  string    `json:"location"`
	SeatId    string    `json:"seatId"`
}
