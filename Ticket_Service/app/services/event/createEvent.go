package eventService

import (
	"fmt"

	commonStructs "github.com/Altair1618/Tubes_IF4031_03/Ticket_Service/app/common/structs"
	"github.com/Altair1618/Tubes_IF4031_03/Ticket_Service/app/configs"
	"github.com/Altair1618/Tubes_IF4031_03/Ticket_Service/app/models"
	"github.com/Altair1618/Tubes_IF4031_03/Ticket_Service/app/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func CreateEventService(payload commonStructs.CreateEventServicePayload) utils.ResponseBody {
	validator := utils.CustomValidator{
		Validator: validator.New(),
	}

	if err := validator.Validate(payload); err != nil {
		return utils.ResponseBody{
			Code:    fiber.StatusBadRequest,
			Message: utils.GetValidationErrorMessages(err)[0].Message,
		}
	}

	event := models.Event{
		EventName: payload.EventName,
		EventTime: payload.EventTime,
		Location:  payload.Location,
	}

	db, _ := configs.GetGormClient()

	result := db.Create(&event)

	if result.Error != nil {
		fmt.Println(result.Error)

		return utils.ResponseBody{
			Code:    500,
			Message: "Error While Inserting Data To Database",
			Data:    nil,
		}
	} else {
		return utils.ResponseBody{
			Code:    200,
			Message: "Event Data Inserted Successfully",
			Data:    nil,
		}
	}
}
