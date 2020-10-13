# How to Contribute

app-netutil is [Apache 2.0 licensed](LICENSE) and accepts contributions via GitHub
pull requests. This document outlines some of the conventions on development
workflow, commit message formatting, contact points and other resources to make
it easier to get your contribution accepted.

Questions or issues can be raised by opening up a GitHub issue.

## Coding Style

Please follows the standard formatting recommendations and language idioms set out
in [Effective Go](https://golang.org/doc/effective_go.html) and in the
[Go Code Review Comments wiki](https://github.com/golang/go/wiki/CodeReviewComments).

All pull requests should have properly formatted code by running the code base
through `gofmt` and `golangci-lint` before submitting.

To install tools locally:
* [golangci-lint](https://golangci-lint.run/usage/install/#local-installation)

To run locally:

```bash
cd $GOPATH/src/github.com/openshift/app-netutil/
make gofmt
make lint
```

To allow `gofmt` to update a file it has detected an error in:

```bash
gofmt -s -w <file>
```

## Certificate of Origin

In order to get a clear contribution chain of trust we use the [signed-off-by language](https://01.org/community/signed-process)
used by the Linux kernel project.

## Format of the patch

Beside the signed-off-by footer, we expect each patch to comply with the following format:

```
Change summary

More detailed explanation of your changes: Why and how.
Wrap it to 72 characters.
See [here] (http://chris.beams.io/posts/git-commit/)
for some more good advices.

Fixes #NUMBER (or URL to the issue)

Signed-off-by: <contributor@foo.com>
```

For example:

```
Fix poorly named identifiers
  
One identifier, fnname, in func.go was poorly named.  It has been renamed
to fnName.  Another identifier retval was not needed and has been removed
entirely.

Fixes #1
    
Signed-off-by: Abc Xyz <abc.xyz@intel.com>
```

## Pull requests

We accept GitHub pull requests.
