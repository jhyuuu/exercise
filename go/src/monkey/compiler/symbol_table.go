package compiler

type SymbolScope string

const (
    GlobalScope SymbolScope = "GLOBAL"
)

type Symbol struct {
    Name  string
    Index int
    Scope SymbolScope
}

type SymbolTable struct {
    store         map[string]Symbol
    numDefinition int
}

func NewSymbolTable() *SymbolTable {
    s := make(map[string]Symbol)
    return &SymbolTable{store: s}
}

func (s *SymbolTable) Define(name string) Symbol {
    symbol := Symbol{Name: name, Index: s.numDefinition, Scope: GlobalScope}
    s.store[name] = symbol
    s.numDefinition++
    return symbol
}

func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
    symbol, ok := s.store[name]
    return symbol, ok
}