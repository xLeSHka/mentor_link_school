package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
	"gitlab.prodcontest.ru/team-14/lotti/internal/repository/group"
	repositoryMentor "gitlab.prodcontest.ru/team-14/lotti/internal/repository/mentor"
	groupService "gitlab.prodcontest.ru/team-14/lotti/internal/service/group"
	mentorService "gitlab.prodcontest.ru/team-14/lotti/internal/service/mentor"
	groupsRoute "gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/handler/group"
	mentorsRoute "gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/handler/mentor"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"path/filepath"
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
var GroupRepository repository.GroupRepository
var validator *Validators.Validator
var UserService service.UserService
var MentorRepository repository.MentorRepository
var GroupService service.GroupService
var MentorService service.MentorService
var routers *ApiRouters.ApiRouters
var profile1JWT string
var unknownJWT string
var profile2JWT string

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
	GroupRepository = group.New(db)
	UserRepository = repositoryUser.New(db)
	MinioRepository = repository2.New(minioClient, config)
	MentorRepository = repositoryMentor.New(db)
	UserService = userService.New(userService.FxOpts{
		JWT:              jwt,
		UsersRepository:  UserRepository,
		MinioRepository:  MinioRepository,
		MentorRepository: MentorRepository,
		Config:           config,
	})
	MentorService = mentorService.New(mentorService.FxOpts{
		UsersRepository:  UserRepository,
		MinioRepository:  MinioRepository,
		MentorRepository: MentorRepository,
		JWT:              jwt,
		Config:           config,
	})
	GroupService = groupService.New(groupService.FxOpts{
		MinioRepository: MinioRepository,
		UserRepository:  UserRepository,
		GroupRepository: GroupRepository,
	})
	routers = ApiRouters.CreateApiRoutes(http3, jwt)

	publicRoute.PublicRoute(routers, db)
	usersRoute.UsersRoute(usersRoute.FxOpts{
		ApiRouter:       routers,
		Validator:       validator,
		MinioRepository: MinioRepository,
		UsersService:    UserService,
	})
	mentorsRoute.MentorsRoute(mentorsRoute.FxOpts{
		ApiRouter:       routers,
		Validator:       validator,
		MinioRepository: MinioRepository,
		UsersService:    UserService,
		MentorService:   MentorService,
	})
	groupsRoute.GroupsRoutes(groupsRoute.FxOpts{
		ApiRouter:       routers,
		Validator:       validator,
		MinioRepository: MinioRepository,
		UserService:     UserService,
		GroupService:    GroupService,
	})
}
func setUp() (func(), chan os.Signal, error) {
	var err error
	profile1JWT, err = jwt.CreateToken(jwtlib.MapClaims{
		"id": profile1.ID,
	}, time.Now().Add(time.Hour*24*7))
	if err != nil {
		log.Fatal(err)
	}
	unknownJWT, err = jwt.CreateToken(jwtlib.MapClaims{
		"id": uuid.New(),
	}, time.Now().Add(time.Hour*24*7))
	if err != nil {
		log.Fatal(err)
	}
	profile2JWT, err = jwt.CreateToken(jwtlib.MapClaims{
		"id": profile2.ID,
	}, time.Now().Add(time.Hour*24*7))
	if err != nil {
		log.Fatal(err)
	}
	err = db2.MigrationUp(config)
	if err != nil {
		return nil, nil, err
	}
	db.Exec("DELETE FROM public.pairs;")
	db.Exec("DELETE FROM public.help_requests;")
	db.Exec("DELETE FROM public.roles;")
	db.Exec("DELETE FROM public.users;")
	db.Exec("DELETE FROM public.groups;")
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
var profile2 models.User = models.User{
	ID:       uuid.New(),
	Name:     "Profile2",
	BIO:      &bio,
	Telegram: "@profile2",
}
var group1 models.Group = models.Group{
	ID:   uuid.New(),
	Name: "Group1",
}
var roleStudent models.Role = models.Role{
	UserID:  profile1.ID,
	GroupID: group1.ID,
	Role:    "student",
}
var roleMentor models.Role = models.Role{
	UserID:  profile2.ID,
	GroupID: group1.ID,
	Role:    "mentor",
}
var pending models.HelpRequest = models.HelpRequest{
	ID:       uuid.New(),
	UserID:   profile1.ID,
	GroupID:  group1.ID,
	MentorID: profile2.ID,
	Goal:     "Some goal",
	BIO:      &bio,
	Status:   "pending",
}
var accepted models.HelpRequest = models.HelpRequest{
	ID:       pending.ID,
	UserID:   profile1.ID,
	GroupID:  group1.ID,
	MentorID: profile2.ID,
	Goal:     "Some goal",
	BIO:      &bio,
	Status:   "accepted",
}
var rejected models.HelpRequest = models.HelpRequest{
	ID:       pending.ID,
	UserID:   profile1.ID,
	GroupID:  group1.ID,
	MentorID: profile2.ID,
	Goal:     "Some goal",
	BIO:      &bio,
	Status:   "rejected",
}
var pair models.Pair = models.Pair{
	UserID:   profile1.ID,
	GroupID:  group1.ID,
	MentorID: profile2.ID,
	Goal:     "Some goal",
}

type resGetProfile struct {
	Name      string  `json:"name"`
	AvatarUrl *string `json:"avatar_url,omitempty"`
	BIO       *string `json:"bio,omitempty"`
}
type respGetMentor struct {
	MentorID  uuid.UUID `json:"mentor_id" binding:"required"`
	GroupID   uuid.UUID `json:"group_id" binding:"required"`
	AvatarUrl *string   `json:"avatar_url,omitempty"`
	Name      string    `json:"name" binding:"required"`
	BIO       *string   `json:"bio,omitempty"`
}
type Pair struct {
	MentorID uuid.UUID `json:"mentor_id"`
	GroupId  uuid.UUID `json:"group_id"`
}
type reqCreateHelp struct {
	Requests []Pair `json:"requests" binding:"required"`
	Goal     string `json:"goal" binding:"required"`
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
				assert.Equal(t, test.Expected.Name, user.Name)
				assert.Equal(t, *test.Expected.BIO, *user.BIO)
			}
		})
	}
}
func TestAvailableMentors(t *testing.T) {
	fn, quit, err := setUp()
	assert.Nil(t, err)
	defer func() {
		quit <- syscall.SIGTERM
		fn()
	}()
	type Test struct {
		Expected     []models.User
		jwt          string
		name         string
		expectedCode int
	}
	tests := []Test{
		{
			Expected:     []models.User{profile2},
			jwt:          profile1JWT,
			name:         "get available profile 2",
			expectedCode: http.StatusOK,
		},
		{
			jwt:          unknownJWT,
			name:         "get unknown user available",
			expectedCode: http.StatusNotFound,
		},
		{
			Expected:     []models.User{},
			jwt:          profile2JWT,
			name:         "available mentors not found",
			expectedCode: http.StatusOK,
		},
	}
	db.Create(&profile1)
	db.Create(&profile2)
	db.Create(&group1)
	db.Create(&roleMentor)
	db.Create(&roleStudent)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url := "/api/user/availableMentors"
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
				var user []*respGetMentor
				err = json.Unmarshal(w.Body.Bytes(), &user)
				assert.Nil(t, err)
				for i, exp := range test.Expected {
					assert.Equal(t, exp.Name, user[i].Name)
					assert.Equal(t, *exp.BIO, *user[i].BIO)
					assert.Equal(t, exp.ID, user[i].MentorID)
					assert.Equal(t, group1.ID, user[i].GroupID)
				}
			}
		})
	}
}
func TestCreateRequest(t *testing.T) {
	fn, quit, err := setUp()
	assert.Nil(t, err)
	defer func() {
		quit <- syscall.SIGTERM
		fn()
	}()
	type Test struct {
		Request      reqCreateHelp
		jwt          string
		name         string
		expectedCode int
	}
	tests := []Test{
		{
			Request: reqCreateHelp{
				Requests: []Pair{
					{
						MentorID: profile2.ID,
						GroupId:  group1.ID,
					},
				},
				Goal: "Some goal",
			},
			jwt:          profile1JWT,
			name:         "succes create",
			expectedCode: http.StatusOK,
		},
		{
			Request: reqCreateHelp{
				Requests: []Pair{
					{
						MentorID: uuid.New(),
						GroupId:  group1.ID,
					},
				},
				Goal: "Some goal",
			},
			jwt:          unknownJWT,
			name:         "unknown user",
			expectedCode: http.StatusNotFound,
		},
		{
			Request: reqCreateHelp{
				Requests: []Pair{
					{
						MentorID: profile1.ID,
						GroupId:  group1.ID,
					},
				},
				Goal: "Some goal",
			},
			jwt:          profile2JWT,
			name:         "user is not student",
			expectedCode: http.StatusBadRequest,
		},
	}
	db.Create(&profile1)
	db.Create(&profile2)
	db.Create(&group1)
	db.Create(&roleMentor)
	db.Create(&roleStudent)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url := "/api/group/requests"
			jsonData, _ := json.Marshal(&test.Request)
			assert.Nil(t, err)
			req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData)) // bytes.NewBuffer(jsonData)
			req.Header.Set("Authorization", "Bearer "+test.jwt)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			http3.ServeHTTP(w, req)
			defer func() {
				err := w.Result().Body.Close()
				assert.Nil(t, err)
			}()
			assert.Equal(t, test.expectedCode, w.Code)
		})
	}
}
func TestUserGetRequests(t *testing.T) {
	fn, quit, err := setUp()
	assert.Nil(t, err)
	defer func() {
		quit <- syscall.SIGTERM
		fn()
	}()
	type Test struct {
		Expected     []models.HelpRequest
		jwt          string
		name         string
		expectedCode int
	}
	tests := []Test{
		{
			Expected:     []models.HelpRequest{pending},
			jwt:          profile1JWT,
			name:         "succes student get requests",
			expectedCode: http.StatusOK,
		},
		{
			jwt:          unknownJWT,
			name:         "unknown user",
			expectedCode: http.StatusNotFound,
		},
		{
			Expected:     []models.HelpRequest{},
			jwt:          profile2JWT,
			name:         "empty student requests",
			expectedCode: http.StatusOK,
		},
	}
	db.Create(&profile1)
	db.Create(&profile2)
	db.Create(&group1)
	db.Create(&roleMentor)
	db.Create(&roleStudent)
	db.Create(&pending)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url := "/api/user/requests"
			assert.Nil(t, err)
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
				var resp []respGetHelp
				err = json.Unmarshal(w.Body.Bytes(), &resp)
				assert.Nil(t, err)
				for i, exp := range test.Expected {
					assert.Equal(t, exp.ID, resp[i].ID)
					assert.Equal(t, exp.MentorID, resp[i].MentorID)
					assert.Equal(t, profile2.Name, resp[i].MentorName)
					assert.Equal(t, exp.Goal, resp[i].Goal)
					assert.Equal(t, exp.Status, resp[i].Status)
				}
			}
		})
	}
}

