# Validation in generate protocol

## Summary

`jk generate` and `jk transform` have a tiny protocol between the
user-supplied script and the library function that will output files
(or to `stdout`): the script must supply a list of objects
representing the generated (or transformed) configuration, along with
pragma indicating e.g., the file into which a value should be written.

This RFC proposes that the protocol is expanded modestly to allow a
validation function to be included with configuration values, so that
e.g., an invocation of `jk generate` will fail (with errors) rather
than produce invalid output.

## Example

This example script combines the use of a (hypothetical) libary
providing a function for constructing configuration, and a custom
validation function, to generate validated configuration.

```js
// config.js
import * as param from '@jkcfg/std/param';
import { generateChart } from '@example.com/charts/mychart';

function checkMyName(value) {
  if (!value.metadata.name.startsWith('my-')) {
    return 'name does not start with my-'
  }
  return 'ok';
}

function addValidation({ ...fields }) {
  return { validate: checkMyName, ...fields };
}

const values = generateChart(param);
export default Promise.resolve(values).then(vals => vals.map(addValidation));
```

When used with

    jk generate ./config.js

this will instantiate the chart, and check each generated value passes
the validation as given in `checkMyName`. If any values fail
validation, the command fails and prints out the validation errors.

## Motivation

To date we have provided a mix of ways of _generating_ configuration:
charts, the short format (`@jkcfg/kubernetes/short`), Kustomize-like
patching, for example. Using one of these libraries with `jk generate`
does not in itself guarantee a usable configuration, however,
because --

 - there may be bugs in the libraries,
 - user-supplied code can introduce problems (e.g., a mistake in a
   template),
 - constraints may come from the user rather than the target system;
   e.g., names used in the configuration must follow a particular
   scheme,

and so on. To have more confidence of getting a usable configuration,
it's necessary to run a validation step over the output before
reporting success or failure.

## Design

This RFC proposes adding a hook to the "generate protocol" such that
values can be given a validation procedure, which will be run against
the value before output.

### Changes to generate protocol

Presently, `jk generate` expects a default export of type (in
TypeScript notation)

```typescript
interface Value {
    path: string;
    value: any | Promise<any>;
    format?: std.Format;
}

type ValueArg = Value[] | Promise<Value[]> | () => Value[];
```

In other words: either an array, or the promise of an array, or a
thunk returning an array, of objects; each of which gives a path and a
value (or promise of a value) and optionally an output format.

The proposal is to add another optional field to each object, for a
validation function:

```
type ValidateFn = (value: any) => string[];

interface Value {
    path: string;
    value: any | Promise<any>;
    format?: std.Format;
    validate?: ValidationFn;
}
```

When this function is present, it is run against the (resolved)
value. If any of the return values from running a validate procedure
indicate a failure, the validation errors are output and the whole
thing fails.

## Backward compatibility

Hitherto, no examples or libraries will supply a validate field, so
`jk generate` will behave exactly as before, until validation
functions are provided.

User code could in principle supply (erroneously!) a validate field --
in that case, generation may fail, either because the it's not a
function, or the return value is unexpected. Since that would be a
mistake, albeit a harmless one, this RFC does not attempt to account
for user code of that nature.

## Drawbacks and limitations

**Does not provide a hook for user-supplied validation**

Since the validation procedure is supplied with the generated values,
this design is most useful if you have a library that generates values
which must have a standard form -- but can't itself guarantee that
user input won't lead to invalid values. For example,
`@jkcfg/chart/template` lets you generate manifests with text
templates you supply yourself; if it accompanied those with schema
validation, it would catch cases where the templates (or other input)
resulted in invalid manifests.

It is less useful in itself for when you want the user to supply their
own validation, absent support for that in libraries (or elsewise
further work). For example, a user might want to have generated
Kubernetes manifests checked against a schema, but also to make sure
some of their own constraints are satisfied, like complying with a
naming scheme or whatever. But this does open the door to that kind of
support in libraries, by providing a hook for validation.

Further work could establish ways of combining validation, so that
user-supplied or third-party validation can be combined with standard
validation.

In the meantime, this does not shut down other routes to validating
configuration -- libraries that work with this design will likely be
adaptable to other modes of use, because the type for validation
functions is quite broad.

**Harder to fit into `jk transform`**

When using `jk transform`, configuration values are supplied rather
than calculated. So there's not the same opportunity to include a
validation function. To use validation with `jk tansform`, the
function would have to be chosen by the user (e.g., as the default
export of a module they name).

**Doesn't let you validate whole configurations**

Some kinds of validation might need to look at more than one of the
values; for example, to make sure that for each Kubernetes Deployment
there is a corresponding Service. This proposal does not provide for
that kind of calculation, since the validate procedure accompanies a
single value, and is only supplied that value when invoked.

## Alternatives

_Explain other designs or formulations that were considered (including
doing nothing, if not already covered above), and why the proposed
design is superior._

**Run validation as a separate command**

E.g., `jk validate`, which validates files (or stdin), given a script
or module with the validate function.

This covers a slightly different use case, that of checking existing
files (including files not generated by `jk`), and is more flexible in
that regard. But it makes the user choose the validation, whereas the
idea in this RFC is for libraries to provide validation along with the
values.

**Leave validation to the user script, rather than running it as part
of generation**

Instead of adding a validate function for each value and running it
automatically, require the user to call it as part of their script:

```js
import { validateAll, outputValues } from '@jkcfg/std/generate';
import { validate } from 'k8s';

const values = generateChart(param);
const validationErrors = validateAll(values, validate);
if (validationErrors === 'ok') {
  outputValues(values);
}
// ... deal with errors otherwise
```

This pushes a little complexity onto the user. The gain is that you
can decide what to do with the validation errors -- you might choose
to log them and move on, for instance. However, as a starting point
including the validation in the generation step per this RFC is a
reasonable default, and does not prevent other schemes.

**Require libraries to do their own validation**

An argument implied in this RFC is that libraries are good at choosing
the validation to be done on values they generate. Why not just
require those libraries to do the validation, before returning the
values?

One reason is that it would make it harder for users to interpose
their own validation, unless that were also built into each library.

## Unresolved questions

_Keep track here of questions that come up while this is a draft.
Ideally, there will be nothing unresolved by the time the RFC is
accepted. It is OK to resolve a question by explaining why it
does not need to be answered_ yet _._

 - Better to use some other type for the validate return value?

       type ValidateFn = (value: any) => boolean | string[]

   (that is, the result is either `true` for "pass validation",
   `false` for "fail non-specifically", or a list of validation
   errors)? Or just `string[]` with an empty list meaning none?

