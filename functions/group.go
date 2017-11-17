package functions

import (
	"errors"
	"fmt"
	"sort"

	"github.com/influxdata/ifql/ifql"
	"github.com/influxdata/ifql/query"
	"github.com/influxdata/ifql/query/execute"
	"github.com/influxdata/ifql/query/plan"
)

const GroupKind = "group"

type GroupOpSpec struct {
	By     []string `json:"by"`
	Keep   []string `json:"keep"`
	Ignore []string `json:"ignore"`
}

func init() {
	ifql.RegisterFunction(GroupKind, createGroupOpSpec)
	query.RegisterOpSpec(GroupKind, newGroupOp)
	plan.RegisterProcedureSpec(GroupKind, newGroupProcedure, GroupKind)
	execute.RegisterTransformation(GroupKind, createGroupTransformation)
}

func createGroupOpSpec(args map[string]ifql.Value, ctx ifql.Context) (query.OperationSpec, error) {
	spec := new(GroupOpSpec)
	if len(args) == 0 {
		return spec, nil
	}

	if value, ok := args["by"]; ok {
		if value.Type != ifql.TArray {
			return nil, fmt.Errorf("'by' argument must be a list of strings got %v", value.Type)
		}
		list := value.Value.(ifql.Array)
		if list.Type != ifql.TString {
			return nil, fmt.Errorf("'by' argument must be a list of strings, got list of %v", list.Type)
		}
		spec.By = list.Elements.([]string)
	}

	if value, ok := args["keep"]; ok {
		if value.Type != ifql.TArray {
			return nil, fmt.Errorf("keep argument must be a list of strings got %v", value.Type)
		}
		list := value.Value.(ifql.Array)
		if list.Type != ifql.TString {
			return nil, fmt.Errorf("keep argument must be a list of strings, got list of %v", list.Type)
		}
		spec.Keep = list.Elements.([]string)
	}
	if value, ok := args["ignore"]; ok {
		if value.Type != ifql.TArray {
			return nil, fmt.Errorf("ignore argument must be a list of strings got %v", value.Type)
		}
		list := value.Value.(ifql.Array)
		if list.Type != ifql.TString {
			return nil, fmt.Errorf("ignore argument must be a list of strings, got list of %v", list.Type)
		}
		spec.Ignore = list.Elements.([]string)
	}
	if len(spec.By) > 0 && len(spec.Ignore) > 0 {
		return nil, errors.New("cannot specify both by and ignore keys")
	}
	return spec, nil
}

func newGroupOp() query.OperationSpec {
	return new(GroupOpSpec)
}

func (s *GroupOpSpec) Kind() query.OperationKind {
	return GroupKind
}

type GroupProcedureSpec struct {
	By     []string
	Ignore []string
	Keep   []string
}

func newGroupProcedure(qs query.OperationSpec) (plan.ProcedureSpec, error) {
	spec, ok := qs.(*GroupOpSpec)
	if !ok {
		return nil, fmt.Errorf("invalid spec type %T", qs)
	}

	p := &GroupProcedureSpec{
		By:     spec.By,
		Ignore: spec.Ignore,
		Keep:   spec.Keep,
	}
	return p, nil
}

func (s *GroupProcedureSpec) Kind() plan.ProcedureKind {
	return GroupKind
}
func (s *GroupProcedureSpec) Copy() plan.ProcedureSpec {
	ns := new(GroupProcedureSpec)

	ns.By = make([]string, len(s.By))
	copy(ns.By, s.By)

	ns.Ignore = make([]string, len(s.Ignore))
	copy(ns.Ignore, s.Ignore)

	ns.Keep = make([]string, len(s.Keep))
	copy(ns.Keep, s.Keep)

	return ns
}

func (s *GroupProcedureSpec) PushDownRule() plan.PushDownRule {
	return plan.PushDownRule{
		Root:    SelectKind,
		Through: []plan.ProcedureKind{LimitKind, RangeKind, FilterKind},
		Match: func(root *plan.Procedure) bool {
			selectSpec := root.Spec.(*SelectProcedureSpec)
			return !selectSpec.AggregateSet
		},
	}
}

