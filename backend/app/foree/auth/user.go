package foree_auth

import "xue.io/go-pay/auth"

// Move the service to app folder
type UserService struct {
	UserRepo *auth.UserRepo
}

func (us *UserService) GetUserPermissions(userId int64) ([]auth.Permission, error) {
	return nil, nil
}

func (us *UserService) UpdateUserStatus(userId int64, status auth.UserStatus) error {
	return us.UserRepo.UpdateStatus(userId, status)
}

func (us *UserService) CreateNewUser() {

}
