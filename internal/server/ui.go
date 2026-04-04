package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Githost</title>
<link href="https://fonts.googleapis.com/css2?family=Libre+Baskerville:ital,wght@0,400;0,700;1,400&family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--mono:'JetBrains Mono',monospace;--serif:'Libre Baskerville',serif}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--serif);line-height:1.6}
.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}.hdr h1{font-family:var(--mono);font-size:.9rem;letter-spacing:2px}.hdr h1 span{color:var(--rust)}
.main{padding:1.5rem;max-width:960px;margin:0 auto}
.stats{display:grid;grid-template-columns:repeat(3,1fr);gap:.5rem;margin-bottom:1rem}
.st{background:var(--bg2);border:1px solid var(--bg3);padding:.6rem;text-align:center;font-family:var(--mono)}
.st-v{font-size:1.3rem;font-weight:700}.st-l{font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.15rem}
.toolbar{display:flex;gap:.5rem;margin-bottom:1rem;align-items:center;flex-wrap:wrap}
.search{flex:1;min-width:180px;padding:.4rem .6rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.search:focus{outline:none;border-color:var(--leather)}
.filter-sel{padding:.4rem .5rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.65rem}
.repo{background:var(--bg2);border:1px solid var(--bg3);padding:.8rem 1rem;margin-bottom:.5rem;transition:border-color .2s}
.repo:hover{border-color:var(--leather)}
.repo-top{display:flex;justify-content:space-between;align-items:flex-start;gap:.5rem}
.repo-name{font-family:var(--mono);font-size:.88rem;font-weight:700;color:var(--rust)}
.repo-desc{font-size:.75rem;color:var(--cd);margin-top:.2rem}
.repo-meta{font-family:var(--mono);font-size:.55rem;color:var(--cm);margin-top:.35rem;display:flex;gap:.6rem;flex-wrap:wrap;align-items:center}
.repo-actions{display:flex;gap:.3rem;flex-shrink:0}
.badge{font-family:var(--mono);font-size:.5rem;padding:.12rem .35rem;text-transform:uppercase;letter-spacing:1px;border:1px solid}
.badge.public{border-color:var(--green);color:var(--green)}.badge.private{border-color:var(--gold);color:var(--gold)}.badge.archived{border-color:var(--cm);color:var(--cm)}
.badge.active{border-color:var(--green);color:var(--green)}.badge.inactive{border-color:var(--cm);color:var(--cm)}
.branch-badge{font-family:var(--mono);font-size:.5rem;padding:.1rem .3rem;background:var(--bg3);color:var(--cd)}
.star{color:var(--gold);font-size:.65rem}
.btn{font-family:var(--mono);font-size:.6rem;padding:.25rem .5rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd);transition:all .2s}
.btn:hover{border-color:var(--leather);color:var(--cream)}.btn-p{background:var(--rust);border-color:var(--rust);color:#fff}
.btn-sm{font-size:.55rem;padding:.2rem .4rem}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.65);z-index:100;align-items:center;justify-content:center}.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:460px;max-width:92vw}
.modal h2{font-family:var(--mono);font-size:.8rem;margin-bottom:1rem;color:var(--rust);letter-spacing:1px}
.fr{margin-bottom:.6rem}.fr label{display:block;font-family:var(--mono);font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.fr input,.fr select,.fr textarea{width:100%;padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.fr input:focus,.fr select:focus,.fr textarea:focus{outline:none;border-color:var(--leather)}
.row2{display:grid;grid-template-columns:1fr 1fr;gap:.5rem}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:1rem}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.85rem}
@media(max-width:600px){.stats{grid-template-columns:repeat(3,1fr)}.row2{grid-template-columns:1fr}.toolbar{flex-direction:column}.search{min-width:100%}}
</style></head><body>
<div class="hdr"><h1><span>&#9670;</span> GITHOST</h1><button class="btn btn-p" onclick="openForm()">+ New Repo</button></div>
<div class="main">
<div class="stats" id="stats"></div>
<div class="toolbar">
<input class="search" id="search" placeholder="Search repositories..." oninput="render()">
<select class="filter-sel" id="vis-filter" onchange="render()"><option value="">All Visibility</option><option value="public">Public</option><option value="private">Private</option></select>
</div>
<div id="repos"></div>
</div>
<div class="modal-bg" id="mbg" onclick="if(event.target===this)closeModal()"><div class="modal" id="mdl"></div></div>
<script>
var A='/api',repos=[],editId=null;

async function load(){var r=await fetch(A+'/repositories').then(function(r){return r.json()});repos=r.repositories||[];renderStats();render();}

