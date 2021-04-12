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
)

// queuesCmd represents the queues command
var (
	repair bool
	mismatch bool

	queuesCmd = &cobra.Command{
		Use:   "queues",
		Short: "View reports and manage queues for a RabbitMQ cluster",
		Long:  `
The queues command provides queue reports and management for the RabbitMQ cluster.
Running the queues command with no flags will generate a report displaying the queue
names and the vhost they are associated with. Be aware that by default you will only
see rpc application queues (e.g rpc-some_unique_identifier) unless you use the -a/--all
flag.`,
		Run: func(cmd *cobra.Command, args []string) {
			queuesDriver()
		},
	}
)

func init() {
	rootCmd.AddCommand(queuesCmd)
	queuesCmd.PersistentFlags().BoolVarP(&includeAll, "all", "a", false, "Bypass white list and include all connections and queues")
	queuesCmd.Flags().BoolVarP(&mismatch, "mismatch", "m", false, "Report connection/queue mismatch issues")
	queuesCmd.Flags().BoolVarP(&repair, "repair", "r", false, "Report connection/queue mismatch issues and attempt to repair them.")
}

func queuesDriver() {
	queues, err := queues()
	cobra.CheckErr(err)
	if mismatch || repair {
		processMismatchReport(queues)
	} else {
		printQueueReport(queues)
	}
}

func queues() ([]*data.Queue, error) {
	queues, err :=  http.NewRabbitClient().Queues()
	//queues, err :=  http.NewRabbitClient("guest", "guest").Queues()
	if err != nil {
		return nil, err
	}
	if !includeAll {
		whiteList := make([]*data.Queue, 0)
		for _,q := range queues {
			if q.IsWhiteListed() {
				whiteList = append(whiteList, q)
			}
		}
		return whiteList, nil
	}
	return queues, nil
}

func processMismatchReport(queues []*data.Queue) {
	connections, err := connections()
	cobra.CheckErr(err)
	bases := make(map[string]*data.QueueReportBase)
	mapConnectionsByProvidedName(bases, connections)
	mapQueuesByName(bases, queues)
	printMismatchReport(bases)
	repairMismatchedQueues(bases)
}

func mapConnectionsByProvidedName(bases map[string]*data.QueueReportBase, connections []*data.Connection) {
	name := ""
	for _, c := range connections {
		name = c.ShortProvidedName()
		base, found := bases[name]
		if !found {
			bases[name] = data.NewQueueReportBase(c)
		} else {
			base.AddConnection(c)
		}
	}
}

func mapQueuesByName(bases map[string]*data.QueueReportBase, queues []*data.Queue) {
	name := ""
	for _, q := range queues {
		name = q.ShortName()
		base, found := bases[name]
		if !found {
			bases[name] = &data.QueueReportBase{Queue: q}
		} else {
			base.Queue = q
		}
	}
}

func printMismatchReport(bases map[string]*data.QueueReportBase) {
	table := newAsciiTable()
	table.SetTitle("Queue/Connection Mismatch")
	table.SetHeaders([]string{"Type", "Name", "Count", "Issue"})
	for _, b := range bases {
		if b.IsMissingConnection() {
			table.AddRow([]string{"Queue", b.Queue.Name, "1", "No Connection"})
		} else if b.IsMissingQueue() {
			table.AddRow([]string{ "Connection", b.Connections()[0].ProvidedName, strconv.Itoa(len(b.Connections())), "No Queue"})
		}
	}
	fmt.Println(table.String())
}

func repairMismatchedQueues(bases map[string]*data.QueueReportBase) {
	if repair {
		cons := make([]*data.Connection, 0)
		for _, b := range bases {
			if b.IsMissingQueue() {
				cons = append(cons, b.Connections()...)
			}
		}
		closeConnections(cons)
	}
}

func printQueueReport(queues []*data.Queue) {
	table := newAsciiTable()
	table.SetTitle(strconv.Itoa(len(queues)) + " Queues")
	table.SetHeaders([]string{"Name", "VHost"})
	for _, q := range queues {
		table.AddRow([]string{q.Name, q.VHost})
	}
	fmt.Println(table.String())
}
