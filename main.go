package main

import (
	"onlyone_smc/api"
	"onlyone_smc/internal/env"
	"os"
)

// @title OnlyOne Smart Contract
// @version 1.0
// @description Documentation Smart Contract
// @termsOfService https://www.bjungle.net/terms/
// @contact.name API Support
// @contact.email info@bjungle.net
// @license.name Software Owner
// @license.url https://www.bjungle.net/terms/licenses
// @host http://172.174.77.149:2054
// @tag.name Credentials
// @tag.description Credentials of OnlyOne Clients
// @tag.name User
// @tag.description Methods of user
// @tag.name Authentication
// @tag.description Methods of Authentication
// @tag.Name Categories
// @tag.description Categories of credentials
// @BasePath /
func main() {
	c := env.NewConfiguration()
	_ = os.Setenv("AWS_ACCESS_KEY_ID", c.Aws.AWSACCESSKEYID)
	_ = os.Setenv("AWS_SECRET_ACCESS_KEY", c.Aws.AWSSECRETACCESSKEY)
	_ = os.Setenv("AWS_DEFAULT_REGION", c.Aws.AWSDEFAULTREGION)

	api.Start(c.App.Port, c.App.ServiceName, c.App.LoggerHttp, c.App.AllowedDomains)
}
