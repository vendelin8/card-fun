// Package api contains API functions available from HTTP.
package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/vendelin8/card-fun/pkg/util"
)

// @title Card Fun Demo
// @version 1.0
// @description This is a card deck simulator server demo.
// @termsOfService http://swagger.io/terms/
// @BasePath /v1

// CreateHandler creates a new deck.
// @Summary Creates a new deck.
// @Description new full deck, or one based on a given set of cards, maybe shuffled.
// @Tags deck
// @Param shuffled query boolean false "if new deck needs to be shuffled"
// @Param cards query string false "komma separated list of card codes to init deck with"
// @Accept */*
// @Produce json
// @Success 200 {object} string
// @Router /create [put]
func CreateHandler(c *fiber.Ctx) error {
	shuffled, ok := util.ParseBool(c, c.Query("shuffled", "false"), "shuffled")
	if !ok {
		return nil
	}
	codes := util.SplitStr(c.Query("cards", ""))
	fmt.Println("create", shuffled, codes) // TODO: implement
	c.WriteString("OK")
	return nil
}

// OpenHandler opens a deck.
// @Summary Opens a deck.
// @Description returns current content of a deck.
// @Tags deck
// @Param deck_id query string true "identifier of the deck about to open"
// @Accept */*
// @Produce json
// @Success 200 {object} string
// @Router /open [get]
func OpenHandler(c *fiber.Ctx) error {
	deckID := c.Query("deck_id")
	fmt.Println("open", deckID) // TODO: implement
	c.WriteString("OK")
	return nil
}

// DrawHandler draws some cards from a deck.
// @Summary Draws some cards from a deck.
// @Description Pops one or multiple cards from the top of a given deck.
// @Tags cards
// @Param deck_id query string true "identifier of the deck about to draw cards from"
// @Param count query integer false "number of cards to draw from the deck"
// @Accept */*
// @Produce json
// @Success 200 {object} string
// @Router /draw [post]
func DrawHandler(c *fiber.Ctx) error {
	deckID := c.Query("deck_id")
	count, ok := util.ParsePosInt(c, c.Query("count", "1"), "count")
	if !ok {
		return nil
	}
	fmt.Println("draw", deckID, count) // TODO: implement
	c.WriteString("OK")
	return nil
}
