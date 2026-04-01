package control

// A struct representing the settable fields of an Ably application.
type NewApp struct {
	// The application ID.
	ID string `json:"id,omitempty"`
	// The application name.
	Name string `json:"name,omitempty"`
	// The application status. Disabled applications will not accept
	// new connections and will return an error to all clients.
	Status string `json:"status,omitempty"`
	// Enforce TLS for all connections. This setting overrides any channel setting.
	TLSOnly bool `json:"tlsOnly"`
	// The Firebase Cloud Messaging key.
	FcmKey string `json:"fcmKey"`
	// The Firebase Service Account key. To use the service account key you must also provide a projectId.
	FcmServiceAccount string `json:"fcmServiceAccount"`
	// The Firebase Project ID. To authenticate with firebase you must also provide a service account key.
	FcmProjectId string `json:"fcmProjectId"`
	// The Apple Push Notification service certificate.
	// This field can only be used to set a new value,
	// it will not be populated by queries.
	ApnsCertificate string `json:"apnsCertificate"`
	// The Apple Push Notification service private key.
	// This field can only be used to set a new value,
	// it will not be populated by queries.
	ApnsPrivateKey string `json:"apnsPrivateKey"`
	// Use the Apple Push Notification service sandbox endpoint.
	ApnsUseSandboxEndpoint bool `json:"apnsUseSandboxEndpoint"`
	// The APNs authentication type. Either "certificate" or "token".
	ApnsAuthType *string `json:"apnsAuthType,omitempty"`
	// The APNs signing key (.p8 file contents) for token-based authentication.
	// Write-only: will not be populated by queries.
	ApnsSigningKey *string `json:"apnsSigningKey,omitempty"`
	// The 10-character Key ID for APNs token-based authentication.
	ApnsSigningKeyId *string `json:"apnsSigningKeyId,omitempty"`
	// The Team ID (issuer key) for APNs token-based authentication.
	ApnsIssuerKey *string `json:"apnsIssuerKey,omitempty"`
	// The bundle ID used as the APNs topic header.
	ApnsTopicHeader *string `json:"apnsTopicHeader,omitempty"`
}

// A struct representing an Ably application.
type App struct {
	// The application ID.
	ID string `json:"id,omitempty"`
	// The ID of your Ably account.
	AccountID string `json:"accountId,omitempty"`
	// The application name.
	Name string `json:"name,omitempty"`
	// The application status. Disabled applications will not accept
	// new connections and will return an error to all clients.
	Status string `json:"status,omitempty"`
	// Enforce TLS for all connections. This setting overrides any channel setting.
	TLSOnly bool `json:"tlsOnly"`
	// The Firebase Cloud Messaging key.
	FcmKey string `json:"fcmKey"`
	// The Firebase Service Account key. To use the service account key you must also provide a projectId.
	FcmServiceAccount string `json:"fcmServiceAccount"`
	// The Firebase Project ID. To authenticate with firebase you must also provide a service account key.
	FcmProjectId string `json:"fcmProjectId"`
	// Whether a Firebase service account key has been configured.
	FcmServiceAccountConfigured bool `json:"fcmServiceAccountConfigured"`
	// The Apple Push Notification service certificate.
	// This field can only be used to set a new value,
	// it will not be populated by queries.
	ApnsCertificate string `json:"apnsCertificate"`
	// The Apple Push Notification service private key.
	// This field can only be used to set a new value,
	// it will not be populated by queries.
	ApnsPrivateKey string `json:"apnsPrivateKey"`
	// Use the Apple Push Notification service sandbox endpoint.
	ApnsUseSandboxEndpoint bool `json:"apnsUseSandboxEndpoint"`
	// The APNs authentication type. Either "certificate" or "token".
	ApnsAuthType *string `json:"apnsAuthType,omitempty"`
	// Whether an APNs certificate has been configured.
	ApnsCertificateConfigured bool `json:"apnsCertificateConfigured"`
	// Whether an APNs signing key has been configured.
	ApnsSigningKeyConfigured bool `json:"apnsSigningKeyConfigured"`
	// The 10-character Key ID for APNs token-based authentication.
	ApnsSigningKeyId *string `json:"apnsSigningKeyId,omitempty"`
	// The Team ID (issuer key) for APNs token-based authentication.
	ApnsIssuerKey *string `json:"apnsIssuerKey,omitempty"`
	// The bundle ID used as the APNs topic header.
	ApnsTopicHeader *string `json:"apnsTopicHeader,omitempty"`
	// Unix timestamp representing the date and time of creation of the app.
	Created int `json:"created"`
	// Unix timestamp representing the date and time of last modification of the app.
	Modified int `json:"modified"`
}

// Apps fetches a list of all your Ably apps.
func (c *Client) Apps() ([]App, error) {
	var apps []App
	err := c.request("GET", "/accounts/"+c.accountID+"/apps", nil, &apps)
	return apps, err
}

// CreateApp creates a new Ably app.
func (c *Client) CreateApp(app *NewApp) (App, error) {
	var out App
	err := c.request("POST", "/accounts/"+c.accountID+"/apps", app, &out)
	return out, err
}

// UpdateApp updates an existing Ably app.
func (c *Client) UpdateApp(id string, app *NewApp) (App, error) {
	var out App
	err := c.request("PATCH", "/apps/"+id, app, &out)
	return out, err
}

// DeleteApp deletes an Ably app.
func (c *Client) DeleteApp(id string) error {
	err := c.request("DELETE", "/apps/"+id, nil, nil)
	return err
}
