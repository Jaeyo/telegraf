package api

import (
	"context"
	"log"
	"os"

	"github.com/influxdata/telegraf/plugins/inputs/zmon/protodata"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type CommandAPI interface {
	CheckCommand() (*protodata.Command, error)
}

type commandAPI struct {
	zmonAccountID  string
	token          string
	serverEndpoint string
}

func NewCommandAPI(zmonAccountID, token, serverEndpoint string) CommandAPI {
	return &commandAPI{
		zmonAccountID:  zmonAccountID,
		token:          token,
		serverEndpoint: serverEndpoint,
	}
}

func (api *commandAPI) createClient() (*grpc.ClientConn, protodata.CommandServerClient, error) {
	conn, err := createConnection(api.serverEndpoint)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create connection to server")
	}

	client := protodata.NewCommandServerClient(conn)
	return conn, client, nil
}

func (api *commandAPI) CheckCommand() (*protodata.Command, error) {
	log.Println("D! ZMON, Check Command")

	conn, client, err := api.createClient()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create grpc client")
	}
	defer conn.Close()

	hostname, err := os.Hostname()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get hostname")
	}

	req := &protodata.CheckCommandRequest{
		Auth: &protodata.TenantAuth{
			ZmonAccountID: api.zmonAccountID,
			Token:         api.token,
		},
		Hostname: hostname,
	}
	resp, err := client.CheckCommand(context.Background(), req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get response from grpc server")
	}
	return resp.Command, nil
}
