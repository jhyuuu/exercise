#include "llvm/ADT/APFloat.h"
#include "llvm/ADT/Optional.h"
#include "llvm/ADT/STLExtras.h"
#include "llvm/IR/BasicBlock.h"
#include "llvm/IR/Constants.h"
#include "llvm/IR/DerivedTypes.h"
#include "llvm/IR/Function.h"
#include "llvm/IR/Instructions.h"
#include "llvm/IR/IRBuilder.h"
#include "llvm/IR/LLVMContext.h"
#include "llvm/IR/LegacyPassManager.h"
#include "llvm/IR/Module.h"
#include "llvm/IR/Type.h"
#include "llvm/IR/Verifier.h"
#include "llvm/Support/FileSystem.h"
#include "llvm/Support/Host.h"
#include "llvm/Support/raw_ostream.h"
#include "llvm/Support/TargetRegistry.h"
#include "llvm/Support/TargetSelect.h"
#include "llvm/Target/TargetMachine.h"
#include "llvm/Target/TargetOptions.h"
#include <algorithm>
#include <cassert>
#include <cctype>
#include <cstdio>
#include <cstdlib>
#include <map>
#include <memory>
#include <string>
#include <system_error>
#include <utility>
#include <vector>

//===----------------------------------------------------------------------===//
// Lexer
//===----------------------------------------------------------------------===//

// The lexer returns tokens [0-255] if it is an unknown character, otherwise one
// of these for known things.

enum Token {
    tok_eof = -1,

    // commands
    tok_def = -2,
    tok_extern = -3,
    tok_identifier = -4,
    tok_number = -5,
    tok_if = -6,
    tok_then = -7,
    tok_else = -8,
    tok_for = -9,
    tok_in = -10,
    // operators
    tok_binary = -11, 
    tok_unary = -12,
    tok_var = -13
};

static std::string IdentifierStr; // Filled in if tok_identifier
static double NumVal;             // Filled in if tok_number

/// gettok - Return the next token from standard input.
static int gettok() {
  static int LastChar = ' ';

  // Skip any whitespace.
  while (isspace(LastChar))
    LastChar = getchar();

  if (isalpha(LastChar)) { // identifier: [a-zA-Z][a-zA-Z0-9]*
    IdentifierStr = LastChar;
    while (isalnum((LastChar = getchar())))
      IdentifierStr += LastChar;

    if (IdentifierStr == "def")
      return tok_def;
    if (IdentifierStr == "extern")
      return tok_extern;
    if (IdentifierStr == "if")
      return tok_if;
    if (IdentifierStr == "then")
      return tok_then;
    if (IdentifierStr == "else")
      return tok_else;
    if (IdentifierStr == "for")
      return tok_for;
    if (IdentifierStr == "in")
      return tok_in;
    if (IdentifierStr == "binary")
      return tok_binary;
    if (IdentifierStr == "unary")
      return tok_unary;
    if (IdentifierStr == "var")
      return tok_var;
    return tok_identifier;
  }

  if (isdigit(LastChar) || LastChar == '.') { // Number: [0-9.]+
    std::string NumStr;
    do {
      NumStr += LastChar;
      LastChar = getchar();
    } while (isdigit(LastChar) || LastChar == '.');

    NumVal = strtod(NumStr.c_str(), nullptr);
    return tok_number;
  }

  if (LastChar == '#') {
    // Comment until end of line.
    do
      LastChar = getchar();
    while (LastChar != EOF && LastChar != '\n' && LastChar != '\r');

    if (LastChar != EOF)
      return gettok();
  }

  // Check for end of file.  Don't eat the EOF.
  if (LastChar == EOF)
    return tok_eof;

  // Otherwise, just return the character as its ascii value.
  int ThisChar = LastChar;
  LastChar = getchar();
  return ThisChar;
}

//===----------------------------------------------------------------------===//
// Abstract Syntax Tree (aka Parse Tree)
//===----------------------------------------------------------------------===//


namespace {

/// ExprAST - Base class for all expression nodes.
class ExprAST {
public:
  virtual ~ExprAST() { };
  virtual llvm::Value* codegen() = 0;
};

/// NumberExprAST - Expression class for numeric literals like "1.0".
class NumberExprAST : public ExprAST {
  double Val;

public:
  NumberExprAST(double Val) : Val(Val) {}
  virtual llvm::Value* codegen() override;
};

/// VariableExprAST - Expression class for referencing a variable, like "a".
class VariableExprAST : public ExprAST {
  std::string Name;

public:
  VariableExprAST(const std::string &Name) : Name(Name) {}
  virtual llvm::Value* codegen() override;
  const std::string& getName() const { return Name; }
};

class UnaryExprAST : public ExprAST {
  char Op;
  std::unique_ptr<ExprAST> Operand;

public:
  UnaryExprAST(char Op,
               std::unique_ptr<ExprAST> Operand)
      : Op(Op), Operand(std::move(Operand)) {}

