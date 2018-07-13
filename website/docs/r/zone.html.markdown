---
layout: "cloudflare"
page_title: "Cloudflare: cloudflare_zone"
sidebar_current: "docs-cloudflare-resource-zone"
description: |-
  Provides a Cloudflare Zone resource.
---

# cloudflare_zone

Provides a Cloudflare Zone resource.

## Example Usage

```hcl
resource "cloudflare_zone" "example" {
  domain = "example.com"
  jump_start = true
  organization_id = "01a7362d577a6c3019a474fd6f485823"
  plan = "enterprise"
}

output "nameservers" {
  value = "${cloudflare_zone.example.name_servers}"
}
```

## Argument Reference

* `domain` - (Required) The domain name for the zone.
* `jump_start` - (Optional) Automatically attempt to fetch existing DNS records (true/false).
* `organization_id` - (Optional) The organization id to create the zone under. If this isn't specified it will be created under the users account.
* `plan` - (Optional) The rate plan to use for the zone subscription. Valid values: `free`, `pro`, `business` or `enterprise`.

## Attributes Referenence

The following attributes are exported:

* `id` - Unique identifier in the API for the load balancer.
* `name_servers` - An array of the name servers to set on your domain to utilize Cloudflare.

## Import

Records can be imported by using the Zone ID (found on the "Overview" page of the Cloudflare dashboard.

```
$ terraform import cloudflare_zone.example "506e3185e9c882d175a2d0cb0093d9f2"
```
