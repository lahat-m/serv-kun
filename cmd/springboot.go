package cmd

import (
	"github.com/lahat-m/serv-kun/internal/compose"
	"github.com/spf13/cobra"
)

var springbootCmd = &cobra.Command{
	Use:   "springboot",
	Short: "Set up and deploy a Spring Boot app shaped like Kuniel's docker-compose.prod.yml",
	Long: `springboot targets projects that follow Kuniel's docker-compose.prod.yml
layout: a compose file with "app", "postgres", "nginx", and "certbot"
services, driven by a .env.prod file. Run these from the project directory
(the one containing docker-compose.prod.yml), or pass --dir.`,
	Args: cobra.NoArgs, // reject mistyped subcommands — see showHelp in root.go
	RunE: showHelp,
}

func init() {
	springbootCmd.PersistentFlags().String("dir", ".", "project directory (containing the compose file)")
	springbootCmd.PersistentFlags().String("compose-file", "docker-compose.prod.yml", "compose file name")
	springbootCmd.PersistentFlags().String("env-file", ".env.prod", "env file name")
}

// projectFromFlags builds a compose.Project from the springboot command's
// persistent --dir/--compose-file/--env-file flags.
func projectFromFlags(cmd *cobra.Command) compose.Project {
	dir, _ := cmd.Flags().GetString("dir")
	file, _ := cmd.Flags().GetString("compose-file")
	envFile, _ := cmd.Flags().GetString("env-file")
	return compose.Project{Dir: dir, File: file, EnvFile: envFile}
}
