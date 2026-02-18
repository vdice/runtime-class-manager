# runtime-class-manager

runtime-class-manager is a Kubernetes operator that manages installation of Wasm shims onto nodes and related Runtimeclasses via [Shim custom resources](https://github.com/spinframework/runtime-class-manager/tree/v{{ CHART_VERSION }}/config/crd/bases/runtime.spinkube.dev_shims.yaml).

## Prerequisites

- [Kubernetes v1.20+](https://kubernetes.io/docs/setup/)
- [Helm v3](https://helm.sh/docs/intro/install/)

## Installing the chart

The following installs the runtime-class-manager chart with the release name `rcm`:

```shell
helm upgrade --install rcm \
  --namespace rcm \
  --create-namespace \
  --version {{ CHART_VERSION }} \
  --wait \
  oci://ghcr.io/spinframework/charts/runtime-class-manager
```

## Post-installation

With runtime-class-manager running, you're ready to create one or more Wasm Shims. See the samples in the [config/samples directory](https://github.com/spinframework/runtime-class-manager/tree/v{{ CHART_VERSION }}/config/samples/).

> Note: Ensure that the `location` for the specified shim binary points to the correct architecture for your Node(s)

For example, here we install the Spin shim on nodes with x86_64 architecture:

```shell
ARCH=x86_64 kubectl apply -f https://raw.githubusercontent.com/spinframework/runtime-class-manager/refs/heads/v{{ CHART_VERSION }}/config/samples/sample_shim_spin_$ARCH.yaml
```

Now when you annotate one or more nodes with a label corresponding to the `nodeSelector` declared in the Shim, runtime-class-manager will install the shim as well as create the corresponding Runtimeclass:

```shell
kubectl label node --all spin=true
```

You are now ready to deploy your Wasm workloads.
