# Kleiner - CLI to set up distribution of other clis

Let's face it, go is an excellent tool to build self contained clis, however
there is always a question of distribution. You could ask your friends to download
new binary from the releases page, however wouldn't it be cooler to have a cli
that can report about updates and update itself?

This has been done a load of time, however kleiner allows to do that at scale
since it's so damn easy to generate yet another cli and have it's distribution
set up immediately. All you need to do is to add your commands and you're done.

We don't need to have another tool for scaffolding, hence please have [cobra-cli](https://github.com/spf13/cobra-cli)
installed

Local cobra config was used using:

```
cobra-cli init --config .cobra.yaml
```

We won't be building binaries ourselves either, the will be done by [goreleaser](https://goreleaser.com/install)

Be sure to have it [installed](https://goreleaser.com/quick-start/), also [generate](https://github.com/settings/tokens/new?scopes=repo,write:packages) a github token and put it
in `GITHUB_TOKEN` variable

## Inspiration

The code is inspired a lot and uses parts of the [flyctl](https://github.com/superfly/flyctl) code. The files
with the borrowings mention this explicitly.

## License

MIT
