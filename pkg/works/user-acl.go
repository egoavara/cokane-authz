package works

import (
	"egoavara.net/authz/pkg/acts"
	"github.com/google/uuid"
	"go.temporal.io/sdk/workflow"
)

type UserAclRequest struct {
	UserId           *uuid.UUID
	InitialObjectIds []uuid.UUID
	Searchers        *UserAclSearcher
}
type UserAclSearcher struct {
	Name                *string
	Labels              *[]string
	WithBlockState      bool
	WithLastAllowTime   bool
	WithLastDenyTime    bool
	WithLastBlockTime   bool
	WithLastUnblockTime bool
	WithLastModifyTime  bool
}

type UserAclResponse struct {
	UserId   uuid.UUID
	ObjectId []uuid.UUID
}

type UserAclAllowSignal struct {
	ObjectIds []uuid.UUID
}

type UserAclSignal struct {
	IsBlocked bool
	expired   bool
}

func (s *UserAclSignal) Listen(ctx workflow.Context) {
	log := workflow.GetLogger(ctx)
	// actx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
	// 	ScheduleToStartTimeout: time.Hour,
	// })
	for {
		selector := workflow.NewSelector(ctx)
		// selector.AddReceive(workflow.GetSignalChannel(ctx, "Block"), func(c workflow.ReceiveChannel, more bool) {
		// 	c.Receive(ctx, nil)
		// 	a.Signal1Received = true
		// 	log.Info("Signal1 Received")
		// })
		// selector.AddReceive(workflow.GetSignalChannel(ctx, "Unblock"), func(c workflow.ReceiveChannel, more bool) {
		// 	c.Receive(ctx, nil)
		// 	a.Signal2Received = true
		// 	log.Info("Signal2 Received")
		// })
		selector.AddReceive(workflow.GetSignalChannel(ctx, "Allow"), func(c workflow.ReceiveChannel, more bool) {
			var signal UserAclAllowSignal
			c.Receive(ctx, &signal)
			log.Info("Allow Received", "Signal", signal)
		})
		// selector.AddReceive(workflow.GetSignalChannel(ctx, "Deny"), func(c workflow.ReceiveChannel, more bool) {
		// 	c.Receive(ctx, nil)
		// 	a.Signal3Received = true
		// 	log.Info("Signal3 Received")
		// })
		// selector.AddReceive(workflow.GetSignalChannel(ctx, "Clear"), func(c workflow.ReceiveChannel, more bool) {
		// 	c.Receive(ctx, nil)
		// 	a.Signal3Received = true
		// 	log.Info("Signal3 Received")
		// })
		selector.AddReceive(workflow.GetSignalChannel(ctx, "Expire"), func(c workflow.ReceiveChannel, more bool) {
			c.Receive(ctx, nil)
			log.Info("Expire Received")
			s.expired = true
		})
		selector.Select(ctx)
	}

}
func (s *UserAclSignal) IsExpired() bool {
	return s.expired
}
func UserAcl(ctx workflow.Context, req *UserAclRequest) (*UserAclResponse, error) {
	var logger = workflow.GetLogger(ctx)
	var signal UserAclSignal
	var act acts.Act
	var userId uuid.UUID

	logger.Info("UserAcl Workflow Started")
	// ========================================
	// 값 초기화
	workflow.SideEffect(ctx, func(ctx workflow.Context) interface{} {
		if req.UserId == nil {
			return uuid.Must(uuid.NewV7())
		}
		return req.UserId
	}).Get(&userId)
	// ========================================
	// 사후 작업 정리
	defer func() {
		// ACL 정리
		newCtx, _ := workflow.NewDisconnectedContext(ctx)
		err := workflow.ExecuteActivity(newCtx, act.ClearAcl, &acts.ClearAclRequest{User: userId.String()}).
			Get(newCtx, nil)
		if err != nil {
			logger.Error("ClearACL Request Failed", "Error", err)
		}
	}()
	// ========================================
	// 쿼리 핸들러 등록
	workflow.SetQueryHandler(ctx, "GetACL", func() (string, error) {
		return "", nil
	})
	// ========================================
	// 시그널 수신
	workflow.Go(ctx, signal.Listen)
	workflow.Await(ctx, signal.IsExpired)
	return &UserAclResponse{
		UserId: userId,
	}, nil
}
