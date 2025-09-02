package integration

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/kulikovroman08/reviewlink-backend/internal/controller"
	"github.com/kulikovroman08/reviewlink-backend/internal/repository/user"
	userService "github.com/kulikovroman08/reviewlink-backend/internal/service/user"
)

type TestSetup struct {
	App *gin.Engine
	DB  *pgxpool.Pool
}

func NewTestSetup() *TestSetup {
	gin.SetMode(gin.TestMode)

	_ = godotenv.Load(".env.test")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbURL := os.Getenv("DB_URL")
	db, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("failed to connect to test DB: %v", err)
	}

	userRepo := user.NewPostgresUserRepository(db)
	userSrv := userService.NewUserService(userRepo)

	app := &controller.Application{
		UserService: userSrv,
	}

	r := gin.Default()
	r.POST("/signup", app.Signup)

	return &TestSetup{
		App: r,
		DB:  db,
	}
}

func (ts *TestSetup) TruncateUsers() {
	_, err := ts.DB.Exec(context.Background(), "TRUNCATE users RESTART IDENTITY CASCADE;")
	if err != nil {
		log.Fatalf("failed to truncate users table: %v", err)
	}
}
