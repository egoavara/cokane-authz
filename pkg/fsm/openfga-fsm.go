package fsm

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/hashicorp/raft"
	openfga "github.com/openfga/go-sdk"
)

var (
	ErrFga = errors.New("openfga error")
)

type OpenFgaFSM struct {
	currentModelId string
	modelLock      sync.RWMutex
	FgaClient      *openfga.APIClient
}

func NewOpenFgaFSM(fgaClient *openfga.APIClient) *OpenFgaFSM {
	return &OpenFgaFSM{
		FgaClient: fgaClient,
	}
}

func (s *OpenFgaFSM) Apply(log *raft.Log) interface{} {
	var command OpenFgaFSMCommand
	err := json.Unmarshal(log.Data, &command)
	fmt.Println("Applying command:", string(log.Data))
	if err != nil {
		return err
	}
	switch command.Type {
	case OpenFgaFSMCommandTypeModel:
		return s.applyModel(command.Model)
	case OpenFgaFSMCommandTypeWrite:
		fallthrough
	case OpenFgaFSMCommandTypeDelete:
		return s.applyNonModel(command)
	default:
		return OpenFgaFSMCommandResult{
			Ok:    false,
			Error: "unknown command type",
		}
	}
}

func (s *OpenFgaFSM) ensureStore() {
	if s.FgaClient.GetStoreId() == "" {
		s.modelLock.Lock()
		defer s.modelLock.Unlock()
		if s.FgaClient.GetStoreId() == "" {
			fga := s.FgaClient
			storeList, _, err := fga.OpenFgaApi.ListStores(context.Background()).PageSize(10).Execute()
			if err != nil {
				panic(err)
			}
			if len(storeList.GetStores()) > 0 {
				s.FgaClient.SetStoreId(storeList.GetStores()[0].Id)
				return
			}
			store, _, err := fga.OpenFgaApi.CreateStore(context.Background()).
				Body(openfga.CreateStoreRequest{
					Name: "default",
				}).
				Execute()
			if err != nil {
				panic(err)
			}
			s.FgaClient.SetStoreId(store.Id)
			return
		}
	}
}
func (s *OpenFgaFSM) applyModel(model *openfga.AuthorizationModel) OpenFgaFSMCommandResult {
	s.ensureStore()
	s.modelLock.Lock()
	defer s.modelLock.Unlock()
	fga := s.FgaClient
	req := openfga.NewWriteAuthorizationModelRequest(model.GetTypeDefinitions(), model.GetSchemaVersion())
	req.SetConditions(model.GetConditions())
	resp, _, err := fga.OpenFgaApi.WriteAuthorizationModel(context.Background()).
		Body(*req).
		Execute()
	if err != nil {
		return OpenFgaFSMCommandResult{
			Ok:    false,
			Error: err.Error(),
		}
	}
	s.currentModelId = resp.GetAuthorizationModelId()
	return OpenFgaFSMCommandResult{
		Ok: true,
	}
}
func (s *OpenFgaFSM) applyNonModel(command OpenFgaFSMCommand) OpenFgaFSMCommandResult {
	s.ensureStore()
	s.modelLock.RLock()
	defer s.modelLock.RUnlock()
	fga := s.FgaClient

	req := openfga.NewWriteRequest()
	req.SetWrites(*openfga.NewWriteRequestWrites(command.Write))
	req.SetDeletes(*openfga.NewWriteRequestDeletes(command.Delete))
	_, _, err := fga.OpenFgaApi.Write(context.Background()).
		Body(*req).
		Execute()
	if err != nil {
		return OpenFgaFSMCommandResult{
			Ok:    false,
			Error: err.Error(),
		}
	}
	return OpenFgaFSMCommandResult{
		Ok: true,
	}
}

func (s *OpenFgaFSM) Snapshot() (raft.FSMSnapshot, error) {
	fmt.Println("Snapshotting")
	// unlock은 snapshot이 release될 때 수행됩니다.
	s.modelLock.Lock()
	return &OpenFgaSnapshot{FSM: s}, nil
}

