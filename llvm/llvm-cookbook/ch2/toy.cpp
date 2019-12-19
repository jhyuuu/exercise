#include <cstdio>
#include <map>
#include <string>
#include <vector>

#include "llvm/IR/Verifier.h"
#include "llvm/IR/DerivedTypes.h"
#include "llvm/IR/IRBuilder.h"
#include "llvm-c/Core.h"
#include "llvm/IR/Module.h"

using namespace llvm;

enum Token_Type {
    EOF_TOKEN = 0,
    DEF_TOKEN,
    IDENTIFIER_TOKEN,
    NUMERIC_TOKEN
};

FILE* file;
static int Numeric_Val;
static std::string Identifier_string;

static int get_token() {
    static int LastChar = ' ';

    while (isspace(LastChar)) {
        LastChar = fgetc(file);
    }

    if (isalpha(LastChar)) {
        Identifier_string = LastChar;
        while (isalnum((LastChar = fgetc(file)))) {
            Identifier_string = LastChar;
        }

        if (Identifier_string == "def")
            return DEF_TOKEN;

        return IDENTIFIER_TOKEN;
    }

    if (isdigit(LastChar)) {
        std::string NumStr;
        do {
            NumStr += LastChar;
            LastChar = fgetc(file);
        } while (isdigit(LastChar));

        Numeric_Val = strtod(NumStr.c_str(), 0);
        return NUMERIC_TOKEN;
    }

    if (LastChar == '#') {
        do {
            LastChar = fgetc(file);
        } while (LastChar != EOF && LastChar != '\n' && LastChar != '\r');

        if (LastChar != EOF)
            return get_token();
    }

    if (LastChar == EOF)
        return EOF_TOKEN;

    int ThisChar = LastChar;
    LastChar = fgetc(file);
    return ThisChar;
}

namespace {

class BaseAST {
  public:
    virtual ~BaseAST() { }
};

class VariableAST : public BaseAST {
  public:
    VariableAST(const std::string& name) : Var_Name(name) { }
  private:
    std::string Var_Name;
};

class NumericAST : public BaseAST {
  public:
    NumericAST(int val) : numeric_val(val) { }
  private:
    int numeric_val;
};

class BinaryAST : public BaseAST {
  public:
    BinaryAST(const std::string& op, BaseAST* lhs, BaseAST* rhs)
        : Bin_Operator(op), LHS(lhs), RHS(rhs) { }
  private:
    std::string Bin_Operator;
    BaseAST* LHS;
    BaseAST* RHS;
};

class FunctionDeclAST : public BaseAST {
  public:
    FunctionDeclAST(const std::string& name, const std::vector<std::string>& args)
        : Func_Name(name), Arguments(args) { }
  private:
    std::string Func_Name;
    std::vector<std::string> Arguments;
};

class FunctionDefnAST : public BaseAST {
  public:
    FunctionDefnAST(FunctionDeclAST* proto, BaseAST* body)
        : Func_Decl(proto), Body(body) { }
  private:
    FunctionDeclAST* Func_Decl;
    BaseAST* Body;
};

class FunctionCallAST : public BaseAST {
  public:
    FunctionCallAST(const std::string& callee, const std::vector<BaseAST*>& args)
        : Function_Callee(callee), Function_Arguments(args) { }
  private:
    std::string Function_Callee;
    std::vector<BaseAST*> Function_Arguments;
};

} // end namespace

static int Current_token;
static int next_token() {
    return Current_token = get_token();
}
static std::map<char, int> Operator_Precedence;
static int getBinOpPrecedence() {
    if(!isascii(Current_token))
      return -1;
	
    if (Operator_Precedence.count(Current_token) != 0) {
        int TokPrec = Operator_Precedence[Current_token];
        return TokPrec;
    }
    return -1;
}

static BaseAST* expression_parser();

static BaseAST* identifier_parser() {
    std::string IdName = Identifier_string;

    next_token();

    // variable
    if (Current_token != '(')
        return new VariableAST(IdName);

    // function call
    next_token();
    
    std::vector<BaseAST*> Args;
    if (Current_token != ')') {
        while (1) {
            BaseAST* Arg = expression_parser();
            if (Arg == NULL) return NULL;
            Args.push_back(Arg);

            if (Current_token == ')') break;
            if (Current_token != ',') return NULL;

            next_token();
        }
    }

    next_token();

    return new FunctionCallAST(IdName, Args);
}