type respGetHelp struct {
	ID         uuid.UUID `json:"id"`
	MentorID   uuid.UUID `json:"mentor_id"`
	MentorName string    `json:"mentor_name"`
	AvatarUrl  *string   `json:"avatar_url,omitempty"`
	Goal       string    `json:"goal"`
	Status     string    `json:"status"`
}
type reqUpdateStatus struct {
	ID     uuid.UUID `json:"id"`
	Status bool      `json:"status"`
}

func TestMentorGetRequests(t *testing.T) {
	fn, quit, err := setUp()
	assert.Nil(t, err)
	defer func() {
		quit <- syscall.SIGTERM
		fn()
	}()
	type Test struct {
		Expected     []models.HelpRequest
		jwt          string
		name         string
		expectedCode int
	}
	tests := []Test{
		{
			Expected:     []models.HelpRequest{pending},
			jwt:          profile2JWT,
			name:         "succes mentor get requests",
			expectedCode: http.StatusOK,
		},
		{
			jwt:          unknownJWT,
			name:         "unknown user",
			expectedCode: http.StatusNotFound,
		},
		{
			Expected:     []models.HelpRequest{},
			jwt:          profile1JWT,
			name:         "empty mentor requests",
			expectedCode: http.StatusOK,
		},
	}
	db.Create(&profile1)
	db.Create(&profile2)
	db.Create(&group1)
	db.Create(&roleMentor)
	db.Create(&roleStudent)
	db.Create(&pending)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url := "/api/mentors/requests"
			assert.Nil(t, err)
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
				var resp []respGetRequest
				err = json.Unmarshal(w.Body.Bytes(), &resp)
				assert.Nil(t, err)
				for i, exp := range test.Expected {
					assert.Equal(t, exp.ID, resp[i].ID)
					assert.Equal(t, exp.UserID, resp[i].UserID)
					assert.Equal(t, profile1.Name, resp[i].Name)
					assert.Equal(t, exp.Goal, resp[i].Goal)
					assert.Equal(t, exp.Status, resp[i].Status)
				}
			}
		})
	}
}
func TestMentorUpdateRequests(t *testing.T) {
	fn, quit, err := setUp()
	assert.Nil(t, err)
	defer func() {
		quit <- syscall.SIGTERM
		fn()
	}()
	type Test struct {
		Expected     models.HelpRequest
		Req          reqUpdateStatus
		jwt          string
		name         string
		expectedCode int
	}
	tests := []Test{
		{
			Expected: accepted,
			Req: reqUpdateStatus{
				ID:     pending.ID,
				Status: true,
			},
			jwt:          profile2JWT,
			name:         "succes mentor update request",
			expectedCode: http.StatusOK,
		},
		{
			Req: reqUpdateStatus{
				ID:     pending.ID,
				Status: true,
			},
			jwt:          unknownJWT,
			name:         "unknown user",
			expectedCode: http.StatusBadRequest,
		},
		{
			Req: reqUpdateStatus{
				ID:     pending.ID,
				Status: true,
			},
			jwt:          profile1JWT,
			name:         "not own request",
			expectedCode: http.StatusBadRequest,
		},
	}
	db.Create(&profile1)
	db.Create(&profile2)
	db.Create(&group1)
	db.Create(&roleMentor)
	db.Create(&roleStudent)
	db.Create(&pending)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url := "/api/mentors/requests"
			jsonData, _ := json.Marshal(test.Req)
			assert.Nil(t, err)
			req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData)) // bytes.NewBuffer(jsonData)
			req.Header.Set("Authorization", "Bearer "+test.jwt)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			http3.ServeHTTP(w, req)
			defer func() {
				err := w.Result().Body.Close()
				assert.Nil(t, err)
			}()
			assert.Equal(t, test.expectedCode, w.Code)
			if test.expectedCode == http.StatusOK {
				var resp models.HelpRequest
				err = db.Model(&models.HelpRequest{}).Where("id = ?", test.Req.ID).First(&resp).Error
				assert.Nil(t, err)
				assert.Equal(t, test.Expected.Status, resp.Status)
				var p models.Pair
				err = db.Model(&models.Pair{}).Where("user_id = ? AND mentor_id = ? AND group_id = ?", resp.UserID, resp.MentorID, resp.GroupID).First(&p).Error
				assert.Nil(t, err)
			}
		})
	}
}