  virtual llvm::Value* codegen() override;
};

class BinaryExprAST : public ExprAST {
  char Op;
  std::unique_ptr<ExprAST> LHS, RHS;

public:
  BinaryExprAST(char Op,
		        std::unique_ptr<ExprAST> LHS,
                std::unique_ptr<ExprAST> RHS)
      : Op(Op), LHS(std::move(LHS)), RHS(std::move(RHS)) {}

  virtual llvm::Value* codegen() override;
};

class CallExprAST : public ExprAST {
  std::string Callee;
  std::vector<std::unique_ptr<ExprAST> > Args;

public:
  CallExprAST(const std::string& callee,
              std::vector<std::unique_ptr<ExprAST> > args)
      : Callee(callee), Args(std::move(args)) { }
  virtual llvm::Value* codegen() override;
};

class IfExprAST : public ExprAST {
    std::unique_ptr<ExprAST> Cond;
    std::unique_ptr<ExprAST> Then;
    std::unique_ptr<ExprAST> Else;

  public:
    IfExprAST(std::unique_ptr<ExprAST> Cond,
              std::unique_ptr<ExprAST> Then,
              std::unique_ptr<ExprAST> Else)
        : Cond(std::move(Cond)),
          Then(std::move(Then)),
          Else(std::move(Else)) { }

    virtual llvm::Value* codegen() override;
};

class ForExprAST : public ExprAST {
    std::string VarName;
    std::unique_ptr<ExprAST> Start;
    std::unique_ptr<ExprAST> End;
    std::unique_ptr<ExprAST> Step;
    std::unique_ptr<ExprAST> Body;

  public:
    ForExprAST(const std::string& VarName,
               std::unique_ptr<ExprAST> Start,
               std::unique_ptr<ExprAST> End,
               std::unique_ptr<ExprAST> Step,
               std::unique_ptr<ExprAST> Body)
        : VarName(VarName),
          Start(std::move(Start)),
          End(std::move(End)),
          Step(std::move(Step)),
          Body(std::move(Body)) { }

    virtual llvm::Value* codegen() override;
};

class VarExprAST : public ExprAST {
    std::vector<std::pair<std::string, std::unique_ptr<ExprAST> > > VarNames;
    std::unique_ptr<ExprAST> Body;

  public:
    VarExprAST(std::vector<std::pair<std::string, std::unique_ptr<ExprAST> > > VarNames,
               std::unique_ptr<ExprAST> Body)
        : VarNames(std::move(VarNames)), Body(std::move(Body)) { }

    virtual llvm::Value* codegen() override;
};

class PrototypeAST {
  std::string Name;
  std::vector<std::string> Args;
  bool IsOperator;
  unsigned Precedence;

public:
  PrototypeAST(const std::string& name,
               std::vector<std::string> Args,
               bool IsOperator = false,
               unsigned Precedence = 0)
      : Name(name), 
        Args(std::move(Args)),
        IsOperator(IsOperator),
        Precedence(Precedence) { }

  const std::string& getName() const { return Name; }

  llvm::Function* codegen();

  bool isUnaryOp()  const { return IsOperator && Args.size() == 1; }
  bool isBinaryOp() const { return IsOperator && Args.size() == 2; }

  char getOperatorName() const { 
      assert(isUnaryOp() || isBinaryOp());
      return Name[Name.size() - 1];
  }

  unsigned getBinaryPrecedence() const { return Precedence; }
};

class FunctionAST {
  std::unique_ptr<PrototypeAST> Proto;
  std::unique_ptr<ExprAST> Body;

public:
  FunctionAST(std::unique_ptr<PrototypeAST> Proto, std::unique_ptr<ExprAST> Body)
      : Proto(std::move(Proto)), Body(std::move(Body)) { }

  llvm::Function* codegen();
};

} // end anonymous namespace

//===----------------------------------------------------------------------===//
// Parser
//===----------------------------------------------------------------------===//

/// CurTok/getNextToken - Provide a simple token buffer.  CurTok is the current
/// token the parser is looking at.  getNextToken reads another token from the
/// lexer and updates CurTok with its results.
static int CurTok;
static int getNextToken() { return CurTok = gettok(); }

