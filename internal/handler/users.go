package handler

import (
	"net/http"

	"github.com/caarvid/armadan/internal/armadan"
	"github.com/caarvid/armadan/internal/utils/response"
	"github.com/caarvid/armadan/web/template/partials"
	"github.com/caarvid/armadan/web/template/views"
)

// TODO: Error handling!

func ManageUsersView(us armadan.UserService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		users, err := us.All(r.Context())

		if err != nil {
			return
		}

		views.ManageUsers(users).Render(r.Context(), w)
	})
}

func EditUser(us armadan.UserService, v armadan.Validator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := v.ValidateIdParam(r, "id")
		if err != nil {
			return
		}

		user, err := us.Get(r.Context(), id)
		if err != nil {
			return
		}

		partials.EditUserRow(*user).Render(r.Context(), w)
	})
}

func CancelEditUser(us armadan.UserService, v armadan.Validator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := v.ValidateIdParam(r, "id")
		if err != nil {
			return
		}

		user, err := us.Get(r.Context(), id)
		if err != nil {
			return
		}

		partials.UserRow(*user).Render(r.Context(), w)
	})
}

func UpdateUser(us armadan.UserService, v armadan.Validator) http.Handler {
	type updateUserData struct {
		Role string `json:"role" validate:"required"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := v.ValidateIdParam(r, "id")
		if err != nil {
			return
		}

		data := updateUserData{}

		if err := v.Validate(r, &data); err != nil {
			return
		}

		user, err := us.UpdateRole(r.Context(), id, data.Role)
		if err != nil {
			return
		}

		response.
			New(w, r, partials.UserRow(*user)).
			WithSuccess("Användare uppdaterad").
			HTML()
	})
}

// func InsertAdminUser(us armadan.UserService, v armadan.Validator) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		data := schema.CreateUserParams{}
//
// 		if err := v.Validate(r, &data); err != nil {
// 			return
// 		}
//
// 		password, err := utils.GenerateHash(data.Password, nil)
// 		if err != nil {
// 			return
// 		}
//
// 		user, err := db.CreateUser(r.Context(), &schema.CreateUserParams{
// 			Email:    data.Email,
// 			Password: password.Encode(),
// 		})
//
// 		if err != nil {
// 			return
// 		}
//
// 		db.UpdateUserRole(r.Context(), &schema.UpdateUserRoleParams{
// 			Role: schema.UsersRoleEnumAdmin,
// 			ID:   user.ID,
// 		})
//
// 		w.WriteHeader(http.StatusNoContent)
// 	})
// }
