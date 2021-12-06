# krateo

Cross platform commandline tool to manage Krateo Platform.

## Syntax

Most of the commands have flags; you can specify these:

- using the short notation (single dash and single letter; i.e.: `-k`)
- using the long notation (double dash and full flag name; i.e.: `--kubeconfig`)
- by specifying environment variables
- from a config file located in `$HOME/.krateo/krateo.yaml`

### Initializing  Krateo Platform

Usage: **`krateo init [flags]`** where:

| Flag               | Description
|:-------------------|:-------------------------------------|
| `-k, --kubeconfig` | absolute path to the kubeconfig file |
| `-v, --verbose`    | print verbose output                 |

by default `--kubeconfig` flag points to your `$HOME/.kube/config` file.

Example:

```sh
$ krateo init
```

### Installing a specific module

Usage: **`krateo install <MODULE> [flags]`** where:

| Flag               | Description
|:-------------------|:---------------------------------------------------|
| `-k, --kubeconfig` | absolute path to the kubeconfig file               |
| `-t, --repo-token` | token for git repository authentication            |
| `-r, --repo-url`   | url of the git repository where the module resides |
| `-v, --verbose`    | print verbose output                               |

By default `--kubeconfig` flag points to your `$HOME/.kube/config` file.

- you can also use environment variables: `KRATEO_REPO_URL` and `KRATEO_REPO_TOKEN`
- you can also specify these values in the config file: `repo-url: ...`, `repo-token: ...`

Example:

```sh
$ krateo install core 
```

#### Constraints and Requirements for loading modules from your private repositories

Your Krateo module project must reside in a git repository named `krateo-module-NAME`.

> For instance, referring to github as git service, the repository url/name pattern is: 
>
> - `https://github.com/USER/krateo-module-NAME.git`
>
> Example, assuming that: 
>
>   - `USER` = `krateoplatformops`
>   - `NAME` = `core` 
>
>   the module files will be stored in the following repository:
>
>   - `https://github.com/krateoplatformops/krateo-module-core.git`

The project tree must follow this structure:

```txt
.
├── cluster
│   ├── composition.yaml
│   └── definition.yaml
├── crossplane.yaml
├── examples
│   ├── krateo-module-NAME.yaml
│   └── krateo-package-module-NAME.yaml
```

