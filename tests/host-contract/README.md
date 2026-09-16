# Published host contract tests

Run `just host-contract` from the repository root with Go 1.27+ and a C
compiler. This separate module pins the supported Terva host v0.137.0 and
uses its actual permission resolver against our real manifest in temporary
homes. It needs module downloads on first run; production still vendors only
the extension SDK. Do not copy installed host config or credentials into tests.

These tests verify host policy, not an interactive approval UI or a complete
Terva CLI launch. The subprocess conformance suite independently checks actual
SDK registrations. Keep both pins aligned when changing the supported host.
