package interfaces

import (
	"marketplace/internal/application"
	"marketplace/internal/domain/entity"
	"marketplace/internal/infrastructure/auth"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// new
// type UserAppInterface interface {
// }

// type UserHandler struct {
// 	UserApp UserAppInterface
// }

type UserInterface interface {
	//Login(login *model.LoginParams) (*model.S2C_Login, error)
	// GetAuthInfo(token string) (*model.AuthInfo, error)
	// Get(userID *model.UserID) (*model.S2C_UserInfo, error)
	// Register(register *model.RegisterParams) (*model.S2C_Login, error)
	// Transfer(fromUserID, toUserID *model.UserID, amount *model.Amount, currencyStr string) error

	SaveUser(c *gin.Context)
	GetUsers(c *gin.Context)
	GetUser(c *gin.Context)
}

// Users struct defines the dependencies that will be used
type Users struct {
	us application.UserAppInterface
	rd auth.AuthInterface
	tk auth.TokenInterface
}

// Users constructor
func NewUsers(us application.UserAppInterface, rd auth.AuthInterface, tk auth.TokenInterface) *Users {
	return &Users{
		us: us,
		rd: rd,
		tk: tk,
	}
}

func (s *Users) SaveUser(c *gin.Context) {
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"invalid_json": "invalid json",
		})
		return
	}
	//validate the request:
	validateErr := user.Validate("")
	if len(validateErr) > 0 {
		c.JSON(http.StatusUnprocessableEntity, validateErr)
		return
	}
	newUser, err := s.us.SaveUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, newUser.PublicUser())
}

func (s *Users) GetUsers(c *gin.Context) {
	users := entity.Users{} //customize user
	var err error
	//us, err = application.UserApp.GetUsers()
	//s.GetUsers()
	users, err = s.us.GetUsers() // interface to infrastructure
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, users.PublicUsers())
}

func (s *Users) GetUser(c *gin.Context) {
	userId, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	user, err := s.us.GetUser(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, user.PublicUser())
}