function renderStats(){
var total=repos.length;
var pub=repos.filter(function(r){return r.visibility==='public'}).length;
var priv=repos.filter(function(r){return r.visibility==='private'}).length;
document.getElementById('stats').innerHTML=[
{l:'Repositories',v:total},{l:'Public',v:pub,c:'var(--green)'},{l:'Private',v:priv,c:'var(--gold)'}
].map(function(x){return '<div class="st"><div class="st-v" style="'+(x.c?'color:'+x.c:'')+'">'+x.v+'</div><div class="st-l">'+x.l+'</div></div>'}).join('');
}

function fmtSize(b){if(!b)return'0 B';if(b<1024)return b+' B';if(b<1048576)return(b/1024).toFixed(1)+' KB';return(b/1048576).toFixed(1)+' MB';}

function render(){
var q=(document.getElementById('search').value||'').toLowerCase();
var vf=document.getElementById('vis-filter').value;
var f=repos;
if(vf)f=f.filter(function(r){return r.visibility===vf});
if(q)f=f.filter(function(r){return(r.name||'').toLowerCase().includes(q)||(r.description||'').toLowerCase().includes(q)});
if(!f.length){document.getElementById('repos').innerHTML='<div class="empty">No repositories found.</div>';return;}
var h='';f.forEach(function(r){
h+='<div class="repo"><div class="repo-top"><div style="flex:1">';
h+='<div class="repo-name">'+esc(r.name)+'</div>';
if(r.description)h+='<div class="repo-desc">'+esc(r.description)+'</div>';
h+='</div><div class="repo-actions">';
h+='<button class="btn btn-sm" onclick="openEdit(''+r.id+'')">Edit</button>';
h+='<button class="btn btn-sm" onclick="del(''+r.id+'')" style="color:var(--red)">&#10005;</button>';
h+='</div></div>';
h+='<div class="repo-meta">';
h+='<span class="badge '+r.visibility+'">'+r.visibility+'</span>';
if(r.default_branch)h+='<span class="branch-badge">&#9745; '+esc(r.default_branch)+'</span>';
if(r.star_count)h+='<span class="star">&#9733; '+r.star_count+'</span>';
h+='<span>'+fmtSize(r.size_bytes)+'</span>';
h+='<span>'+ft(r.created_at)+'</span>';
h+='</div></div>';
});
document.getElementById('repos').innerHTML=h;
}

async function del(id){if(!confirm('Delete this repository?'))return;await fetch(A+'/repositories/'+id,{method:'DELETE'});load();}

function formHTML(repo){
var i=repo||{name:'',description:'',default_branch:'main',visibility:'private'};
var isEdit=!!repo;
var h='<h2>'+(isEdit?'EDIT REPOSITORY':'NEW REPOSITORY')+'</h2>';
h+='<div class="fr"><label>Name *</label><input id="f-name" value="'+esc(i.name)+'" placeholder="my-project"></div>';
h+='<div class="fr"><label>Description</label><input id="f-desc" value="'+esc(i.description)+'" placeholder="What is this repo about?"></div>';
h+='<div class="row2"><div class="fr"><label>Default Branch</label><input id="f-branch" value="'+esc(i.default_branch||'main')+'"></div>';
h+='<div class="fr"><label>Visibility</label><select id="f-vis">';
['public','private'].forEach(function(v){h+='<option value="'+v+'"'+(i.visibility===v?' selected':'')+'>'+v.charAt(0).toUpperCase()+v.slice(1)+'</option>';});
h+='</select></div></div>';
h+='<div class="acts"><button class="btn" onclick="closeModal()">Cancel</button><button class="btn btn-p" onclick="submit()">'+(isEdit?'Save':'Create')+'</button></div>';
return h;
}

function openForm(){editId=null;document.getElementById('mdl').innerHTML=formHTML();document.getElementById('mbg').classList.add('open');document.getElementById('f-name').focus();}
function openEdit(id){var r=null;for(var j=0;j<repos.length;j++){if(repos[j].id===id){r=repos[j];break;}}if(!r)return;editId=id;document.getElementById('mdl').innerHTML=formHTML(r);document.getElementById('mbg').classList.add('open');}
function closeModal(){document.getElementById('mbg').classList.remove('open');editId=null;}

async function submit(){
var name=document.getElementById('f-name').value.trim();
if(!name){alert('Name is required');return;}
var body={name:name,description:document.getElementById('f-desc').value.trim(),default_branch:document.getElementById('f-branch').value.trim()||'main',visibility:document.getElementById('f-vis').value};
if(editId){await fetch(A+'/repositories/'+editId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}
else{await fetch(A+'/repositories',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}
closeModal();load();
}

function ft(t){if(!t)return'';try{return new Date(t).toLocaleDateString('en-US',{month:'short',day:'numeric',year:'numeric'})}catch(e){return t;}}
function esc(s){if(!s)return'';var d=document.createElement('div');d.textContent=s;return d.innerHTML;}
document.addEventListener('keydown',function(e){if(e.key==='Escape')closeModal();});
load();
</script></body></html>`
