package browse

import (
	"container/list"
	"regexp"
	"strings"

	"github.com/c2stack/c2g/c2"
	"github.com/c2stack/c2g/node"
)

// Role-based Access Control
//   see https://en.wikipedia.org/wiki/Role-based_access_control
type Rbac struct {
	Roles map[string]*Role
}

type Role struct {
	Id     string
	Access *list.List
}

func NewRole() *Role {
	return &Role{
		Access: list.New(),
	}
}

type AccessControl struct {
	Path        string
	Permissions Permission
	pathRegx    *regexp.Regexp
}

func (self *AccessControl) Matches(targetPath string) (bool, error) {
	if self.pathRegx == nil {
		var err error
		if self.pathRegx, err = regexp.Compile(self.Path); err != nil {
			return false, err
		}
	}

	return self.pathRegx.MatchString(targetPath), nil
}

type Permission uint

const (
	Read Permission = 1 << iota
	Write
	Execute
	Subscribe
)

func DecodePermission(s []string) Permission {
	var p Permission
	for _, pstr := range s {
		switch pstr {
		case "r":
			p = p | Read
		case "w":
			p = p | Write
		case "x":
			p = p | Execute
		case "s":
			p = p | Subscribe
		}

	}
	return p
}

func EncodePermission(p Permission) []string {
	encoded := [4]string{}
	var j int
	if (p & Read) == Read {
		encoded[j] = "r"
		j++
	}
	if (p & Write) == Write {
		encoded[j] = "w"
		j++
	}
	if (p & Execute) == Execute {
		encoded[j] = "x"
		j++
	}
	if (p & Subscribe) == Subscribe {
		encoded[j] = "s"
		j++
	}
	return encoded[0:j]
}

var UnauthorizedError = c2.NewErrC("Unauthorized", 401)

func (self *Role) CheckListPreConstraints(r *node.ListRequest) (bool, error) {
	if r.IsNavigation() {
		return true, nil
	}
	p := Read
	if r.New {
		p = Write
	}
	if r.Key != nil {
		return self.check(r.Selection.Path.StringNoModule()+"="+node.EncodeKey(r.Key), p)
	}
	return true, nil
}

func (self *Role) CheckContainerPreConstraints(r *node.ChildRequest) (bool, error) {
	if r.IsNavigation() {
		return true, nil
	}
	p := Read
	if r.New {
		p = Write
	}
	return self.check(r.Selection.Path.StringNoModule()+"/"+r.Meta.GetIdent(), p)
}

func (self *Role) CheckFieldPreConstraints(r *node.FieldRequest, hnd *node.ValueHandle) (bool, error) {
	if r.IsNavigation() {
		return true, nil
	}
	p := Read
	if r.Write {
		p = Write
	}
	return self.check(r.Selection.Path.StringNoModule()+"/"+r.Meta.GetIdent(), p)
}

func (self *Role) CheckActionPreConstraints(r *node.ActionRequest) (bool, error) {
	return self.check(r.Selection.Path.StringNoModule(), Execute)
}

func (self *Role) check(targetPath string, p Permission) (bool, error) {
	// HACK: occasional leading path messing things up, find out why this is inconsistent
	targetPath = strings.TrimLeft(targetPath, "/")

	i := self.Access.Front()
	for i != nil {
		ac := i.Value.(*AccessControl)
		if found, err := ac.Matches(targetPath); found {
			if (ac.Permissions & p) == p {
				return true, nil
			}
		} else if err != nil {
			return false, err
		}
		i = i.Next()
	}
	return false, UnauthorizedError
}

type NoAccess struct {
}

func (self NoAccess) CheckListPreConstraints(r *node.ListRequest) (bool, error) {
	return false, UnauthorizedError
}

func (self NoAccess) CheckContainerPreConstraints(r *node.ChildRequest) (bool, error) {
	return false, UnauthorizedError
}

func (self NoAccess) CheckFieldPreConstraints(r *node.FieldRequest, hnd *node.ValueHandle) (bool, error) {
	return false, UnauthorizedError
}
