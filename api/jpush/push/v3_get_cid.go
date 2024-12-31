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

package push

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/calvinit/jiguang-sdk-go/api"
)

// 获取推送唯一标识 (CID)
//  - 功能说明：CID 是用于防止 API 调用端重试造成服务端的重复推送而定义的一个推送参数。用户使用一个 CID 推送后，再次使用相同的 CID 进行推送，则会直接返回第一次成功推送的结果，不会再次进行推送。
//  CID 的有效期为 1 天，格式为：{appkey}-{uuid}，在使用 CID 之前，必须通过接口获取你的 CID 池。
//	- 调用地址：GET `/v3/push/cid?type=push&count={count}`，如 count < 1 则自动重置为 1。
//  - 接口文档：https://docs.jiguang.cn/jpush/server/push/rest_api_v3_push_advanced#%E8%8E%B7%E5%8F%96%E6%8E%A8%E9%80%81%E5%94%AF%E4%B8%80%E6%A0%87%E8%AF%86cid
func (p *apiv3) GetCidForPush(ctx context.Context, count int) (*CidGetResult, error) {
	if p == nil {
		return nil, api.ErrNilJPushPushAPIv3
	}

	if count < 1 {
		count = 1
	}

	req := &api.Request{
		Method: http.MethodGet,
		Proto:  p.proto,
		URL:    p.host + "/v3/push/cid?type=push&count=" + strconv.Itoa(count),
		Auth:   p.auth,
	}
	resp, err := p.client.Request(ctx, req)
	if err != nil {
		return nil, err
	}

	result := &CidGetResult{Response: resp}
	err = json.Unmarshal(resp.RawBody, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
