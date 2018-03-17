package preve

type Source struct {
	Repo string `json:"repo" validate:"repository"`
	When string `json:"when" validate:"pr_action"`

	BaseURL string `json:"base_url" validate:"omitempty,url"` // for GitHub Enterprise
}
