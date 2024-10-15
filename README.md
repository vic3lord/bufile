# Bufile

Bufile is a CLI to generate [Linkerd](https://linkerd.io) service profiles
using [Buf Schema Registry](https://buf.build/docs/registry).

## Usage

Here's an example protobuf file that includes the right annotations:

https://github.com/vic3lord/bufile/blob/d8f727f2d61a7a3a65064a95683f0fdb7885aadc/proto/testbufile/v1/testbufile.proto#L1-L23

```shell
bufile --config <path-to-bufile.json>
```
