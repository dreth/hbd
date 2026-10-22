import assert from 'node:assert/strict';
import { createRequire } from 'node:module';

// GHSA-vfj7-8cjw-p6xm: reject excessive parsed nesting before recursive
// compilation/expansion, using the same input-validation contract as maxLength.
export function verifyBracesSecurity(packageJsonPath) {
  const require = createRequire(packageJsonPath);
  const braces = require('./index.js');
  assert.deepEqual(braces.expand('src/{app,lib}/file-{1..3}.js'), [
    'src/app/file-1.js', 'src/app/file-2.js', 'src/app/file-3.js',
    'src/lib/file-1.js', 'src/lib/file-2.js', 'src/lib/file-3.js',
  ]);
  assert.equal(braces.compile('a/{b,c}/d'), 'a/(b|c)/d');
  for (const [open, close] of [['{', '}'], ['(', ')']]) {
    const boundary = open.repeat(100) + 'x' + close.repeat(100);
    for (const operation of ['parse', 'compile', 'expand', 'stringify']) {
      assert.doesNotThrow(() => braces[operation](boundary));
      for (const depth of [101, 4000]) {
        const pattern = open.repeat(depth) + 'a,b' + close.repeat(depth);
        assert.throws(() => braces[operation](pattern), {
          name: 'SyntaxError', message: 'Input nesting exceeds maximum depth (100)',
        });
      }
    }
  }
  const mixed = '{('.repeat(51) + 'x' + ')}'.repeat(51);
  assert.throws(() => braces.compile(mixed), { name: 'SyntaxError' });
  // Escaped and quoted brace characters are text, rather than nested ASTs.
  assert.doesNotThrow(() => braces.compile('\\{'.repeat(1000)));
  assert.doesNotThrow(() => braces.compile('"' + '{'.repeat(1000) + '"'));
}
