Package `queryparser` parses the `Conditions = xyz` query.

## Spec
```
# Given data:
#   /r/load-avg: {"Data": {"LoadAvg1m": 0.904296875}}

# 1. Every part must be defined inside parenthesis.
((/r/load-avg.LoadAvg1m > 0.5) && (/r/load-avg.LoadAvg1m < 10))

# 2. Boolean operator: &&, ||
# 3. Numerical operator: typical
# 4. String operator: ==, !=
```