package actions

import (
	"fmt"
	"go-fiber-auth/controllers"
	"go-fiber-auth/models"
	"go-fiber-auth/utix"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func Updateuserdata(c *fiber.Ctx) error {

	var userinfo models.User
	var err error
	id := c.Params("id")

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	fmt.Println(user)

	userinfo, err = controllers.GetByID("_id", id)
	if err != nil {
		c.
			Status(http.StatusUnprocessableEntity).
			JSON(utix.NewJError(err))

	}
	if userinfo.ID.Hex() == claims["Id"] && userinfo.ID.Hex() == claims["Issuer"] {
		fmt.Println("both claims match")
		file, err := c.FormFile("attachment")
		if err != nil {
			fmt.Println(err, "ERRRRR")
			return c.Status(422).JSON(fiber.Map{"errors": [1]string{"We were not able upload your attachment"}})

		}
		// fmt.Print(file)

		c.SaveFile(file, fmt.Sprintf("./uploads/%s", file.Filename))

		return c.
			Status(http.StatusOK).
			JSON(fiber.Map{"message": "you file uplaoded ggz "})
	} else {

		return c.
			Status(http.StatusBadGateway).
			SendString("stop tryna access other profile")

	}

	//return c.SendString("UNAUTHORIZED")

}
