package apps

import (
	"fmt"
	"log"
	"os"

	execute "github.com/alexellis/go-execute/pkg/v1"
	"github.com/alexellis/k3sup/pkg/config"
	"github.com/alexellis/k3sup/pkg/env"
	"github.com/spf13/cobra"
	"gopkg.in/src-d/go-git.v4"
)

func getTraefik2Repo(repoPath, path string) error {

	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:      repoPath,
		Progress: os.Stdout,
	})
	if err != nil {
		return err
	}
	return nil
}

func installTraefik2(parts ...string) (execute.ExecResult, error) {

	task := execute.ExecTask{
		Command:     "helm",
		Args:        parts,
		StreamStdio: true,
	}
	res, err := task.Execute()
	if err != nil {
		return res, err
	}
	if res.ExitCode != 0 {
		return res, fmt.Errorf("exit code %d, stderr: %s", res.ExitCode, res.Stderr)
	}
	return res, nil
}

func MakeInstallTraefik2() *cobra.Command {
	var traefik2 = &cobra.Command{
		Use:          "traefik2",
		Short:        "Install traefik2",
		Long:         "Install traefik2",
		Example:      `  k3sup app install traefik2`,
		SilenceUsage: true,
	}
	fmt.Printf(traefikstarted)
	traefik2.RunE = func(command *cobra.Command, args []string) error {
		kubeConfigPath := getDefaultKubeconfig()

		if command.Flags().Changed("kubeconfig") {
			kubeConfigPath, _ = command.Flags().GetString("kubeconfig")
		}
		fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)
		//arch := getNodeArchitecture()
		//fmt.Printf("Node architecture: %q\n", arch)

		userPath, err := config.InitUserDir()
		if err != nil {
			return err
		}
		_, clientOS := env.GetClientArch()
		fmt.Printf("Client: %q\n", clientOS)

		log.Printf("User dir established as: %s\n", userPath)
		traefikPath := userPath + "traefik-helm-chart"

		if _, err = os.Stat(traefikPath); !os.IsNotExist(err) {
			os.RemoveAll(traefikPath)

		}
		fmt.Printf("Downloading Repository in: %s\n", traefikPath)
		err = getTraefik2Repo("https://github.com/containous/traefik-helm-chart.git", traefikPath)

		if err != nil {
			return err
		}
		fmt.Printf("Installing Helm Chart\n")
		_, err = installTraefik2("install", traefikPath+"/traefik")
		if err != nil {
			return err

		}
		os.Remove(traefikPath)
		return nil
	}

	return traefik2
}

const traefikstarted = `
_______              __ _ _           ___  
|__   __|            / _(_) |         |__ \ 
   | |_ __ __ _  ___| |_ _| | __ __   __ ) |
   | | '__/ _´ |/ _ \  _| | |/ / \ \ / // / 
   | | | | (_| |  __/ | | |   <   \ V // /_ 
   |_|_|  \__,_|\___|_| |_|_|\_\   \_/|____|
`
