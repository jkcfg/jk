# Built-in JSON Schema validation

## Summary

This RFC proposes that validation with JSON Schema is built into the
jk runtime.

Since the values `jk` scripts produce for output are in general
JSON-serialisable (even if not destined for a JSON file), adding
schema support would cover a broad set of validation
requirements. JSON schemas exist, or can be derived, for many APIs,
including the Kubernetes API.

Implementing this proposal would mean adding a modest API, via the RPC
mechanism, for dealing with schemas and values.

## Example

```js
import { read, log } from '@jkcfg/std';
import { dir, info } from '@jkcfg/fs';
import { withModuleRef } from '@jkcfg/std/resource';
import { validateByResource } from '@jkcfg/std/validate/schema';

// validate all the YAML files in $PWD, against the schema in
// 'schema.json', relative to this module.

const schemaPath = 'schema.json';

const d = dir('.');
const yamls = d.files.filter(f => f.path.endsWith('.yaml'));

// File reads and validation are both async, but we can do them concurrently per file.
async function validateFile(path) {
  const obj = await read(path);
  // withModuleRef is needed so that validateByResource can resolve the path relative to this module
  const validation = await withModuleRef(ref => validateByResource(obj, schemaFile, ref));
  return { path, validation };
}

const results = Promise.all(yamls.map(f => validateFile(f.path)));
results.then(function(rs) {
  for (const { path, validation } of rs) {
    if (validation == 'ok') {
      log(`${path}: ok`);
      return;
    }
    // otherwise, it's a list of errors
    for (err of validation) {
      log(`${path}: error: ${err}`);
    }
  }
});
```

## Motivation

`jk` and its libraries provide various means for generating
configuration, but there is in general no guarantee that the output is
usable. For that reason, to become a more rounded tool, validation
should be considered a core capability for `jk`.

