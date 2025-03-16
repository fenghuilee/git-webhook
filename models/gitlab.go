package models

type GitLabWebhook struct {
	ObjectKind string `json:"object_kind"`
	Project    struct {
		Name string `json:"name"`
	} `json:"project"`
	Ref        string `json:"ref"`
	Repository struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"repository"`
	Commits []struct {
		ID       string   `json:"id"`
		Message  string   `json:"message"`
		Added    []string `json:"added"`
		Modified []string `json:"modified"`
		Removed  []string `json:"removed"`
	} `json:"commits"`
}
