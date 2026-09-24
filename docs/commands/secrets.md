# Secrets Command

The `secrets` command retrieves or generates the `secrets.yaml` bundle for the cluster.

The secrets bundle contains all sensitive cluster material including private keys, tokens, and certificates. If the file is SOPS-encrypted, it is decrypted automatically. [vals](https://github.com/helmfile/vals) references are also resolved if present.

## New Secrets

When no existing secrets bundle is found, a new one is generated and stored automatically. Where it is stored depends on the configuration:

- **Default (no `secretsProvider`):** the bundle is written to the local filesystem, next to `topf.yaml` as `secrets.yaml` (or the path set via `secretsPath`). It is SOPS-encrypted on write if a corresponding SOPS config is found; otherwise it is stored as plaintext.
- **With a `secretsProvider`:** the bundle is sent to the configured provider binary instead (see [secrets provider](../providers.md#secrets-provider)).

Regardless of whether the bundle was loaded from storage or freshly generated, `topf secrets` also prints it to stdout; storage happens automatically, so no output redirection is needed to create the file. The printed bundle is **not redacted**, even when `--redact` is enabled: it contains the actual secrets, so take care when piping or copying the output.

## Confirmation

When no existing secrets bundle is found and a new one needs to be generated, topf will prompt for confirmation before creating and storing it (unless the global `--confirm=false` flag is set, see [global flags](../configuration.md#global-flags)). This prevents accidental secret generation in interactive usage. In CI/CD pipelines, use `--confirm=false` to skip the prompt.

## Example Usage

```bash
# Load existing secrets, or generate and store a new bundle (prompts first)
topf secrets

# Generate secrets without confirmation (e.g. in CI)
topf secrets --confirm=false
```
