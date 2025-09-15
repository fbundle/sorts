class IndentedLine(pydantic.BaseModel):
    indentation: int
    line: list[str]

def main(source: str):
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
    
    print(indented_lines)





if __name__ == "__main__":
    main(open(sys.argv[0]).read())