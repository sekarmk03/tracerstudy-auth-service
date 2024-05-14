package authorization

type AccessibleRoles map[string]map[string][]uint32

/*
	1. Super Admin
	2. Admin
	3. Manager
	4. Executive
	5. Prodi
	6. Alumni
	7. Pengguna Alumni
*/

const (
	BasePath = "tracer_study_grpc"
	AuthSvc  = "AuthService"
	UserSvc  = "UserService"
)

var roles = AccessibleRoles{
	"/" + BasePath + "." + AuthSvc + "/": {
		"RegisterUser": {1, 2},
	},
	"/" + BasePath + "." + UserSvc + "/": {
		"GetAllUser":  {1, 2},
		"GetUserById": {1, 2},
		"CreateUser":  {1, 2},
		"UpdateUser":  {1, 2},
		"DeleteUser":  {1, 2},
	},
}

func GetAccessibleRoles() map[string][]uint32 {
	routes := make(map[string][]uint32)

	for service, methods := range roles {
		for method, methodRoles := range methods {
			route := service + method
			routes[route] = methodRoles
		}
	}

	return routes
}
