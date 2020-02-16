package cmd

import (
	"fmt"
	"strings"

	"github.com/alexellis/k3sup/cmd/apps"
	"github.com/spf13/cobra"
)

func MakeApps() *cobra.Command {
	var command = &cobra.Command{
		Use:   "app",
		Short: "Manage Kubernetes apps",
		Long:  `Manage Kubernetes apps`,
		Example: `  k3sup app install
  k3sup app info`,
		SilenceUsage: false,
	}

	var install = &cobra.Command{
		Use:   "install",
		Short: "Install a Kubernetes app",
		Example: `  k3sup app install [APP]
  k3sup app install openfaas --help
  k3sup app install inlets-operator --token-file $HOME/do
  k3sup app install --help`,
		SilenceUsage: true,
	}

	install.PersistentFlags().String("kubeconfig", "kubeconfig", "Local path for your kubeconfig file")

	install.RunE = func(command *cobra.Command, args []string) error {

		if len(args) == 0 {
			fmt.Printf("You can install: %s\n%s\n\n", strings.TrimRight("\n - "+strings.Join(getApps(), "\n - "), "\n - "),
				`Run k3sup app install NAME --help to see configuration options.`)
			return nil
		}

		return nil
	}

	info := &cobra.Command{
		Use:   "info",
		Short: "Find info about a Kubernetes app",
		Long:  "Find info about how to use the installed Kubernetes app",
		Example: `  k3sup app info [APP]
  k3sup app info openfaas
  k3sup app info inlets-operator
  k3sup app info
  k3sup app info --help`,
		SilenceUsage: true,
	}

	info.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			fmt.Println("You can get info about: openfaas, inlets-operator")
			return nil
		}

		if len(args) != 1 {
			return fmt.Errorf("you can only get info about exactly one installed app")
		}

		appName := args[0]

		switch appName {
		case "inlets-operator":
			fmt.Println(apps.InletsOperatorInfoMsg)
		case "openfaas":
			fmt.Println(apps.OpenFaaSInfoMsg)
		default:
			return fmt.Errorf("no info or no app available for %s", appName)
		}

		return nil
	}

	command.AddCommand(install)
	install.AddCommand(apps.MakeInstallOpenFaaS())
	install.AddCommand(apps.MakeInstallMetricsServer())
	install.AddCommand(apps.MakeInstallInletsOperator())
	install.AddCommand(apps.MakeInstallCertManager())
	install.AddCommand(apps.MakeInstallOpenFaaSIngress())
	install.AddCommand(apps.MakeInstallNginx())
	install.AddCommand(apps.MakeInstallChart())
	install.AddCommand(apps.MakeInstallTiller())
	install.AddCommand(apps.MakeInstallLinkerd())
	install.AddCommand(apps.MakeInstallCronConnector())
	install.AddCommand(apps.MakeInstallKafkaConnector())
	install.AddCommand(apps.MakeInstallMinio())
	install.AddCommand(apps.MakeInstallPostgresql())
	install.AddCommand(apps.MakeInstallKubernetesDashboard())
	install.AddCommand(apps.MakeInstallIstio())
	install.AddCommand(apps.MakeInstallTraefik2())

	command.AddCommand(info)

	return command
}

func getApps() []string {
	return []string{"openfaas",
		"nginx-ingress",
		"cert-manager",
		"openfaas-ingress",
		"inlets-operator",
		"metrics-server",
		"chart",
		"tiller",
		"linkerd",
		"cron-connector",
		"kafka-connector",
		"minio",
		"postgresql",
		"kubernetes-dashboard",
		"istio",
		"traefik2",
	}
}
