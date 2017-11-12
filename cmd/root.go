package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dfs",
	Short: "a file server, upload and download",
	Long:  `DFS is a file upload/download server. You can use it to many computer.`,
	//	Run: func(cmd *cobra.Command, args []string) { },
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show version for DFS",
	Long:  `show version for DFS`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v0.1.0")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("debug", "d", "false", "debug, default false")
	rootCmd.AddCommand(versionCmd)

	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	viper.Set("default.configfile", "/etc/dfs.json")
	viper.Set("default.server.host", "127.0.0.1")
	viper.Set("default.server.port", "8058")
	viper.Set("default.server.updir", "/data/dfs")
	viper.Set("default.client.host", "127.0.0.1")
	viper.Set("default.client.port", "8058")
	viper.Set("default.client.downdir", "/tmp/download")
}
