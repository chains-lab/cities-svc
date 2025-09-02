package rest

//type Policies struct {
//	Service  func(h http.Handler) http.Handler
//	Sysadmin func(h http.Handler) http.Handler
//	//SysadminOrMayor      func(h http.Handler) http.Handler
//	//Mayor                func(h http.Handler) http.Handler
//	//ModeratorGovOrHigher func(h http.Handler) http.Handler
//	//AnyGov               func(h http.Handler) http.Handler
//}
//
//func (s *Service) buildPolicies() Policies {
//	return Policies{
//		Service: mdlv.ServiceGrant(constant.ServiceName, s.cfg.JWT.Service.SecretKey),
//
//		Sysadmin: mdlv.RoleGrant(meta.UserCtxKey, map[string]bool{
//			roles.Admin:     true,
//			roles.SuperUser: true,
//		}),
//
//		SysadminOrMayor: m.SysadminOrGov(
//			map[string]bool{
//				constant.CityGovRoleMayor: true,
//			},
//			map[string]bool{
//				roles.Admin:     true,
//				roles.SuperUser: true,
//			},
//		),
//
//		Mayor: m.Govs(constant.CityGovRoleMayor),
//
//		ModeratorGovOrHigher: m.Govs(
//			constant.CityGovRoleMayor,
//			constant.CityGovRoleAdvisor,
//			constant.CityGovRoleModerator,
//		),
//
//		AnyGov: m.Govs(constant.GetAllCityGovsRoles()...),
//	}
//}
