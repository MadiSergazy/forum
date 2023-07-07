package handler

import (
	"net/http"
)

func (h *Handler) LikePost(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
}
