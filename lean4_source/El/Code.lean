structure Beta (α: Type) where
  cmd: α
  args: List α

structure Annot (α: Type) where
  name: String
  type: α

structure Binding (α: Type) where
  name: String
  value: α

structure Typeof (α: Type) where
  value: α

structure Inhabited (α: Type) where
  value: α

structure Pi (α: Type) where
  params: List (Annot α)
  body: α

structure Inductive (α: Type) where
  type: Annot α
  cons: List (Annot α)

inductive Code where
  | name: String → Code
  | beta: Beta Code → Code
  | annot: Annot Code → Code
  | binding: Binding Code → Code
  | typeof: Typeof Code → Code
  | inh: Inhabited Code → Code
  | pi: Pi Code → Code
  | ind: Inductive Code → Code