func (s *GroupProcedureSpec) PushDown(root *plan.Procedure, dup func() *plan.Procedure) {
	selectSpec := root.Spec.(*SelectProcedureSpec)
	if selectSpec.GroupingSet {
		root = dup()
		selectSpec = root.Spec.(*SelectProcedureSpec)
		selectSpec.OrderByTime = false
		selectSpec.GroupingSet = false
		selectSpec.MergeAll = false
		selectSpec.GroupKeys = nil
		selectSpec.GroupIgnore = nil
		selectSpec.GroupKeep = nil
		return
	}
	selectSpec.GroupingSet = true
	// TODO implement OrderByTime
	//selectSpec.OrderByTime = true

	// Merge all series into a single group if we have no specific grouping dimensions.
	selectSpec.MergeAll = len(s.By) == 0 && len(s.Ignore) == 0
	selectSpec.GroupKeys = s.By
	selectSpec.GroupIgnore = s.Ignore
	selectSpec.GroupKeep = s.Keep
}

func createGroupTransformation(id execute.DatasetID, mode execute.AccumulationMode, spec plan.ProcedureSpec, ctx execute.Context) (execute.Transformation, execute.Dataset, error) {
	s, ok := spec.(*GroupProcedureSpec)
	if !ok {
		return nil, nil, fmt.Errorf("invalid spec type %T", spec)
	}
	cache := execute.NewBlockBuilderCache()
	d := execute.NewDataset(id, mode, cache)
	t := NewGroupTransformation(d, cache, s)
	return t, d, nil
}

type groupTransformation struct {
	d     execute.Dataset
	cache execute.BlockBuilderCache

	keys   []string
	ignore []string
	keep   []string

	// Ignoring is true of len(keys) == 0 && len(ignore) > 0
	ignoring bool
}

func NewGroupTransformation(d execute.Dataset, cache execute.BlockBuilderCache, spec *GroupProcedureSpec) *groupTransformation {
	t := &groupTransformation{
		d:        d,
		cache:    cache,
		keys:     spec.By,
		ignore:   spec.Ignore,
		keep:     spec.Keep,
		ignoring: len(spec.By) == 0 && len(spec.Ignore) > 0,
	}
	sort.Strings(t.keys)
	sort.Strings(t.ignore)
	sort.Strings(t.keep)
	return t
}

func (t *groupTransformation) RetractBlock(id execute.DatasetID, meta execute.BlockMetadata) {
	//TODO(nathanielc): Investigate if this can be smarter and not retract all blocks with the same time bounds.
	t.cache.ForEachBuilder(func(bk execute.BlockKey, builder execute.BlockBuilder) {
		if meta.Bounds().Equal(builder.Bounds()) {
			t.d.RetractBlock(bk)
		}
	})
}

func (t *groupTransformation) Process(id execute.DatasetID, b execute.Block) {
	isFanIn := false
	var tags execute.Tags
	if t.ignoring {
		// Assume we can fan in, we check for the false condition below
		isFanIn = true
		blockTags := b.Tags()
		tags = make(execute.Tags, len(blockTags))
		cols := b.Cols()
		for _, c := range cols {
			if c.IsTag {
				found := false
				for _, tag := range t.ignore {
					if tag == c.Label {
						found = true
						break
					}
				}
				if !found {
					if !c.IsCommon {
						isFanIn = false
						break
					}
					tags[c.Label] = blockTags[c.Label]
				}
			}
		}
	} else {
		tags, isFanIn = b.Tags().Subset(t.keys)
	}
	if isFanIn {
		t.processFanIn(b, tags)
	} else {
		t.processFanOut(b)
	}
}

