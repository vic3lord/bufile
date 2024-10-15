# Bufile

Bufile is a CLI to generate [Linkerd](https://linkerd.io) service profiles
using [Buf Schema Registry](https://buf.build/docs/registry).

## Usage

Here's an example protobuf file that includes the right annotations:

https://github.com/vic3lord/bufile/blob/main/proto/testbufile/v1/testbufile.proto

```shell
bufile --config <path-to-bufile.json>
```
