package model

import (
	"time"

	"github.com/openshift/osin"
)

// Collections names
const (
	AuthorizeDataCollectionName = "authorizeData"
	AccessDataCollectionName    = "accessData"
)

// AccessData represents an access grant (tokens, expiration, client, etc)
type AccessData struct {
	// Client information
	Client *Service

	// Authorize data, for authorization code
	AuthorizeData *osin.AuthorizeData

	// Previous access data, for refresh token
	AccessData *osin.AccessData

	// Access token
	AccessToken string

	// Refresh Token. Can be blank
	RefreshToken string

	// Token expiration in seconds
	ExpiresIn int32

	// Requested scope
	Scope string

	// Redirect Uri from request
	RedirectUri string

	// Date created
	CreatedAt time.Time

	UserData *UserData
}

type UserData struct {
	ID            string
	ApplicationID string
	ClientID      string
	IsDevice      bool
}

func NewAccessData(o *osin.AccessData) *AccessData {
	accessData := new(AccessData)
	accessData.AccessData = o.AccessData
	accessData.AccessToken = o.AccessToken
	accessData.AuthorizeData = o.AuthorizeData
	//accessData.Client = o.Client.(*Service)
	accessData.CreatedAt = o.CreatedAt
	accessData.ExpiresIn = o.ExpiresIn
	accessData.RedirectUri = o.RedirectUri
	accessData.RefreshToken = o.RefreshToken
	accessData.Scope = o.Scope
	//accessData.UserData = o.UserData.(*UserData)

	return accessData
}

func (o *AccessData) ToOsinAccessData() *osin.AccessData {
	accessData := new(osin.AccessData)
	accessData.AccessData = o.AccessData
	accessData.AccessToken = o.AccessToken
	accessData.AuthorizeData = o.AuthorizeData
	accessData.Client = o.Client
	accessData.CreatedAt = o.CreatedAt
	accessData.ExpiresIn = o.ExpiresIn
	accessData.RedirectUri = o.RedirectUri
	accessData.RefreshToken = o.RefreshToken
	accessData.Scope = o.Scope
	accessData.UserData = o.UserData

	return accessData
}

type Service struct {
	ID                  string `datastore:"ID"`
	ApplicationID       string `datastore:"applicationID"`
	ClientID            string `datastore:"clientID"`
	Description         string `datastore:"description"`
	Secret              string `datastore:"secret"`
	CallbackURL         string `datastore:"callbackURL"`
	MicroserviceDNS     string `datastore:"microserviceDNS"`
	CreatedAt           int64  `datastore:"createdAt"`
	UpdatedAt           int64  `datastore:"updatedAt"`
	Deleted             bool   `datastore:"deleted"`
	MicroserviceContext string `datastore:"microserviceContext"`
	ACLs                *Acls  `datastore:"acls"`
}

type Acls struct {
	RegularPolicies *RegularPolicies
}

type RegularPolicies struct {
	Read  []string // ALL | NONE | ANONYMOUS | SELF | an ID... GET:[/datacollection/v2/*]:SELF,POST:[/datacollection/v2/*]:ALL
	Write []string
	Admin []string
}

func (s *Service) GetId() string {
	return s.ID
}

func (s *Service) GetSecret() string {
	return s.Secret
}

func (s *Service) GetRedirectUri() string {
	return s.MicroserviceDNS
}

func (s *Service) GetUserData() interface{} {
	return s.ACLs
}

// Implement the ClientSecretMatcher interface
func (s *Service) ClientSecretMatches(secret string) bool {
	return s.Secret == secret
}

func (s *Service) CopyFrom(client osin.Client) {
	s.ID = client.GetId()
	s.Secret = client.GetSecret()
	s.MicroserviceDNS = client.GetRedirectUri()
	s.ACLs = (client.GetUserData()).(*Acls)
}
