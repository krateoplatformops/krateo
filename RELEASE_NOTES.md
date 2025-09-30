## Release 2.5.1

## Removed Charts
- portal v0.1.10: Removed

## eventrouter v0.5.5
### What's Changed


**Full Changelog**: https://github.com/krateoplatformops/eventrouter/compare/0.5.5...0.5.5


## resource-tree-handler v0.3.0
### What's Changed


**Full Changelog**: https://github.com/krateoplatformops/resource-tree-handler/compare/0.3.0...0.3.0


## oasgen-provider v0.6.0
### What's Changed

### ✨ Features
- feat: added additionalstatusfields in RestDefinition ([link](https://github.com/krateoplatformops/oasgen-provider/pull/72)) by @vicentinileonardo
- feat: make the status field of resources in oasgen type safe ([link](https://github.com/krateoplatformops/oasgen-provider/pull/74)) by @vicentinileonardo
- feat: oas2jsonschema code section refactor + Configuration CRD ([link](https://github.com/krateoplatformops/oasgen-provider/pull/76)) by @vicentinileonardo

### 📚 Documentation
- docs: readme update - added informations about pending state ([link](https://github.com/krateoplatformops/oasgen-provider/pull/69)) by @matteogastaldello

### 🔧 Other Changes
- ci: fix broken step ([link](https://github.com/krateoplatformops/oasgen-provider/pull/67)) by @matteogastaldello


**Full Changelog**: https://github.com/krateoplatformops/oasgen-provider/compare/0.5.3...0.6.0


## snowplow v0.19.0
### What's Changed

### ✨ Features
- feat: API call should proxy also user defined content type during POS… ([link](https://github.com/krateoplatformops/snowplow/pull/123)) by @lucasepe
- feat: Add precooked payload for PATCH in resourcerefs response ([link](https://github.com/krateoplatformops/snowplow/pull/125)) by @lucasepe
- feat: handle custom JQ modules from local path KRA-480 ([link](https://github.com/krateoplatformops/snowplow/pull/127)) by @lucasepe
- feat: fix fetching spec.resourceRefs after updating to map ([link](https://github.com/krateoplatformops/snowplow/pull/131)) by @lucasepe
- feat: bump plumbing dependency ([link](https://github.com/krateoplatformops/snowplow/pull/135)) by @lucasepe
- feat: expose JQ engine via endpoint KRA-705 ([link](https://github.com/krateoplatformops/snowplow/pull/137)) by @lucasepe
- feat: add traceId to widget status KRA-709 ([link](https://github.com/krateoplatformops/snowplow/pull/139)) by @lucasepe
- feat: inject variables in the RESTAction context at call time ([link](https://github.com/krateoplatformops/snowplow/pull/141)) by @lucasepe
- feat: align widgetdata crd validation with kube apiserver strict mode ([link](https://github.com/krateoplatformops/snowplow/pull/143)) by @lucasepe
- feat: change how resourcerefs are returned to fe KRA-707 ([link](https://github.com/krateoplatformops/snowplow/pull/147)) by @lucasepe
- feat: expose pagination data to restaction filter ([link](https://github.com/krateoplatformops/snowplow/pull/149)) by @lucasepe
- feat: bump plumbing deps to use http client retry wrapper ([link](https://github.com/krateoplatformops/snowplow/pull/152)) by @lucasepe

### 🐛 Bug Fixes
- fix: allow patch verb in build uri function ([link](https://github.com/krateoplatformops/snowplow/pull/133)) by @lucasepe
- fix: handle preserve unknow fields ([link](https://github.com/krateoplatformops/snowplow/pull/145)) by @lucasepe

### 🔧 Other Changes
- 128 restactions and widgets pagination ([link](https://github.com/krateoplatformops/snowplow/pull/129)) by @lucasepe


**Full Changelog**: https://github.com/krateoplatformops/snowplow/compare/0.12.2...0.19.0


## cratedb-chart v0.1.4
### What's Changed

### ✨ Features
- feat: added default storageClassName selection + minor general fixes ([link](https://github.com/krateoplatformops/cratedb-chart/pull/10)) by @vicentinileonardo


**Full Changelog**: https://github.com/krateoplatformops/cratedb-chart/compare/0.1.3...0.1.4


## finops-composition-definition-parser v0.1.2
### What's Changed


**Full Changelog**: https://github.com/krateoplatformops/finops-composition-definition-parser/compare/0.1.2...0.1.2


## core-provider v0.25.2
### What's Changed

### ✨ Features
- feat: add path field on composition status ([link](https://github.com/krateoplatformops/core-provider/pull/182)) by @matteogastaldello

### 🔧 Other Changes
- chore: added resources field on status ([link](https://github.com/krateoplatformops/core-provider/pull/172)) by @matteogastaldello
- 173 add support for service for composition dynamic controller ([link](https://github.com/krateoplatformops/core-provider/pull/174)) by @matteogastaldello
- KRA-763 - change logging library ([link](https://github.com/krateoplatformops/core-provider/pull/176)) by @matteogastaldello
- chore: lib update ([link](https://github.com/krateoplatformops/core-provider/pull/178)) by @matteogastaldello
- 179 libs update logs revision ([link](https://github.com/krateoplatformops/core-provider/pull/180)) by @matteogastaldello


**Full Changelog**: https://github.com/krateoplatformops/core-provider/compare/0.24.7...0.25.2


## opa-chart v0.1.0
### What's Changed

### 🐛 Bug Fixes
- fix: service port name ([link](https://github.com/krateoplatformops/opa-chart/pull/4)) by @FrancescoL96

### 🔧 Other Changes
- refactor: value file now follows Krateo's standard ([link](https://github.com/krateoplatformops/opa-chart/pull/2)) by @FrancescoL96
- Let clusterrole and clusterrolebinding depend on the release namespace ([link](https://github.com/krateoplatformops/opa-chart/pull/6)) by @braghettos


**Full Changelog**: https://github.com/krateoplatformops/opa-chart/commits/0.1.0


## frontend v0.0.48
### What's Changed

### ✨ Features
- feat: add declarative success and error messages in action rest ([link](https://github.com/krateoplatformops/frontend/pull/5)) by @fedepini
- feat: update resources refs handling ([link](https://github.com/krateoplatformops/frontend/pull/7)) by @braghettos

### 🔧 Other Changes
- refactor: updated getEndpointUrl function and error handling ([link](https://github.com/krateoplatformops/frontend/pull/6)) by @fedepini


**Full Changelog**: https://github.com/krateoplatformops/frontend/compare/0.0.19...0.0.48


## finops-operator-exporter v0.4.3
### What's Changed


**Full Changelog**: https://github.com/krateoplatformops/finops-operator-exporter/compare/0.4.3...0.4.3


## finops-database-handler-uploader v0.1.0
### What's Changed

### 🐛 Bug Fixes
- fix: error in configuration parsing for dbhandler URL override ([link](https://github.com/krateoplatformops/finops-database-handler-uploader/pull/4)) by @FrancescoL96

### 🔧 Other Changes
- test: added tests for update and delete functions of the controller ([link](https://github.com/krateoplatformops/finops-database-handler-uploader/pull/2)) by @FrancescoL96


**Full Changelog**: https://github.com/krateoplatformops/finops-database-handler-uploader/commits/0.1.0


## finops-moving-window-policy-chart v0.1.2
### What's Changed

### 🔧 Other Changes
- refactor: updated example composition kind ([link](https://github.com/krateoplatformops/finops-moving-window-policy-chart/pull/8)) by @FrancescoL96


**Full Changelog**: https://github.com/krateoplatformops/finops-moving-window-policy-chart/compare/0.1.1...0.1.2


## finops-operator-focus v0.4.3
### What's Changed


**Full Changelog**: https://github.com/krateoplatformops/finops-operator-focus/compare/0.4.3...0.4.3


## authn v0.20.1
### What's Changed


**Full Changelog**: https://github.com/krateoplatformops/authn/compare/0.20.1...0.20.1


## finops-database-handler v0.4.5
### What's Changed

### ✨ Features
- feat: added new notebook code api ([link](https://github.com/krateoplatformops/finops-database-handler/pull/41)) by @FrancescoL96

### 🐛 Bug Fixes
- fix: added try catch on list notebooks query ([link](https://github.com/krateoplatformops/finops-database-handler/pull/43)) by @FrancescoL96


**Full Changelog**: https://github.com/krateoplatformops/finops-database-handler/compare/0.4.4...0.4.5


## etcd-chart v11.1.3
### What's Changed


**Full Changelog**: https://github.com/krateoplatformops/etcd-chart/compare/11.1.3...11.1.3


## smithery v0.8.0
### What's Changed

### ✨ Features
- feat: update resourcesRefs schema to handle slicing ([link](https://github.com/krateoplatformops/smithery/pull/26)) by @lucasepe
- feat: add resource enum key ([link](https://github.com/krateoplatformops/smithery/pull/28)) by @lucasepe


**Full Changelog**: https://github.com/krateoplatformops/smithery/compare/0.6.1...0.8.0


## finops-notebooks-chart v0.1.1
### What's Changed

### 🐛 Bug Fixes
- fix: deletion job checks for notebooks before deleting endpoint ([link](https://github.com/krateoplatformops/finops-notebooks-chart/pull/2)) by @FrancescoL96


**Full Changelog**: https://github.com/krateoplatformops/finops-notebooks-chart/compare/0.1.0...0.1.1


## eventsse v0.5.3
### What's Changed


**Full Changelog**: https://github.com/krateoplatformops/eventsse/compare/0.5.3...0.5.3


## finops-operator-scraper v0.4.2
### What's Changed


**Full Changelog**: https://github.com/krateoplatformops/finops-operator-scraper/compare/0.4.2...0.4.2


## portal v1.0.0
### What's Changed

### ✨ Features
- feat: improve dashboard ([link](https://github.com/krateoplatformops-blueprints/portal/pull/5)) by @braghettos
- feat: adopt allowedResouces array ([link](https://github.com/krateoplatformops-blueprints/portal/pull/7)) by @braghettos
- feat: adjust adoption of allowedResources array ([link](https://github.com/krateoplatformops-blueprints/portal/pull/9)) by @braghettos
- feat: update cyberjoker RBAC ([link](https://github.com/krateoplatformops-blueprints/portal/pull/15)) by @braghettos

### 🐛 Bug Fixes
- fix: update rbac for cyberjoker user ([link](https://github.com/krateoplatformops-blueprints/portal/pull/13)) by @braghettos

### 📚 Documentation
- docs: update README.md ([link](https://github.com/krateoplatformops-blueprints/portal/pull/11)) by @braghettos


**Full Changelog**: https://github.com/krateoplatformops-blueprints/portal/compare/0.0.1...1.0.0






## Release 2.6.0

## Removed Charts
Nothing removed

## core-provider v0.25.2
### What's Changed


**Full Changelog**: https://github.com/krateoplatformops/core-provider/compare/0.25.2...0.25.2


## finops-composition-definition-parser v0.1.2
### What's Changed


**Full Changelog**: https://github.com/krateoplatformops/finops-composition-definition-parser/compare/0.1.2...0.1.2


## finops-operator-focus v0.4.3
### What's Changed


**Full Changelog**: https://github.com/krateoplatformops/finops-operator-focus/compare/0.4.3...0.4.3


## finops-database-handler v0.4.5
### What's Changed


**Full Changelog**: https://github.com/krateoplatformops/finops-database-handler/compare/0.4.5...0.4.5


## etcd-chart v11.1.3
### What's Changed


**Full Changelog**: https://github.com/krateoplatformops/etcd-chart/compare/11.1.3...11.1.3


## finops-notebooks-chart v0.1.1
### What's Changed


**Full Changelog**: https://github.com/krateoplatformops/finops-notebooks-chart/compare/0.1.1...0.1.1


## opa-chart v0.1.0
### What's Changed

### 🐛 Bug Fixes
- fix: service port name ([link](https://github.com/krateoplatformops/opa-chart/pull/4)) by @FrancescoL96

### 🔧 Other Changes
- refactor: value file now follows Krateo's standard ([link](https://github.com/krateoplatformops/opa-chart/pull/2)) by @FrancescoL96
- Let clusterrole and clusterrolebinding depend on the release namespace ([link](https://github.com/krateoplatformops/opa-chart/pull/6)) by @braghettos


**Full Changelog**: https://github.com/krateoplatformops/opa-chart/commits/0.1.0


## smithery v0.10.1
### What's Changed

### ✨ Features
- feat: upgrade crdgen dependency to v0.5.0 ([link](https://github.com/krateoplatformops/smithery/pull/30)) by @lucasepe
- feat: upgrade go version in dockerfile ([link](https://github.com/krateoplatformops/smithery/pull/32)) by @lucasepe
- feat: add extra slice properties for pagination ([link](https://github.com/krateoplatformops/smithery/pull/34)) by @lucasepe

### 🔧 Other Changes
- 35 use shared crdgen from jsonschema package ([link](https://github.com/krateoplatformops/smithery/pull/36)) by @lucasepe


**Full Changelog**: https://github.com/krateoplatformops/smithery/compare/0.8.0...0.10.1


## resource-tree-handler v0.3.0
### What's Changed


**Full Changelog**: https://github.com/krateoplatformops/resource-tree-handler/compare/0.3.0...0.3.0


## finops-operator-exporter v0.5.0
### What's Changed

### ✨ Features
- feat: added support for single exporter ([link](https://github.com/krateoplatformops/finops-operator-exporter/pull/88)) by @FrancescoL96


**Full Changelog**: https://github.com/krateoplatformops/finops-operator-exporter/compare/0.4.3...0.5.0


## frontend v0.0.57
### What's Changed

### 🔧 Other Changes
- ci: update smithery version to 0.9.1 ([link](https://github.com/krateoplatformops/frontend/pull/9)) by @braghettos
- ci: move crds to dedicate chart ([link](https://github.com/krateoplatformops/frontend/pull/12)) by @braghettos
- Widget pagination ([link](https://github.com/krateoplatformops/frontend/pull/10)) by @kandros


**Full Changelog**: https://github.com/krateoplatformops/frontend/compare/0.0.48...0.0.57


## oasgen-provider v0.7.0
### What's Changed

### ✨ Features
- feat: added recursion guard + fixes ([link](https://github.com/krateoplatformops/oasgen-provider/pull/82)) by @vicentinileonardo

### 📚 Documentation
- docs: updated readme ([link](https://github.com/krateoplatformops/oasgen-provider/pull/78)) by @vicentinileonardo
- docs: updated readme ([link](https://github.com/krateoplatformops/oasgen-provider/pull/80)) by @vicentinileonardo


**Full Changelog**: https://github.com/krateoplatformops/oasgen-provider/compare/0.6.0...0.7.0


## snowplow v0.20.2
### What's Changed

### ✨ Features
- feat: upgrade plumbing deps to use retry http client ([link](https://github.com/krateoplatformops/snowplow/pull/154)) by @lucasepe
- feat: upgrade swagger and replace all perPage params ([link](https://github.com/krateoplatformops/snowplow/pull/158)) by @lucasepe

### 🐛 Bug Fixes
- fix: replace perpage with perPage in resourcesref url ([link](https://github.com/krateoplatformops/snowplow/pull/160)) by @lucasepe

### 🔧 Other Changes
- 155 pagination ([link](https://github.com/krateoplatformops/snowplow/pull/156)) by @lucasepe


**Full Changelog**: https://github.com/krateoplatformops/snowplow/compare/0.19.0...0.20.2


## cratedb-chart v0.1.4
### What's Changed

### ✨ Features
- feat: added default storageClassName selection + minor general fixes ([link](https://github.com/krateoplatformops/cratedb-chart/pull/10)) by @vicentinileonardo


**Full Changelog**: https://github.com/krateoplatformops/cratedb-chart/compare/0.1.3...0.1.4


## authn v0.20.1
### What's Changed


**Full Changelog**: https://github.com/krateoplatformops/authn/compare/0.20.1...0.20.1


## eventsse v0.5.3
### What's Changed


**Full Changelog**: https://github.com/krateoplatformops/eventsse/compare/0.5.3...0.5.3


## finops-database-handler-uploader v0.1.0
### What's Changed


**Full Changelog**: https://github.com/krateoplatformops/finops-database-handler-uploader/compare/0.1.0...0.1.0


## portal v1.1.0
### What's Changed

### ✨ Features
- feat: remove widgets presentation logic from restaction ([link](https://github.com/krateoplatformops-blueprints/portal/pull/18)) by @braghettos


**Full Changelog**: https://github.com/krateoplatformops-blueprints/portal/compare/1.0.0...1.1.0


## finops-moving-window-policy-chart v0.1.2
### What's Changed


**Full Changelog**: https://github.com/krateoplatformops/finops-moving-window-policy-chart/compare/0.1.2...0.1.2


## finops-operator-scraper v0.4.2
### What's Changed


**Full Changelog**: https://github.com/krateoplatformops/finops-operator-scraper/compare/0.4.2...0.4.2


## eventrouter v0.5.5
### What's Changed


**Full Changelog**: https://github.com/krateoplatformops/eventrouter/compare/0.5.5...0.5.5




