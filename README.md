# Kubeswitch Kind Plugin

This is a store plugin for [kubeswitch](https://github.com/danielfoehrKn/kubeswitch) to get contexts for your local [kind](https://kind.sigs.k8s.io/) clusters.

## Usage

### Build

```bash
go build .
```

### Install

Assuming `~/.local/bin` is in your path:

```bash
cp ks-kind-plugin ~/.local/bin/ks-kind-plugin
```

### Configure

Add the plugin to your kubeswitch configuration, usually in `~/.kube/switch-config.yaml`.
It goes in `kubeconfigStores`. See [switch-config.yaml](switch-config.yaml) for an example.

### Test

You can test just the output of this plugin by running this in the repo:

```bash
kubeswitch --config-path switch-config.yaml
```

