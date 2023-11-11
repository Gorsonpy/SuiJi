package service

import (
	"time"

	"github.com/XZ0730/runFzu/biz/dal/db"
	"github.com/XZ0730/runFzu/biz/model/multiledger"
	"github.com/XZ0730/runFzu/pkg/errno"
	"github.com/cloudwego/kitex/pkg/klog"
	"golang.org/x/sync/errgroup"
)

func (m *MultiLedgerService) CreateMultiLedger(uid int64, req *multiledger.CreateMLRequest) (code int64, msg string) {
	ml := db.NewMultiLedger(req.GetMultiLedgerName(), req.GetDescription(), req.GetPassword())
	if err := db.CreateMultiLedger(ml); err != nil {
		klog.Error("[mul_l]error:", err.Error())
		return errno.CreateError.ErrorCode, errno.CreateError.ErrorMsg
	}
	if err := db.CreateM_user(ml.MultiLedgerId, uid); err != nil {
		klog.Error("[mul_l]error:", err.Error())
		return errno.CreateError.ErrorCode, errno.CreateError.ErrorMsg
	}
	return errno.StatusSuccessCode, errno.StatusSuccessMsg
}

func (m *MultiLedgerService) JoinMultiledger(uid int64, pwd string) (code int64, msg string) {

	id, err := db.GetMultiLedgerByPassword(pwd)
	if err != nil {
		klog.Error("[multi]error:", err.Error())
		return errno.GetError.ErrorCode, errno.GetError.ErrorMsg
	}
	if err = db.JudgeM_user(id, uid); err == nil {
		klog.Error("[multi]error: user exist")
		return errno.UserExistedError.ErrorCode, errno.UserExistedError.ErrorMsg
	}
	if err = db.CreateM_user(id, uid); err != nil {
		klog.Error("[multi]error:", err.Error())
		return errno.CreateError.ErrorCode, errno.CreateError.ErrorMsg
	}

	return errno.StatusSuccessCode, errno.StatusSuccessMsg
}

func (m *MultiLedgerService) InsertMlConsumption(uid int64, req *multiledger.InsertMlConsumReq) (code int64, msg string) {
	if err := db.JudgeConsumption(req.GetConsId()); err != nil {
		klog.Error("[mul_consumption]error:", err.Error())
		return errno.NotExistError.ErrorCode, errno.NotExistError.ErrorMsg
	}

	if err := db.JudgeM_consumption(req.GetMultiLedgerId(), uid, req.GetConsId()); err == nil {
		klog.Error("[mul_consumption]error: consumption have exist")
		return errno.ExistError.ErrorCode, errno.ExistError.ErrorMsg
	}

	if err := db.CreateM_Consumption(req.GetMultiLedgerId(), uid, req.GetConsId()); err != nil {
		klog.Error("[mul_consumption]error:", err.Error())
		return errno.CreateError.ErrorCode, errno.CreateError.ErrorMsg
	}

	return errno.StatusSuccessCode, errno.StatusSuccessMsg
}

func (m *MultiLedgerService) GetMulConsumption(mid int64) ([]*multiledger.ConsumptionModel, int64, string) {
	cm := make([]*multiledger.ConsumptionModel, 0)

	cl, err := db.GetMl_Consumption(mid)
	if err != nil {
		klog.Error("[mul_consumption]error:", err.Error())
		return nil, errno.GetError.ErrorCode, errno.GetError.ErrorMsg
	}
	var eg errgroup.Group
	for _, val := range cl {
		tmp := val
		eg.Go(func() error {
			vo_g := new(multiledger.ConsumptionModel)
			vo_g.ConsumptionId = tmp.ConsumptionId
			vo_g.ConsumptionName = tmp.ConsumptionName
			vo_g.Description = tmp.Description
			vo_g.Amount = tmp.Amount
			vo_g.TypeId = tmp.TypeId
			vo_g.Store = tmp.Store
			vo_g.ConsumeTime = tmp.ConsumeTime.Format(time.DateTime)
			vo_g.Credential = tmp.Credential
			cm = append(cm, vo_g)
			return nil
		})
	}
	if err = eg.Wait(); err != nil {
		klog.Info("[multi_ledger]get error:", err.Error())
		return nil, errno.GetError.ErrorCode, errno.GetError.ErrorMsg
	}
	return cm, errno.StatusSuccessCode, errno.StatusSuccessMsg
}
