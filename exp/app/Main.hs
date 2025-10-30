module Main where

import Prelude hiding (lookup)

type Id = String

data Exp = Type         -- type
  | Var Id              -- variable
  | App Exp Exp         -- application
  | Abs Id Exp          -- abstraction
  | Let Id Exp Exp Exp  -- let binding (?? name value type return)
  | Pi Id Exp Exp       -- pi type

data Val = VType        -- type
  | VGen Int            -- ??
  | VApp Val Val        -- application
  | VClos Env Exp       -- closure

type Env = [(Id, Val)]  -- context/frame/environment

update :: Env -> Id -> Val -> Env
update env x u = (x, u) : env

lookup :: Id -> Env -> Val
lookup x ((y, u) : env) = 
  if x == y then u
  else lookup x env

lookup x [] = error ("lookup" ++ x)

-- a short way of writing the whnf algorithm
app :: Val -> Val -> Val
eval :: Env -> Exp -> Val

app u v =
  case u of
    VClos env (Abs x e) -> eval (update env x v) e
    _ -> VApp u v

eval env e =
  case e of
    Var x         -> lookup x env
    App e1 e2     -> app (eval env e1) (eval env e2)
    Let x e1 _ e3 -> eval (update env x (eval env e1)) e3
    Type          -> VType
    _             -> VClos env e

whnf :: Val -> Val
whnf v =
  case v of
    VApp u w      -> app (whnf u) (whnf v)
    VClos env e   -> eval env e       -- wouldn't return the same thing?
    _             -> v

-- the conversion algorithm; the integer is
-- used to represent the introduction of a fresh variable

eqVal :: (Int, Val, Val) -> Bool
eqVal (k, u1, u2) = 
  case (whnf u1, whnf u2) of
    (VType, VType)                                        -> True
    (VApp t1 w1, VApp t2 w2)                              ->
      eqVal (k, t1, t2) && eqVal (k, w1, w2)
    (VGen k1, VGen k2)                                    -> k1 == k2
    (VClos env1 (Abs x1 e1), VClos env2 (Abs x2 e2))      ->
      let v = VGen k
      in eqVal  ( k + 1,
                  VClos (update env1 x1 v) e1,
                  VClos (update env2 x2 v) e2
                )
    (VClos env1 (Pi x1 a1 b1), VClos env2 (Pi x2 a2 b2))  ->
      let v = VGen k
      in  eqVal (k, VClos env1 a1, VClos env2 a2) 
          &&
          eqVal ( k + 1, 
                  VClos (update env1 x1 v) b1,
                  VClos (update env2 x2 v) b2
                )
    _ -> False

-- type checking and type inference
checkExp :: (Int, Env, Env) -> Exp -> Val -> Bool
inferExp :: (Int, Env, Env) -> Exp -> Val
checkType :: (Int, Env, Env) -> Exp -> Bool

checkType (k, rho, gamma) e = checkExp (k, rho, gamma) e VType
checkExp (k, rho, gamma) e v =
  case e of 
    Abs x n ->
      case whnf v of
        VClos env (Pi y a b) ->
          let v = VGen k
          in checkExp ( k+1,
                        update rho x v,
                        update gamma x (VClos env a)
                      )
                      n (VClos (update env y v) b)
        _ -> error "expected Pi"
    Pi x a b ->
      case whnf v of 
        VType ->  checkType (k, rho, gamma) a
                  &&
                  checkType  ( k,
                              update rho x (VGen k),
                              update gamma x (VClos rho a)
                            ) b
        _ -> error "expected Type"
    Let x e1 e2 e3 ->
      checkType (k, rho, gamma) e2 
      &&
      checkExp  ( k,
                  update rho x (eval rho e1),
                  update gamma x (eval rho e2)
                ) e3 v
    _ -> eqVal (k, inferExp (k, rho, gamma) e, v)

inferExp (k, rho, gamma) e = 
  case e of
    Var id -> lookup id gamma
    App e1 e2 ->
      case whnf (inferExp (k, rho, gamma) e1) of
        VClos env (Pi x a b) ->
          if checkExp (k, rho, gamma) e2 (VClos env a) then
            VClos (update env x (VClos rho e2)) b
          else
            error "application error"
        _ -> error "application, expected Pi"
    Type -> VType
    _ -> error "cannot infer type"

typecheck :: Exp -> Exp -> Bool
typecheck m a =
  checkType (0, [], []) a
  &&
  checkExp (0, [], []) m (VClos [] a)

test :: Bool
test = typecheck  (Abs "A" (Abs "x" (Var "x")))
                  (Pi "A" Type (Pi "x" (Var "A") (Var "A")))


main :: IO ()
main = print test


