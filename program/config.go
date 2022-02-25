package program

import (
	"github.com/hashicorp/hcl/v2/hclsimple"
)

type Config struct {
	Programs []Program `hcl:"program,block"`
}

func ParseHCL(path string) (Config, error) {
	c := Config{}
	err := hclsimple.DecodeFile(path, nil, &c)
	return c, err
}
