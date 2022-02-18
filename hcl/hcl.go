package hcl

import (
	p "github.com/gulducat/go-run-programs/program"
	"github.com/hashicorp/hcl/v2/hclsimple"
)

type Config struct {
	Programs []p.Program `hcl:"program,block"`
}

func ParseHCL(path string) (Config, error) {
	c := Config{}
	err := hclsimple.DecodeFile(path, nil, &c)
	return c, err
}

func RunFromHCL(path string) (func(), error) {
	c, err := ParseHCL(path)
	if err != nil {
		return func() {}, err
	}
	return p.RunInBackground(c.Programs...)
}
