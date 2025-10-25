import Std
namespace EL2.Term

class NameMap M α where
  size: M → Nat
  set: M → String → α → M
  get?: M → String → Option α

instance : NameMap (Std.HashMap String α) α where
  size := Std.HashMap.size
  set := Std.HashMap.insert
  get? := Std.HashMap.get?


end EL2.Term