/// BinopPrecedence - This holds the precedence for each binary operator that is
/// defined.
static std::map<char, int> BinopPrecedence = {
    {'=', 2},
    {'<', 10},
    {'+', 20},
    {'-', 20},
    {'*', 40},
};

/// GetTokPrecedence - Get the precedence of the pending binary operator token.
static int GetTokPrecedence() {
    if (!isascii(CurTok))
        return -1;

    int TokPrec = BinopPrecedence[CurTok];
    if (TokPrec <= 0)
        return -1;

    return TokPrec;
}

/// LogError* - These are little helper functions for error handling.
std::unique_ptr<ExprAST> LogError(const char* str) {
    fprintf(stderr, "Error: %s\n", str);
    return nullptr;
}

std::unique_ptr<PrototypeAST> LogErrorP(const char* str) {
    LogError(str);
    return nullptr;
}

static std::unique_ptr<ExprAST> ParseExpression();

/// numberexpr ::= number
static std::unique_ptr<ExprAST> ParserNumberExpr() {
    std::unique_ptr<NumberExprAST> result = std::make_unique<NumberExprAST>(NumVal);
    getNextToken();
    return std::move(result);
}

/// parenexpr ::= '(' expression ')'
static std::unique_ptr<ExprAST> ParseParenExpr() {
    getNextToken(); // eat (
    std::unique_ptr<ExprAST> V = ParseExpression();
    if (!V)
        return nullptr;

    if (CurTok != ')')
        return LogError("expected ')'");

    getNextToken(); // eat )
    return std::move(V);
}

/// identifierexpr
///   ::= identifier
///   ::= identifier '(' expression* ')'
static std::unique_ptr<ExprAST> ParseIdentifierExpr() {
    std::string name = IdentifierStr;

    getNextToken();

    if (CurTok != '(') // Simple variable ref.
        return std::make_unique<VariableExprAST>(name);

    // function call
    getNextToken();
    std::vector<std::unique_ptr<ExprAST> > Args;

    if (CurTok != ')') {
        // not empty args
        while (true) {
            std::unique_ptr<ExprAST> Arg = ParseExpression();
            if (Arg)
                Args.push_back(std::move(Arg));
            else
                return nullptr;

            if (CurTok == ')')
                break;

            if (CurTok != ',')
                return LogError("Expected ')' or ',' in argument list");
            
            getNextToken();
        }
    }

    getNextToken();
    
    return std::make_unique<CallExprAST>(name, std::move(Args));
}

/// ifexpr ::= 'if' expression 'then' expression 'else' expression
std::unique_ptr<ExprAST> ParseIfExpr() {
    getNextToken(); // eat if

    std::unique_ptr<ExprAST> Cond = ParseExpression();
    if (!Cond)
        return nullptr;

    if (CurTok != tok_then)
        return LogError("Except then");

    getNextToken(); // eat then

    std::unique_ptr<ExprAST> Then = ParseExpression();
    if (!Then)
        return nullptr;    

    if (CurTok != tok_else)
        return LogError("Except else");

    getNextToken(); // eat else

    std::unique_ptr<ExprAST> Else = ParseExpression();
    if (!Then)
        return nullptr;

    return std::make_unique<IfExprAST>(std::move(Cond),
                                       std::move(Then),
                                       std::move(Else));
}

/// forexpr ::= 'for' identifier '=' expr ',' expr (',' expr)? 'in' expression
static std::unique_ptr<ExprAST> ParseForExpr() {
    getNextToken(); // eat for

    if (CurTok != tok_identifier)
        return LogError("expected identifier after for");
    
    std::string VarName = IdentifierStr;
    getNextToken(); // eat identifier

    if (CurTok != '=')
        return LogError("expected '=' after for");
    getNextToken(); // eat =

    std::unique_ptr<ExprAST> Start = ParseExpression();
    if (!Start)
        return nullptr;

    if (CurTok != ',')
        return LogError("expected ',' after for start value");
    getNextToken(); // eat ,

    std::unique_ptr<ExprAST> End = ParseExpression();
    if (!End)
        return nullptr;

    std::unique_ptr<ExprAST> Step;
    if (CurTok == ',') {
        getNextToken(); // eat ,
        Step = ParseExpression();
        if (!Step)
            return nullptr;
    }

    if (CurTok != tok_in)
        return LogError("expected 'in' after for");
    getNextToken(); // eat 'in'.       

    std::unique_ptr<ExprAST> Body = ParseExpression();
    if (!Body)
        return nullptr;

    return std::make_unique<ForExprAST>(
        VarName,
        std::move(Start),
        std::move(End),
        std::move(Step),
        std::move(Body)
    );
}

