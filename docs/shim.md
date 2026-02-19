## Shim

Runtime Class Manager operates on `Shim` custom resources based on the [Shim CRD](../config/crd/bases/runtime.spinkube.dev_shims.yaml).

Whenever a Shim is created, updated or deleted, Runtime-Class-Manager will perform the necessary actions, e.g. creating, updating or removing the associated [RuntimeClass](./runtimeclass.md), installing or removing shim binaries on any [Nodes](https://kubernetes.io/docs/concepts/architecture/nodes/) labeled with a corresponding to its `nodeSelector`, etc.

### Configuration

For full, detailed configuration options, see the [Shim CRD](../config/crd/bases/runtime.spinkube.dev_shims.yaml). Here we point out a few pertinent items.

* `spec.nodeSelector`: The label key and value applied to Nodes where this particular shim should be installed
* `spec.fetchStrategy`: The strategy for fetching the shim binary
  * `spec.fetchStrategy.anonHttp`: Fetch the shim binary from a specified URL. This is the legacy option.
  * `spec.fetchStrategy.platforms`: A list of per-OS/architecture artifact entries. Each entry specifies `os`, `arch`, `location`, and an optional `sha256` digest. The controller selects the matching entry for each target node. This is the current recommended strategy.
* `spec.containerdRuntimeOptions`: Options specific to the shim that should be added to the containerd configuration

### Operation

You may observe the "install" and "uninstall" [Jobs](https://kubernetes.io/docs/concepts/workloads/controllers/job/) responsible for downloading and installing (or uninstalling) the shim binary. These will run on every Node that matches the Shim's `nodeSelector`.
