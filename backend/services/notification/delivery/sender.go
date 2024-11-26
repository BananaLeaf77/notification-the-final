package delivery

import (
	"notification/config"
	"notification/domain"
	"notification/middleware"

	"github.com/gofiber/fiber/v2"
)

type senderHandler struct {
	suc domain.SenderUseCase
}

func NewSenderDelivery(app *fiber.App, uc domain.SenderUseCase) {
	handler := &senderHandler{
		suc: uc,
	}

	route := app.Group("/sender")
	route.Post("/send-mass", handler.sendMassHandler)
}

func NewSenderDeliveryDeploy(app *fiber.App, uc domain.SenderUseCase) {
	handler := &senderHandler{
		suc: uc,
	}

	route := app.Group("/sender")
	route.Post("/send-mass", middleware.AuthRequired(), middleware.RoleRequired("admin", "staff"), handler.sendMassHandler)
}

func (h *senderHandler) sendMassHandler(c *fiber.Ctx) error {
	var payload struct {
		IDs       []int `json:"ids"`
		SubjectID int   `json:"subject_id"`
	}

	userToken := c.Locals("user").(*domain.Claims)
	userID := userToken.UserID

	if err := c.BodyParser(&payload); err != nil {
		config.PrintLogInfo(&userToken.Username, fiber.StatusBadRequest, "sendMassHandler")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "invalid request body",
			"success": false,
		})
	}

	if err := h.suc.SendMass(c.Context(), &payload.IDs, &userID, payload.SubjectID); err != nil {
		config.PrintLogInfo(&userToken.Username, fiber.StatusInternalServerError, "sendMassHandler")

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "failed to send notifications",
			"detail":  err.Error(),
		})
	}

	config.PrintLogInfo(&userToken.Username, fiber.StatusOK, "sendMassHandler")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "notifications sent successfully",
		"success": true,
	})
}
