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

type FriendsRequest struct {
	Names []string `json:"names"`
}

type FriendsResponse struct {
	Users []*ent.User `json:"users"`
	Error string      `json:"error"`
}

func Friends(gr repo.GameRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := &FriendsResponse{}
		req := &FriendsRequest{}
		ctx := repo.NewContext(r.Context(), gr)
		var err error
		vars := mux.Vars(r)

		if r.Method == http.MethodGet {
			res, err = getFriends(ctx, vars["name"])
			if err != nil {
				log.Printf("getFriends got error: %v", err)
			}
		} else if r.Method == http.MethodPost {
			if err = bindJSON(r, req); err != nil {
				log.Printf("bindJSON got error: %v", err)
			} else {
				res, err = addFriends(ctx, vars["name"], req)
				if err != nil {
					log.Printf("createUser got error: %v", err)
				}
			}
		}

		if res == nil {
			res = &FriendsResponse{Error: err.Error()}
		}
		fmt.Fprint(w, marshalJson(res))
	}
}

func getFriends(ctx context.Context, name string) (*FriendsResponse, error) {
	gr := repo.FromContext(ctx)
	user, err := gr.SelectUserByName(ctx, name)
	if err != nil {
		return nil, err
	}
	friends, err := gr.UserFriends(ctx, user)
	if err != nil {
		return nil, err
	}
	return &FriendsResponse{
		Users: friends,
	}, nil
}

func addFriends(ctx context.Context, name string, req *FriendsRequest) (*FriendsResponse, error) {
	gr := repo.FromContext(ctx)
	user, err := gr.SelectUserByName(ctx, name)
	if err != nil {
		return nil, err
	}
	us, err := gr.SelectUsersByName(ctx, req.Names...)
	if err != nil {
		return nil, err
	}
	if err := gr.AddFriends(ctx, user, us); err != nil {
		return nil, err
	}
	return &FriendsResponse{
		Users: us,
	}, nil
}
