package main

import (
	"auth-rbac/internal/config"
	"auth-rbac/internal/rbac"
	"auth-rbac/internal/session"
	"auth-rbac/internal/user"
	"auth-rbac/pkg/logger"
	"html/template"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Handler struct {
	logger    logger.Logger
	cfg       config.Config
	rbac      *rbac.RBAC
	User      user.Users
	Session   session.Sessions
	templates map[string]*template.Template
}

func newHandler(newLogger logger.Logger, config config.Config, rbac *rbac.RBAC, user user.Users, session session.Sessions, templates map[string]*template.Template) *Handler {
	return &Handler{
		logger:    newLogger,
		cfg:       config,
		rbac:      rbac,
		User:      user,
		Session:   session,
		templates: templates,
	}
}

func (h *Handler) routers(r *chi.Mux) *chi.Mux {
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/signup", h.signUpHelper)
		r.Post("/signup", h.signUp)
		r.Get("/signin", h.signInHelper)
		r.Post("/signin", h.signIn)
		r.Get("/foo", h.foo)
		r.Get("/bar", h.bar)
		r.Get("/sigma", h.sigma)
	})

	return r
}

func parseTemplates() map[string]*template.Template {
	var templates map[string]*template.Template

	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	templates["signin"] = template.Must(template.ParseFiles("./template/signin/index.html", "./template/signin/base.html"))
	templates["signup"] = template.Must(template.ParseFiles("./template/signup/index.html", "./template/signup/base.html"))

	return templates
}

func (h *Handler) renderTemplate(w io.Writer, name string, viewModel interface{}) {
	tmpl, ok := h.templates[name]
	if !ok {
		h.logger.Fatalf("can't find template")

		return
	}

	err := tmpl.ExecuteTemplate(w, "base", viewModel)
	if err != nil {
		h.logger.Fatalf("can not execute tamplate with template and viewmodel: %s: %s", viewModel, err)

		return
	}
}

func (h *Handler) foo(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("Authorization")
	if err != nil {
		http.Error(w, "login to get access", http.StatusBadRequest)

		return
	}

	token := cookie.Value

	userID, err := session.DecodeToken(token)
	if err != nil {
		http.Error(w, "can not decode token", http.StatusInternalServerError)

		return
	}

	ses, err := h.Session.FindByToken(token)
	if err != nil {
		http.Error(w, "can not find session", http.StatusBadRequest)

		return
	}

	if !session.CheckValidSes(token, ses) {
		http.Error(w, "session timeout", http.StatusUnauthorized)

		return
	}

	u, err := h.User.FindUserByID(userID)
	if err != nil {
		http.Error(w, "can not find user", http.StatusBadRequest)

		return
	}

	if err := h.rbac.Authorize(r, u.Role, "foo", "get"); err != nil {
		http.Error(w, "you have no permission", http.StatusForbidden)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) bar(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("Authorization")
	if err != nil {
		http.Error(w, "login to get access", http.StatusBadRequest)

		return
	}

	token := cookie.Value

	userID, err := session.DecodeToken(token)
	if err != nil {
		http.Error(w, "can not decode token", http.StatusInternalServerError)

		return
	}

	ses, err := h.Session.FindByToken(token)
	if err != nil {
		http.Error(w, "can not find session", http.StatusBadRequest)

		return
	}

	if !session.CheckValidSes(token, ses) {
		http.Error(w, "session timeout", http.StatusUnauthorized)

		return
	}

	u, err := h.User.FindUserByID(userID)
	if err != nil {
		http.Error(w, "can not find user", http.StatusBadRequest)

		return
	}

	if err := h.rbac.Authorize(r, u.Role, "bar", "get"); err != nil {
		http.Error(w, "you have no permission", http.StatusForbidden)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) sigma(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("Authorization")
	if err != nil {
		http.Error(w, "login to get access", http.StatusUnauthorized)

		return
	}

	token := cookie.Value

	userID, err := session.DecodeToken(token)
	if err != nil {
		http.Error(w, "can not decode token", http.StatusInternalServerError)

		return
	}

	ses, err := h.Session.FindByToken(token)
	if err != nil {
		http.Error(w, "can not find session", http.StatusBadRequest)

		return
	}

	if !session.CheckValidSes(token, ses) {
		http.Error(w, "session timeout", http.StatusUnauthorized)

		return
	}

	u, err := h.User.FindUserByID(userID)
	if err != nil {
		http.Error(w, "can not find user", http.StatusBadRequest)

		return
	}

	if err := h.rbac.Authorize(r, u.Role, "sigma", "get"); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) signUpHelper(w http.ResponseWriter, r *http.Request) {
	h.renderTemplate(w, "signup", "")
}

func (h *Handler) signInHelper(w http.ResponseWriter, r *http.Request) {
	h.renderTemplate(w, "signin", "")
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	pass := r.FormValue("pass")

	u, err := h.User.FindUserByLogin(login)
	if err != nil {
		http.Error(w, "can not find user", http.StatusNotFound)

		return
	}

	hp, err := user.HashPass(pass)
	if err != nil {
		h.logger.Debugf("failed to hash pass")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	if u.Password == hp {
		token, ses, err := session.CreateSes(u, h.cfg.Session.Life)
		if err != nil {
			h.logger.Errorf("failed to create ses", err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		err = h.Session.Upsert(ses)
		if err != nil {
			h.logger.Errorf("failed to create ses", err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		dur, err := time.ParseDuration(h.cfg.Cookies.Life)
		if err != nil {
			h.logger.Errorf("failed to parse duration", err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		cookie := http.Cookie{Name: "Authorization", Value: token, Expires: time.Now().Add(dur)}
		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "wrong password", http.StatusConflict)
	}
}

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	var u user.User

	u.Login = r.FormValue("login")
	u.Password = r.FormValue("pass")

	err := user.CheckValidUser(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	hashedPass, err := user.HashPass(u.Password)
	if err != nil {
		h.logger.Errorf("failed to hash pass", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	u.Password = hashedPass
	u.CreatedAt = time.Now()
	u.Role = "client"

	err = h.User.Create(&u)
	if err != nil {
		http.Error(w, "user already exist", http.StatusBadRequest)

		return
	}

	w.WriteHeader(http.StatusCreated)
}
