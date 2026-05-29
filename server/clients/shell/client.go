package shell

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Command(name string, args []string) *Command {
	return &Command{Name: name, Args: args}
}
