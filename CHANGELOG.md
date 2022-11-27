# Changelog

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic
Versioning](http://semver.org/spec/v2.0.0.html).

## [v0.0.51] - 2022-11-27

### Fixed

- Go server: Fix an issue with int arrays with a format in parameters.

## [v0.0.50] - 2022-11-04

### Changed

- Go client/server: Current OneOf strategy thrown out for `any`
  until such time that we can improve it to work with primitives.

## [v0.0.49] - 2022-10-21

### Fixed

- Go server/client: Bump the opt package to not switch json packages mid-encode.
- Go client: Fix a bug where you could not use WithDebug with the generated
  client due to empty url.

## [v0.0.48] - 2022-10-20

### Fixed

- Go server/client: Bump the opt and json package to have symmetrical behavior
  in edge cases in encoding.

## [v0.0.47] - 2022-10-16

### Changed

- Go server/client: Now use dfs for outputting embedded structs. This creates
  more smaller files, previously this would output embedded items in the same
  file as their parents which was nice, but had corner cases which would break
  generation.

## [v0.0.46] - 2022-09-25

### Changed

- Go server: Now uses pointers for returns where applicable to allow
  `return nil, err`.

## [v0.0.45] - 2022-09-24

### Changed

- Go server: Change json payload reading in server to keep the payload in a buf
  and pass it into r.Body again so the body could be read again if necessary.
  This prolongs the duration of the bytes.Buffer pool alloc slightly.

## [v0.0.44] - 2022-09-24

### Fixed

- openapi3: Fix it such that path parameters are required to be provided in
  either the path or the operation or it will fail validation.

## [v0.0.43] - 2022-09-24

### Changed

- Go client: Change url handling to be much less annoying. Types are only
  generated when necessary (there are more than 1 url or there's a parameterized
  url). The client now carries its own url builder as the default to be used.
- Go client/server: Change response handling to be much more streamlined. Oneof
  interface enforcers are only created and used when there's more than one type.
  When there's only a single type both the server and client simply use that
  type inline.
- Go client/server: io.ReadCloser is used when non-json mime types are used for
  responses, previously this was disallowed by the spec validation.

## [v0.0.42] - 2022-09-08

### Fixed

- Fix parameter ordering when using param refs

## [v0.0.41] - 2022-09-07

### Fixed

- Fix parameter reference validation for url params

## [v0.0.40] - 2022-09-07

### Fixed

- Go client: Use the limiter passed in to constructor

## [v0.0.39] - 2022-08-29

### Fixed

- Go client: Allow response to introspect the response on response code failures
- Go client: Client now produces error that contains the status code

## [v0.0.38] - 2022-08-28

### Fixed

- Fix bad switch in status code handling for generated client

## [v0.0.37] - 2022-08-14

### Changed

- Revamped Go server and client parameter handling to be able to meet the spec
  more closely. In the spec there are objects, arrays, and primitives allowed in
  parameters. In general many of these are very useful because of body payloads
  being sufficient for anything complicated and the serialization methods they
  use are bizarre. However it can make sense to have an array of values
  especially for query/header/cookie values where there are multiple values
  allowed in the HTTP spec. Given that, these now support arrays and arrays of
  enums (which are still only able to be strings).

## [v0.0.36] - 2022-08-12

### Added

- Add `decimaltype` parameter, behaves exactly like the timetype parameter in that
  it will replace `string` types in Go with `decimal.Decimal` from the
  shopspring decimal package, but this disables string validation for it.
- When decimal format is specified on a string type it will validate the string
  when `decimaltype!=shopspring`

### Changed

- Change `uuid` format to be controlled by a `uuidtype` parameter to the
  generator. Behaves exactly the same as `decimaltype` and `timetype` but uses
  google's uuid.UUID package. String validation is disabled when this param
  is active.

## [v0.0.35] - 2022-08-02

- Fix empty validation cases due to `format`

## [v0.0.34] - 2022-08-01

### Fixed

- Fix a few bugs in validation of optional and nullable fields. It will now
  generate code that checks for the presence of the value before attempting
  to validate it.

## [v0.0.33] - 2022-08-01

### Fixed

- Fix bug where validation of non-required refs did not deref omit
  properly.

## [v0.0.32] - 2022-07-31

### Added

- Added CI

### Fixed

- Use github.com/aarondl/json fork to avoid asymmetric payload issues
- Set application/json content type headers for json payloads
- Fix client generation error
- Fix lint errors

## [v0.0.31] - 2022-07-30

### Fixed

- Fix enum values null/omit/omitnullability
- Fix error message for required field missing in props error

## [v0.0.30] - 2022-07-30

### Added

- Add validation for required fields who do not appear in an object type's
  property list.

## [v0.0.29] - 2022-07-27

### Changed

- Upgraded yaml to v3, this change came with a restriction that all keys for
  yaml objects must be keys as in: `map[string]any`, because the spec is also
  able to be JSON, this property must hold anyway.

## [v0.0.28] - 2022-07-27

### Added

- Added type safe base urls to client methods

### Fixed

- Do not try to generate body params for non-application/json request bodies
  and also remove validation to prevent non-application/json bodies, this allows
  a handler to handle the body in whichever way they desire, while still using
  the open api spec to its full potential.

## [v0.0.27] - 2022-07-21

### Fixed

- Fix local clients debug output.

## [v0.0.26] - 2022-07-21

### Added

- Add Go client constructor (NewLocalClient) to enable testing against
  an http server.
- Add support for enums to operation parameters.

### Fixed

- Fix omit/null/omitnull handling to only be used at references to the types
  and never inside types themselves. Although openapi3 defines a type as being
  able to be null, it never matters until you need to use that type somewhere.

## [v0.0.25] - 2022-07-20

### Fixed

- Upgrade deps
- Fix omit/null/omitnull around an object's fields

## [v0.0.24] - 2022-07-19

### Added

- Add go client

### Fixed

- Fix array validation
- Fix support for inline array request bodies
