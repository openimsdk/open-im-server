package third

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/protocol/third"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func IsNotFound(err error) bool {
	switch errs.Unwrap(err) {
	case redis.Nil, mongo.ErrNoDocuments:
		return true
	default:
		return false
	}
}

func (t *thirdServer) db2pbApplication(val *model.Application) *third.ApplicationVersion {
	return &third.ApplicationVersion{
		Id:         val.ID.Hex(),
		Platform:   val.Platform,
		Version:    val.Version,
		Url:        val.Url,
		Text:       val.Text,
		Force:      val.Force,
		Latest:     val.Latest,
		CreateTime: val.CreateTime.UnixMilli(),
	}
}

func (t *thirdServer) LatestApplicationVersion(ctx context.Context, req *third.LatestApplicationVersionReq) (*third.LatestApplicationVersionResp, error) {
	res, err := t.applicationDatabase.LatestVersion(ctx, req.Platform)
	if err == nil {
		return &third.LatestApplicationVersionResp{Version: t.db2pbApplication(res)}, nil
	} else if IsNotFound(err) {
		return &third.LatestApplicationVersionResp{}, nil
	} else {
		return nil, err
	}
}

func (t *thirdServer) AddApplicationVersion(ctx context.Context, req *third.AddApplicationVersionReq) (*third.AddApplicationVersionResp, error) {
	if err := authverify.CheckAdmin(ctx, t.config.Share.IMAdminUserID); err != nil {
		return nil, err
	}
	val := &model.Application{
		ID:         primitive.NewObjectID(),
		Platform:   req.Platform,
		Version:    req.Version,
		Url:        req.Url,
		Text:       req.Text,
		Force:      req.Force,
		Latest:     req.Latest,
		CreateTime: time.Now(),
	}
	if err := t.applicationDatabase.AddVersion(ctx, val); err != nil {
		return nil, err
	}
	return &third.AddApplicationVersionResp{}, nil
}

func (t *thirdServer) UpdateApplicationVersion(ctx context.Context, req *third.UpdateApplicationVersionReq) (*third.UpdateApplicationVersionResp, error) {
	if err := authverify.CheckAdmin(ctx, t.config.Share.IMAdminUserID); err != nil {
		return nil, err
	}
	oid, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, errs.ErrArgs.WrapMsg("invalid id " + err.Error())
	}
	update := make(map[string]any)
	putUpdate(update, "platform", req.Platform)
	putUpdate(update, "version", req.Version)
	putUpdate(update, "url", req.Url)
	putUpdate(update, "text", req.Text)
	putUpdate(update, "force", req.Force)
	putUpdate(update, "latest", req.Latest)
	if err := t.applicationDatabase.UpdateVersion(ctx, oid, update); err != nil {
		return nil, err
	}
	return &third.UpdateApplicationVersionResp{}, nil
}

func (t *thirdServer) DeleteApplicationVersion(ctx context.Context, req *third.DeleteApplicationVersionReq) (*third.DeleteApplicationVersionResp, error) {
	if err := authverify.CheckAdmin(ctx, t.config.Share.IMAdminUserID); err != nil {
		return nil, err
	}
	ids := make([]primitive.ObjectID, 0, len(req.Id))
	for _, id := range req.Id {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, errs.ErrArgs.WrapMsg("invalid id " + err.Error())
		}
		ids = append(ids, oid)
	}
	if err := t.applicationDatabase.DeleteVersion(ctx, ids); err != nil {
		return nil, err
	}
	return &third.DeleteApplicationVersionResp{}, nil
}

func (t *thirdServer) PageApplicationVersion(ctx context.Context, req *third.PageApplicationVersionReq) (*third.PageApplicationVersionResp, error) {
	total, res, err := t.applicationDatabase.PageVersion(ctx, req.Platform, req.Pagination)
	if err != nil {
		return nil, err
	}
	return &third.PageApplicationVersionResp{
		Total:    total,
		Versions: datautil.Slice(res, t.db2pbApplication),
	}, nil
}
