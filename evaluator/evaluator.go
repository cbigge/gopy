package evaluator

import (
	"fmt"
	"gopy/ast"
	"gopy/interpreter"
)

var (
	TRUE = &interpreter.Bool{Val: true}
	FALSE = &interpreter.Bool{Val: false}
)

var builtins = map[string]*interpreter.Builtin{
	"print": {
		Fn: func(args ...interpreter.Item) interpreter.Item {
			var result string
			for _, arg := range args {
				result += arg.Visit()
			}
			return &interpreter.Str{Val: result}
		},
	},
}

func Evaluate(node ast.Node, env *interpreter.Environment) interpreter.Item {
	switch node := node.(type) {
	case *ast.Program:
		return evaluateStmts(node.Stmts, env)
	case *ast.ExprStmt:
		return Evaluate(node.Expr, env)
	case *ast.CallExpr:
		fn := Evaluate(node.Func, env)
		if fn.Type() == interpreter.ERR {
			return fn
		}
		args := evaluateExprs(node.Args, env)
		if len(args) == 1 && args[0].Type() == interpreter.ERR {
			return args[0]
		}
		//return applyFn(fn, args)
	case *ast.VarStmt:
		v := Evaluate(node.Value, env)
		if v.Type() == interpreter.ERR {
			return v
		}
		env.Store(node.Ident.Val, v)
		return v
	case *ast.Identifier:
		return evaluateIdent(node, env)
	case *ast.PrefixExpr:
		expr := Evaluate(node.Expr, env)
		if expr.Type() == interpreter.ERR {
			return expr
		}
		return evaluatePrefixExpr(node.Op, expr)
	case *ast.InfixExpr:
		l := Evaluate(node.Left, env)
		if l.Type() == interpreter.ERR {
			return l
		}
		r := Evaluate(node.Right, env)
		if r.Type() == interpreter.ERR {
			return r
		}
		return evaluateInfixExpr(node.Op, l, r)
	case *ast.BlockStmt:
		return evaluateStmts(node.Stmts, env)
	case *ast.IfExpr:
		return evaluateIfExpr(node, env)
	case *ast.IntLiteral:
		return &interpreter.Int{Val: node.Value}
	case *ast.StrLiteral:
		return &interpreter.Str{Val: node.Value}
	}
	return nil
}

func applyFn(fn interpreter.Item, args []interpreter.Item) interpreter.Item {
	fun, ok := fn.(*interpreter.Builtin)
	if !ok {
		return newErr("not a function: %s", fn.Type())
	}
	return fun.Fn(args...)
}

func evaluateStmts(stmts []ast.Stmt, env *interpreter.Environment) interpreter.Item {
	var result interpreter.Item
	for _, stmt := range stmts {
		result = Evaluate(stmt, env)
	}
	return result
}

func evaluateExprs(e []ast.Expr, env *interpreter.Environment) []interpreter.Item {
	var result []interpreter.Item
	for _, expr := range e {
		eval := Evaluate(expr, env)
		if eval.Type() == interpreter.ERR {
			return []interpreter.Item{eval}
		}
		result = append(result, eval)
	}
	return result
}

func evaluateIdent(i *ast.Identifier, env *interpreter.Environment) interpreter.Item {
	if val, ok := env.Get(i.Val); ok {
		return val
	}
	if builtin, ok := builtins[i.Val]; ok {
		return builtin
	}
	return newErr("identifier not found: " + i.Val)
}

func evaluatePrefixExpr(op string, expr interpreter.Item) interpreter.Item {
	switch op {
	case "-":
		return evaluateNegateOpExpr(expr)
	default:
		return newErr("unknown operator: %s%s", op, expr.Type())
	}
}

func evaluateNegateOpExpr(expr interpreter.Item) interpreter.Item {
	if expr.Type() != interpreter.INT {
		return nil
	}
	val := expr.(*interpreter.Int).Val
	return &interpreter.Int{Val: -val}
}

func evaluateInfixExpr(op string, l interpreter.Item, r interpreter.Item) interpreter.Item {
	switch {
	case l.Type() == interpreter.INT && r.Type() == interpreter.INT:
		return evaluateIntInfixExpr(op, l, r)
	case l.Type() == interpreter.STR && r.Type() == interpreter.STR,
			l.Type() == interpreter.INT && r.Type() == interpreter.STR,
			l.Type() == interpreter.STR && r.Type() == interpreter.INT:
		return evaluateStrInfixExpr(op, l, r)
	default:
		return newErr("unknown operator: %s %s %s", l.Type(), op, r.Type())
	}
}

func evaluateIntInfixExpr(op string, l interpreter.Item, r interpreter.Item) interpreter.Item {
	left := l.(*interpreter.Int).Val
	right := r.(*interpreter.Int).Val
	switch op {
	case "+":
		return &interpreter.Int{Val: left+right}
	case "-":
		return &interpreter.Int{Val: left-right}
	case "*":
		return &interpreter.Int{Val: left*right}
	case "/":
		return &interpreter.Int{Val: left/right}
	case "<":
		if left < right {
			return TRUE
		} else {
			return FALSE
		}
	case ">":
		if left > right {
			return TRUE
		} else {
			return FALSE
		}
	case "==":
		if left == right {
			return TRUE
		} else {
			return FALSE
		}
	case "!=":
		if left != right {
			return TRUE
		} else {
			return FALSE
		}
	default:
		return newErr("unknown operator: %s", op)
	}
}

func evaluateStrInfixExpr(op string, l interpreter.Item, r interpreter.Item) interpreter.Item {
	left := l.Visit()
	right := r.Visit()
	switch op {
	case "+":
		return &interpreter.Str{Val: fmt.Sprintf("%s%s", left,right)}
	default:
		return newErr("unknown operator: %s", op)
	}
}

func evaluateIfExpr(ie *ast.IfExpr, env *interpreter.Environment) interpreter.Item {
	cond := Evaluate(ie.Cond, env)
	if isTrue(cond) {
		return Evaluate(ie.Pass, env)
	} else if ie.Fail != nil {
		return Evaluate(ie.Fail, env)
	} else {
		return nil
	}
}

func isTrue(item interpreter.Item) bool {
	switch item {
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return false
	}
}

func newErr(f string, e ...interface{}) *interpreter.Error {
	return &interpreter.Error{Err: fmt.Sprintf(f, e...)}
}