package acts

import (
	"context"

	openfga "github.com/openfga/go-sdk"
)

type ClearAclRequest struct {
	User string `json:"user"`
}
type ClearAclResponse struct {
	Deletes []openfga.TupleKeyWithoutCondition `json:"deletes"`
}

func isRemain(read *openfga.ReadResponse) bool {
	tuples, ok := read.GetTuplesOk()
	if !ok {
		return false
	}
	return len(*tuples) > 0
}
func isContinue(read *openfga.ReadResponse) bool {
	_, ok := read.GetContinuationTokenOk()
	return ok
}
func appendDeletes(out []openfga.TupleKeyWithoutCondition, deletes []openfga.Tuple) []openfga.TupleKeyWithoutCondition {
	for _, delete := range deletes {
		key := delete.GetKey()
		out = append(out, openfga.TupleKeyWithoutCondition{
			User:     key.GetUser(),
			Relation: key.GetRelation(),
			Object:   key.GetObject(),
		})
	}
	return out
}

func (s *Act) deleteKeys(deletes []openfga.Tuple) error {
	deleteKeys := make([]openfga.TupleKeyWithoutCondition, 0, len(deletes))
	for _, delete := range deletes {
		key := delete.GetKey()
		deleteKeys = append(deleteKeys, openfga.TupleKeyWithoutCondition{
			User:     key.GetUser(),
			Relation: key.GetRelation(),
			Object:   key.GetObject(),
		})
	}
	req := openfga.NewWriteRequest()
	req.SetDeletes(*openfga.NewWriteRequestDeletes(deleteKeys))
	_, _, err := s.OpenFga.OpenFgaApi.Write(context.Background()).
		Body(*req).
		Execute()
	if err != nil {
		return err
	}
	return nil
}

func (s *Act) ClearAcl(ctx context.Context, acls ClearAclRequest) (*ClearAclResponse, error) {
	req := openfga.NewReadRequest()
	req.SetPageSize(1000)
	req.SetTupleKey(openfga.ReadRequestTupleKey{
		User: &acls.User,
	})
	res, _, err := s.OpenFga.OpenFgaApi.Read(ctx).
		Body(*req).
		Execute()
	if err != nil {
		return nil, err
	}
	var result ClearAclResponse
	for isRemain(&res) {
		s.deleteKeys(res.GetTuples())
		result.Deletes = appendDeletes(result.Deletes, res.GetTuples())
		if isContinue(&res) {
			req.SetContinuationToken(res.GetContinuationToken())
			res, _, err = s.OpenFga.OpenFgaApi.Read(ctx).
				Body(*req).
				Execute()
			if err != nil {
				return nil, err
			}
		}
	}
	return &result, nil
}
