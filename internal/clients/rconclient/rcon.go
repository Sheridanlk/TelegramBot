package rconclient

import (
	"TelegramBot/lib/e"
	"time"

	"github.com/gorcon/rcon"
)

type Client struct {
	adress   string
	password string
	timeout  time.Duration

	conn *rcon.Conn
}

func New(addr string, password string, timeout time.Duration) *Client {
	return &Client{
		adress:   addr,
		password: password,
		timeout:  timeout,
	}
}

func (c *Client) Connect() error {
	if c.conn != nil {
		return nil
	}
	conn, err := rcon.Dial(
		c.adress,
		c.password,
		rcon.SetDialTimeout(c.timeout),
		rcon.SetDeadline(c.timeout),
	)
	if err != nil {
		return e.Wrap("can't connect to rcon", err)
	}
	c.conn = conn
	return nil
}

func (c *Client) Execute(command string) (resp string, err error) {
	if c.conn == nil {
		err := c.Connect()
		if err != nil {
			return "", err
		}
	}

	resp, err = c.conn.Execute(command)
	if err != nil {
		return "", e.Wrap("can't execute command", err)
	}
	return resp, nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil
		return err
	}
	return nil
}
