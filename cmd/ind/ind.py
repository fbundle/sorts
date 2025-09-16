from typing import *
from dataclasses import dataclass

@dataclass
class Type:
    Name: str
    Param: Optional[str] = None

Arrow = List[Type]

@dataclass
class Inductive:
    Type: Type
    Constructor: Dict[str, Arrow]

    def generate_go(self) -> str:
        ...

if __name__ == "__main__":
    ind = Inductive(
        Type=Type(Name="List", Param="T"),
        Constructor={
            "Nil": [],
            "Cons": [Type(Name="T"), Type(Name="List", Param="T"), Type(Name="List", Param="T")],
        },
    )
    print(ind)

