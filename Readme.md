# go-run-programs

Go run programs in the background and check until they are ready.

## Go API

The original purpose of this was to run programs in the background to test against,
specifically in a way that works the same in Windows and Unix for simple multi-platform CI.

<!-- TODO: godoc and example_test.go -->

```go
import (
    "github.com/gulducat/go-run-programs/program"
)

func TestThatMyThingTalksToTheOtherThing(t *testing.T) {
    t.Log("starting the other-thing in the background")
    p := program.Program{
        Name:    "OtherThing",
        Command: "other-thing server run",
        Check:   "other-thing is-running",
    }
    stop, err := program.RunInBackground(context.Background(), p)
    t.Cleanup(stop) // stop it when tests are done.
    if err != nil {
        t.Fatal(err)
    }
    
    t.Log("ok ready to run tests now")
}
```

## CLI

It can also be run as a CLI for cute little local dev environments or such.

```go
go-run-programs hashicorp.hcl
```

```hcl
# hashicorp.hcl

program "consul" {
  command = "consul agent -dev"
  check   = "consul members"
}
  
program "nomad" {
  command = "nomad agent -dev"
  check   = "nomad node status"
}
  
program "vault" {
  command = "vault server -dev"
  check   = "vault status"
  env = {
    # default client is https, but vault server -dev is http
    VAULT_ADDR = "http://localhost:8200"
  }
}
```
