package response

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/caarvid/armadan/web/template/partials"
)

type ResponseBuilder struct {
	res http.ResponseWriter
	req *http.Request
	buf *bytes.Buffer
}

func New(res http.ResponseWriter, req *http.Request, comp templ.Component) *ResponseBuilder {
	buf := &bytes.Buffer{}

	comp.Render(req.Context(), buf)

	return &ResponseBuilder{
		buf: buf,
		res: res,
		req: req,
	}
}

func (rb *ResponseBuilder) WithPartial(comp templ.Component) *ResponseBuilder {
	comp.Render(rb.req.Context(), rb.buf)

	return rb
}

func (rb *ResponseBuilder) WithSuccess(msg string) *ResponseBuilder {
	partials.SuccessToast(msg).Render(rb.req.Context(), rb.buf)

	return rb
}

func (rb *ResponseBuilder) WithError(msg string) *ResponseBuilder {
	partials.ErrorToast(msg).Render(rb.req.Context(), rb.buf)

	return rb
}

func (rb *ResponseBuilder) HTML() {
	rb.res.Write(rb.buf.Bytes())
}

func ResetPasswordEmailSent(w http.ResponseWriter, r *http.Request, email string) {
	w.Header().Add("HX-Retarget", "#forgot-password-message")
	w.Header().Add("HX-Reselect", "#forgot-password-message")
	w.WriteHeader(http.StatusOK)

	partials.ForgotPasswordMessage(fmt.Sprintf("En återställningslänk har skickats till %s", email)).Render(r.Context(), w)
}

func ResetPasswordMessage(w http.ResponseWriter, r *http.Request, msg, t string) {
	status := http.StatusOK

	if t == "error" {
		status = http.StatusUnprocessableEntity
	}

	w.Header().Add("HX-Retarget", "#reset-password-message")
	w.Header().Add("HX-Reselect", "#reset-password-message")
	w.WriteHeader(status)

	partials.ResetPasswordMessage(msg, t).Render(r.Context(), w)
}

func InvalidCredentialsError(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("HX-Retarget", "#login-error")
	w.Header().Add("HX-Reselect", "#login-error")
	w.WriteHeader(http.StatusUnprocessableEntity)

	partials.LoginError("Fel email eller lösenord").Render(r.Context(), w)
}

func LoginValidationError(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("HX-Retarget", "#login-error")
	w.Header().Add("HX-Reselect", "#login-error")
	w.WriteHeader(http.StatusUnprocessableEntity)

	partials.LoginError("Ogiltigt format").Render(r.Context(), w)
}

func GeneralLoginError(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("HX-Retarget", "#login-error")
	w.Header().Add("HX-Reselect", "#login-error")
	w.WriteHeader(http.StatusInternalServerError)

	partials.LoginError("Något gick fel, försök igen senare").Render(r.Context(), w)
}

func ErrorMessage(w http.ResponseWriter, r *http.Request, msg string) {
	w.WriteHeader(http.StatusInternalServerError)

	partials.ErrorToast(msg).Render(r.Context(), w)
}

func GeneralError(w http.ResponseWriter, r *http.Request) {
	ErrorMessage(w, r, "Något gick fel")
}