/// varexpr ::= 'var' identifier ('=' expression)?
//                    (',' identifier ('=' expression)?)* 'in' expression
static std::unique_ptr<ExprAST> ParseVarExpr() {
    getNextToken(); // eat var

    std::vector<std::pair<std::string, std::unique_ptr<ExprAST> > > VarNames;

    // At least one variable name is required.
    if (CurTok != tok_identifier)
        return LogError("expected identifier after var");

    while (true) {
        std::string Name = IdentifierStr;
        getNextToken();

        std::unique_ptr<ExprAST> Init = nullptr;
        if (CurTok == '=') {
            getNextToken();

            Init = ParseExpression();
            if (!Init)
                return nullptr;
        }

        VarNames.push_back(std::make_pair(Name, std::move(Init)));

        if (CurTok != ',')
            break;
        getNextToken(); // eat ,

        if (CurTok != tok_identifier)
            return LogError("expected identifier list after var");
    }

    // At this point, we have to have 'in'.
    if (CurTok != tok_in)
        return LogError("expected 'in' keyword after 'var'");
    getNextToken(); // eat in

    std::unique_ptr<ExprAST> Body = ParseExpression();
    if (!Body)
        return nullptr;

    return std::make_unique<VarExprAST>(std::move(VarNames), std::move(Body));
}

/// primary
///   ::= identifierexpr
///   ::= numberexpr
///   ::= parenexpr
static std::unique_ptr<ExprAST> ParsePrimary() {
    switch (CurTok) {
        case tok_identifier:
            return ParseIdentifierExpr();
        case tok_number:
            return ParserNumberExpr();
        case tok_if:
            return ParseIfExpr();
        case tok_for:
            return ParseForExpr();
        case tok_var:
            return ParseVarExpr();
        case '(':
            return ParseParenExpr();
        default:
            char buf[128];
            sprintf(buf, "unknown token when expecting an expression: '%c'", (char)CurTok);
            return LogError(buf);
    }
}

static std::unique_ptr<ExprAST> ParseUnary() {
    if (!isascii(CurTok) || CurTok == '(' || CurTok == ',')
        return ParsePrimary();
    
    int Op = CurTok;
    getNextToken();
    std::unique_ptr<ExprAST> Operand = ParseUnary();
    if (!Operand)
        return nullptr;
    return std::make_unique<UnaryExprAST>(Op, std::move(Operand));
}

static std::unique_ptr<ExprAST> ParseBinOpRHS(int ExprPrec,
                                              std::unique_ptr<ExprAST> LHS) {
    while (true) {
        int TokPrec = GetTokPrecedence();

        // 传入左侧优先级和LHS，如果当前运算符优先级低于左侧优先级，
        // 则LHS应该向左归约，即LHS与更左的LHS merge
        // 因此需要“特别的+1”，使得从左向右计算。
        // 测试用例： a + b * c + d
        if (TokPrec < ExprPrec)
            return LHS;

        int BinOp = CurTok;
        getNextToken();

        auto RHS = ParseUnary();
        if (!RHS)
            return nullptr;

        int NextPrec = GetTokPrecedence();
        if (NextPrec > TokPrec)
            RHS = ParseBinOpRHS(TokPrec + 1, std::move(RHS));
            if (!RHS)
                return nullptr;

        LHS = std::make_unique<BinaryExprAST>(BinOp, 
                                              std::move(LHS),
                                              std::move(RHS));
    }
}

/// expression
///   ::= primary binoprhs
///
static std::unique_ptr<ExprAST> ParseExpression() {
    auto LHS = ParseUnary();
    if (!LHS)
        return nullptr;

    return ParseBinOpRHS(0, std::move(LHS));
}

static std::unique_ptr<PrototypeAST> ParsePrototype() {
    std::string FnName;

    unsigned Kind = 0; // 0 = identifier, 1 = unary, 2 = binary.
    unsigned BinaryPrecedence = 30;

    switch (CurTok) {
        case tok_identifier:
            Kind = 0;
            FnName = IdentifierStr;
            getNextToken();
            break;
        case tok_unary:
            getNextToken();
            if (!isascii(CurTok))
                return LogErrorP("Expected unary operator"); 
            FnName = "unary";
            FnName += (char)CurTok;
            Kind = 1;
            getNextToken();
            break;
        case tok_binary: 
            getNextToken();
            if (!isascii(CurTok))
                return LogErrorP("Expected binary operator");
            FnName = "binary";
            FnName += (char)CurTok;
            Kind = 2;
            getNextToken();

            if (CurTok == tok_number) {
                if (NumVal < 1 || NumVal > 100)
                    return LogErrorP("Invalid precedence: must be 1..100");
                BinaryPrecedence = NumVal;
                getNextToken();
            }
            break;
        default:
            return LogErrorP("Expected function name in prototype");
    }

    if (CurTok != '(')
        return LogErrorP("Expected '(' in prototype");

    std::vector<std::string> ArgNames;
    while (getNextToken() == tok_identifier) {
        ArgNames.push_back(IdentifierStr);
    }

    if (CurTok != ')')
        return LogErrorP("Expected ')' in prototype");

    getNextToken(); // eat ')'

    if (Kind && ArgNames.size() != Kind)
        return LogErrorP("Invalid number of operands for operator");

    return std::make_unique<PrototypeAST>(
        FnName,
        ArgNames,
        Kind != 0,
        BinaryPrecedence
    );
}

