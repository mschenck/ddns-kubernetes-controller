# ddns-kubernetes-controller

[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/ddns-kubernetes-controller)](https://artifacthub.io/packages/search?repo=ddns-kubernetes-controller)

## Overview

Dynamic DNS (kubernetes controller) will manage any number of DNS records and update them to resolve to your external IP.

This controller introduces a `DdnsRecord` custom resource.  Each resource will configure:
1. The Record (DNS record) to dynamically manage
2. The Zone (DNS Zone) to manage a Record for
3. the Provider (which the given Zone is served from)
4. The TTL of the record (and the interval in which it is updated)

The controller will, watch each `DdnsRecord` and at the interval (of the `DdnsRecord.TTL`) will:
1. Perform an IP lookup (currently via ipify.com) to determine the public IP of _this_ cluster
2. Perform a DNS lookup for the `DdnsRecord`'s `Record.Zone` to see what it resolves to.
3. If they don't match, the Record is updated to _this_ cluster's external IP address.

## Configuration

- Each DNS record is configured via the `DdnsRecord` custom resource
- Each Provider is configured via the `ddns-config` secret

### DdnsRecord

These are the custom resources for managing all of the DNS records you want to be mapped to _this_ cluster's public IP (i.e. NAT gateway address).

See [sample ddns_v1_ddnsrecord](./config/samples/ddns_v1_ddnsrecord.yaml)

The spec attributes are (for the given "foo.example.com"):
- `record`: The Zone record entry you wish to manage, i.e. "foo"
- `zone`: The DNS Zone, i.e. "example.com"
- `ttl`: The TTL for the zone record (and the refresh interval for the controller), i.e. "30s"
- `provider`: The provider (vendor) hosting the Zone, i.e. "aws"

### ddns-config

This is how secret credentials your zone(s) provider(s) are configured.  The format is:

```
<vendor>
  <config key1>: <config value1>
  <config key2>: <config value2>
# I.e.
aws:
  AWS_ACCESS_KEY_ID: <access key ID>
  AWS_SECRET_ACCESS_KEY: <secret accesss key>
```

## Helm

Add repository:

        helm repo add ddns-kubernetes-controller https://mschenck.github.io/ddns-kubernetes-controller

Install chart:

        helm install ddns-kubernetes-controller ddns-kubernetes-controller/ddns-kubernetes-controller

