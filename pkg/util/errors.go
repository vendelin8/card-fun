package util

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
)

var ErrUnknown = errors.New("Unknown error happened. Please retry in a couple of minutes")

// Err is the error hlder for all responses through the app.
type Err struct {
	Code int
	Msg  string
}

// CheckErr checks if the given parameter has an error. In that case, prints it to the
// response and returns false.
func CheckErr(c *fiber.Ctx, err error) bool {
	if err == nil {
		return true
	}
	wrap400(c, err.Error())
	return false
}

// wrap400 JSON encodes the given error to the response.
func wrap400(c *fiber.Ctx, msg string) {
	if err := c.Status(400).JSON(Err{400, msg}); err != nil {
		log.Printf("%v: failed to set JSON content %s", err, msg)
	}
}
