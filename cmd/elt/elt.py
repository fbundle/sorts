import pydantic
import sys
from typing import Tuple, List

Token = str

class IndentedLine(pydantic.BaseModel):
    indentation: int
    line: list[Token]

def prepare_tokenize(source: str) -> str:
    source = source.replace("(", " ( ")
    source = source.replace(")", " ) ")
    source = source.replace("=>", " => ")
    source = source.replace("=", " = ")
    source = source.replace("  = > ", " => ") # fix clash
    source = source.replace("->", " -> ")
    source = source.replace(":", " : ")
    source = source.replace("|", " | ") # | starts a block, don't add space before that
    return source

def prepare(source: str) -> list[IndentedLine]:
    lines: list[str] = source.split("\n")

    indented_lines: list[IndentedLine] = []
    for line in lines:
        line = line.split("#")[0] # remove comment
        line = line.rstrip()
        if len(line.strip()) == 0:
            continue

        
        indentation: int = 0
        for ch in list(line):
            if ch == " ":
                indentation += 1
            elif ch == "\t":
                indentation += 2
            else:
                break
        line = prepare_tokenize(line).split()
        indented_lines.append(IndentedLine(
            indentation=indentation,
            line=line,
        ))
    
    return indented_lines

def find_start_of_block(line: list[Token]) -> int | None:
    i = len(line)-1
    while i >= 0:
        if line[i] == "let":
            return i
        if line[i] == "match":
            return i 
        if line[i] == "(":
            return i
        if line[i] == "lambda":
            return i

        i -= 1
    return None

def parse(indentation_stack: list[int], code: list[Token], indented_lines: list[IndentedLine]) -> tuple[list[Token], list[IndentedLine]]:
    # print(f"call: {indentation_stack} {code} {indented_lines}")

    if len(indented_lines) == 0:
        code.append("__ELT_CLOSE_BLOCK")
        return code, indented_lines
    
    line = indented_lines[0].line

    if len(indented_lines) == 1:
        code.append("__ELT_LINE_BREAK")
        code.extend(line)

        return parse(
            indentation_stack=indentation_stack,
            code=code,
            indented_lines=indented_lines[1:],
        )

    

    curr_ind = indentation_stack[-1]
    next_ind = indented_lines[1].indentation
    assert curr_ind >= indented_lines[0].indentation
    if curr_ind < indented_lines[0].indentation:
        return parse(
            indentation_stack=indentation_stack[:-1],
            code=code,
            indented_lines=indented_lines,
        )
    
    
    if curr_ind == next_ind:
        code.append("__ELT_LINE_BREAK")
        code.extend(line)

        return parse(
            indentation_stack=indentation_stack,
            code=code,
            indented_lines=indented_lines[1:],
        )
    elif curr_ind > next_ind:
        code.append("__ELT_LINE_BREAK")
        code.extend(line)
        code.append("__ELT_CLOSE_BLOCK")
        return parse(
            indentation_stack=indentation_stack[:-1],
            code=code,
            indented_lines=indented_lines[1:],
        )
    elif curr_ind < next_ind:
        # new block
        i = find_start_of_block(line)
        if i is None:
            raise RuntimeError(f"cannot find start of block: {line}")
        
        for j in range(0, i):
            code.append(line[j])

        child: list[Token] = ["__ELT_OPEN_BLOCK"]
        for j in range(i, len(line)):
            child.append(line[j])

        child, indented_lines = parse(
            indentation_stack=indentation_stack + [next_ind],
            code=child,
            indented_lines=indented_lines[1:]
        )
        code.append(child)

        return parse(
            indentation_stack=indentation_stack,
            code=code,
            indented_lines=indented_lines,
        )
    else:
        raise RuntimeError("unreachable")


# ---------- Rendering helpers ----------

INFIX_OPS = {"⊕", "⊗"}

def transform_arrow_in_parens(tokens: List[str]) -> List[str]:
    out: List[str] = []
    i = 0
    while i < len(tokens):
        if tokens[i] == "(":
            # find next close ) at same level (no deep nesting expected in input)
            j = i + 1
            depth = 1
            while j < len(tokens) and depth > 0:
                if tokens[j] == "(":
                    depth += 1
                elif tokens[j] == ")":
                    depth -= 1
                j += 1
            # j is index after matching ')'
            inner = tokens[i+1:j-1] if j-1 > i+1 else []
            if "->" in inner and depth == 0:
                k = inner.index("->")
                left = [t for t in inner[:k] if t]
                right = [t for t in inner[k+1:] if t]
                out.append("(")
                out.append("->")
                out.extend(left)
                out.extend(right)
                out.append(")")
            else:
                out.extend(tokens[i:j])
            i = j
        else:
            out.append(tokens[i])
            i += 1
    return out

def join_tokens_with_parens(tokens: List[str]) -> str:
    out: List[str] = []
    for tok in tokens:
        if tok == "(":
            out.append("(")
        elif tok == ")":
            out.append(")")
        else:
            if len(out) == 0 or out[-1] == "(":
                out.append(tok)
            else:
                out.append(" " + tok)
    return "".join(out)

