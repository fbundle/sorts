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

    def type_name(self) -> str:
        return self.name
    
    def go_type(self) -> Tuple[str, str]:
        if self.param is None:
            return self.name, self.name
        else:
            return f"{self.name}[{self.param} any]", f"{self.name}[{self.param}]"

    def go_package(self) -> str:
        return self.name.lower()

Arrow = List[Type]

@dataclass
class Inductive:
    itype: Type
    constructor: Dict[str, Arrow]

    def repr(self) -> str:
        constructor_list: List[str] = []
        for name, arrow in self.constructor.items():
            arrow_str = ": " + " -> ".join(map(lambda t: t.type_sig(), arrow))
            constructor = constructor_template.format(
                name=name,
                arrow=arrow_str,
            )
            constructor_list.append(constructor)

        return repr_template.format(
            type_sig=self.itype.type_sig(),
            constructor_list="\n".join(constructor_list),
        )

    def generate_go(self) -> str:
        package_name = self.itype.go_package()
        type_name = self.itype.type_name()
        type_def, type_call = self.itype.go_type()
        
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

go_constructor_template = """

type {type_def} struct {

}

func (o {type_call})
"""



if __name__ == "__main__":
    s = Inductive(
        itype=Type(name="List", param="T"),
        constructor={
            "Nil": [Type(name="List", param="T")],
            "Cons": [Type(name="T"), Type(name="List", param="T"), Type(name="List", param="T")],
        },
    ).generate_go()
    print(s)

    s = Inductive(
        itype=Type(name="Nat"),
        constructor={
            "Zero": [Type(name="Nat")],
            "Succ": [Type(name="Nat"), Type(name="Nat")],
        },
    ).generate_go()
    print(s)

