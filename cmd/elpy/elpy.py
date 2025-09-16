from __future__ import annotations

from typing import *
from dataclasses import dataclass

@dataclass
class ELS:
    code: str | list[str]
    def __call__(self, s: ELS) -> ELS:
        return ELS(code=[self.code, s.code])
    
    def __repr__(self) -> str:
        return str(self.code)





Any_2 = ELS(code="Any_2")
inh = ELS(code="inh")

def arrow(*args: ELS) -> ELS:
    if len(args) == 0:
        raise RuntimeError("arrow empty")
    output, args = args[-1], args[:-1]
    while len(args) > 0:
        param, args = args[-1], args[:-1]
        output = ELS(code=["->", param.code, output.code])
    
    return output


# main program

def main():
    Nat = inh(Any_2)
    n0 = inh(Nat)
    succ = inh(arrow(Nat, Nat))




    return Nat, n0, succ


if __name__ == "__main__":
    output = main()
    for els in output:
        print(els)








