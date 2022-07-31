# Changelog

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic
Versioning](http://semver.org/spec/v2.0.0.html).

## [v0.0.32] - 2022-07-31

### Fixed

- Use github.com/aarondl/json fork to avoid asymmetric payload issues
- Set application/json content type headers for json payloads
- Fix client generation error

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