type respGetRequest struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	AvatarUrl *string   `json:"avatar_url,omitempty"`
	Name      string    `json:"name"`
	Goal      string    `json:"goal"`
	Status    string    `json:"status"`
}

func TestMentorGetStudents(t *testing.T) {
	fn, quit, err := setUp()
	assert.Nil(t, err)
	defer func() {
		quit <- syscall.SIGTERM
		fn()
	}()
	type Test struct {
		Expected     []models.User
		jwt          string
		name         string
		expectedCode int
	}
	tests := []Test{
		{
			Expected:     []models.User{profile1},
			jwt:          profile2JWT,
			name:         "succes mentor get students",
			expectedCode: http.StatusOK,
		},
		{
			jwt:          unknownJWT,
			name:         "unknown user",
			expectedCode: http.StatusNotFound,
		},
		{
			Expected:     []models.User{},
			jwt:          profile1JWT,
			name:         "empty mentor students",
			expectedCode: http.StatusOK,
		},
	}
	db.Create(&profile1)
	db.Create(&profile2)
	db.Create(&group1)
	db.Create(&roleMentor)
	db.Create(&roleStudent)
	db.Create(&accepted)
	db.Create(&pair)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url := "/api/mentors/students"
			assert.Nil(t, err)
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
				var resp []respGetMyStudent
				err = json.Unmarshal(w.Body.Bytes(), &resp)
				assert.Nil(t, err)
				for i, exp := range test.Expected {
					assert.Equal(t, exp.ID, resp[i].StudentID)
					assert.Equal(t, profile1.Name, resp[i].Name)
				}
			}
		})
	}
}

