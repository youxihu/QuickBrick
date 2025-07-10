package domain

type Config struct {
	Port        string     `yaml:"port"`
	SecretToken string     `yaml:"secret_token"`
	DB          DBConfig   `yaml:"database"`
	Pipelines   []Pipeline `yaml:"pipelines"`
}

type DBConfig struct {
	Driver   string `yaml:"driver"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Database string `yaml:"database"`
}

type Pipeline struct {
	Name      string `yaml:"name"`
	Type      string `yaml:"type"`       // frontend/backend
	Env       string `yaml:"env"`        // beta/prod
	EventType string `yaml:"event_type"` // push/tag
	Script    string `yaml:"script"`
}

// PushEvent 完整的 GitLab Push Hook 结构
type PushEvent struct {
	ObjectKind        string `json:"object_kind"`
	EventName         string `json:"event_name"` // "push"
	Ref               string `json:"ref"`        // branch or tag
	TotalCommitsCount int    `json:"total_commits_count"`

	UserEmail string `json:"user_email"`
	UserName  string `json:"user_name"`

	Project struct {
		Name string `json:"name"`
		URL  string `json:"web_url"`
	} `json:"project"`

	Commits []struct {
		ID      string `json:"id"`
		Message string `json:"message"`
		Author  struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"author"`
		Timestamp string `json:"timestamp"`
		URL       string `json:"url"`
	} `json:"commits"`
}