// processFanIn assumes that all rows of b will be placed in the same builder.
func (t *groupTransformation) processFanIn(b execute.Block, tags execute.Tags) {
	builder, new := t.cache.BlockBuilder(blockMetadata{
		tags:   tags,
		bounds: b.Bounds(),
	})
	if new {
		// Determine columns of new block.

		// Add existing columns, skipping tags.
		for _, c := range b.Cols() {
			if !c.IsTag {
				builder.AddCol(c)
			}
		}

		// Add tags.
		execute.AddTags(tags, builder)

		// Add columns for tags that are to be kept.
		for _, k := range t.keep {
			builder.AddCol(execute.ColMeta{
				Label: k,
				Type:  execute.TString,
				IsTag: true,
			})
		}
	}

	// Construct map of builder column index to block column index.
	builderCols := builder.Cols()
	blockCols := b.Cols()
	colMap := make([]int, len(builderCols))
	for j, c := range builderCols {
		for nj, nc := range blockCols {
			if c.Label == nc.Label {
				colMap[j] = nj
				break
			}
		}
	}

	execute.AppendBlock(b, builder, colMap)
}

type tagMeta struct {
	idx      int
	isCommon bool
}

// processFanOut assumes each row of b could end up in a different builder.
func (t *groupTransformation) processFanOut(b execute.Block) {
	cols := b.Cols()
	tagMap := make(map[string]tagMeta, len(cols))
	for j, c := range cols {
		if c.IsTag {
			ignoreTag := false
			for _, tag := range t.ignore {
				if tag == c.Label {
					ignoreTag = true
					break
				}
			}
			byTag := false
			for _, tag := range t.keys {
				if tag == c.Label {
					byTag = true
					break
				}
			}
			keepTag := false
			for _, tag := range t.keep {
				if tag == c.Label {
					keepTag = true
					break
				}
			}
			if (t.ignoring && !ignoreTag) || byTag || keepTag {
				tagMap[c.Label] = tagMeta{
					idx:      j,
					isCommon: (t.ignoring && !keepTag) || (!t.ignoring && byTag),
				}
			}
		}
	}

	// Iterate over each row and append to specific builder
	b.Times().DoTime(func(ts []execute.Time, rr execute.RowReader) {
		for i := range ts {
			tags := t.determineRowTags(tagMap, i, rr)
			builder, new := t.cache.BlockBuilder(blockMetadata{
				tags:   tags,
				bounds: b.Bounds(),
			})
			if new {
				// Add existing columns, skipping tags.
				for _, c := range cols {
					if !c.IsTag {
						builder.AddCol(c)
						continue
					}
					if meta, ok := tagMap[c.Label]; ok {
						j := builder.AddCol(execute.ColMeta{
							Label:    c.Label,
							Type:     execute.TString,
							IsTag:    true,
							IsCommon: meta.isCommon,
						})
						if meta.isCommon {
							builder.SetCommonString(j, tags[c.Label])
						}
					}
				}
			}
			// Construct map of builder column index to block column index.
			builderCols := builder.Cols()
			colMap := make([]int, len(builderCols))
			for j, c := range builderCols {
				for nj, nc := range cols {
					if c.Label == nc.Label {
						colMap[j] = nj
						break
					}
				}
			}

			// Add row to builder
			execute.AppendRow(i, rr, builder, colMap)
		}
	})
}

func (t *groupTransformation) determineRowTags(tagMap map[string]tagMeta, i int, rr execute.RowReader) execute.Tags {
	cols := rr.Cols()
	tags := make(execute.Tags, len(cols))
	for t, meta := range tagMap {
		if meta.isCommon {
			tags[t] = rr.AtString(i, meta.idx)
		}
	}
	return tags
}

func (t *groupTransformation) UpdateWatermark(id execute.DatasetID, mark execute.Time) {
	t.d.UpdateWatermark(mark)
}
func (t *groupTransformation) UpdateProcessingTime(id execute.DatasetID, pt execute.Time) {
	t.d.UpdateProcessingTime(pt)
}
func (t *groupTransformation) Finish(id execute.DatasetID, err error) {
	t.d.Finish(err)
}
func (t *groupTransformation) SetParents(ids []execute.DatasetID) {
}

type blockMetadata struct {
	tags   execute.Tags
	bounds execute.Bounds
}

func (m blockMetadata) Tags() execute.Tags {
	return m.tags
}
func (m blockMetadata) Bounds() execute.Bounds {
	return m.bounds
}