type respGetMyStudent struct {
	StudentID uuid.UUID `json:"student_id" binding:"required"`
	AvatarUrl *string   `json:"avatar_url,omitempty"`
	Name      string    `json:"name" binding:"required"`
}

func TestStudentGetMentors(t *testing.T) {
	fn, quit, err := setUp()
	assert.Nil(t, err)
	defer func() {
		quit <- syscall.SIGTERM
		fn()
	}()
	type Test struct {
		Expected     []models.User
		jwt          string
		name         string
		expectedCode int
	}
	tests := []Test{
		{
			Expected:     []models.User{profile2},
			jwt:          profile1JWT,
			name:         "succes mentor get students",
			expectedCode: http.StatusOK,
		},
		{
			jwt:          unknownJWT,
			name:         "unknown user",
			expectedCode: http.StatusNotFound,
		},
		{
			Expected:     []models.User{},
			jwt:          profile2JWT,
			name:         "empty mentor students",
			expectedCode: http.StatusOK,
		},
	}
	db.Create(&profile1)
	db.Create(&profile2)
	db.Create(&group1)
	db.Create(&roleMentor)
	db.Create(&roleStudent)
	db.Create(&accepted)
	db.Create(&pair)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url := "/api/user/mentors"
			assert.Nil(t, err)
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
				var resp []respGetMyMentor
				err = json.Unmarshal(w.Body.Bytes(), &resp)
				assert.Nil(t, err)
				for i, exp := range test.Expected {
					assert.Equal(t, exp.ID, resp[i].MentorID)
					assert.Equal(t, profile2.Name, resp[i].Name)
				}
			}
		})
	}
}

