package scenarios

type Difficulty string

const (
	Easy    Difficulty = "easy"
	Medium  Difficulty = "medium"
	Hard    Difficulty = "hard"
	Extreme Difficulty = "extreme"
)

type Incident struct {
	ID          string
	Name        string
	Description string
	Difficulty  Difficulty
	TargetType  string // "docker", "kubernetes", "github-actions"
}

// GetAvailableIncidents returns a catalog of pre-defined failures that Dojo can inject
func GetAvailableIncidents() []Incident {
	return []Incident{
		{
			ID:          "k8s-oom",
			Name:        "OOMKilled Pod",
			Description: "Sets memory limits drastically low causing the pod to crash with OOMKilled",
			Difficulty:  Medium,
			TargetType:  "kubernetes",
		},
		{
			ID:          "docker-zombie",
			Name:        "Zombie Build Stage",
			Description: "Removes a crucial dependency in the final Docker image stage",
			Difficulty:  Hard,
			TargetType:  "docker",
		},
		{
			ID:          "k8s-typo",
			Name:        "Manifest Typo",
			Description: "Introduces a YAML syntax error or invalid field in a manifest",
			Difficulty:  Easy,
			TargetType:  "kubernetes",
		},
	}
}
