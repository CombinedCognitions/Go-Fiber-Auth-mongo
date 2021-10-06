package routes

import (
	"go-fiber-auth/auth"
	"go-fiber-auth/utix"
	"net/http"
	"os"

	jwtware "github.com/gofiber/jwt/v3"

	"github.com/gofiber/fiber/v2"
)

func Install(app *fiber.App) {

	secret := os.Getenv("JWT_SECRET_KEY")
	app.Post("/SignUp", auth.SignUp)
	app.Post("/LoginIn", auth.Login)

	private := app.Group("/private")
	private.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(secret),
		SigningMethod: "HS256",
		TokenLookup:   "header:Authorization",
		ErrorHandler: func(c *fiber.Ctx, e error) error {
			return c.
				Status(http.StatusUnauthorized).
				JSON(utix.NewJError(e))
		},
	}),
	)
	//private.Get("/", auth.GetUser)

	private.Get("/:id", auth.RequestInfoByID)
	private.Post("/x", auth.GetUser)
	private.Post("/checkjwt", auth.CheckJwt)

}
