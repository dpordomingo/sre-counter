# sre-counter

Example app for [SRE: Escalabilidad en entornos de alto rendimiento](https://trainingit.es/curso-sre-escalabilidad/) course.

## Install and use

If you have `Go`:

```bash
$ export GO111MODULE=on
$ go get -u github.com/dpordomingo/sre-counter/cmd/redis-counter
$ REDIS_HOST_PORT="hostname:6379" \
    SERVER_PORT=8090 \
    INSTANCE_NAME="instance name" \
    `go env GOPATH`/bin/redis-counter
```

And navigate to:
- `http://hostname:8090/counter`, to see the counter and increment it.
- `http://hostname:8090/reset`, to reset the counter.


For help, use:

```bash
$ redis-counter --help
```

## Build from sources

```bash
$ go build -o build/redis-counter cmd/redis-counter/*
```

## License

MIT License, see [LICENSE](./LICENSE.md).