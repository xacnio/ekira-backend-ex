package main

import (
	_ "ekira-backend/docs" // load API Docs files (Swagger)
	"ekira-backend/pkg/configs"
	"ekira-backend/pkg/middleware"
	"ekira-backend/pkg/routes"
	"ekira-backend/pkg/utils"
	"ekira-backend/platform/database"
	"fmt"
	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload" // load .env file automatically
	"github.com/stripe/stripe-go/v74"
	"os"
)

// @title E-Kira API
// @version 1.0
// @description This is an auto-generated API Docs.
// @termsOfService http://swagger.io/terms/
// @contact.name   API Support
// @host      api.e-kira.tk
// @BasePath  /v1
// @securityDefinitions.apikey Authentication
// @in header
// @name Authorization
// @description "Type 'Bearer TOKEN' to correctly set the API Key"
func main() {
	// Load global .ENV file config
	err := configs.LoadConfig()
	if err != nil {
		fmt.Println(err)
	}

	// Stripe key initialization
	stripe.Key = os.Getenv("STRIPE_KEY")

	// Timezone initialization
	utils.LoadTimezone()

	// Create folder if not exists
	folder := "public/photos"
	os.MkdirAll(folder, os.ModePerm)
	folder = "public/photos/rental-house"
	os.MkdirAll(folder, os.ModePerm)
	folder = "public/photos/profile"
	os.MkdirAll(folder, os.ModePerm)
	folder = "logs"
	os.MkdirAll(folder, os.ModePerm)

	// Check redis connection
	rds := database.NewRConnection()
	err = rds.RPing()
	if err != nil {
		panic(err)
	}

	// Migrate database tables
	db, err := database.OpenDBConnection()
	if err != nil {
		panic(err)
	}
	sqlDb, _ := db.DB.DB()
	err = database.Migrate(db.DB)
	if err != nil {
		fmt.Println(err)
	}
	err = sqlDb.Close()
	if err != nil {
		fmt.Println(err)
	}

	// Define Fiber config
	config := configs.FiberConfig()

	// Define a new Fiber app with config
	app := fiber.New(config)

	// Middlewares
	middleware.FiberMiddleware(app) // Register Fiber's middleware for app.

	// Routes
	routes.SwaggerRoute(app)        // Register a route for API Docs (Swagger).
	routes.AddressRoutes(app)       // Register a route group for address routes.
	routes.RentalHouseRoutes(app)   // Register a route group for rental house routes.
	routes.AuthRoutes(app)          // Register an auth routes for app.
	routes.UserRoutes(app)          // Register a route group for user routes.
	routes.ReservationRoutes(app)   // Register a route group for reservation routes.
	routes.PaymentRoutes(app)       // Register a route group for payment routes.
	routes.StripeWebhookRoutes(app) // Register a route group for stripe webhook routes.
	routes.WalletRoutes(app)        // Register a route group for wallet routes.
	routes.NotFoundRoute(app)       // Register route for 404 Error.

	// Start server
	utils.StartServer(app)
}
