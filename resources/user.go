package resources

// User represents a user entity
type User struct {
	ID       int
	Username string
	Password string
	Role     UserRole
}

func (u *User) hasAccess(p Permission) bool {
	return true
}

// UserRole represents a User's type, which can be one of the defined types
type UserRole struct {
	ID          int
	Name        string
	Permissions []Permission
}

// Item represents a resource
type Item struct {
	ID   int
	Name string
}

// Permission represents allowed action on a given resource
type Permission struct {
	Item   Item
	Action string
}
