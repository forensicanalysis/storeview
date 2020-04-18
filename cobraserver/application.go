package cobraserver

import (
	"fmt"
	"io"
	"os"

	"github.com/markbates/pkger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Command struct {
	Name       string
	Route      string
	Method     string
	SetupFlags func(f *pflag.FlagSet)
	Handler    func(io.Writer, io.Reader, *pflag.FlagSet) error
}

func Application(name string, width, height int, staticPath pkger.Dir, uiRoot bool, commands ...*Command) *cobra.Command {
	var rootCmd *cobra.Command

	uiCmd := UICommand(name, width, height, staticPath, commands...)
	if uiRoot {
		rootCmd = uiCmd
		rootCmd.Use = name
	} else {
		rootCmd = &cobra.Command{Use: name}
		rootCmd.AddCommand(uiCmd)
	}

	rootCmd.AddCommand(Commandline(commands...)...)
	rootCmd.AddCommand(ServeCommand(staticPath, commands...))

	// parse config
	var cfgFile string
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/."+name+".yaml)")
	initConfig(name, cfgFile)

	return rootCmd
}

// initConfig reads in config file and ENV variables if set.
func initConfig(name, cfgFile string) {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".view" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName("." + name)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func setupFlags(command *Command, flagset *pflag.FlagSet, formatDefault string) {
	flagset.String("format", formatDefault, "Set output format")
	if command.SetupFlags != nil {
		command.SetupFlags(flagset)
	}
}
