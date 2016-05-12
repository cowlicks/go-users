package users

type UserDB interface {
	CreateUserTable() error
	UserExists(username string) (bool, error)
	CreateUser(uc UserCredentials) error
	VerifyCredentials(uc UserCredentials) bool
	UpdateUser(old_creds, new_creds UserCredentials) error
	DeleteUser(uc UserCredentials) error
}
