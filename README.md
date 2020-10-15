= Split batch processor

This is an example processor, accompanying part of the blog post: ["Extending the OpenTelemetry Collector with your own components"](https://medium.com/p/64c10cf675db).

Do not use this processor in your own OpenTelemetry Collector distribution! If you need the feature provided by this processor, use the module [`github.com/open-telemetry/opentelemetry-collector-contrib/pkg/batchpertrace`](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/master/pkg/batchpertrace) instead.

== Running

```console
$ go run ./cmd --config example.yaml
```

== Distribution

To build a distribution with this module using the [OpenTelemetry Collector builder](https://github.com/observatorium/opentelemetry-collector-builder), run:

```console
$ opentelemetry-collector-builder --config builder.yaml
```