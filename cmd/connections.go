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
	"fmt"
	"github.com/spf13/cobra"
	"rabbitcli/data"
	"rabbitcli/http"
	"strconv"
	"time"
)

var (
	duplicates bool
	close      bool
	report     bool
	// connectionsCmd represents the connections command
	connectionsCmd = &cobra.Command{
		Use:   "connections",
		Short: "View reports and manage connections for a RAbbitMQ cluster",
		Long: `
The connections command provides connection reports and management for a RabbitMQ cluster.
Running the connections command with no flags will generate a report displaying the connections 
provided name, RabbitMQ name, the instance connected to and the time it connected. Be aware that
by default you will only see taxes application connections (e.g taxes-some_unique_identifier)
unless you use the -a/--all flag which will create reports using all connections. The duplicate
command reports on the provided name (all RabbitMQ names are unique). If you add the -c/--close
command while using the -d/--duplicates flag, all of the duplicate connections will be closed
(rabbitcli connections -d -c). If you use the -c/--close flag only, you will be prompted to make 
sure you want to close all the connections to the set.`,
		Run: func(cmd *cobra.Command, args []string) {
			connectionsDriver()
		},
	}
)

func init() {
	rootCmd.AddCommand(connectionsCmd)
	connectionsCmd.PersistentFlags().BoolVarP(&includeAll, "all", "a", false, "Bypass white list and include all connections")
	connectionsCmd.Flags().BoolVarP(&close, "close", "c", false, "Close connections")
	connectionsCmd.Flags().BoolVarP(&duplicates, "duplicates", "d", false, "Report duplicate connections and counts")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// connectionsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// connectionsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func connectionsDriver() {
	connections, err := connections()
	cobra.CheckErr(err)
	if duplicates {
		processDuplicates(connections)
	} else if close {
		if promptUserYN("This will close all the connections in localhost, continue") {
			closeConnections(connections)
		}
	} else {
		printConnectionsReport(connections)
	}
}

func connections() ([]*data.Connection, error) {
	connections, err := http.NewRabbitClient().Connections()
	//connections,err := http.NewRabbitClient("guest", "guest").Connections()
	if err != nil {
		return nil, err
	}
	if !includeAll {
		whiteList := make([]*data.Connection, 0)
		for _, c := range connections {
			if c.IsWhiteListed() {
				whiteList = append(whiteList, c)
			}
		}
		return whiteList, nil
	}
	return connections, nil
}

func processDuplicates(connections []*data.Connection) {
	dups := findDuplicateConnections(connections)
	if close {
		dc := make([]*data.Connection, 0)
		for _, dup := range dups {
			dc = append(dc, dup.Connections...)
		}
		closeConnections(dc)
	} else {
		printDuplicateConnectionsReport(dups)
	}
}

func findDuplicateConnections(connections []*data.Connection) []*data.ConnectionReportBase {
	m := mapConnectionCountsByProvidedName(connections)
	return findDuplicateConnectionsByProviderName(m)
}

func duplicateConnectionCount(bases []*data.ConnectionReportBase) int {
	c := 0
	for _, b := range bases {
		c-- // hopefully one of the connections is not a dup
		c += b.ConnectionCount()
	}
	return c
}

func mapConnectionCountsByProvidedName(connections []*data.Connection) map[string]*data.ConnectionReportBase {
	m := make(map[string]*data.ConnectionReportBase)
	n := ""
	for _, c := range connections {
		n = c.ProvidedName
		b, found := m[n]
		if !found {
			m[n] = data.NewConnectionReportBase(c)
		} else {
			b.AddConnection(c)
		}
	}
	return m
}

func findDuplicateConnectionsByProviderName(bases map[string]*data.ConnectionReportBase) []*data.ConnectionReportBase {
	m := make([]*data.ConnectionReportBase, 0)
	for _, c := range bases {
		if c.ConnectionCount() > 1 {
			m = append(m, c)
		}
	}
	return m
}

func closeConnections(connections []*data.Connection) {
	table := newAsciiTable()
	table.SetTitle("Closed Connections")
	table.SetHeaders([]string{"Provided Name", "Rabbit Name", "Result"})
	for _, c := range connections {
		err := closeConnection(c)
		if err != nil {
			table.AddRow([]string{c.ProvidedName, c.RabbitName, err.Error()})
		} else {
			table.AddRow([]string{c.ProvidedName, c.RabbitName, "OK"})
		}
	}
	fmt.Println(table.String())
}

func closeConnection(connection *data.Connection) error {
	return http.NewRabbitClient().CloseConnection(connection)
	//return http.NewRabbitClient("guest", "guest").CloseConnection(connection)
}

func printDuplicateConnectionsReport(bases []*data.ConnectionReportBase) {
	dupes := 0
	table := newAsciiTable()
	table.SetTitle(strconv.Itoa(duplicateConnectionCount(bases)) + " Duplicate Connections")
	table.SetHeaders([]string{"Name", "Count"})
	for _, c := range bases {
		dupes += c.ConnectionCount()
		table.AddRow([]string{c.Name, strconv.Itoa(c.ConnectionCount())})
	}
	fmt.Println(table.String())
}

func printConnectionsReport(connections []*data.Connection) {
	table := newAsciiTable()
	table.SetTitle(strconv.Itoa(len(connections)) + " Connections")
	table.SetHeaders([]string{"Provided Name", "Rabbit Name", "Node", "Connected Time"})
	for _, c := range connections {
		table.AddRow([]string{c.ProvidedName, c.RabbitName, c.ShortNodeName(), time.Unix(c.ConnectedTime, 0).Format(time.RFC822)})
	}
	fmt.Println(table.String())
}
