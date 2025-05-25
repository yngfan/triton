# Triton-io/Triton

[![License](https://img.shields.io/badge/license-Apache%202-4EB1BA.svg)](https://www.apache.org/licenses/LICENSE-2.0.html)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/2908/badge)](https://bestpractices.coreinfrastructure.org/en/projects/2908)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-v2.0%20adopted-ff69b4.svg)](./CODE_OF_CONDUCT.md)

English | [简体中文](./README-zh_CN.md)

## Introduction

Triton provides a Cloud-Native DeployFlow, which is safe, controllable, and policy-rich.

For more introduction details, see [Overview docs](./docs/README.md)

## Key Features

- **Canary Deploy**

   Triton support `Canary` deploy strategy and provide a Kubernetes service-based traffic pull-in and pull-out solution. 

- **Deploy in batches**

    Triton can deploy your application in several batches and you can pause/resume in deploy process.

- **REST && GRPC API**

   Triton provides many APIs to make deploy easy, such as `Next`, `Cancel`, `Pause`, `Resume`, `Scale`, `Gets`, `Restart` etc. 

- **Selective pods to scale/restart**

  Users can scale or restart selective pods.

- **Base on OpenKruise**

  Triton use [OpenKruise](https://openkruise.io/en-us/docs/what_is_openkruise.html) as workloads which have more powerful capabilities.

- **REST && GRPC API**  
  提供完整的部署生命周期管理接口，主要包括以下能力：

  | 操作类型        | 主要方法                          | 功能描述                                                                 |
  |----------------|----------------------------------|----------------------------------------------------------------------|
  | 部署流程控制    | Next/Pause/Resume/Cancel         | 控制分批次部署流程（推进到下一批次/暂停部署/恢复暂停/终止部署流程）                |
  | 应用扩缩容      | Scale/ScaleIn/ScaleOut           | 支持按指定副本数或百分比进行精确扩缩容，可选择特定 Pod 进行操作                   |
  | 状态查询        | GetDeployStatus/ListDeploys      | 获取实时部署状态、历史部署记录及详情元数据                                   |
  | 异常恢复        | Rollback/Restart                 | 支持快速回滚到历史版本或重启指定批次/全部 Pod                                |
  | 系统管理        | HealthCheck/GetVersion            | 系统健康状态检查及版本信息查询                                            |

- **...**

## Quick Start

For a Kubernetes cluster with its version higher than v1.13, you can simply install Triton with helm v3.1.0+:
```bash
 helm install triton.io https://github.com/triton-io/triton/releases/download/v0.1.1/triton-0.1.1.tgz
```

Note that installing this chart directly means it will use the default template values for the triton-manager.

For more details, see [installation doc](./docs/installation/README.md).

## Documentation

We provide [**tutorials**](./docs/tutorial/README.md) for Triton controllers to demonstrate how to use them.

We also provide [**debug guide**](./docs/debug/README.md) for Triton developers to help to debug. 

## Users


## Contributing

You are warmly welcome to hack on Triton. We have prepared a detailed guide [CONTRIBUTING.md](CONTRIBUTING.md).

## Community

Active communication channels:


## RoadMap
* [ ] Support custom traffic pull-in
* [ ] Provide helm to install DeployFlow
* [ ] REST & GRPC API doc
* [ ] .......

## License

Triton is licensed under the Apache License, Version 2.0. See [LICENSE](./LICENSE.md) for the full license text.


