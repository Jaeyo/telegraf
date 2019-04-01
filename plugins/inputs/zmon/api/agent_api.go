package api

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/influxdata/telegraf/plugins/inputs/zmon/protodata"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type AgentAPI interface {
	NotifyAgentUp(time.Duration) error
	NotifyAgentDown() error
	ChangeSyncStatus(string) error
}

type agentAPI struct {
	zmonAccountID  string
	token          string
	serverEndpoint string
}

func NewAgentAPI(zmonAccountID, token, serverEndpoint string) AgentAPI {
	return &agentAPI{
		zmonAccountID:  zmonAccountID,
		token:          token,
		serverEndpoint: serverEndpoint,
	}
}

func (api *agentAPI) createClient() (*grpc.ClientConn, protodata.AgentServerClient, error) {
	conn, err := createConnection(api.serverEndpoint)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create connection to server")
	}

	client := protodata.NewAgentServerClient(conn)
	return conn, client, nil
}

func (api *agentAPI) NotifyAgentUp(interval time.Duration) error {
	log.Println("D! ZMON, Notify Agent Up")

	conn, client, err := api.createClient()
	if err != nil {
		return errors.Wrap(err, "failed to create grpc client")
	}
	defer conn.Close()

	hostname, err := os.Hostname()
	if err != nil {
		return errors.Wrap(err, "failed to get hostname")
	}

	req := &protodata.NotifyAgentUpRequest{
		Auth: &protodata.TenantAuth{
			ZmonAccountID: api.zmonAccountID,
			Token:         api.token,
		},
		Hostname: hostname,
		Interval: int32(interval / time.Second),
	}
	_, err = client.NotifyAgentUp(context.Background(), req)
	return err
}

func (api *agentAPI) NotifyAgentDown() error {
	log.Println("D! ZMON, Notify Agent Down")

	conn, client, err := api.createClient()
	if err != nil {
		return errors.Wrap(err, "failed to create grpc client")
	}
	defer conn.Close()

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	req := &protodata.NotifyAgentDownRequest{
		Auth: &protodata.TenantAuth{
			ZmonAccountID: api.zmonAccountID,
			Token:         api.token,
		},
		Hostname: hostname,
	}
	_, err = client.NotifyAgentDown(context.Background(), req)
	return err
}

func (api *agentAPI) ChangeSyncStatus(status string) error {
	log.Println("D! ZMON, Change Sync Status")

	conn, client, err := api.createClient()
	if err != nil {
		return errors.Wrap(err, "failed to create grpc client")
	}
	defer conn.Close()

	hostname, err := os.Hostname()
	if err != nil {
		return errors.Wrap(err, "failed to get hostname")
	}

	req := &protodata.ChangeSyncStatusRequest{
		Auth: &protodata.TenantAuth{
			ZmonAccountID: api.zmonAccountID,
			Token:         api.token,
		},
		Hostname: hostname,
		Status:   status,
	}
	_, err = client.ChangeSyncStatus(context.Background(), req)
	return err
}
