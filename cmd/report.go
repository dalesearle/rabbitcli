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
	"github.com/spf13/viper"
	"rabbitcli/data"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Create a report on a RabbitMQ cluster",
	Long: `
The report command generates the following RabbitMQ cluster information:

* The number of connections.
* The number of queues.
* The node leader of quorum queues (if -a/--all is flagged).
* The number of duplicate connections.
* The number of connections on each instance in the cluster

The report will be run for taxes connections and rpc queues by default, to include
all connections and queues use -a/--all.`,
	Run: func(cmd *cobra.Command, args []string) {
		reportDriver()
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
	reportCmd.PersistentFlags().BoolVarP(&includeAll, "all", "a", false, "Bypass white list and include all connections and queues")
}

func reportDriver() {
	processReport()
}

func processReport() {
	r := make([]*data.Report,0)
	table := newAsciiTable()
	table.SetTitle("Cluster Report")
	table.SetHeaders([]string{"Report", "Value"})
	r = append(r, data.NewReport("Cluster", data.NewWorkingCluster(viper.Get(data.WORKING_CLUSTER).(string)).Name))
	r = append(r, reports()...)
	sort.Sort(data.ByName(r))
	for _, rp := range r {
		table.AddRow(rp.TableRow())
	}
	fmt.Println(table.String())
}

func reports() []*data.Report {
	var c *data.Connection
	cpn := ""
	csqs := "Connections - CSQS"
	dupes := "Connections - Duplicates"
	mq := "Connections - Missing Queue"
	taxes := "Connections - TAXES"
	m := make(map[string]int)
	reports := make([]*data.Report, 0)

	queues, err := queues()
	cobra.CheckErr(err)
	connections, err := connections()
	cobra.CheckErr(err)
	incrementReport(m, "Connections", len(connections))
	incrementReport(m, "Queues", len(queues))
	incrementReport(m, mq, 0)
	incrementReport(m, dupes, 0)
	bases := make(map[string]*data.QueueReportBase)
	mapConnectionsByProvidedName(bases, connections)
	mapQueuesByName(bases, queues)

	for _,c := range connections {
		incrementReport(m, c.ShortNodeName(), 1)
	}
	for _, b := range bases {
		if !b.IsMissingQueue() && b.Queue.HasLeader() {
			q := b.Queue
			reports = append(reports, data.NewReport("Queues - Leader - "+q.Name, q.ShortLeaderName()))
		}
		if !b.IsMissingConnection() {
			l := len(b.Connections())
			if l > 1 {
				incrementReport(m, dupes, l - 1)
			}
			c = b.Connections()[0]
			cpn = c.ProvidedName
			if strings.HasPrefix(cpn, "taxes") {
				incrementReport(m, taxes, 1)
				if b.IsMissingQueue() {
					incrementReport(m, mq, 1)
				}
			} else if strings.HasPrefix(cpn, "csqs") {
				incrementReport(m, csqs, 1)
			}
		}
	}
	for k,v := range m {
		reports = append(reports, data.NewReport(k, strconv.Itoa(v)))
	}
	return reports
}

func incrementReport(m map[string]int, key string, v int) {
	c,found := m[key]
	if !found {
		m[key] = v
	} else {
		m[key] = c + v
	}
}
