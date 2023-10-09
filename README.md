# opg-go-common

[![codecov](https://codecov.io/gh/ministryofjustice/opg-go-common/branch/main/graph/badge.svg?token=BFGR5FBQ0T)](https://codecov.io/gh/ministryofjustice/opg-go-common)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ministryofjustice/opg-go-common)](https://pkg.go.dev/github.com/ministryofjustice/opg-go-common)

Common go packages used by OPG: Managed by opg-org-infra &amp; Terraform

## Using common UI components

Components that include a template element make this available to the consuming service with Golang's [embed directive](https://pkg.go.dev/embed). 
This initialises the content of a file or directory to a variable (`string` or `embed.FS`), which can then be accessed in the usual manner.
The consuming service can import and parse this variable into the template map.

### Pagination

Typical use of the pagination component is as follows, with the position of the component in the id and aria label updated accordingly:

```html
<div class="govuk-grid-row">
  <nav id="bottom-pagination" aria-label="Bottom pagination">
    {{ template "pagination" .Pagination }}
  </nav>
</div>
```