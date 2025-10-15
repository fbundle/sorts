namespace Code

-- Frame γ δ is type of a persistent map with key of type String and value of type δ
class Frame (γ: Type) (δ: Type) where
  set: γ → String → δ → γ
  get?: γ → String → Option δ

end Code
