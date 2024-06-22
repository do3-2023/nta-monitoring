# ☸️ Monitoring Exercise

This repository contains the code for an assignment as part of the course Monitoring at Polytech Montpellier.

## Table of contents

1. [Preface](#preface)
2. [Setup](#setup)
3. [WASM](#wasm)
4. [License](#license)

## Preface

This assignment follows the previous one entitled Orchestration exercise. Please look at this first.

## Setup

This assignment is based on the previous one. You can find the code in the [orchestration-exercise](ihttps://github.com/do3-2023/nta-kube).

As my kubernetes distribution, I used RKE2. You can find the installation instructions [here](https://rke2.io/).

My setup consists of 2 nodes, one master and one worker. Both of them are running on local debian VMs.

I used Cilium as my CNI plugin. You can find the installation instructions [here](https://docs.cilium.io/en/v1.10/gettingstarted/k8s-install-default/).

The config I used for this cluster is in infra/cilium.yaml.

For PV and PVC, I used the local-path provisioner. You can find the installation instructions [here](https://github.com/rancher/local-path-provisioner).

## WASM

I made a second version of the app using WebAssembly. You can find the code in the [wasm branch](https://github.com/do3-2023/nta-monitoring/tree/spinkube).

This version has been made using Spin and TinyGo.

To build the WebAssembly file, you can use the following command:

```bash
spin build
```

The app is then deployed on k8s using Spinkube. You can find the installation instructions [here](https://www.spinkube.dev/).

To access the app, I port-forwarded the service to my local machine:

```bash
kubectl port-forward svc/nta-kube-api 3000:3000
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
