package handler

import (
	"fmt"
	"log"
	"mq/academy/repo"
	"mq/academy/session"
	"net/http"
)

type LoginRequest struct {
	Name   string `json:"name"`
	Passwd string `json:"passwd"`
}

type LoginResponse struct {
	Session *session.Session `json:"session"`
	Error   string           `json:"error"`
}

func Login(gr repo.GameRepo, smgr session.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &LoginRequest{}
		res := &LoginResponse{}
		ctx := r.Context()
		errorHandler := func(format string, err error) {
			if err != nil {
				log.Printf(format, err)
				res.Error = err.Error()
			} else {
				log.Printf(format)
				res.Error = format
			}
			fmt.Fprintf(w, marshalJson(res))
		}

		if err := bindJSON(r, req); err != nil {
			errorHandler("bindJSON got error: %v", err)
			return
		}
		user, err := gr.SelectUserByName(ctx, req.Name)
		if err != nil {
			errorHandler("querying user got error: %v", err)
			return
		}
		acc, err := gr.UserAccount(ctx, user)
		if err != nil {
			errorHandler("querying userAccount gor error: %v", err)
			return
		}
		if acc.Passwd != req.Passwd {
			errorHandler("passwd not matched", nil)
			return
		}
		sess, err := smgr.Set(user)
		if err != nil {
			errorHandler("creating session gor error: %v", err)
		}

		res.Session = sess
		fmt.Fprintf(w, marshalJson(res))
	}
}
