/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package command

import (
	"configcenter/src/common/blog"
	"configcenter/src/storage/mongodb"
	"configcenter/src/storage/mongodb/options/aggregateopt"
	"configcenter/src/storage/rpc"
	"configcenter/src/storage/tmserver/core"
	"configcenter/src/storage/types"
)

func init() {
	core.GCommands.SetCommand(types.OPAggregateCode, &aggregate{})
}

var _ core.SetDBProxy = (*aggregate)(nil)

type aggregate struct {
	dbProxy mongodb.Client
}

func (d *aggregate) SetDBProxy(db mongodb.Client) {
	d.dbProxy = db
}
func (d *aggregate) Execute(ctx core.ContextParams, decoder rpc.Request) (*types.OPReply, error) {

	msg := types.OPAggregateOperation{}
	reply := &types.OPReply{}
	if err := decoder.Decode(&msg); nil != err {
		reply.Message = err.Error()
		return reply, err
	}
	blog.V(4).Infof("[MONGO OPERATION] %+v", &msg)

	opt := aggregateopt.One{}
	//opt.Sort = msg.Sort

	err := d.dbProxy.Collection(msg.Collection).AggregateOne(ctx, msg.Pipiline, &opt, &reply.Docs)
	if nil == err {
		reply.Success = true
	} else {
		reply.Message = err.Error()
	}

	return reply, err
}
