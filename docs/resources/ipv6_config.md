---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "st-cdnetworks_ipv6_config Resource - st-cdnetworks"
subcategory: ""
description: |-
  Update DNS region IP version, available value: 'V6'
---

# st-cdnetworks_ipv6_config (Resource)

Update DNS region IP version, available value: 'V6'

## Example Usage

```terraform
resource "st-cdnetworks_ipv6_config" "test" {
  domain_id   = "5048000"
  enable_ipv6 = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `domain_id` (String) Domain id
- `enable_ipv6` (Boolean) Ipv6