static std::unique_ptr<FunctionAST> ParseDefinition() {
    getNextToken();

    auto Proto = ParsePrototype();
    if (!Proto)
        return nullptr;

    auto E = ParseExpression();
    if (!E)
        return nullptr;

    return std::make_unique<FunctionAST>(std::move(Proto), std::move(E));
}

static std::unique_ptr<FunctionAST> ParseTopLevelExpr() {
    auto E = ParseExpression();
    if (!E)
        return nullptr;

    auto Proto = std::make_unique<PrototypeAST>("__anonymous_expr",
                                                std::vector<std::string>());
    return std::make_unique<FunctionAST>(std::move(Proto), std::move(E));
}

/// external ::= 'extern' prototype
static std::unique_ptr<PrototypeAST> ParseExtern() {
    getNextToken();
    return ParsePrototype();
}

//===----------------------------------------------------------------------===//
// Code Generation
//===----------------------------------------------------------------------===//

static llvm::LLVMContext TheContext; // 保存类型表和常量值表
static llvm::IRBuilder<> Builder(TheContext); // 用于生成LLVM指令
static std::unique_ptr<llvm::Module> TheModule; // 用于保存IR
static std::map<std::string, llvm::AllocaInst*> NamedValues;
static std::map<std::string, std::unique_ptr<PrototypeAST> > FunctionProtos;

llvm::Value* LogErrorV(const char* str) {
    LogError(str);
    return nullptr;
}

llvm::Function* getFunction(std::string Name) {
    // First, see if the function has already been added to the current module.
    if (auto* F = TheModule->getFunction(Name))
        return F;

    // If not, check whether we can codegen the declaration from some existing
    // prototype.
    auto FI = FunctionProtos.find(Name);
    if (FI != FunctionProtos.end())
        return FI->second->codegen();

    // If no existing prototype exists, return null.
    return nullptr;
}

/// CreateEntryBlockAlloca - Create an alloca instruction in the entry block of
/// the function.  This is used for mutable variables etc.
static llvm::AllocaInst* CreateEntryBlockAlloca(llvm::Function* TheFunction,
                                                const std::string& VarName) {
    llvm::IRBuilder<> TmpBuilder(&TheFunction->getEntryBlock(),
                                 TheFunction->getEntryBlock().begin());
    return TmpBuilder.CreateAlloca(llvm::Type::getDoubleTy(TheContext), nullptr, VarName);
}

llvm::Value* NumberExprAST::codegen() {
    return llvm::ConstantFP::get(TheContext, llvm::APFloat(Val));
}

llvm::Value* VariableExprAST::codegen() {
    llvm::Value* V = NamedValues[Name];
    if (!V) {
        char buf[128];
        sprintf(buf, "Unknown variable name: '%s'", Name.c_str());
        LogErrorV(buf);
    }
    return Builder.CreateLoad(V, Name.c_str());
}

llvm::Value* UnaryExprAST::codegen() {
    llvm::Value* OperandV = Operand->codegen();

    llvm::Function* F = getFunction(std::string("unary") + Op);
    assert (F && "unary operator not found!");

    return Builder.CreateCall(F, OperandV, "unop");
}

