package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/Validators"
	db2 "gitlab.prodcontest.ru/team-14/lotti/internal/connetions/db"
	minio2 "gitlab.prodcontest.ru/team-14/lotti/internal/connetions/minio"
	config2 "gitlab.prodcontest.ru/team-14/lotti/internal/pkg/config"
	"gitlab.prodcontest.ru/team-14/lotti/internal/repository"
	repository2 "gitlab.prodcontest.ru/team-14/lotti/internal/repository/minio"
	repositoryUser "gitlab.prodcontest.ru/team-14/lotti/internal/repository/user"
	"gitlab.prodcontest.ru/team-14/lotti/internal/service"
	userService "gitlab.prodcontest.ru/team-14/lotti/internal/service/user"
	http2 "gitlab.prodcontest.ru/team-14/lotti/internal/transport/http"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/handler/ApiRouters"
	publicRoute "gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/handler/public"
	usersRoute "gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/handler/user"
	jwt2 "gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/pkg/jwt"

	"gorm.io/gorm"
)

var quit = make(chan os.Signal, 1)
var config config2.Config
var db *gorm.DB
var http3 *gin.Engine
var minioClient *minio.Client
var rdb *redis.Client
var jwt *jwt2.JWT
var UserRepository repository.UsersRepository
var MinioRepository repository.MinioRepository
var validator *Validators.Validator
var UserService service.UserService

var routers *ApiRouters.ApiRouters
var profile1JWT string
var unknownJWT string

func init() {
	gin.SetMode(gin.TestMode)
	var err error
	config, err = config2.New()
	if err != nil {
		log.Fatal(err)
	}

	db, err = db2.New(config)
	if err != nil {
		log.Fatal(err)
	}
	err = db2.MigrationUp(config)
	if err != nil {
		log.Fatal(err)
	}
	minioClient, err = minio2.New(config)
	if err != nil {
		log.Fatal(err)
	}
	http3 = gin.Default()
	http3.Use(gin.Recovery())
	http3.Use(http2.CORSMiddleware())
	rdb = redis.NewClient(&redis.Options{
		Addr: config.RedisHost + ":" + strconv.Itoa(int(config.RedisPort)),
		DB:   0,
	})
	validator = Validators.New()
	err = rdb.Ping(context.Background()).Err()

	if err != nil {
		log.Fatal(err)
	}
	jwt = jwt2.New(config)
	minioClient, err = minio2.New(config)
	if err != nil {
		log.Fatal(err)
	}
	profile1JWT, err = jwt.CreateToken(jwtlib.MapClaims{
		"id": profile1.ID,
	}, time.Now().Add(time.Hour*6))
	if err != nil {
		log.Fatal(err)
	}
	unknownJWT, err = jwt.CreateToken(jwtlib.MapClaims{
		"id": uuid.New(),
	}, time.Now().Add(time.Hour*6))
	if err != nil {
		log.Fatal(err)
	}
	UserRepository = repositoryUser.New(db)
	MinioRepository = repository2.New(minioClient, config)

	UserService = userService.New(userService.FxOpts{
		JWT:             jwt,
		UsersRepository: UserRepository,
		MinioRepository: MinioRepository,
		Config:          config,
	})
	routers = ApiRouters.CreateApiRoutes(http3, jwt)

	publicRoute.PublicRoute(routers)
	usersRoute.UsersRoute(usersRoute.FxOpts{
		ApiRouter:       routers,
		Validator:       validator,
		MinioRepository: MinioRepository,
		UsersService:    UserService,
	})
}
func setUp() (func(), chan os.Signal, error) {
	err := db2.MigrationUp(config)
	if err != nil {
		return nil, nil, err
	}
	db.Exec("DELETE FROM public.users;")
	db.Exec("DELETE FROM public.groups;")
	db.Exec("DELETE FROM public.pairs;")
	db.Exec("DELETE FROM public.help_requests;")
	db.Exec("DELETE FROM public.roles;")
	srv := &http.Server{
		Addr:    config.ServerAddress,
		Handler: http3,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()
	quit := make(chan os.Signal, 1)
	fn := func() {
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("Server shutdown:", err)
		}
		select {
		case <-ctx.Done():
			log.Println("timeout of 1s")
		}
		log.Println("server exiting")
	}
	log.Println("server started")
	return fn, quit, nil
}

type Name struct {
	Name string `json:"name"`
}

func TestLogin(t *testing.T) {
	fn, quit, err := setUp()
	assert.Nil(t, err)
	defer func() {
		quit <- syscall.SIGTERM
		fn()
	}()
	type Test struct {
		Name         Name
		name         string
		expectedCode int
	}
	tests := []Test{
		{
			Name:         Name{"User1"},
			name:         "create user",
			expectedCode: http.StatusOK,
		},
		{
			Name:         Name{"User1"},
			name:         "login user",
			expectedCode: http.StatusOK,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(test.Name)
			url := "/api/user/auth/sign-in"
			req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData)) // bytes.NewBuffer(jsonData)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			http3.ServeHTTP(w, req)
			defer func() {
				err := w.Result().Body.Close()
				assert.Nil(t, err)
			}()
			assert.Equal(t, test.expectedCode, w.Code)
			if test.expectedCode == http.StatusOK {
				var user models.User
				err := db.Model(&models.User{}).First(&user, test.Name).Error
				assert.Nil(t, err)
				assert.Equal(t, test.Name.Name, user.Name)
			}
		})
	}
}

var bio string = "Profile 1 bio"
var profile1 models.User = models.User{
	ID:       uuid.New(),
	Name:     "Profile1",
	BIO:      &bio,
	Telegram: "@profile1",
}

type resGetProfile struct {
	Name      string  `json:"name"`
	AvatarUrl *string `json:"avatar_url,omitempty"`
	BIO       *string `json:"bio,omitempty"`
}

func TestGetProfile(t *testing.T) {
	fn, quit, err := setUp()
	assert.Nil(t, err)
	defer func() {
		quit <- syscall.SIGTERM
		fn()
	}()
	type Test struct {
		Expected     models.User
		jwt          string
		name         string
		expectedCode int
	}
	tests := []Test{
		{
			Expected:     profile1,
			jwt:          profile1JWT,
			name:         "get profile 1",
			expectedCode: http.StatusOK,
		},
		{
			jwt:          unknownJWT,
			name:         "get unknown profile",
			expectedCode: http.StatusNotFound,
		},
	}
	db.Create(&profile1)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url := "/api/user/profile"
			req, _ := http.NewRequest(http.MethodGet, url, nil) // bytes.NewBuffer(jsonData)
			req.Header.Set("Authorization", "Bearer "+test.jwt)

			w := httptest.NewRecorder()
			http3.ServeHTTP(w, req)
			defer func() {
				err := w.Result().Body.Close()
				assert.Nil(t, err)
			}()
			assert.Equal(t, test.expectedCode, w.Code)
			if test.expectedCode == http.StatusOK {
				var user resGetProfile
				err := json.Unmarshal(w.Body.Bytes(), &user)
				assert.Nil(t, err)
				assert.Equal(t, profile1.Name, user.Name)
				assert.Equal(t, *profile1.BIO, *user.BIO)
			}
		})
	}
}
