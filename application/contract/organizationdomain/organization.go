package organization

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

type Organization struct {
	foundation.BaseContract
	Name   string
	Public bool
	// member ref => role ref
	Members map[core.RecordRef]core.RecordRef
}

func (o *Organization) GetName() (string, error) {
	return o.Name, nil
}

func (o *Organization) IsPublic() (bool, error) {
	return o.Public, nil
}

func NewOrganization(name string, isPublic bool) (*Organization, error) {
	return &Organization{
		Name:   name,
		Public: isPublic,
	}, nil
}

func AddMember() {

}

func RemoveMember() {

}

func ReplaceMember() {

}

func JoinDomain() {

}

func LeaveDomain() {

}

func AddRole() {

}

func RemoveRole() {

}

func AssignRoleToMember() {

}

func ChangeMemberRole() {

}