static BaseAST* numeric_parser() {
    BaseAST* result = new  NumericAST(Numeric_Val);
    next_token();
    return result;
}

static FunctionDeclAST* func_decl_parser() {
    if (Current_token != DEF_TOKEN)
        return NULL;

    std::string FnName = Identifier_string;

    next_token();

    if (Current_token != '(') return NULL;

    std::vector<std::string> arg_names;
    while (next_token() == IDENTIFIER_TOKEN) {
        // TODO: how to parser ','
        arg_names.push_back(Identifier_string);
    }

    if (Current_token != ')') return NULL;

    next_token();

    return new FunctionDeclAST(FnName, arg_names);
}

static FunctionDefnAST* func_defn_parser() {
    FunctionDeclAST* Decl = func_decl_parser();
    if (Decl == NULL) return NULL;

    BaseAST* Body = expression_parser();
    if (Body != NULL) {
        return new FunctionDefnAST(Decl, Body);
    }

    next_token();

    return NULL;
}

static FunctionDefnAST* top_level_parser() {
    BaseAST* E = expression_parser();
    if (E != NULL) {
        FunctionDeclAST *Func_Decl = new FunctionDeclAST("", std::vector<std::string>());
        return new FunctionDefnAST(Func_Decl, E);
    }
    return NULL;
}

static BaseAST* paran_parser() {
    next_token();
    BaseAST* expr = expression_parser();
    if (expr == NULL) return NULL;

    if (Current_token != ')')
        return NULL;

    return expr;
}

static BaseAST* Base_Parser() {
    switch (Current_token) {
        case IDENTIFIER_TOKEN: return identifier_parser();
        case NUMERIC_TOKEN:    return numeric_parser();
        case '(':              return paran_parser();
        default: return NULL;
    }
}

static BaseAST* binary_op_parser(int Old_Prec, BaseAST* LHS) {
    while (1) {
        int Operator_Prec = getBinOpPrecedence();

        if (Operator_Prec < Old_Prec) {
            return LHS;
        }

        int BinOp = Current_token;

        next_token();

        BaseAST* RHS = Base_Parser();
        if (RHS == NULL) return NULL;

        int Next_Prec = getBinOpPrecedence();
        if (Operator_Prec < Next_Prec) {
            RHS = binary_op_parser(Operator_Prec+1, RHS);
            if (RHS == NULL) return NULL;
        }

        LHS = new BinaryAST(std::to_string(BinOp), LHS, RHS);
    }
}

static BaseAST* expression_parser() {
    BaseAST* LHS = Base_Parser();
    if (!LHS) return NULL;
    return binary_op_parser(0, LHS);
}

static void init_precedence() {
  Operator_Precedence['-'] = 1;
  Operator_Precedence['+'] = 2;
  Operator_Precedence['/'] = 3;
  Operator_Precedence['*'] = 4;
}

static Module *Module_Ob;
// static IRBuilder<> Builder(getGlobalContext());

static void HandleDefn() {
    if (FunctionDefnAST* F = func_defn_parser()) {
    } else {
        next_token();
    }
}

static void HandleTopExpression() {
    if (FunctionDefnAST* F = top_level_parser()) {
    } else {
        next_token();
    }
}

static void Driver() {
    while (1) {
        switch (Current_token) {
            case EOF_TOKEN: return;
            case ';': next_token(); break;
            case DEF_TOKEN: HandleDefn(); break;
            default: HandleTopExpression(); break;
        }
    }
}

int main(int argc, char* argv[]) {
    // LLVMContext& Context = getGlobalContext();
    const LLVMContextRef& ContextRef = LLVMGetGlobalContext();
    LLVMContext &Context = *unwrap(ContextRef);
    init_precedence();

    file = fopen(argv[1], "r");

    if (file == NULL) {
        printf("Could not open or access file '%s'.", argv[1]);
    }


    next_token();
    Module_Ob = new Module("my compiler", Context);
    Driver();
    Module_Ob->dump();
    return 0;
}