package role

type Roles struct {
	ID   int
	Name string
}

type HasRole struct {
	UUID   string
	RoleID int
}
