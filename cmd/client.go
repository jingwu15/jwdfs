package cmd

import (
	client "dfs/client"
	util "dfs/lib"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	os "os"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "upload file to DFS",
	Long:  `upload file to DFS`,
}

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "upload file to DFS",
	Long:  `upload file to DFS`,
	Run: func(cmd *cobra.Command, args []string) {
		mergeViperClient()

		client.Upload(viper.Get("client.file").(string), viper.Get("client.filekey").(string), "http://127.0.0.1:8058/upload")
		fmt.Println("upload done")
	},
}

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "download file from DFS",
	Long:  "download file from DFS",
	Run: func(cmd *cobra.Command, args []string) {
		client.Download(viper.Get("client.filekey").(string), "http://127.0.0.1:8058/download")
		fmt.Println("download done")
	},
}

// downCmd represents the client command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "the info for file from DFS",
	Long:  "the info for file from DFS",
	Run: func(cmd *cobra.Command, args []string) {
		client.Info(viper.Get("client.filekey").(string), "http://127.0.0.1:8058/info")
		fmt.Println("info done")
	},
}

var filekey string

func init() {
	clientCmd.PersistentFlags().StringP("config-file", "", "/etc/dfs.json", "the file key")
	clientCmd.PersistentFlags().StringP("host", "H", "", "the host of server")
	clientCmd.PersistentFlags().StringP("port", "P", "", "the port of server")
	clientCmd.PersistentFlags().StringP("file-key", "", "", "the file key")
	upCmd.Flags().StringP("file", "f", "", "the file name, include the path")
	downCmd.Flags().StringP("down-dir", "", "", "the path that dowoload file")
	downCmd.Flags().StringP("file-name", "", "", "download file to the file name, absolute path")

	viper.BindPFlag("configfile", clientCmd.PersistentFlags().Lookup("config-file"))
	viper.BindPFlag("host", clientCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("port", clientCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("downdir", downCmd.Flags().Lookup("down-dir"))

	viper.BindPFlag("client.filekey", clientCmd.PersistentFlags().Lookup("file-key"))
	viper.BindPFlag("client.file", upCmd.Flags().Lookup("file"))
	viper.BindPFlag("client.filename", downCmd.Flags().Lookup("file-name"))
	//viper.SetDefault("configfile", "/etc/dfs.json")

	clientCmd.AddCommand(upCmd)
	clientCmd.AddCommand(downCmd)
	clientCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(clientCmd)
}

func mergeViperClient() {
	//加载配置文件
	configfile := viper.Get("configfile").(string)
	viper.SetConfigFile(configfile)
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("the config file", configfile, "not exists")
		os.Exit(1)
	}

	//如果命令行有配置参数，则使用，否则使用默认值
	if util.IsEmpty(viper.Get("host")) {
		if !viper.IsSet("client.host") {
			viper.Set("client.host", viper.Get("default.client.host").(string))
		}
	} else {
		viper.Set("client.host", viper.Get("host").(string))
	}
	if util.IsEmpty(viper.Get("port")) {
		if !viper.IsSet("client.port") {
			viper.Set("client.port", viper.Get("default.client.port").(string))
		}
	} else {
		viper.Set("client.port", viper.Get("port").(string))
	}
	fmt.Println(viper.Get("downdir").(string))
	if util.IsEmpty(viper.Get("downdir")) {
		if !viper.IsSet("client.downdir") {
			viper.Set("client.downdir", viper.Get("default.client.downdir").(string))
		}
	} else {
		viper.Set("client.downdir", viper.Get("downdir").(string))
	}

}
