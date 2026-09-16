# Published host contract tests

Run `just host-contract` with Go 1.27+, Bash and a C compiler. Windows requires
Git Bash. This separate module pins Terva v0.137.0 and downloads its dependencies
on first use. Production builds use the vendored extension SDK.

The tests run the published permission resolver against the real manifest in
all five approval modes, with explicit user allow and deny rules. They also
launch the extension through the published `extdriver`, using the manifest's
Bash command. They check registration, command dispatch, session changes,
network saves, oversized-image errors and shutdown.

`WEB_CONFORMANCE_INSTALL` selects an extracted archive for driver tests;
otherwise the test builds a temporary installation. Fixtures use temporary
homes and synthetic settings. Do not copy installed credentials into them.

These tests do not cover the interactive approval UI or the full CLI install,
upgrade and rollback process. The root subprocess conformance suite separately
checks SDK wire behavior. Keep the module pin aligned with the supported host.
