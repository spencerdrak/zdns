/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/zmap/zdns/pkg/zdns"
)

var cfgFile string
var gc zdns.GlobalConf

const envPrefix = "ZDNS"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "zdns",
	Short: "High-speed, low-drag DNS lookups",
	Long: `ZDNS is a library and CLI tool for making very fast DNS requests. It's built upon
https://github.com/zmap/dns (and in turn https://github.com/miekg/dns) for constructing
and parsing raw DNS packets. 

ZDNS also includes its own recursive resolution and a cache to further optimize performance.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:

	Run: func(cmd *cobra.Command, args []string) {
		log.Info(viper.AllKeys())
		log.Info(args)
		log.Info(gc.AlexaFormat)
		log.Info(viper.GetBool("alexa"))
		log.Info(cmd.Flags().GetBool("alexa"))
		// if len(args) < 1 {
		// 	log.Fatal("No lookup module specified. Valid modules: ", zdns.ValidlookupsString(), ".")
		// }
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.zdns.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	rootCmd.PersistentFlags().IntVar(&gc.Threads, "threads", 1000, "number of lightweight go threads")
	rootCmd.PersistentFlags().IntVar(&gc.GoMaxProcs, "go-processes", 0, "number of OS processes (GOMAXPROCS)")
	rootCmd.PersistentFlags().StringVar(&gc.NamePrefix, "prefix", "", "name to be prepended to what's passed in (e.g., www.)")
	rootCmd.PersistentFlags().StringVar(&gc.NameOverride, "override-name", "", "name overrides all passed in names")
	rootCmd.PersistentFlags().BoolVar(&gc.AlexaFormat, "alexa", false, "is input file from Alexa Top Million download")
	rootCmd.PersistentFlags().BoolVar(&gc.MetadataFormat, "metadata-passthrough", false, "if input records have the form 'name,METADATA', METADATA will be propagated to the output")
	rootCmd.PersistentFlags().BoolVar(&gc.IterativeResolution, "iterative", false, "Perform own iteration instead of relying on recursive resolver")
	rootCmd.PersistentFlags().StringVar(&gc.InputFilePath, "input-file", "-", "names to read")
	rootCmd.PersistentFlags().StringVar(&gc.OutputFilePath, "output-file", "-", "where should JSON output be saved")
	rootCmd.PersistentFlags().StringVar(&gc.MetadataFilePath, "metadata-file", "", "where should JSON metadata be saved")
	rootCmd.PersistentFlags().StringVar(&gc.LogFilePath, "log-file", "", "where should JSON logs be saved")

	rootCmd.PersistentFlags().StringVar(&gc.ResultVerbosity, "result-verbosity", "normal", "Sets verbosity of each output record. Options: short, normal, long, trace")
	rootCmd.PersistentFlags().StringVar(&gc.IncludeInOutput, "include-fields", "", "Comma separated list of fields to additionally output beyond result verbosity. Options: class, protocol, ttl, resolver, flags")

	rootCmd.PersistentFlags().IntVar(&gc.Verbosity, "verbosity", 3, "log verbosity: 1 (lowest)--5 (highest)")
	rootCmd.PersistentFlags().IntVar(&gc.Retries, "retries", 1, "how many times should zdns retry query if timeout or temporary failure")
	rootCmd.PersistentFlags().IntVar(&gc.MaxDepth, "max-depth", 10, "how deep should we recurse when performing iterative lookups")
	rootCmd.PersistentFlags().IntVar(&gc.CacheSize, "cache-size", 10000, "how many items can be stored in internal recursive cache")
	rootCmd.PersistentFlags().BoolVar(&gc.TCPOnly, "tcp-only", false, "Only perform lookups over TCP")
	rootCmd.PersistentFlags().BoolVar(&gc.UDPOnly, "udp-only", false, "Only perform lookups over UDP")
	rootCmd.PersistentFlags().BoolVar(&gc.NameServerMode, "name-server-mode", false, "Treats input as nameservers to query with a static query rather than queries to send to a static name server")

	rootCmd.PersistentFlags().String("name-servers", "", "List of DNS servers to use. Can be passed as comma-delimited string or via @/path/to/file. If no port is specified, defaults to 53.")
	rootCmd.PersistentFlags().String("local-addr", "", "comma-delimited list of local addresses to use")
	rootCmd.PersistentFlags().String("local-interface", "", "local interface to use")
	rootCmd.PersistentFlags().String("conf-file", "/etc/resolv.conf", "config file for DNS servers")
	rootCmd.PersistentFlags().Int("timeout", 15, "timeout for resolving an individual name")
	rootCmd.PersistentFlags().Int("iteration-timeout", 4, "timeout for resolving a single iteration in an iterative query")
	rootCmd.PersistentFlags().String("class", "INET", "DNS class to query. Options: INET, CSNET, CHAOS, HESIOD, NONE, ANY. Default: INET.")
	rootCmd.PersistentFlags().Bool("nanoseconds", false, "Use nanosecond resolution timestamps")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".zdns" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".zdns")
	}

	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
	// Bind the current command's flags to viper
	bindFlags(rootCmd, viper.GetViper())
}

// Reference: https://github.com/carolynvs/stingoftheviper/blob/main/main.go
// For how to make cobra/viper sync up, and still use custom struct
// Bind each cobra flag to its associated viper configuration (config file and environment variable)
func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		log.Info(f.Name, ": ", f.Value)
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --alexa to ZDNS_ALEXA
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			v.BindEnv(f.Name, fmt.Sprintf("%s_%s", envPrefix, envVarSuffix))
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}