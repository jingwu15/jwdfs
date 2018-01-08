package cmd

import (
	client "jwdfs/client"
	util "jwdfs/lib"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	os "os"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "upload file to JW-DFS",
	Long:  `upload file to JW-DFS`,
}

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "upload file to JW-DFS",
	Long:  `upload file to JW-DFS`,
	Run: func(cmd *cobra.Command, args []string) {
		mergeViperClient()
		if util.IsEmpty(viper.Get("client.file")) {
			fmt.Println("please set file params, eg: --file /tmp/tester.jpg\n")
			clientCmd.Usage()
			os.Exit(1)
		}

		response := client.Upload(viper.Get("client.file").(string),
			viper.Get("client.filekey").(string),
			viper.Get("client.api").(string)+"/upload")
		fmt.Println(response)
	},
}

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "download file from JW-DFS",
	Long:  "download file from JW-DFS",
	Run: func(cmd *cobra.Command, args []string) {
		mergeViperClient()
		response := client.Download(viper.Get("client.filekey").(string),
			viper.Get("client.api").(string)+"/download")
		fmt.Println(response)
	},
}

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "the info for file from JW-DFS",
	Long:  "the info for file from JW-DFS",
	Run: func(cmd *cobra.Command, args []string) {
		mergeViperClient()
		response := client.Info(viper.Get("client.filekey").(string),
			viper.Get("client.api").(string)+"/info")
		fmt.Println(response)
	},
}

var filekey string

func init() {
	clientCmd.PersistentFlags().StringP("config-file", "", "/etc/jwdfs.json", "the file key")
	clientCmd.PersistentFlags().StringP("host", "", "", "the host of server")
	clientCmd.PersistentFlags().StringP("port", "", "", "the port of server")
	clientCmd.PersistentFlags().StringP("file-key", "", "", "the file key")
	upCmd.Flags().StringP("file", "", "", "the file name, include the path")
	downCmd.Flags().StringP("down-dir", "", "", "the path that dowoload file")
	downCmd.Flags().StringP("down-file", "", "", "download file to the file name, absolute path")

	viper.BindPFlag("configfile", clientCmd.PersistentFlags().Lookup("config-file"))
	viper.BindPFlag("host", clientCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("port", clientCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("downdir", downCmd.Flags().Lookup("down-dir"))

	viper.BindPFlag("client.filekey", clientCmd.PersistentFlags().Lookup("file-key"))
	viper.BindPFlag("client.file", upCmd.Flags().Lookup("file"))
	viper.BindPFlag("client.downfile", downCmd.Flags().Lookup("down-file"))
	//viper.SetDefault("configfile", "/etc/jwdfs.json")

	clientCmd.AddCommand(upCmd)
	clientCmd.AddCommand(downCmd)
	clientCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(clientCmd)
}

func mergeViperClient() {
	//加载配置文件
	configfile := viper.Get("configfile").(string)
	viper.SetConfigFile(configfile)
	viper.ReadInConfig()

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
	if util.IsEmpty(viper.Get("downdir")) {
		if !viper.IsSet("client.downdir") {
			viper.Set("client.downdir", viper.Get("default.client.downdir").(string))
		}
	} else {
		viper.Set("client.downdir", viper.Get("downdir").(string))
	}
	if util.IsEmpty(viper.Get("client.filekey")) {
		fmt.Println("please set file-key params, eg: --file-key tester/tester.jpg\n")
		clientCmd.Usage()
		os.Exit(1)
	}
	viper.Set("client.api", "http://"+viper.Get("client.host").(string)+":"+viper.Get("client.port").(string))
}
