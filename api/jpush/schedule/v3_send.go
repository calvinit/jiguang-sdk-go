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
	"github.com/calvinit/jiguang-sdk-go/api/jpush/push/send"
	"github.com/calvinit/jiguang-sdk-go/jiguang"
)

func (s *apiv3) ScheduleSend(ctx context.Context, param *SendParam) (*SendResult, error) {
	return s.CustomScheduleSend(ctx, param)
}

func (s *apiv3) CustomScheduleSend(ctx context.Context, param interface{}) (*SendResult, error) {
	if s == nil {
		return nil, api.ErrNilJPushScheduleAPIv3
	}

	if param == nil {
		return nil, errors.New("`param` cannot be nil")
	}

	req := &api.Request{
		Method: http.MethodPost,
		Proto:  s.proto,
		URL:    s.host + "/v3/schedules",
		Auth:   s.auth,
		Body:   param,
	}
	resp, err := s.client.Request(ctx, req)
	if err != nil {
		return nil, err
	}

	result := &SendResult{Response: resp}
	err = json.Unmarshal(resp.RawBody, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ↓↓↓ 这是为了方便 SDK 的使用者，提供了一些共享模型的别名定义。↓↓↓

// 任务推送参数。
type Push = send.Param

// ↑↑↑ 这是为了方便 SDK 的使用者，提供了一些共享模型的别名定义。↑↑↑

type SendParam struct {
	CID     string   `json:"cid,omitempty"` // 【可选】用于防止 API 调用端重试造成服务端的重复推送而定义的一个标识符，可通过 GetCidForSchedulePush 接口获取。
	Name    string   `json:"name"`          // 【必填】任务名称，长度最大 255 字节，数字、字母、下划线、汉字。
	Enabled bool     `json:"enabled"`       // 【必填】任务当前状态。
	Trigger *Trigger `json:"trigger"`       // 【必填】任务触发条件。
	Push    *Push    `json:"push"`          // 【必填】任务推送参数。
}

// 任务触发条件。
type Trigger struct {
	Single     *Single     `json:"single,omitempty"`     // 【可选】定时任务，单次触发执行。
	Periodical *Periodical `json:"periodical,omitempty"` // 【可选】定期任务，周期触发执行。
}

// 定时任务，单次触发条件。
type Single struct {
	Time jiguang.LocalDateTime `json:"time"` // 【必填】最晚时间不能超过一年。
}

// 定期任务，周期触发条件。
type Periodical struct {
	StartTime jiguang.LocalDateTime `json:"start"`           // 【必填】有效起始时间。
	EndTime   jiguang.LocalDateTime `json:"end"`             // 【必填】有效结束时间。
	Time      jiguang.LocalTime     `json:"time"`            // 【必填】任务执行时间。
	TimeUnit  jiguang.TimeUnit      `json:"time_unit"`       // 【必填】任务执行最小时间单位，有 jiguang.TimeUnitDay, jiguang.TimeUnitWeek, jiguang.TimeUnitMonth 三种。
	Frequency int                   `json:"frequency"`       // 【必填】任务执行频次，与 TimeUnit 的乘积共同表示的定期任务的执行周期，目前支持的最大值为 100。
	Point     []string              `json:"point,omitempty"` // 【可选】任务执行点，当 TimeUnit 为 jiguang.TimeUnitDay 时，此参数无效。
}

type SendResult struct {
	*api.Response `json:"-"`
	Error         *api.CodeError `json:"error,omitempty"`
	ScheduleID    string         `json:"schedule_id,omitempty"` // 任务 ID。
	Name          string         `json:"name,omitempty"`        // 任务名称。
}

func (rs *SendResult) IsSuccess() bool {
	return rs != nil && rs.StatusCode/100 == 2 && rs.Error.IsSuccess()
}