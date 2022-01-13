# Proposed ZDNS Refactoring

NOTE: This document is a draft and still being updated, revised, edited, etc. If you have any questions, please contact me!

## High-Level Code Changes

- Use `spf13/cobra` as CLI framework
    - This will help with gathering the options from the command line and then passing them to the CLI. It won't have much to do with the refactoring of the library itself.
    - This will also help with reducing some of the error checking that goes on in `main.go` currently. 
- Migrate all packages from `modules` to `pkg` directory
    - We'll do a few things to ensure a cleaner end user experience:
        - decouple our logic from the logic of any deps that we use (e.g., this will likely require some work on the miekg/dns side of things)
        - explicitly version our deps
    - As I perform this migration, I'll determine and consult the team as needed on what should live where. Majority of this functionality will go into `pkg` but I think there are some pieces that will end up in `internal`.

### External API

This section will likely change, however, it's what the current vision will be. No plan survives first contact, but, failing to plan is planning to fail.

The libary will receive some sort of configuration object, run the tasks, and then exit. All of the goroutines spawned within a call to the library must be terminated within that same call. This contains the asynchronous behavior to a single function, providing these benefits:

- Eliminating data races
- Clean interface boundary
- Easier to test
- Easier to read/understand

The library will accept a modified `Lookup` interface and will then return the results of the lookup. Depending on the structure, we may want to use Go 1.18's generics here, but I'll likely not know for sure until I'm deeper into the refactoring. The idea here is that the user (CLI or otherwise) will just have create an struct that satisfies this object, pass it through to the library and the library will handle the rest. The library will likely only have a few (or one) functions, I envision the following:

- `DoLookup` - takes the interface and runs the lookup(s) specified within.

The packages that currently live in `modules` will be moved into the `pkg` directory to follow a more standardized go structure. Some functionality may be moved to `internal` to prevent use by library consumers. 

The `Lookup` interface described above will contain the options and parameters that are being passed into the library from the CLI. This will be instantiated in the form of a `DNSClient`. In this way, a user (e.g., the CLI) will be able to create this `DNSClient` and then call its `DoLookup` method with the appropriate parameters. As much of the complexity as possible will be contained behind this layer of abstraction, however, some things will not be able to live there. We will likely need to allow for passing sockets to the library so that the DNS queries can be made more efficiently and we aren't waiting on creation/destruction of a socket each time a request is made.

### Likely Sticking points
- clean ways of handling concurrency in the library
- clean shutdown of goroutines
- clean method of passing parameters

## CI Changes
- Remain with GitHub Actions. Continue to expand testing and checks to include useful info for PRs as well as checks before merging.

## Versioning System
- Each PR into master gets a new subminor (patch) version. For example, v1.0.1 -> v1.0.2 after a successful PR. 
- We're going to start with major version v0 (if possible and approved). This gives us the time to do things "the right way" and mess around with the API for a few weeks until it's stabilized. Then we go to v1 and start committing to stability.
- At any point, we can deploy a new minor or major version - the CI actions will just take the latest and go from there.
- `zmap/dns` will require a bit of a different versioning scheme. 
    - My recommendation is to follow the same scheme as `miekg/dns`, even if we are missing big chunks. This'll make it clearer to us and others what we're basing off of.
    - It may also make sense to vendor our own version of miekg/dns, to make the project go-gettable

## Revendoring DNS library

This tool relies heavily on the `miekg/dns` library. I propose that if it's possible to revendor this under the `zmap` project, then we should do so to ensure that we have consistent control and versioning over our fork, but also so that we can make the project "go gettable". The idea is to simply install go, run `go get github.com/zmap/znds` and be off to the races. This is not possible in the current state of things.

We may or may not be able to do this, however, I'll speak with the team on this aspect.

## Deployment options

- Vishal mentioned that we want to run this from several vantage points. This means we'll want portability, and to not have to deal with rebuilding each time.
- Original thought was to have this bundled in a docker image, however, that's extra machinery that's not necessary assuming that we can make this "go gettable".

## Testing Changes

- Increase quanitity and quality of tests 
- I imagine three types of tests
    - Unit tests: Test each function/file. 
        - Investigate code coverage for this
    - Integration tests: I envision this primarily serving the purpose of ensuring _completeness_ but not necessarily _correctness_. By this I mean that our integration tests will ensure the system doesn't break, but maybe not that it handles every single edge case.
        - This might be a suitable place to poll a third-party DNS like Google's DNS or something similar.
    - E2E or System tests: These will check our correctness.
        - I'm thinking of a somewhat extension testing framework. Idea being that we somehow (maybe into K8s) deploy our own `bind` server and then via some kind of config file ensure that it responds (or doesn't respond) in the correct way.
    - Regression tests: These tests would ideally be run on a nightly or weekly basis. They'll be the most expensive and long-running, thus, we won't want to run them each time, but we should have a suite of them.