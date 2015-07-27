## Readers

There are two ways to gather server data in ResourceD:

**1. Using a scripting language.**

Use `GoStruct = "Shell"` and specify the script's path in `[GoStructFields] Command`.

Example: [darwin-memory.toml](https://github.com/resourced/resourced/blob/master/tests/data/resourced-configs/readers/darwin-memory.toml)


**2. Using Go natively.**

ResourceD provides a lot of [native Go readers](https://github.com/resourced/resourced/tree/master/readers).

To use them, define `GoStruct` field with the name of the struct.

To find out the names, look at `func init()` on each of the Go file in [readers](https://github.com/resourced/resourced/tree/master/readers) directory.
