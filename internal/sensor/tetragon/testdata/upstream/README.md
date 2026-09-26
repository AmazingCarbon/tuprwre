# Upstream Tetragon event samples

`filename-access.jsonl` is a real `process_kprobe` event exported by Cilium
Tetragon, copied verbatim (only whitespace was removed so it fits on one
JSONL line) from the recorded output in Tetragon's documentation:

- source: https://github.com/cilium/tetragon/blob/main/docs/content/en/docs/use-cases/filename-access.md
- upstream project: cilium/tetragon (Apache-2.0)
- retrieved: tag `main`, 2026-09

It is a `security_file_permission` kprobe with `int_arg: 2` (MAY_WRITE) on a
real `file_arg`:

```json
"args": [
  {"file_arg": {"path": "/etc/passwd", "permission": "-rw-r--r--"}},
  {"int_arg": 2}
]
```

`file_arg` is the wire form of Tetragon's `KprobeFile` message
(`api/v1/tetragon/tetragon.proto`), whose fields `mount`, `path`, `flags` and
`permission` are declared directly — it is **not** wrapped in a nested `File`
message. `TestUpstreamFileArg` asserts the adapter reads this real shape, so a
regression to an assumed shape fails in CI without any environment variable.

The sample carries no secrets: it is Tetragon's own published documentation
output (Pod/container identifiers only).
