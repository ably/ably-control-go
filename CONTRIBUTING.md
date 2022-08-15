# Contributing to ably-control-go

1. Fork `github.com/ably/ably-control-go`
2. create your feature branch: `git checkout -b my-new-feature`
3. commit your changes (`git commit -am 'Add some feature'`)
4. ensure your code is formatted with `go fmt`
5. ensure you have added suitable tests and the test suite is passing.
6. push to the branch: `git push origin my-new-feature`
7. create a new Pull Request.

## Running Tests

The tests are ran against a real Ably account. For this to work `ABLY_ACCOUNT_TOKEN`
must be set with a valid [access token](https://ably.com/docs/control-api#creating-access-token).

Additionally, `ABLY_CONTROL_URL` can be set to an alternative control API endpoint.

The tests can then be ran with:

```
go test -v
``` 

## Release process

This library uses [semantic versioning](http://semver.org/). For each release, the following needs to be done:

1. Make sure the tests are passing in CI for the branch you're building
2. Create a new branch for the release, for example `release/1.2.3`
3. Update the CHANGELOG.md with any customer-affecting changes since the last release and add this to the git index
4. Replace all references of the current version number with the new version number and add this to the git index
5. Create a PR for the release branch
6. Once the PR is approved, merge it into `main`
7. Run `git tag <VERSION_NUMBER>` with the new version and push the tag to git
8. Create the release on Github, from the new tag, including populating the release notes
9. Update the [Ably Changelog](https://changelog.ably.com/) (via [headwayapp](https://headwayapp.co/)) with these changes (again, you can just copy the notes you added to the CHANGELOG)
