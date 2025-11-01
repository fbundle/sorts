
def test1 := -- Nat and Vec
  id
  $ ind "Nat" [] [
    ("zero", [], var "Nat"),
    ("succ", [("prev", var "Nat")], var "Nat"),
  ]
  $ ind "Vec0" [("n", var "Nat"), ("T", typ 0)] [
    (
      "nil",
      [("T", typ 0)],
      curry (var "Vec0") [var "zero", var "T"],
    ),
    (
      "push",
      [
        ("n", var "Nat"), ("T", typ 0),
        ("v", curry (var "Vec0") [var "n", var "T"]),
        ("x", var "T"),
      ],
      curry (var "Vec0") [curry (var "succ") [var "n"], var "T"],
    ),
  ]
  $ bnd "one" (curry (var "succ") [var "zero"]) (var "Nat")
  $ bnd "empty_vec" (curry (var "nil") [var "Nat"]) (curry (var "Vec0") [var "zero", var "Nat"])
  $ bnd "single_vec" (curry (var "push") [var "empty_vec"]) (curry (var "Vec0") [var "one", var "Nat"])
  $ ann (var "single_vec") (curry (var "Vec0") [var "one", var "Nat"])
