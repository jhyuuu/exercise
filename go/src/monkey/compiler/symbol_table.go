package compiler

type SymbolScope string

const (
    LocalScope  SymbolScope = "LOCAL"
    GlobalScope SymbolScope = "GLOBAL"
)

type Symbol struct {
    Name  string
    Index int
    Scope SymbolScope
}

type SymbolTable struct {
    Outer         *SymbolTable

    store         map[string]Symbol
    numDefinition int
}

func NewSymbolTable() *SymbolTable {
    s := make(map[string]Symbol)
    return &SymbolTable{store: s}
}

func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
    s := NewSymbolTable()
    s.Outer = outer
    return s
}

func (s *SymbolTable) Define(name string) Symbol {
    symbol := Symbol{Name: name, Index: s.numDefinition}
    if s.Outer == nil {
        symbol.Scope = GlobalScope
    } else {
        symbol.Scope = LocalScope
    }

    s.store[name] = symbol
    s.numDefinition++
    return symbol
}

func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
    obj, ok := s.store[name]
    if !ok && s.Outer != nil {
        obj, ok = s.Outer.Resolve(name)
        return obj, ok
    }
    return obj, ok
}