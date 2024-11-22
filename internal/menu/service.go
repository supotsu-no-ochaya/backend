package menu

type MenuService struct {
	db *MenuDb
}

func (m *MenuService) getAllMenuItems() {

}

func (m *MenuService) createItem() {

}

func NewMenuService(db *MenuDb) MenuService {
	return MenuService{db: db}
}
