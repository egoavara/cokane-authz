package acts

import (
	"context"
	"encoding/json"

	openfga "github.com/openfga/go-sdk"
)

type ModifyAclsRequest struct {
	Writes  []openfga.TupleKey                 `json:"writes"`
	Deletes []openfga.TupleKeyWithoutCondition `json:"deletes"`
}
type ModifyAclsResponse struct {
	AuthorizationModelId string                             `json:"authorization_model_id"`
	Writes               []openfga.TupleKey                 `json:"writes"`
	Deletes              []openfga.TupleKeyWithoutCondition `json:"deletes"`
	Raw                  map[string]interface{}             `json:"-"`
}

func (s *Act) ModifyAcls(ctx context.Context, acls ModifyAclsRequest) (*ModifyAclsResponse, error) {
	req := openfga.NewWriteRequest()
	req.SetWrites(*openfga.NewWriteRequestWrites(acls.Writes))
	req.SetDeletes(*openfga.NewWriteRequestDeletes(acls.Deletes))
	res, _, err := s.OpenFga.OpenFgaApi.Write(ctx).
		Body(*req).
		Execute()
	if err != nil {
		return nil, err
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	var result ModifyAclsResponse
	err = json.Unmarshal(resBytes, &result)
	if err != nil {
		return nil, err
	}
	result.Raw = res
	return &result, nil
}
