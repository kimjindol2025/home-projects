const { Lexer } = require('./dist/lexer.js');

const code = `fn test(x: int, y: int): bool {
  return x >= 0 and y <= 10
}`;

console.log("Testing and/or keywords:");
console.log(code);
console.log("\n--- Tokenization ---");

const lexer = new Lexer(code);
const { tokens } = lexer.tokenize();

console.log("Looking for AND token:");
tokens.forEach((t, i) => {
  if (t.lexeme === 'and' || t.lexeme === 'or' || t.type === 'AND' || t.type === 'OR') {
    console.log(`  ${i}: ${t.type} = "${t.lexeme}"`);
  }
});

const andToken = tokens.find(t => t.lexeme === 'and' && t.type === 'AND');
if (andToken) {
  console.log("\n✅ 'and' keyword recognized as AND token!");
} else {
  const andLike = tokens.find(t => t.lexeme === 'and');
  if (andLike) {
    console.log(`\n❌ 'and' tokenized as ${andLike.type}, not AND`);
  } else {
    console.log("\n❌ 'and' not found in tokens");
  }
}
