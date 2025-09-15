#!/usr/bin/env python3
"""
Python to Emacs Lisp Transpiler

Converts Python-like syntax to Emacs Lisp syntax with support for:
- let blocks with proper indentation
- variable assignments
- lambda functions with type annotations
- match expressions with pattern matching
- infix operators using curly braces
- comments
"""

import re
import sys
from typing import List, Tuple, Optional

class PythonToElispTranspiler:
    def __init__(self):
        self.indent_level = 0
        self.in_let_block = False
        self.let_vars = []
        
    def transpile(self, python_code: str) -> str:
        """Main transpilation method"""
        lines = python_code.split('\n')
        result_lines = []
        
        i = 0
        while i < len(lines):
            line = lines[i].rstrip()
            
            # Skip empty lines
            if not line.strip():
                result_lines.append('')
                i += 1
                continue
                
            # Handle comments
            if line.strip().startswith('#'):
                result_lines.append(f"; {line.strip()[1:].strip()}")
                i += 1
                continue
            
            # Handle let block start
            if line.strip() == 'let':
                result_lines.append('(let')
                self.in_let_block = True
                self.let_vars = []
                i += 1
                continue
            
            # Handle let block end (empty line or end of file)
            if self.in_let_block and (not line.strip() or i == len(lines) - 1):
                if self.let_vars:
                    # Add all variable assignments
                    for var_line in self.let_vars:
                        result_lines.append(f"    {var_line}")
                    result_lines.append('')
                    result_lines.append('    Unit_0')
                result_lines.append(')')
                self.in_let_block = False
                self.let_vars = []
                i += 1
                continue
            
            # Handle variable assignments in let block
            if self.in_let_block and '=' in line and not line.strip().startswith('#'):
                transpiled_line = self._transpile_assignment(line)
                self.let_vars.append(transpiled_line)
                i += 1
                continue
            
            # Handle lambda functions
            if 'lambda' in line:
                transpiled_line = self._transpile_lambda(line, lines, i)
                if self.in_let_block:
                    self.let_vars.append(transpiled_line)
                else:
                    result_lines.append(transpiled_line)
                i += 1
                continue
            
            # Handle print statements
            if line.strip().startswith('print '):
                expr = line.strip()[6:].strip()
                transpiled_expr = self._transpile_expression(expr)
                transpiled_line = f"(:= _ (inspect {transpiled_expr}))"
                if self.in_let_block:
                    self.let_vars.append(transpiled_line)
                else:
                    result_lines.append(transpiled_line)
                i += 1
                continue
            
            # Handle other expressions
            transpiled_line = self._transpile_expression(line)
            if self.in_let_block:
                self.let_vars.append(transpiled_line)
            else:
                result_lines.append(transpiled_line)
            i += 1
        
        return '\n'.join(result_lines)
    
    def _transpile_assignment(self, line: str) -> str:
        """Convert Python assignment to Lisp assignment"""
        # Remove leading whitespace
        line = line.strip()
        
        # Split on first '=' to handle complex expressions
        parts = line.split('=', 1)
        if len(parts) != 2:
            return line
            
        var_name = parts[0].strip()
        expr = parts[1].strip()
        
        # Transpile the expression
        transpiled_expr = self._transpile_expression(expr)
        
        return f"(:= {var_name} {transpiled_expr})"
    
    def _transpile_lambda(self, line: str, lines: List[str], line_idx: int) -> str:
        """Convert lambda function to Lisp lambda"""
        # Extract lambda parameters and body
        lambda_match = re.match(r'(\w+)\s*=\s*lambda\s*\(([^)]+)\)', line)
        if not lambda_match:
            return line
            
        func_name = lambda_match.group(1)
        params = lambda_match.group(2).strip()
        
        # Parse parameters with type annotations
        param_parts = []
        if params:
            for param in params.split(','):
                param = param.strip()
                if ':' in param:
                    name, type_ = param.split(':', 1)
                    param_parts.append(f"(: {name.strip()} {type_.strip()})")
                else:
                    param_parts.append(param)
        
        # Find the lambda body (next lines until we hit a non-indented line)
        body_lines = []
        i = line_idx + 1
        while i < len(lines) and lines[i].strip() and (lines[i].startswith('   ') or lines[i].startswith('\t')):
            body_lines.append(lines[i].strip())
            i += 1
        
        # Transpile the body
        transpiled_body = self._transpile_lambda_body(body_lines)
        
        # Create the lambda expression
        if len(param_parts) == 1:
            params_str = param_parts[0]
        else:
            params_str = f"({' '.join(param_parts)})"
        
        return f"(:= {func_name} (=> {params_str} {transpiled_body}))"
    
    def _transpile_lambda_body(self, body_lines: List[str]) -> str:
        """Transpile lambda function body"""
        if not body_lines:
            return "Unit_0"
        
        # Check if it's a match expression
        if body_lines[0].strip().startswith('match '):
            return self._transpile_match(body_lines)
        
        # Handle simple expressions
        if len(body_lines) == 1:
            return self._transpile_expression(body_lines[0])
        
        # Handle multiple expressions (wrap in progn)
        transpiled_lines = []
        for line in body_lines:
            transpiled_lines.append(self._transpile_expression(line))
        
        return f"(progn {' '.join(transpiled_lines)})"
    
    def _transpile_match(self, body_lines: List[str]) -> str:
        """Convert match expression to Lisp match"""
        if not body_lines:
            return "Unit_0"
        
        # Extract the match expression
        match_line = body_lines[0].strip()
        match_expr = match_line[6:].strip()  # Remove 'match '
        
        # Extract pattern cases
        cases = []
        i = 1
        while i < len(body_lines):
            line = body_lines[i].strip()
            if line.startswith('| ') and '=>' in line:
                # Parse pattern case
                case_parts = line[2:].split('=>', 1)
                pattern = case_parts[0].strip()
                expr = case_parts[1].strip()
                
                transpiled_expr = self._transpile_expression(expr)
                cases.append(f"(=> {pattern} {transpiled_expr})")
            i += 1
        
        # Create match expression
        cases_str = ' '.join(cases)
        return f"(match {match_expr} {cases_str})"
    
    def _transpile_expression(self, expr: str) -> str:
        """Transpile general expressions"""
        expr = expr.strip()
        
        # Handle infix operators with curly braces
        if any(op in expr for op in ['⊕', '⊗', '+', '-', '*', '/', '==', '!=', '<', '>', '<=', '>=']):
            return self._transpile_infix_operators(expr)
        
        # Handle function calls
        if '(' in expr and ')' in expr:
            return self._transpile_function_call(expr)
        
        # Handle type annotations
        if ':' in expr and '->' in expr:
            return self._transpile_type_annotation(expr)
        
        # Handle simple expressions
        return expr
    
    def _transpile_infix_operators(self, expr: str) -> str:
        """Convert infix operators to Lisp syntax with curly braces"""
        # Define operator precedence (higher number = higher precedence)
        operators = [
            ('⊗', 3), ('*', 3), ('/', 3),
            ('⊕', 2), ('+', 2), ('-', 2),
            ('==', 1), ('!=', 1), ('<', 1), ('>', 1), ('<=', 1), ('>=', 1)
        ]
        
        # Sort by precedence (lowest first)
        operators.sort(key=lambda x: x[1])
        
        for op, _ in operators:
            if op in expr:
                # Split by operator and wrap in curly braces
                parts = expr.split(op)
                if len(parts) > 1:
                    # Transpile each part
                    transpiled_parts = [self._transpile_expression(part.strip()) for part in parts]
                    return f"{{{' '.join(transpiled_parts)}}}"
        
        return expr
    
    def _transpile_function_call(self, expr: str) -> str:
        """Convert function calls to Lisp syntax"""
        # Simple function call: func(args) -> (func args)
        if '(' in expr and ')' in expr:
            paren_start = expr.find('(')
            paren_end = expr.rfind(')')
            
            func_name = expr[:paren_start].strip()
            args_str = expr[paren_start+1:paren_end].strip()
            
            if args_str:
                # Parse arguments
                args = []
                current_arg = ""
                paren_count = 0
                
                for char in args_str:
                    if char == '(':
                        paren_count += 1
                        current_arg += char
                    elif char == ')':
                        paren_count -= 1
                        current_arg += char
                    elif char == ',' and paren_count == 0:
                        args.append(self._transpile_expression(current_arg.strip()))
                        current_arg = ""
                    else:
                        current_arg += char
                
                if current_arg.strip():
                    args.append(self._transpile_expression(current_arg.strip()))
                
                return f"({func_name} {' '.join(args)})"
            else:
                return f"({func_name})"
        
        return expr
    
    def _transpile_type_annotation(self, expr: str) -> str:
        """Convert type annotations to Lisp syntax"""
        # Handle function types: (A -> B) -> (-> A B)
        if '->' in expr:
            # Replace (A -> B) with (-> A B)
            expr = re.sub(r'\(([^)]+)\s*->\s*([^)]+)\)', r'(-> \1 \2)', expr)
        
        return expr

def main():
    """Main function to run the transpiler"""
    if len(sys.argv) != 3:
        print("Usage: python transpiler.py input.py output.el")
        sys.exit(1)
    
    input_file = sys.argv[1]
    output_file = sys.argv[2]
    
    try:
        with open(input_file, 'r') as f:
            python_code = f.read()
        
        transpiler = PythonToElispTranspiler()
        elisp_code = transpiler.transpile(python_code)
        
        with open(output_file, 'w') as f:
            f.write(elisp_code)
        
        print(f"Successfully transpiled {input_file} to {output_file}")
        
    except FileNotFoundError:
        print(f"Error: File {input_file} not found")
        sys.exit(1)
    except Exception as e:
        print(f"Error: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()