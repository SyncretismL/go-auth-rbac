package rbac

import (
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

type RBAC map[string]Resource

type Resource map[string]Endpoint

type Endpoint map[string]Permission

type Permission struct {
	Allow bool `yaml:"allow"`
}

func FromFile(path string) (*RBAC, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	rbac := &RBAC{}
	err = yaml.Unmarshal(f, rbac)
	if err != nil {
		return nil, err
	}

	return rbac, nil
}

func (rbac RBAC) Authorize(r *http.Request, role, resource, endpoint string) error {
	permission, exists := rbac[role][resource][endpoint]
	if !exists {
		return errors.New("unknown role")
	}

	if !permission.Allow {
		return errors.New("you have no permission")
	}

	return nil
}
