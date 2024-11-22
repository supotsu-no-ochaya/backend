package user

type UserService struct {
	db *UserDb
}

func (svc *UserService) register(name string) {

}

func NewUserService(db *UserDb) UserService {
	return UserService{db: db}
}
