## Releasing a new version of kube-linter

This doc contains all steps required to release a new version of kube-linter!

### Create a release tag

Decide on the version number of the next release and create a git tag for it:

```shell
$ git pull --tags
$ git checkout main
$ git tag <new.release.version> # NOTE: PLEASE DO NOT USE THE PREFIX "v" FOR THE TAG!
$ git push origin <new.release.version>
```

This will trigger a workflow that:
- Pushes docker images to <registry>/stackrox/kube-linter:<new.release.version>.
- Uploads latest built assets to the draft release.

### Publish the release notes

Kube-linter uses the GitHub action [release-drafter](https://github.com/release-drafter/release-drafter).
This will create a draft release upon each commit on main.

You should see the draft release under [releases](https://github.com/stackrox/kube-linter/releases).

Ensure you update the following:
- The title and flag should reflect the new release version.
- The compare link should reflect the new release version: `https://github.com/stackrox/kube-linter/compare/0.3.0...<new.release.version>`

If you have made the required updates, review the contents once more.

If everything checks out, you can publish the release!

## Troubleshooting

### Wrong tag used

In case a wrong tag was used (i.e. prefixed with `v`, typo etc.), you could do the following:
- Delete the published release.
- Delete the existing tag.
- Tag the main commit with the new tag. After the workflow is run successfully, you can use the newly created draft 
  release for publishing.

### Issues with kube-linter-action

For some releases, it may happen that there's an issue with the kube-linter action. Most probably, the issue will be
due to the uploaded assets.

In this case, you could:
- Remove / re-upload the assets manually, making the required changes for the release as a temporary workaround.
