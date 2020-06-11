package bank

func Must(userid string, permission string) {
	_, _ = FindUserByID(userid)

	// if !isAllowed(userid, permission) {
	// 	e := errors.New()
	// 	e.Error = err
	// 	e.UserError = errors.UserError{http.StatusForbidden, "Permission Denied."}
	// 	panic(e)
	// }
	return
}
