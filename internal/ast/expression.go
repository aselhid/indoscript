package ast

type ExprVisitor interface {
	VisitBinaryExpr(expr BinaryExpr) any
	VisitUnaryExpr(expr UnaryExpr) any
	VisitPrimaryExpr(expr PrimaryExpr) any
	VisitGroupExpr(expr GroupExpr) any
	VisitVarExpr(expr VarExpr) any
	VisitLogicalExpr(expr LogicalExpr) any
	VisitCallExpr(expr CallExpr) any
}

type Expr interface {
	Accept(visitor ExprVisitor) any
}

type BinaryExpr struct {
	Left     Expr
	Right    Expr
	Operator Token
}

func (e BinaryExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitBinaryExpr(e)
}

func NewBinaryExpr(left Expr, operator Token, right Expr) BinaryExpr {
	return BinaryExpr{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

type LogicalExpr struct {
	Left     Expr
	Right    Expr
	Operator Token
}

func (e LogicalExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitLogicalExpr(e)
}

func NewLogicalExpr(left Expr, operator Token, right Expr) LogicalExpr {
	return LogicalExpr{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

type UnaryExpr struct {
	Right    Expr
	Operator Token
}

func (e UnaryExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitUnaryExpr(e)
}

func NewUnaryExpr(operator Token, right Expr) UnaryExpr {
	return UnaryExpr{
		Operator: operator,
		Right:    right,
	}
}

type PrimaryExpr struct {
	Literal any
}

func (e PrimaryExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitPrimaryExpr(e)
}

func NewPrimaryExpr(literal any) PrimaryExpr {
	return PrimaryExpr{
		Literal: literal,
	}
}

type GroupExpr struct {
	expression Expr
}

func (e GroupExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitGroupExpr(e)
}

func NewGroupExpr(expression Expr) GroupExpr {
	return GroupExpr{expression: expression}
}

type VarExpr struct {
	Identifier Token
}

func (e VarExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitVarExpr(e)
}

func NewVarExpr(identifier Token) VarExpr {
	return VarExpr{Identifier: identifier}
}

type CallExpr struct {
	Callee      Expr
	Arguments   []Expr
	Parenthesis Token
}

func (e CallExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitCallExpr(e)
}

func NewCallExpr(callee Expr, arguments []Expr, parenthesis Token) CallExpr {
	return CallExpr{
		Callee:      callee,
		Arguments:   arguments,
		Parenthesis: parenthesis,
	}
}
