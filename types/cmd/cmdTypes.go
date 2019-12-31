package cmdTypes

type AuthContext struct {
	ContextName string
	Username    string
	Password    string
	Url         string
	Insecure    bool
	Timeout     int
}
