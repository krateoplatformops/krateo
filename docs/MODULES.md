# Modules

## Constraints and Requirements for loading modules from your private repositories

Your Krateo module project must reside in a git repository named `krateo-module-NAME`.

> For instance, referring to github as git service, the repository url/name pattern is:
>
> - `https://github.com/USER/krateo-module-NAME.git`
>
> Example, assuming that:
>
> - `USER` = `krateoplatformops`
> - `NAME` = `core`
>
> the module files will be stored in the following repository:
>
> - `https://github.com/krateoplatformops/krateo-module-core.git`

## Default repo url

⚠️ If you don't specify the repo url will be used the default one: https://github.com/krateoplatformops/

## Folder structure

The project tree must follow this structure:

```text
.
├── cluster
│   ├── composition.yaml
│   └── definition.yaml
├── crossplane.yaml
├── defaults
│   ├── krateo-module-NAME.yaml
│   └── krateo-package-module-NAME.yaml
├── examples
│   └── krateo-module-NAME.yaml
```

## Configure a specific module

Usage: **`krateo config <MODULE> [flags]`**

## Install a specific module

Usage: **`krateo install <MODULE> [flags]`**

## Flags

| Flag               | Description                                        |
| :----------------- | :------------------------------------------------- |
| `-t, --repo-token` | token for git repository authentication            |
| `-r, --repo-url`   | url of the git repository where the module resides |
| `-v, --verbose`    | print verbose output                               |

- you can also use environment variables: `KRATEO_REPO_URL` and `KRATEO_REPO_TOKEN`
- you can also specify these values in the config file: `repo-url: ...`, `repo-token: ...`

Example:

```sh
$ krateo install core
```
