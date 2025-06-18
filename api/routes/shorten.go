package routes

import (
	"os"
	"strconv"
	"time"
	"url-shortner/database"
	"url-shortner/helpers"

	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"custom_short,omitempty"`
	Expiry      time.Duration `json:"expiry,omitempty"` // Optional field for Expiry
}

type response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"custom_short,omitempty"`
	Expiry          time.Duration `json:"expiry,omitempty"` // Optional field for Expiry
	XRateRemaining  int           `json:"x_rate_remaining"`
	XRateLimitReset time.Duration `json:"x_rate_limit_reset"`
}

func ShortenURL(c *fiber.Ctx) error {
	body := new(request)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	//implement rate limiting logic here
	r2 := database.CreateClient(1)
	defer r2.Close()

	val, err := r2.Get(database.Ctx, c.IP()).Result()
	if err == redis.Nil {
		// If the key does not exist, set it with an initial value
		_ = r2.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else {
		val, _ = r2.Get(database.Ctx, c.IP()).Result()
		valInt, _ := strconv.Atoi(val)
		if valInt <= 0 {
			limit, _ := r2.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":            "Rate limit exceeded",
				"rate_limit_reset": limit / time.Nanosecond / time.Minute,
			})
		}
	}

	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid URL format",
		})
	}

	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "URL cannot be a domain",
		})
	}

	body.URL = helpers.EnforceHTTP(body.URL)

	//logic for shorten the url
	var id string
	if body.CustomShort == "" {
		// Generate a random ID if CustomShort is not provided
		// Using a UUID library to generate a random ID

		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	r := database.CreateClient(0)
	defer r.Close()

	val, err = r.Get(database.Ctx, id).Result()
	if val != "" {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Custom short URL already exists",
		})
	}
	if body.Expiry == 0 {
		body.Expiry = 24 * 3600 * time.Second // Default expiry time
	}
	err = r.Set(database.Ctx, id, body.URL, body.Expiry).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error, cannot connect to db",
		})
	}

	r2.Decr(database.Ctx, c.IP())

	//create the repsonse
	resp := response{
		URL:             body.URL,
		CustomShort:     "",
		Expiry:          body.Expiry,
		XRateRemaining:  10,
		XRateLimitReset: 30,
	}

	val, _ = r2.Get(database.Ctx, c.IP()).Result()
	resp.XRateRemaining, _ = strconv.Atoi(val)
	ttl, _ := r2.TTL(database.Ctx, c.IP()).Result()
	resp.XRateLimitReset = ttl / time.Nanosecond / time.Minute

	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id
	return c.Status(fiber.StatusCreated).JSON(resp)
}
