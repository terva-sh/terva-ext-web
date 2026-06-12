# Release-cut targets — maintainer-only tooling, imported by the main
# justfile via `import?` so the public tree (which excludes this file
# and scripts/release.sh) keeps a working justfile. Modeled on terva's
# flow, minus versions/tags/binaries: the curated public history IS
# the release.
#
# Typical cut:
#   just release-cut            prepare the curation worktree + worklist
#   (author the curated commits in ../zot-web-release-cut)
#   just release-verify         prove the public tree is publishable
#   just release-publish        push release to origin + the staging gate
#   (test + inspect in the gate, then push from inside it to go live)

# Prepare a curation worktree for a release cut from main.
release-cut:
    ./scripts/release.sh cut

# Verify the curated tree: identity, scrub, toolchain.
release-verify:
    ./scripts/release.sh verify

# Push release to origin (backup) and the staging gate; drop the cut/N marker.
release-publish:
    ./scripts/release.sh publish

# Show the in-progress cut's state and worktree status.
release-status:
    ./scripts/release.sh status

# Print the curation worklist for the in-progress cut.
release-worklist:
    ./scripts/release.sh worklist

# Tear down an in-progress cut and restore branches.
release-abort:
    ./scripts/release.sh abort
