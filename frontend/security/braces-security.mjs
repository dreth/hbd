import {readFileSync,writeFileSync,readdirSync,realpathSync,existsSync} from 'node:fs';
import {createHash} from 'node:crypto';
import {join,resolve} from 'node:path';
import {spawnSync} from 'node:child_process';
import {verifyBracesSecurity} from './verify-braces-security.mjs';
const root=resolve(import.meta.dirname,'..');
const hash=bytes=>createHash('sha256').update(bytes).digest('hex');
const before='e572166565f15fa6ad9865ae49d678218e32aabfd1b3720f6d0d43d39800d310';
const after='73544e77a0f07970e39812ae537b622da789c626323568ccf244dea52dcc550c';
const copies=new Set();
for (const entry of readdirSync(join(root,'node_modules'),{recursive:true,withFileTypes:true})) {
 if(entry.name!=='braces'||!(entry.isDirectory()||entry.isSymbolicLink()))continue;
 const p=realpathSync(join(entry.parentPath,entry.name));
 if(existsSync(join(p,'package.json'))&&JSON.parse(readFileSync(join(p,'package.json'))).name==='braces')copies.add(p);
}
if(!copies.size)throw Error('Expected braces dependency was not found');
for(const p of copies){
 const metadata=JSON.parse(readFileSync(join(p,'package.json')));
 if(metadata.version!=='3.0.3')throw Error('Unattested braces version');
 const file=join(p,'lib/parse.js');let s=readFileSync(file,'utf8');
 if(process.argv.includes('--patch')&&hash(s)===before){
  for(const line of ["      block = push({ type: 'paren', nodes: [] });","      block = push(brace);"]){
   if(s.split(line).length!==2)throw Error('Unexpected parser source');
   s=s.replace(line,line+"\n      if (stack.length > 100) {\n        throw new SyntaxError('Input nesting exceeds maximum depth (100)');\n      }");
  }
  if(hash(s)!==after)throw Error('Security patch integrity mismatch');
  writeFileSync(file,s);
 }
 if(hash(readFileSync(file))!==after)throw Error('Unpatched braces parser');
 verifyBracesSecurity(join(p,'package.json'));
}
console.log(`Verified bounded brace/parenthesis parsing in ${copies.size} installed copies.`);
if(process.argv.includes('--audit')){
 const r=spawnSync('npm',['audit','--json'],{cwd:root,encoding:'utf8'});
 if(r.error||![0,1].includes(r.status))throw Error('npm audit could not complete');
 const audit=JSON.parse(r.stdout);
 if(audit.auditReportVersion!==2||!audit.vulnerabilities||!audit.metadata)throw Error('Unexpected audit response');
 const findings=audit.vulnerabilities;
 function verified(name,seen=new Set()){
  if(seen.has(name))return false;
  const entry=findings[name];if(!entry||!Array.isArray(entry.via)||!entry.via.length)return false;
  const next=new Set([...seen,name]);
  return entry.via.every(a=>typeof a==='string'?verified(a,next):name==='braces'&&a.name==='braces'&&a.url==='https://github.com/advisories/GHSA-vfj7-8cjw-p6xm');
 }
 const unresolved=Object.keys(findings).filter(p=>!verified(p));
 if(unresolved.length){console.error(JSON.stringify(unresolved.map(p=>({package:p,...findings[p]}))));process.exit(1);}
 console.log('Audit passed: any remaining advisory metadata refers only to the verified local braces security backport.');
}
