## Authentication

You can protect the ResourceD endpoints behind token-based authentication.

To use this feature, simply generate a long hash and put it under `access-tokens` directory.

You can define one token per file, or one token per line, or multiple tokens separated by comma per line.

Examples: https://github.com/resourced/resourced/blob/master/tests/resourced-configs/access-tokens
