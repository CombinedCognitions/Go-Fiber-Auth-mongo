package auth

import (
	"fmt"
	"go-fiber-auth/controllers"
	"go-fiber-auth/models"
	"go-fiber-auth/security"
	"go-fiber-auth/utix"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/asaskevich/govalidator.v9"
)

func SignUp(ctx *fiber.Ctx) error {

	var newuser models.User

	err := ctx.BodyParser(&newuser)
	if err != nil {
		ctx.
			Status(http.StatusUnprocessableEntity).
			JSON(utix.NewJError(err))
	}

	newuser.Email = utix.NormalizeEmail(newuser.Email)

	if !govalidator.IsEmail(newuser.Email) {

		return ctx.
			Status(http.StatusBadRequest).
			JSON(utix.NewJError(utix.ErrUnknown))

	}

	exist, er := controllers.GetByEmail(newuser.Email)
	utix.CheckErorr(er)

	if exist.Email == "" {

		if strings.TrimSpace(newuser.Password) == "" {
			return ctx.
				Status(http.StatusBadRequest).
				JSON(utix.NewJError(utix.ErrEmptyPassword))

		}
		if len(newuser.Password) < 5 {
			fmt.Println("short password")
			return ctx.
				Status(http.StatusBadRequest).
				JSON(utix.NewJError(utix.ErrShortPassword))
		}

		newuser.Password, err = security.EncryptPassword(newuser.Password)
		if err != nil {
			return ctx.
				Status(http.StatusBadRequest).
				JSON(utix.NewJError(err))

		}

		newuser.CreatedAt = time.Now()
		newuser.UpdatedAt = newuser.CreatedAt
		newuser.ID = primitive.NewObjectID()
		err := controllers.Save(&newuser)
		if err != nil {
			return ctx.
				Status(http.StatusBadRequest).
				JSON(utix.NewJError(err))
		}

		token, exp, err := security.NewToken(newuser)
		if err != nil {
			return err
		}

		fmt.Println(token)
		fmt.Println(exp)

		return ctx.
			Status(http.StatusCreated).
			JSON(fiber.Map{"token": token, "exp": exp, "user": models.User{Email: newuser.Email, ID: newuser.ID, CreatedAt: newuser.CreatedAt}})
	}

	if exist.Email != "" {

		err = utix.ErrEmailAlreadyExists

		return ctx.Status(http.StatusBadRequest).JSON(utix.NewJError(err))

	}

	return err
}
func Login(ctx *fiber.Ctx) error {

	var input models.User

	err := ctx.BodyParser(&input)
	if err != nil {
		ctx.
			Status(http.StatusUnprocessableEntity).
			JSON(utix.NewJError(err))
	}

	input.Email = utix.NormalizeEmail(input.Email)

	if !govalidator.IsEmail(input.Email) {
		return ctx.
			Status(http.StatusBadRequest).
			JSON(utix.NewJError(utix.ErrInvalidEmail))

	}

	userinfo, er := controllers.GetByEmail(input.Email)
	utix.CheckErorr(er)
	fmt.Println(userinfo)

	if er != nil {
		log.Println("login failed")
		return ctx.
			Status(http.StatusUnauthorized).
			JSON(utix.NewJError(utix.ErrIncorrectEmail))

	}

	err = security.VerifyPassword(userinfo.Password, input.Password)
	if err == nil {

		token, exp, err := security.NewToken(userinfo)
		if err != nil {
			return ctx.
				Status(http.StatusUnprocessableEntity).
				JSON(utix.NewJError(err))

		}
		cookie := new(fiber.Cookie)
		cookie.Name = "jwt"
		cookie.SameSite = "None"
		cookie.Value = fmt.Sprintf("Bearer %s", token)
		cookie.HTTPOnly = true

		cookie.Expires = time.Now().Add(24 * time.Hour)

		ctx.Cookie(cookie)

		return ctx.
			Status(http.StatusOK).
			JSON(fiber.Map{"token": fmt.Sprintf("Bearer %s", token), "exp": exp, "user": models.User{Email: userinfo.Email, ID: userinfo.ID, CreatedAt: userinfo.CreatedAt}})

	}
	if err != nil {

		log.Println(err)

		return ctx.
			Status(http.StatusUnauthorized).
			JSON(utix.NewJError(utix.ErrIncorrectPassword))

	}

	return er
}

func GetUser(ctx *fiber.Ctx) error {
	var input models.User

	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	fmt.Println(user, claims)

	err := ctx.BodyParser(&input)
	if err != nil {
		ctx.
			Status(http.StatusUnprocessableEntity).
			JSON(utix.NewJError(err))
	}

	datafromDB, err := controllers.GetByEmail(input.Email)

	fmt.Println(datafromDB)
	if err != nil {
		return ctx.
			Status(http.StatusUnprocessableEntity).
			JSON(utix.NewJError(err))

	}

	if input.ID.Hex() == claims["Id"] && input.ID.Hex() == claims["Issuer"] {
		fmt.Println("both claims match")
		return ctx.
			Status(http.StatusOK).
			JSON(fiber.Map{"data": datafromDB, "message": "damm you unlcoked private routes"})
	}
	fmt.Println("sad that dint work lol")
	return ctx.
		Status(http.StatusUnprocessableEntity).
		JSON(utix.NewJError(utix.ErrInvalidCredentials))

}

func CheckJwt(ctx *fiber.Ctx) error {
	var input models.User

	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	fmt.Println(user, claims)

	err := ctx.BodyParser(&input)
	if err != nil {
		ctx.
			Status(http.StatusUnprocessableEntity).
			JSON(utix.NewJError(err))
	}

	// fmt.Println(datafromDB)
	if err != nil {
		return ctx.
			Status(http.StatusUnprocessableEntity).
			JSON(utix.NewJError(err))

	}

	if input.ID.Hex() == claims["Id"] && input.ID.Hex() == claims["Issuer"] {
		fmt.Println("both claims match")
		return ctx.
			Status(http.StatusOK).
			JSON(fiber.Map{"message": "LOGGEDIN"})
	}
	fmt.Println("sad that dint work lol")
	return ctx.
		Status(http.StatusUnprocessableEntity).
		JSON(utix.NewJError(utix.ErrLogout))

}

func RequestInfoByID(ctx *fiber.Ctx) error {

	var userinfo models.User
	id := ctx.Params("id")

	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	fmt.Println(user)

	userinfo, err := controllers.GetByID("_id", id)
	if err != nil {
		ctx.
			Status(http.StatusUnprocessableEntity).
			JSON(utix.NewJError(err))

	}
	if userinfo.ID.Hex() == claims["Id"] && userinfo.ID.Hex() == claims["Issuer"] {
		fmt.Println("both claims match", userinfo.ID, claims["Id"])
		return ctx.
			Status(http.StatusOK).
			JSON(userinfo)
	}

	return ctx.SendString("UNAUTHORIZED")
}

func CheckExpiredToken(ctx *fiber.Ctx) error {

	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	// claimtime := claims["exp"]
	nice, err := fmt.Println(claims["Id"], time.Now().Unix())

	fmt.Println(nice)
	// if claimtime > time.Now().Unix() {

	// }

	return err
}