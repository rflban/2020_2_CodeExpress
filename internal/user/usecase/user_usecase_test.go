package usecase_test

import (
	"database/sql"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/consts"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/models"
	. "github.com/go-park-mail-ru/2020_2_CodeExpress/internal/tools/error_response"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/user/mock_user"
	"github.com/go-park-mail-ru/2020_2_CodeExpress/internal/user/usecase"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserUsecase_CreateUser_Passed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_user.NewMockUserRep(ctrl)
	mockUsecase := usecase.NewUserUsecase(mockRepo)

	user := &models.User{
		Name:     "Username",
		Email:    "email@mail.ru",
		Password: "qwertyuiop123",
	}

	mockRepo.
		EXPECT().
		SelectByNameOrEmail(gomock.Eq(user.Name), gomock.Eq(user.Email)).
		Return([]*models.User{}, nil)

	mockRepo.
		EXPECT().
		Insert(gomock.Eq(user.Name), gomock.Eq(user.Email), gomock.Eq(user.Password)).
		DoAndReturn(func(name string, email string, password string) (*models.User, error) {
			user.ID = 1
			return user, nil
		})

	createdUser, err := mockUsecase.Create(user.Name, user.Email, user.Password)
	assert.Nil(t, err)
	assert.Equal(t, user.Name, createdUser.Name)
	assert.Equal(t, user.Email, createdUser.Email)
	assert.Equal(t, user.Password, createdUser.Password)
}

func TestUserUsecase_CreateUser_FailedName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_user.NewMockUserRep(ctrl)
	mockUsecase := usecase.NewUserUsecase(mockRepo)

	user := &models.User{
		Name:     "Username",
		Email:    "email@mail.ru",
		Password: "qwertyuiop123",
	}

	existingUser := &models.User{
		ID:       1,
		Name:     "Username",
		Email:    "email2@mail.ru",
		Password: "qwertyuiop123",
	}

	mockRepo.
		EXPECT().
		SelectByNameOrEmail(gomock.Eq(user.Name), gomock.Eq(user.Email)).
		Return([]*models.User{existingUser}, nil)

	createdUser, err := mockUsecase.Create(user.Name, user.Email, user.Password)
	assert.Equal(t, err, NewErrorResponse(ErrNameAlreadyExist, nil))
	assert.Nil(t, createdUser)
}

func TestUserUsecase_CreateUser_FailedEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_user.NewMockUserRep(ctrl)
	mockUsecase := usecase.NewUserUsecase(mockRepo)

	user := &models.User{
		Name:     "Username",
		Email:    "email@mail.ru",
		Password: "qwertyuiop123",
	}

	existingUser := &models.User{
		ID:       1,
		Name:     "Username2",
		Email:    "email@mail.ru",
		Password: "qwertyuiop123",
	}

	mockRepo.
		EXPECT().
		SelectByNameOrEmail(gomock.Eq(user.Name), gomock.Eq(user.Email)).
		Return([]*models.User{existingUser}, nil)

	createdUser, err := mockUsecase.Create(user.Name, user.Email, user.Password)
	assert.Equal(t, err, NewErrorResponse(ErrEmailAlreadyExist, nil))
	assert.Nil(t, createdUser)
}

func TestUserUsecase_GetById_Passed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_user.NewMockUserRep(ctrl)
	mockUsecase := usecase.NewUserUsecase(mockRepo)

	expectedUser := &models.User{
		ID:       1,
		Name:     "Username",
		Email:    "email@mail.ru",
		Password: "qwertyuiop123",
	}

	mockRepo.
		EXPECT().
		SelectById(gomock.Eq(expectedUser.ID)).
		Return(expectedUser, nil)

	user, err := mockUsecase.GetById(expectedUser.ID)
	assert.Nil(t, err)
	assert.Equal(t, user, expectedUser)
}

