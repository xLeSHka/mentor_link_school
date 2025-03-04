package tests

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"gitlab.prodcontest.ru/team-14/lotti/internal/app/Validators"
	db2 "gitlab.prodcontest.ru/team-14/lotti/internal/connetions/db"
	minio2 "gitlab.prodcontest.ru/team-14/lotti/internal/connetions/minio"
	"gitlab.prodcontest.ru/team-14/lotti/internal/models"
	config2 "gitlab.prodcontest.ru/team-14/lotti/internal/pkg/config"
	"gitlab.prodcontest.ru/team-14/lotti/internal/repository"
	"gitlab.prodcontest.ru/team-14/lotti/internal/repository/group"
	repositoryMentor "gitlab.prodcontest.ru/team-14/lotti/internal/repository/mentor"
	repositoryMinio "gitlab.prodcontest.ru/team-14/lotti/internal/repository/minio"
	repositoryUser "gitlab.prodcontest.ru/team-14/lotti/internal/repository/user"
	"gitlab.prodcontest.ru/team-14/lotti/internal/service"
	groupService "gitlab.prodcontest.ru/team-14/lotti/internal/service/group"
	mentorService "gitlab.prodcontest.ru/team-14/lotti/internal/service/mentor"
	userService "gitlab.prodcontest.ru/team-14/lotti/internal/service/user"
	http2 "gitlab.prodcontest.ru/team-14/lotti/internal/transport/http"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/handler/ApiRouters"
	groupsRoute "gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/handler/group"
	mentorsRoute "gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/handler/mentor"
	publicRoute "gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/handler/public"
	usersRoute "gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/handler/user"
	"gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/handler/ws"
	jwt2 "gitlab.prodcontest.ru/team-14/lotti/internal/transport/http/pkg/jwt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var quit = make(chan os.Signal, 1)
var config config2.Config
var db *gorm.DB
var http3 *gin.Engine
var minioClient *minio.Client
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
var profile3JWT string
var wsconn *ws.WebSocket

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

	validator = Validators.New()

	jwt = jwt2.New(config)
	minioClient, err = minio2.New(config)
	if err != nil {
		log.Fatal(err)
	}
	wsconn = ws.New()
	GroupRepository = group.New(db)
	UserRepository = repositoryUser.New(db)
	MinioRepository = repositoryMinio.New(minioClient, config)
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

	publicRoute.PublicRoute(routers, db, wsconn)
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
		UsersService:    UserService,
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
	profile3JWT, err = jwt.CreateToken(jwtlib.MapClaims{
		"id": profile3.ID,
	}, time.Now().Add(time.Hour*24*7))
	wsconn = ws.New()
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

var bio string = "Profile 1 bio"
var profile1 models.User = models.User{
	ID:       uuid.New(),
	Name:     "Profile1",
	BIO:      &bio,
	Telegram: "@profile1",
}
var profile3 models.User = models.User{
	ID:       uuid.New(),
	Name:     "Profile3",
	BIO:      &bio,
	Telegram: "@profile3",
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

var roleOwner models.Role = models.Role{
	UserID:  profile3.ID,
	GroupID: group1.ID,
	Role:    "owner",
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
var stats models.GroupStat = models.GroupStat{
	StudentsCount:        1,
	MentorsCount:         1,
	HelpRequestCount:     1,
	AcceptedRequestCount: 1,
	RejectedRequestCount: 0,
	Conversion:           100,
}
var inviteCode = "uhfdc"
var group2 models.Group = models.Group{
	ID:         uuid.New(),
	Name:       "Group2",
	InviteCode: &inviteCode,
}

type reqEditGroup struct {
	Name string `json:"name" binding:"required"`
}
type resGetProfile struct {
	Name      string  `json:"name"`
	AvatarUrl *string `json:"avatar_url,omitempty"`
	BIO       *string `json:"bio,omitempty"`
}
type respGetMentor struct {
	MentorID  uuid.UUID `json:"mentor_id" binding:"required"`
	GroupIDs  []string  `json:"group_id" binding:"required"`
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
type respGetHelp struct {
	ID         uuid.UUID `json:"id"`
	MentorID   uuid.UUID `json:"mentor_id"`
	GroupIDs   []string  `json:"group_ids"`
	MentorName string    `json:"mentor_name"`
	AvatarUrl  *string   `json:"avatar_url,omitempty"`
	Goal       string    `json:"goal"`
	Status     string    `json:"status"`
}
type reqUpdateStatus struct {
	ID     uuid.UUID `json:"id"`
	Status bool      `json:"status"`
}
type respGetRequest struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	GroupIDs  []string  `json:"group_ids"`
	AvatarUrl *string   `json:"avatar_url,omitempty"`
	Name      string    `json:"name"`
	Goal      string    `json:"goal"`
	Status    string    `json:"status"`
}

type respGetMyMentor struct {
	MentorID  uuid.UUID `json:"mentor_id" binding:"required"`
	GroupIDs  []string  `json:"group_ids" binding:"required"`
	AvatarUrl *string   `json:"avatar_url,omitempty"`
	Name      string    `json:"name" binding:"required"`
}
type respGetMyStudent struct {
	StudentID uuid.UUID `json:"student_id" binding:"required"`
	GroupIDs  []string  `json:"group_ids" binding:"required"`
	AvatarUrl *string   `json:"avatar_url,omitempty"`
	Name      string    `json:"name" binding:"required"`
}
type respOtherProfile struct {
	Telegram string  `json:"telegram"`
	BIO      *string `json:"bio,omitempty"`
}
type reqEditUser struct {
	Name     string `json:"name" binding:"required"`
	Telegram string `json:"telegram,required"`
	BIO      string `json:"bio,required"`
}

type respStat struct {
	StudentsCount        int64   `json:"students_count"`
	MentorsCount         int64   `json:"mentors_count"`
	HelpRequestCount     int64   `json:"help_request_count"`
	AcceptedRequestCount int64   `json:"accepted_request_count"`
	RejectedRequestCount int64   `json:"rejected_request_count"`
	Conversion           float64 `json:"conversion"`
}
type respGetMember struct {
	UserID    uuid.UUID `json:"user_id" binding:"required"`
	AvatarUrl *string   `json:"avatar_url,omitempty"`
	Name      string    `json:"name" binding:"required"`
	Role      string    `uri:"role" binding:"required"`
}
type reqUpdateRole struct {
	Role string `json:"role" binding:"required"`
	ID   string `json:"id" binding:"required,uuid"`
}
