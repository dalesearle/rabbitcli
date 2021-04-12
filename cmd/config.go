/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"rabbitcli/data"
	"strings"
)

// cfgCmd represents the cfg command
var (
	appendCluster  bool
	delete         string
	initialize     bool
	list           bool
	passport       bool
	workingCluster string
	showWorking    bool

	configCmd = &cobra.Command{
		Use:   "config",
		Short: "Manage rabbitcli configuration",
		Long: `
The config command is used to manage cluster connection information.  The 
config must be set up before you can use the rabbitcli tool 
(rabbitcli config -i). The configuration file is called .rabbitcli.yaml
and is located in the users HOME directory. A cluster is configured with a
name, host, password, protocol (http/https) and user used to connect to the 
given host. Use a clusters name to set eh working cluster.
An example cluster configuration:

lasvegas: 
  host: lvrabbit.taxhawk.lv
  password: encrypted
  protocol: https
  user: encrypted
`,
		Run: func(cmd *cobra.Command, args []string) {
			configDriver()
		},
	}
)

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().BoolVarP(&appendCluster, "append", "n", false, "Add a cluster")
	configCmd.Flags().StringVarP(&delete, "delete", "d", "", "Delete a cluster")
	configCmd.Flags().BoolVarP(&initialize, "initialize", "i", false, "Initialize configuration")
	configCmd.Flags().BoolVarP(&list, "list", "l", false, "List configured clusters")
	configCmd.Flags().BoolVarP(&passport, "passport", "p", false, "Set user name and password for a cluster")
	configCmd.Flags().StringVarP(&workingCluster, "set", "s", "", "Set the working cluster")
	configCmd.Flags().BoolVarP(&showWorking, "working", "w", false, "Show working cluster")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cfgCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cfgCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func configDriver() {
	if appendCluster {
		processAppend()
	} else if delete != "" {
		processDelete()
	} else if initialize {
		processInitialize()
	} else if list {
		processList()
	} else if passport {
		processPassport()
	} else if workingCluster != "" {
		processNewWorkingCluster()
	} else if showWorking {
		processShowWorking()
	}
}

func processInitialize() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	cfgFile := home + "/.rabbitcli.yaml"
	_, err = os.OpenFile(cfgFile, os.O_RDONLY|os.O_EXCL|os.O_CREATE, 0644)
	if os.IsExist(err) {
		if !promptUserYN("Config file already exists, do you want to delete it and start over") {
			os.Exit(0)
		}
	}
	resetViper()
	if promptUserYN("Set configuration defaults for localhost, lasvegas and test") {
		lh := data.NewCluster("localhost")
		lh.Host = "localhost:15672"
		lh.Protocol = "http"
		lh.WriteToViper()
		lv := data.NewCluster("lasvegas")
		lv.Host = "lvrabbit.taxhawk.lv"
		lv.Protocol = "https"
		lv.WriteToViper()
		t := data.NewCluster("test")
		t.Host = "prvrabbit.taxhawk.prv"
		t.Protocol = "https"
		t.WriteToViper()
		lhu, lhp := createPassport(lh.Name)
		lh.SetPassword(lhp)
		lh.SetUserName(lhu)
		lh.WriteToViper()
		lvu, lvp := createPassport(lv.Name)
		lv.SetPassword(lvp)
		lv.SetUserName(lvu)
		lv.WriteToViper()
		tu, tp := createPassport(t.Name)
		t.SetPassword(tp)
		t.SetUserName(tu)
		t.WriteToViper()
		workingCluster := promptUser("Which of these would you like to set as your working cluster (localhost, lasvegas, test)? ")
		viper.Set(data.WORKING_CLUSTER, workingCluster)
	}
	viper.WriteConfig()
}

func resetViper() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	viper.Reset()
	viper.AddConfigPath(home)
	viper.SetConfigName(".rabbitcli")
	viper.SetConfigType("yaml")
}

func processAppend() {
	cc := allCurrentClusters()
	name := promptUser("Enter the cluster name: ")
	host := promptUser("Enter the host name: ")
	protocol := "http"
	if promptUserYN("Does this connect using https? ") {
		protocol = "https"
	}
	u, p := createPassport(name)
	nc := data.NewCluster(name)
	nc.Host = host
	nc.Protocol = protocol
	nc.WriteToViper()
	nc.SetPassword(p)
	nc.SetUserName(u)
	nc.WriteToViper()
	for _, cl := range cc {
		cl.SetUserName(cl.UserName())
		cl.SetPassword(cl.Password())
		cl.WriteToViper()
	}
	viper.WriteConfig()
}

func allCurrentClusters() []*data.Cluster {
	clusters := make([]*data.Cluster, 0)
	for _, k := range viper.AllKeys() {
		if strings.Contains(k, data.PROTOCOL) {
			cl := data.NewWorkingCluster(k)
			cl.Password()
			cl.UserName()
			clusters = append(clusters, cl)
		}
	}
	return clusters
}

func processList() {
	for _, k := range viper.AllKeys() {
		if strings.HasSuffix(k, "."+data.PROTOCOL) {
			c := data.NewWorkingCluster(k)
			fmt.Println(c.String())
		}
	}
}

func processNewWorkingCluster() {
	if viper.InConfig(workingCluster) {
		viper.Set(data.WORKING_CLUSTER, workingCluster)
		viper.WriteConfig()
	} else {
		fmt.Println("Cluster " + workingCluster + " does not exist")
	}
}

func processShowWorking() {
	c := data.NewWorkingCluster(viper.GetString(data.WORKING_CLUSTER))
	fmt.Println(c.String())
}

func processPassport() {
	name := promptUser("Resetting credentials, enter the cluster name: ")
	verifyClusterName(name)
	c := data.NewWorkingCluster(name)
	u, p := createPassport(name)
	c.SetPassword(p)
	c.SetUserName(u)
	c.WriteToViper()
	viper.WriteConfig()
}

func createPassport(id string) ([]byte, []byte) {
	u := promptUserSecret("Enter your " + id + " cluster user name: ")
	p := promptUserSecret("Enter your " + id + " cluster password: ")
	return u, p
}

func verifyClusterName(name string) {
	if viper.Get(name+"."+data.PROTOCOL) == nil {
		cobra.CheckErr(errors.New("cluster with name " + name + " does not exist"))
	}
}

func processDelete() {
	clusters := allCurrentClusters()
	resetViper()
	for _,c := range clusters {
		if c.Name != delete {
			c.WriteToViper()
		}
 	}
	for _,c := range clusters {
		if c.Name != delete {
			c.SetPassword(c.Password())
			c.SetUserName(c.UserName())
			c.WriteToViper()
		}
	}
	viper.WriteConfig()
}
