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
		// Kubernetes Scenarios
		{
			ID:          "k8s-oomkilled",
			Name:        "[K8s] OOMKilled Pod",
			Description: "Sets memory limits drastically low, causing the pod to crash with OOMKilled under load.",
			Difficulty:  Medium,
			TargetType:  "kubernetes",
		},
		{
			ID:          "k8s-imagepullbackoff",
			Name:        "[K8s] ImagePullBackOff",
			Description: "Simulates a typo in the container image tag or missing pull secrets.",
			Difficulty:  Easy,
			TargetType:  "kubernetes",
		},
		{
			ID:          "k8s-crashloop",
			Name:        "[K8s] CrashLoopBackOff",
			Description: "Removes a critical environment variable required by the application to start.",
			Difficulty:  Medium,
			TargetType:  "kubernetes",
		},
		{
			ID:          "k8s-service-selector",
			Name:        "[K8s] Unroutable Service",
			Description: "Changes the Service selector label to mismatch the Deployment pods (classic 502/504 error).",
			Difficulty:  Hard,
			TargetType:  "kubernetes",
		},
		{
			ID:          "k8s-liveness-timeout",
			Name:        "[K8s] Liveness Probe Failure",
			Description: "Sets the liveness probe timeout too low, causing Kubernetes to constantly restart healthy pods.",
			Difficulty:  Extreme,
			TargetType:  "kubernetes",
		},
		
		// Docker Scenarios
		{
			ID:          "docker-missing-deps",
			Name:        "[Docker] Zombie Runtime",
			Description: "Removes a crucial system dependency (like ca-certificates or libc) in the final image stage.",
			Difficulty:  Hard,
			TargetType:  "docker",
		},
		{
			ID:          "docker-permissions",
			Name:        "[Docker] Permission Denied",
			Description: "Changes the USER instruction, causing the entrypoint script to fail due to lack of file permissions.",
			Difficulty:  Medium,
			TargetType:  "docker",
		},
		
		// CI/CD Scenarios
		{
			ID:          "ci-missing-secret",
			Name:        "[CI/CD] Missing Deployment Secret",
			Description: "Alters the pipeline YAML to reference an undefined repository secret.",
			Difficulty:  Easy,
			TargetType:  "github-actions",
		},
		{
			ID:          "ci-test-flake",
			Name:        "[CI/CD] Flaky Tests Config",
			Description: "Removes cache configurations, causing pipeline times to skyrocket and tests to randomly timeout.",
			Difficulty:  Hard,
			TargetType:  "github-actions",
		},
		
		// Terraform Scenarios
		{
			ID:          "tf-syntax-error",
			Name:        "[Terraform] Invalid Syntax",
			Description: "Injects a subtle missing bracket or invalid HCL syntax.",
			Difficulty:  Easy,
			TargetType:  "terraform",
		},
		{
			ID:          "tf-missing-var",
			Name:        "[Terraform] Missing Required Variable",
			Description: "Removes a required variable definition, causing terraform plan to fail.",
			Difficulty:  Medium,
			TargetType:  "terraform",
		},
		
		// Infrastructure / Compose Scenarios
		{
			ID:          "nginx-bad-gateway",
			Name:        "[Nginx] 502 Bad Gateway",
			Description: "Modifies the nginx.conf proxy_pass directive to point to a non-existent upstream service.",
			Difficulty:  Medium,
			TargetType:  "docker-compose",
		},
		{
			ID:          "postgres-max-connections",
			Name:        "[PostgreSQL] Connection Refused",
			Description: "Lowers max_connections in Postgres to 1, causing the app to immediately fail to connect.",
			Difficulty:  Hard,
			TargetType:  "docker-compose",
		},
	}
}
