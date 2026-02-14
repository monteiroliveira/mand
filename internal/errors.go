package internal

type SyntaxError struct{}

func (e *SyntaxError) Error() string {
	return "Syntax Error"
}

func SetSyntaxError() *SyntaxError {
	return &SyntaxError{}
}

type SemanticError struct{}

func (e *SemanticError) Error() string {
	return "Semantic Error"
}

func SetSemanticError() *SemanticError {
	return &SemanticError{}
}
