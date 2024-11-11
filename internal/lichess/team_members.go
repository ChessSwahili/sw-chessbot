package lichess

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const team_members_url = "https://lichess.org/api/team/nyumbani-mates/users"

type Member struct {
	ID       string `json:"id"`
	Username string `json:"name"`
}

type InsertMemberParams struct {
	LichessID string `json:"lichess_id"`
	Username  string `json:"username"`
}

func FetchTeamMembers() []InsertMemberParams {

	var dt []InsertMemberParams
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", team_members_url, nil)

	if err != nil {
		slog.Error("failed to create request", "error", err)
		return dt
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("LICHESS_TOKEN")))

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("failed to fetch team members", "err", err)
		return dt
	}

	defer resp.Body.Close()

	results := json.NewDecoder(resp.Body)

	for {

		var member Member

		err := results.Decode(&member)

		if err != nil {
			if err != io.EOF {
				slog.Error("we got an error while reading", "err", err)
			}

			break
		}

		dt = append(dt, InsertMemberParams{
			LichessID: member.ID,
			Username:  member.Username,
		})

	}

	// add by force lichess username & Id /;
	dt = append(dt, InsertMemberParams{
		LichessID: "herald18",
		Username:  "herald18",
	})

	return dt

}
