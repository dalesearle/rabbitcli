package cmd

import (
	"rabbitcli/data"
	"testing"
)

var connectionTests = []struct {
	connections []*data.Connection
	expected    int
}{
	{[]*data.Connection{}, 0},
	{[]*data.Connection{&data.Connection{RabbitName: "RN1", ProvidedName: "PN1"}}, 1},
	{[]*data.Connection{&data.Connection{RabbitName: "RN1", ProvidedName: "PN1"}, &data.Connection{RabbitName: "RN2", ProvidedName: "PN2"}, &data.Connection{RabbitName: "RN3", ProvidedName: "PN3"}}, 3},
	{[]*data.Connection{&data.Connection{RabbitName: "RN1", ProvidedName: "PN1"}, &data.Connection{RabbitName: "RN2", ProvidedName: "PN2"}, &data.Connection{RabbitName: "RN2", ProvidedName: "PN2"}}, 2},
}

var duplicateTests = []struct {
	bases    map[string]*data.ConnectionReportBase
	expected int
}{
	{map[string]*data.ConnectionReportBase{}, 0},
	{map[string]*data.ConnectionReportBase{"PN1": &data.ConnectionReportBase{Name: "PN1", Connections: []*data.Connection{&data.Connection{ProvidedName: "PN1", RabbitName: "RN1"}}}}, 0},
	{map[string]*data.ConnectionReportBase{"PN1": &data.ConnectionReportBase{Name: "PN1", Connections: []*data.Connection{&data.Connection{ProvidedName: "PN1", RabbitName: "RN1"},&data.Connection{ProvidedName: "PN1", RabbitName: "RN2"}}}}, 1},
}

func TestConnectionCountsByProvidedName(t *testing.T) {
	l := 0
	for i, test := range connectionTests {
		c := mapConnectionCountsByProvidedName(test.connections)
		l = len(c)
		if l != test.expected {
			t.Fatalf("test %d failed, expected %d got %d", i, test.expected, l)
		}
	}
}

func TestFindDuplicateConnections(t *testing.T) {
	var counters []*data.ConnectionReportBase
	var l int
	for i, test := range duplicateTests {
		counters = findDuplicateConnectionsByProviderName(test.bases)
		l = len(counters)
		if l != test.expected {
			t.Fatalf("test %d failed, expected %d got %d", i, test.expected, l)
		}
	}
}
