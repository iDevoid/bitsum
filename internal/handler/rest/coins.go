package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/iDevoid/bitsum/internal/constant/model"
	"github.com/iDevoid/bitsum/internal/module/coins"
)

//go:generate mockgen -destination=../../../mocks/coins/handler_mock.go -package=coins_mock -source=coins.go

// CoinsHandler contains the function of handler for domain coins
type CoinsHandler interface {
	// Payment is the handler for payment
	Payment(ctx *fiber.Ctx) error
	// Receive is the handler for receive/increment amount
	Receive(ctx *fiber.Ctx) error
	History(ctx *fiber.Ctx) error
	Wallet(ctx *fiber.Ctx) error
}

type coinsHandler struct {
	coinsCase coins.Usecase
}

// CoinsInit is to initialize the rest handler for domain coins
func CoinsInit(coinsCase coins.Usecase) CoinsHandler {
	return &coinsHandler{
		coinsCase,
	}
}

// Payment is the handler for payment to decrement amount
func (ch *coinsHandler) Payment(ctx *fiber.Ctx) error {
	var body model.Coin
	err := json.Unmarshal(ctx.Body(), &body)
	if err != nil || body.Amount == 0 || body.DateTime.IsZero() {
		ctx.Status(http.StatusBadRequest)
		return errors.New(http.StatusText(http.StatusBadRequest))
	}

	response := model.Response{
		Success: true,
		Message: "Payment Success!",
	}
	err = ch.coinsCase.Pay(ctx.UserContext(), &body)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		response = model.Response{
			Success: false,
			Message: err.Error(),
		}
	}

	raw, errMar := json.Marshal(response)
	if errMar != nil {
		return err
	}
	return ctx.Send(raw)
}

// Receive is the handler for receive/increment amount
func (ch *coinsHandler) Receive(ctx *fiber.Ctx) error {
	var body model.Coin
	err := json.Unmarshal(ctx.Body(), &body)
	if err != nil || body.Amount == 0 || body.DateTime.IsZero() {
		ctx.Status(http.StatusBadRequest)
		return errors.New(http.StatusText(http.StatusBadRequest))
	}

	response := model.Response{
		Success: true,
		Message: fmt.Sprintf("%f coin received successfully!", body.Amount),
	}
	err = ch.coinsCase.Receive(ctx.UserContext(), &body)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		response = model.Response{
			Success: false,
			Message: err.Error(),
		}
	}

	raw, err := json.Marshal(response)
	if err != nil {
		return err
	}
	return ctx.Send(raw)
}

// History is the handler to show the transactional data history between 2 dates
func (ch *coinsHandler) History(ctx *fiber.Ctx) error {
	var body model.FilterDate
	err := json.Unmarshal(ctx.Body(), &body)
	if err != nil || body.StartDateTime.IsZero() || body.EndDateTime.IsZero() {
		ctx.Status(http.StatusBadRequest)
		return errors.New(http.StatusText(http.StatusBadRequest))
	}

	res, err := ch.coinsCase.HistoryTransaction(ctx.UserContext(), &body)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		ctx.Response().SetBodyString(err.Error())
		return err
	}

	raw, err := json.Marshal(res)
	if err != nil {
		return err
	}
	return ctx.Send(raw)
}

// Wallet is to show the data of the lastest transaction date and the latest amount
func (ch *coinsHandler) Wallet(ctx *fiber.Ctx) error {
	res, err := ch.coinsCase.Balance(ctx.UserContext())
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		ctx.Response().SetBodyString(err.Error())
		return err
	}

	raw, err := json.Marshal(res)
	if err != nil {
		return err
	}

	return ctx.Send(raw)
}
