package routes

import (
	"github.com/caarvid/armadan/internal/handlers"
	"github.com/labstack/echo/v4"
)

func Register(app *echo.Echo, handler *handlers.Handler) {
	api := app.Group("/api")
	admin := app.Group("/admin")

	// Views
	app.GET("/", hxPushPath(handler.HomeView))
	app.GET("/schedule", hxPushPath(handler.ScheduleView))
	app.GET("/login", hxPushPath(handler.LoginView))
	app.GET("/forgot-password", hxPushPath(handler.ForgotPasswordView))
	admin.GET("", handler.AdminView)

	// Posts
	admin.GET("/posts", hxPushPath(handler.ManagePostsView))
	admin.GET("/posts/new", hxPushPath(handler.CreatePostView))
	admin.GET("/posts/edit/:id", handler.EditPost)
	admin.GET("/posts/edit/:id/cancel", handler.CancelEditPost)

	api.PUT("/posts/:id", handler.UpdatePost)
	api.POST("/posts", handler.InsertPost)
	api.DELETE("/posts/:id", handler.DeletePost)

	// Results
	admin.GET("/results", hxPushPath(handler.ManageResultsView))

	// Courses
	admin.GET("/courses", hxPushPath(handler.ManageCoursesView))
	admin.GET("/courses/new", hxPushPath(handler.CreateCourseView))
	admin.GET("/courses/tee/new", handler.GetEmptyTeeForm)
	admin.GET("/courses/tee/remove", handler.RemoveEmptyTeeForm)
	admin.GET("/courses/edit/:id", handler.EditCourse)
	admin.GET("/courses/edit/:id/cancel", handler.CancelEditCourse)

	api.POST("/courses", handler.InsertCourse)
	api.PUT("/courses/:id", handler.UpdateCourse)
	api.DELETE("/courses/:id", handler.DeleteCourse)
	api.DELETE("/courses/tee/:id", handler.RemoveTee)

	// Weeks
	admin.GET("/weeks", hxPushPath(handler.ManageWeeksView))
	admin.GET("/weeks/tees", handler.CourseTees)
	admin.GET("/weeks/:id/edit", handler.EditWeek)
	admin.GET("/weeks/:id/edit/cancel", handler.CancelEditWeek)

	api.POST("/weeks", handler.InsertWeek)
	api.PUT("/weeks/:id", handler.UpdateWeek)

	// Players
	admin.GET("/players", hxPushPath(handler.ManagePlayersView))
	admin.GET("/players/:id/edit", handler.EditPlayer)
	admin.GET("/players/:id/edit/cancel", handler.CancelEditPlayer)

	api.POST("/players", handler.InsertPlayer)
	api.PUT("/players/:id", handler.UpdatePlayer)
	api.DELETE("/players/:id", handler.DeletePlayer)
}
