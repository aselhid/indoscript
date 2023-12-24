package ast

type StmtVisitor interface {
	VisitPrintStmt(stmt PrintStmt)
	VisitExprStmt(stmt ExprStmt)
	VisitVarStmt(stmt VarStmt)
	VisitBlockStmt(stmt BlockStmt)
	VisitIfStmt(stmt IfStmt)
	VisitWhileStmt(stmt WhileStmt)
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

type BlockStmt struct {
	Statements []Stmt
}

func (s BlockStmt) Accept(visitor StmtVisitor) {
	visitor.VisitBlockStmt(s)
}

func NewBlockStmt(statements []Stmt) BlockStmt {
	return BlockStmt{Statements: statements}
}

type IfStmt struct {
	Condition Expr
	ThenStmt  BlockStmt
	ElseStmt  BlockStmt
}

func (s IfStmt) Accept(visitor StmtVisitor) {
	visitor.VisitIfStmt(s)
}

func NewIfStmt(condition Expr, thenStmt BlockStmt, elseStmt BlockStmt) IfStmt {
	return IfStmt{
		Condition: condition,
		ThenStmt:  thenStmt,
		ElseStmt:  elseStmt,
	}
}

type WhileStmt struct {
	Condition Expr
	Stmt      BlockStmt
}

func (s WhileStmt) Accept(visitor StmtVisitor) {
	visitor.VisitWhileStmt(s)
}

func NewWhileStmt(condition Expr, blockStmt BlockStmt) WhileStmt {
	return WhileStmt{
		Condition: condition,
		Stmt:      blockStmt,
	}
}
