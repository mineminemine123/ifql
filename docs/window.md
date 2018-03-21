
IFQL has an unanswered design question:

> Should aggregate/selection operations be arguments to the window function? Or should they be their own functions?


This document attempts to capture the existing conversations around this question and to clearly state the implications that each option holds.

### Disclaimer:

I am heavily in favor of the current model as opposed to changing the way the window function works.
I have tried to be objective in my analysis but I am declaring my bias.

If you are in a hurry, there is a one sentence TL;DR at the end of this document.

## The Problem

I postulate that whether window functions accept an aggregate function or whether aggregate functions standalone, is actually a question of which underlying types should be used to represent the data.

As such this question really boils down to the type system we use to model the data on which IFQL operates.
The type system exposes itself in IFQL the language and the query engine.
If we can clearly state what type is at play in any expression in an IFQL script and declare how a function operates on that type then I believe we will have solved this design question.

The current confusion around how the window function operates in relation with aggregate functions stems from the underlying type system not being explicitly defined anywhere, or clearly expressed to users.

I will explain the two type systems at play and how that arises to the two different manifestations in IFQL syntax.
Then we can make a decision on which type system we like best for expressing the operations. 

## The vector type system

The vector type system treats everything as a vector of data. There are two kinds of vectors:

* Column vectors - Column vectors contain a list of values, the values all represent the same point in time and each value comes from a different series.
* Row vectors - Row vectors contain a list of values, the values all represent data from the same series and each value comes from a different time.

A table is either a set of row vectors or a set of column vectors. This means that a table has series along is row dimension and time along its column dimension.

The window function accepts three arguments: 

