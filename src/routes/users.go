package routes

import "net/http"

var usersRoutes = []route{
	{
		URI:            "/users",
		Method:         http.MethodGet,
		Function:       func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("teste")) },
		Authentication: false,
	},
}
