package microutils

import (
	"context"

	micro "github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
)

type Client struct {
	client client.Client
}

func NewClient() *Client {
	service := micro.NewService()
	service.Init()
	c := service.Client()
	return &Client{
		client: c,
	}
}

func (c Client) Call(resource, endpoint string, request, response interface{}) (err error) {
	microRequest := c.client.NewRequest(resource, endpoint, request, client.WithContentType("application/json"))
	return c.client.Call(context.TODO(), microRequest, response)
}

