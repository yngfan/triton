/*
Copyright 2021 The Triton Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"github.com/triton-io/triton/pkg/kube/controller/cloneset"
	"github.com/triton-io/triton/pkg/kube/controller/deployflow"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

/*
*
切片用于存储所有控制器的初始化函数。
每个控制器都有一个初始化函数，该函数接受一个 manager.Manager 作为参数，并返回一个 error。
这个切片用于存储所有控制器的初始化函数，以便在主函数中调用它们。
主函数会遍历这个切片，调用每个控制器的初始化函数，将控制器添加到管理器中。
如果控制器的 CRD 未安装，则会跳过该控制器的初始化。
如果控制器的 CRD 安装失败，则会返回错误。
如果控制器的初始化函数返回错误，则会跳过该控制器的初始化。
如果控制器的初始化函数返回 nil，则会将控制器添加到管理器中。
*/
var controllerAddFuncs []func(manager.Manager) error

func init() {
	// 将 deployflow 控制器的 Add 方法注册到控制器列表
	controllerAddFuncs = append(controllerAddFuncs, deployflow.Add)
	// 将 cloneset 控制器的 Add 方法注册到控制器列表
	controllerAddFuncs = append(controllerAddFuncs, cloneset.Add)
}

func SetupWithManager(m manager.Manager) error {
	for _, f := range controllerAddFuncs {
		if err := f(m); err != nil {
			if kindMatchErr, ok := err.(*meta.NoKindMatchError); ok {
				klog.Infof("CRD %v is not installed, its controller will perform noops!", kindMatchErr.GroupKind)
				continue
			}
			return err
		}
	}
	return nil
}
