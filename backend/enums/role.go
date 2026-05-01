package enums

type Role int

const (
	User Role = iota
	Admin
	SuperAdmin
)

func (w Role) String() string {
	return [...]string{"User", "Admin", "SuperAdmin"}[w]
}

func (w Role) EnumIndex() int {
	return int(w)
}