One very general-purpose means of validation is by checking against a
data model or schema. In the case of JSON-serialisable JavaScript
objects -- which is effectively the data model of `jk` -- a reasonable
choice of schema language is JSON Schema. Schemas using JSON Schema
(and the closely related OpenAPI data model, more on which see below)
are available for lots of APIs, and it has support in tooling e.g.,
[VSCode](https://code.visualstudio.com/docs/languages/json#_json-schemas-and-settings).

Building in schema validation means libraries are free to rely on it,
encouraging more widespread validation.

## Design

### Discussion

**Synchronisation**

The expectation might be that validation is synchronous, since it is
used mostwhere as a predicate. There's good reasons why it should be
asynchronous, though:

 - validating with a file means a file read, and those are
   asynchronous, so the precedent exists;
 - resolving references may involve network requests, and those are
   naturally represented as asynchronous.

**Resolving schemas**

 - `$schema` values will likely refer to public URIs (but this RFC
   says nothing about respecting `$schema`; that could follow in
   another RFC, perhaps).
 - schemas will often refer to other schemas with `$ref: <json
   pointer>`, which can come from the same file, or from a file
   relative to the current file, or a URI (though it's not required to
   be a network locator, i.e., it's up to the library how to process
   it)
 - libraries will want to include their schemata in the distributed
   package, so will want to be able to load them as a resource.

**Caching**

Usually, caching is strictly an optimisation, and as such it is not
necessary to discuss it by way of specification. But it's worth
mentioning that Go libraries for schema checking tend to come with a
cache for external schemas; and some schemas will be used extensively
during a run (e.g., the schema for a Kubernetes Deployment). It would
be good to work in sympathy with any caching in the library, where
possible.

### API in JavaScript

The result of a schema check is either that everything was OK, or that
there were some specific problems.

```typescript
type ValidationResult = 'ok' | string[];
```

Either `'ok'` or `string[]` indicates a successful validation call, so
either can appear as the resolution to a promise. A promise will be
rejected if there's an error in the validation process itself -- for
example, a file cannot be found.

The obvious mode of use is to supply the value and the schema, as
JavaScript objects:

```typescript
validateBySchema: (value: any, schemaObj: any) => ValidationResult;
```

This isn't necessarily the most useful though; it will often be more
useful to be able to refer to a file, either an input (relative to the
`jk` invocation) or a resource (relative to the module). At first
glance, you'd expect this latter could be implemented by simply
reading the file with `read`, then using `validateBySchema`; but,
resolving `$ref`s may need a base path, and this will only be
available if the path (or module reference) is supplied to the
runtime.

```typescript
validateByFile: (value: any, path: string) => Promise<ValidationResult>
validateByResource: (value: any, path: string) => Promise<ValidationResult>
```

Getting the module reference will require a bit of extra machinery in
the runtime. Adding a module reference argument to the RPC protocol,
and making the value available via the `@jkcfg/std/resource` magic
import would be one way to enable it. There are surely similar,
general schemes -- the important point being that it should remain an
internal mechanism as much as possible.

### Changes to the runtime

The generated resource module in `std/resource.go` needs to provide a
way for `validateByResource(...)` to refer to resources; i.e., to the
base path of the _importing_ module. To date, the generated module
exports functions that close over the hash referring to the
module. But this approach doesn't suffice: we don't want to put
validateByResource (and future procedures with a similar requirement)
in the generated module.

A general mechanism is to make the module reference itself
available. This means that module X can import `@jkcfg/std/resource`,
then supply its own module reference (-> access to its resources) to a
procedure from elsewhere. The module reference is just a value, though
an opaque one -- it can be passed around, and e.g., used with
`read(...)` to read files from the module's directory. But that's
already true with `@jkcfg/std/resource#read` itself -- a module can
"give away" access to its files by supplying that to another module,
by design.

```typescript
type ModuleRef = string;
withModuleRef<T>: ((ref: ModuleRef) => T) => T;
```

## Backward compatibility

If schema checking via `$schema` is on by default, it may mean some
values that previously worked would break. But since those values are
invalid, this ought to be considered a good thing, like a new compiler
release finding a type error it was overlooking before.

## Drawbacks and limitations

The main downsides to adding JSON schema validation to the runtime
are:

 - it's a bit more code, and another dependency
 - it identifies JSON schema as a preferred schema language, which if
   it turns out to be a poor choice, might need deprecation (and
   nobody likes dealing with deprecation).

## Alternatives

**Leave schema validation entirely out of scope**

In other words, do nothing. This has the benefit of requiring no
engineering, but it does leave the responsibility -- if we still
recognise validation as important -- of showing how to do things
elsewise.

It also means that the burden of implementing (schema) validation
always rests on the end user -- it can't easily be built into
libraries, or be done automatically or by default.

**Use OpenAPI's version of JSON schema**

OpenAPI (formerly Swagger) has its own formulation of data schema,
which is [adapted from JSON
Schema](https://swagger.io/docs/specification/data-models/).

The OpenAPI schema is not a subset of JSON schema -- it has keywords
that are not in JSON Schema, and it alters the meaning of some
keywords. However, it does appear possible to [translate from OpenAPI
data model to JSON
schema](https://github.com/instrumenta/openapi2jsonschema).

To some extent this decision is a matter of taste -- there's no
decisive point in favour of one or other. In the future both varieties
could be built in, with a switch in the API to specify which to
use. Or, [the specs will converge
again](https://apisyouwonthate.com/blog/solving-openapi-and-json-schema-divergence).

**Don't build it into the runtime, but use a plugin mechanism**

The argument here would be that to keep the runtime lean, validation
should be opt-in and "user pays". A (hypothetical) plugin system that
let binaries be distributed with libraries would enable a json-schema
library, that other libraries could then use for their own validation,
etc.

However: with schema validation built in, `jk` can have a command that
will validate directly against a schema:

    jk validate --schema ./definition.json -R ./config/

whereas supporting schema validation only via a library would mean _at
least_ that the user has to fetch the dependencies first, and probably
write a shim (similar to the example right at the top) to use the
library.

This comes down to judging whether JSON schema validation will be
important and useful enough to include as a core part of the
runtime. In its favour:

 - it's broadly applicable; that is, would be useful for validating
   most if not all values generated in jk;
 - it's on a standards track, so it is reasonable to expect it to
   stick around
 - we think validation is a core use of jk

Ultimately, this RFC argues that schema validation is a force
multiplier -- it can enhance other features (like guarding config
generation or transformations) -- and making this convenient is worthy
of being a first-class feature.

**Use a JavaScript implementation of JSON Schema, from libraries**

There is [at least one decent-looking
implementation](https://www.npmjs.com/package/jsonschema) of JSON
Schema in JavaScript, with no apparent dependency on Node.JS standard
lib. Instead of building JSON schema into the runtime, libraries that
were interested in using JSON schema to validate things could depend
on that implementation.

The argument against this runs similarly to that above regarding
plugins: each library using schema validation would have to depend on
the jsonschema library; and, there would be no route to ad-hoc schema
validation.
