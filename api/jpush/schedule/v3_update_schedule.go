/*
 *
 * Copyright 2024 calvinit/jiguang-sdk-go authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package schedule

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/calvinit/jiguang-sdk-go/api"
)

func (s *apiv3) UpdateSchedule(ctx context.Context, scheduleID string, param *UpdateParam) (*UpdateResult, error) {
	if s == nil {
		return nil, api.ErrNilJPushScheduleAPIv3
	}

	if scheduleID == "" {
		return nil, errors.New("`scheduleID` cannot be empty")
	}
	if param == nil {
		return nil, errors.New("`param` cannot be nil")
	}

	req := &api.Request{
		Method: http.MethodPut,
		Proto:  s.proto,
		URL:    s.host + "/v3/schedules/" + scheduleID,
		Auth:   s.auth,
		Body:   param,
	}
	resp, err := s.client.Request(ctx, req)
	if err != nil {
		return nil, err
	}

	result := &UpdateResult{Response: resp}
	if !resp.IsNoContent() {
		err = json.Unmarshal(resp.RawBody, result)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

type UpdateParam struct {
	Name    string   `json:"name,omitempty"`    // 【可选】任务名称，长度最大 255 字节，数字、字母、下划线、汉字。
	Enabled *bool    `json:"enabled,omitempty"` // 【可选】任务当前状态。
	Trigger *Trigger `json:"trigger,omitempty"` // 【可选】任务触发条件。
	Push    *Push    `json:"push,omitempty"`    // 【可选】任务推送参数。
}

type UpdateResult struct {
	*api.Response   `json:"-"`
	Error           *api.CodeError `json:"error,omitempty"`
	*scheduleDetail `json:"-"`     // 定时任务详情
}

func (rs *UpdateResult) UnmarshalJSON(data []byte) error {
	var aux map[string]json.RawMessage

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if codeError, ok := aux["error"]; ok {
		if err := json.Unmarshal(codeError, &rs.Error); err != nil {
			return err
		}
		delete(aux, "error")
	}

	var detail scheduleDetail
	if err := json.Unmarshal(data, &detail); err != nil {
		return err
	}
	rs.scheduleDetail = &detail

	return nil
}

func (rs UpdateResult) MarshalJSON() ([]byte, error) {
	if rs.Error != nil {
		data := make(map[string]*api.CodeError, 1)
		data["error"] = rs.Error
		return json.Marshal(data)
	}
	return json.Marshal(rs.scheduleDetail)
}

func (rs *UpdateResult) IsSuccess() bool {
	return rs != nil && rs.StatusCode/100 == 2 && rs.Error.IsSuccess()
}
