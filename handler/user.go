package handler

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"mq/academy/ent"
	"mq/academy/repo"
	"net/http"
)

type UsersRequest struct {
	Passwd string `json:"passwd"`
	Email  string `json:"email"`
}

type UsersResponse struct {
	User    *ent.User        `json:"user"`
	Account *ent.UserAccount `json:"account"`
	Error   string           `json:"error"`
}

func Users(gr repo.GameRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := &UsersResponse{}
		req := &UsersRequest{}
		ctx := repo.NewContext(r.Context(), gr)
		var err error
		vars := mux.Vars(r)

		if r.Method == http.MethodGet {
			res, err = getUser(ctx, vars["name"])
			if err != nil {
				log.Printf("getUser got error: %v", err)
			}
		} else if r.Method == http.MethodPost {
			if err = bindJSON(r, req); err != nil {
				log.Printf("bindJSON got error: %v", err)
			} else {
				res, err = createUser(ctx, vars["name"], req)
				if err != nil {
					log.Printf("createUser got error: %v", err)
				}
			}
		}

		if res == nil {
			res = &UsersResponse{Error: err.Error()}
		}
		fmt.Fprint(w, marshalJson(res))
	}
}

func getUser(ctx context.Context, name string) (*UsersResponse, error) {
	gr := repo.FromContext(ctx)
	user, err := gr.SelectUserByName(ctx, name)
	if err != nil {
		return nil, err
	}
	acc, err := gr.UserAccount(ctx, user)
	if err != nil {
		return nil, err
	}
	return &UsersResponse{
		User:    user,
		Account: acc,
	}, nil
}

func createUser(ctx context.Context, name string, req *UsersRequest) (*UsersResponse, error) {
	gr := repo.FromContext(ctx)

	res := &UsersResponse{}
	var err error
	if err := gr.WithTx(ctx, func(gr repo.GameRepo) error {
		res.User, err = gr.CreateUser(ctx, name)
		if err != nil {
			return err
		}
		res.Account, err = gr.CreateUserAccount(ctx, res.User, req.Passwd, req.Email)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return res, nil
}
