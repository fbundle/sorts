# el

The goal of this project is to implement minimal dependent type

## the first time EL2 can type check dependent type

![screenshot1.png](https://raw.githubusercontent.com/fbundle/sorts/refs/heads/master/screenshots/screenshot1.png)

in the code, it was able to verify that the type of `append_if_empty` matches the dependent type 
```lean
(m: Nat) (T: U_1) (vec: Vec m T) (val: T): (
  match m with
    | zero => Vec (succ m) T
    | succ _ => Vec m T
)
```
