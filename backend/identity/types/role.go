package types

type UserRole string
const(
	RoleUnregistered 	UserRole = "UNREGISTERED"
	RoleUser  			UserRole = "USER"
	RoleAdmin 			UserRole = "ADMIN"
)
