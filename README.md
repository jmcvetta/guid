# guid - A library for generating globally unique identifiers

Derived from [noeqd][], which in turn is based on [snowflake][].

## Motivation

GUIDs (Globally Unique IDs) are useful for a number a obvious reasons:
database keys, logging, etc.

Generating GUIDs with pure randomness is not always ideal because it doesn't
cluster well, produces terrible locality, and no insight as to when it was
generated.

This network service should also have these properties (Differences from [snowflake][]):

* easy distribution with *no dependencies* and little to *no setup*
* dirt simple wire-protocol (trivial to implement clients without added dependencies and complexity)
* low memory footprint (starts and stays around ~1MB)
* zero configuration
* reduced network IO when multiple keys are needed at once

## Glossary of terms to follow

* `GUID`: Globally Unique Identifier
* `datacenter`: A facility used to house computer systems.
* `worker`: A single process with a worker and datacenter ID combination unique to their cohort.
* `datacenter-id`: An integer representing a particular datacenter.
* `worker-id`: An integer representing a particular worker.
* `machine-id`: The comination of `datacenter-id` and `worker-id`
* `twepoch`: custom epoch (same as [snowflake][])

## Important Note

Reliability, and guarantees depend on:

**System clock depedency and skew protection:** - (From [snowflake][] README and slightly modified)

You should use NTP to keep your system clock accurate. Noeq protects from
non-monotonic clocks, i.e. clocks that run backwards. If your clock is running
fast and NTP tells it to repeat a few milliseconds, Noeq will refuse to
generate ids until a time that is after the last time we generated an id. Even
better, run in a mode where ntp won't move the clock backwards. See
<http://wiki.dovecot.org/TimeMovedBackwards#Time_synchronization> for tips on how
to do this.

**Avoiding the reuse of a worker-id + datacenter-id too quickly**

It's important to know that a newly born process has no way of tracking its
previous life and where it left of. This means time could have moved
backwards while it was dead.

It's important to **not** use the same worker-id + datacenter-id without
telling the new process when to start generating new IDs to avoid duplicates.

It is only safe to reuse the same worker-id + datacenter-id when you can
guarantee the current time is greater than the time of death. You can use the
`-t` option to specifiy this.

You may have up to 1024 machine ids. It's generally safe to not reuse them
until you've reached this limit.


# How it works

## GUID generation and guarantees

GUIDs are represented as 64bit integers and are composed of (as described by the [snowflake][] README):

* time - 41 bits (millisecond precision with a custom epoch gives us 69 years)
* configured machine id - 10 bits - gives us up to 1024 machines
* sequence number - 12 bits - rolls over every 4096 per machine (with protection to avoid rollover in the same ms)

## Sorting - Time Ordered

*Strictly sorted*:

* GUIDs generated in a single request by a worker will strictly sort.
* GUIDs generated one second or longer apart, by more than one worker, will strictly sort.
* GUIDs generated over multiple requests by the same worker, will strictly sort.

*Roughly sorted*:

* GUIDs generated by multiple workers within a second could roughly sort.

An example of roughly sorted:

If client A requests three GUIDs from worker A in one request, and client B
requests three GUIDs from worker B in another request, and both requests are
processed within the same second, together they may sort like:

		GUID-A1
		GUID-A2
		GUID-B1
		GUID-B2
		GUID-A3
		GUID-B3

NOTE: The A GUIDs will strictly sort, as will B's.


## Status

[![Build Status](https://drone.io/github.com/jmcvetta/guid/status.png)](https://drone.io/github.com/jmcvetta/guid/latest)


## Contributing

This is Github. You know the drill. Please make sure you keep your changes in a
branch other than `master` and in nice, clean, atomic commits. If you modify a
`.go` file, please use `gofmt` with no parameters to format it; then hit the
pull-request button.

## Issues

These are tracked in this repos Github [issues tracker](http://github.com/jmcvetta/guid/issues).

## Thank you

I want to make sure I give Blake Mizerany
([@bmizerany](http://twitter.com/bmizerany)) at Heroku and the Snowflake team
at Twitter as much credit as possible. The heart of this program is their
doing.

## LICENSE

Copyright (C) 2012 by Jason McVetta

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE. 

[Doozer]: http://github.com/ha/doozerd
[snowflake]: http://github.com/twitter/snowflake
[noeqd]: http://github.com/bmizerany/noeqd