llvm::Value* BinaryExprAST::codegen() {
    // Special case '=' because we don't want to emit the LHS as an expression.
    if (Op == '=') {
        VariableExprAST* LHSE = static_cast<VariableExprAST*>(LHS.get());
        if (!LHSE)
            return LogErrorV("destination of '=' must be a variable");

        // Codegen the RHS.
        llvm::Value* Val = RHS->codegen();
        if (!Val)
            return nullptr;

        // Look up the name.
        llvm::Value* Variable = NamedValues[LHSE->getName()];
        if (!Variable)
            return LogErrorV("Unknown variable name");

        Builder.CreateStore(Val, Variable);
        return Val;
    }

    llvm::Value* L = LHS->codegen();
    llvm::Value* R = RHS->codegen();
    if (!L || !R)
        return nullptr;

    switch (Op) {
        case '+': return Builder.CreateFAdd(L, R, "addtmp");
        case '-': return Builder.CreateFSub(L, R, "subtmp");
        case '*': return Builder.CreateFMul(L, R, "multmp");
        case '<': 
            L = Builder.CreateFCmpULT(L, R, "cmptmp");
            return Builder.CreateUIToFP(L, llvm::Type::getDoubleTy(TheContext), "booltmp");
        default:
            break;
    }

    llvm::Function* F = getFunction(std::string("binary") + Op);
    assert (F && "binary operator not found!");

    llvm::Value* Ops[2] = {L, R};
    return Builder.CreateCall(F, Ops, "binop");
}

llvm::Value* CallExprAST::codegen() {
    // Look up the name in the global module table.
    llvm::Function* CalleeF = getFunction(Callee);
    if (!CalleeF) {
        char buf[128];
        sprintf(buf, "Unknown function referenced: '%s'", Callee.c_str());
        return LogErrorV(buf);
    }

    if (CalleeF->arg_size() != Args.size()) {
        char buf[128];
        sprintf(buf, "Incorrect # arguments passed when call function: '%s'", Callee.c_str());
        return LogErrorV(buf);
    }

    std::vector<llvm::Value*> ArgsV;
    for (unsigned i = 0, e = Args.size(); i < e; ++i) {
        ArgsV.push_back(Args[i]->codegen());
        if (!ArgsV.back())
            return nullptr;
    }
    return Builder.CreateCall(CalleeF, ArgsV, "calltmp");
}

llvm::Value* IfExprAST::codegen() {
    llvm::Value* CondV = Cond->codegen();
    if (!CondV) {
        return nullptr;
    }

    // Convert condition to a bool by comparing non-equal to 0.0.
    CondV = Builder.CreateFCmpONE( //Ordered and Not Equal
        CondV,
        llvm::ConstantFP::get(TheContext, llvm::APFloat(0.0)), "ifcond");

    llvm::Function* TheFunction = Builder.GetInsertBlock()->getParent();

    // Create blocks for the then and else cases.  
    // Insert the 'then' block at the end of the function.
    llvm::BasicBlock* ThenBB  = llvm::BasicBlock::Create(TheContext, "then", TheFunction);
    llvm::BasicBlock* ElseBB  = llvm::BasicBlock::Create(TheContext, "else");
    llvm::BasicBlock* MergeBB = llvm::BasicBlock::Create(TheContext, "ifcont");

    Builder.CreateCondBr(CondV, ThenBB, ElseBB);

    // Emit then block.
    Builder.SetInsertPoint(ThenBB);
    llvm::Value* ThenV = Then->codegen();
    if (!ThenV)
        return nullptr;

    Builder.CreateBr(MergeBB);
    ThenBB = Builder.GetInsertBlock();

    // Emit else block.
    TheFunction->getBasicBlockList().push_back(ElseBB);
    Builder.SetInsertPoint(ElseBB);
    llvm::Value* ElseV = Else->codegen();
    if (!ElseV)
        return nullptr;

    Builder.CreateBr(MergeBB);
    ElseBB = Builder.GetInsertBlock();

    // Emit merge block.
    TheFunction->getBasicBlockList().push_back(MergeBB);
    Builder.SetInsertPoint(MergeBB);
    llvm::PHINode* PN = Builder.CreatePHI(llvm::Type::getDoubleTy(TheContext), 2, "iftmp");

    PN->addIncoming(ThenV, ThenBB);
    PN->addIncoming(ElseV, ElseBB);
    return PN;
}