func TestUserUsecase_GetById_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_user.NewMockUserRep(ctrl)
	mockUsecase := usecase.NewUserUsecase(mockRepo)

	mockRepo.
		EXPECT().
		SelectById(gomock.Eq(uint64(1))).
		Return(nil, sql.ErrNoRows)

	user, err := mockUsecase.GetById(uint64(1))
	assert.Equal(t, err, NewErrorResponse(ErrNotAuthorized, sql.ErrNoRows))
	assert.Nil(t, user)
}

func TestUserUsecase_LoginUser_Passed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_user.NewMockUserRep(ctrl)
	mockUsecase := usecase.NewUserUsecase(mockRepo)

	expectedUser := &models.User{
		ID:       1,
		Name:     "Username",
		Email:    "email@mail.ru",
		Password: "qwertyuiop123",
	}

	mockRepo.
		EXPECT().
		SelectByLogin(gomock.Eq(expectedUser.Name)).
		Return(expectedUser, nil)

	mockRepo.
		EXPECT().
		SelectByLogin(gomock.Eq(expectedUser.Email)).
		Return(expectedUser, nil)

	user, err := mockUsecase.GetUserByLogin(expectedUser.Name, "qwertyuiop123") //TODO: копипаст кода? Сделать TestCase?
	assert.Nil(t, err)
	assert.Equal(t, user, expectedUser)

	user, err = mockUsecase.GetUserByLogin(expectedUser.Email, "qwertyuiop123")
	assert.Nil(t, err)
	assert.Equal(t, user, expectedUser)
}

func TestUserUsecase_LoginUser_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_user.NewMockUserRep(ctrl)
	mockUsecase := usecase.NewUserUsecase(mockRepo)

	mockRepo.
		EXPECT().
		SelectByLogin(gomock.Eq("Username")).
		Return(nil, sql.ErrNoRows)

	mockRepo.
		EXPECT().
		SelectByLogin(gomock.Eq("email@mail.ru")).
		Return(nil, sql.ErrNoRows)

	user, err := mockUsecase.GetUserByLogin("Username", "qwertyuiop123")
	assert.Equal(t, err, NewErrorResponse(ErrIncorrectLoginOrPassword, nil))
	assert.Nil(t, user)

	user, err = mockUsecase.GetUserByLogin("email@mail.ru", "qwertyuiop123")
	assert.Equal(t, err, NewErrorResponse(ErrIncorrectLoginOrPassword, nil))
	assert.Nil(t, user)
}

func TestUserUsecase_LoginUser_FailedPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_user.NewMockUserRep(ctrl)
	mockUsecase := usecase.NewUserUsecase(mockRepo)

	expectedUser := &models.User{
		ID:       1,
		Name:     "Username",
		Email:    "email@mail.ru",
		Password: "qwertyuiop123",
	}

	mockRepo.
		EXPECT().
		SelectByLogin(gomock.Eq(expectedUser.Name)).
		Return(expectedUser, nil)

	mockRepo.
		EXPECT().
		SelectByLogin(gomock.Eq(expectedUser.Email)).
		Return(expectedUser, nil)

	user, err := mockUsecase.GetUserByLogin(expectedUser.Name, "qwertyuiop")
	assert.Equal(t, err, NewErrorResponse(ErrIncorrectLoginOrPassword, nil))
	assert.Nil(t, user)

	user, err = mockUsecase.GetUserByLogin(expectedUser.Email, "qwertyuiop")
	assert.Equal(t, err, NewErrorResponse(ErrIncorrectLoginOrPassword, nil))
	assert.Nil(t, user)
}

