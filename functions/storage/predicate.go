package storage

import (
	"fmt"

	"github.com/influxdata/ifql/ast"
	"github.com/influxdata/ifql/functions/storage/internal/pb"
	"github.com/influxdata/ifql/semantic"
	"github.com/pkg/errors"
)

func ToStoragePredicate(f *semantic.FunctionExpression) (*pb.Predicate, error) {
	if len(f.Params) != 1 {
		return nil, errors.New("storage predicate functions must have exactly one parameter")
	}

	root, err := toStoragePredicate(f.Body.(semantic.Expression), f.Params[0].Key.Name)
	if err != nil {
		return nil, err
	}

	return &pb.Predicate{
		Root: root,
	}, nil
}

func toStoragePredicate(n semantic.Expression, objectName string) (*pb.Node, error) {
	switch n := n.(type) {
	case *semantic.LogicalExpression:
		left, err := toStoragePredicate(n.Left, objectName)
		if err != nil {
			return nil, errors.Wrap(err, "left hand side")
		}
		right, err := toStoragePredicate(n.Right, objectName)
		if err != nil {
			return nil, errors.Wrap(err, "right hand side")
		}
		children := []*pb.Node{left, right}
		switch n.Operator {
		case ast.AndOperator:
			return &pb.Node{
				NodeType: pb.NodeTypeLogicalExpression,
				Value:    &pb.Node_Logical_{Logical: pb.LogicalAnd},
				Children: children,
			}, nil
		case ast.OrOperator:
			return &pb.Node{
				NodeType: pb.NodeTypeLogicalExpression,
				Value:    &pb.Node_Logical_{Logical: pb.LogicalOr},
				Children: children,
			}, nil
		default:
			return nil, fmt.Errorf("unknown logical operator %v", n.Operator)
		}
	case *semantic.BinaryExpression:
		left, err := toStoragePredicate(n.Left, objectName)
		if err != nil {
			return nil, errors.Wrap(err, "left hand side")
		}
		right, err := toStoragePredicate(n.Right, objectName)
		if err != nil {
			return nil, errors.Wrap(err, "right hand side")
		}
		children := []*pb.Node{left, right}
		op, err := toComparisonOperator(n.Operator)
		if err != nil {
			return nil, err
		}
		return &pb.Node{
			NodeType: pb.NodeTypeComparisonExpression,
			Value:    &pb.Node_Comparison_{Comparison: op},
			Children: children,
		}, nil
	case *semantic.StringLiteral:
		return &pb.Node{
			NodeType: pb.NodeTypeLiteral,
			Value: &pb.Node_StringValue{
				StringValue: n.Value,
			},
		}, nil
	case *semantic.IntegerLiteral:
		return &pb.Node{
			NodeType: pb.NodeTypeLiteral,
			Value: &pb.Node_IntegerValue{
				IntegerValue: n.Value,
			},
		}, nil
	case *semantic.BooleanLiteral:
		return &pb.Node{
			NodeType: pb.NodeTypeLiteral,
			Value: &pb.Node_BooleanValue{
				BooleanValue: n.Value,
			},
		}, nil
	case *semantic.FloatLiteral:
		return &pb.Node{
			NodeType: pb.NodeTypeLiteral,
			Value: &pb.Node_FloatValue{
				FloatValue: n.Value,
			},
		}, nil
	case *semantic.RegexpLiteral:
		return &pb.Node{
			NodeType: pb.NodeTypeLiteral,
			Value: &pb.Node_RegexValue{
				RegexValue: n.Value.String(),
			},
		}, nil
	case *semantic.MemberExpression:
		// Sanity check that the object is the objectName identifier
		if ident, ok := n.Object.(*semantic.IdentifierExpression); !ok || ident.Name != objectName {
			return nil, fmt.Errorf("unknown object %q", n.Object)
		}
		if n.Property == "_value" {
			return &pb.Node{
				NodeType: pb.NodeTypeFieldRef,
				Value: &pb.Node_FieldRefValue{
					FieldRefValue: "_value",
				},
			}, nil
		}
		return &pb.Node{
			NodeType: pb.NodeTypeTagRef,
			Value: &pb.Node_TagRefValue{
				TagRefValue: n.Property,
			},
		}, nil
	case *semantic.DurationLiteral:
		return nil, errors.New("duration literals not supported in storage predicates")
	case *semantic.DateTimeLiteral:
		return nil, errors.New("time literals not supported in storage predicates")
	default:
		return nil, fmt.Errorf("unsupported semantic expression type %T", n)
	}
}

func toComparisonOperator(o ast.OperatorKind) (pb.Node_Comparison, error) {
	switch o {
	case ast.EqualOperator:
		return pb.ComparisonEqual, nil
	case ast.NotEqualOperator:
		return pb.ComparisonNotEqual, nil
	case ast.RegexpMatchOperator:
		return pb.ComparisonRegex, nil
	case ast.NotRegexpMatchOperator:
		return pb.ComparisonNotRegex, nil
	case ast.StartsWithOperator:
		return pb.ComparisonStartsWith, nil
	case ast.LessThanOperator:
		return pb.ComparisonLess, nil
	case ast.LessThanEqualOperator:
		return pb.ComparisonLessEqual, nil
	case ast.GreaterThanOperator:
		return pb.ComparisonGreater, nil
	case ast.GreaterThanEqualOperator:
		return pb.ComparisonGreaterEqual, nil
	default:
		return 0, fmt.Errorf("unknown operator %v", o)
	}
}
