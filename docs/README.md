First comparison round:

```console
goos: darwin
goarch: amd64
pkg: github.com/damienstanton/twocrawlers/A
BenchmarkA-4      229148              5228 ns/op             224 B/op          3 allocs/op
PASS
ok      github.com/damienstanton/twocrawlers/A  1.327s
goos: darwin
goarch: amd64
pkg: github.com/damienstanton/twocrawlers/B
BenchmarkB-4       93612             12434 ns/op             933 B/op          8 allocs/op
PASS
ok      github.com/damienstanton/twocrawlers/B  1.383s
```

© 2020 Damien Stanton

See LICENSE for details.