type respGetMyMentor struct {
	MentorID  uuid.UUID `json:"mentor_id" binding:"required"`
	AvatarUrl *string   `json:"avatar_url,omitempty"`
	Name      string    `json:"name" binding:"required"`
}

func TestMentorUploadImage(t *testing.T) {
	fn, quit, err := setUp()
	assert.Nil(t, err)
	defer func() {
		quit <- syscall.SIGTERM
		fn()
	}()
	type Test struct {
		UserID       uuid.UUID
		jwt          string
		name         string
		expectedCode int
		Formfile     string
		Imagename    string
	}
	tests := []Test{
		{
			jwt:          profile2JWT,
			name:         "add image to profile2",
			expectedCode: http.StatusOK,
			Imagename:    "go.jpg",
			Formfile:     "image",
		},
		{
			jwt:          profile2JWT,
			name:         "failed upload gif",
			expectedCode: http.StatusBadRequest,
			Imagename:    "go.gif",
			Formfile:     "image",
		},
		{
			jwt:          unknownJWT,
			name:         "unknown user",
			expectedCode: http.StatusNotFound,
			Imagename:    "go.jpg",
			Formfile:     "image",
		},
	}
	db.Create(&profile2)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			image, err := os.Open("./image/" + test.Imagename)
			assert.Nil(t, err)
			defer image.Close()

			buff := &bytes.Buffer{}
			buffWriter := io.Writer(buff)

			ext := filepath.Ext(test.Imagename)
			imageURL := profile2.ID.String() + ext

			formWriter := multipart.NewWriter(buffWriter)

			formPart, err := formWriter.CreateFormFile(test.Formfile, imageURL)
			assert.Nil(t, err)
			_, err = io.Copy(formPart, image)
			assert.Nil(t, err)
			formWriter.Close()

			req, _ := http.NewRequest(http.MethodPost, "/api/user/uploadAvatar", buff)
			req.Header.Set("Content-Type", formWriter.FormDataContentType())
			req.Header.Set("Authorization", "Bearer "+test.jwt)

			w := httptest.NewRecorder()
			http3.ServeHTTP(w, req)
			defer func() {
				err := w.Result().Body.Close()
				assert.Nil(t, err)
			}()
			assert.Equal(t, test.expectedCode, w.Code)
			if test.expectedCode == http.StatusOK {
				_, err := MinioRepository.GetImage(imageURL)
				assert.Nil(t, err)
				ext := filepath.Ext(test.Imagename)
				imagename := fmt.Sprintf("%s%s", profile2.ID.String(), ext)

				MinioRepository.DeleteImage(imagename)
			}
		})
	}
}