llvm::Value* ForExprAST::codegen() {
    llvm::Function* TheFunction = Builder.GetInsertBlock()->getParent();
    llvm::AllocaInst* Alloca = CreateEntryBlockAlloca(TheFunction, VarName);
    
    llvm::Value* StartV = Start->codegen();
    if (!StartV)
        return nullptr;
    Builder.CreateStore(StartV, Alloca);

    llvm::BasicBlock* LoopBB = llvm::BasicBlock::Create(TheContext, "loop", TheFunction);
    Builder.CreateBr(LoopBB);
    Builder.SetInsertPoint(LoopBB);

    llvm::AllocaInst* OldVal = NamedValues[VarName];
    NamedValues[VarName] = Alloca;

    if (!Body->codegen())
        return nullptr;

    // // Emit the step value
    llvm::Value* StepV = nullptr;
    if (Step) {
        StepV = Step->codegen();
        if (!StepV)
            return nullptr;
    } else {
        StepV = llvm::ConstantFP::get(TheContext, llvm::APFloat(1.0));
    }

    // Compute the end condition.
    llvm::Value* EndV = End->codegen();
    if (!EndV)
        return nullptr;

    llvm::Value* CurVal  = Builder.CreateLoad(Alloca, VarName.c_str());
    llvm::Value* NextVal = Builder.CreateFAdd(CurVal, StepV, "nextvar");
    Builder.CreateStore(NextVal, Alloca);

    EndV = Builder.CreateFCmpONE(EndV,
        llvm::ConstantFP::get(TheContext, llvm::APFloat(0.0)), "loopcond");

    llvm::BasicBlock* AfterBB = llvm::BasicBlock::Create(TheContext, "afterloop", TheFunction);

    Builder.CreateCondBr(EndV, LoopBB, AfterBB);

    // Any new code will be inserted in AfterBB.
    Builder.SetInsertPoint(AfterBB);

    if (OldVal)
        NamedValues[VarName] = OldVal;
    else
        NamedValues.erase(VarName);

    return llvm::ConstantFP::getNullValue(llvm::Type::getDoubleTy(TheContext));
}

llvm::Value* VarExprAST::codegen() {
    std::vector<llvm::AllocaInst*> OldBindings;
    
    llvm::Function* TheFunction = Builder.GetInsertBlock()->getParent();

    for (unsigned i = 0, e = VarNames.size(); i < e; ++i) {
        const std::string& VarName = VarNames[i].first;
        ExprAST* Init = VarNames[i].second.get();

        // Emit the initializer before adding the variable to scope, this prevents
        // the initializer from referencing the variable itself, and permits stuff
        // like this:
        //  var a = 1 in
        //    var a = a in ...   # refers to outer 'a'.
        llvm::Value* InitV = nullptr;
        if (Init) {
            InitV = Init->codegen();
            if (!InitV)
                return nullptr;
        } else {
            InitV = llvm::ConstantFP::get(TheContext, llvm::APFloat(0.0));
        }

        llvm::AllocaInst* Alloca = CreateEntryBlockAlloca(TheFunction, VarName);
        Builder.CreateStore(InitV, Alloca);

        OldBindings.push_back(NamedValues[VarName]);
        NamedValues[VarName] = Alloca;
    }

    llvm::Value* BodyV = Body->codegen();
    if (!BodyV)
        return nullptr;

    for (unsigned i = 0, e = OldBindings.size(); i < e; ++i)
        NamedValues[VarNames[i].first] = OldBindings[i];

    return BodyV;
}

llvm::Function* PrototypeAST::codegen() {
    std::vector<llvm::Type*> Doubles(Args.size(), llvm::Type::getDoubleTy(TheContext));

    llvm::FunctionType* FT =
        llvm::FunctionType::get(llvm::Type::getDoubleTy(TheContext), Doubles, false);

    llvm::Function* F = 
        llvm::Function::Create(FT, llvm::Function::ExternalLinkage, Name, TheModule.get());

    unsigned Idx = 0;
    for (auto& Arg : F->args())
        Arg.setName(Args[Idx++]);

    return F;
}

llvm::Function* FunctionAST::codegen() {
    // Transfer ownership of the prototype to the FunctionProtos map, but keep a
    // reference to it for use below.
    auto& P = *Proto;
    std::string Name = Proto->getName();
    FunctionProtos[Name] = std::move(Proto);
    llvm::Function* TheFunction = getFunction(Name);

    if (!TheFunction)
        return nullptr;

    if (P.isBinaryOp())
        BinopPrecedence[P.getOperatorName()] = P.getBinaryPrecedence();

    if (!TheFunction->empty()) {
        char buf[128];
        sprintf(buf, "Function '%s' cannot be redefined.", Proto->getName().c_str());
        return (llvm::Function*)LogErrorV(buf);
    }

    llvm::BasicBlock* BB = llvm::BasicBlock::Create(TheContext, "entry", TheFunction);
    Builder.SetInsertPoint(BB);

    NamedValues.clear();
    // for (auto& Arg : TheFunction->args())
        // NamedValues[Arg.getName()] = &Arg;
    for (auto& Arg : TheFunction->args()) {
        llvm::AllocaInst* Alloca = CreateEntryBlockAlloca(TheFunction, Arg.getName());
        Builder.CreateStore(&Arg, Alloca);
        NamedValues[Arg.getName()] = Alloca;
    }

    if (llvm::Value* RetVal = Body->codegen()) {
        // Finish off the function.
        Builder.CreateRet(RetVal);

        // Validate the generated code, checking for consistency.
        llvm::verifyFunction(*TheFunction);

        return TheFunction;
    }

    // Error reading body, remove function.
    TheFunction->eraseFromParent();
    if (P.isBinaryOp())
        BinopPrecedence.erase(P.getOperatorName());
    return nullptr;
}

