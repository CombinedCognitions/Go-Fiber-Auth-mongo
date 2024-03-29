package main

import (
	"go-fiber-auth/controllers"
	"go-fiber-auth/routes"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"gopkg.in/asaskevich/govalidator.v9"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Panicln(err)
	}
	govalidator.SetFieldsRequiredByDefault(true)

}

func main() {

	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024,
	})

	app.Use(logger.New())
	app.Static("/storage", "./uploads")
	app.Use(cors.New(cors.Config{
		AllowMethods:     "POST, GET, OPTIONS, PUT, DELETE",
		AllowOrigins:     "*",
		AllowHeaders:     "Origin, Content-Type, Accept, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
		AllowCredentials: true,
		ExposeHeaders:    "true",
	}))

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "helll world"})
	})

	routes.Install(app)
	defer controllers.Close()

	// controllers.Save(&newUser)
	// res, err := controllers.GetByEmail("akashs2000@gmail.com")
	// utix.CheckErorr(err)
	// fmt.Println(res["password"])
	// res, _ := controllers.GetByKey("password", "bcrytitpls")
	// if erro != nil {
	// 	fmt.Println(erro)
	// }
	// fmt.Println(res.Email)

	// alldocs := controllers.GetAll()

	// // fmt.Println(alldocs[0])
	// for i := 0; i < len(alldocs); i++ {

	// 	fmt.Println(alldocs[i])

	// }

	// res, err := controllers.Delete("6136273639da1dca9045e348")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(res.DeletedCount)

	//var newuser models.User

	//controllers.Update("email", "headflake@gmail.com", "username", "funcissogoood")

	// has, erer := security.EncryptPassword("sass")
	// utix.CheckErorr(erer)
	// fmt.Println(has, "this is the hash")

	// serror := security.VerifyPassword(has, "sass")
	// if serror == nil {
	// 	fmt.Println("PASSWORD CORRECT")
	// }
	// utix.CheckErorr(serror)

	// token, err := security.NewToken("6136e3cfc24585d52e907467")
	// utix.CheckErorr(err)
	// fmt.Println("_________ token")
	// fmt.Println(token)
	// newUser, err := controllers.GetByKey("email", "akashd2000@gmail.com")
	// fmt.Println(newUser)
	// utix.CheckErorr(err)
	log.Fatal(app.Listen(":8080"))

}
