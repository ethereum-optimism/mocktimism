# Mocktimism contributing guide

Add [Mocktimism] specifications.

[Mocktimism TODO]: https://github.com/ethereum-optimism/optimism/blob/develop/specs/superchain-upgrades.md

# Mocktimism contributing guide

🎈 Thanks for your help improving the project! We are so happy to have you!

**No contribution is too small and all contributions are valued.**

Mocktimism is a WIP and thus there are plenty of ways to contribute, in particular we appreciate support in the following areas:

- Implementing issues. You can start off with those tagged ["good first issue"](https://github.com/ethereum-optimism/mocktimism/contribute) which are meant as introductory issues for external contributors.
- Reporting issues. For security issues see [Security policy](https://github.com/ethereum-optimism/.github/blob/master/SECURITY.md).

Note that we have a [Code of Conduct](https://github.com/ethereum-optimism/.github/blob/master/CODE_OF_CONDUCT.md), please follow it in all your interactions with the project.

## Workflow for Pull Requests

🚨 Before making any non-trivial change, please first respond to or open an issue describing the change to solicit feedback and guidance. This will increase the likelihood of the PR getting merged.

In general, the smaller the diff the easier it will be for us to review quickly.

In order to contribute, fork the appropriate branch, for non-breaking changes to production that is `develop`. 

Additionally, if you are writing a new feature, please ensure you add appropriate test cases.

We recommend using the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) format on commit messages.

Unless your PR is ready for immediate review and merging, please mark it as 'draft' (or simply do not open a PR yet).

**Bonus:** Add comments to the diff under the "Files Changed" tab on the PR page to clarify any sections where you think we might have questions about the approach taken.

### Response time:
We aim to provide a meaningful response to all PRs and issues from external contributors within 2 business days.

## Development Quick Start

### Dependencies

You'll need the following:

* [Git](https://git-scm.com/downloads)
* [Docker](https://docs.docker.com/get-docker/)
* [Docker Compose](https://docs.docker.com/compose/install/)
* [Go](https://go.dev/dl/)
* [Foundry](https://getfoundry.sh)

### Setup

Clone the repository and open it:

```bash
git clone git@github.com:ethereum-optimism/mocktimism.git
cd mocktimism
```

### Running tests

```
make test
```

### Running build

```
make build
```

Or build with docker

```
make docker
```

### Running linters

Go tidy to clean up go dependencies

```
make tidy
```

Linter with autofix

```
make lint
```

### Writing Docs

User docs are in [./docs](docs/).
