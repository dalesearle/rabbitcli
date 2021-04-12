package data

import "strings"

var whiteListedQueues = map[string]struct{}{
	"rpc": {},
}

type Queue struct {
	Name string `json:"Name"`
	Leader string `json"leader"`
	VHost string `json"vhost"`
}

func (q *Queue) HasLeader() bool {
	return q.Leader != ""
}

func (q *Queue) IsWhiteListed() bool {
	parts := strings.Split(q.Name, "-")
	if len(parts) > 1 {
		_,found := whiteListedQueues[parts[0]]
		return found
	}
	return false
}

func (q *Queue) ShortName() string {
	parts := strings.Split(q.Name, "-")
	if len(parts) > 1 {
		return parts[1]
	}
	return q.Name
}

func (q *Queue) ShortLeaderName() string {
	return strings.Split(q.Leader, ".")[0]
}

func (q *Queue) String() string {
	return q.Name
}


/*
Queue JSON (3/22/2021)
{
	"arguments": {},
	"auto_delete": true,
	"backing_queue_status": {
		"avg_ack_egress_rate": 1.1744731253615061e-7,
		"avg_ack_ingress_rate": 1.1744731253615061e-7,
		"avg_egress_rate": 1.1744731253615061e-7,
		"avg_ingress_rate": 1.1744731253615061e-7,
		"delta": [
			"delta",
			"undefined",
			0,
			0,
			"undefined"
		],
		"len": 0,
		"mode": "default",
		"next_seq_id": 6003,
		"q1": 0,
		"q2": 0,
		"q3": 0,
		"q4": 0,
		"target_ram_count": "infinity"
	},
	"consumer_utilisation": null,
	"consumers": 1,
	"durable": false,
	"effective_policy_definition": {},
	"exclusive": true,
	"exclusive_consumer_tag": null,
	"garbage_collection": {
		"fullsweep_after": 65535,
		"max_heap_size": 0,
		"min_bin_vheap_size": 46422,
		"min_heap_size": 233,
		"minor_gcs": 23467
	},
	"head_message_timestamp": null,
	"idle_since": "2021-03-22 15:14:30",
	"memory": 13664,
	"message_bytes": 0,
	"message_bytes_paged_out": 0,
	"message_bytes_persistent": 0,
	"message_bytes_ram": 0,
	"message_bytes_ready": 0,
	"message_bytes_unacknowledged": 0,
	"message_stats": {
		"ack": 6003,
		"ack_details": {
			"rate": 0.0
		},
		"deliver": 6003,
		"deliver_details": {
			"rate": 0.0
		},
		"deliver_get": 6003,
		"deliver_get_details": {
			"rate": 0.0
		},
		"deliver_no_ack": 0,
		"deliver_no_ack_details": {
			"rate": 0.0
		},
		"get": 0,
		"get_details": {
			"rate": 0.0
		},
		"get_empty": 0,
		"get_empty_details": {
			"rate": 0.0
		},
		"get_no_ack": 0,
		"get_no_ack_details": {
			"rate": 0.0
		},
		"publish": 6003,
		"publish_details": {
			"rate": 0.0
		},
		"redeliver": 0,
		"redeliver_details": {
			"rate": 0.0
		}
	},
	"messages": 0,
	"messages_details": {
		"rate": 0.0
	},
	"messages_paged_out": 0,
	"messages_persistent": 0,
	"messages_ram": 0,
	"messages_ready": 0,
	"messages_ready_details": {
		"rate": 0.0
	},
	"messages_ready_ram": 0,
	"messages_unacknowledged": 0,
	"messages_unacknowledged_details": {
		"rate": 0.0
	},
	"messages_unacknowledged_ram": 0,
	"Name": "rpc-lvAppD_2020App01",
	"node": "rabbit@rabbitmq-1.rabbitmq-headless.csqs.svc.cluster.local",
	"operator_policy": null,
	"policy": null,
	"recoverable_slaves": null,
	"reductions": 15584446,
	"reductions_details": {
		"rate": 0.0
	},
	"single_active_consumer_tag": null,
	"state": "running",
	"type": "classic",
	"vhost": "CSQS"
}
*/
