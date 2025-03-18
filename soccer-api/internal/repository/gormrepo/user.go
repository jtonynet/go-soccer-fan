package gormrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/database"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/entity"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/model"
	"github.com/jtonynet/go-soccer-fan/soccer-api/internal/util"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	db *gorm.DB
}

func NewUser(gConn *database.GormConn) *User {
	return &User{
		db: gConn.GetDB(),
	}
}

func (u *User) FindByUID(ctx context.Context, uid uuid.UUID) (*entity.User, error) {
	var uModel model.User

	if err := u.db.WithContext(ctx).First(&uModel, uid).Error; err != nil {
		return nil, fmt.Errorf("user not found: %s", uid.String())
	}

	return &entity.User{
		ID:       uModel.ID,
		UID:      uModel.UID,
		UserName: uModel.Username,
		Password: "",
		Name:     uModel.Name,
		Email:    uModel.Email,
	}, nil
}

func (u *User) FindByUserName(ctx context.Context, userName string) (*entity.User, error) {
	var uModel model.User

	if err := u.db.WithContext(ctx).Where("username = ?", userName).First(&uModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found: %s", userName)
		}
		return nil, err
	}

	return &entity.User{
		ID:       uModel.ID,
		UID:      uModel.UID,
		UserName: uModel.Username,
		Password: "",
		Name:     uModel.Name,
		Email:    uModel.Email,
	}, nil
}

func (u *User) Create(ctx context.Context, eUser *entity.User) (*entity.User, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(eUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	eUser.Password = string(hashedPassword)

	uModel := model.User{
		UID:      eUser.UID,
		Username: eUser.UserName,
		Password: eUser.Password,
		Name:     eUser.Name,
		Email:    eUser.Email,
	}

	err = u.db.WithContext(ctx).Create(&uModel).Error
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:       uModel.ID,
		UID:      uModel.UID,
		UserName: uModel.Username,
		Password: "",
		Name:     uModel.Name,
		Email:    uModel.Email,
	}, nil
}

func (u *User) Login(ctx context.Context, userName, password string) (string, error) {
	var err error

	uModel := model.User{}

	err = u.db.WithContext(ctx).Model(User{}).Where("username = ?", userName).Take(&uModel).Error

	if err != nil {
		return "", err
	}

	// err = verifyPassword(password, uModel.Password)
	err = bcrypt.CompareHashAndPassword([]byte(uModel.Password), []byte(password))

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	token, err := util.GenerateToken(uModel.ID)

	if err != nil {
		return "", err
	}

	return token, nil
}

// func verifyPassword(password, hashedPassword string) error {
// 	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
// }
