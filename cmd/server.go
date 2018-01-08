package cmd

import (
	"fmt"
	strings "strings"

	util "jwdfs/lib"
	server "jwdfs/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "the JW-DFS Server",
	Long:  `the JW-DFS Server`,
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start server",
	Long:  `start Server`,
	Run: func(cmd *cobra.Command, args []string) {
		mergeViperServer()

		server.Start()
		fmt.Println("started")
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop Server",
	Long:  `stop Server`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("stoped")
	},
}

func init() {
	serverCmd.PersistentFlags().StringP("config-file", "", "/etc/jwdfs.json", "the config file")
	startCmd.Flags().StringP("host", "H", "", "the host ip, default 127.0.0.1")
	startCmd.Flags().StringP("port", "P", "", "the port, default 8058")
	startCmd.Flags().StringP("up-dir", "", "", "the base path to store file for server, default /tmp/up")
	viper.BindPFlag("configfile", serverCmd.PersistentFlags().Lookup("config-file"))
	viper.BindPFlag("host", startCmd.Flags().Lookup("host"))
	viper.BindPFlag("port", startCmd.Flags().Lookup("port"))
	viper.BindPFlag("updir", startCmd.Flags().Lookup("up-dir"))
	//viper.SetDefault("server.host", "127.0.0.1")

	serverCmd.AddCommand(startCmd)
	serverCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(serverCmd)
}

func mergeViperServer() {
	//加载配置文件
	configfile := viper.Get("configfile").(string)
	viper.SetConfigFile(configfile)
	viper.ReadInConfig()

	//如果命令行有配置参数，则使用，否则使用默认值
	if util.IsEmpty(viper.Get("host")) {
		if !viper.IsSet("server.host") {
			viper.Set("server.host", viper.Get("default.server.host").(string))
		}
	} else {
		viper.Set("server.host", viper.Get("host").(string))
	}
	if util.IsEmpty(viper.Get("port")) {
		if !viper.IsSet("server.port") {
			viper.Set("server.port", viper.Get("default.server.port").(string))
		}
	} else {
		viper.Set("server.port", viper.Get("port").(string))
	}
	if util.IsEmpty(viper.Get("updir")) {
		if !viper.IsSet("server.updir") {
			viper.Set("server.updir", viper.Get("default.server.updir").(string))
		}
	} else {
		viper.Set("server.updir", viper.Get("updir").(string))
	}
	if !strings.HasSuffix(viper.Get("server.updir").(string), "/") {
		viper.Set("server.updir", viper.Get("server.updir").(string)+"/")
	}
}
