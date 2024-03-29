package bookingService

import (
	"encoding/json"
	"fmt"

	commonStructs "github.com/Altair1618/Tubes_IF4031_03/Ticket_Service/app/common/structs"
	"github.com/Altair1618/Tubes_IF4031_03/Ticket_Service/app/configs"
	"github.com/Altair1618/Tubes_IF4031_03/Ticket_Service/app/models"
	"github.com/Altair1618/Tubes_IF4031_03/Ticket_Service/app/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func RequestBookingService(payload commonStructs.RequestBookingServicePayload) utils.ResponseBody {
	// Check if ticket exists
	db, _ := configs.GetGormClient()

	var ticket models.Ticket
	result := db.First(&ticket, "id = ?", payload.TicketId)

	if result.Error != nil {
		fmt.Println(result.Error)

		url, err := utils.GeneratePDF(false, payload.UserId, payload.BookingId.String(), commonStructs.FailedPDFPayload{
			ErrorMessage: result.Error.Error(),
		})

		if err != nil {
			fmt.Println(err)

			return utils.ResponseBody{
				Code:    fiber.StatusInternalServerError,
				Message: "Error While Fetching Data From Database",
				Data:    fiber.Map{"status": "FAILED"},
			}
		}

		return utils.ResponseBody{
			Code:    fiber.StatusInternalServerError,
			Message: "Error While Fetching Data From Database",
			Data:    fiber.Map{"status": "FAILED", "pdf_url": url},
		}
	}

	if result.RowsAffected == 0 {
		url, err := utils.GeneratePDF(false, payload.UserId, payload.BookingId.String(), commonStructs.FailedPDFPayload{
			ErrorMessage: "Ticket Not Found",
		})

		if err != nil {
			fmt.Println(err)

			return utils.ResponseBody{
				Code:    fiber.StatusNotFound,
				Message: "Ticket Not Found",
				Data:    fiber.Map{"status": "FAILED"},
			}
		}

		return utils.ResponseBody{
			Code:    fiber.StatusNotFound,
			Message: "Ticket Not Found",
			Data:    fiber.Map{"status": "FAILED", "pdf_url": url},
		}
	}

	// Simulate Failed External Call
	if !utils.SimulateProbability(20) {
		url, err := utils.GeneratePDF(false, payload.UserId, payload.BookingId.String(), commonStructs.FailedPDFPayload{
			ErrorMessage: "Payment Process Failed",
		})

		if err != nil {
			fmt.Println(err)

			return utils.ResponseBody{
				Code:    fiber.StatusInternalServerError,
				Message: "Failed To Book Ticket",
				Data:    fiber.Map{"status": "FAILED"},
			}
		}

		return utils.ResponseBody{
			Code:    fiber.StatusInternalServerError,
			Message: "Failed To Book Ticket",
			Data:    fiber.Map{"status": "FAILED", "pdf_url": url},
		}
	}

	// Check if the ticket is booked or ongoing
	fmt.Println(ticket.Status)
	if ticket.Status == commonStructs.Booked || ticket.Status == commonStructs.Ongoing {
		return utils.ResponseBody{
			Code:    fiber.StatusAccepted,
			Message: "Ticket Already Booked, You Are In The Waiting List",
			Data:    fiber.Map{"status": "BOOKED"},
		}
	}

	// Update ticket status
	ticket.Status = commonStructs.Ongoing

	result = db.Save(&ticket)

	if result.Error != nil {
		fmt.Println(result.Error)

		url, err := utils.GeneratePDF(false, payload.UserId, payload.BookingId.String(), commonStructs.FailedPDFPayload{
			ErrorMessage: result.Error.Error(),
		})

		if err != nil {
			fmt.Println(err)

			return utils.ResponseBody{
				Code:    fiber.StatusInternalServerError,
				Message: "Error While Updating Data From Database",
				Data:    fiber.Map{"status": "FAILED"},
			}
		}

		return utils.ResponseBody{
			Code:    fiber.StatusInternalServerError,
			Message: "Error While Updating Data From Database",
			Data:    fiber.Map{"status": "FAILED", "pdf_url": url},
		}
	}

	url := fmt.Sprintf("%s/api/v1/invoice", viper.Get("PAYMENT_BASE_URL"))
	agent := fiber.Post(url).
		Set("Content-Type", "application/json").
		Set("Authorization", fmt.Sprintf(payload.Token)).
		JSON(fiber.Map{
			"ticketId": ticket.Id,
		})

	statusCode, body, errs := agent.Bytes()
	if len(errs) > 0 {
		fmt.Println(errs)

		url, err := utils.GeneratePDF(false, payload.UserId, payload.BookingId.String(), commonStructs.FailedPDFPayload{
			ErrorMessage: "Error While Generating Payment",
		})

		if err != nil {
			fmt.Println(err)

			return utils.ResponseBody{
				Code:    fiber.StatusInternalServerError,
				Message: "Error While Generating Payment",
				Data:    fiber.Map{"status": "FAILED"},
			}
		}

		return utils.ResponseBody{
			Code:    statusCode,
			Message: "Error While Generating Payment",
			Data:    fiber.Map{"status": "FAILED", "pdf_url": url},
		}
	}

	var response fiber.Map
	err := json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)

		url, err := utils.GeneratePDF(false, payload.UserId, payload.BookingId.String(), commonStructs.FailedPDFPayload{
			ErrorMessage: "Error While Generating Payment",
		})

		if err != nil {
			fmt.Println(err)

			return utils.ResponseBody{
				Code:    fiber.StatusInternalServerError,
				Message: "Error While Generating Payment",
				Data:    fiber.Map{"status": "FAILED"},
			}
		}

		return utils.ResponseBody{
			Code:    fiber.StatusInternalServerError,
			Message: "Error While Generating Payment",
			Data:    fiber.Map{"status": "FAILED", "pdf_url": url},
		}
	}

	return utils.ResponseBody{
		Code:    fiber.StatusOK,
		Message: "Ticket Booked Successfully",
		Data:    fiber.Map{"status": "ONGOING"},
	}
}
