# Runtime Class Manager release process

The vast majority of the release process is handled by GitHub workflows triggered by a tag push.
See [release.yaml](./.github/workflows/release.yml) for the main workflow.

First, let's start by setting a `$TAG` environment variable to the version we'll release, so that we can reference it later.

For example, if the most recent release was `v0.1.0` and subsequent changes on
the `main` branch include the usual collection of patches, fixes and features,
we'd cut `v0.2.0`:

```console
export TAG=v0.2.0 # CHANGEME
```

Next, be sure that CI is green for the current commit on the `main` branch.

Then, to push the tag, do the following:

```console
git checkout main
git remote add upstream git@github.com:spinframework/runtime-class-manager
git pull upstream main
git tag --sign $TAG --message "Runtime Class Manager $TAG"
git push upstream $TAG
```

Observe that the [Release Workflow run for the tag](https://github.com/spinframework/runtime-class-manager/actions/workflows/release.yml) completed successfully.

Next, you'll need to update the documentation:

```console
git clone git@github.com:spinframework/spinkube-docs
cd spinkube-docs
```

Change all references from the previous version to the new version.

Contribute those changes and open a PR.

As an optional step, you can run a set of smoke tests to ensure the latest release works as expected.

Finally, announce the new release on the #spinkube CNCF Slack channel.
