# Design how meta queries work in IFQL

## SHOW DATABASES

Databases are now buckets and buckets are not solely owned by the query engine.
As such the API tier will provide a method for requesting the list of buckets.

## SHOW RETENTION POLICIES

Retention policies are no longer part of 2.0

## SHOW MEASUREMENTS

A mesurement is now a tag, see SHOW TAG VALUES

## SHOW FIELD KEYS

A field key is now a tag, see SHOW TAG VALUES

## SHOW TAG KEYS

Tag keys are represented as the column labels.
Column have a `kind` property, which is one of `time`,`tag`,`value`.
The `tags` method will produce a new table that has a row for each `tag` column on the input table.

```
from(db:"foo")
    |> group(keep:[*])
    |> tags()
```

Example data for merged table, table is bounded on the interval [0,3).

| _time:time | host:tag | region:tag | _value:value |
| ----- | ---- | ------ | ------ |
| 0     | A    |  east  |  9     |
| 1     | B    |  east  |  5     |
| 2     | A    |  west  |  3     |

Example data for tags table, table is bounded on the interval [0,3).

| _time | _value |
| ----- | ------ |
| 3     | host   |
| 3     | region |

NOTE: There is no row corresponding to the `_time` and `_value` columns, this is because those columns are not `tag` columns as defined by their `kind`.


## SHOW TAG VALUES

Producing a list of tag values is a `disticnt` operation.

```
SHOW TAG VALUES ON "foo" WITH KEY = "host"
```

```
from(db:"foo")
    |> group(keep:["host"])
    |> distinct(column:"host")
```

Example data for merged table, table is bounded on the interval [0,3).

| _time | host | region | _value |
| ----- | ---- | ------ | ------ |
| 0     | A    |  east  |  9     |
| 1     | B    |  east  |  5     |
| 2     | A    |  west  |  3     |

Example data for distinct table, table is bounded on the interval [0,3).

| _time | _value |
| ----- | ------ |
| 3     | A      |
| 3     | B      |

The WHERE clause becomes a filter step.

```
SHOW TAG VALUES ON "foo" WITH KEY = "host" WHERE "region" = 'east'
```

```
from(db:"foo")
    |> filter(fn:(r) => r.region = "east")
    |> group(keep:["host"])
    |> distinct(column:"host")
```


```
SHOW MEASUREMENTS
```

```
from(db:"foo")
    |> group(keep:["_measurement"])
    |> distinct(column:"_measurement")
```

```
SHOW FIELD KEYS
```

```
from(db:"foo")
    |> group(keep:["_field"])
    |> distinct(column:"_field")
```

Now it is also possible to perform these meta queries for bounded time intervals

```
from(db:"foo")
    |> group(keep:["_field"])
    |> range(start:-1m)
    |> distinct(column:"_field")
```

We can define these helper functions:

```
tagValues = (tag, table=<-) =>
    table
        |> group(keep:[tag])
        |> distinct(column:tag)

fieldKeys = (table=<-) => table |> tagValues(tag:"_field")
measurements = (table=<-) => table |> tagValues(tag:"_measurement")
```

Question: Introducing time means that the actual data will be read and a distinct operation performed on the data.
With time unbounded queries the index can be consulted to produce the tag key and values.
What physical operations does the storage support for quering this data?



## SHOW SERIES

Do we still need this in IFQL? I am not sure what the use case for this lower level meta query is in InfluxQL.


