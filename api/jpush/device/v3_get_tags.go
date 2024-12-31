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

package device

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/calvinit/jiguang-sdk-go/api"
)

// 查询标签列表
//  - 功能说明：获取当前应用的所有标签列表，每个平台最多返回 100 个。
//	- 调用地址：GET `/v3/tags`
//  - 接口文档：https://docs.jiguang.cn/jpush/server/push/rest_api_v3_device#%E6%9F%A5%E8%AF%A2%E6%A0%87%E7%AD%BE%E5%88%97%E8%A1%A8
func (d *apiv3) GetTags(ctx context.Context) (*TagsGetResult, error) {
	if d == nil {
		return nil, api.ErrNilJPushDeviceAPIv3
	}

	req := &api.Request{
		Method: http.MethodGet,
		Proto:  d.proto,
		URL:    d.host + "/v3/tags",
		Auth:   d.auth,
	}
	resp, err := d.client.Request(ctx, req)
	if err != nil {
		return nil, err
	}

	result := &TagsGetResult{Response: resp}
	err = json.Unmarshal(resp.RawBody, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type TagsGetResult struct {
	*api.Response `json:"-"`
	Error         *api.CodeError `json:"error,omitempty"`
	Tags          []string       `json:"tags,omitempty"` // 标签列表
}

func (rs *TagsGetResult) IsSuccess() bool {
	return rs != nil && rs.StatusCode/100 == 2 && rs.Error.IsSuccess()
}