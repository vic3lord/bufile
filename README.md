# Bufile

Bufile is a CLI to generate [Linkerd](https://linkerd.io) service profiles
using [Buf Schema Registry](https://buf.build/docs/registry).

## Usage

Here's an example protobuf file that includes the right annotations:

```proto
 syntax = "proto3"; 
  
 package testbufile.v1; 
  
 import "bufile/v1/bufile.proto"; 
  
 service GreetingService { 
   rpc Say(SayRequest) returns (SayResponse) { 
     option idempotency_level = IDEMPOTENT; 
     option (bufile.v1.linkerd_timeout) = "30s"; 
   } 
 } 
  
 message SayRequest { 
   string name = 1; 
 } 
  
 message SayResponse { 
   string greeting = 1; 
 }
```

```shell
bufile --config <path-to-bufile.json>
```