//===----------------------------------------------------------------------===//
// Top-Level parsing and JIT Driver
//===----------------------------------------------------------------------===//

static void InitializeModuleAndPassManager() {
    // Open a new module
    TheModule = std::make_unique<llvm::Module>("my cool jit", TheContext);
}

static void HandleDefinition() {
    if (auto FnAST = ParseDefinition()) {
        if (llvm::Function* FnIR = FnAST->codegen()) {
            fprintf(stderr, "Read function definition: ");
            FnIR->print(llvm::errs());
            fprintf(stderr, "\n");
        }
    } else {
        getNextToken();
    }
}

static void HandleExtern() {
    if (auto ProtoAST = ParseExtern()) {
        if (llvm::Function* FnIR = ProtoAST->codegen()) {
            fprintf(stderr, "Read extern: ");
            FnIR->print(llvm::errs());
            fprintf(stderr, "\n");
            FunctionProtos[ProtoAST->getName()] = std::move(ProtoAST);
        }
    } else {
        getNextToken();
    }
}

static void HandleTopLevelExpression() {
  // Evaluate a top-level expression into an anonymous function.
    if (auto FnAST = ParseTopLevelExpr()) {
        FnAST->codegen();
    } else {
        // Skip token for error recovery.
        getNextToken();
    }
}

static void MainLoop() {
    while (true) {
        fprintf(stderr, "ready> ");
        switch (CurTok) {
            case tok_eof    : return;
            case ';'        : getNextToken(); break;
            case tok_def    : HandleDefinition(); break;
            case tok_extern : HandleExtern(); break;
            default: 
                HandleTopLevelExpression(); 
                break;
        }
    }
}

//===----------------------------------------------------------------------===//
// "Library" functions that can be "extern'd" from user code.
//===----------------------------------------------------------------------===//

#ifdef _WIN32
#define DLLEXPORT __declspec(dllexport)
#else
#define DLLEXPORT
#endif

/// putchard - putchar that takes a double and returns 0.
extern "C" DLLEXPORT double putchard(double X) {
    fputc((char)X, stderr);
    return 0;
}

/// printd - printf that takes a double prints it as "%f\n", returning 0.
extern "C" DLLEXPORT double printd(double X) {
    fprintf(stderr, "%f\n", X);
    return 0;
}

int main(int argc, char* argv[]) {
    fprintf(stderr, "ready> ");
    getNextToken();

    InitializeModuleAndPassManager();

    MainLoop();

    // Initialize the target registry etc.
    llvm::InitializeAllTargetInfos();
    llvm::InitializeAllTargets();
    llvm::InitializeAllTargetMCs();
    llvm::InitializeAllAsmParsers();
    llvm::InitializeAllAsmPrinters();

    auto TargetTriple = llvm::sys::getDefaultTargetTriple();
    TheModule->setTargetTriple(TargetTriple);

    std::string Error;
    auto Target = llvm::TargetRegistry::lookupTarget(TargetTriple, Error);

    // Print an error and exit if we couldn't find the requested target.
    // This generally occurs if we've forgotten to initialise the
    // TargetRegistry or we have a bogus target triple.
    if (!Target) {
        llvm::errs() << Error;
        return 1;
    }

    auto CPU = "generic";
    auto Features = "";

    llvm::TargetOptions opt;
    auto RM = llvm::Optional<llvm::Reloc::Model>();
    auto TargetMachine =
        Target->createTargetMachine(TargetTriple, CPU, Features, opt, RM);

    TheModule->setDataLayout(TargetMachine->createDataLayout());

    auto Filename = "output.o";
    std::error_code EC;
    llvm::raw_fd_ostream dest(Filename, EC, llvm::sys::fs::OF_None);

    if (EC) {
        llvm::errs() << "Could not open file: " << EC.message();
        return 1;
    }

    llvm::legacy::PassManager pass;
    auto FileType = llvm::CGFT_ObjectFile;

    if (TargetMachine->addPassesToEmitFile(pass, dest, nullptr, FileType)) {
        llvm::errs() << "TheTargetMachine can't emit a file of this type";
        return 1;
    }

    pass.run(*TheModule);
    dest.flush();

    llvm::outs() << "Wrote " << Filename << "\n";

    return 0;
}
