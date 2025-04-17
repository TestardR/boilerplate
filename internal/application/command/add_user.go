package command

type AddUser struct {
	username string
}

func NewAddUser(username string) AddUser {
	return AddUser{
		username: username,
	}
}

func (c AddUser) Username() string {
	return c.username
}
