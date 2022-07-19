package cmd

import (
	"fmt"
	"strings"

	"github.com/krateoplatformops/krateo/internal/osutils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	appName    = "krateo"
	appSummary = "Run your Resources on Every Cloud"
	banner     = `┏┓             ┏┓
┃┃┏┓ ┏━┓ ┏━━┓ ┏┛┗┓ ┏━━┓ ┏━━┓
┃┗┛┛ ┃┏┛ ┃┏┓┃ ┗┓┏┛ ┃┃━┫ ┃┏┓┃
┃┏┓┓ ┃┃  ┃┏┓┃  ┃┗┓ ┃┃━┫ ┃┗┛┃
┗┛┗┛ ┗┛  ┗┛┗┛  ┗━┛ ┗━━┛ ┗━━┛`
	envPrefix = "KRATEO"
)

func KrateoCtl(ver, build string) *cobra.Command {
	cmd := &cobra.Command{
		DisableSuggestions:    true,
		DisableFlagsInUseLine: true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		Use:           fmt.Sprintf("%s <COMMAND>", appName),
		Short:         appSummary,
		Long:          fmt.Sprintf("%s\n%s\n", banner, appSummary),
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd)
		},
	}

	cmd.AddCommand(newCmdVersion(ver, build))
	cmd.AddCommand(newInitCmd())
	cmd.AddCommand(newUninstallCmd())

	return cmd
}

func initializeConfig(cmd *cobra.Command) error {
	dir, err := osutils.GetAppDir(appName)
	if err != nil {
		return err
	}

	v := viper.New()

	// Set the base name of the config file, without the file extension.
	v.SetConfigName(appName)

	// Set as many paths as you like where viper should look for the
	// config file. We are only looking in the current working directory.
	v.AddConfigPath(dir)

	// Attempt to read the config file, gracefully ignoring errors
	// caused by a config file not being found. Return an error
	// if we cannot parse the config file.
	if err := v.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	v.SetEnvPrefix(envPrefix)

	// Bind to environment variables
	// Works great for simple config names, but needs help for names
	// like --favorite-color which we fix in the bindFlags function
	v.AutomaticEnv()

	// Bind the current command's flags to viper
	bindFlags(cmd, v)

	return nil
}

// Bind each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			//nolint:errcheck
			v.BindEnv(f.Name, fmt.Sprintf("%s_%s", envPrefix, envVarSuffix))
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			//nolint:errcheck
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
