package migrate

import "testing"

func TestSchemaBuild_Upgrade(t *testing.T) {
	dsn := "username:passwordk@tcp(host:port)/dbname"
	build := NewSchemaBuild()
	//build.Upgrade("./.model", dsn)
	build.Upgrade("./.model/user.json", dsn)
}
