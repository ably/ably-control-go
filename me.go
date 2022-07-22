package control

type Me struct {
	Token   Token   `json:"token"`
	User    User    `json:"user"`
	Account Account `json:"account"`
}

type Token struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Capabilities []string `json:"capabilities"`
}

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type Account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (c *Client) Me() (Me, error) {
	var me Me
	err := c.request("GET", "/me", nil, &me)
	if err != nil {
		return me, err
	}
	return me, nil
}
