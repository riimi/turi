package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"mq/academy/ent"
	"mq/academy/repo"
	"net/http"
)

type MaybeFriendsResponse struct {
	MaybeFriends []*ent.User `json:"maybeFriends"`
}

func MaybeFriends(gr repo.GameRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := &MaybeFriendsResponse{}
		//ctx := repo.NewContext(r.Context(), gr)
		ctx := r.Context()
		var err error
		vars := mux.Vars(r)

		user, err := gr.SelectUserByName(ctx, vars["name"])
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err), http.StatusInternalServerError)
			return
		}
		us, err := gr.UserMaybeFriends(ctx, user)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err), http.StatusInternalServerError)
			return
		}

		res.MaybeFriends = us
		fmt.Fprint(w, marshalJson(res))
	}
}
