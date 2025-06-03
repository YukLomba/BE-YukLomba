package util

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// parseUUIDParam parses UUID from path param and returns error if invalid/missing
func ParseCtxParam(ctx *fiber.Ctx, param string) (uuid.UUID, error) {
	id := ctx.Params(param)
	if id == "" {
		return uuid.Nil, fmt.Errorf("missing '%s' parameter", param)
	}
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid UUID format for '%s'", param)
	}
	return parsedID, nil
}