* a table,
* an aggregate function,
* an interval on which to window the data (this could be a set of arguments defining how data should be windowed, but for simplicity let's assume its a single argument).

The window function returns a table as a set of column vectors.

 The window function does the following:

1. Split each row vector into multiple row vectors based on the window interval
2. Apply the aggregate function to each row vector producing a scalar.
3. Combine the scalars from all series into a column vector for a specific instant in time.
4. Returns a table which is the set of column vectors.

### Example IFQL syntax:

```
from(db:”foo”)
    |> range(start:-1h)
    |> window(every: 1m, fn: (row) => mean(row))
```

Interestingly this model would apply to the group function as well, expect the group function operates on columns:

 
```
from(db:”foo”)
    |> range(start:-1h)
    |> group(by: [“host”, “region”], fn: (col) => mean(col))
```

In this type system the aggregate functions would not operate on tables but only on either row or column vectors.

### Some questions:

> How do you perform both a group by aggregate and window aggregate in the same step?

What if I want to compute the mean RAM usage across all hosts by region over the last 5m.
I could compute the mean by each host first over the last 5m, and then compute the mean of the host means.
But that is not numerically the same result.

If the answer is that group by should not accept an aggregate function then why? Why is window different?

> How are series with multiple values modeled?

I am not sure.

## The block type system

The current implementation of IFQL uses a type system I am calling the block type system because the fundamental type is a block.

A block is a table of data.
The columns are named and represent the time, tags and fields of the data.
The rows represent each record in the dataset.
A block closely resembles a traditional SQL relation (aka table).

A stream is a set of blocks.
A stream may be unbounded meaning that there may be potentially an infinite number of blocks in the set.
The blocks within a stream have not ordering.
NOTE: We currently call a stream a `table`, which is probably contributing to the confusion in the data model.

A stream is data in motion, a block is data at rest.

Every function in IFQL consumes a stream and produces a stream.
There is no relation between the number of blocks that enter a function with the number of blocks that leave a function.

Specifically the window function may consume any number of blocks from the input stream and output any number of blocks on the output stream.

Since the window function is consuming an arbitrary number of blocks and producing an arbitrary number of blocks, the question arises: When do we know a block is complete and can be emitted on the output stream? The answer is to use triggers.
Triggers indicate when a block should be materialized and emitted on the output stream.
The conditions that a trigger can use vary, but for simplicity I’ll explain only the “lagged trigger”.
A lagged trigger “fires” when the current processing time + some fixed amount of lag has elapsed past the event time of the window.
Each outgoing block being produced by the window function has its own trigger.

The window function does the following:

1. Receive an incoming block
2. Assign each row from the incoming block to any number of outgoing blocks.
3. Check the triggers for all outgoing blocks.
4. Emit blocks who’s trigger has “fired”.


The group by function actually does the same thing as the window function except that step 2 uses a different method of assigning input rows to output blocks.
The window function assigns output blocks by the time of the row, the group by function assigns outgoing blocks by the tags of the row.

In fact all functions in this model perform steps 1, 3 and 4.
This is not specific to windowing or grouping. A function need only define its behavior in step 2.

An aggregate/selection function consumes any number of input blocks and always outputs exactly one output block per group represented in the input blocks.
Each incoming block produces a single row in the output block.
The output block time bounds represent the entire time range of the query.
In contrast the grouping bounds of the output block remain unchanged. 

Using the same example as the vector system, this is the current IFQL syntax:

```
from(db:”foo”)
    |> range(start:-1h)
    |> window(every: 1m)
    |> mean()
```

Similarly the group by becomes

 
```
from(db:”foo”)
    |> range(start:-1h)
    |> group(by: [“host”, “region”])
    |> mean()
```

And finally we can now express computing the mean of RAM usage for all hosts by region over the last 5m.

```
from(db:”foo”)
    |> range(start:-1h)
    |> group(by: [“region”])
    |> window(every:5m)
    |> mean()
```


### Some questions:

I think the use of the parameter name `table` has possibly caused some confusion with how the underlying model is understood.
Maybe we should use `stream` since in reality each function is accepting a stream of tables and outputting a stream of tables.
In addition we would rename block to `table` to make it clear.

> Why do aggregate/selection functions produce only a single block per group instead of having a one-to-one relation with input blocks?

I think we need to discuss this more and possibly explore what it would mean to change this so that aggregate functions produce a block for every input block.


## Window/Aggregate with the Block type system

I want to directly address why the design of having the window function accept an aggregate function as an argument is incompatible with the block type system.

First, since the aggregate functions can be expressed in the type system independent of the window function, why would we want to force coupling the concepts.

In addition, if we were to force the coupling then the window function would need to implement internal to itself the logic to handle internal block boundaries and aggregation instead of letting the existing framework handle the logic.


## Concerns with the current model

There have been some great feedback about the current model. I want to address that feedback directly.
The sources of the feedback have been documented below in the NOTES section.

> How do I know when an IFQL function is getting a stream of blocks instead of just a single block?

IFQL functions always receives a stream of blocks, that stream may have many blocks or just one but can always be thought of as a stream.

> What does the window function return if I yield its results directly?

The window function will produce a stream of blocks, each block will represent a single interval defined by the window function.

> How do I perform aggregations along the time dimension vs the series dimension?

An aggregate function always aggregates the values in a block.
Which data is contained within a block can be controlled via the window and group functions.

> Why does an aggregate/selection function change the time bounds of a block but not the series bounds?

We initially thought that it made intuitive sense to remove the time bounds with aggregates, but maybe we need to revisit that decision.

## Comparison

The following is a comparison between the two type systems.
I’ll address various topics for each.

### Inherited Partitions

The block model relies on the user chaining functions together to correctly shape the data such that aggregate functions operate on the desired subset of data.
This makes it so that functions up the chain from the aggregate directly affect how the aggregate is performed.
While in the vector system it’s always coupled how the data is partitioned and aggregated in the same function call.
This means that users of IFQL with the block system need to keep mental track of how they have partitioned the data throughout the query flow.
This is how TICKscript works and has been a point of confusion for some users.
I think if we were to expose via interactive IFQL execution what the data partitions are at each step this would no longer be a significant issue.

### Multiple fields

I haven’t figured out a way to represent a point with multiple fields in the vector system.
In the block system a new field is simply a new column on the block.

### Flexibility in Windowing behaviors

We plan to support more flexible windowing behaviors in the future, beyond simple fixed sized intervals.
For example session based windowing is useful for click stream analytics.
Session based windowing works by defining window boundaries between large gaps in activity, this way each users “session” is in its own window.
With the block model this is straightforward as each session is placed in its own block.
It’s not clear to me how this would work in the vector model since it’s sort of implied that row vectors are the same size, but maybe that constraint could be relaxed.

## Conclusions

Two type systems (aka data models) have been presented and shown how they represent themselves as behaviors in the IFQL syntax.
The vector model closely resembles the Prometheus model.
The block model closely resembles the Apache Beam model, and is the current implementation.

My conclusion is that adding the aggregate functions to the window function represents a type system quite different from the current implementation.
Additionally I also conclude that the current type system has not been well defined and has some inconsistencies that are causing confusion. 

I propose that we address the inconsistencies in the current system instead of changing the underlying system completely.
I believe that if we correctly communicate the type system and expose it via interactive querying that users will be able to understand.

### Recommended Actions

* Rename block to table
* Rename the `table` argument in IFQL functions to `stream`
* Explore what it means to have aggregate/selection functions use a one-to-one block mapping instead of modifying the block time bounds. This may lead to more verbose IFQL scripts, but may also make the block partitioning more visible to the user.
* Explore ways to expose the current partitioning of blocks when using the IFQL REPL.



# TL;DR

Keep IFQL the way it is but improve the vocabulary and visibility into the types being used between IFQL functions so that it is clear what is happening.



# Notes:

This is a list of notes and sources used to compile these ideas.

Original github issue discussing this question from Dan Cech of Grafana after the initial release of IFQL https://github.com/influxdata/ifql/issues/144


Quote from Adam Anthony on the confusion when first learning IFQL:

> @nathaniel a stream to me is a sequence of more or less simple objects. so for IFQL, it's a series of points.
it all sort of collapses for me when sometimes I am getting a stream of tables instead of a stream of [points]


At InfluxDays NYC I (Nathaniel) had a long conversation with Tom Willkie about how the window operation functions.
He was confused about how it partitioned the data.
His perspective was that there are two dimensions on the data: time and series, and that different types of aggregate operations worked on each dimension independently.
This is reflected in Prometheus’s API where they have functions `sum` and `sum_over_time`, where the `sum` function sums along the series dimension and the `sum_over_time` function sums over time.
This is also reflected in the Prometheus type system of instant vectors vs range vectors.


Tyler Akidau talk on a data model to unify stream and table(relational) data models:
 Slides: https://s.apache.org/streaming-sql-big-data-spain
 Video: https://www.youtube.com/watch?v=UlPsp7LaA38 


