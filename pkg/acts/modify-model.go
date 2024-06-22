package acts

import (
	"context"

	openfga "github.com/openfga/go-sdk"
)

type ModifyModelResponse struct {
}

func (s *Act) ModifyModel(ctx context.Context, model openfga.WriteAuthorizationModelRequest) (*openfga.WriteAuthorizationModelResponse, error) {
	req := openfga.NewWriteAuthorizationModelRequestWithDefaults()
	req.SetSchemaVersion(model.SchemaVersion)
	req.SetTypeDefinitions(model.TypeDefinitions)
	if model.Conditions != nil {
		req.SetConditions(*model.Conditions)
	}
	res, _, err := s.OpenFga.OpenFgaApi.WriteAuthorizationModel(ctx).
		Body(*req).
		Execute()
	if err != nil {
		return nil, err
	}
	return &res, nil
}
