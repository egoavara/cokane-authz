// Code generated by ogen, DO NOT EDIT.

package manage

import (
	"github.com/go-faster/errors"
)

type GetRaftNodeOK struct {
	Nodes []RaftMetaNode `json:"nodes"`
}

// GetNodes returns the value of Nodes.
func (s *GetRaftNodeOK) GetNodes() []RaftMetaNode {
	return s.Nodes
}

// SetNodes sets the value of Nodes.
func (s *GetRaftNodeOK) SetNodes(val []RaftMetaNode) {
	s.Nodes = val
}

// NewOptBool returns new OptBool with value set to v.
func NewOptBool(v bool) OptBool {
	return OptBool{
		Value: v,
		Set:   true,
	}
}

// OptBool is optional bool.
type OptBool struct {
	Value bool
	Set   bool
}

// IsSet returns true if OptBool was set.
func (o OptBool) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptBool) Reset() {
	var v bool
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptBool) SetTo(v bool) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptBool) Get() (v bool, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptBool) Or(d bool) bool {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// Ref: #/components/schemas/RaftBaseNode
type RaftBaseNode struct {
	ID   string `json:"id"`
	Addr string `json:"addr"`
}

// GetID returns the value of ID.
func (s *RaftBaseNode) GetID() string {
	return s.ID
}

// GetAddr returns the value of Addr.
func (s *RaftBaseNode) GetAddr() string {
	return s.Addr
}

// SetID sets the value of ID.
func (s *RaftBaseNode) SetID(val string) {
	s.ID = val
}

// SetAddr sets the value of Addr.
func (s *RaftBaseNode) SetAddr(val string) {
	s.Addr = val
}

// Merged schema.
// Ref: #/components/schemas/RaftMetaNode
type RaftMetaNode struct {
	ID     string     `json:"id"`
	Addr   string     `json:"addr"`
	Status RaftStatus `json:"status"`
	Role   RaftRole   `json:"role"`
}

// GetID returns the value of ID.
func (s *RaftMetaNode) GetID() string {
	return s.ID
}

// GetAddr returns the value of Addr.
func (s *RaftMetaNode) GetAddr() string {
	return s.Addr
}

// GetStatus returns the value of Status.
func (s *RaftMetaNode) GetStatus() RaftStatus {
	return s.Status
}

// GetRole returns the value of Role.
func (s *RaftMetaNode) GetRole() RaftRole {
	return s.Role
}

// SetID sets the value of ID.
func (s *RaftMetaNode) SetID(val string) {
	s.ID = val
}

// SetAddr sets the value of Addr.
func (s *RaftMetaNode) SetAddr(val string) {
	s.Addr = val
}

// SetStatus sets the value of Status.
func (s *RaftMetaNode) SetStatus(val RaftStatus) {
	s.Status = val
}

// SetRole sets the value of Role.
func (s *RaftMetaNode) SetRole(val RaftRole) {
	s.Role = val
}

// Raft node role.
// Ref: #/components/schemas/RaftRole
type RaftRole string

const (
	RaftRoleNonVoter RaftRole = "non-voter"
	RaftRoleVoter    RaftRole = "voter"
)

// AllValues returns all RaftRole values.
func (RaftRole) AllValues() []RaftRole {
	return []RaftRole{
		RaftRoleNonVoter,
		RaftRoleVoter,
	}
}

// MarshalText implements encoding.TextMarshaler.
func (s RaftRole) MarshalText() ([]byte, error) {
	switch s {
	case RaftRoleNonVoter:
		return []byte(s), nil
	case RaftRoleVoter:
		return []byte(s), nil
	default:
		return nil, errors.Errorf("invalid value: %q", s)
	}
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *RaftRole) UnmarshalText(data []byte) error {
	switch RaftRole(data) {
	case RaftRoleNonVoter:
		*s = RaftRoleNonVoter
		return nil
	case RaftRoleVoter:
		*s = RaftRoleVoter
		return nil
	default:
		return errors.Errorf("invalid value: %q", data)
	}
}

// Raft node status.
// Ref: #/components/schemas/RaftStatus
type RaftStatus string

const (
	RaftStatusLeader   RaftStatus = "leader"
	RaftStatusFollower RaftStatus = "follower"
	RaftStatusNonVoter RaftStatus = "non-voter"
)

// AllValues returns all RaftStatus values.
func (RaftStatus) AllValues() []RaftStatus {
	return []RaftStatus{
		RaftStatusLeader,
		RaftStatusFollower,
		RaftStatusNonVoter,
	}
}

// MarshalText implements encoding.TextMarshaler.
func (s RaftStatus) MarshalText() ([]byte, error) {
	switch s {
	case RaftStatusLeader:
		return []byte(s), nil
	case RaftStatusFollower:
		return []byte(s), nil
	case RaftStatusNonVoter:
		return []byte(s), nil
	default:
		return nil, errors.Errorf("invalid value: %q", s)
	}
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *RaftStatus) UnmarshalText(data []byte) error {
	switch RaftStatus(data) {
	case RaftStatusLeader:
		*s = RaftStatusLeader
		return nil
	case RaftStatusFollower:
		*s = RaftStatusFollower
		return nil
	case RaftStatusNonVoter:
		*s = RaftStatusNonVoter
		return nil
	default:
		return errors.Errorf("invalid value: %q", data)
	}
}
