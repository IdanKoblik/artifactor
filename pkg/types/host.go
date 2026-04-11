package types

type HostType int

const (
	Gitlab HostType = iota
	Github
)

var hostName = map[HostType]string{
	Gitlab: "gitlab",
	Github: "github",
}

var hostTypeFromName = map[string]HostType{
	"gitlab": Gitlab,
	"github": Github,
}

func HostTypeFromString(s string) (HostType, bool) {
	t, ok := hostTypeFromName[s]
	return t, ok
}

type Host struct {
	Id            int
	Url           string
	Type          HostType
	ApplicationId string
	Secret        string
}

func (h HostType) String() string {
	return hostName[h]
}