func TestUserUsecase_UpdateProfile_Passed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_user.NewMockUserRep(ctrl)
	mockUsecase := usecase.NewUserUsecase(mockRepo)

	user := &models.User{
		ID:    1,
		Name:  "Username2",
		Email: "email2@mail.ru",
	}

	mockRepo.
		EXPECT().
		SelectById(gomock.Eq(user.ID)).
		Return(user, nil)

	mockRepo.
		EXPECT().
		SelectByNameOrEmail(gomock.Eq(user.Name), gomock.Eq(user.Email)).
		Return([]*models.User{}, nil)

	mockRepo.
		EXPECT().
		Update(gomock.Eq(user)).
		Return(nil)

	_, err := mockUsecase.UpdateProfile(user.ID, user.Name, user.Email)
	assert.Nil(t, err)
}

func TestUserUsecase_UpdateProfile_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_user.NewMockUserRep(ctrl)
	mockUsecase := usecase.NewUserUsecase(mockRepo)

	updatedUser := &models.User{
		ID:    1,
		Name:  "Username",
		Email: "email@mail.ru",
	}

	existingUser1 := &models.User{ //"Пользователь с данным Name уже существует"
		ID:    2,
		Name:  "Username",
		Email: "email2@mail.ru",
	}

	existingUser2 := &models.User{ //"Пользователь с данным Email уже существует"
		ID:    3,
		Name:  "Username2",
		Email: "email@mail.ru",
	}

	mockRepo.
		EXPECT().
		SelectById(gomock.Eq(updatedUser.ID)).
		Return(updatedUser, nil)

	mockRepo.
		EXPECT().
		SelectByNameOrEmail(gomock.Eq(updatedUser.Name), gomock.Eq(updatedUser.Email)).
		Return([]*models.User{existingUser1}, nil)

	_, err := mockUsecase.UpdateProfile(updatedUser.ID, updatedUser.Name, updatedUser.Email)
	assert.Equal(t, err, NewErrorResponse(ErrNameAlreadyExist, nil))

	mockRepo.
		EXPECT().
		SelectById(gomock.Eq(updatedUser.ID)).
		Return(updatedUser, nil)

	mockRepo.
		EXPECT().
		SelectByNameOrEmail(gomock.Eq(updatedUser.Name), gomock.Eq(updatedUser.Email)).
		Return([]*models.User{existingUser2}, nil)

	_, err = mockUsecase.UpdateProfile(updatedUser.ID, updatedUser.Name, updatedUser.Email)
	assert.Equal(t, err, NewErrorResponse(ErrEmailAlreadyExist, nil))
}

func TestUserUsecase_UpdatePassword_Passed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_user.NewMockUserRep(ctrl)
	mockUsecase := usecase.NewUserUsecase(mockRepo)

	user := &models.User{
		ID:       1,
		Password: "qwertyuiop",
	}

	updatedUser := &models.User{
		ID:       1,
		Password: "qwertyuiop123",
	}

	mockRepo.
		EXPECT().
		SelectById(gomock.Eq(updatedUser.ID)).
		Return(user, nil)

	mockRepo.
		EXPECT().
		Update(gomock.Eq(user)).
		Return(nil)

	err := mockUsecase.UpdatePassword(updatedUser.ID, user.Password, updatedUser.Password)
	assert.Nil(t, err)
	//assert.Equal(t, resultUser, updatedUser)
}

func TestUserUsecase_UpdatePassword_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_user.NewMockUserRep(ctrl)
	mockUsecase := usecase.NewUserUsecase(mockRepo)

	user := &models.User{
		ID:       1,
		Password: "qwertyuiop",
	}

	mockRepo.
		EXPECT().
		SelectById(gomock.Eq(user.ID)).
		Return(user, nil)

	err := mockUsecase.UpdatePassword(user.ID, user.Password, user.Password)
	assert.Equal(t, err, NewErrorResponse(ErrNewPasswordIsOld, nil))

	mockRepo.
		EXPECT().
		SelectById(gomock.Eq(user.ID)).
		Return(user, nil)

	err = mockUsecase.UpdatePassword(user.ID, "qwertyuiop123", user.Password)
	assert.Equal(t, err, NewErrorResponse(ErrWrongOldPassword, nil))
}
