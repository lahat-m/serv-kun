// Package scaffold generates the infrastructure files a Kuniel-shaped
// Spring Boot + Postgres project needs: docker-compose.prod.yml, nginx
// config, the Postgres backup script, a GitHub Actions deploy workflow, and
// a DEPLOY.md — so a new project gets the whole shape from `springboot
// init` instead of copy-pasting another project's files.
package scaffold

// File is one generated file, relative to the project root.
type File struct {
	Path    string
	Content string
	Exec    bool // true if the file should be made executable
}

// Generate renders every scaffolded file for the given params.
func Generate(p Params) ([]File, error) {
	compose, err := renderCompose(p)
	if err != nil {
		return nil, err
	}
	nginxConf, err := renderNginxConf(p)
	if err != nil {
		return nil, err
	}
	deployDoc, err := renderDeployDoc(p)
	if err != nil {
		return nil, err
	}

	return []File{
		{Path: "docker-compose.prod.yml", Content: compose},
		{Path: "nginx/templates/app.conf.template", Content: nginxConf},
		{Path: "nginx/conf.d/hardening.conf", Content: HardeningConf},
		{Path: "scripts/backup-postgres.sh", Content: BackupScript, Exec: true},
		{Path: ".github/workflows/deploy.yml", Content: DeployWorkflow},
		{Path: "DEPLOY.md", Content: deployDoc},
	}, nil
}
