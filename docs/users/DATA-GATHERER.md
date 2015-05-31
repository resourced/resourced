## Data Gatherer

In general there are two ways to gather server data in ResourceD:

**1. Using a scripting language.**

Simply specify the script's path in `Command` field. Example: [darwin-memory.toml](https://github.com/resourced/resourced/blob/master/tests/data/config-reader/darwin-memory.toml)


**2. Using Go natively.**

ResourceD provides a lot of [native Go readers](https://github.com/resourced/resourced/tree/master/readers).

To use them, define `GoStruct` field with the name of the struct. You can see the full list of legit names [here](https://github.com/resourced/resourced/blob/master/readers/base.go#L12).

At the moment, there's no way to add your own custom Go reader. But we are open to ideas, feel free to submit pull request.
