package ast

type StmtVisitor interface {
	VisitPrintStmt(stmt PrintStmt)
	VisitExprStmt(stmt ExprStmt)
	VisitVarStmt(stmt VarStmt)
}

type Stmt interface {
	Accept(visitor StmtVisitor)
}

type PrintStmt struct {
	Expression Expr
}

func (s PrintStmt) Accept(visitor StmtVisitor) {
	visitor.VisitPrintStmt(s)
}

func NewPrintStmt(expression Expr) PrintStmt {
	return PrintStmt{Expression: expression}
}

type ExprStmt struct {
	Expression Expr
}

func (s ExprStmt) Accept(visitor StmtVisitor) {
	visitor.VisitExprStmt(s)
}

func NewExprStmt(expression Expr) ExprStmt {
	return ExprStmt{Expression: expression}
}

type VarStmt struct {
	Identifier Token
	Expression Expr
}

func (s VarStmt) Accept(visitor StmtVisitor) {
	visitor.VisitVarStmt(s)
}

func NewVarStmt(identifier Token, expression Expr) VarStmt {
	return VarStmt{Identifier: identifier, Expression: expression}
}
