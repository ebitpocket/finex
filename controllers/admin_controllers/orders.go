package admin_controllers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/nusa-exchange/finex/config"
	"github.com/nusa-exchange/finex/controllers/entities"
	"github.com/nusa-exchange/finex/controllers/helpers"
	"github.com/nusa-exchange/finex/controllers/queries"
	"github.com/nusa-exchange/finex/models"
	"github.com/nusa-exchange/finex/types"
	"github.com/nusa-exchange/pkg"
)

func CancelOrder(c *fiber.Ctx) error {
	uuid, err := uuid.Parse(c.Params("uuid"))
	if err != nil {
		return c.Status(422).JSON(helpers.Errors{
			Errors: []string{"admin.order.cancel_error"},
		})
	}

	var order *models.Order

	result := config.DataBase.Where("uuid = ?", uuid).First(&order)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.Status(404).JSON(helpers.Errors{
			Errors: []string{"record.not_found"},
		})
	}

	// Doing cancel
	config.KafkaProducer.Produce("matching", map[string]interface{}{
		"action": pkg.ActionCancel,
		"order":  order.ToMatchingAttributes(),
	})

	return c.Status(200).JSON(order.ToJSON())
}

func CancelAllOrders(c *fiber.Ctx) error {
	var orders []*models.Order
	params := new(queries.CancelOrderParams)

	if err := c.BodyParser(params); err != nil {
		return c.Status(500).JSON(helpers.Errors{
			Errors: []string{"server.method.invalid_message_body"},
		})
	}

	tx := config.DataBase.Where("state = ?", models.StateWait)

	if len(params.Market) > 0 {
		tx = tx.Where("market_id = ?", params.Market)
	}

	if len(params.Side) > 0 {
		var nSide models.OrderSide

		if params.Side == types.TypeBuy {
			nSide = models.SideBuy
		} else if params.Side == types.TypeSell {
			nSide = models.SideSell
		} else {
			return c.Status(422).JSON(helpers.Errors{
				Errors: []string{"admin.orders.invalid_side"},
			})
		}

		tx = tx.Where("type = ?", nSide)
	}

	tx.Find(&orders)

	for _, order := range orders {
		// Doing cancel
		config.KafkaProducer.Produce("matching", map[string]interface{}{
			"action": pkg.ActionCancel,
			"order":  order.ToMatchingAttributes(),
		})
	}

	var ordersJSON []entities.OrderEntity

	for _, order := range orders {
		ordersJSON = append(ordersJSON, order.ToJSON())
	}

	return c.Status(201).JSON(ordersJSON)
}
