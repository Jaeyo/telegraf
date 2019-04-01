package api

import (
	"context"
	"log"
	"os"

	"github.com/influxdata/telegraf/plugins/inputs/zmon/protodata"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type TelegrafConfigAPI interface {
	GetTelegrafConfig(int) (*protodata.TelegrafConfig, error)
}

type telegrafConfigAPI struct {
	zmonAccountID  string
	token          string
	serverEndpoint string
}

func NewTelegrafConfigAPI(zmonAccountID, token, serverEndpoint string) TelegrafConfigAPI {
	return &telegrafConfigAPI{
		zmonAccountID:  zmonAccountID,
		token:          token,
		serverEndpoint: serverEndpoint,
	}
}

func (api *telegrafConfigAPI) createClient() (*grpc.ClientConn, protodata.TelegrafConfigServerClient, error) {
	conn, err := createConnection(api.serverEndpoint)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create connection to server")
	}

	client := protodata.NewTelegrafConfigServerClient(conn)
	return conn, client, nil
}

func (api *telegrafConfigAPI) GetTelegrafConfig(telegrafConfigID int) (*protodata.TelegrafConfig, error) {
	log.Println("D! ZMON, Get Telegraf Config")

	conn, client, err := api.createClient()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create grpc client")
	}
	defer conn.Close()

	hostname, err := os.Hostname()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get hostname")
	}

	req := &protodata.GetTelegrafConfigOnAgentRequest{
		Auth: &protodata.TenantAuth{
			ZmonAccountID: api.zmonAccountID,
			Token:         api.token,
		},
		Hostname:         hostname,
		TelegrafConfigID: int32(telegrafConfigID),
	}
	resp, err := client.GetTelegrafConfigOnAgent(context.Background(), req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get response from grpc server")
	}
	return resp.TelegrafConfig, nil
}
