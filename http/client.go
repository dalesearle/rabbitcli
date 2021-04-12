package http

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"rabbitcli/data"
)

type rabbitClient struct {
	cluster *data.Cluster
}

func NewRabbitClient() *rabbitClient {
	return &rabbitClient{
		cluster: data.NewWorkingCluster(viper.GetString(data.WORKING_CLUSTER)),
	}
}

func (c *rabbitClient) CloseConnection(connection *data.Connection) error {
	req, err := http.NewRequest("DELETE", c.cluster.Protocol + "://" + c.cluster.Host + "/api/connections/" + connection.RabbitName, nil)
	if err != nil {
		return err
	}
	req.Header.Add("X-Reason", "rabbitcli")
	c.doRequest(req, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *rabbitClient) Connections() ([]*data.Connection, error) {
	connectionData := make([]*data.Connection, 0)

	req, err := http.NewRequest("GET", c.cluster.Protocol + "://" + c.cluster.Host + "/api/connections", nil)
	if err != nil {
		return nil, err
	}
	err = c.doRequest(req, &connectionData)
	if err != nil {
		return nil, err
	}
	return connectionData, nil
}

func (c *rabbitClient) Queues() ([]*data.Queue, error) {
	queueData := make([]*data.Queue, 0)

	req, err := http.NewRequest("GET", c.cluster.Protocol + "://" + c.cluster.Host + "/api/queues", nil)
	if err != nil {
		return nil, err
	}
	err = c.doRequest(req, &queueData)
	if err != nil {
		return nil, err
	}
	return queueData, nil
}

func (c *rabbitClient) doRequest(req *http.Request, target interface{}) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	req.SetBasicAuth(string(c.cluster.UserName()), string(c.cluster.Password()))
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		if target != nil {
			err = json.Unmarshal(body, target)
			if err != nil {
				return err
			}
		}
	} else {
		apiError := &data.ApiError{}
		err = json.Unmarshal(body, apiError)
		if err == nil {
			err = errors.New(apiError.Reason)
		}
		return err
	}
	return nil
}
