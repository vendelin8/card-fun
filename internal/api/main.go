// Package api contains API functions available from HTTP.
package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vendelin8/card-fun/internal/db"
	"github.com/vendelin8/card-fun/pkg/deck"
	"github.com/vendelin8/card-fun/pkg/deck/french52"
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
	d := &french52.French52{}
	d.Shuffled = shuffled
	d.Codes = util.SplitStr(c.Query("cards", ""))
	if !util.CheckErr(c, d.New()) {
		return nil
	}
	ctx, cancel := util.CtxTimeout()
	defer cancel()
	if ok = util.CheckErr(c, db.StoreDeck(ctx, (*deck.Deck)(d))); !ok {
		return nil
	}
	return c.JSON(d.Details)
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
	d := &french52.French52{}
	d.ID = c.Query("deck_id")
	ctx, cancel := util.CtxTimeout()
	defer cancel()
	if !util.CheckErr(c, db.All(ctx, (*deck.Deck)(d))) {
		return nil
	}
	d.Resolve()
	return c.JSON(d)
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
	d := &french52.French52{}
	d.ID = c.Query("deck_id")
	n, ok := util.ParsePosInt(c, c.Query("count", "1"), "count")
	if !ok {
		return nil
	}
	ctx, cancel := util.CtxTimeout()
	defer cancel()
	if !util.CheckErr(c, db.Draw(ctx, (*deck.Deck)(d), n)) {
		return nil
	}
	d.Resolve()
	return c.JSON(map[string]interface{}{"cards": d.Cards})
}
