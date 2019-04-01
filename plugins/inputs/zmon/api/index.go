package api

import "google.golang.org/grpc"

func createConnection(serverEndpoint string) (*grpc.ClientConn, error) {
	return grpc.Dial(serverEndpoint, grpc.WithInsecure())
}
