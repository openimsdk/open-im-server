package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/openimsdk/protocol/auth"
	"github.com/openimsdk/protocol/third"
	"github.com/openimsdk/tools/errs"
)

type Api struct {
	Api    string
	UserID string
	Secret string
	Token  string
	Client *http.Client
}

func (a *Api) apiPost(ctx context.Context, path string, req any, resp any) error {
	operationID, _ := ctx.Value("operationID").(string)
	if operationID == "" {
		return errs.New("call api operationID is empty")
	}
	reqBody, err := json.Marshal(req)
	if err != nil {
		return err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, a.Api+path, bytes.NewReader(reqBody))
	if err != nil {
		return err
	}
	DefaultRequestHeader(request.Header)
	request.ContentLength = int64(len(reqBody))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("operationID", operationID)
	if a.Token != "" {
		request.Header.Set("token", a.Token)
	}
	response, err := a.Client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("api %s status %s body %s", path, response.Status, body)
	}
	var baseResponse struct {
		ErrCode int             `json:"errCode"`
		ErrMsg  string          `json:"errMsg"`
		ErrDlt  string          `json:"errDlt"`
		Data    json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &baseResponse); err != nil {
		return err
	}
	if baseResponse.ErrCode != 0 {
		return fmt.Errorf("api %s errCode %d errMsg %s errDlt %s", path, baseResponse.ErrCode, baseResponse.ErrMsg, baseResponse.ErrDlt)
	}
	if resp != nil {
		if err := json.Unmarshal(baseResponse.Data, resp); err != nil {
			return err
		}
	}
	return nil
}

func (a *Api) GetAdminToken(ctx context.Context) (string, error) {
	req := auth.GetAdminTokenReq{
		UserID: a.UserID,
		Secret: a.Secret,
	}
	var resp auth.GetAdminTokenResp
	if err := a.apiPost(ctx, "/auth/get_admin_token", &req, &resp); err != nil {
		return "", err
	}
	return resp.Token, nil
}

func (a *Api) GetPartLimit(ctx context.Context) (*third.PartLimitResp, error) {
	var resp third.PartLimitResp
	if err := a.apiPost(ctx, "/object/part_limit", &third.PartLimitReq{}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (a *Api) InitiateMultipartUpload(ctx context.Context, req *third.InitiateMultipartUploadReq) (*third.InitiateMultipartUploadResp, error) {
	var resp third.InitiateMultipartUploadResp
	if err := a.apiPost(ctx, "/object/initiate_multipart_upload", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (a *Api) CompleteMultipartUpload(ctx context.Context, req *third.CompleteMultipartUploadReq) (string, error) {
	var resp third.CompleteMultipartUploadResp
	if err := a.apiPost(ctx, "/object/complete_multipart_upload", req, &resp); err != nil {
		return "", err
	}
	return resp.Url, nil
}
