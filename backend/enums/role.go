package enums

type Role int

const (
	Client Role = iota
	Admin
	SuperAdmin
)

func (w Role) String() string {
	return [...]string{"Client", "Admin", "SuperAdmin"}[w]
}

func (w Role) EnumIndex() int {
	return int(w)
}
