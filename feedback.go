package main

import (
	"net/http"

	"github.com/JedBeom/wbot/model"
)

func SkillFeedback(w http.ResponseWriter, r *http.Request) {
	history, ok := r.Context().Value("history").(model.History)
	if !ok {
		w.WriteHeader(400)
		return
	}

	f := model.Feedback{
		HistoryID: history.ID,
		UserID:    history.UserID,
		Text:      history.Params["feedback"],
	}

	err := f.Create(db)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	writeOK(w, `{"version": "2.0"}`)
}
