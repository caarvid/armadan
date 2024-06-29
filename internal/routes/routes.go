package routes

import (
	"github.com/caarvid/armadan/internal/handlers"
	m "github.com/caarvid/armadan/internal/middleware"
	"github.com/caarvid/armadan/internal/schema"
	"github.com/labstack/echo/v4"
)

func Register(app *echo.Echo, handler *handlers.Handler, db *schema.Queries) {
	authorize := m.Authorize(db)

	api := app.Group("/api")
	admin := app.Group("/admin", authorize)
	auth := app.Group("/auth")

	// Views
	app.GET("/", m.HxPushPath(handler.HomeView), authorize)
	app.GET("/schedule", m.HxPushPath(handler.ScheduleView), authorize)
	app.GET("/login", m.HxPushPath(handler.LoginView))
	app.GET("/forgot-password", m.HxPushPath(handler.ForgotPasswordView))
	app.GET("/reset-password", m.HxPushPath(handler.ResetPasswordView))
	admin.GET("", handler.AdminView)

	// Posts
	admin.GET("/posts", m.HxPushPath(handler.ManagePostsView))
	admin.GET("/posts/new", m.HxPushPath(handler.CreatePostView))
	admin.GET("/posts/edit/:id", handler.EditPost)
	admin.GET("/posts/edit/:id/cancel", handler.CancelEditPost)

	api.PUT("/posts/:id", handler.UpdatePost)
	api.POST("/posts", handler.InsertPost)
	api.DELETE("/posts/:id", handler.DeletePost)

	// Results
	admin.GET("/results", m.HxPushPath(handler.ManageResultsView))

	// Courses
	admin.GET("/courses", m.HxPushPath(handler.ManageCoursesView))
	admin.GET("/courses/new", m.HxPushPath(handler.CreateCourseView))
	admin.GET("/courses/tee/new", handler.GetEmptyTeeForm)
	admin.GET("/courses/tee/remove", handler.RemoveEmptyTeeForm)
	admin.GET("/courses/edit/:id", handler.EditCourse)
	admin.GET("/courses/edit/:id/cancel", handler.CancelEditCourse)

	api.POST("/courses", handler.InsertCourse)
	api.PUT("/courses/:id", handler.UpdateCourse)
	api.DELETE("/courses/:id", handler.DeleteCourse)
	api.DELETE("/courses/tee/:id", handler.RemoveTee)

	// Weeks
	admin.GET("/weeks", m.HxPushPath(handler.ManageWeeksView))
	admin.GET("/weeks/tees", handler.CourseTees)
	admin.GET("/weeks/:id/edit", handler.EditWeek)
	admin.GET("/weeks/:id/edit/cancel", handler.CancelEditWeek)

	api.POST("/weeks", handler.InsertWeek)
	api.PUT("/weeks/:id", handler.UpdateWeek)

	// Players
	admin.GET("/players", m.HxPushPath(handler.ManagePlayersView))
	admin.GET("/players/:id/edit", handler.EditPlayer)
	admin.GET("/players/:id/edit/cancel", handler.CancelEditPlayer)

	api.POST("/players", handler.InsertPlayer)
	api.PUT("/players/:id", handler.UpdatePlayer)
	api.DELETE("/players/:id", handler.DeletePlayer)

	// Auth
	auth.POST("/login", handler.Login)
	auth.GET("/logout", handler.Logout)
}
