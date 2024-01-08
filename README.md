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

## License

MIT