func (s *OpenFgaFSM) Restore(snapshot io.ReadCloser) error {
	s.ensureStore()
	fmt.Println("Restoring snapshot")
	s.modelLock.Lock()
	defer s.modelLock.Unlock()
	defer snapshot.Close()

	fga := s.FgaClient
	reader := bufio.NewReader(snapshot)
	// 인증 모델을 복원
	modelJson, err := reader.ReadBytes('\n')
	if err != nil {
		return err
	}
	model := openfga.ReadAuthorizationModelResponse{}
	err = json.Unmarshal(modelJson, &model)
	if err != nil {
		return err
	}
	if !model.HasAuthorizationModel() {
		return errors.New("model not found")
	}
	modelReq := openfga.NewWriteAuthorizationModelRequest(model.GetAuthorizationModel().TypeDefinitions, model.GetAuthorizationModel().SchemaVersion)
	modelReq.SetConditions(*model.GetAuthorizationModel().Conditions)
	modelResp, _, err := fga.OpenFgaApi.WriteAuthorizationModel(context.Background()).
		Body(*modelReq).
		Execute()
	if err != nil {
		return err
	}
	// 복원된 모델 id를 저장
	s.currentModelId = modelResp.GetAuthorizationModelId()
	// 튜플을 복원
	writeReq := openfga.NewWriteRequestWithDefaults()
	writeReq.SetAuthorizationModelId(s.currentModelId)
	chunk := make([]openfga.TupleKey, 0, 1000)
	for {
		tupleJson, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		tuple := openfga.Tuple{}
		err = json.Unmarshal(tupleJson, &tuple)
		if err != nil {
			return err
		}
		chunk = append(chunk, tuple.GetKey())
		if len(chunk) == 1000 {
			writeReq.Writes.SetTupleKeys(chunk)
			_, _, err := fga.OpenFgaApi.Write(context.Background()).
				Body(*writeReq).
				Execute()
			if err != nil {
				return err
			}
			chunk = chunk[:0]
		}
	}
	if len(chunk) > 0 {
		writeReq.Writes.SetTupleKeys(chunk)
		_, _, err := fga.OpenFgaApi.Write(context.Background()).
			Body(*writeReq).
			Execute()
		if err != nil {
			return err
		}
	}
	return nil
}

type OpenFgaSnapshot struct {
	FSM *OpenFgaFSM
}

func (s *OpenFgaSnapshot) Persist(sink raft.SnapshotSink) error {
	fmt.Println("Persisting snapshot")
	s.FSM.ensureStore()
	fga := s.FSM.FgaClient
	modelId := s.FSM.currentModelId
	model, _, err := fga.OpenFgaApi.ReadAuthorizationModel(context.Background(), modelId).Execute()
	if err != nil {
		return err
	}
	modelJson, err := model.MarshalJSON()
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(sink)
	defer sink.Close()
	// 인증 모델에 대한 정보를 스냅샷에 저장
	_, err = sink.Write(modelJson)
	if err != nil {
		sink.Cancel()
		return err
	}
	writer.Flush()
	// 인증 모델의 실제 데이터(aka. tuple)를 스냅샷에 저장
	req := openfga.NewReadRequestWithDefaults()
	req.SetPageSize(100)
	for {
		res, _, err := fga.OpenFgaApi.Read(context.Background()).
			Body(*req).
			Execute()
		if err != nil {
			sink.Cancel()
			return err
		}
		// 모든 튜플을 스냅샷에 저장
		tuples := res.GetTuples()
		if len(tuples) == 0 {
			break
		}
		for _, tuple := range tuples {
			tupleJson, err := tuple.MarshalJSON()
			if err != nil {
				sink.Cancel()
				return err
			}
			writer.WriteString("\n")
			_, err = sink.Write(tupleJson)
			if err != nil {
				sink.Cancel()
				return err
			}
			writer.Flush()
		}
		// 다음 페이지가 있는지 확인
		continuationToken, ok := res.GetContinuationTokenOk()
		if !ok {
			break
		}
		req.SetContinuationToken(*continuationToken)
	}
	return nil
}
func (s *OpenFgaSnapshot) Release() {
	s.FSM.modelLock.Unlock()
}

const (
	OpenFgaFSMCommandTypeModel  = "model"
	OpenFgaFSMCommandTypeWrite  = "write"
	OpenFgaFSMCommandTypeDelete = "delete"
)

type OpenFgaFSMCommand struct {
	Type   string                             `json:"type"`
	Model  *openfga.AuthorizationModel        `json:"model"`
	Write  []openfga.TupleKey                 `json:"write"`
	Delete []openfga.TupleKeyWithoutCondition `json:"delete"`
}

func NewOpenFgaFSMCommandModel(model *openfga.AuthorizationModel) *OpenFgaFSMCommand {
	return &OpenFgaFSMCommand{
		Type:  OpenFgaFSMCommandTypeModel,
		Model: model,
	}
}
func NewOpenFgaFSMCommandWrite(keys []openfga.TupleKey) *OpenFgaFSMCommand {
	return &OpenFgaFSMCommand{
		Type:  OpenFgaFSMCommandTypeWrite,
		Write: keys,
	}
}
func NewOpenFgaFSMCommandDelete(keys []openfga.TupleKeyWithoutCondition) *OpenFgaFSMCommand {
	return &OpenFgaFSMCommand{
		Type:   OpenFgaFSMCommandTypeDelete,
		Delete: keys,
	}
}

type OpenFgaFSMCommandResult struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}
