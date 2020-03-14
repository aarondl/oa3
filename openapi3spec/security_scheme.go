package openapi3spec

// SecurityScheme that can be used by the operations. Supported
// schemes are HTTP authentication, an API key (either as a header or as a query
// parameter), OAuth2's common flows (implicit, password, application and access
// code) as defined in RFC6749, and OpenID Connect Discovery.
type SecurityScheme struct {
	Type             string     `json:"type,omitempty" yaml:"type,omitempty"`
	Description      *string    `json:"description,omitempty" yaml:"description,omitempty"`
	Name             string     `json:"name,omitempty" yaml:"name,omitempty"`
	In               string     `json:"in,omitempty" yaml:"in,omitempty"`
	Scheme           string     `json:"scheme,omitempty" yaml:"scheme,omitempty"`
	BearerFormat     *string    `json:"bearerFormat,omitempty" yaml:"bearerFormat,omitempty"`
	Flows            OAuthFlows `json:"flows,omitempty" yaml:"flows,omitempty"`
	OpenIDConnectURL string     `json:"openIdConnectUrl,omitempty" yaml:"openIdConnectUrl,omitempty"`
}

// OAuthFlows allows configuration of supported oauthflows
type OAuthFlows struct {
	Implicit          *OAuthFlow `json:"implicit,omitempty" yaml:"implicit,omitempty"`
	Password          *OAuthFlow `json:"password,omitempty" yaml:"password,omitempty"`
	ClientCredentials *OAuthFlow `json:"clientCredentials,omitempty" yaml:"clientCredentials,omitempty"`
	AuthorizationCode *OAuthFlow `json:"authorizationCode,omitempty" yaml:"authorizationCode,omitempty"`

	Extensions `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// OAuthFlow configuration details
type OAuthFlow struct {
	AuthorizationURL string            `json:"authorizationUrl,omitempty" yaml:"authorizationUrl,omitempty"`
	TokenURL         string            `json:"tokenUrl,omitempty" yaml:"tokenUrl,omitempty"`
	RefreshURL       *string           `json:"refreshUrl,omitempty" yaml:"refreshUrl,omitempty"`
	Scopes           map[string]string `json:"scopes,omitempty" yaml:"scopes,omitempty"`

	Extensions `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// SecuritySchemeRef refers to a security scheme
type SecuritySchemeRef struct {
	Ref string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	*SecurityScheme
}
