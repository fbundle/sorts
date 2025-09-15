class IndentedLine(pydantic.BaseModel):
    indentation: int
    line: list[str]

def prepare(source: str) -> list[IndentedLine]:
    lines: list[str] = source.split("\n")

    indented_lines: list[IndentedLine] = []
    for line in lines:
        indentation: int = 0
        for ch in line:
            if ch == " ":
                indentation += 1
            elif ch == "\t":
                indentation += 2
            else:
                break
        indented_lines.append(Line(
            indentation=indentation,
            line=line.split(),
        ))
    
    return indented_lines




def main(source: str):
    indented_lines = prepare(source)





if __name__ == "__main__":
    main(open(sys.argv[0]).read())