import pydantic
import sys

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
    if len(indented_lines) == 0:
        code.append("_CLOSE_BLOCK")
        return code, indented_lines
    
    line = indented_lines[0].line

    if len(indented_lines) == 1:
        code.append("_LINE_BREAK")
        code.extend(line)

        return parse(
            indentation_stack=indentation_stack,
            code=code,
            indented_lines=indented_lines[1:],
        )

    

    curr_ind = indentation_stack[-1]
    next_ind = indented_lines[1].indentation
    assert curr_ind == indented_lines[0].indentation, f"{indentation_stack} {indented_lines[0]} {code}"
    
    
    if curr_ind == next_ind:
        code.append("_LINE_BREAK")
        code.extend(line)

        return parse(
            indentation_stack=indentation_stack,
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

        child: list[Token] = ["_OPEN_BLOCK"]
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

    elif curr_ind > next_ind:
        code.append("_CLOSE_BLOCK")
        return parse(
            indentation_stack=indentation_stack[:-1],
            code=code,
            indented_lines=indented_lines,
        )

    else:
        raise RuntimeError("unreachable")


def main(source: str):
    indented_lines = prepare(source)

    output = parse(
        indentation_stack=[0],
        code=["_OPEN_BLOCK"],
        indented_lines=indented_lines,
    )
    print(output)





if __name__ == "__main__":
    main(open(sys.argv[1]).read())