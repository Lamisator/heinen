package main

var indexHTML = `<!DOCTYPE html>
<html lang="de">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Heinen – Das Zahnquiz</title>
<style>
@font-face{font-family:'Outfit';font-style:normal;font-weight:300;font-display:swap;src:url('/fonts/Outfit-Light.woff2') format('woff2')}
@font-face{font-family:'Outfit';font-style:normal;font-weight:400;font-display:swap;src:url('/fonts/Outfit-Regular.woff2') format('woff2')}
@font-face{font-family:'Outfit';font-style:normal;font-weight:600;font-display:swap;src:url('/fonts/Outfit-SemiBold.woff2') format('woff2')}
@font-face{font-family:'Outfit';font-style:normal;font-weight:800;font-display:swap;src:url('/fonts/Outfit-ExtraBold.woff2') format('woff2')}
@font-face{font-family:'Outfit';font-style:normal;font-weight:900;font-display:swap;src:url('/fonts/Outfit-Black.woff2') format('woff2')}
@font-face{font-family:'Space Mono';font-style:normal;font-weight:400;font-display:swap;src:url('/fonts/SpaceMono-Regular.woff2') format('woff2')}
@font-face{font-family:'Space Mono';font-style:normal;font-weight:700;font-display:swap;src:url('/fonts/SpaceMono-Bold.woff2') format('woff2')}
:root{--bg:#0a0a0f;--bg2:#12121a;--surface:#1e1e2e;--border:#2a2a3e;--text:#e4e4ef;--text2:#8888a0;--accent:#ff3366;--accent2:#ff6b9d;--correct:#00e68a;--wrong:#ff3366;--gold:#ffd700;--tooth-white:#f0ece0;--tooth-dead:#1a1a1a}
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:'Outfit',sans-serif;background:var(--bg);color:var(--text);min-height:100vh;overflow-x:hidden}
body::before{content:'';position:fixed;inset:0;pointer-events:none;z-index:0;background:radial-gradient(ellipse at 20% 50%,rgba(255,51,102,.06) 0%,transparent 50%),radial-gradient(ellipse at 80% 20%,rgba(255,107,157,.04) 0%,transparent 50%),radial-gradient(ellipse at 50% 80%,rgba(255,215,0,.03) 0%,transparent 50%)}
.app{position:relative;z-index:1;max-width:1100px;margin:0 auto;padding:20px;min-height:100vh}
.header{text-align:center;padding:20px 0 12px}.logo{font-family:'Space Mono',monospace;font-size:2.6rem;font-weight:700;letter-spacing:-2px;cursor:pointer;background:linear-gradient(135deg,var(--accent),var(--gold));-webkit-background-clip:text;-webkit-text-fill-color:transparent}.logo-sub{font-size:.75rem;color:var(--text2);letter-spacing:4px;text-transform:uppercase}
.screen{display:none}.screen.active{display:block}
.ig{margin-bottom:14px;text-align:left}.ig label{display:block;font-size:.7rem;color:var(--text2);text-transform:uppercase;letter-spacing:2px;margin-bottom:5px}.ig input,.ig select{width:100%;padding:12px 14px;background:var(--surface);border:1px solid var(--border);border-radius:10px;color:var(--text);font-family:'Outfit',sans-serif;font-size:.95rem;outline:none;transition:border-color .2s}.ig input:focus,.ig select:focus{border-color:var(--accent)}.ig input[type=checkbox]{width:auto;margin-right:8px;vertical-align:middle}
.btn{display:inline-block;padding:12px 28px;border:none;border-radius:10px;font-family:'Outfit',sans-serif;font-size:.9rem;font-weight:600;cursor:pointer;transition:all .2s;text-transform:uppercase;letter-spacing:1px}.btn-p{background:linear-gradient(135deg,var(--accent),#cc2952);color:#fff;width:100%;margin-top:6px}.btn-p:hover{transform:translateY(-2px);box-shadow:0 6px 20px rgba(255,51,102,.3)}.btn-s{background:var(--surface);color:var(--text);border:1px solid var(--border);width:100%;margin-top:6px}.btn-s:hover{border-color:var(--accent)}.btn-sm{padding:8px 16px;font-size:.75rem;width:auto;margin-top:0}.btn-d{background:rgba(255,51,102,.12);color:var(--wrong);border:1px solid rgba(255,51,102,.3);width:auto;margin-top:6px}.btn-d:hover{background:rgba(255,51,102,.25)}
.toast{position:fixed;bottom:28px;left:50%;transform:translateX(-50%) translateY(100px);background:var(--surface);border:1px solid var(--border);padding:12px 22px;border-radius:10px;font-size:.85rem;z-index:1000;transition:transform .3s;pointer-events:none}.toast.show{transform:translateX(-50%) translateY(0)}.toast.error{border-color:var(--wrong);color:var(--wrong)}.toast.success{border-color:var(--correct);color:var(--correct)}
.nav-bar{display:flex;justify-content:flex-end;gap:10px;padding:8px 0;position:absolute;top:16px;right:20px;z-index:50}.nav-btn{background:var(--surface);border:1px solid var(--border);color:var(--text2);padding:6px 14px;border-radius:8px;font-size:.75rem;cursor:pointer;font-family:'Outfit',sans-serif;transition:all .2s}.nav-btn:hover{border-color:var(--accent);color:var(--text)}
.login-container{max-width:380px;margin:60px auto;text-align:center}.login-title{font-size:1.2rem;font-weight:600;margin-bottom:20px}
.dc{max-width:500px;margin:30px auto;text-align:center}.da{display:flex;flex-direction:column;gap:10px;max-width:320px;margin:0 auto}
.admin-container{max-width:700px;margin:20px auto}.as{background:var(--surface);border:1px solid var(--border);border-radius:14px;padding:20px;margin-bottom:20px}.ast{font-size:.8rem;color:var(--text2);text-transform:uppercase;letter-spacing:2px;margin-bottom:14px}
.user-row{display:flex;align-items:center;justify-content:space-between;padding:10px 12px;background:var(--bg2);border-radius:8px;margin-bottom:8px;flex-wrap:wrap;gap:8px}.user-info{display:flex;align-items:center;gap:10px}
.ub{font-size:.6rem;padding:2px 8px;border-radius:4px;text-transform:uppercase;letter-spacing:1px}.ub.admin{background:rgba(255,215,0,.15);color:var(--gold)}.ub.user{background:rgba(136,136,160,.15);color:var(--text2)}.ub.hash{background:rgba(100,150,200,.15);color:#6496c8;margin-left:4px}.ub.hash.bcrypt{background:rgba(0,200,100,.15);color:#00c864}.ub.hash.sha256{background:rgba(255,180,100,.15);color:#ffb464}.ub.hash.legacy{background:rgba(255,100,100,.15);color:#ff6464}
.user-actions{display:flex;gap:6px;flex-wrap:wrap}.ib{background:none;border:1px solid var(--border);color:var(--text2);padding:4px 10px;border-radius:6px;cursor:pointer;font-size:.75rem;font-family:'Outfit',sans-serif;transition:all .2s}.ib:hover{border-color:var(--accent);color:var(--text)}.ib.danger:hover{border-color:var(--wrong);color:var(--wrong)}.ib:disabled{opacity:.3;cursor:default;pointer-events:none}
.sound-row{display:flex;gap:10px;align-items:center;margin-bottom:10px;flex-wrap:wrap}.sound-label{font-size:.8rem;color:var(--text2);min-width:140px}.sound-status{font-size:.75rem;color:var(--correct)}.sound-status.none{color:var(--text2)}
.add-user-form{display:flex;gap:10px;align-items:flex-end;flex-wrap:wrap}.add-user-form .ig{flex:1;min-width:120px;margin-bottom:0}
.cc{max-width:520px;margin:20px auto}.sp{background:var(--surface);border:1px solid var(--border);border-radius:14px;padding:20px;margin-bottom:20px}.spt{font-size:.8rem;color:var(--text2);text-transform:uppercase;letter-spacing:2px;margin-bottom:14px}.sg{display:grid;grid-template-columns:1fr 1fr;gap:12px}.sg .ig{margin-bottom:0}.sg input,.sg select{padding:9px 11px;font-size:.9rem}
.lobby-container{max-width:800px;margin:16px auto}.gcd{text-align:center;background:var(--surface);border:1px solid var(--border);border-radius:14px;padding:20px;margin-bottom:18px}.gcl{font-size:.7rem;color:var(--text2);text-transform:uppercase;letter-spacing:3px;margin-bottom:4px}.il{font-family:'Space Mono',monospace;font-size:.7rem;color:var(--gold);word-break:break-all;cursor:pointer;padding:8px 12px;background:var(--bg2);border-radius:6px;margin-top:8px;display:inline-block}.il:hover{text-decoration:underline}
.pst{font-size:.8rem;color:var(--text2);text-transform:uppercase;letter-spacing:2px;margin-bottom:12px}.pg{display:grid;grid-template-columns:repeat(auto-fill,minmax(170px,1fr));gap:12px;margin-bottom:20px}
.pc{background:var(--bg2);border:1px solid var(--border);border-radius:12px;padding:14px;text-align:center;position:relative;transition:all .3s}.pc.eliminated{opacity:.45;border-color:var(--wrong)}.pc.disconnected{opacity:.35}.pc.is-me{border-color:var(--accent);border-width:2px}.pc.q-answered{border-color:var(--correct);background:rgba(0,230,138,.05)}.pn{font-weight:600;font-size:.9rem;margin-bottom:3px}.ps{font-size:.6rem;color:var(--text2);text-transform:uppercase;letter-spacing:1px;margin-bottom:8px}.ps.host{color:var(--gold)}.ps.elim{color:var(--wrong)}.ps.delegate{color:var(--accent2)}.ans-progress{font-size:.75rem;color:var(--text2);margin-bottom:10px;text-align:center}
.teeth-wrap{position:relative;padding:8px 0;min-height:60px}.teeth-rows{display:flex;flex-direction:column;align-items:center;gap:2px}.tooth-row{display:flex;justify-content:center;gap:2px}.tooth{width:14px;height:20px;transition:all .4s cubic-bezier(.68,-.55,.265,1.55)}.tooth.upper{border-radius:3px 3px 7px 7px}.tooth.lower{border-radius:7px 7px 3px 3px}.tooth.alive{background:linear-gradient(180deg,var(--tooth-white) 0%,#d8d4c8 100%);box-shadow:0 2px 6px rgba(0,0,0,.3),inset 0 1px 0 rgba(255,255,255,.4)}.tooth.dead{background:var(--tooth-dead);box-shadow:inset 0 2px 4px rgba(0,0,0,.5)}.tooth.upper.just-lost{animation:tf-down .7s ease-in forwards}.tooth.lower.just-lost{animation:tf-up .7s ease-in forwards}
@keyframes tf-down{0%{transform:translateY(0) rotate(0);opacity:1}40%{transform:translateY(12px) rotate(-10deg);opacity:.9}100%{transform:translateY(20px) rotate(15deg);opacity:.5;background:var(--tooth-dead)}}
@keyframes tf-up{0%{transform:translateY(0) rotate(0);opacity:1}40%{transform:translateY(-12px) rotate(10deg);opacity:.9}100%{transform:translateY(-20px) rotate(-15deg);opacity:.5;background:var(--tooth-dead)}}
.heinen-overlay{position:absolute;inset:0;display:flex;align-items:center;justify-content:center;pointer-events:none;z-index:10}.ht{font-family:'Space Mono',monospace;font-weight:900;white-space:nowrap;padding:4px 10px;border-radius:6px;background:rgba(10,10,15,.85);letter-spacing:1px}.ht.lost{font-size:.75rem;color:#ff4477;text-shadow:0 0 16px rgba(255,51,102,.8),0 0 30px rgba(255,51,102,.4);animation:hp-lost 1.2s ease-in-out infinite}.ht.dead{font-size:.8rem;color:#ff2255;text-shadow:0 0 20px rgba(255,0,60,.9),0 0 40px rgba(255,0,60,.5);animation:hp-dead 1.8s ease-in-out infinite}
@keyframes hp-lost{0%,100%{opacity:.85;transform:scale(1)}50%{opacity:1;transform:scale(1.08)}}@keyframes hp-dead{0%,100%{opacity:.8;transform:scale(1)}50%{opacity:1;transform:scale(1.12)}}
.kick-btn{position:absolute;top:5px;right:5px;background:none;border:none;color:var(--text2);cursor:pointer;font-size:.7rem;padding:2px 5px;border-radius:4px;opacity:0;transition:opacity .2s}.pc:hover .kick-btn{opacity:1}.kick-btn:hover{color:var(--wrong);background:rgba(255,51,102,.1)}
.intro-overlay{position:fixed;inset:0;z-index:100;background:var(--bg);display:flex;flex-direction:column;align-items:center;justify-content:center}.intro-title{font-family:'Space Mono',monospace;font-size:5rem;font-weight:700;letter-spacing:-3px;background:linear-gradient(135deg,var(--accent),var(--gold),var(--accent2));background-size:200% 200%;-webkit-background-clip:text;-webkit-text-fill-color:transparent;animation:ig 2s ease infinite,iss .8s cubic-bezier(.68,-.55,.265,1.55) forwards}@keyframes ig{0%,100%{background-position:0% 50%}50%{background-position:100% 50%}}@keyframes iss{0%{transform:scale(.3);opacity:0}100%{transform:scale(1);opacity:1}}.intro-slogan{font-size:1.1rem;color:var(--text2);margin-top:20px;text-align:center;max-width:520px;line-height:1.6;animation:ifl .8s ease .5s both}@keyframes ifl{0%{opacity:0;transform:translateY(15px)}100%{opacity:1;transform:translateY(0)}}
.qc{max-width:680px;margin:16px auto;text-align:center}.qh{display:flex;justify-content:space-between;align-items:center;margin-bottom:18px;padding:0 8px}.qn{font-size:.75rem;color:var(--text2);text-transform:uppercase;letter-spacing:2px}.timer{font-family:'Space Mono',monospace;font-size:1.5rem;font-weight:700;color:var(--text);transition:color .3s}.timer.urgent{color:var(--wrong);animation:tp .5s ease-in-out infinite}@keyframes tp{0%,100%{transform:scale(1)}50%{transform:scale(1.1)}}.qt{font-size:1.35rem;font-weight:600;margin-bottom:24px;line-height:1.4;padding:0 8px}
.og{display:grid;grid-template-columns:1fr 1fr;gap:10px;margin-bottom:24px}.ob{padding:15px 16px;background:var(--surface);border:2px solid var(--border);border-radius:12px;color:var(--text);font-family:'Outfit',sans-serif;font-size:.9rem;font-weight:500;cursor:pointer;transition:all .2s;text-align:left}.ob:hover:not(.sel):not(.dis){border-color:var(--accent);background:rgba(255,51,102,.05);transform:translateY(-1px)}.ob.sel{border-color:var(--gold)!important;background:rgba(255,215,0,.12)!important;box-shadow:0 0 14px rgba(255,215,0,.35)}.ob.sel .ol{color:var(--gold)}.ob.correct{border-color:var(--correct)!important;background:rgba(0,230,138,.12)!important;box-shadow:0 0 14px rgba(0,230,138,.35)!important}.ob.wrong{border-color:var(--wrong)!important;background:rgba(255,51,102,.18)!important;box-shadow:0 0 22px rgba(255,51,102,.7),0 0 40px rgba(255,51,102,.3)!important;animation:wrong-pulse 1.2s ease-in-out infinite}@keyframes wrong-pulse{0%,100%{box-shadow:0 0 22px rgba(255,51,102,.7),0 0 40px rgba(255,51,102,.3)}50%{box-shadow:0 0 30px rgba(255,51,102,.9),0 0 55px rgba(255,51,102,.5)}}.ob.dis{cursor:default;opacity:.85}.ob.spectator{cursor:default;opacity:.4;pointer-events:none}.ol{font-family:'Space Mono',monospace;font-weight:700;margin-right:8px;color:var(--text2)}.ob.correct .ol{color:var(--correct)}.ob.wrong .ol{color:var(--wrong)}
.ri{text-align:center;margin:16px 0;font-size:.9rem;color:var(--text2)}.end-container{text-align:center;max-width:600px;margin:30px auto}.winner-display{font-size:1.8rem;font-weight:800;margin:14px 0;background:linear-gradient(135deg,var(--gold),#ffaa00);-webkit-background-clip:text;-webkit-text-fill-color:transparent}.no-winner{font-size:1.3rem;color:var(--text2);margin:14px 0}.end-sub{font-size:.9rem;color:var(--text2);margin-bottom:20px}
.loading-container{text-align:center;padding:50px 20px}
.tooth-stage{width:220px;height:230px;margin:0 auto 22px;perspective:780px;perspective-origin:50% 42%;--r:96px}
.tooth-iso{position:relative;width:100%;height:100%;transform-style:preserve-3d;transform:rotateX(-16deg)}
.tooth-rot{position:absolute;inset:0;transform-style:preserve-3d;animation:tooth-spin 4s linear infinite;will-change:transform}
@keyframes tooth-spin{from{transform:rotateY(0deg)}to{transform:rotateY(-360deg)}}
.tooth-face{position:absolute;inset:0;display:flex;align-items:center;justify-content:center}
.tooth-face svg{width:130px;height:170px;display:block;overflow:visible;filter:drop-shadow(0 8px 14px rgba(0,0,0,.55)) drop-shadow(0 0 3px rgba(255,215,0,.25))}
.orbit-rot{position:absolute;inset:0;transform-style:preserve-3d;animation:orbit-spin 3s linear infinite;will-change:transform;pointer-events:none}
@keyframes orbit-spin{from{transform:rotateY(0deg)}to{transform:rotateY(360deg)}}
.orbit-dot{position:absolute;left:50%;top:50%;width:8px;height:8px;margin-left:-4px;margin-top:-4px;background:radial-gradient(circle at 38% 32%,#fff 0%,#fff 40%,#dcdcec 78%,#9a9ab8 100%);border-radius:50%;box-shadow:0 0 9px rgba(255,255,255,.6);opacity:var(--o,1);animation:dot-orient 3s linear infinite;will-change:transform}
@keyframes dot-orient{from{transform:rotateY(var(--a,0deg)) translateX(var(--r,96px)) rotateY(calc(-1 * var(--a,0deg))) rotateY(0deg)}to{transform:rotateY(var(--a,0deg)) translateX(var(--r,96px)) rotateY(calc(-1 * var(--a,0deg))) rotateY(-360deg)}}
.orbit-dot.head{width:11px;height:11px;margin-left:-5.5px;margin-top:-5.5px;background:radial-gradient(circle at 38% 32%,#fff 0%,#fff 55%,#e6e6f6 80%,#b8b8d0 100%);box-shadow:0 0 16px #fff,0 0 30px rgba(255,255,255,.55)}
.loading-text{color:var(--text2);font-size:.85rem;letter-spacing:.5px}
.tooth-spinner-mini{display:block;width:150px;height:150px;margin:0 auto}
.tooth-spinner-mini .tooth-stage{width:150px;height:150px;margin:0;perspective:560px;--r:64px}
.tooth-spinner-mini .tooth-face svg{width:88px;height:116px}
.tooth-spinner-mini .orbit-dot{width:6px;height:6px;margin-left:-3px;margin-top:-3px}
.tooth-spinner-mini .orbit-dot.head{width:9px;height:9px;margin-left:-4.5px;margin-top:-4.5px}
.sv{background:var(--surface);border:1px solid var(--border);border-radius:14px;padding:20px;margin-bottom:18px}.sv-grid{display:grid;grid-template-columns:1fr 1fr;gap:8px 20px}.sv-item{display:flex;justify-content:space-between;padding:6px 0;border-bottom:1px solid var(--border)}.sv-label{font-size:.75rem;color:var(--text2);text-transform:uppercase;letter-spacing:1px}.sv-value{font-size:.85rem;color:var(--text);font-weight:600}
.mute-btn{z-index:90;background:var(--surface);border:1px solid var(--border);border-radius:50%;width:44px;height:44px;display:none;align-items:center;justify-content:center;cursor:pointer;font-size:1.2rem;transition:all .2s;color:var(--text)}.mute-btn:hover{border-color:var(--accent)}.mute-btn.muted{color:var(--text2)}
.host-controls{text-align:center;margin-top:16px}.tutorial-content{max-width:700px;margin:20px auto;background:var(--surface);border:1px solid var(--border);border-radius:14px;padding:30px;line-height:1.7}.tutorial-content h1{font-size:1.5rem;margin-bottom:16px;color:var(--gold)}.tutorial-content h2{font-size:1.1rem;margin:20px 0 10px;color:var(--accent)}.tutorial-content h3{font-size:.95rem;margin:16px 0 8px;color:var(--text)}.tutorial-content p{margin-bottom:12px;color:var(--text2)}.tutorial-content ul,.tutorial-content li{color:var(--text2);margin-left:20px;margin-bottom:6px}.tutorial-content strong{color:var(--text)}
.ai-row{display:flex;gap:10px;align-items:flex-end;margin-bottom:10px;flex-wrap:wrap}.ai-row .ig{flex:1;min-width:150px;margin-bottom:0}
.lobby-card{background:var(--surface);border:1px solid var(--border);border-radius:12px;padding:16px;margin-bottom:12px;display:flex;justify-content:space-between;align-items:center;flex-wrap:wrap;gap:10px}.lobby-card:hover{border-color:var(--accent)}.lc-info{flex:1}.lc-name{font-weight:600;font-size:1rem}.lc-details{font-size:.75rem;color:var(--text2);margin-top:4px}.lc-badge{font-size:.6rem;padding:2px 8px;border-radius:4px;text-transform:uppercase;letter-spacing:1px}.lc-badge.open{background:rgba(0,230,138,.15);color:var(--correct)}.lc-badge.pw{background:rgba(255,215,0,.15);color:var(--gold)}
.log-table{width:100%;border-collapse:collapse;font-family:'Space Mono',monospace;font-size:.7rem}
.log-table thead{position:sticky;top:0;background:var(--surface);z-index:1}
.log-table th{text-align:left;padding:8px 10px;color:var(--text2);text-transform:uppercase;font-size:.6rem;letter-spacing:1px;border-bottom:1px solid var(--border);font-weight:600;font-family:'Outfit',sans-serif}
.log-table td{padding:6px 10px;border-bottom:1px solid var(--border);color:var(--text2);vertical-align:top;word-break:break-word}
.log-table tr:hover td{background:rgba(255,255,255,.02)}
.log-lvl{display:inline-block;font-size:.55rem;font-weight:700;padding:2px 6px;border-radius:3px;text-transform:uppercase;letter-spacing:.5px;font-family:'Space Mono',monospace}
.log-lvl.DEBUG{background:rgba(136,136,160,.2);color:var(--text2)}
.log-lvl.INFO{background:rgba(0,150,255,.15);color:#5aafff}
.log-lvl.WARN{background:rgba(255,215,0,.15);color:var(--gold)}
.log-lvl.ERROR{background:rgba(255,51,102,.18);color:var(--wrong)}
.log-action{color:var(--text);font-weight:600}
.log-empty{padding:30px;text-align:center;color:var(--text2);font-size:.85rem}
.qe-card{background:var(--bg2);border:1px solid var(--border);border-radius:12px;padding:16px;margin-bottom:10px;transition:border-color .2s}.qe-card:hover{border-color:var(--accent)}.qe-question{font-size:.95rem;font-weight:600;color:var(--text);margin:10px 0;line-height:1.5}.qe-opts{display:grid;grid-template-columns:1fr 1fr;gap:6px;margin-bottom:8px}.qe-opt{display:flex;align-items:flex-start;gap:7px;padding:7px 10px;background:var(--surface);border:1px solid var(--border);border-radius:8px;font-size:.78rem;color:var(--text2);line-height:1.4}.qe-opt.cor{border-color:var(--correct);background:rgba(0,230,138,.08);color:var(--text)}.qe-ol{font-family:'Space Mono',monospace;font-weight:700;font-size:.68rem;color:var(--text2);min-width:14px;padding-top:1px}.qe-opt.cor .qe-ol{color:var(--correct)}.qe-tag{display:inline-block;font-size:.58rem;padding:2px 7px;border-radius:4px;text-transform:uppercase;letter-spacing:.5px;font-weight:700}.qe-tag.leicht{background:rgba(0,150,255,.15);color:#5aafff}.qe-tag.mittel{background:rgba(255,215,0,.15);color:var(--gold)}.qe-tag.schwer{background:rgba(255,140,0,.15);color:#ff8c00}.qe-tag.extrem{background:rgba(255,51,102,.18);color:var(--wrong)}.qe-tag.topic{background:rgba(136,136,160,.15);color:var(--text2)}.qe-tag.ai-openai{background:rgba(16,163,127,.15);color:#10a37f}.qe-tag.ai-anthropic{background:rgba(210,90,40,.15);color:#d25a28}.qe-tag.ai-manual{background:rgba(136,136,160,.15);color:var(--text2)}.qe-modal-overlay{display:none;position:fixed;inset:0;background:rgba(0,0,0,.65);z-index:500;align-items:center;justify-content:center;overflow-y:auto}.qe-modal-box{background:var(--bg2);border:1px solid var(--border);border-radius:14px;padding:24px;max-width:560px;width:92%;box-shadow:0 8px 40px rgba(0,0,0,.5);margin:20px auto}.qe-modal-title{font-family:'Space Mono',monospace;font-size:.85rem;font-weight:700;letter-spacing:2px;background:linear-gradient(135deg,var(--accent),var(--gold));-webkit-background-clip:text;-webkit-text-fill-color:transparent;margin-bottom:18px}.qe-date{font-size:.62rem;color:var(--text2);font-family:'Space Mono',monospace;margin-top:6px}.qe-filters{display:flex;gap:8px;flex-wrap:wrap;margin-bottom:12px}.qe-filters input,.qe-filters select{padding:8px 12px;background:var(--bg);border:1px solid var(--border);border-radius:8px;color:var(--text);font-family:'Outfit',sans-serif;font-size:.85rem;flex:1;min-width:110px;outline:none;transition:border-color .2s}.qe-filters input:focus,.qe-filters select:focus{border-color:var(--accent)}.qe-count{font-size:.72rem;color:var(--text2);margin-bottom:10px}.qe-pagination{display:flex;gap:6px;justify-content:center;margin-top:12px;flex-wrap:wrap}.qe-stats{display:flex;gap:8px;flex-wrap:wrap;margin-bottom:14px;padding:12px 14px;background:var(--bg);border:1px solid var(--border);border-radius:10px}.qe-stat{display:flex;flex-direction:column;align-items:center;gap:2px;padding:0 10px;border-right:1px solid var(--border)}.qe-stat:last-child{border-right:none}.qe-stat-val{font-family:'Space Mono',monospace;font-size:1.1rem;font-weight:700;color:var(--text)}.qe-stat-lbl{font-size:.58rem;color:var(--text2);text-transform:uppercase;letter-spacing:1px}
#screen-questions{display:none;position:fixed;inset:0;background:var(--bg);z-index:200;overflow-y:auto;overflow-x:hidden}.qe-page-header{position:sticky;top:0;z-index:10;background:rgba(10,10,15,.97);backdrop-filter:blur(8px);border-bottom:1px solid var(--border);padding:13px 24px;display:flex;align-items:center;gap:14px;flex-wrap:wrap}.qe-page-title{font-family:'Space Mono',monospace;font-size:.85rem;font-weight:700;letter-spacing:2px;background:linear-gradient(135deg,var(--accent),var(--gold));-webkit-background-clip:text;-webkit-text-fill-color:transparent;white-space:nowrap}.qe-body{padding:18px 24px;max-width:1600px;margin:0 auto}.qe-grid{display:grid;grid-template-columns:1fr;gap:10px}.qe-edit-row{display:flex;gap:7px;flex-wrap:wrap;margin-bottom:10px;align-items:center}.qe-edit-field{padding:6px 10px;background:var(--bg2);border:1px solid var(--border);border-radius:7px;color:var(--text);font-family:'Outfit',sans-serif;font-size:.82rem;outline:none;transition:border-color .2s}.qe-edit-field:focus{border-color:var(--accent)}.qe-edit-ta{width:100%;padding:8px 12px;background:var(--bg2);border:1px solid var(--border);border-radius:8px;color:var(--text);font-family:'Outfit',sans-serif;font-size:.9rem;resize:vertical;min-height:54px;outline:none;margin-bottom:10px;line-height:1.5;box-sizing:border-box}.qe-edit-ta:focus{border-color:var(--accent)}.qe-opt-edit{display:flex;align-items:center;gap:8px;padding:7px 10px;background:var(--bg2);border:1px solid var(--border);border-radius:8px;transition:border-color .2s}.qe-opt-edit.sel-cor{border-color:var(--correct);background:rgba(0,230,138,.06)}.qe-opt-edit input[type=radio]{accent-color:var(--correct);width:13px;height:13px;cursor:pointer;flex-shrink:0}.qe-opt-edit input[type=text]{flex:1;background:none;border:none;color:var(--text);font-family:'Outfit',sans-serif;font-size:.83rem;outline:none;min-width:0}.qe-tab-bar{display:flex;gap:0;border-bottom:1px solid var(--border);margin-bottom:0}.qe-tab-btn{padding:10px 20px;background:none;border:none;border-bottom:2px solid transparent;color:var(--text2);font-family:'Outfit',sans-serif;font-size:.85rem;font-weight:600;cursor:pointer;transition:all .2s;margin-bottom:-1px}.qe-tab-btn:hover{color:var(--text)}.qe-tab-btn.active{color:var(--accent);border-bottom-color:var(--accent)}.qe-topics-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(180px,1fr));gap:10px;margin-top:4px}.qe-topic-card{background:var(--bg2);border:1px solid var(--border);border-radius:10px;padding:14px 16px;cursor:pointer;transition:all .2s}.qe-topic-card:hover{border-color:var(--accent);transform:translateY(-2px)}.qe-topic-card-name{font-size:.9rem;font-weight:600;margin-bottom:6px;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}.qe-topic-card-count{font-family:'Space Mono',monospace;font-size:1.15rem;font-weight:700;color:var(--text);margin-bottom:8px}.qe-topic-card-diffs{display:flex;gap:4px;flex-wrap:wrap}.qe-tag.opts{background:rgba(136,136,160,.15);color:var(--text2)}.q-report-btn{background:rgba(255,51,102,.08);border:1px solid rgba(255,51,102,.25);color:rgba(255,100,130,.7);padding:5px 14px;border-radius:8px;font-size:.72rem;font-family:'Outfit',sans-serif;cursor:pointer;transition:all .2s;letter-spacing:.5px}.q-report-btn:hover:not(:disabled){background:rgba(255,51,102,.16);border-color:rgba(255,51,102,.5);color:var(--wrong)}.q-report-btn:disabled{opacity:.45;cursor:default}
.qe-dupe-pair{background:var(--bg2);border:1px solid var(--border);border-radius:12px;padding:16px;margin-bottom:14px;transition:border-color .2s}.qe-dupe-pair:hover{border-color:var(--accent)}.qe-dupe-score{display:inline-block;font-family:'Space Mono',monospace;font-size:.7rem;font-weight:700;padding:2px 8px;border-radius:4px;background:rgba(255,215,0,.12);color:var(--gold);margin-bottom:10px}.qe-dupe-cols{display:grid;grid-template-columns:1fr 1fr;gap:12px;margin-bottom:10px}.qe-dupe-col{background:var(--surface);border:1px solid var(--border);border-radius:9px;padding:12px}.qe-dupe-col-label{font-size:.62rem;color:var(--text2);text-transform:uppercase;letter-spacing:2px;margin-bottom:6px}.qe-dupe-actions{display:flex;gap:6px;flex-wrap:wrap;justify-content:flex-end}@media(max-width:700px){.qe-dupe-cols{grid-template-columns:1fr}}@media(min-width:1100px){.qe-grid{grid-template-columns:1fr 1fr}}@media(max-width:700px){.qe-body{padding:12px}}
@media(max-width:600px){.logo{font-size:2rem}.intro-title{font-size:3rem}.og{grid-template-columns:1fr}.sg{grid-template-columns:1fr}.qt{font-size:1.1rem}.pg{grid-template-columns:repeat(auto-fill,minmax(140px,1fr))}.add-user-form{flex-direction:column}.ai-row{flex-direction:column}.nav-bar{position:static;justify-content:center;margin-bottom:10px}.sv-grid{grid-template-columns:1fr}}
</style>
</head>
<body>
<div class="app">
<div class="header"><div class="logo" onclick="goHome()">HEINEN</div><div class="logo-sub">Das Zahnquiz</div></div>
<div class="nav-bar" id="nav-bar"></div>

<div id="screen-home" class="screen active">
<div style="max-width:600px;margin:20px auto">
<div style="display:flex;gap:10px;justify-content:center;margin-bottom:20px;flex-wrap:wrap">
<button class="btn btn-p" style="width:auto;padding:12px 24px" id="home-login-btn" onclick="showScreen('login')">Anmelden</button>
<button class="btn btn-s" style="width:auto;padding:12px 24px;display:none" id="home-panel-btn" onclick="showScreen('dashboard')">User Panel</button>
</div>
<div class="pst">Offene Lobbys</div>
<div id="lobby-list"><div style="color:var(--text2);text-align:center;padding:20px">Keine offenen Lobbys vorhanden.</div></div>
<div style="margin-top:10px;text-align:center"><div class="ig"><label>Beitritt per Einladungscode</label><div style="display:flex;gap:8px"><input type="text" id="home-invite" placeholder="Code eingeben..." style="flex:1"/><button class="btn btn-sm btn-p" onclick="joinByCode()">Beitreten</button></div></div></div>
<div class="pst" style="margin-top:30px">Spielanleitung</div>
<div class="tutorial-content" id="home-tutorial"></div>
</div></div>

<div id="screen-login" class="screen"><div class="login-container"><div class="login-title">Anmelden</div>
<div class="ig"><label>Benutzername</label><input type="text" id="login-user" autocomplete="username"/></div>
<div class="ig"><label>Passwort</label><input type="password" id="login-pass" autocomplete="current-password"/></div>
<button class="btn btn-p" onclick="doLogin()">Anmelden</button>
<button class="btn btn-s" id="passkey-login-btn" onclick="doPasskeyLogin()" style="margin-top:6px">🔑 Mit Passkey anmelden</button>
<button class="btn btn-s" onclick="goHome()" style="margin-top:10px">Zurück</button></div></div>

<div id="screen-dashboard" class="screen"><div class="dc"><div style="font-size:1rem;color:var(--text2);margin-bottom:24px" id="welcome-text"></div>
<div class="da"><button class="btn btn-p" onclick="showScreen('create')">Neues Spiel erstellen</button>
<button class="btn btn-s" id="btn-admin" onclick="loadAdmin();showScreen('admin')" style="display:none">Admin Control Panel</button>
<button class="btn btn-s" onclick="showScreen('pw')">Passwort ändern</button>
<button class="btn btn-s" onclick="loadPasskeys();showScreen('passkeys')">🔑 Passkeys verwalten</button>
<button class="btn btn-s" onclick="goHome()">Zurück zur Startseite</button></div></div></div>

<div id="screen-passkeys" class="screen"><div class="cc">
<div class="sp"><div class="spt">Passkeys</div>
<div style="font-size:.82rem;color:var(--text2);margin-bottom:14px">Passkeys erlauben passwortlosen Login per Fingerabdruck, Gesichtserkennung oder Sicherheitsschlüssel.</div>
<div id="passkeys-list" style="margin-bottom:14px"></div>
<div class="ig"><label>Name des neuen Passkeys</label><input type="text" id="passkey-name" placeholder="z.B. MacBook, iPhone…" maxlength="50"/></div>
<button class="btn btn-p" onclick="doRegisterPasskey()">+ Passkey registrieren</button>
</div>
<button class="btn btn-s" onclick="showScreen('dashboard')" style="margin-top:10px">Zurück</button>
</div></div>

<div id="screen-pw" class="screen"><div class="login-container"><div class="login-title">Passwort ändern</div>
<div class="ig"><label>Altes Passwort</label><input type="password" id="pw-old"/></div>
<div class="ig"><label>Neues Passwort</label><input type="password" id="pw-new"/></div>
<div class="ig"><label>Bestätigen</label><input type="password" id="pw-new2"/></div>
<button class="btn btn-p" onclick="changePw()">Ändern</button>
<button class="btn btn-s" onclick="showScreen('dashboard')" style="margin-top:10px">Zurück</button></div></div>

<div id="screen-admin" class="screen"><div class="admin-container">
<div class="as"><div class="ast">Vergangene Spiele</div>
<div id="games-list-view"><div id="games-list-container" style="background:var(--bg);border:1px solid var(--border);border-radius:8px;max-height:420px;overflow:auto"><div style="padding:20px;color:var(--text2);text-align:center">Lade...</div></div></div>
<div id="games-detail-view" style="display:none"><button class="btn btn-sm btn-s" onclick="backToGamesList()" style="margin-bottom:12px;width:auto">&larr; Zurück zur Übersicht</button><div id="games-detail-content"></div></div></div>

<div class="as"><div class="ast">Fragen-Datenbank</div>
<div style="display:flex;align-items:center;gap:14px;flex-wrap:wrap">
<button class="btn btn-s" onclick="showQuestions()" style="max-width:280px">Fragen-Explorer &#8599;</button>
<span id="qe-entry-stats" style="font-size:.75rem;color:var(--text2)"></span>
</div></div>

<div class="as"><div class="ast">KI-Konfiguration</div>
<div class="ai-row"><div class="ig"><label>Anbieter</label><select id="ai-provider" onchange="onProviderChange()"><option value="openai">OpenAI</option><option value="anthropic">Anthropic (Claude)</option></select></div><div class="ig"><label>Modell</label><select id="ai-model"></select></div></div>
<div class="ai-row" id="ai-openai-row"><div class="ig"><label>OpenAI API-Key</label><input type="password" id="ai-openai-key" placeholder="sk-..."/></div></div>
<div class="ai-row" id="ai-anthropic-row" style="display:none"><div class="ig"><label>Anthropic API-Key</label><input type="password" id="ai-anthropic-key" placeholder="sk-ant-..."/></div></div>
<div class="ai-row"><div class="ig"><label>Intro-Dauer (Sek.)</label><input type="number" id="ai-intro-delay" value="4" min="1" max="30" style="max-width:100px"/></div></div>
<div style="display:flex;gap:10px;margin-top:8px"><button class="btn btn-sm btn-s" onclick="testAI()">API testen</button></div>
<div class="api-status" id="ai-status" style="font-size:.8rem;margin-top:8px"></div></div>

<div class="as"><div class="ast">Sounds (global)</div><div id="sound-sections"></div>
<div class="ast" style="margin-top:18px">Lautstärke</div>
<div id="vol-sliders"></div></div>

<div class="as"><div class="ast">Logs</div>
<div style="display:flex;gap:10px;margin-bottom:10px;flex-wrap:wrap;align-items:center"><button class="btn btn-sm btn-s" onclick="loadLogs()" style="width:auto">Aktualisieren</button><a href="/api/logs/export" class="btn btn-sm btn-s" style="text-decoration:none;text-align:center;width:auto">Exportieren</a><button class="btn btn-sm btn-d" onclick="clearLogs()" style="width:auto">Logs leeren</button></div>
<div style="display:flex;gap:10px;margin-bottom:10px;flex-wrap:wrap"><input type="text" id="log-search" placeholder="Suche..." style="flex:1;min-width:140px;padding:8px 12px;background:var(--bg);border:1px solid var(--border);border-radius:8px;color:var(--text);font-family:'Outfit',sans-serif;font-size:.85rem" oninput="loadLogs()"/><select id="log-level" style="padding:8px 12px;background:var(--bg);border:1px solid var(--border);border-radius:8px;color:var(--text);font-family:'Outfit',sans-serif;font-size:.85rem" onchange="loadLogs()"><option value="">Alle Levels</option><option value="DEBUG">DEBUG</option><option value="INFO">INFO</option><option value="WARN">WARN</option><option value="ERROR">ERROR</option></select></div>
<div id="log-table-container" style="background:var(--bg);border:1px solid var(--border);border-radius:8px;max-height:500px;overflow:auto"><div style="padding:20px;color:var(--text2);text-align:center">Lade Logs...</div></div>
<div id="log-meta" style="font-size:.7rem;color:var(--text2);margin-top:6px"></div></div>

<div class="as"><div class="ast">Benutzer*innen</div><div id="users-list"></div>
<div style="margin-top:14px"><div class="ast">Neue*r Benutzer*in</div>
<div class="add-user-form"><div class="ig"><label>Benutzername</label><input type="text" id="new-user-name"/></div>
<div class="ig"><label>Passwort</label><input type="password" id="new-user-pass"/></div>
<div class="ig"><label>Rolle</label><select id="new-user-role"><option value="0">Benutzer*in</option><option value="1">Admin</option></select></div>
<button class="btn btn-sm btn-p" style="margin-top:0;white-space:nowrap" onclick="addUser()">Hinzufügen</button></div></div></div>
<button class="btn btn-p" onclick="saveAIConfig()">Einstellungen speichern</button>
<button class="btn btn-s" onclick="showScreen('dashboard')">Zurück</button></div></div>

<div id="screen-create" class="screen"><div class="cc">
<div class="sp"><div class="spt">Spiel konfigurieren</div>
<div class="ig"><label>Dein Anzeigename</label><input type="text" id="host-name" placeholder="Name..." maxlength="20"/></div>
<div class="sg">
<div class="ig"><label>Thema</label><input type="text" id="create-topic" value="Allgemeinwissen"/></div>
<div class="ig"><label>Schwierigkeit</label><select id="create-diff"><option value="leicht">Leicht</option><option value="mittel" selected>Mittel</option><option value="schwer">Schwer</option><option value="extrem">Extrem</option></select></div>
<div class="ig"><label>Spielmodus</label><select id="create-mode" onchange="toggleQC()"><option value="classic">Klassisch</option><option value="elimination">Elimination</option><option value="kfo_battle_royale">KFO Battle Royale</option><option value="kfo_singleplayer">KFO Singleplayer</option></select></div>
<div class="ig" id="qc-group"><label>Anzahl Fragen</label><input type="number" id="create-questions" value="10" min="1" max="50"/></div>
<div class="ig" id="sd-group" style="display:none"><label>Start-Schwierigkeit</label><select id="create-startdiff"><option value="leicht" selected>Leicht</option><option value="mittel">Mittel</option><option value="schwer">Schwer</option><option value="extrem">Extrem</option></select></div>
<div class="ig"><label>Zeit pro Frage (Sek.)</label><input type="number" id="create-time" value="20" min="5" max="120"/></div>
<div class="ig"><label>Antwortmöglichkeiten</label><select id="create-options"><option value="2">2</option><option value="3">3</option><option value="4" selected>4</option></select></div>
<div class="ig"><label>Anzahl Zähne</label><input type="number" id="create-teeth" value="5" min="1" max="20"/></div>
</div>
<div class="ig" style="margin-top:8px"><label><input type="checkbox" id="create-tutorial" checked/>Tutorial vor Spielstart anzeigen</label></div>
<div class="ig" style="margin-top:4px"><label><input type="checkbox" id="create-playintro" checked/>Intro-Sound abspielen</label></div>
<div class="ig" style="margin-top:4px"><label><input type="checkbox" id="create-websearch"/>Internet-Recherche für Fragen (nur OpenAI, langsamer)</label></div>
<div class="sg" style="margin-top:8px">
<div class="ig"><label>Lobby-Name</label><input type="text" id="create-lobbyname" placeholder="Zufällig..." maxlength="30"/></div>
<div class="ig"><label>Lobby-Modus</label><select id="create-lobbymode" onchange="toggleLobbyPw()"><option value="invite">Nur Einladung</option><option value="password">Mit Passwort</option><option value="open">Offen</option></select></div>
</div>
<div class="ig" id="lobby-pw-group" style="display:none"><label>Lobby-Passwort</label><input type="password" id="create-lobbypw"/></div>
</div>
<button class="btn btn-p" onclick="createGame()">Spiel erstellen</button>
<button class="btn btn-s" onclick="showScreen('dashboard')" style="margin-top:10px">Zurück</button></div></div>

<div id="screen-join" class="screen"><div class="login-container"><div class="login-title">Spiel beitreten</div>
<div class="ig"><label>Dein Name</label><input type="text" id="join-name" placeholder="Name..." maxlength="20"/></div>
<div class="ig" id="join-pw-group" style="display:none"><label>Lobby-Passwort</label><input type="password" id="join-pw"/></div>
<button class="btn btn-p" onclick="doJoin()">Beitreten</button>
<button class="btn btn-s" onclick="goHome()" style="margin-top:10px">Zurück</button></div></div>

<div id="screen-lobby" class="screen"><div class="lobby-container">
<div class="gcd"><div class="gcl">Einladungslink</div><div class="il" id="invite-link" onclick="copyInvite()"></div>
<div style="font-size:.6rem;color:var(--text2);margin-top:6px">Klicken zum Kopieren</div>
<div id="qr-code" style="margin:16px auto 0;display:inline-block"></div></div>
<div id="lobby-settings"></div>
<div class="pst" id="players-count">Spieler*innen (0)</div>
<div class="pg" id="lobby-players"></div>
<div id="start-btn-container" style="display:none;text-align:center"><button class="btn btn-p" onclick="startGame()" style="max-width:280px">Spiel starten</button></div></div></div>

<div id="screen-loading" class="screen"><div class="loading-container"><div class="tooth-mount"></div><div class="loading-text">Fragen werden generiert...</div></div></div>
<div id="screen-refill" class="screen"><div class="loading-container"><div class="tooth-mount"></div><div class="loading-text">Neue Fragen werden nachgeladen...</div></div></div>
<div id="screen-tutorial" class="screen"><div class="tutorial-content" id="tutorial-content"></div><div style="text-align:center;margin-top:20px" id="tutorial-actions"></div></div>
<div id="screen-intro" class="screen"><div class="intro-overlay"><div class="intro-title">HEINEN</div><div class="intro-slogan">Das einzige Spiel mit Spaß-Garantie!<br>Kein Spaß? Geld zurück!</div></div></div>
<div id="screen-game" class="screen">
<div class="qc"><div class="qh"><div class="qn" id="q-counter">Frage 1/10</div><div class="timer" id="q-timer">20</div></div>
<div id="q-report-area" style="text-align:right;margin-bottom:4px;min-height:28px"></div>
<div class="qt" id="q-text"></div><div class="og" id="q-options"></div></div>
<div class="ri" id="results-info" style="display:none"></div>
<div class="pst" id="game-players-title">Spieler*innen</div><div class="pg" id="game-players"></div>
<div class="host-controls" id="host-controls" style="display:none"><button class="btn btn-d" onclick="endGameEarly()">Spiel vorzeitig beenden</button></div></div>
<div id="screen-end" class="screen"><div class="end-container"><div id="end-content"></div><div id="end-players"></div><div id="end-actions" style="margin-top:20px"></div></div></div>
<div id="screen-error" class="screen"><div class="end-container" style="max-width:500px">
<div style="font-size:3rem;margin-bottom:16px">&#9888;</div>
<div style="font-size:1.3rem;font-weight:700;color:var(--wrong);margin-bottom:14px">Fehler bei der Fragengenerierung</div>
<div id="error-msg" style="background:var(--bg2);border:1px solid var(--border);border-radius:10px;padding:16px;margin-bottom:20px;color:var(--text2);font-size:.85rem;font-family:'Space Mono',monospace;word-break:break-word"></div>
<button class="btn btn-p" onclick="leaveGame()" style="max-width:280px;margin:0 auto">Zur Lobby zurückkehren</button>
</div></div>
</div>

<div id="screen-questions">
<div class="qe-page-header">
<button class="nav-btn" onclick="hideQuestions()" style="flex-shrink:0">&#8592; Zurück</button>
<div class="qe-page-title">FRAGEN-EXPLORER</div>
<div class="qe-tab-bar">
<button class="qe-tab-btn active" id="qe-tab-btn-questions" onclick="showQETab('questions')">Fragen</button>
<button class="qe-tab-btn" id="qe-tab-btn-topics" onclick="showQETab('topics')">Themen</button>
<button class="qe-tab-btn" id="qe-tab-btn-dupes" onclick="showQETab('dupes')">Duplikate</button>
<button class="qe-tab-btn" id="qe-tab-btn-reported" onclick="showQETab('reported')">Gemeldete Fragen</button>
</div>
<div id="qe-stats" class="qe-stats" style="display:none;margin-bottom:0;flex:1;border:none;padding:0 0 0 12px;background:none;min-width:0"></div>
</div>
<div class="qe-body">
<div id="qe-view-questions">
<div class="qe-filters">
<input type="text" id="qe-search" placeholder="&#128269; Frage suchen oder #ID..." oninput="qeDebounce()"/>
<input type="text" id="qe-topic" placeholder="&#127991; Thema..." oninput="qeDebounce()"/>
<select id="qe-diff" onchange="loadQuestions(1)"><option value="">Alle Schwierigkeiten</option><option value="leicht">Leicht</option><option value="mittel">Mittel</option><option value="schwer">Schwer</option><option value="extrem">Extrem</option></select>
<select id="qe-provider" onchange="loadQuestions(1)"><option value="">Alle KI-Anbieter</option><option value="openai">OpenAI</option><option value="anthropic">Anthropic</option><option value="Manuell">Manuell</option></select>
<select id="qe-numopts" onchange="loadQuestions(1)"><option value="">Alle Optionszahlen</option><option value="2">2 Optionen</option><option value="3">3 Optionen</option><option value="4">4 Optionen</option></select>
</div>
<div style="display:flex;gap:8px;flex-wrap:wrap;margin-bottom:12px">
<button class="btn btn-sm btn-p" onclick="qeGenOpen()" style="width:auto">+ KI generieren</button>
<button class="btn btn-sm btn-s" onclick="qeAddOpen()" style="width:auto">+ Manuell hinzuf&uuml;gen</button>
</div>
<div id="qe-meta" class="qe-count"></div>
<div id="qe-results" class="qe-grid"></div>
<div id="qe-pagination" class="qe-pagination"></div>
<div style="margin-top:16px;display:flex;gap:8px;flex-wrap:wrap">
<button class="btn btn-sm btn-d" onclick="qeDeleteFiltered()" style="width:auto">Gefilterte l&ouml;schen</button>
<button class="btn btn-sm btn-d" onclick="qeDeleteAll()" style="width:auto">Alle l&ouml;schen</button>
</div>
</div>
<div id="qe-view-topics" style="display:none">
<div id="qe-topics-meta" class="qe-count" style="margin-bottom:12px"></div>
<div id="qe-topics-grid" class="qe-topics-grid"></div>
</div>
<div id="qe-view-dupes" style="display:none">
<div id="qe-dupes-meta" class="qe-count" style="margin-bottom:12px"></div>
<div style="display:flex;gap:8px;align-items:center;margin-bottom:10px;flex-wrap:wrap">
<label style="font-size:.78rem;color:var(--text2)">Schwelle:
<input type="range" id="qe-dupe-threshold" min="0.6" max="1.0" step="0.01" value="0.75" style="width:120px;vertical-align:middle;cursor:pointer;accent-color:var(--accent)" oninput="document.getElementById('qe-dupe-thr-val').textContent=parseFloat(this.value).toFixed(2)"/>
<span id="qe-dupe-thr-val" style="font-family:'Space Mono',monospace;font-size:.78rem;color:var(--text);min-width:32px;display:inline-block">0.75</span>
</label>
<button class="btn btn-sm btn-s" onclick="loadDuplicates()">Suchen</button>
</div>
<div style="display:flex;gap:8px;flex-wrap:wrap;margin-bottom:6px;align-items:center">
<button class="btn btn-sm btn-s" id="qe-dupe-sa-btn" onclick="qeDupeToggleSameAnswer()">&#9654; Gleiche richtige Antwort</button>
<button class="btn btn-sm btn-d" id="qe-dupe-del-btn" style="display:none" onclick="qeDeleteDupeRandom()">&#127922; Jeweils eine l&ouml;schen (Zufall)</button>
</div>
<div style="display:flex;gap:8px;flex-wrap:wrap;margin-bottom:14px;align-items:center">
<span style="font-size:.78rem;color:var(--text2)">&#196;hnlichkeit &ge;</span>
<input type="number" id="qe-dupe-auto-thr" value="95" min="60" max="100" step="1" class="qe-edit-field" style="width:58px;padding:4px 8px;text-align:center" oninput="qeDupeAutoThrUpdate()"/>
<span id="qe-dupe-auto-thr-lbl" style="font-size:.78rem;color:var(--text2)">% &ndash; <b id="qe-dupe-auto-count" style="color:var(--text)">?</b> Paare</span>
<button class="btn btn-sm btn-d" id="qe-dupe-sim-btn" onclick="qeDeleteDupeSimilar()" style="display:none">&#127922; Jeweils eine l&ouml;schen (Zufall)</button>
</div>
<div id="qe-dupes-list"></div>
</div>
<div id="qe-view-reported" style="display:none">
<div id="qe-reported-meta" class="qe-count" style="margin-bottom:12px"></div>
<div id="qe-reported-list"></div>
</div>
</div>
</div>

<div id="qe-gen-modal" class="qe-modal-overlay">
<div class="qe-modal-box">
<div class="qe-modal-title">KI-FRAGEN GENERIEREN</div>
<div id="qe-gen-form">
<div style="margin-bottom:10px"><label style="font-size:.78rem;color:var(--text2);display:block;margin-bottom:4px">Thema</label><input type="text" id="qe-gen-topic" class="qe-edit-field" style="width:100%;box-sizing:border-box" placeholder="z.B. Geschichte, Wissenschaft, Geographie…"/></div>
<div style="display:flex;gap:10px;margin-bottom:10px">
<div style="flex:1"><label style="font-size:.78rem;color:var(--text2);display:block;margin-bottom:4px">Schwierigkeit</label><select id="qe-gen-diff" class="qe-edit-field" style="width:100%"><option value="leicht">Leicht</option><option value="mittel" selected>Mittel</option><option value="schwer">Schwer</option><option value="extrem">Extrem</option></select></div>
<div style="flex:1"><label style="font-size:.78rem;color:var(--text2);display:block;margin-bottom:4px">Fragen / Batch</label><input type="number" id="qe-gen-count" class="qe-edit-field" style="width:100%;box-sizing:border-box" value="10" min="1" max="50" oninput="qeGenUpdateTotal()"/></div>
<div style="flex:1"><label style="font-size:.78rem;color:var(--text2);display:block;margin-bottom:4px">Batches</label><input type="number" id="qe-gen-batches" class="qe-edit-field" style="width:100%;box-sizing:border-box" value="1" min="1" max="20" oninput="qeGenUpdateTotal()"/></div>
<div style="flex:1"><label style="font-size:.78rem;color:var(--text2);display:block;margin-bottom:4px">Antwortoptionen</label><select id="qe-gen-numopts" class="qe-edit-field" style="width:100%"><option value="2">2</option><option value="3">3</option><option value="4" selected>4</option></select></div>
</div>
<div id="qe-gen-total" style="font-size:.72rem;color:var(--text2);margin-bottom:2px;display:none">&#8594; 1 Batch &times; 10 Fragen = <b style="color:var(--text)">10 Fragen</b> gesamt</div>
<div style="margin-top:12px"><label style="font-size:.82rem;color:var(--text2);display:flex;align-items:center;gap:7px;cursor:pointer"><input type="checkbox" id="qe-gen-websearch" style="accent-color:var(--accent);width:14px;height:14px"/>Internet-Recherche f&uuml;r Fragen (nur OpenAI, langsamer)</label></div>
<div style="display:flex;gap:8px;margin-top:14px">
<button class="btn btn-p" onclick="qeGenSubmit()" style="flex:1">Generieren</button>
<button class="btn btn-s" onclick="qeGenClose()" style="flex:1">Abbrechen</button>
</div>
</div>
<div id="qe-gen-loading" style="display:none;text-align:center;padding:12px 0 4px">
<div class="tooth-mount tooth-spinner-mini" style="margin:0 auto 8px"></div>
<div style="background:var(--bg);border-radius:6px;height:5px;overflow:hidden;margin:0 0 12px;width:100%"><div id="qe-gen-bar" style="height:100%;background:linear-gradient(90deg,var(--accent),var(--gold));border-radius:6px;width:5%;transition:width 2.2s ease"></div></div>
<div id="qe-gen-progress" class="loading-text"></div>
</div>
<div id="qe-gen-error" style="display:none">
<div style="font-size:2.5rem;text-align:center;margin-bottom:10px">&#9888;</div>
<div style="font-size:.95rem;font-weight:700;color:var(--wrong);text-align:center;margin-bottom:12px">Fehler bei der Generierung</div>
<div id="qe-gen-errmsg" style="background:var(--bg);border:1px solid var(--border);border-radius:10px;padding:14px;margin-bottom:16px;color:var(--text2);font-size:.75rem;font-family:'Space Mono',monospace;word-break:break-word;max-height:200px;overflow-y:auto;white-space:pre-wrap"></div>
<div style="display:flex;gap:8px">
<button class="btn btn-s" onclick="qeGenShowForm()" style="flex:1">&#8592; Zurück</button>
<button class="btn btn-d" onclick="qeGenClose()" style="flex:1">Schließen</button>
</div>
</div>
</div>
</div>

<div id="qe-add-modal" class="qe-modal-overlay">
<div class="qe-modal-box">
<div class="qe-modal-title">FRAGE HINZUF&Uuml;GEN</div>
<div class="qe-edit-row" style="margin-bottom:12px">
<select id="qe-add-diff" class="qe-edit-field"><option value="leicht">Leicht</option><option value="mittel" selected>Mittel</option><option value="schwer">Schwer</option><option value="extrem">Extrem</option></select>
<input type="text" id="qe-add-topic" class="qe-edit-field" placeholder="Thema…" style="flex:1;min-width:80px"/>
<select id="qe-add-numopts" class="qe-edit-field" onchange="qeAddUpdateOpts()"><option value="2">2 Opt.</option><option value="3">3 Opt.</option><option value="4" selected>4 Opt.</option></select>
</div>
<textarea id="qe-add-text" class="qe-edit-ta" placeholder="Fragetext…"></textarea>
<div id="qe-add-opts" class="qe-opts" style="margin-bottom:10px;display:grid;grid-template-columns:1fr;gap:6px"></div>
<div style="display:flex;gap:8px;margin-top:16px">
<button class="btn btn-p" onclick="qeAddSubmit()" style="flex:1">Hinzuf&uuml;gen</button>
<button class="btn btn-s" onclick="qeAddClose()" style="flex:1">Abbrechen</button>
</div>
</div>
</div>

<div id="audio-controls" style="position:fixed;bottom:20px;right:20px;z-index:90;display:none;align-items:center;gap:6px;background:var(--surface);border:1px solid var(--border);border-radius:22px;padding:4px 10px">
<button class="mute-btn" id="mute-btn" onclick="toggleMute()" title="Stumm schalten" style="position:static;display:flex;width:32px;height:32px;font-size:1rem;border:none;background:none">&#128266;</button>
<input type="range" id="bg-vol-slider" min="0" max="1" step="0.05" value="0.2" style="width:80px;cursor:pointer" oninput="adjustBgVol(this.value)"/>
</div>
<div class="toast" id="toast"></div>
<svg width="0" height="0" style="position:absolute;width:0;height:0" aria-hidden="true" focusable="false">
<defs>
<radialGradient id="tooth-grad" cx="34%" cy="26%" r="92%" fx="30%" fy="20%">
<stop offset="0" stop-color="#fffceb"/>
<stop offset=".35" stop-color="#f4ead0"/>
<stop offset=".75" stop-color="#cbbf99"/>
<stop offset="1" stop-color="#7d7152"/>
</radialGradient>
<linearGradient id="tooth-shade" x1="-50" y1="0" x2="50" y2="0" gradientUnits="userSpaceOnUse">
<stop offset="0" stop-color="rgba(0,0,0,0)"/>
<stop offset=".55" stop-color="rgba(0,0,0,0)"/>
<stop offset="1" stop-color="rgba(60,40,10,.45)"/>
</linearGradient>
<symbol id="tooth-icon" viewBox="-50 -65 100 130">
<path d="M-38 -52 C-44 -50 -46 -38 -46 -22 L-46 5 C-46 12 -40 17 -32 17 L-24 17 L-30 55 C-30 60 -25 62 -22 58 L-12 22 C-10 18 -6 17 0 17 C6 17 10 18 12 22 L22 58 C25 62 30 60 30 55 L24 17 L32 17 C40 17 46 12 46 5 L46 -22 C46 -38 44 -50 38 -52 C30 -58 18 -60 0 -60 C-18 -60 -30 -58 -38 -52 Z" fill="url(#tooth-grad)" stroke="#5a5037" stroke-width="1.8" stroke-linejoin="round"/>
<path d="M-38 -52 C-44 -50 -46 -38 -46 -22 L-46 5 C-46 12 -40 17 -32 17 L-24 17 L-30 55 C-30 60 -25 62 -22 58 L-12 22 C-10 18 -6 17 0 17 C6 17 10 18 12 22 L22 58 C25 62 30 60 30 55 L24 17 L32 17 C40 17 46 12 46 5 L46 -22 C46 -38 44 -50 38 -52 C30 -58 18 -60 0 -60 C-18 -60 -30 -58 -38 -52 Z" fill="url(#tooth-shade)" stroke="none"/>
<path d="M-30 -52 C-38 -45 -40 -30 -34 -16 C-20 -20 -16 -38 -23 -55 C-26 -54 -28 -53 -30 -52 Z" fill="rgba(255,255,255,.55)" stroke="none"/>
<path d="M-2 22 C-1 30 -1 40 -3 55 M2 22 C1 30 1 40 3 55" stroke="rgba(70,55,30,.35)" stroke-width="1" fill="none" stroke-linecap="round"/>
</symbol>
</defs>
</svg>
<audio id="bg-audio" preload="auto" loop></audio>
<audio id="gen-audio" preload="auto" loop></audio>
<script src="https://cdnjs.cloudflare.com/ajax/libs/qrcodejs/1.0.0/qrcode.min.js"></script>
<script>
let ws=null,myId='',inviteCode='',gameState=null,selectedAnswer=-1,timerInterval=null,currentTimeLeft=0,shuffleAnimating=false;
let currentUser=null,joinPending='',joinNeedsPw=false,joinLobbyPw='';
let globalSounds={},bgMuted=false,bgStarted=false,tutorialHtml='';
const DL={leicht:'Leicht',mittel:'Mittel',schwer:'Schwer',extrem:'Extrem'};
const ML={classic:'Klassisch',elimination:'Elimination',kfo_battle_royale:'KFO Battle Royale',kfo_singleplayer:'KFO Singleplayer'};
const openaiModels=['gpt-5.4','gpt-5.4-mini','gpt-5.4-nano','gpt-5','gpt-4o','gpt-4o-mini','gpt-4.1'];
const anthropicModels=['claude-opus-4-6','claude-sonnet-4-6','claude-haiku-4-5-20251001'];
const soundDefs=[{key:'intro_sound',label:'Intro-Sound',id:'file-intro'},{key:'background_sound',label:'Background-Song',id:'file-bg'},{key:'wrong_sound',label:'Falsch-Sound',id:'file-wrong'},{key:'answer_sound',label:'Antwort-Sound',id:'file-answer'},{key:'hurry_sound',label:'Zeit-läuft-ab-Sound',id:'file-hurry'},{key:'timeout_sound',label:'Zeit-abgelaufen-Sound',id:'file-timeout'},{key:'question_sound',label:'Nächste-Frage-Sound',id:'file-question'},{key:'allwrong_sound',label:'Alle-falsch-Sound',id:'file-allwrong'},{key:'allcorrect_sound',label:'Alle-richtig-Sound',id:'file-allcorrect'},{key:'generating_sound',label:'Generierungs-Musik',id:'file-generating'}];
const volDefs=[{id:'vol-intro',key:'vol_intro',label:'Intro',def:'0.6'},{id:'vol-bg',key:'vol_background',label:'Hintergrund',def:'0.2'},{id:'vol-wrong',key:'vol_wrong',label:'Falsch',def:'0.6'},{id:'vol-answer',key:'vol_answer',label:'Antwort',def:'0.6'},{id:'vol-hurry',key:'vol_hurry',label:'Zeit läuft ab',def:'0.5'},{id:'vol-timeout',key:'vol_timeout',label:'Zeit abgelaufen',def:'0.6'},{id:'vol-question',key:'vol_question',label:'Nächste Frage',def:'0.5'},{id:'vol-allwrong',key:'vol_allwrong',label:'Alle falsch',def:'0.6'},{id:'vol-allcorrect',key:'vol_allcorrect',label:'Alle richtig',def:'0.6'},{id:'vol-generating',key:'vol_generating',label:'Generierung',def:'0.4'}];
const SOUND_MAP={intro_sound:'introSound',background_sound:'backgroundSound',wrong_sound:'wrongSound',answer_sound:'answerSound',hurry_sound:'hurrySound',timeout_sound:'timeoutSound',question_sound:'questionSound',allwrong_sound:'allwrongSound',allcorrect_sound:'allcorrectSound',generating_sound:'generatingSound'};
function buildToothSpinner(){
  const tooth='<div class="tooth-face"><svg viewBox="-50 -65 100 130"><use href="#tooth-icon"/></svg></div>';
  let dots='<div class="orbit-dot head"></div>';
  // Trail dots: negative --a angles (wrapper rotates +Y, so the trail behind the head sits at negative wrapper-local angles)
  const cfg=[[8,.85],[16,.7],[24,.58],[32,.48],[40,.4],[50,.32],[60,.25],[72,.18],[86,.12],[102,.07],[120,.04]];
  cfg.forEach(c=>{dots+='<div class="orbit-dot" style="--a:-'+c[0]+'deg;--o:'+c[1]+'"></div>'});
  return '<div class="tooth-stage"><div class="tooth-iso"><div class="tooth-rot">'+tooth+'</div><div class="orbit-rot">'+dots+'</div></div></div>'
}
function mountToothSpinners(){const html=buildToothSpinner();document.querySelectorAll('.tooth-mount').forEach(el=>{if(!el.firstChild)el.innerHTML=html})}
function getCookie(name){const m=document.cookie.match('(^|;)\\s*'+name+'\\s*=\\s*([^;]+)');return m?m[2]:''}
async function apiFetch(url,opts={}){const method=opts.method||'GET';const headers=opts.headers||{};if(method!=='GET'&&!opts.nocsrf){headers['X-CSRF-Token']=getCookie('heinen_csrf')}return fetch(url,{...opts,headers})}

(async function(){
  mountToothSpinners();
  const p=new URLSearchParams(location.search);const inv=p.get('join');
  await loadSounds();await loadTutorial();
  if(inv){joinPending=inv;showScreen('join');return}
  // Try auto-reconnect
  const ri=sessionStorage.getItem('h_invite'),rn=sessionStorage.getItem('h_name');
  if(ri&&rn){connectWS();const ck=setInterval(()=>{if(ws&&ws.readyState===1){clearInterval(ck);send('reconnect',{Name:rn,InviteCode:ri})}},100)}
  try{const r=await fetch('/api/me');if(r.ok){currentUser=await r.json()}}catch(e){}
  showScreen('home');updateNav();loadLobbies();
})();

async function loadSounds(){try{const r=await fetch('/api/global-sounds');const d=await r.json();globalSounds={introSound:d.intro_sound||'',backgroundSound:d.background_sound||'',wrongSound:d.wrong_sound||'',answerSound:d.answer_sound||'',hurrySound:d.hurry_sound||'',timeoutSound:d.timeout_sound||'',questionSound:d.question_sound||'',allwrongSound:d.allwrong_sound||'',allcorrectSound:d.allcorrect_sound||'',generatingSound:d.generating_sound||'',volIntro:parseFloat(d.vol_intro)||0.6,volBg:parseFloat(d.vol_background)||0.2,volWrong:parseFloat(d.vol_wrong)||0.6,volAnswer:parseFloat(d.vol_answer)||0.6,volHurry:parseFloat(d.vol_hurry)||0.5,volTimeout:parseFloat(d.vol_timeout)||0.6,volQuestion:parseFloat(d.vol_question)||0.5,volAllwrong:parseFloat(d.vol_allwrong)||0.6,volAllcorrect:parseFloat(d.vol_allcorrect)||0.6,volGenerating:parseFloat(d.vol_generating)||0.4}}catch(e){}}
async function loadTutorial(){try{const r=await fetch('/api/tutorial');const d=await r.json();tutorialHtml=markdownToHtml(d.content||'')}catch(e){}}
async function loadLobbies(){try{const r=await fetch('/api/lobbies');const lobbies=await r.json();const el=document.getElementById('lobby-list');
  if(!lobbies||lobbies.length===0){el.innerHTML='<div style="color:var(--text2);text-align:center;padding:20px">Keine offenen Lobbys vorhanden.</div>';return}
  el.innerHTML=lobbies.map(l=>'<div class="lobby-card"><div class="lc-info"><div class="lc-name">'+esc(l.lobbyName)+' <span class="lc-badge '+(l.lobbyMode==='open'?'open':'pw')+'">'+(l.lobbyMode==='open'?'Offen':'Passwort')+'</span></div><div class="lc-details">'+esc(l.topic)+' \u2022 '+(ML[l.mode]||l.mode)+' \u2022 '+l.players+' Spieler*in(nen)</div></div><button class="btn btn-sm btn-p" onclick="joinLobby(\''+l.inviteCode+'\',\''+(l.lobbyMode==='password'?'pw':'open')+'\')">Beitreten</button></div>').join('')
}catch(e){}}

function markdownToHtml(md){md=md.replace(/\\\*/g,'\u2605STAR\u2605');md=md.replace(/^### (.+)$/gm,'<h3>$1</h3>').replace(/^## (.+)$/gm,'<h2>$1</h2>').replace(/^# (.+)$/gm,'<h1>$1</h1>').replace(/\*\*(.+?)\*\*/g,'<strong>$1</strong>').replace(/\*(.+?)\*/g,'<em>$1</em>').replace(/^- (.+)$/gm,'<li>$1</li>').replace(/(<li>[\s\S]*?<\/li>)/g,function(m){return '<ul>'+m+'</ul>'}).replace(/<\/ul>\s*<ul>/g,'').replace(/^(?!<[hulo])(.*\S.*)$/gm,'<p>$1</p>');md=md.replace(/\u2605STAR\u2605/g,'*');return md}

function updateNav(){const n=document.getElementById('nav-bar');const hpb=document.getElementById('home-panel-btn');const hlb=document.getElementById('home-login-btn');
  if(currentUser){n.innerHTML='<span style="color:var(--text2);font-size:.75rem;align-self:center">'+esc(currentUser.username)+'</span><button class="nav-btn" onclick="doLogout()">Abmelden</button>';if(hpb)hpb.style.display='inline-block';if(hlb)hlb.style.display='none'}
  else{n.innerHTML='';if(hpb)hpb.style.display='none';if(hlb)hlb.style.display='inline-block'}}
function goHome(){hideQuestions();if(gameState&&['question','results','intro','loading','refill','tutorial'].includes(gameState.phase)){if(!confirm('Spiel wirklich verlassen?'))return}
  gameState=null;stopBg();stopGenSound();if(ws){try{ws.close()}catch(e){}}ws=null;myId='';inviteCode='';sessionStorage.removeItem('h_invite');sessionStorage.removeItem('h_name');
  history.replaceState(null,'','/');showScreen('home');updateNav();loadLobbies()}
function showScreen(name){document.querySelectorAll('.screen').forEach(s=>s.classList.remove('active'));const el=document.getElementById('screen-'+name);if(el)el.classList.add('active');
  if(name==='dashboard'&&currentUser){document.getElementById('welcome-text').textContent='Willkommen, '+currentUser.username+'!';document.getElementById('btn-admin').style.display=currentUser.isAdmin?'block':'none'}
  if(name==='home'){document.getElementById('home-tutorial').innerHTML=tutorialHtml||'';updateNav()}
  document.getElementById('audio-controls').style.display=['game','intro','tutorial'].includes(name)?'flex':'none'}
function toggleQC(){const m=document.getElementById('create-mode').value;const isKFO=m==='kfo_battle_royale'||m==='kfo_singleplayer';document.getElementById('qc-group').style.display=m==='classic'?'block':'none';document.getElementById('sd-group').style.display=isKFO?'block':'none'}
function toggleLobbyPw(){document.getElementById('lobby-pw-group').style.display=document.getElementById('create-lobbymode').value==='password'?'block':'none'}

async function doLogin(){const u=document.getElementById('login-user').value.trim(),p=document.getElementById('login-pass').value;if(!u||!p){showToast('Felder ausfüllen',1);return}
  const r=await apiFetch('/api/login',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({username:u,password:p})});const d=await r.json();if(!r.ok){showToast(d.error||'Fehler',1);return}currentUser=d;showScreen('dashboard');updateNav()}
async function doLogout(){await apiFetch('/api/logout',{method:'POST'});currentUser=null;showScreen('home');updateNav()}
async function changePw(){const o=document.getElementById('pw-old').value,n=document.getElementById('pw-new').value,n2=document.getElementById('pw-new2').value;if(!o||!n){showToast('Felder ausfüllen',1);return}if(n!==n2){showToast('Nicht überein',1);return}
  const r=await apiFetch('/api/change-password',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({OldPassword:o,NewPassword:n})});const d=await r.json();if(d.error){showToast(d.error,1);return}showToast('Geändert',0);showScreen('dashboard')}

// Admin
async function loadAdmin(){
  const ar=await fetch('/api/ai-config');const ac=await ar.json();document.getElementById('ai-provider').value=ac.provider||'openai';onProviderChange();document.getElementById('ai-model').value=ac.model||'';
  document.getElementById('ai-openai-key').value='';document.getElementById('ai-openai-key').placeholder=ac.openaiKey||'sk-...';document.getElementById('ai-anthropic-key').value='';document.getElementById('ai-anthropic-key').placeholder=ac.anthropicKey||'sk-ant-...';document.getElementById('ai-intro-delay').value=ac.introDelay||'4';document.getElementById('ai-status').textContent='';
  // Volumes
  document.getElementById('vol-sliders').innerHTML=volDefs.map(v=>{const val=ac[v.key]||v.def;return '<div class="sound-row"><span class="sound-label">'+v.label+'</span><input type="range" id="'+v.id+'" min="0" max="1" step="0.05" value="'+val+'" style="flex:1" oninput="document.getElementById(\''+v.id+'-val\').textContent=parseFloat(this.value).toFixed(2)"/><span id="'+v.id+'-val" style="font-size:.75rem;color:var(--text2);width:35px;text-align:right">'+val+'</span><button class="btn btn-sm btn-s" onclick="previewVol(\''+v.key.replace('vol_','')+'_sound\',\''+v.id+'\')">&#9654;</button></div>'}).join('');
  // Sounds
  const sr=await fetch('/api/sounds');const sd=await sr.json();document.getElementById('sound-sections').innerHTML=soundDefs.map(s=>{const exists=sd[s.key];return '<div style="margin-bottom:14px"><div class="sound-row"><span class="sound-label">'+s.label+'</span><span class="sound-status'+(exists?'':' none')+'">'+(exists?'\u2714 Hochgeladen':'Nicht hinterlegt')+'</span></div><div class="sound-row"><input type="file" id="'+s.id+'" accept=".mp3,.wav" style="font-size:.8rem;color:var(--text2)"/><button class="btn btn-sm btn-p" onclick="uploadSound(\''+s.key+'\',\''+s.id+'\')">Hochladen</button><button class="btn btn-sm btn-s" onclick="previewSound(\''+s.key+'\')">Probe</button><button class="btn btn-sm btn-d" onclick="deleteSound(\''+s.key+'\')">Entfernen</button></div></div>'}).join('');
  const ur=await fetch('/api/users');const users=await ur.json();renderUsers(users||[]);loadLogs();loadAdminGames();
  try{const qr=await fetch('/api/admin/questions?page=1&limit=1');const qd=await qr.json();const el=document.getElementById('qe-entry-stats');if(el&&qd.totalAll!=null)el.textContent=qd.totalAll.toLocaleString('de-DE')+' Fragen gespeichert'}catch(e){}
}
function onProviderChange(){const p=document.getElementById('ai-provider').value;document.getElementById('ai-openai-row').style.display=p==='openai'?'flex':'none';document.getElementById('ai-anthropic-row').style.display=p==='anthropic'?'flex':'none';const sel=document.getElementById('ai-model');const models=p==='anthropic'?anthropicModels:openaiModels;const cur=sel.value;sel.innerHTML=models.map(m=>'<option value="'+m+'">'+m+'</option>').join('');if(models.includes(cur))sel.value=cur}
async function saveAIConfig(){const b={Provider:document.getElementById('ai-provider').value,Model:document.getElementById('ai-model').value,IntroDelay:document.getElementById('ai-intro-delay').value};
  volDefs.forEach(v=>{const el=document.getElementById(v.id);if(el)b[v.key.split('_').map((w,i)=>i?w[0].toUpperCase()+w.slice(1):w.charAt(0).toUpperCase()+w.slice(1)).join('')]=el.value});
  const ok=document.getElementById('ai-openai-key').value.trim();if(ok)b.OpenaiKey=ok;const ak=document.getElementById('ai-anthropic-key').value.trim();if(ak)b.AnthropicKey=ak;
  await apiFetch('/api/ai-config',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(b)});showToast('Gespeichert',0);await loadSounds();loadAdmin()}
async function testAI(){const st=document.getElementById('ai-status');st.textContent='Teste...';const p=document.getElementById('ai-provider').value;const k=p==='anthropic'?document.getElementById('ai-anthropic-key').value.trim():document.getElementById('ai-openai-key').value.trim();const m=document.getElementById('ai-model').value;const r=await apiFetch('/api/test-ai',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({Provider:p,Key:k,Model:m})});const d=await r.json();st.textContent=d.message;st.style.color=d.ok?'var(--correct)':'var(--wrong)'}
function renderUsers(users){document.getElementById('users-list').innerHTML=users.map(u=>{const dis=u.isOnlyAdmin;const hashClass=u.passwordHashType||'unknown';return '<div class="user-row"><div class="user-info"><span>'+esc(u.username)+'</span><span class="ub '+(u.isAdmin?'admin':'user')+'">'+(u.isAdmin?'Admin':'User')+'</span><span class="ub hash '+hashClass+'">'+hashClass+'</span></div><div class="user-actions">'+(u.isAdmin?'<button class="ib"'+(dis?' disabled':' onclick="toggleAdmin('+u.id+',false)"')+'>Degradieren</button>':'<button class="ib" onclick="toggleAdmin('+u.id+',true)">Zum Admin</button>')+'<button class="ib" onclick="resetUserPassword('+u.id+',\''+esc(u.username)+'\')">PW zurücksetzen</button><button class="ib danger"'+(dis?' disabled':' onclick="deleteUser('+u.id+',\''+esc(u.username)+'\')"')+'>Löschen</button></div></div>'}).join('')}
async function resetUserPassword(id,name){if(!confirm('Passwort für '+name+' auf ein zufälliges zurücksetzen?'))return;const chars='ABCDEFGHJKLMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz23456789';let pw='';const arr=new Uint8Array(8);crypto.getRandomValues(arr);arr.forEach(b=>pw+=chars[b%chars.length]);const r=await apiFetch('/api/users',{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify({id,password:pw})});const d=await r.json();if(d.error){showToast(d.error,1);return}alert('Neues Passwort für '+name+':\n\n'+pw+'\n\nDieses Passwort wird nur einmal angezeigt.');loadAdmin()}
async function addUser(){const n=document.getElementById('new-user-name').value.trim(),p=document.getElementById('new-user-pass').value.trim(),a=document.getElementById('new-user-role').value==='1';if(!n||!p){showToast('Felder ausfüllen',1);return}const r=await apiFetch('/api/users',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({username:n,password:p,isAdmin:a})});const d=await r.json();if(d.error){showToast(d.error,1);return}document.getElementById('new-user-name').value='';document.getElementById('new-user-pass').value='';showToast('Angelegt',0);loadAdmin()}
async function deleteUser(id,name){if(!confirm(name+' löschen?'))return;await apiFetch('/api/users',{method:'DELETE',headers:{'Content-Type':'application/json'},body:JSON.stringify({id})});loadAdmin()}
async function toggleAdmin(id,make){await apiFetch('/api/users',{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify({id,isAdmin:make})});loadAdmin()}

function uploadSound(type,inputId){const input=document.getElementById(inputId);const f=input&&input.files[0];if(!f){showToast('Datei auswählen',1);return}if(f.size>20*1024*1024){showToast('Max 20 MB',1);return}
  const fd=new FormData();fd.append('type',type);fd.append('file',f);const xhr=new XMLHttpRequest();xhr.setRequestHeader('X-CSRF-Token',getCookie('heinen_csrf'));xhr.onload=async function(){try{const d=JSON.parse(xhr.responseText);if(d.ok){showToast('Hochgeladen!',0);input.value='';await loadSounds();loadAdmin()}else showToast(d.error||'Fehler',1)}catch(e){showToast('Serverfehler',1)}};xhr.onerror=()=>showToast('Netzwerkfehler',1);xhr.open('POST','/api/sounds');xhr.send(fd)}
async function deleteSound(type){await apiFetch('/api/sounds',{method:'DELETE',headers:{'Content-Type':'application/json'},body:JSON.stringify({type})});showToast('Entfernt',0);await loadSounds();loadAdmin()}
let pAudio=null;async function previewSound(type){await loadSounds();const url=globalSounds[SOUND_MAP[type]];if(!url){showToast('Nicht hinterlegt',1);return}if(pAudio){pAudio.pause();pAudio=null}pAudio=new Audio(url);pAudio.volume=0.6;const p=pAudio.play();if(p&&p.catch)p.catch(e=>showToast('Fehler',1));setTimeout(()=>{if(pAudio){pAudio.pause();pAudio=null}},10000)}
function previewVol(soundKey,sliderId){const url=globalSounds[SOUND_MAP[soundKey]];if(!url){showToast('Nicht hinterlegt',1);return}const vol=parseFloat(document.getElementById(sliderId).value);if(pAudio){pAudio.pause();pAudio=null}pAudio=new Audio(url);pAudio.volume=vol;const p=pAudio.play();if(p&&p.catch)p.catch(()=>{});setTimeout(()=>{if(pAudio){pAudio.pause();pAudio=null}},5000)}
async function loadLogs(){
  try{
    const search=document.getElementById('log-search').value;
    const level=document.getElementById('log-level').value;
    const params=new URLSearchParams();
    if(search)params.set('search',search);
    if(level)params.set('level',level);
    const r=await fetch('/api/logs?'+params.toString());
    const d=await r.json();
    const entries=d.entries||[];
    const c=document.getElementById('log-table-container');
    if(entries.length===0){c.innerHTML='<div class="log-empty">Keine Einträge gefunden.</div>';document.getElementById('log-meta').textContent='0 Einträge';return}
    let html='<table class="log-table"><thead><tr><th>Zeit</th><th>Level</th><th>IP</th><th>User</th><th>Aktion</th><th>Details</th></tr></thead><tbody>';
    for(const e of entries){
      html+='<tr><td>'+esc(e.timestamp)+'</td><td><span class="log-lvl '+esc(e.level)+'">'+esc(e.level)+'</span></td><td>'+esc(e.ip)+'</td><td>'+esc(e.user)+'</td><td class="log-action">'+esc(e.action)+'</td><td>'+esc(e.details||'')+'</td></tr>';
    }
    html+='</tbody></table>';
    c.innerHTML=html;
    document.getElementById('log-meta').textContent=entries.length+' Einträge'+(d.total>entries.length?' (begrenzt auf 1000)':'');
  }catch(e){document.getElementById('log-table-container').innerHTML='<div class="log-empty">Fehler beim Laden.</div>'}
}
async function clearLogs(){if(!confirm('Alle Logs unwiderruflich löschen?'))return;await apiFetch('/api/logs',{method:'DELETE'});showToast('Geleert',0);loadLogs()}

async function loadAdminGames(){
  const container=document.getElementById('games-list-container');
  if(!container)return;
  try{
    const r=await fetch('/api/admin/games');
    if(!r.ok){container.innerHTML='<div class="log-empty">Fehler beim Laden.</div>';return}
    const gs=await r.json();
    if(!gs||!gs.length){container.innerHTML='<div class="log-empty">Noch keine Spiele aufgezeichnet.</div>';return}
    let html='<table class="log-table"><thead><tr><th>Datum</th><th>Host</th><th>Lobby</th><th>Modus</th><th>Thema</th><th>Spieler</th><th>Runden</th><th>Gewinner</th></tr></thead><tbody>';
    for(const g of gs){
      const dt=new Date(g.startedAt*1000).toLocaleString('de-DE');
      const mode=ML[g.mode]||g.mode;
      const winners=(g.winners||[]).join(', ')||'\u2013';
      html+='<tr class="game-row" data-gid="'+esc(g.id)+'" onclick="loadGameDetail(this.dataset.gid)" style="cursor:pointer">';
      html+='<td>'+esc(dt)+'</td><td>'+esc(g.hostUser)+'</td><td>'+esc(g.lobbyName)+'</td>';
      html+='<td>'+esc(mode)+'</td><td>'+esc(g.topic)+'</td>';
      html+='<td>'+g.numPlayers+'</td><td>'+g.totalRounds+'</td>';
      html+='<td style="color:var(--gold)">'+esc(winners)+'</td></tr>'
    }
    html+='</tbody></table>';
    container.innerHTML=html
  }catch(e){container.innerHTML='<div class="log-empty">Fehler beim Laden.</div>'}
}

async function loadGameDetail(id){
  document.getElementById('games-list-view').style.display='none';
  document.getElementById('games-detail-view').style.display='block';
  const c=document.getElementById('games-detail-content');
  c.innerHTML='<div style="padding:20px;color:var(--text2);text-align:center">Lade Details...</div>';
  try{
    const r=await fetch('/api/admin/games/'+encodeURIComponent(id));
    if(!r.ok){c.innerHTML='<div class="log-empty">Fehler beim Laden.</div>';return}
    const g=await r.json();
    renderGameDetail(g,c)
  }catch(e){c.innerHTML='<div class="log-empty">Fehler beim Laden.</div>'}
}

function renderGameDetail(g,container){
  const dt=new Date(g.startedAt*1000).toLocaleString('de-DE');
  const dtEnd=g.endedAt?new Date(g.endedAt*1000).toLocaleString('de-DE'):'\u2013';
  const durSec=g.endedAt?(g.endedAt-g.startedAt):0;
  const durStr=g.endedAt?(Math.floor(durSec/60)+' Min. '+(durSec%60)+' Sek.'):'\u2013';
  const winners=(g.winners||[]).join(', ')||'\u2013';
  const mode=ML[g.mode]||g.mode;
  const diff=DL[g.difficulty]||g.difficulty||'\u2013';
  const endMap={'end':'Beendet','error':'Fehler','abandoned':'Abgebrochen'};
  const endStatus=endMap[g.endPhase]||(g.endPhase||'\u2013');
  let html='<div style="background:var(--bg2);border-radius:10px;padding:16px;margin-bottom:16px"><div class="sv-grid">';
  html+='<div class="sv-item"><span class="sv-label">Start</span><span class="sv-value">'+esc(dt)+'</span></div>';
  html+='<div class="sv-item"><span class="sv-label">Ende</span><span class="sv-value">'+esc(dtEnd)+'</span></div>';
  html+='<div class="sv-item"><span class="sv-label">Dauer</span><span class="sv-value">'+esc(durStr)+'</span></div>';
  html+='<div class="sv-item"><span class="sv-label">Status</span><span class="sv-value">'+esc(endStatus)+'</span></div>';
  html+='<div class="sv-item"><span class="sv-label">Host</span><span class="sv-value">'+esc(g.hostUser)+'</span></div>';
  html+='<div class="sv-item"><span class="sv-label">Lobby</span><span class="sv-value">'+esc(g.lobbyName)+'</span></div>';
  html+='<div class="sv-item"><span class="sv-label">Modus</span><span class="sv-value">'+esc(mode)+'</span></div>';
  html+='<div class="sv-item"><span class="sv-label">Thema</span><span class="sv-value">'+esc(g.topic)+'</span></div>';
  html+='<div class="sv-item"><span class="sv-label">Schwierigkeit</span><span class="sv-value">'+esc(diff)+'</span></div>';
  html+='<div class="sv-item"><span class="sv-label">Z\u00e4hne (Start)</span><span class="sv-value">'+g.numTeeth+'</span></div>';
  html+='<div class="sv-item"><span class="sv-label">Zeit/Frage</span><span class="sv-value">'+g.timePerQ+'s</span></div>';
  html+='<div class="sv-item"><span class="sv-label">Spieler</span><span class="sv-value">'+g.numPlayers+'</span></div>';
  html+='<div class="sv-item" style="grid-column:1/-1"><span class="sv-label">Gewinner</span><span class="sv-value" style="color:var(--gold)">'+esc(winners)+'</span></div>';
  html+='</div></div>';
  const rounds=g.rounds||[];
  if(!rounds.length){html+='<div class="log-empty">Keine Runden aufgezeichnet.</div>';container.innerHTML=html;return}
  const playerOrder=[],playerNames={};
  for(const round of rounds){
    for(const pr of (round.players||[])){
      if(!playerNames[pr.playerId]){playerNames[pr.playerId]=pr.playerName;playerOrder.push(pr.playerId)}
    }
  }
  const OL=['A','B','C','D','E','F'];
  html+='<div style="overflow-x:auto"><table class="log-table"><thead><tr>';
  html+='<th style="min-width:36px;text-align:center">#</th><th>Schw.</th><th style="min-width:220px">Frage &amp; Antworten</th>';
  for(const pid of playerOrder)html+='<th>'+esc(playerNames[pid])+'</th>';
  html+='</tr></thead><tbody>';
  for(const round of rounds){
    const pData={};
    for(const pr of (round.players||[]))pData[pr.playerId]=pr;
    const opts=round.options||[];
    html+='<tr><td style="font-family:\'Space Mono\',monospace;font-weight:700;text-align:center">'+(round.number+1)+'</td>';
    html+='<td>'+esc(DL[round.difficulty]||round.difficulty||'\u2013')+'</td>';
    html+='<td><div style="font-weight:600;margin-bottom:4px;font-size:.85rem">'+esc(round.question)+'</div><div style="font-size:.7rem;color:var(--text2)">';
    for(let i=0;i<opts.length;i++){
      const cor=(i===round.correctAnswer);
      html+='<span style="margin-right:10px'+(cor?';color:var(--correct);font-weight:700':'')+'">'+OL[i]+': '+esc(opts[i])+'</span>'
    }
    html+='</div></td>';
    for(const pid of playerOrder){
      const pr=pData[pid];
      if(!pr){html+='<td style="color:var(--text2);text-align:center">\u2013</td>';continue}
      const rc=pr.result==='correct'?'var(--correct)':(pr.result==='timeout'?'var(--gold)':'var(--wrong)');
      const rl=pr.result==='correct'?'richtig':(pr.result==='timeout'?'Zeit':'falsch');
      const al=pr.answer>=0&&pr.answer<OL.length?OL[pr.answer]:'\u2013';
      const elim=pr.eliminated?'<span style="color:var(--wrong);font-size:.6rem"> elim.</span>':'';
      html+='<td style="min-width:90px">';
      html+='<div style="color:'+rc+';font-size:.8rem;font-weight:600">'+al+' \u2013 '+rl+'</div>';
      html+='<div style="font-size:.65rem;color:var(--text2)">'+pr.teethBefore+' \u2192 '+pr.teethAfter+' Z\u00e4hne'+elim+'</div></td>'
    }
    html+='</tr>'
  }
  html+='</tbody></table></div>';
  const finalTeeth={},finalElim={};
  for(const round of rounds){
    for(const pr of (round.players||[])){finalTeeth[pr.playerId]=pr.teethAfter;finalElim[pr.playerId]=pr.eliminated||false}
  }
  const standings=playerOrder.map(pid=>({id:pid,name:playerNames[pid]||'?',teeth:finalTeeth[pid]||0,elim:finalElim[pid]})).sort((a,b)=>b.teeth-a.teeth);
  if(standings.length){
    html+='<div style="margin-top:16px"><div class="ast">Endstand</div><div style="display:flex;gap:10px;flex-wrap:wrap;margin-top:10px">';
    for(let i=0;i<standings.length;i++){
      const s=standings[i];
      const isW=(g.winners||[]).includes(s.name);
      html+='<div style="background:var(--bg2);border:1px solid '+(isW?'var(--gold)':'var(--border)')+';border-radius:10px;padding:10px 16px;min-width:100px;text-align:center">';
      html+='<div style="font-size:.65rem;color:var(--text2);text-transform:uppercase;letter-spacing:1px;margin-bottom:4px">'+(i+1)+'. Platz</div>';
      html+='<div style="font-weight:700;margin-bottom:4px'+(isW?';color:var(--gold)':'')+'">'+esc(s.name)+'</div>';
      html+='<div style="font-size:.8rem;color:var(--text2)">'+s.teeth+' Z\u00e4hne</div></div>'
    }
    html+='</div></div>'
  }
  container.innerHTML=html
}
function backToGamesList(){
  document.getElementById('games-detail-view').style.display='none';
  document.getElementById('games-list-view').style.display='block'
}

// ── Fragen-Explorer ──────────────────────────────────────────────────────────
let qePage=1,qeTimer=null,qeData={},qeTab='questions',qeDupePairs=[],qeDupeSkipped=new Set(),qReportedCurrent=false;
const OLQ=['A','B','C','D','E','F'];
const DIFF_OPTS=['leicht','mittel','schwer','extrem'];
const DIFF_COLORS={leicht:'#5aafff',mittel:'var(--gold)',schwer:'#ff8c00',extrem:'var(--wrong)'};
function showQuestions(){
  if(!currentUser||!currentUser.isAdmin){showToast('Kein Zugriff',1);return}
  document.getElementById('screen-questions').style.display='block';
  document.body.style.overflow='hidden';
  showQETab('questions')
}
function hideQuestions(){
  document.getElementById('screen-questions').style.display='none';
  document.body.style.overflow=''
}
function showQETab(tab){
  qeTab=tab;
  ['questions','topics','dupes','reported'].forEach(t=>{
    const v=document.getElementById('qe-view-'+t);if(v)v.style.display=t===tab?'block':'none';
    const b=document.getElementById('qe-tab-btn-'+t);if(b)b.className='qe-tab-btn'+(t===tab?' active':'');
  });
  if(tab==='questions')loadQuestions(1);
  else if(tab==='topics')loadTopics();
  else if(tab==='reported')loadReportedQuestions();
}
async function loadTopics(){
  const gridEl=document.getElementById('qe-topics-grid'),metaEl=document.getElementById('qe-topics-meta');
  if(gridEl)gridEl.innerHTML='<div style="color:var(--text2);grid-column:1/-1;padding:20px 0">Lade\u2026</div>';
  try{
    const r=await fetch('/api/admin/questions?view=topics');
    if(!r.ok){if(gridEl)gridEl.innerHTML='<div style="color:var(--wrong);grid-column:1/-1;padding:20px 0">Fehler beim Laden.</div>';return}
    const d=await r.json();
    const topics=d.topics||[];
    if(metaEl)metaEl.textContent=topics.length+' Themen';
    if(!topics.length){if(gridEl)gridEl.innerHTML='<div style="color:var(--text2);grid-column:1/-1;padding:20px 0">Keine Themen gefunden.</div>';return}
    if(gridEl)gridEl.innerHTML=topics.map(t=>{
      const diffs=DIFF_OPTS.filter(k=>t.byDiff&&t.byDiff[k]);
      return '<div class="qe-topic-card" onclick="qeSelectTopic(\''+esc(t.topic)+'\')">'
        +'<div class="qe-topic-card-name" title="'+esc(t.topic)+'">'+esc(t.topic)+'</div>'
        +'<div class="qe-topic-card-count">'+t.count.toLocaleString('de-DE')+' Fragen</div>'
        +'<div class="qe-topic-card-diffs">'+diffs.map(k=>'<span class="qe-tag '+k+'" style="font-size:.52rem">'+DL[k]+' '+t.byDiff[k]+'</span>').join('')+'</div>'
        +'</div>'
    }).join('')
  }catch(e){if(gridEl)gridEl.innerHTML='<div style="color:var(--wrong);grid-column:1/-1;padding:20px 0">Fehler beim Laden.</div>'}
}
function qeSelectTopic(topic){
  showQETab('questions');
  const el=document.getElementById('qe-topic');if(el)el.value=topic;
  loadQuestions(1)
}
function qeDebounce(){clearTimeout(qeTimer);qeTimer=setTimeout(()=>loadQuestions(1),300)}
async function loadQuestions(page){
  qePage=page||1;
  const gv=id=>(document.getElementById(id)||{}).value||'';
  const search=gv('qe-search'),topic=gv('qe-topic'),diff=gv('qe-diff'),provider=gv('qe-provider'),numopts=gv('qe-numopts');
  const params=new URLSearchParams();
  if(search)params.set('search',search);if(topic)params.set('topic',topic);
  if(diff)params.set('difficulty',diff);if(provider)params.set('ai_provider',provider);
  if(numopts)params.set('num_options',numopts);
  params.set('page',qePage);
  const resEl=document.getElementById('qe-results'),metaEl=document.getElementById('qe-meta');
  if(resEl)resEl.innerHTML='<div class="log-empty" style="grid-column:1/-1">Lade\u2026</div>';
  try{
    const r=await fetch('/api/admin/questions?'+params.toString());
    if(!r.ok){if(resEl)resEl.innerHTML='<div class="log-empty" style="grid-column:1/-1">Fehler beim Laden.</div>';return}
    const d=await r.json();
    const statsEl=document.getElementById('qe-stats');
    if(statsEl&&d.totalAll!=null){
      const bd=d.byDifficulty||{};
      statsEl.style.display='flex';
      statsEl.innerHTML='<div class="qe-stat"><span class="qe-stat-val">'+d.totalAll.toLocaleString('de-DE')+'</span><span class="qe-stat-lbl">Gesamt</span></div>'
        +DIFF_OPTS.filter(k=>bd[k]).map(k=>'<div class="qe-stat"><span class="qe-stat-val" style="color:'+DIFF_COLORS[k]+'">'+bd[k].toLocaleString('de-DE')+'</span><span class="qe-stat-lbl">'+DL[k]+'</span></div>').join('');
      const entryEl=document.getElementById('qe-entry-stats');
      if(entryEl)entryEl.textContent=d.totalAll.toLocaleString('de-DE')+' Fragen gespeichert'
    }
    const totalPages=Math.ceil((d.total||0)/(d.limit||20));
    if(metaEl)metaEl.textContent=(d.total||0).toLocaleString('de-DE')+(search||topic||diff||provider?' Treffer':' Fragen')+(d.total!==d.totalAll&&d.totalAll?' ('+d.totalAll.toLocaleString('de-DE')+' gesamt)':'');
    const qs=d.questions||[];
    if(!qs.length){if(resEl)resEl.innerHTML='<div class="log-empty" style="grid-column:1/-1">Keine Fragen gefunden.</div>';const p=document.getElementById('qe-pagination');if(p)p.innerHTML='';return}
    qeData={};for(const q of qs)qeData[q.id]=q;
    if(resEl)resEl.innerHTML=qs.map(q=>renderQECard(q)).join('');
    const pag=document.getElementById('qe-pagination');
    if(!pag||totalPages<=1){if(pag)pag.innerHTML='';return}
    let ph='';
    if(qePage>1)ph+='<button class="btn btn-sm btn-s" onclick="loadQuestions('+(qePage-1)+')">&#8249;</button>';
    const ps=Math.max(1,qePage-2),pe=Math.min(totalPages,qePage+2);
    if(ps>1)ph+='<button class="btn btn-sm btn-s" onclick="loadQuestions(1)">1</button>'+(ps>2?'<span style="color:var(--text2);padding:0 4px">&hellip;</span>':'');
    for(let i=ps;i<=pe;i++)ph+='<button class="btn btn-sm '+(i===qePage?'btn-p':'btn-s')+'" onclick="loadQuestions('+i+')">'+i+'</button>';
    if(pe<totalPages)ph+=(pe<totalPages-1?'<span style="color:var(--text2);padding:0 4px">&hellip;</span>':'')+'<button class="btn btn-sm btn-s" onclick="loadQuestions('+totalPages+')">'+totalPages+'</button>';
    if(qePage<totalPages)ph+='<button class="btn btn-sm btn-s" onclick="loadQuestions('+(qePage+1)+')">&#8250;</button>';
    pag.innerHTML=ph
  }catch(e){if(resEl)resEl.innerHTML='<div class="log-empty" style="grid-column:1/-1">Fehler beim Laden.</div>'}
}
function renderQECard(q){
  const dt=new Date(q.askedAt*1000).toLocaleString('de-DE');
  const diffCls={leicht:'leicht',mittel:'mittel',schwer:'schwer',extrem:'extrem'}[q.difficulty]||'';
  const diffTag=q.difficulty?'<span class="qe-tag '+diffCls+'">'+(DL[q.difficulty]||q.difficulty)+'</span>':'';
  const topicTag=q.topic?'<span class="qe-tag topic">'+esc(q.topic)+'</span>':'';
  const aiCls=q.aiProvider==='openai'?'ai-openai':q.aiProvider==='anthropic'?'ai-anthropic':q.aiProvider==='Manuell'?'ai-manual':'';
  const aiLabel=q.aiProvider?(q.aiProvider+(q.aiModel?'\u00b7'+q.aiModel:'')):'';
  const aiTag=aiLabel?'<span class="qe-tag '+aiCls+'">'+esc(aiLabel)+'</span>':'';
  const optsTag=q.numOptions?'<span class="qe-tag opts">'+q.numOptions+' Optionen</span>':'';
  const opts=(q.options||[]).map((o,i)=>'<div class="qe-opt'+(i===q.correctAnswer?' cor':'')+'">'+'<span class="qe-ol">'+OLQ[i]+'</span><span>'+esc(o)+'</span></div>').join('');
  return '<div class="qe-card" data-qid="'+q.id+'">'
    +'<div style="display:flex;justify-content:space-between;align-items:flex-start;gap:8px;margin-bottom:8px">'
    +'<div style="display:flex;gap:5px;flex-wrap:wrap">'+diffTag+topicTag+aiTag+optsTag+'</div>'
    +'<div style="display:flex;gap:5px;flex-shrink:0">'
    +'<button class="ib" style="font-size:.65rem;padding:2px 9px" onclick="qeEdit('+q.id+')" title="Bearbeiten">&#9998;</button>'
    +'<button class="ib danger" style="font-size:.65rem;padding:2px 9px" onclick="qeDeleteOne('+q.id+')" title="L\u00f6schen">&#10005;</button>'
    +'</div></div>'
    +'<div class="qe-question">'+esc(q.text)+'</div>'
    +'<div class="qe-opts">'+opts+'</div>'
    +'<div class="qe-date">&#128336; '+esc(dt)+' \u00b7 ID\u202f'+q.id+'</div>'
    +'</div>'
}
function qeEdit(id){
  const card=document.querySelector('[data-qid="'+id+'"]');
  if(!card)return;
  const q=qeData[id];if(!q)return;
  const diffOpts=DIFF_OPTS.map(d=>'<option value="'+d+'"'+(d===q.difficulty?' selected':'')+'>'+DL[d]+'</option>').join('');
  const provOpts=['openai','anthropic','Manuell'].map(p=>'<option value="'+p+'"'+(p===q.aiProvider?' selected':'')+'>'+p+'</option>').join('');
  const numOptsOpts=[2,3,4].map(n=>'<option value="'+n+'"'+((q.numOptions||4)===n?' selected':'')+'>'+n+' Opt.</option>').join('');
  const optsHtml=(q.options||[]).map((o,i)=>'<div class="qe-opt-edit'+(i===q.correctAnswer?' sel-cor':'')+'" id="qe-orow-'+i+'-'+id+'">'
    +'<input type="radio" name="qe-cor-'+id+'" value="'+i+'"'+(i===q.correctAnswer?' checked':'')+' onchange="qeUpdateCor('+id+')">'
    +'<span class="qe-ol" style="font-family:\'Space Mono\',monospace;font-weight:700;font-size:.68rem;color:var(--text2);min-width:14px;flex-shrink:0">'+OLQ[i]+'</span>'
    +'<input type="text" id="qe-opt'+i+'-'+id+'" value="'+esc(o)+'" style="flex:1;background:none;border:none;color:var(--text);font-family:Outfit,sans-serif;font-size:.83rem;outline:none;min-width:0" placeholder="Option '+OLQ[i]+'\u2026"/>'
    +'</div>').join('');
  card.innerHTML=
    '<div style="display:flex;justify-content:space-between;align-items:flex-start;gap:8px;margin-bottom:10px">'
    +'<div class="qe-edit-row" style="flex:1;margin-bottom:0;gap:6px">'
    +'<select id="qe-ediff-'+id+'" class="qe-edit-field">'+diffOpts+'</select>'
    +'<input type="text" id="qe-etopic-'+id+'" class="qe-edit-field" value="'+esc(q.topic||'')+'" placeholder="Thema\u2026" style="flex:1;min-width:80px"/>'
    +'<select id="qe-eprov-'+id+'" class="qe-edit-field">'+provOpts+'</select>'
    +'<input type="text" id="qe-emodel-'+id+'" class="qe-edit-field" value="'+esc(q.aiModel||'')+'" placeholder="KI-Modell\u2026" style="min-width:90px"/>'
    +'<select id="qe-enumopts-'+id+'" class="qe-edit-field">'+numOptsOpts+'</select>'
    +'</div>'
    +'<div style="display:flex;gap:5px;flex-shrink:0;margin-left:8px">'
    +'<button class="ib" style="font-size:.75rem;padding:3px 11px;color:var(--correct);border-color:rgba(0,230,138,.5)" onclick="qeSave('+id+')" title="Speichern">&#10003;</button>'
    +'<button class="ib" style="font-size:.75rem;padding:3px 11px" onclick="qeCancel('+id+')" title="Abbrechen">&#10005;</button>'
    +'</div></div>'
    +'<textarea id="qe-etext-'+id+'" class="qe-edit-ta">'+esc(q.text)+'</textarea>'
    +'<div class="qe-opts">'+optsHtml+'</div>'
}
function qeUpdateCor(id){
  document.querySelectorAll('[name="qe-cor-'+id+'"]').forEach((r,i)=>{
    const row=document.getElementById('qe-orow-'+i+'-'+id);
    if(row)row.className='qe-opt-edit'+(r.checked?' sel-cor':'')
  })
}
function qeCancel(id){
  const card=document.querySelector('[data-qid="'+id+'"]');
  if(!card||!qeData[id])return;
  card.outerHTML=qeData[id]._context==='reported'?renderReportedCard(qeData[id]):renderQECard(qeData[id])
}
async function qeSave(id){
  const card=document.querySelector('[data-qid="'+id+'"]');
  if(!card||!qeData[id])return;
  const gv=sid=>(document.getElementById(sid)||{}).value||'';
  const text=gv('qe-etext-'+id),difficulty=gv('qe-ediff-'+id),topic=gv('qe-etopic-'+id);
  const aiProvider=gv('qe-eprov-'+id),aiModel=gv('qe-emodel-'+id);
  const numOptions=parseInt(gv('qe-enumopts-'+id))||0;
  const options=[];let correctAnswer=qeData[id].correctAnswer;
  document.querySelectorAll('[name="qe-cor-'+id+'"]').forEach((r,i)=>{
    options.push(gv('qe-opt'+i+'-'+id));if(r.checked)correctAnswer=i
  });
  try{
    const r=await apiFetch('/api/admin/questions',{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify({id,text,options,correctAnswer,numOptions,difficulty,topic,aiProvider,aiModel})});
    const d=await r.json();
    if(!d.ok){showToast(d.error||'Fehler beim Speichern',1);return}
    const ctx=qeData[id]._context;
    const updated=Object.assign({},qeData[id],{text,options,correctAnswer,numOptions,difficulty,topic,aiProvider,aiModel,_context:ctx});
    qeData[id]=updated;
    card.outerHTML=ctx==='reported'?renderReportedCard(updated):renderQECard(updated);
    showToast('Gespeichert',0)
  }catch(e){showToast('Netzwerkfehler',1)}
}
async function qeDeleteOne(id){
  if(!confirm('Diese Frage l\u00f6schen?'))return;
  const ctx=qeData[id]&&qeData[id]._context;
  const r=await apiFetch('/api/admin/questions',{method:'DELETE',headers:{'Content-Type':'application/json'},body:JSON.stringify({id})});
  const d=await r.json();
  if(d.ok){showToast('Gel\u00f6scht',0);if(ctx==='reported')loadReportedQuestions();else loadQuestions(qePage)}else showToast(d.error||'Fehler',1)
}
async function qeDeleteFiltered(){
  const gv=id=>(document.getElementById(id)||{}).value||'';
  const search=gv('qe-search'),topic=gv('qe-topic'),diff=gv('qe-diff'),provider=gv('qe-provider');
  if(!search&&!topic&&!diff&&!provider){showToast('Bitte zuerst filtern',1);return}
  if(!confirm('Alle gefilterten Fragen l\u00f6schen?'))return;
  const r=await apiFetch('/api/admin/questions',{method:'DELETE',headers:{'Content-Type':'application/json'},body:JSON.stringify({deleteAll:true,topic,difficulty:diff,aiProvider:provider})});
  const d=await r.json();
  if(d.ok){showToast(d.deleted+' Fragen gel\u00f6scht',0);loadQuestions(1)}else showToast(d.error||'Fehler',1)
}
async function qeDeleteAll(){
  if(!confirm('ALLE Fragen unwiderruflich l\u00f6schen?'))return;
  const r=await apiFetch('/api/admin/questions',{method:'DELETE',headers:{'Content-Type':'application/json'},body:JSON.stringify({deleteAll:true})});
  const d=await r.json();
  if(d.ok){showToast(d.deleted+' Fragen gel\u00f6scht',0);loadQuestions(1)}else showToast(d.error||'Fehler',1)
}

// ── KI-Fragen generieren ─────────────────────────────────────────────────────
let _qeGenTimer=null;
function qeGenUpdateTotal(){
  const count=parseInt((document.getElementById('qe-gen-count')||{}).value)||10;
  const batches=parseInt((document.getElementById('qe-gen-batches')||{}).value)||1;
  const el=document.getElementById('qe-gen-total');
  if(!el)return;
  if(batches>1){
    el.style.display='';
    el.innerHTML='&#8594; '+batches+' Batches &times; '+count+' Fragen = <b style="color:var(--text)">'+(batches*count)+' Fragen</b> gesamt';
  }else{el.style.display='none'}
}
function qeGenOpen(){
  qeGenShowForm();
  document.getElementById('qe-gen-modal').style.display='flex';
  setTimeout(()=>document.getElementById('qe-gen-topic').focus(),50)
}
function qeGenClose(){
  clearInterval(_qeGenTimer);_qeGenTimer=null;
  stopGenSound();
  document.getElementById('qe-gen-modal').style.display='none';
  qeGenShowForm()
}
function qeGenShowForm(){
  stopGenSound();
  document.getElementById('qe-gen-form').style.display='';
  document.getElementById('qe-gen-loading').style.display='none';
  document.getElementById('qe-gen-error').style.display='none'
}
function qeGenShowLoading(){
  document.getElementById('qe-gen-form').style.display='none';
  document.getElementById('qe-gen-loading').style.display='';
  document.getElementById('qe-gen-error').style.display='none';
  clearInterval(_qeGenTimer);_qeGenTimer=null;
  const barEl=document.getElementById('qe-gen-bar');
  if(barEl){barEl.style.transition='none';barEl.style.width='0%'}
  startGenSound()
}
function qeGenSetBatchProgress(b,total){
  clearInterval(_qeGenTimer);
  const basePct=Math.round(b/total*90);
  const capPct=Math.round((b+1)/total*90)-1;
  const barEl=document.getElementById('qe-gen-bar');
  const progEl=document.getElementById('qe-gen-progress');
  if(barEl){barEl.style.transition='width .5s ease';barEl.style.width=basePct+'%'}
  const bLabel=total>1?' \u2013 Batch '+(b+1)+'\u202f/\u202f'+total:'';
  const steps=['Verbinde mit KI\u2026','Erstelle Anfrage\u2026','Generiere Fragen\u2026','Pr\u00fcfe Antworten\u2026','Noch einen Moment\u2026'];
  let si=0,pct=basePct;
  if(progEl)progEl.textContent=steps[0]+bLabel;
  _qeGenTimer=setInterval(()=>{
    pct=Math.min(pct+(capPct-basePct)*0.28+1,capPct);
    if(barEl){barEl.style.transition='width 2s ease';barEl.style.width=pct+'%'}
    si=(si+1)%steps.length;
    if(progEl)progEl.textContent=steps[si]+bLabel
  },2400)
}
function qeGenShowError(msg){
  clearInterval(_qeGenTimer);_qeGenTimer=null;
  stopGenSound();
  document.getElementById('qe-gen-form').style.display='none';
  document.getElementById('qe-gen-loading').style.display='none';
  document.getElementById('qe-gen-error').style.display='';
  const el=document.getElementById('qe-gen-errmsg');
  if(el)el.textContent=msg||'Unbekannter Fehler'
}
async function qeGenSubmit(){
  const gv=id=>(document.getElementById(id)||{}).value||'';
  const topic=gv('qe-gen-topic').trim()||'Allgemeinwissen';
  const difficulty=gv('qe-gen-diff')||'mittel';
  const count=parseInt(gv('qe-gen-count'))||10;
  const numOptions=parseInt(gv('qe-gen-numopts'))||4;
  const webSearch=!!(document.getElementById('qe-gen-websearch')||{}).checked;
  const batches=Math.min(Math.max(parseInt(gv('qe-gen-batches'))||1,1),20);
  qeGenShowLoading();
  let totalGenerated=0,totalFiltered=0,batchErrors=0;
  for(let b=0;b<batches;b++){
    qeGenSetBatchProgress(b,batches);
    try{
      const r=await apiFetch('/api/admin/questions/generate',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({topic,difficulty,count,numOptions,webSearch})});
      const d=await r.json();
      if(d.ok){totalGenerated+=d.count;totalFiltered+=(d.filtered||0);}else batchErrors++;
    }catch(e){batchErrors++}
  }
  clearInterval(_qeGenTimer);_qeGenTimer=null;
  const barEl=document.getElementById('qe-gen-bar');
  const progEl=document.getElementById('qe-gen-progress');
  if(barEl){barEl.style.transition='width .4s ease';barEl.style.width='100%'}
  if(totalGenerated===0){
    qeGenShowError('Alle '+batches+' Batches fehlgeschlagen.\nBitte KI-Konfiguration prüfen.');
    return
  }
  const filterNote=totalFiltered>0?' ('+totalFiltered+' Duplikat'+(totalFiltered!==1?'e':'')+' gefiltert)':'';
  const note=batchErrors?' ('+batchErrors+' Batch'+(batchErrors>1?' fehlgeschlagen':'es fehlgeschlagen')+')':'';
  if(progEl)progEl.textContent=totalGenerated+' Fragen gespeichert \u2713'+filterNote+note;
  setTimeout(()=>{
    showToast(totalGenerated+' Fragen generiert'+filterNote+note,batchErrors?1:0);
    qeGenClose();loadQuestions(1)
  },500)
}

// ── Frage manuell hinzufügen ──────────────────────────────────────────────────
function qeAddOpen(){
  document.getElementById('qe-add-text').value='';
  document.getElementById('qe-add-topic').value='';
  document.getElementById('qe-add-diff').value='mittel';
  document.getElementById('qe-add-numopts').value='4';
  qeAddUpdateOpts();
  document.getElementById('qe-add-modal').style.display='flex';
  setTimeout(()=>document.getElementById('qe-add-text').focus(),50)
}
function qeAddClose(){document.getElementById('qe-add-modal').style.display='none'}
function qeAddUpdateOpts(){
  const n=parseInt((document.getElementById('qe-add-numopts')||{}).value)||4;
  let h='';
  for(let i=0;i<n;i++){
    h+='<div class="qe-opt-edit" id="qe-add-orow-'+i+'">'
      +'<input type="radio" name="qe-add-cor" value="'+i+'"'+(i===0?' checked':'')+' onchange="qeAddUpdateCor()">'
      +'<span class="qe-ol">'+OLQ[i]+'</span>'
      +'<input type="text" id="qe-add-opt'+i+'" style="flex:1;background:none;border:none;color:var(--text);font-family:Outfit,sans-serif;font-size:.83rem;outline:none;min-width:0" placeholder="Option '+OLQ[i]+'\u2026"/>'
      +'</div>';
  }
  const el=document.getElementById('qe-add-opts');if(el)el.innerHTML=h;
  qeAddUpdateCor()
}
function qeAddUpdateCor(){
  const n=parseInt((document.getElementById('qe-add-numopts')||{}).value)||4;
  document.querySelectorAll('[name="qe-add-cor"]').forEach((r,i)=>{
    const row=document.getElementById('qe-add-orow-'+i);
    if(row)row.className='qe-opt-edit'+(r.checked?' sel-cor':'')
  })
}
async function qeAddSubmit(){
  const gv=id=>(document.getElementById(id)||{}).value||'';
  const text=(document.getElementById('qe-add-text')||{}).value||'';
  if(!text.trim()){showToast('Fragetext erforderlich',1);return}
  const topic=gv('qe-add-topic').trim();
  const difficulty=gv('qe-add-diff');
  const numOptions=parseInt(gv('qe-add-numopts'))||4;
  const options=[];let correctAnswer=0;
  for(let i=0;i<numOptions;i++){
    const v=gv('qe-add-opt'+i).trim();
    if(!v){showToast('Alle Optionen ausf\u00fcllen',1);return}
    options.push(v)
  }
  document.querySelectorAll('[name="qe-add-cor"]').forEach((r,i)=>{if(r.checked)correctAnswer=i});
  try{
    const r=await apiFetch('/api/admin/questions',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({text,options,correctAnswer,numOptions,difficulty,topic})});
    const d=await r.json();
    if(d.ok){showToast('Frage hinzugef\u00fcgt',0);qeAddClose();loadQuestions(1)}
    else showToast(d.error||'Fehler',1)
  }catch(e){showToast('Netzwerkfehler',1)}
}

// ── Duplikat-Erkennung ───────────────────────────────────────────────────────
let qeDupeSameAnswerFilter=false;
function qeDupeSameAnswer(p){
  const ac=(p.a.options||[])[p.a.correctAnswer];
  const bc=(p.b.options||[])[p.b.correctAnswer];
  return ac!=null&&bc!=null&&ac===bc
}
function qeDupeAutoThrUpdate(){
  const thr=Math.min(100,Math.max(60,parseInt((document.getElementById('qe-dupe-auto-thr')||{}).value)||95))/100;
  const pairs=qeDupePairs.filter((p,i)=>!qeDupeSkipped.has(i)&&p.score>=thr);
  const countEl=document.getElementById('qe-dupe-auto-count');
  if(countEl)countEl.textContent=pairs.length;
  const simBtn=document.getElementById('qe-dupe-sim-btn');
  if(simBtn)simBtn.style.display=pairs.length?'':'none'
}
function qeDupeResetFilterUI(){
  qeDupeSameAnswerFilter=false;
  const btn=document.getElementById('qe-dupe-sa-btn');
  if(btn){btn.className='btn btn-sm btn-s';btn.textContent='\u25ba Gleiche richtige Antwort'}
  const del=document.getElementById('qe-dupe-del-btn');
  if(del)del.style.display='none'
}
function qeDupeToggleSameAnswer(){
  qeDupeSameAnswerFilter=!qeDupeSameAnswerFilter;
  const btn=document.getElementById('qe-dupe-sa-btn');
  if(btn){btn.className='btn btn-sm '+(qeDupeSameAnswerFilter?'btn-p':'btn-s');btn.textContent=(qeDupeSameAnswerFilter?'\u25bc':'\u25ba')+' Gleiche richtige Antwort'}
  renderDupeList()
}
async function qeDeleteDupeRandom(){
  const pairs=qeDupePairs.filter((p,i)=>!qeDupeSkipped.has(i)&&qeDupeSameAnswer(p));
  if(!pairs.length){showToast('Keine Paare zum L\u00f6schen',1);return}
  if(!confirm(pairs.length+' Paar'+(pairs.length!==1?'e':'')+': Jeweils eine zuf\u00e4llige Frage l\u00f6schen?'))return;
  const toDelete=[...new Set(pairs.map(p=>Math.random()<0.5?p.a.id:p.b.id))];
  let deleted=0,failed=0;
  for(const id of toDelete){
    try{
      const r=await apiFetch('/api/admin/questions',{method:'DELETE',headers:{'Content-Type':'application/json'},body:JSON.stringify({id})});
      const d=await r.json();
      if(d.ok)deleted++;else failed++
    }catch(e){failed++}
  }
  showToast(deleted+' Fragen gel\u00f6scht'+(failed?' ('+failed+' Fehler)':''),failed?1:0);
  loadDuplicates()
}
async function qeDeleteDupeSimilar(){
  const thr=Math.min(100,Math.max(60,parseInt((document.getElementById('qe-dupe-auto-thr')||{}).value)||95))/100;
  const pairs=qeDupePairs.filter((p,i)=>!qeDupeSkipped.has(i)&&p.score>=thr);
  if(!pairs.length){showToast('Keine Paare in diesem Bereich',1);return}
  const pct=Math.round(thr*100);
  if(!confirm(pairs.length+' Paar'+(pairs.length!==1?'e':'')+' mit \u2265\u202f'+pct+'\u202f% \u00c4hnlichkeit: Jeweils eine zuf\u00e4llige Frage l\u00f6schen?'))return;
  const toDelete=[...new Set(pairs.map(p=>Math.random()<0.5?p.a.id:p.b.id))];
  let deleted=0,failed=0;
  for(const id of toDelete){
    try{
      const r=await apiFetch('/api/admin/questions',{method:'DELETE',headers:{'Content-Type':'application/json'},body:JSON.stringify({id})});
      const d=await r.json();
      if(d.ok)deleted++;else failed++
    }catch(e){failed++}
  }
  showToast(deleted+' Fragen gel\u00f6scht'+(failed?' ('+failed+' Fehler)':''),failed?1:0);
  loadDuplicates()
}
async function loadDuplicates(){
  const listEl=document.getElementById('qe-dupes-list'),metaEl=document.getElementById('qe-dupes-meta');
  if(listEl)listEl.innerHTML='<div class="log-empty">Suche Duplikate\u2026</div>';
  if(metaEl)metaEl.textContent='';
  qeDupeResetFilterUI();
  const thr=parseFloat((document.getElementById('qe-dupe-threshold')||{}).value||'0.85').toFixed(2);
  try{
    const r=await fetch('/api/admin/questions?view=duplicates&threshold='+thr);
    if(!r.ok){if(listEl)listEl.innerHTML='<div class="log-empty" style="color:var(--wrong)">Fehler beim Laden.</div>';return}
    const d=await r.json();
    // Deduplicate by normalised pair key (smaller ID first) – safety net against backend edge cases
    const seenKeys=new Set();
    qeDupePairs=(d.pairs||[]).filter(p=>{
      const key=Math.min(p.a.id,p.b.id)+'-'+Math.max(p.a.id,p.b.id);
      if(seenKeys.has(key))return false;
      seenKeys.add(key);return true
    });
    qeDupeSkipped=new Set();
    renderDupeList();
    qeDupeAutoThrUpdate()
  }catch(e){if(listEl)listEl.innerHTML='<div class="log-empty" style="color:var(--wrong)">Netzwerkfehler.</div>'}
}
function renderDupeCardContent(q){
  const diffCls={leicht:'leicht',mittel:'mittel',schwer:'schwer',extrem:'extrem'}[q.difficulty]||'';
  const diffTag=q.difficulty?'<span class="qe-tag '+diffCls+'">'+(DL[q.difficulty]||q.difficulty)+'</span>':'';
  const topicTag=q.topic?'<span class="qe-tag topic">'+esc(q.topic)+'</span>':'';
  const opts=(q.options||[]).map((o,i)=>'<div class="qe-opt'+(i===q.correctAnswer?' cor':'')+'">'+'<span class="qe-ol">'+OLQ[i]+'</span><span>'+esc(o)+'</span></div>').join('');
  return '<div style="display:flex;gap:5px;flex-wrap:wrap;margin-bottom:7px">'+diffTag+topicTag+'</div>'
    +'<div class="qe-question" style="font-size:.88rem">'+esc(q.text)+'</div>'
    +'<div class="qe-opts" style="margin-top:6px">'+opts+'</div>'
}
function renderDupeList(){
  const listEl=document.getElementById('qe-dupes-list'),metaEl=document.getElementById('qe-dupes-meta');
  if(!listEl)return;
  let visible=qeDupePairs.filter((p,i)=>!qeDupeSkipped.has(i));
  if(qeDupeSameAnswerFilter)visible=visible.filter(qeDupeSameAnswer);
  const delBtn=document.getElementById('qe-dupe-del-btn');
  if(delBtn)delBtn.style.display=(qeDupeSameAnswerFilter&&visible.length)?'':'none';
  qeDupeAutoThrUpdate();
  if(metaEl){
    const sk=qeDupeSkipped.size;
    const fNote=qeDupeSameAnswerFilter?' \u2022 Filter aktiv: gleiche richtige Antwort':'';
    metaEl.textContent=visible.length+' Paar'+(visible.length!==1?'e':'')+(sk?' ('+sk+' \u00fcbersprungen)':'')+fNote
  }
  if(!visible.length){listEl.innerHTML='<div class="log-empty">'+(qeDupeSameAnswerFilter?'Keine Paare mit gleicher richtiger Antwort gefunden.':'Keine Duplikate gefunden.')+'</div>';return}
  listEl.innerHTML=qeDupePairs.map((p,i)=>{
    if(qeDupeSkipped.has(i))return '';
    if(qeDupeSameAnswerFilter&&!qeDupeSameAnswer(p))return '';
    const pct=Math.round(p.score*100);
    const dupeLabel=p.score>=0.95?'Exaktes Duplikat':p.score>=0.85?'Sehr \u00e4hnlich':'M\u00f6glicherweise Duplikat';
    return '<div class="qe-dupe-pair" id="qe-dupe-'+i+'">'
      +'<span class="qe-dupe-score">'+pct+'% \u2014 '+dupeLabel+'</span>'
      +'<div class="qe-dupe-cols">'
      +'<div class="qe-dupe-col"><div class="qe-dupe-col-label">Frage A &middot; ID '+p.a.id+'</div>'+renderDupeCardContent(p.a)+'</div>'
      +'<div class="qe-dupe-col"><div class="qe-dupe-col-label">Frage B &middot; ID '+p.b.id+'</div>'+renderDupeCardContent(p.b)+'</div>'
      +'</div>'
      +'<div class="qe-dupe-actions">'
      +'<button class="ib" onclick="qeDupeAction('+i+',\'keepA\')">\u2714 A behalten</button>'
      +'<button class="ib" onclick="qeDupeAction('+i+',\'keepB\')">\u2714 B behalten</button>'
      +'<button class="ib danger" onclick="qeDupeAction('+i+',\'both\')">\u2716 Beide l\u00f6schen</button>'
      +'<button class="ib" style="color:var(--text2)" onclick="qeDupeAction('+i+',\'skip\')">\u2192 \u00dcberspringen</button>'
      +'</div></div>'
  }).join('')
}
async function markPairReviewed(idA,idB){
  try{await apiFetch('/api/admin/questions',{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify({markPairA:idA,markPairB:idB})})}catch(e){}
}
async function qeDupeAction(pairIdx,action){
  const pair=qeDupePairs[pairIdx];if(!pair)return;
  if(action==='skip'){
    await markPairReviewed(pair.a.id,pair.b.id);
    qeDupeSkipped.add(pairIdx);renderDupeList();return
  }
  const toDelete=[];
  if(action==='keepA')toDelete.push(pair.b.id);
  else if(action==='keepB')toDelete.push(pair.a.id);
  else if(action==='both'){toDelete.push(pair.a.id);toDelete.push(pair.b.id)}
  for(const id of toDelete){
    try{
      const r=await apiFetch('/api/admin/questions',{method:'DELETE',headers:{'Content-Type':'application/json'},body:JSON.stringify({id})});
      const d=await r.json();
      if(!d.ok){showToast(d.error||'Fehler beim L\u00f6schen',1);return}
    }catch(e){showToast('Netzwerkfehler',1);return}
  }
  if(action!=='both')await markPairReviewed(pair.a.id,pair.b.id);
  const deletedIds=new Set(toDelete);
  qeDupePairs=qeDupePairs.filter((p,i)=>{
    if(i===pairIdx)return false;
    if(deletedIds.has(p.a.id)||deletedIds.has(p.b.id))return false;
    return true;
  });
  qeDupeSkipped=new Set();
  showToast(toDelete.length===2?'Beide gel\u00f6scht':'Gel\u00f6scht',0);
  renderDupeList()
}

async function reportQuestion(){
  if(qReportedCurrent)return;
  const q=gameState&&gameState.question;if(!q)return;
  qReportedCurrent=true;
  const btn=document.querySelector('#q-report-area .q-report-btn');if(btn)btn.disabled=true;
  try{await apiFetch('/api/questions/report',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({text:q.text})})}catch(e){}
  showToast('Danke f\u00fcr den Hinweis!',0)
}
async function loadReportedQuestions(){
  const listEl=document.getElementById('qe-reported-list'),metaEl=document.getElementById('qe-reported-meta');
  if(listEl)listEl.innerHTML='<div class="log-empty">Lade\u2026</div>';
  try{
    const r=await fetch('/api/admin/questions?view=reported');
    if(!r.ok){if(listEl)listEl.innerHTML='<div class="log-empty" style="color:var(--wrong)">Fehler beim Laden.</div>';return}
    const d=await r.json();
    const qs=d.questions||[];
    if(metaEl)metaEl.textContent=qs.length+' gemeldete Frage'+(qs.length!==1?'n':'');
    if(!qs.length){if(listEl)listEl.innerHTML='<div class="log-empty">Keine gemeldeten Fragen.</div>';return}
    if(listEl)listEl.innerHTML='<div class="qe-grid">'+qs.map(q=>renderReportedCard(q)).join('')+'</div>'
  }catch(e){if(listEl)listEl.innerHTML='<div class="log-empty" style="color:var(--wrong)">Fehler beim Laden.</div>'}
}
function renderReportedCard(q){
  const dt=new Date(q.askedAt*1000).toLocaleString('de-DE');
  const diffCls={leicht:'leicht',mittel:'mittel',schwer:'schwer',extrem:'extrem'}[q.difficulty]||'';
  const diffTag=q.difficulty?'<span class="qe-tag '+diffCls+'">'+(DL[q.difficulty]||q.difficulty)+'</span>':'';
  const topicTag=q.topic?'<span class="qe-tag topic">'+esc(q.topic)+'</span>':'';
  const reportBadge='<span class="qe-tag" style="background:rgba(255,51,102,.18);color:var(--wrong)">\u26a0 '+q.reportCount+' Meldung'+(q.reportCount!==1?'en':'')+'</span>';
  const opts=(q.options||[]).map((o,i)=>'<div class="qe-opt'+(i===q.correctAnswer?' cor':'')+'">'+'<span class="qe-ol">'+OLQ[i]+'</span><span>'+esc(o)+'</span></div>').join('');
  qeData[q.id]=Object.assign({},q,{_context:'reported'});
  return '<div class="qe-card" data-qid="'+q.id+'">'
    +'<div style="display:flex;justify-content:space-between;align-items:flex-start;gap:8px;margin-bottom:8px">'
    +'<div style="display:flex;gap:5px;flex-wrap:wrap">'+diffTag+topicTag+reportBadge+'</div>'
    +'<div style="display:flex;gap:5px;flex-shrink:0">'
    +'<button class="ib" style="font-size:.65rem;padding:2px 9px" onclick="qeEdit('+q.id+')" title="Bearbeiten">&#9998;</button>'
    +'<button class="ib danger" style="font-size:.65rem;padding:2px 9px" onclick="qeDeleteOne('+q.id+')" title="L\u00f6schen">&#10005;</button>'
    +'</div></div>'
    +'<div class="qe-question">'+esc(q.text)+'</div>'
    +'<div class="qe-opts">'+opts+'</div>'
    +'<div class="qe-date" style="display:flex;justify-content:space-between;align-items:center">'
    +'<span>&#128336; '+esc(dt)+' \u00b7 ID\u202f'+q.id+'</span>'
    +'<button class="ib" style="font-size:.65rem;padding:2px 9px" onclick="qeClearReport('+q.id+')">Meldung zur\u00fccksetzen</button>'
    +'</div></div>'
}
async function qeClearReport(id){
  try{
    const r=await apiFetch('/api/admin/questions',{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify({id,clearReport:true})});
    const d=await r.json();
    if(!d.ok){showToast(d.error||'Fehler',1);return}
    showToast('Meldung zur\u00fcckgesetzt',0);loadReportedQuestions()
  }catch(e){showToast('Netzwerkfehler',1)}
}

// ── Audio (iOS-compatible) ──
// iOS requires user gesture to unlock AudioContext and <audio> elements.
// Strategy:
//   - On first user gesture, we unlock the bg-audio element + pool with a silent WAV
//   - We also create a tiny AudioContext to unlock the Web Audio API
//   - The real background music src is set ONLY in startBg() and never replaced
//     unless the actual track URL changed, to preserve iOS unlock state
let audioUnlocked=false;
const audioPool=[];const audioPoolSize=6;
const SILENT_WAV='data:audio/wav;base64,UklGRiQAAABXQVZFZm10IBAAAAABAAEAQB8AAIA+AAACABAAZGF0YQAAAAA=';
function initAudioPool(){for(let i=0;i<audioPoolSize;i++){const a=new Audio();a.preload='auto';audioPool.push(a)}}
initAudioPool();

function unlockAudio(){
  if(audioUnlocked)return;audioUnlocked=true;
  // Unlock the bg-audio element by briefly playing a silent WAV on it
  const bg=document.getElementById('bg-audio');
  const prevSrc=bg.src;
  bg.src=SILENT_WAV;bg.load();
  const p=bg.play();
  if(p&&p.then){p.then(()=>{bg.pause();bg.currentTime=0;if(prevSrc)bg.src=prevSrc}).catch(()=>{})}
  // Unlock gen-audio too (used during question generation)
  const ga=document.getElementById('gen-audio');
  if(ga){const prev2=ga.src;ga.src=SILENT_WAV;ga.load();const gp=ga.play();if(gp&&gp.then){gp.then(()=>{ga.pause();ga.currentTime=0;if(prev2)ga.src=prev2}).catch(()=>{})}}
  // Unlock pool elements
  audioPool.forEach(a=>{a.src=SILENT_WAV;const pp=a.play();if(pp&&pp.then)pp.then(()=>{a.pause();a.currentTime=0}).catch(()=>{})});
  // Unlock Web Audio API
  try{const ctx=new (window.AudioContext||window.webkitAudioContext)();const buf=ctx.createBuffer(1,1,22050);const src=ctx.createBufferSource();src.buffer=buf;src.connect(ctx.destination);src.start(0)}catch(e){}
}
// Register unlock on any user gesture (must be non-passive to fire before other handlers if needed)
['touchstart','touchend','click','keydown'].forEach(evt=>{document.addEventListener(evt,unlockAudio,{capture:true,passive:true})});

let sfxIdx=0;
function playSound(url,vol){
  if(!url)return;
  const a=audioPool[sfxIdx%audioPoolSize];sfxIdx++;
  a.src=url;a.volume=vol||0.6;a.currentTime=0;
  const p=a.play();if(p&&p.catch)p.catch(()=>{});
}

// Absolute URL resolver so we can compare properly
function absURL(u){const a=document.createElement('a');a.href=u;return a.href}
let bgCurrentSrc='';

function tryPlayBg(retries){
  const a=document.getElementById('bg-audio');
  const p=a.play();
  if(p&&p.catch){p.catch(()=>{
    if(retries>0){setTimeout(()=>tryPlayBg(retries-1),200)}
  })}
}

function startBg(){
  if(!globalSounds.backgroundSound||bgStarted)return;
  const a=document.getElementById('bg-audio');
  const targetAbs=absURL(globalSounds.backgroundSound);
  // Only change src if the TRACK actually changed, to preserve iOS unlock
  if(bgCurrentSrc!==targetAbs){
    a.src=globalSounds.backgroundSound;
    bgCurrentSrc=targetAbs;
    a.load();
  }
  const v=globalSounds.volBg||0.2;a.volume=v;
  const sl=document.getElementById('bg-vol-slider');if(sl)sl.value=v;
  a.loop=true;bgStarted=true;
  if(!bgMuted){
    // Wait until the element is ready to play, then try
    if(a.readyState>=2){tryPlayBg(3)}
    else{
      const onReady=()=>{a.removeEventListener('canplay',onReady);tryPlayBg(3)};
      a.addEventListener('canplay',onReady);
      // Also attempt immediately — may succeed on some browsers
      tryPlayBg(3);
    }
  }
}
function stopBg(){const a=document.getElementById('bg-audio');a.pause();a.currentTime=0;bgStarted=false}
let genStarted=false,genCurrentSrc='';
function startGenSound(){
  if(!globalSounds.generatingSound)return;
  const a=document.getElementById('gen-audio');if(!a)return;
  const target=absURL(globalSounds.generatingSound);
  if(genCurrentSrc!==target){a.src=globalSounds.generatingSound;genCurrentSrc=target;a.load()}
  a.volume=globalSounds.volGenerating||0.4;a.loop=true;genStarted=true;
  if(bgMuted)return;
  const tryPlay=(n)=>{const p=a.play();if(p&&p.catch)p.catch(()=>{if(n>0)setTimeout(()=>tryPlay(n-1),200)})};
  if(a.readyState>=2)tryPlay(3);else{const onReady=()=>{a.removeEventListener('canplay',onReady);tryPlay(3)};a.addEventListener('canplay',onReady);tryPlay(3)}
}
function stopGenSound(){const a=document.getElementById('gen-audio');if(!a)return;try{a.pause();a.currentTime=0}catch(e){}genStarted=false}
function toggleMute(){bgMuted=!bgMuted;const a=document.getElementById('bg-audio'),g=document.getElementById('gen-audio'),b=document.getElementById('mute-btn');if(bgMuted){a.pause();if(g)g.pause();b.innerHTML='&#128264;';b.classList.add('muted')}else{if(bgStarted)tryPlayBg(3);if(g&&genStarted){const p=g.play();if(p&&p.catch)p.catch(()=>{})}b.innerHTML='&#128266;';b.classList.remove('muted')}}
function adjustBgVol(v){const a=document.getElementById('bg-audio');a.volume=parseFloat(v);if(bgMuted&&parseFloat(v)>0){bgMuted=false;const b=document.getElementById('mute-btn');b.innerHTML='&#128266;';b.classList.remove('muted');if(bgStarted)tryPlayBg(3)}}

// WebSocket
let wsKeepalive=null;
function connectWS(){if(ws&&ws.readyState<=1)return;const proto=location.protocol==='https:'?'wss:':'ws:';ws=new WebSocket(proto+'//'+location.host+'/ws');
  ws.onopen=()=>{if(wsKeepalive)clearInterval(wsKeepalive);wsKeepalive=setInterval(()=>{if(ws&&ws.readyState===1)ws.send(JSON.stringify({type:'ping',payload:null}))},20000)};
  ws.onclose=()=>{if(wsKeepalive){clearInterval(wsKeepalive);wsKeepalive=null}
    // Auto-reconnect if we have session info
    const ri=sessionStorage.getItem('h_invite'),rn=sessionStorage.getItem('h_name');
    if(ri&&rn){setTimeout(()=>{connectWS();const ck=setInterval(()=>{if(ws&&ws.readyState===1){clearInterval(ck);send('reconnect',{Name:rn,InviteCode:ri})}},100)},2000)}};
  ws.onerror=()=>{};ws.onmessage=e=>{try{handleMessage(JSON.parse(e.data))}catch(ex){}}}
function send(t,p){if(ws&&ws.readyState===1)ws.send(JSON.stringify({type:t,payload:p}))}
function handleMessage(msg){switch(msg.type){
  case 'joined':case 'reconnected':myId=msg.payload.playerId;inviteCode=msg.payload.inviteCode;sessionStorage.setItem('h_invite',inviteCode);break;
  case 'state':gameState=msg.payload;if(!shuffleAnimating)renderGame();break;case 'error':showToast(msg.payload.message,1);break;
  case 'kicked':showToast(msg.payload.message,1);gameState=null;stopBg();stopGenSound();sessionStorage.removeItem('h_invite');sessionStorage.removeItem('h_name');showScreen('home');break;
  case 'shuffle_animation':shuffleAnimation(msg.payload);break}}

function joinLobby(code,mode){joinPending=code;joinNeedsPw=mode==='pw';document.getElementById('join-pw-group').style.display=joinNeedsPw?'block':'none';showScreen('join')}
function joinByCode(){const code=document.getElementById('home-invite').value.trim();if(!code){showToast('Code eingeben',1);return};joinPending=code;joinNeedsPw=false;document.getElementById('join-pw-group').style.display='none';showScreen('join')}
function createGame(){const name=document.getElementById('host-name').value.trim();if(!name){showToast('Namen eingeben',1);return}
  if(!currentUser){showToast('Bitte zuerst anmelden',1);return}
  sessionStorage.setItem('h_name',name);connectWS();const check=setInterval(()=>{if(ws&&ws.readyState===1){clearInterval(check);send('create_game',{Name:name,Username:currentUser?currentUser.username:''});
    setTimeout(()=>send('update_settings',{Topic:document.getElementById('create-topic').value||'Allgemeinwissen',Difficulty:document.getElementById('create-diff').value,StartDifficulty:(document.getElementById('create-startdiff')||{}).value||'leicht',Mode:document.getElementById('create-mode').value,NumQuestions:parseInt(document.getElementById('create-questions').value)||10,TimePerQ:parseInt(document.getElementById('create-time').value)||20,NumOptions:parseInt(document.getElementById('create-options').value)||4,NumTeeth:parseInt(document.getElementById('create-teeth').value)||5,ShowTutorial:document.getElementById('create-tutorial').checked,WebSearch:document.getElementById('create-websearch').checked,PlayIntro:document.getElementById('create-playintro').checked,LobbyName:document.getElementById('create-lobbyname').value.trim()||'',LobbyMode:document.getElementById('create-lobbymode').value,LobbyPassword:document.getElementById('create-lobbypw').value||''}),200)}},100)}
function doJoin(){const name=document.getElementById('join-name').value.trim();if(!name){showToast('Namen eingeben',1);return}
  sessionStorage.setItem('h_name',name);const pw=joinNeedsPw?document.getElementById('join-pw').value:'';connectWS();const check=setInterval(()=>{if(ws&&ws.readyState===1){clearInterval(check);send('join_game',{Name:name,InviteCode:joinPending,Password:pw})}},100)}
function startGame(){send('start_game',{})}
function skipTutorial(){send('skip_tutorial',{})}
function endGameEarly(){if(confirm('Spiel wirklich beenden?'))send('end_game',{})}
function canChangeAnswer(){if(!gameState||!gameState.settings||!gameState.settings.allowAnswerChange)return false;return gameState.players.some(p=>p.alive&&p.connected&&!p.answered)}
function submitAnswer(i){if(selectedAnswer===i)return;if(selectedAnswer>=0&&!canChangeAnswer())return;selectedAnswer=i;playSound(globalSounds.answerSound,globalSounds.volAnswer);send('answer',{answer:i});renderOptions()}
function kickPlayer(pid){if(confirm('Entfernen?'))send('kick_player',{playerId:pid})}
function transferHost(pid){if(confirm('Host-Rechte übertragen?'))send('transfer_host',{playerId:pid})}
function playAgain(){send('play_again',{})}
function leaveGame(){send('leave_game',{})}
function copyInvite(){navigator.clipboard.writeText(location.origin+'/?join='+inviteCode).then(()=>showToast('Kopiert!',0))}
function sendLS(){
  const isDelegated=gameState&&gameState.delegatedTo===myId&&myId!==gameState.hostId;
  const s={};
  const allIds={Topic:'l-topic',Difficulty:'l-diff',StartDifficulty:'l-startdiff',Mode:'l-mode',NumQuestions:'l-questions',TimePerQ:'l-time',NumOptions:'l-options',NumTeeth:'l-teeth',LobbyName:'l-lobbyname',LobbyMode:'l-lobbymode',LobbyPassword:'l-lobbypw'};
  const delegateIds={Topic:'l-topic',Difficulty:'l-diff',StartDifficulty:'l-startdiff'};
  const ids=isDelegated?delegateIds:allIds;
  for(const[k,id] of Object.entries(ids)){const el=document.getElementById(id);if(!el)continue;if(['NumQuestions','TimePerQ','NumOptions','NumTeeth'].includes(k))s[k]=parseInt(el.value)||0;else s[k]=el.value}
  if(!isDelegated){const tc=document.getElementById('l-tutorial');if(tc)s.ShowTutorial=tc.checked;
  const ws2=document.getElementById('l-websearch');if(ws2)s.WebSearch=ws2.checked;
  const pi=document.getElementById('l-playintro');if(pi)s.PlayIntro=pi.checked;
  const ac=document.getElementById('l-allowchange');if(ac)s.AllowAnswerChange=ac.checked;}
  send('update_settings',s)}

function amIAlive(){if(!gameState)return true;const me=gameState.players.find(p=>p.id===myId);return me?me.alive:true}

function delegateSettings(pid){send('delegate_settings',{playerId:pid})}

function shuffleAnimation(data){
  shuffleAnimating=true;
  const players=data.players||[];const winnerId=data.winnerId;
  const names=players.map(p=>p.name);if(!names.length){shuffleAnimating=false;renderGame();return}
  const overlay=document.createElement('div');
  overlay.id='shuffle-overlay';
  overlay.style.cssText='position:fixed;inset:0;background:rgba(10,10,15,.96);z-index:9999;display:flex;flex-direction:column;align-items:center;justify-content:center;gap:24px;animation:so-in .3s ease';
  overlay.innerHTML='<style>@keyframes so-in{from{opacity:0}to{opacity:1}}@keyframes so-glow{0%,100%{text-shadow:0 0 20px rgba(255,215,0,.6)}50%{text-shadow:0 0 40px rgba(255,215,0,1),0 0 80px rgba(255,215,0,.4)}}</style>'
    +'<div style="font-size:.75rem;color:var(--text2);letter-spacing:3px;text-transform:uppercase">Auswahl wird delegiert...</div>'
    +'<div id="shuffle-name" style="font-family:\'Space Mono\',monospace;font-size:2.6rem;font-weight:700;color:var(--accent);text-align:center;min-height:3.5rem;display:flex;align-items:center;justify-content:center;transition:opacity .08s;max-width:90vw;word-break:break-word"></div>'
    +'<div id="shuffle-sub" style="font-size:.85rem;color:var(--text2);display:none">darf Thema &amp; Schwierigkeit wählen</div>';
  document.body.appendChild(overlay);
  const nameEl=document.getElementById('shuffle-name');
  const subEl=document.getElementById('shuffle-sub');
  const winnerName=(players.find(p=>p.id===winnerId)||{}).name||'';
  const totalSteps=28;let step=0;
  function nextStep(){
    if(step<totalSteps){
      nameEl.textContent=names[Math.floor(Math.random()*names.length)];
      const delay=40+Math.pow(step/totalSteps,2)*320;
      step++;setTimeout(nextStep,delay);
    } else {
      nameEl.textContent=winnerName;
      nameEl.style.color='var(--gold)';
      nameEl.style.animation='so-glow 1.5s ease-in-out infinite';
      subEl.style.display='block';
      setTimeout(()=>{
        overlay.style.transition='opacity .5s';overlay.style.opacity='0';
        setTimeout(()=>{overlay.remove();shuffleAnimating=false;renderGame()},500);
      },2500);
    }
  }
  nextStep();
}

let lastPhase='';
function renderGame(){if(!gameState)return;const ph=gameState.phase,isHost=myId===gameState.hostId;
  if(ph!=='loading'&&ph!=='refill'&&genStarted)stopGenSound();
  if(ph==='lobby'){showScreen('lobby');renderLobby(isHost)}
  else if(ph==='loading'){showScreen('loading');startGenSound()}else if(ph==='refill'){showScreen('refill');startGenSound()}
  else if(ph==='tutorial'){showScreen('tutorial');renderTutorial(isHost)}
  else if(ph==='intro'){showScreen('intro');playSound(globalSounds.introSound,globalSounds.volIntro);startBg()}
  else if(ph==='question'){showScreen('game');if(lastPhase!=='question'){selectedAnswer=-1;playSound(globalSounds.questionSound,globalSounds.volQuestion);renderQuestion();}renderGamePlayers(isHost);startTimer()}
  else if(ph==='results'){showScreen('game');renderResults();renderGamePlayers(isHost);clearTimer();
    if(gameState.allWrong&&gameState.players.filter(p=>p.alive&&p.connected).length>=2)playSound(globalSounds.allwrongSound||globalSounds.wrongSound,globalSounds.allwrongSound?globalSounds.volAllwrong:globalSounds.volWrong);
    else if(gameState.allCorrect)playSound(globalSounds.allcorrectSound,globalSounds.volAllcorrect);
    else if(gameState.someoneLost)playSound(globalSounds.wrongSound,globalSounds.volWrong)}
  else if(ph==='end'){clearTimer();showScreen('end');renderEnd(isHost);stopBg()}
  else if(ph==='error'){clearTimer();stopBg();showScreen('error');document.getElementById('error-msg').textContent=gameState.errorMsg||'Unbekannter Fehler'}
  document.getElementById('host-controls').style.display=(isHost&&(ph==='question'||ph==='results'))?'block':'none';lastPhase=ph}

function renderTutorial(isHost){document.getElementById('tutorial-content').innerHTML=tutorialHtml||'<p>Lade Tutorial...</p>';document.getElementById('tutorial-actions').innerHTML=isHost?'<button class="btn btn-p" onclick="skipTutorial()" style="max-width:280px">Weiter</button>':'<div style="color:var(--text2);font-size:.85rem">Warte auf den*die Host...</div>'}

function renderLobby(isHost){
  const invBox=document.querySelector('.gcd');if(invBox)invBox.style.display=isHost?'':'none';
  if(isHost){const link=location.origin+'/?join='+inviteCode;document.getElementById('invite-link').textContent=link;
  const qrEl=document.getElementById('qr-code');qrEl.innerHTML='';try{new QRCode(qrEl,{text:link,width:180,height:180,colorDark:'#e4e4ef',colorLight:'#0a0a0f'})}catch(e){}}
  const s=gameState.settings,players=gameState.players||[],spLocked=players.filter(p=>p.connected).length>1;
  const delegatedTo=gameState.delegatedTo||'';
  const delegatedPlayer=delegatedTo?players.find(p=>p.id===delegatedTo):null;
  const amIDelegated=delegatedTo===myId&&!isHost;
  if(isHost){
    const dOpts=['leicht','mittel','schwer','extrem'].map(d=>'<option value="'+d+'"'+(s.difficulty===d?' selected':'')+'>'+DL[d]+'</option>').join('');
    const sdOpts=['leicht','mittel','schwer','extrem'].map(d=>'<option value="'+d+'"'+((s.startDifficulty||'leicht')===d?' selected':'')+'>'+DL[d]+'</option>').join('');
    const mOpts='<option value="classic"'+(s.mode==='classic'?' selected':'')+'>Klassisch</option><option value="elimination"'+(s.mode==='elimination'?' selected':'')+'>Elimination</option><option value="kfo_battle_royale"'+(s.mode==='kfo_battle_royale'?' selected':'')+'>KFO Battle Royale</option><option value="kfo_singleplayer"'+(s.mode==='kfo_singleplayer'?' selected':'')+(spLocked?' disabled':'')+'>KFO Singleplayer'+(spLocked?' (gesperrt)':'')+'</option>';
    const lmOpts='<option value="invite"'+(s.lobbyMode==='invite'?' selected':'')+'>Nur Einladung</option><option value="password"'+(s.lobbyMode==='password'?' selected':'')+'>Mit Passwort</option><option value="open"'+(s.lobbyMode==='open'?' selected':'')+'>Offen</option>';
    const isKFO=s.mode==='kfo_battle_royale'||s.mode==='kfo_singleplayer',isEndless=s.mode!=='classic';
    const oOpts=[2,3,4].map(n=>'<option value="'+n+'"'+((s.numOptions||4)===n?' selected':'')+'>'+n+'</option>').join('');
    const topicDiffHtml=delegatedTo
      ?'<div class="ig" style="grid-column:1/-1"><div style="background:rgba(255,107,157,.1);border:1px solid rgba(255,107,157,.3);border-radius:8px;padding:8px 12px;font-size:.75rem;color:var(--accent2);margin-bottom:8px">&#9997; '+esc(delegatedPlayer?delegatedPlayer.name:'...')+' w\u00e4hlt Thema &amp; Schwierigkeit</div>'
       +'<div class="sv-grid" style="margin-top:4px"><div class="sv-item"><span class="sv-label">Thema</span><span class="sv-value">'+esc(s.topic)+'</span></div>'
       +(!isKFO?'<div class="sv-item"><span class="sv-label">Schwierigkeit</span><span class="sv-value">'+(DL[s.difficulty]||s.difficulty)+'</span></div>':'')
       +(isKFO?'<div class="sv-item"><span class="sv-label">Start-Schwierigkeit</span><span class="sv-value">'+(DL[s.startDifficulty||'leicht']||s.startDifficulty)+'</span></div>':'')
       +'</div></div>'
      :'<div class="ig"><label>Thema</label><input type="text" id="l-topic" value="'+esc(s.topic)+'" onchange="sendLS()"/></div>'
       +(!isKFO?'<div class="ig"><label>Schwierigkeit</label><select id="l-diff" onchange="sendLS()">'+dOpts+'</select></div>':'')
       +(isKFO?'<div class="ig"><label>Start-Schwierigkeit</label><select id="l-startdiff" onchange="sendLS()">'+sdOpts+'</select></div>':'');
    const otherPlayers=players.filter(p=>p.id!==gameState.hostId&&p.connected);
    const delegateHtml=delegatedTo
      ?'<div class="spt" style="margin-top:18px;margin-bottom:10px">Auswahl delegiert</div><div style="display:flex;align-items:center;gap:10px;flex-wrap:wrap"><span style="color:var(--accent2);font-size:.85rem">'+esc(delegatedPlayer?delegatedPlayer.name:'...')+' darf w\u00e4hlen</span><button class="btn btn-s btn-sm" onclick="delegateSettings(\'\')">Zur\u00fcknehmen</button></div>'
      :otherPlayers.length
        ?'<div class="spt" style="margin-top:18px;margin-bottom:10px">Auswahl delegieren</div><div style="display:flex;gap:8px;flex-wrap:wrap;align-items:center"><select id="l-delegate-pid" style="flex:1;min-width:120px;padding:9px 11px;background:var(--bg2);border:1px solid var(--border);border-radius:10px;color:var(--text);font-family:Outfit,sans-serif;font-size:.9rem;outline:none">'+otherPlayers.map(p=>'<option value="'+p.id+'">'+esc(p.name)+'</option>').join('')+'</select><button class="btn btn-s btn-sm" onclick="delegateSettings(document.getElementById(\'l-delegate-pid\').value)">Delegieren</button><button class="btn btn-s btn-sm" onclick="delegateSettings(\'random\')">\u{1F3B2} Zuf\u00e4llig</button></div>'
        :'';
    document.getElementById('lobby-settings').innerHTML='<div class="sp"><div class="spt">Einstellungen</div><div class="sg">'
      +'<div class="ig"><label>Lobby-Name</label><input type="text" id="l-lobbyname" value="'+esc(s.lobbyName||'')+'" onchange="sendLS()" maxlength="30"/></div>'
      +'<div class="ig"><label>Lobby-Modus</label><select id="l-lobbymode" onchange="sendLS()">'+lmOpts+'</select></div>'
      +(s.lobbyMode==='password'?'<div class="ig"><label>Lobby-Passwort</label><input type="password" id="l-lobbypw" value="'+esc(s.lobbyPassword||'')+'" onchange="sendLS()"/></div>':'')
      +topicDiffHtml
      +'<div class="ig"><label>Modus</label><select id="l-mode" onchange="sendLS()">'+mOpts+'</select></div>'
      +(!isEndless?'<div class="ig"><label>Fragen</label><input type="number" id="l-questions" value="'+s.numQuestions+'" min="1" max="50" onchange="sendLS()"/></div>':'')
      +'<div class="ig"><label>Zeit (Sek.)</label><input type="number" id="l-time" value="'+s.timePerQuestion+'" min="5" max="120" onchange="sendLS()"/></div>'
      +'<div class="ig"><label>Optionen</label><select id="l-options" onchange="sendLS()">'+oOpts+'</select></div>'
      +'<div class="ig"><label>Z\u00e4hne</label><input type="number" id="l-teeth" value="'+s.numTeeth+'" min="1" max="20" onchange="sendLS()"/></div>'
      +'</div><div class="ig" style="margin-top:8px"><label><input type="checkbox" id="l-tutorial"'+(s.showTutorial?' checked':'')+' onchange="sendLS()"/>Tutorial anzeigen</label></div><div class="ig" style="margin-top:4px"><label><input type="checkbox" id="l-playintro"'+(s.playIntro!==false?' checked':'')+' onchange="sendLS()"/>Intro-Sound abspielen</label></div><div class="ig" style="margin-top:4px"><label><input type="checkbox" id="l-websearch"'+(s.webSearch?' checked':'')+' onchange="sendLS()"/>Internet-Recherche f\u00fcr Fragen</label></div><div class="ig" style="margin-top:4px"><label><input type="checkbox" id="l-allowchange"'+(s.allowAnswerChange!==false?' checked':'')+' onchange="sendLS()"/>Antwort \u00e4nderbar, solange nicht alle geantwortet haben</label></div>'
      +delegateHtml+'</div>';
  }else if(amIDelegated){const isKFO2=s.mode==='kfo_battle_royale'||s.mode==='kfo_singleplayer';
    const dOpts=['leicht','mittel','schwer','extrem'].map(d=>'<option value="'+d+'"'+(s.difficulty===d?' selected':'')+'>'+DL[d]+'</option>').join('');
    const sdOpts=['leicht','mittel','schwer','extrem'].map(d=>'<option value="'+d+'"'+((s.startDifficulty||'leicht')===d?' selected':'')+'>'+DL[d]+'</option>').join('');
    document.getElementById('lobby-settings').innerHTML='<div class="sv">'
      +'<div style="background:rgba(255,107,157,.1);border:1px solid rgba(255,107,157,.3);border-radius:8px;padding:10px 14px;margin-bottom:14px;font-size:.8rem;color:var(--accent2);text-align:center;font-weight:600">&#9997; Du darfst Thema &amp; Schwierigkeit w\u00e4hlen</div>'
      +'<div class="sg">'
      +'<div class="ig"><label>Thema</label><input type="text" id="l-topic" value="'+esc(s.topic)+'" onchange="sendLS()"/></div>'
      +(!isKFO2?'<div class="ig"><label>Schwierigkeit</label><select id="l-diff" onchange="sendLS()">'+dOpts+'</select></div>':'')
      +(isKFO2?'<div class="ig"><label>Start-Schwierigkeit</label><select id="l-startdiff" onchange="sendLS()">'+sdOpts+'</select></div>':'')
      +'</div><div class="spt" style="margin-top:14px">Weitere Einstellungen</div><div class="sv-grid">'
      +'<div class="sv-item"><span class="sv-label">Lobby</span><span class="sv-value">'+esc(s.lobbyName||'-')+'</span></div>'
      +'<div class="sv-item"><span class="sv-label">Modus</span><span class="sv-value">'+(ML[s.mode]||'Klassisch')+'</span></div>'
      +(s.mode==='classic'?'<div class="sv-item"><span class="sv-label">Fragen</span><span class="sv-value">'+s.numQuestions+'</span></div>':'')
      +'<div class="sv-item"><span class="sv-label">Zeit</span><span class="sv-value">'+s.timePerQuestion+' Sek.</span></div>'
      +'<div class="sv-item"><span class="sv-label">Z\u00e4hne</span><span class="sv-value">'+s.numTeeth+'</span></div>'
      +'</div></div>';
  }else{const isKFO2=s.mode==='kfo_battle_royale'||s.mode==='kfo_singleplayer';
    const delegateNotice=delegatedTo&&delegatedPlayer?'<div style="background:rgba(255,107,157,.08);border:1px solid rgba(255,107,157,.25);border-radius:8px;padding:8px 12px;margin-bottom:12px;font-size:.75rem;color:var(--accent2)">&#9997; '+esc(delegatedPlayer.name)+' w\u00e4hlt Thema &amp; Schwierigkeit</div>':'';
    document.getElementById('lobby-settings').innerHTML='<div class="sv">'+delegateNotice+'<div class="spt">Spieleinstellungen</div><div class="sv-grid">'
      +'<div class="sv-item"><span class="sv-label">Lobby</span><span class="sv-value">'+esc(s.lobbyName||'-')+'</span></div>'
      +'<div class="sv-item"><span class="sv-label">Thema</span><span class="sv-value">'+esc(s.topic)+'</span></div>'
      +(!isKFO2?'<div class="sv-item"><span class="sv-label">Schwierigkeit</span><span class="sv-value">'+(DL[s.difficulty]||'Mittel')+'</span></div>':'')
      +(isKFO2?'<div class="sv-item"><span class="sv-label">Start-Schwierigkeit</span><span class="sv-value">'+(DL[s.startDifficulty||'leicht']||'Leicht')+'</span></div>':'')
      +'<div class="sv-item"><span class="sv-label">Modus</span><span class="sv-value">'+(ML[s.mode]||'Klassisch')+'</span></div>'
      +(s.mode==='classic'?'<div class="sv-item"><span class="sv-label">Fragen</span><span class="sv-value">'+s.numQuestions+'</span></div>':'')
      +'<div class="sv-item"><span class="sv-label">Zeit</span><span class="sv-value">'+s.timePerQuestion+' Sek.</span></div>'
      +'<div class="sv-item"><span class="sv-label">Z\u00e4hne</span><span class="sv-value">'+s.numTeeth+'</span></div>'
      +'</div></div>'}
  document.getElementById('players-count').textContent='Spieler*innen ('+players.length+')';
  document.getElementById('lobby-players').innerHTML=players.map(p=>{
    const isH=p.id===gameState.hostId,isMe=p.id===myId,isD=delegatedTo&&p.id===delegatedTo;
    return '<div class="pc'+(p.connected?'':' disconnected')+(isMe?' is-me':'')+'">'
      +(isHost&&!isH?'<button class="kick-btn" onclick="kickPlayer(\''+p.id+'\')" title="Entfernen">&#10005;</button><button class="kick-btn" style="top:5px;right:25px" onclick="transferHost(\''+p.id+'\')" title="Host \u00fcbertragen">&#128081;</button>':'')
      +'<div class="pn">'+esc(p.name)+(isMe?' (Du)':'')+'</div>'
      +'<div class="ps'+(isH?' host':isD?' delegate':'')+'">'+
        (isH?'Host':isD?'W\u00e4hlt Thema':(p.connected?'Bereit':'Getrennt'))
      +'</div>'+renderTeeth(p,false,false,false)+'</div>';
  }).join('');
  document.getElementById('start-btn-container').style.display=(isHost&&players.length>=1)?'block':'none'}

function renderTeeth(p,anim,showLost,showDead){let upper=[],lower=[];for(let i=0;i<p.maxTeeth;i++){const alive=i<p.teeth,lost=anim&&p.justLost&&i===p.teeth;const row=i%2===0?'upper':'lower';const cls='tooth '+row+(alive?' alive':' dead')+(lost?' just-lost':'');if(i%2===0)upper.push('<div class="'+cls+'"></div>');else lower.push('<div class="'+cls+'"></div>')}
  let h='<div class="teeth-wrap"><div class="teeth-rows"><div class="tooth-row">'+upper.join('')+'</div><div class="tooth-row">'+lower.join('')+'</div></div>';
  if(showDead&&p.eliminated)h+='<div class="heinen-overlay"><div class="ht dead">ES HAT SICH AUSGEHEINT!</div></div>';else if(showLost&&p.justLost&&!p.eliminated)h+='<div class="heinen-overlay"><div class="ht lost">DU WURDEST GEHEINT!</div></div>';return h+'</div>'}

function qCounterText(q){const m=gameState.settings.mode,isEndless=m!=='classic';let t=isEndless?'Frage '+(q.index+1):'Frage '+(q.index+1)+' / '+gameState.totalQuestions;if(gameState.currentDifficulty)t+=' \u2022 '+(DL[gameState.currentDifficulty]||gameState.currentDifficulty);return t}
function renderQuestion(){const q=gameState.question;if(!q)return;document.getElementById('q-counter').textContent=qCounterText(q);document.getElementById('q-text').textContent=q.text;document.getElementById('results-info').style.display='none';qReportedCurrent=false;const ra=document.getElementById('q-report-area');if(ra)ra.innerHTML='<button class="q-report-btn" onclick="reportQuestion()">Problem mit dieser Frage?</button>';renderOptions()}
function renderOptions(){const q=gameState.question;if(!q)return;const L=['A','B','C','D','E','F'];const alive=amIAlive();const changeable=selectedAnswer>=0&&canChangeAnswer();
  document.getElementById('q-options').innerHTML=q.options.map((o,i)=>{if(!alive)return '<button class="ob spectator"><span class="ol">'+L[i]+'</span> '+esc(o)+'</button>';const sel=selectedAnswer===i?' sel':'',dis=selectedAnswer>=0&&!changeable?' dis':'';return '<button class="ob'+sel+dis+'" '+(selectedAnswer<0||changeable?'onclick="submitAnswer('+i+')"':'')+'><span class="ol">'+L[i]+'</span> '+esc(o)+'</button>'}).join('')}
function renderResults(){const q=gameState.question,r=gameState.results;if(!q||!r)return;document.getElementById('q-counter').textContent=qCounterText(q);document.getElementById('q-text').textContent=q.text;const L=['A','B','C','D','E','F'];
  document.getElementById('q-options').innerHTML=q.options.map((o,i)=>{let c='ob dis';if(i===r.correctAnswer)c+=' correct';else if(selectedAnswer===i)c+=' wrong';return '<button class="'+c+'"><span class="ol">'+L[i]+'</span> '+esc(o)+'</button>'}).join('');
  const my=r.playerResults[myId],info=document.getElementById('results-info');info.style.display='block';
  if(my==='correct')info.innerHTML='<span style="color:var(--correct)">Richtige Antwort!</span>';else if(my==='wrong')info.innerHTML='<span style="color:var(--wrong)">Falsche Antwort – Zahn verloren!</span>';else if(my==='timeout'){info.innerHTML='<span style="color:var(--wrong)">Zeit abgelaufen – Zahn verloren!</span>'}else if(!amIAlive())info.innerHTML='<span style="color:var(--text2)">Du bist ausgeschieden.</span>';document.getElementById('q-timer').textContent='\u2014'}
function renderGamePlayers(isHost){const players=gameState.players||[],ph=gameState.phase;document.getElementById('game-players-title').textContent='Spieler*innen ('+players.filter(p=>p.alive).length+' aktiv)';
  let progressHtml='';
  if(ph==='question'){const active=players.filter(p=>p.alive&&p.connected);const answeredCount=active.filter(p=>p.answered).length;const total=active.length;
    const pct=total>0?Math.round(answeredCount/total*100):0;
    progressHtml='<div class="ans-progress" style="position:sticky;bottom:0;background:linear-gradient(to bottom,transparent,var(--bg) 40%);padding:16px 0 2px;margin-top:10px;margin-bottom:0">'+answeredCount+'\u00a0von\u00a0'+total+' haben geantwortet'
      +'<div style="height:4px;background:var(--border);border-radius:2px;margin-top:5px;overflow:hidden"><div style="height:100%;width:'+pct+'%;background:var(--correct);border-radius:2px;transition:width .4s"></div></div></div>'}
  document.getElementById('game-players').innerHTML=players.map(p=>{let ec='';if(!p.alive)ec=' eliminated';if(!p.connected)ec+=' disconnected';if(p.id===myId)ec+=' is-me';
    if(ph==='question'&&p.alive&&p.connected&&p.answered)ec+=' q-answered';
    let status='';
    if(!p.alive)status='<span style="color:var(--wrong)">Ausgeschieden</span>';
    else if(!p.connected)status='<span style="color:var(--text2)">Getrennt</span>';
    else if(ph==='question'&&p.answered)status='<span style="color:var(--correct);font-size:.8rem;font-weight:600;text-transform:none;letter-spacing:0">&#10003; Beantwortet</span>';
    else if(ph==='question')status='<span style="color:var(--text2);font-size:.8rem;text-transform:none;letter-spacing:0">&#8987; Wartet\u2026</span>';
    else status='';
    return '<div class="pc'+ec+'"><div class="pn">'+esc(p.name)+(p.id===myId?' (Du)':'')+'</div><div class="ps">'+status+'</div>'+renderTeeth(p,ph==='results',true,true)+'</div>'}).join('')+progressHtml}
function renderEnd(isHost){const c=document.getElementById('end-content'),players=gameState.players||[],winners=gameState.winners||[];
  const m=gameState.settings.mode,isEndless=m!=='classic',isSP=m==='kfo_singleplayer';
  if(isSP){c.innerHTML='<div style="font-size:2.5rem;margin-bottom:8px">&#127942;</div><div class="winner-display">'+(gameState.finalScore||0)+' Fragen geschafft!</div><div class="end-sub">Schwierigkeit: '+(DL[gameState.currentDifficulty||'leicht']||'?')+'</div>'}
  else if(winners.length===1)c.innerHTML='<div style="font-size:2.5rem;margin-bottom:8px">&#127942;</div><div class="winner-display">'+esc(winners[0].name)+' gewinnt!</div><div class="end-sub">'+(isEndless?'Letzte*r mit Zähnen!':'Die meisten Zähne behalten!')+'</div>';
  else if(winners.length>1)c.innerHTML='<div style="font-size:2.5rem;margin-bottom:8px">&#127942;</div><div class="winner-display">Gleichstand!</div><div class="end-sub">Gewinner*innen: '+winners.map(w=>esc(w.name)).join(', ')+'</div>';
  else c.innerHTML='<div style="font-size:2.5rem;margin-bottom:8px">&#128128;</div><div class="no-winner">Alle ausgeheint!</div><div class="end-sub">Keine*r hat gewonnen.</div>';
  const wIds=new Set(winners.map(w=>w.id));const sorted=[...players].sort((a,b)=>b.teeth-a.teeth);
  document.getElementById('end-players').innerHTML='<div class="pst">Endergebnis</div><div class="pg">'+sorted.map(p=>{const w=wIds.has(p.id);return '<div class="pc'+(!p.alive?' eliminated':'')+'" style="'+(w?'border-color:var(--gold);box-shadow:0 0 16px rgba(255,215,0,.2)':'')+'"><div class="pn">'+(w?'&#128081; ':'')+esc(p.name)+' <span style="color:var(--text2);font-size:.75rem">('+p.teeth+'/'+p.maxTeeth+')</span></div><div class="ps'+(w?' host':(!p.alive?' elim':''))+'">'+
    (w?'Gewinner*in!':(!p.alive?'Ausgeschieden':'Überlebt'))+'</div>'+renderTeeth(p,false,false,false)+'</div>'}).join('')+'</div>';
  document.getElementById('end-actions').innerHTML=isHost?'<button class="btn btn-p" onclick="playAgain()" style="max-width:280px;margin:0 auto">Nochmal spielen</button>':'<div style="color:var(--text2)">Warte auf den*die Host...</div>';
  sessionStorage.removeItem('h_invite');sessionStorage.removeItem('h_name')}

function startTimer(){clearTimer();currentTimeLeft=gameState.timeLeft||0;updateTimer();timerInterval=setInterval(()=>{currentTimeLeft--;if(currentTimeLeft<0)currentTimeLeft=0;
  if(!gameState||gameState.phase!=='question'){clearTimer();return}
  updateTimer();if(currentTimeLeft>=1&&currentTimeLeft<=5){playSound(globalSounds.hurrySound,globalSounds.volHurry)}
  if(currentTimeLeft<=0){playSound(globalSounds.timeoutSound,globalSounds.volTimeout);clearTimer()}},1000)}
function clearTimer(){if(timerInterval){clearInterval(timerInterval);timerInterval=null}}
function updateTimer(){const el=document.getElementById('q-timer');el.textContent=currentTimeLeft;el.classList.toggle('urgent',currentTimeLeft<=5)}
function esc(s){const d=document.createElement('div');d.textContent=s||'';return d.innerHTML}
function showToast(msg,isErr){const el=document.getElementById('toast');el.textContent=msg;el.className='toast show'+(isErr?' error':' success');setTimeout(()=>el.className='toast',3000)}

document.getElementById('login-pass').addEventListener('keydown',e=>{if(e.key==='Enter')doLogin()});
document.getElementById('login-user').addEventListener('keydown',e=>{if(e.key==='Enter')doLogin()});
document.getElementById('join-name').addEventListener('keydown',e=>{if(e.key==='Enter')doJoin()});

// ── Passkey helpers ──────────────────────────────────────────────────────────
function b64url(buf){return btoa(String.fromCharCode(...new Uint8Array(buf))).replace(/\+/g,'-').replace(/\//g,'_').replace(/=/g,'')}
function fromB64url(s){const pad=s+'==='.slice((s.length+3)%4);const bin=atob(pad.replace(/-/g,'+').replace(/_/g,'/'));const b=new Uint8Array(bin.length);for(let i=0;i<bin.length;i++)b[i]=bin.charCodeAt(i);return b.buffer}

// ── Passkey Login ────────────────────────────────────────────────────────────
async function doPasskeyLogin(){
  if(!window.PublicKeyCredential){showToast('Passkeys nicht unterstützt',1);return}
  try{
    const r1=await fetch('/api/passkey/login/begin',{method:'POST'});
    if(!r1.ok){const d=await r1.json();showToast(d.error||'Fehler',1);return}
    const sid=r1.headers.get('X-WA-Session');
    const opts=await r1.json();
    opts.publicKey.challenge=fromB64url(opts.publicKey.challenge);
    if(opts.publicKey.allowCredentials){
      opts.publicKey.allowCredentials=opts.publicKey.allowCredentials.map(c=>({...c,id:fromB64url(c.id)}))
    }
    const assertion=await navigator.credentials.get(opts);
    const body={
      id:assertion.id,
      rawId:b64url(assertion.rawId),
      type:assertion.type,
      response:{
        authenticatorData:b64url(assertion.response.authenticatorData),
        clientDataJSON:b64url(assertion.response.clientDataJSON),
        signature:b64url(assertion.response.signature),
        userHandle:assertion.response.userHandle?b64url(assertion.response.userHandle):null
      }
    };
    const r2=await fetch('/api/passkey/login/finish?sid='+encodeURIComponent(sid),{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
    const d2=await r2.json();
    if(!r2.ok){showToast(d2.error||'Fehler',1);return}
    currentUser=d2;showScreen('dashboard');updateNav()
  }catch(e){
    if(e.name!=='NotAllowedError')showToast('Passkey-Fehler: '+e.message,1)
  }
}

// ── Passkey Verwaltung ───────────────────────────────────────────────────────
async function loadPasskeys(){
  const el=document.getElementById('passkeys-list');
  if(!el)return;
  el.innerHTML='<div style="color:var(--text2);font-size:.82rem">Lade\u2026</div>';
  try{
    const r=await fetch('/api/passkey/credentials');
    const creds=await r.json();
    if(!creds||!creds.length){el.innerHTML='<div style="color:var(--text2);font-size:.82rem">Noch keine Passkeys registriert.</div>';return}
    el.innerHTML=creds.map(c=>'<div class="user-row"><div class="user-info"><span>\uD83D\uDD11 '+esc(c.name)+'</span><span style="font-size:.68rem;color:var(--text2);margin-left:8px">'+esc(c.createdAt)+'</span></div><div class="user-actions"><button class="ib danger" onclick="deletePasskey(\''+esc(c.id)+'\')">Löschen</button></div></div>').join('')
  }catch(e){el.innerHTML='<div style="color:var(--wrong);font-size:.82rem">Fehler beim Laden.</div>'}
}

async function doRegisterPasskey(){
  if(!window.PublicKeyCredential){showToast('Passkeys nicht unterstützt',1);return}
  const name=document.getElementById('passkey-name').value.trim()||'Passkey';
  try{
    const r1=await fetch('/api/passkey/register/begin',{method:'POST'});
    if(!r1.ok){const d=await r1.json();showToast(d.error||'Fehler',1);return}
    const sid=r1.headers.get('X-WA-Session');
    const opts=await r1.json();
    opts.publicKey.challenge=fromB64url(opts.publicKey.challenge);
    opts.publicKey.user.id=fromB64url(opts.publicKey.user.id);
    if(opts.publicKey.excludeCredentials){
      opts.publicKey.excludeCredentials=opts.publicKey.excludeCredentials.map(c=>({...c,id:fromB64url(c.id)}))
    }
    const cred=await navigator.credentials.create(opts);
    const body={
      id:cred.id,
      rawId:b64url(cred.rawId),
      type:cred.type,
      response:{
        attestationObject:b64url(cred.response.attestationObject),
        clientDataJSON:b64url(cred.response.clientDataJSON)
      }
    };
    const r2=await fetch('/api/passkey/register/finish?sid='+encodeURIComponent(sid)+'&name='+encodeURIComponent(name),{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
    const d2=await r2.json();
    if(!r2.ok){showToast(d2.error||'Fehler',1);return}
    showToast('Passkey registriert!',0);document.getElementById('passkey-name').value='';loadPasskeys()
  }catch(e){
    if(e.name!=='NotAllowedError')showToast('Passkey-Fehler: '+e.message,1)
  }
}

async function deletePasskey(id){
  if(!confirm('Passkey wirklich löschen?'))return;
  const r=await fetch('/api/passkey/credentials',{method:'DELETE',headers:{'Content-Type':'application/json'},body:JSON.stringify({ID:id})});
  const d=await r.json();
  if(!r.ok){showToast(d.error||'Fehler',1);return}
  showToast('Passkey gelöscht',0);loadPasskeys()
}
</script>
</body>
</html>` + ""
