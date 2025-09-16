from typing import *
from dataclasses import dataclass

@dataclass
class Type:
    name: str
    param: Optional[str] = None
    def type_sig(self) -> str:
        if self.param is None:
            return self.name
        else:
            return f"{self.name} {self.param}"


Arrow = List[Type]

@dataclass
class Inductive:
    itype: Type
    constructor: Dict[str, Arrow]

    def repr(self) -> str:
        constructor_list: List[str] = []
        for name, arrow in self.constructor.items():
            constructor = constructor_template.format(
                name=name,
                arrow=": " + "->".join(map(lambda t: t.type_sig(), arrow)),
            )
            constructor_list.append(constructor)

        return repr_template.format(
            type_sig=self.itype.type_sig(),
            constructor_list="\n".join(constructor_list),
        )

    def generate_go(self) -> str:
        package_name=self.itype.name.lower()
        type_name = self.itype.name
        type_def = self.itype.name
        if self.itype.param is not None:
            type_def += f"[{self.itype.param} any]"
        
        return go_template.format(
            repr=self.repr(),
            package_name=package_name,
            type_def=type_def,
            type_name=type_name,
        )

repr_template = """
inductive {type_sig}
{constructor_list}
"""

constructor_template = "  | {name} {arrow}"

go_template = """
/*
auto generated code from
{repr}
*/


package {package_name}

type {type_def} interface {{
    attr{type_name}()
}}
"""



if __name__ == "__main__":
    s = Inductive(
        itype=Type(name="List", param="T"),
        constructor={
            "Nil": [],
            "Cons": [Type(name="T"), Type(name="List", param="T"), Type(name="List", param="T")],
        },
    ).generate_go()
    print(s)

    s = Inductive(
        itype=Type(name="Nat"),
        constructor={
            "Zero": [],
            "Succ": [Type(name="Nat"), Type(name="Nat")],
        },
    ).generate_go()
    print(s)

