/*
 *
 * Copyright 2025 cavlabs/jiguang-sdk-go authors.
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

package jums

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/cavlabs/jiguang-sdk-go/api"
	"github.com/cavlabs/jiguang-sdk-go/api/jums/message"
)

// # 模板消息 - 广播发送
//   - 功能说明：模板消息广播发送。
//   - 调用地址：POST `/v1/template/broadcast`
//   - 接口文档：[docs.jiguang.cn]
//
// [docs.jiguang.cn]: https://docs.jiguang.cn/jums/server/rest_api_jums_template_message
func (u *apiv1) TemplateBroadcastSend(ctx context.Context, param *TemplateBroadcastSendParam) (*TemplateBroadcastSendResult, error) {
	if u == nil {
		return nil, api.ErrNilJUmsAPIv1
	}

	if param == nil {
		return nil, errors.New("`param` cannot be nil")
	}

	req := &api.Request{
		Method: http.MethodPost,
		Proto:  u.proto,
		URL:    u.host + "/v1/template/broadcast",
		Auth:   u.auth,
		Body:   param,
	}
	resp, err := u.client.Request(ctx, req)
	if err != nil {
		return nil, err
	}

	result := &TemplateBroadcastSendResult{Response: resp}
	err = json.Unmarshal(resp.RawBody, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type TemplateBroadcastSendParam struct {
	// 【必填】模板 ID。
	TemplateID int64 `json:"template_id"`
	// 【可选】模板参数，需要替换的参数名和参数值的键值对。
	TemplateParams map[string]interface{} `json:"template_para,omitempty"`
	// 【可选】APP 通道的相关参数，模板中有 APP 通道时必填。
	AppParams *message.AppParams `json:"app_para,omitempty"`
	// 【可选】发送策略 ID，如果是同时发送可传 0 或不传。当使用自定义通道 ID 发送时，该字段无效。
	//
	// 1. 如不需要进行补发，仅进行单通道、多通道同时发送，则不需填写策略 ID，或者设置为 0。
	//
	// 2. 在 官网控制台-渠道-发送策略 中创建一个补发策略后，调 API 时可使用策略 ID 进行指定。
	RuleID int `json:"rule_id,omitempty"`
	// 【可选】可选参数，用于黑白名单 ID、提交人等信息的填写。
	Option *message.Option `json:"option,omitempty"`
	// 【可选】回调参数。
	//
	// 调 API 发送消息时，可以指定 Callback 参数，方便用户临时变更回调 URL 或者回调带上其自定义参数，满足其日常业务需求。详细使用说明请阅读 [消息回调设置]。
	//
	// 此功能仅针对极光 VIP 用户提供，提供「目标有效/无效、提交成功/失败、送达成功/失败、点击、撤回成功/失败」9 种消息状态，需在官网控制台设置所需回调的状态。
	//
	// [消息回调设置]: https://docs.jiguang.cn/jums/advanced/callback
	Callback *message.Callback `json:"callback,omitempty"`
}

type TemplateBroadcastSendResult = TemplateSendResult
