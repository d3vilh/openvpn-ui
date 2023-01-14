#!/usr/bin/env python3
import passlib.hash
import base64
def f(p):
  h = passlib.hash.pbkdf2_sha256.hash(p)
  print('  {"%s", "%s"},' % (p,h))

f('')
f('a')
f('ab')
f('abc')
f('abcd')
f('abcde')
f('abcdef')
f('abcdefg')
f('abcdefgh')
f('abcdefghi')
f('abcdefghij')
f('abcdefghijk')
f('abcdefghijkl')
f('abcdefghijklm')
f('abcdefghijklmn')
f('abcdefghijklmno')
f('abcdefghijklmnop')
f('qrstuvwxyz012345')
f('67890./')
f('ABCDEFGHIJKLMNOP')
f('QRSTUVWXYZ012345')
for i in range(70):
    f(('password'*10)[0:i])
