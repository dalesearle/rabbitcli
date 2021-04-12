package data

import "strings"

/*
The connection provided name convention is app-some_unique_identifier.  For connection reports we currently only
want to look at the taxes app, which conveniently ar labeled as taxes-some_unique_identifier, amazing.  Lets white
list connections we are interested in, the white list may be ignored by the user by using -a/--all flag
 */

var whiteListedApps = map[string]struct{}{
	"taxes": {},
}

type Connection struct {
	ConnectedTime int64 `json:"connected_at"`
	Node string `json:"node"`
	ProvidedName string `json:"user_provided_name"`
	RabbitName string `json:"Name"`
}


func (c *Connection) IsWhiteListed() bool {
	parts := strings.Split(c.ProvidedName, "-")
	if len(parts) > 1 {
		_,found := whiteListedApps[parts[0]]
		return found
	}
	return false
}

func (c *Connection) ShortProvidedName() string {
	parts := strings.Split(c.ProvidedName, "-")
	if len(parts) > 0 {
		return parts[1]
	}
	return c.ProvidedName
}

func (c *Connection) ShortNodeName() string {
	return strings.Split(c.Node, ".")[0]
}

func (c *Connection) String() string {
	return c.RabbitName + ":" + c.ProvidedName
}

/*
Connection JSON (3/22/2021)
{
	"auth_mechanism": "PLAIN",
	"channel_max": 2047,
	"channels": 1,
	"client_properties": {
		"capabilities": {
			"authentication_failure_close": true,
			"basic.nack": true,
			"connection.blocked": true,
			"consumer_cancel_notify": true,
			"exchange_exchange_bindings": true,
			"publisher_confirms": true
		},
		"connection_name": "csqs-1616168344572",
		"copyright": "Copyright (c) 2007-2020 VMware, Inc. or its affiliates.",
		"information": "Licensed under the MPL. See https://www.rabbitmq.com/",
		"platform": "Java",
		"product": "RabbitMQ",
		"version": "5.9.0"
	},
	"connected_at": 1616168347011,
	"frame_max": 131072,
	"garbage_collection": {
		"fullsweep_after": 65535,
		"max_heap_size": 0,
		"min_bin_vheap_size": 46422,
		"min_heap_size": 233,
		"minor_gcs": 763
	},
	"host": "10.42.3.77",
	"name": "10.42.5.235:46812 -> 10.42.3.77:5671",
	"node": "rabbit@rabbitmq-1.rabbitmq-headless.csqs.svc.cluster.local",
	"peer_cert_issuer": null,
	"peer_cert_subject": null,
	"peer_cert_validity": null,
	"peer_host": "10.42.5.235",
	"peer_port": 46812,
	"port": 5671,
	"protocol": "AMQP 0-9-1",
	"recv_cnt": 62040,
	"recv_oct": 18206434,
	"recv_oct_details": {
		"rate": 578.0
	},
	"reductions": 15918583,
	"reductions_details": {
		"rate": 561.4
	},
	"send_cnt": 59050,
	"send_oct": 18327961,
	"send_oct_details": {
		"rate": 549.8
	},
	"send_pend": 0,
	"ssl": true,
	"ssl_cipher": "aes_256_cbc",
	"ssl_hash": "sha384",
	"ssl_key_exchange": "ecdhe_rsa",
	"ssl_protocol": "tlsv1.2",
	"state": "running",
	"timeout": 60,
	"type": "network",
	"user": "CSQSprod",
	"user_provided_name": "csqs-1616168344572",
	"user_who_performed_action": "CSQSprod",
	"vhost": "CSQS"
}
 */
