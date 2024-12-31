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

func (s *apiv3) GetSchedule(ctx context.Context, scheduleID string) (*ScheduleGetResult, error) {
	if s == nil {
		return nil, api.ErrNilJPushScheduleAPIv3
	}

	if scheduleID == "" {
		return nil, errors.New("`scheduleID` cannot be empty")
	}

	req := &api.Request{
		Method: http.MethodGet,
		Proto:  s.proto,
		URL:    s.host + "/v3/schedules/" + scheduleID,
		Auth:   s.auth,
	}
	resp, err := s.client.Request(ctx, req)
	if err != nil {
		return nil, err
	}

	result := &ScheduleGetResult{Response: resp}
	err = json.Unmarshal(resp.RawBody, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type ScheduleGetResult struct {
	*api.Response `json:"-"`
	Error         *api.CodeError `json:"error,omitempty"`
	*Schedule     `json:"-"`     // 定时任务详情
}

type scheduleDetail = SendParam

// 定时任务详情。
type Schedule struct {
	ScheduleID string `json:"schedule_id,omitempty"` // 任务 ID
	scheduleDetail
}

func (rs *ScheduleGetResult) UnmarshalJSON(data []byte) error {
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

	var detail Schedule
	if err := json.Unmarshal(data, &detail); err != nil {
		return err
	}
	rs.Schedule = &detail

	return nil
}

func (rs ScheduleGetResult) MarshalJSON() ([]byte, error) {
	if rs.Error != nil {
		data := make(map[string]*api.CodeError, 1)
		data["error"] = rs.Error
		return json.Marshal(data)
	}
	return json.Marshal(rs.scheduleDetail)
}

func (rs *ScheduleGetResult) IsSuccess() bool {
	return rs != nil && rs.StatusCode/100 == 2 && rs.Error.IsSuccess()
}
