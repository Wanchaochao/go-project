package acl

const (
	ModuleAdmin  = "admin"
	ModuleApplet = "applet"
)

const (
	AuthorityNo   = 0
	AuthorityRead = 1
	AuthorityAll  = 2
)

type Module struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

var Modules = []*Module{
	{Key: ModuleAdmin, Name: "账号权限"},
	{Key: ModuleApplet, Name: "小程序运营"},
}

var AllAuthority = make(Authority)

func init() {
	for _, v := range Modules {
		AllAuthority[v.Key] = AuthorityAll
	}
}