def render_param(tokens: List[str]) -> str:
    # tokens like: ( : x Nat ) or ( x : Nat ) depending on tokenizer
    # We expect: ( x : Type )
    assert tokens[0] == "("
    assert tokens[-1] == ")"
    inner = [t for t in tokens[1:-1] if t != ":"]
    assert len(inner) == 2, f"bad lambda param: {tokens}"
    name, typ = inner
    return f"(: {name} {typ})"

def is_infix(tokens: List[str]) -> bool:
    return any(t in INFIX_OPS for t in tokens)

def render_simple_expr(tokens: List[str]) -> str:
    tokens = transform_arrow_in_parens(tokens)
    if len(tokens) == 0:
        return "Unit_0"
    if is_infix(tokens):
        return "{" + " ".join(tokens) + "}"
    if tokens[0] == "(":
        # Detect arrow type inside parens: ( A -> B ) -> (-> A B)
        inner = tokens[1:-1]
        if "->" in inner:
            i = inner.index("->")
            left = " ".join(inner[:i]).strip()
            right = " ".join(inner[i+1:]).strip()
            parts_left = [t for t in inner[:i] if t]
            parts_right = [t for t in inner[i+1:] if t]
            if len(parts_left) >= 1 and len(parts_right) >= 1:
                return "(-> " + " ".join(parts_left) + " " + " ".join(parts_right) + ")"
        return join_tokens_with_parens(tokens)
    if len(tokens) == 1:
        return tokens[0]
    s = f"({" ".join(tokens)})"
    s = s.replace("( ->", "(->")
    return s

def split_at(tokens: List[str], marker: str) -> Tuple[List[str], List[str]]:
    if marker in tokens:
        i = tokens.index(marker)
        return tokens[:i], tokens[i+1:]
    return tokens, []


# ---------- Statement/Expression emitters over prepared lines ----------

def emit_statements(lines: List[IndentedLine], start: int, base_ind: int) -> Tuple[List[str], int]:
    out: List[str] = []
    i = start
    while i < len(lines):
        ind = lines[i].indentation
        toks = lines[i].line
        if ind < base_ind:
            break
        if ind > base_ind:
            # should not happen: nested handled by owners
            break

        # let block
        if len(toks) == 1 and toks[0] == "let":
            # open let
            out.append("(let")
            # body
            body, j = emit_statements(lines, i+1, lines[i+1].indentation)
            out.extend(("    " + stmt) if not stmt.startswith("(let") and not stmt == ")" else stmt for stmt in body)
            out.append("    Unit_0")
            out.append(")")
            i = j
            continue

        # print statement
        if len(toks) >= 2 and toks[0] == "print":
            expr = render_simple_expr(toks[1:])
            out.append(f"(:= _ (inspect {expr}))")
            i += 1
            continue

        # assignment or lambda
        if "=" in toks:
            left, right = split_at(toks, "=")
            assert len(left) == 1, f"unexpected assignment lhs: {left}"
            name = left[0]

            # lambda form
            if len(right) >= 1 and right[0] == "lambda":
                # expect: lambda ( x : T )
                # grab param paren group
                assert "(" in right and ")" in right, f"lambda missing param parens: {right}"
                pstart = right.index("(")
                # find matching ) from pstart
                pend = len(right) - 1 - right[::-1].index(")")
                param = render_param(right[pstart:pend+1])

                # body is nested block at next line with greater indent
                body_expr, j = emit_expr_block(lines, i+1, lines[i+1].indentation)
                out.append(f"(:= {name} (=> {param} {body_expr}))")
                i = j
                continue

            # normal assignment
            expr = render_simple_expr(right)
            out.append(f"(:= {name} {expr})")
            i += 1
            continue

        # if we reach here, skip unknown line silently
        i += 1

    return out, i

def emit_expr_block(lines: List[IndentedLine], start: int, base_ind: int) -> Tuple[str, int]:
    # Handles a single expression defined by a nested block (used for lambda bodies)
    i = start
    # Expect either a match or a single expression/assignment-like line
    if i >= len(lines) or lines[i].indentation < base_ind:
        return "Unit_0", i

    ind = lines[i].indentation
    toks = lines[i].line

    # match expression header: match x with
    if len(toks) >= 3 and toks[0] == "match" and toks[-1] == "with":
        subject = render_simple_expr(toks[1:-1])
        i += 1
        # collect cases
        cases: List[str] = []
        while i < len(lines) and lines[i].indentation >= base_ind:
            ctoks = lines[i].line
            if ctoks[0] != "|":
                break
            # split at =>
            pat_toks, rhs_toks = split_at(ctoks[1:], "=>")
            pat = f"({" ".join(pat_toks)})" if len(pat_toks) > 1 else pat_toks[0]
            rhs = render_simple_expr(rhs_toks)
            cases.append(f"(=> {pat} {rhs})")
            i += 1
        expr = "(match " + subject + "\n        " + "\n        ".join(cases) + "\n    )"
        return expr, i

    # otherwise, render a simple expression from this line
    expr = render_simple_expr(toks)
    return expr, i + 1


def transpile(source: str) -> str:
    lines = prepare(source)
    stmts, _ = emit_statements(lines, 0, 0)
    return "\n".join(stmts) + "\n"


def main():
    if len(sys.argv) < 2:
        print("usage: python transpiler.py <input.py>")
        sys.exit(1)
    source = open(sys.argv[1]).read()
    output = transpile(source)
    print(output)


if __name__ == "__main__":
    main()


