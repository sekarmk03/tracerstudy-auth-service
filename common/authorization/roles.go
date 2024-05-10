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
	BasePath      = "tracer_study_grpc"
	KabKotaSvc    = "KabKotaService"
	MhsBiodataSvc = "MhsBiodataService"
	PktsSvc       = "PKTSService"
	ProvinsiSvc   = "ProvinsiService"
	RespondenSvc  = "RespondenService"
	UserStudySvc  = "UserStudyService"
	ProdiSvc      = "ProdiService"
)

var roles = AccessibleRoles{
	"/" + BasePath + "." + KabKotaSvc + "/": {
		"GetAllKabKota":     {1, 2, 3, 4, 5, 6, 7},
		"GetKabKotaByIdWil": {1, 2, 3, 4, 5, 6, 7},
		"CreateKabKota":     {1, 2, 3, 4, 5, 6, 7}, // should be only 1, 2
		"UpdateKabKota":     {1, 2, 3, 4, 5, 6, 7}, // should be only 1, 2
		"DeleteKabKota":     {1, 2, 3, 4, 5, 6, 7}, // should be only 1, 2
	},
	"/" + BasePath + "." + MhsBiodataSvc + "/": {
		"FetchMhsBiodataByNim": {1, 2, 3, 4, 5, 6, 7},
	},
	"/" + BasePath + "." + PktsSvc + "/": {
		"GetAllPKTS":         {1, 2, 3, 4, 5, 6, 7},
		"GetPKTSByNim":       {1, 2, 3, 4, 5, 6, 7},
		"CreatePKTS":         {1, 2, 5, 6},
		"UpdatePKTS":         {1, 2, 5, 6},
		"GetNimByDataAtasan": {1, 2, 3, 4, 5, 6, 7},
	},
	"/" + BasePath + "." + ProdiSvc + "/": {
		"GetAllProdi":    {1, 2, 3, 4, 5, 6, 7},
		"GetAllFakultas": {1, 2, 3, 4, 5, 6, 7},
		"GetProdiByKode": {1, 2, 3, 4, 5, 6, 7},
		"CreateProdi":    {1, 2, 3, 4, 5, 6, 7}, // should be only 1, 2
		"UpdateProdi":    {1, 2, 3, 4, 5, 6, 7}, // should be only 1, 2
		"DeleteProdi":    {1, 2, 3, 4, 5, 6, 7}, // should be only 1, 2
	},
	"/" + BasePath + "." + ProvinsiSvc + "/": {
		"GetAllProvinsi":     {1, 2, 3, 4, 5, 6, 7},
		"GetProvinsiByIdWil": {1, 2, 3, 4, 5, 6, 7},
		"CreateProvinsi":     {1, 2, 3, 4, 5, 6, 7}, // should be only 1, 2
		"UpdateProvinsi":     {1, 2, 3, 4, 5, 6, 7}, // should be only 1, 2
		"DeleteProvinsi":     {1, 2, 3, 4, 5, 6, 7}, // should be only 1, 2
	},
	"/" + BasePath + "." + RespondenSvc + "/": {
		"GetAllResponden":         {1, 2, 3, 4, 5, 6, 7},
		"GetRespondenByNim":       {1, 2, 3, 4, 5, 6, 7},
		"UpdateRespondenFromSiak": {1, 2, 5, 6},
		"CreateResponden":         {1, 2, 5, 6},
		"UpdateResponden":         {1, 2, 5, 6},
		"GetRespondenByNimList":   {1, 2, 3, 4, 5, 6, 7},
	},
	"/" + BasePath + "." + UserStudySvc + "/": {
		"GetAllUserStudy":   {1, 2, 3, 4, 5, 6, 7},
		"GetUserStudyByNim": {1, 2, 3, 4, 5, 6, 7},
		"CreateUserStudy":   {1, 2, 7},
		"UpdateUserStudy":   {1, 2, 7},
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
