package control

// The Me struct contains information about the token the current user authenticates with.
type Me struct {
	// The access token used to authenticate.
	Token Token `json:"token"`
	// The user accociated with the used token.
	User User `json:"user"`
	// The account accociated with the used token.
	Account Account `json:"account"`
}

// An access token used to authenticate with the Control API.
type Token struct {
	// The ID of the token.
	ID string `json:"id"`
	// The name of the token.
	Name string `json:"name"`
	// The capabilities of the token.
	Capabilities []string `json:"capabilities"`
}

// User associated with the used token and account.
type User struct {
	// The ID of the user.
	ID int `json:"id"`
	// The user's email address.
	Email string `json:"email"`
}

// The account detials of an Ably account.
type Account struct {
	// The ID of the account.
	ID string `json:"id"`
	// The name of the account.
	Name string `json:"name"`
}

// Me fetches information about the token the current user authenticates with.
func (c *Client) Me() (Me, error) {
	var me Me
	err := c.request("GET", "/me", nil, &me)
	if err != nil {
		return me, err
	}
	return me, nil
}
