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
.pc{background:var(--bg2);border:1px solid var(--border);border-radius:12px;padding:14px;text-align:center;position:relative;transition:all .3s}.pc.eliminated{opacity:.45;border-color:var(--wrong)}.pc.disconnected{opacity:.35}.pc.is-me{border-color:var(--accent);border-width:2px}.pn{font-weight:600;font-size:.9rem;margin-bottom:3px}.ps{font-size:.6rem;color:var(--text2);text-transform:uppercase;letter-spacing:1px;margin-bottom:8px}.ps.host{color:var(--gold)}.ps.elim{color:var(--wrong)}.ps.delegate{color:var(--accent2)}
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
.loading-container{text-align:center;padding:50px 20px}.loading-spinner{width:40px;height:40px;border:3px solid var(--border);border-top-color:var(--accent);border-radius:50%;animation:spin .8s linear infinite;margin:0 auto 14px}@keyframes spin{to{transform:rotate(360deg)}}.loading-text{color:var(--text2);font-size:.85rem}
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
<button class="btn btn-s" onclick="goHome()" style="margin-top:10px">Zurück</button></div></div>

<div id="screen-dashboard" class="screen"><div class="dc"><div style="font-size:1rem;color:var(--text2);margin-bottom:24px" id="welcome-text"></div>
<div class="da"><button class="btn btn-p" onclick="showScreen('create')">Neues Spiel erstellen</button>
<button class="btn btn-s" id="btn-admin" onclick="loadAdmin();showScreen('admin')" style="display:none">Admin Control Panel</button>
<button class="btn btn-s" onclick="showScreen('pw')">Passwort ändern</button>
<button class="btn btn-s" onclick="goHome()">Zurück zur Startseite</button></div></div></div>

<div id="screen-pw" class="screen"><div class="login-container"><div class="login-title">Passwort ändern</div>
<div class="ig"><label>Altes Passwort</label><input type="password" id="pw-old"/></div>
<div class="ig"><label>Neues Passwort</label><input type="password" id="pw-new"/></div>
<div class="ig"><label>Bestätigen</label><input type="password" id="pw-new2"/></div>
<button class="btn btn-p" onclick="changePw()">Ändern</button>
<button class="btn btn-s" onclick="showScreen('dashboard')" style="margin-top:10px">Zurück</button></div></div>

<div id="screen-admin" class="screen"><div class="admin-container">
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

<div id="screen-loading" class="screen"><div class="loading-container"><div class="loading-spinner"></div><div class="loading-text">Fragen werden generiert...</div></div></div>
<div id="screen-refill" class="screen"><div class="loading-container"><div class="loading-spinner"></div><div class="loading-text">Neue Fragen werden nachgeladen...</div></div></div>
<div id="screen-tutorial" class="screen"><div class="tutorial-content" id="tutorial-content"></div><div style="text-align:center;margin-top:20px" id="tutorial-actions"></div></div>
<div id="screen-intro" class="screen"><div class="intro-overlay"><div class="intro-title">HEINEN</div><div class="intro-slogan">Das einzige Spiel mit Spaß-Garantie!<br>Kein Spaß? Geld zurück!</div></div></div>
<div id="screen-game" class="screen">
<div class="qc"><div class="qh"><div class="qn" id="q-counter">Frage 1/10</div><div class="timer" id="q-timer">20</div></div>
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
<div id="audio-controls" style="position:fixed;bottom:20px;right:20px;z-index:90;display:none;align-items:center;gap:6px;background:var(--surface);border:1px solid var(--border);border-radius:22px;padding:4px 10px">
<button class="mute-btn" id="mute-btn" onclick="toggleMute()" title="Stumm schalten" style="position:static;display:flex;width:32px;height:32px;font-size:1rem;border:none;background:none">&#128266;</button>
<input type="range" id="bg-vol-slider" min="0" max="1" step="0.05" value="0.2" style="width:80px;cursor:pointer" oninput="adjustBgVol(this.value)"/>
</div>
<div class="toast" id="toast"></div>
<audio id="bg-audio" preload="auto" loop></audio>
<script src="https://cdnjs.cloudflare.com/ajax/libs/qrcodejs/1.0.0/qrcode.min.js"></script>
<script>
let ws=null,myId='',inviteCode='',gameState=null,selectedAnswer=-1,timerInterval=null,currentTimeLeft=0,shuffleAnimating=false;
let currentUser=null,joinPending='',joinNeedsPw=false,joinLobbyPw='';
let globalSounds={},bgMuted=false,bgStarted=false,tutorialHtml='';
const DL={leicht:'Leicht',mittel:'Mittel',schwer:'Schwer',extrem:'Extrem'};
const ML={classic:'Klassisch',elimination:'Elimination',kfo_battle_royale:'KFO Battle Royale',kfo_singleplayer:'KFO Singleplayer'};
const openaiModels=['gpt-5.4','gpt-5.4-mini','gpt-5.4-nano','gpt-5','gpt-4o','gpt-4o-mini','gpt-4.1'];
const anthropicModels=['claude-opus-4-6','claude-sonnet-4-6','claude-haiku-4-5-20251001'];
const soundDefs=[{key:'intro_sound',label:'Intro-Sound',id:'file-intro'},{key:'background_sound',label:'Background-Song',id:'file-bg'},{key:'wrong_sound',label:'Falsch-Sound',id:'file-wrong'},{key:'answer_sound',label:'Antwort-Sound',id:'file-answer'},{key:'hurry_sound',label:'Zeit-läuft-ab-Sound',id:'file-hurry'},{key:'timeout_sound',label:'Zeit-abgelaufen-Sound',id:'file-timeout'},{key:'question_sound',label:'Nächste-Frage-Sound',id:'file-question'},{key:'allwrong_sound',label:'Alle-falsch-Sound',id:'file-allwrong'},{key:'allcorrect_sound',label:'Alle-richtig-Sound',id:'file-allcorrect'}];
const volDefs=[{id:'vol-intro',key:'vol_intro',label:'Intro',def:'0.6'},{id:'vol-bg',key:'vol_background',label:'Hintergrund',def:'0.2'},{id:'vol-wrong',key:'vol_wrong',label:'Falsch',def:'0.6'},{id:'vol-answer',key:'vol_answer',label:'Antwort',def:'0.6'},{id:'vol-hurry',key:'vol_hurry',label:'Zeit läuft ab',def:'0.5'},{id:'vol-timeout',key:'vol_timeout',label:'Zeit abgelaufen',def:'0.6'},{id:'vol-question',key:'vol_question',label:'Nächste Frage',def:'0.5'},{id:'vol-allwrong',key:'vol_allwrong',label:'Alle falsch',def:'0.6'},{id:'vol-allcorrect',key:'vol_allcorrect',label:'Alle richtig',def:'0.6'}];
function getCookie(name){const m=document.cookie.match('(^|;)\\s*'+name+'\\s*=\\s*([^;]+)');return m?m[2]:''}
async function apiFetch(url,opts={}){const method=opts.method||'GET';const headers=opts.headers||{};if(method!=='GET'&&!opts.nocsrf){headers['X-CSRF-Token']=getCookie('heinen_csrf')}return fetch(url,{...opts,headers})}

(async function(){
  const p=new URLSearchParams(location.search);const inv=p.get('join');
  await loadSounds();await loadTutorial();
  if(inv){joinPending=inv;showScreen('join');return}
  // Try auto-reconnect
  const ri=sessionStorage.getItem('h_invite'),rn=sessionStorage.getItem('h_name');
  if(ri&&rn){connectWS();const ck=setInterval(()=>{if(ws&&ws.readyState===1){clearInterval(ck);send('reconnect',{Name:rn,InviteCode:ri})}},100)}
  try{const r=await fetch('/api/me');if(r.ok){currentUser=await r.json()}}catch(e){}
  showScreen('home');updateNav();loadLobbies();
})();

async function loadSounds(){try{const r=await fetch('/api/global-sounds');const d=await r.json();globalSounds={introSound:d.intro_sound||'',backgroundSound:d.background_sound||'',wrongSound:d.wrong_sound||'',answerSound:d.answer_sound||'',hurrySound:d.hurry_sound||'',timeoutSound:d.timeout_sound||'',questionSound:d.question_sound||'',allwrongSound:d.allwrong_sound||'',allcorrectSound:d.allcorrect_sound||'',volIntro:parseFloat(d.vol_intro)||0.6,volBg:parseFloat(d.vol_background)||0.2,volWrong:parseFloat(d.vol_wrong)||0.6,volAnswer:parseFloat(d.vol_answer)||0.6,volHurry:parseFloat(d.vol_hurry)||0.5,volTimeout:parseFloat(d.vol_timeout)||0.6,volQuestion:parseFloat(d.vol_question)||0.5,volAllwrong:parseFloat(d.vol_allwrong)||0.6,volAllcorrect:parseFloat(d.vol_allcorrect)||0.6}}catch(e){}}
async function loadTutorial(){try{const r=await fetch('/api/tutorial');const d=await r.json();tutorialHtml=markdownToHtml(d.content||'')}catch(e){}}
async function loadLobbies(){try{const r=await fetch('/api/lobbies');const lobbies=await r.json();const el=document.getElementById('lobby-list');
  if(!lobbies||lobbies.length===0){el.innerHTML='<div style="color:var(--text2);text-align:center;padding:20px">Keine offenen Lobbys vorhanden.</div>';return}
  el.innerHTML=lobbies.map(l=>'<div class="lobby-card"><div class="lc-info"><div class="lc-name">'+esc(l.lobbyName)+' <span class="lc-badge '+(l.lobbyMode==='open'?'open':'pw')+'">'+(l.lobbyMode==='open'?'Offen':'Passwort')+'</span></div><div class="lc-details">'+esc(l.topic)+' \u2022 '+(ML[l.mode]||l.mode)+' \u2022 '+l.players+' Spieler*in(nen)</div></div><button class="btn btn-sm btn-p" onclick="joinLobby(\''+l.inviteCode+'\',\''+(l.lobbyMode==='password'?'pw':'open')+'\')">Beitreten</button></div>').join('')
}catch(e){}}

function markdownToHtml(md){md=md.replace(/\\\*/g,'\u2605STAR\u2605');md=md.replace(/^### (.+)$/gm,'<h3>$1</h3>').replace(/^## (.+)$/gm,'<h2>$1</h2>').replace(/^# (.+)$/gm,'<h1>$1</h1>').replace(/\*\*(.+?)\*\*/g,'<strong>$1</strong>').replace(/\*(.+?)\*/g,'<em>$1</em>').replace(/^- (.+)$/gm,'<li>$1</li>').replace(/(<li>[\s\S]*?<\/li>)/g,function(m){return '<ul>'+m+'</ul>'}).replace(/<\/ul>\s*<ul>/g,'').replace(/^(?!<[hulo])(.*\S.*)$/gm,'<p>$1</p>');md=md.replace(/\u2605STAR\u2605/g,'*');return md}

function updateNav(){const n=document.getElementById('nav-bar');const hpb=document.getElementById('home-panel-btn');const hlb=document.getElementById('home-login-btn');
  if(currentUser){n.innerHTML='<span style="color:var(--text2);font-size:.75rem;align-self:center">'+esc(currentUser.username)+'</span><button class="nav-btn" onclick="doLogout()">Abmelden</button>';if(hpb)hpb.style.display='inline-block';if(hlb)hlb.style.display='none'}
  else{n.innerHTML='';if(hpb)hpb.style.display='none';if(hlb)hlb.style.display='inline-block'}}
function goHome(){if(gameState&&['question','results','intro','loading','refill','tutorial'].includes(gameState.phase)){if(!confirm('Spiel wirklich verlassen?'))return}
  gameState=null;stopBg();if(ws){try{ws.close()}catch(e){}}ws=null;myId='';inviteCode='';sessionStorage.removeItem('h_invite');sessionStorage.removeItem('h_name');
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
  const ur=await fetch('/api/users');const users=await ur.json();renderUsers(users||[]);loadLogs()
}
function onProviderChange(){const p=document.getElementById('ai-provider').value;document.getElementById('ai-openai-row').style.display=p==='openai'?'flex':'none';document.getElementById('ai-anthropic-row').style.display=p==='anthropic'?'flex':'none';const sel=document.getElementById('ai-model');const models=p==='anthropic'?anthropicModels:openaiModels;const cur=sel.value;sel.innerHTML=models.map(m=>'<option value="'+m+'">'+m+'</option>').join('');if(models.includes(cur))sel.value=cur}
async function saveAIConfig(){const b={Provider:document.getElementById('ai-provider').value,Model:document.getElementById('ai-model').value,IntroDelay:document.getElementById('ai-intro-delay').value};
  volDefs.forEach(v=>{const el=document.getElementById(v.id);if(el)b[v.key.split('_').map((w,i)=>i?w[0].toUpperCase()+w.slice(1):w.charAt(0).toUpperCase()+w.slice(1)).join('')]=el.value});
  const ok=document.getElementById('ai-openai-key').value.trim();if(ok)b.OpenaiKey=ok;const ak=document.getElementById('ai-anthropic-key').value.trim();if(ak)b.AnthropicKey=ak;
  await apiFetch('/api/ai-config',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(b)});showToast('Gespeichert',0);await loadSounds();loadAdmin()}
async function testAI(){const st=document.getElementById('ai-status');st.textContent='Teste...';const p=document.getElementById('ai-provider').value;const k=p==='anthropic'?document.getElementById('ai-anthropic-key').value.trim():document.getElementById('ai-openai-key').value.trim();const m=document.getElementById('ai-model').value;const r=await apiFetch('/api/test-ai',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({Provider:p,Key:k,Model:m})});const d=await r.json();st.textContent=d.message;st.style.color=d.ok?'var(--correct)':'var(--wrong)'}
function renderUsers(users){document.getElementById('users-list').innerHTML=users.map(u=>{const dis=u.isOnlyAdmin;const hashClass=u.passwordHashType||'unknown';return '<div class="user-row"><div class="user-info"><span>'+esc(u.username)+'</span><span class="ub '+(u.isAdmin?'admin':'user')+'">'+(u.isAdmin?'Admin':'User')+'</span><span class="ub hash '+hashClass+'">'+hashClass+'</span></div><div class="user-actions">'+(u.isAdmin?'<button class="ib"'+(dis?' disabled':' onclick="toggleAdmin('+u.id+',false)"')+'>Degradieren</button>':'<button class="ib" onclick="toggleAdmin('+u.id+',true)">Zum Admin</button>')+'<button class="ib danger"'+(dis?' disabled':' onclick="deleteUser('+u.id+',\''+esc(u.username)+'\')"')+'>Löschen</button></div></div>'}).join('')}
async function addUser(){const n=document.getElementById('new-user-name').value.trim(),p=document.getElementById('new-user-pass').value.trim(),a=document.getElementById('new-user-role').value==='1';if(!n||!p){showToast('Felder ausfüllen',1);return}const r=await apiFetch('/api/users',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({username:n,password:p,isAdmin:a})});const d=await r.json();if(d.error){showToast(d.error,1);return}document.getElementById('new-user-name').value='';document.getElementById('new-user-pass').value='';showToast('Angelegt',0);loadAdmin()}
async function deleteUser(id,name){if(!confirm(name+' löschen?'))return;await apiFetch('/api/users',{method:'DELETE',headers:{'Content-Type':'application/json'},body:JSON.stringify({id})});loadAdmin()}
async function toggleAdmin(id,make){await apiFetch('/api/users',{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify({id,isAdmin:make})});loadAdmin()}

function uploadSound(type,inputId){const input=document.getElementById(inputId);const f=input&&input.files[0];if(!f){showToast('Datei auswählen',1);return}if(f.size>20*1024*1024){showToast('Max 20 MB',1);return}
  const fd=new FormData();fd.append('type',type);fd.append('file',f);const xhr=new XMLHttpRequest();xhr.setRequestHeader('X-CSRF-Token',getCookie('heinen_csrf'));xhr.onload=async function(){try{const d=JSON.parse(xhr.responseText);if(d.ok){showToast('Hochgeladen!',0);input.value='';await loadSounds();loadAdmin()}else showToast(d.error||'Fehler',1)}catch(e){showToast('Serverfehler',1)}};xhr.onerror=()=>showToast('Netzwerkfehler',1);xhr.open('POST','/api/sounds');xhr.send(fd)}
async function deleteSound(type){await apiFetch('/api/sounds',{method:'DELETE',headers:{'Content-Type':'application/json'},body:JSON.stringify({type})});showToast('Entfernt',0);await loadSounds();loadAdmin()}
let pAudio=null;async function previewSound(type){await loadSounds();const map={intro_sound:'introSound',background_sound:'backgroundSound',wrong_sound:'wrongSound',answer_sound:'answerSound',hurry_sound:'hurrySound',timeout_sound:'timeoutSound',question_sound:'questionSound',allwrong_sound:'allwrongSound',allcorrect_sound:'allcorrectSound'};const url=globalSounds[map[type]];if(!url){showToast('Nicht hinterlegt',1);return}if(pAudio){pAudio.pause();pAudio=null}pAudio=new Audio(url);pAudio.volume=0.6;const p=pAudio.play();if(p&&p.catch)p.catch(e=>showToast('Fehler',1));setTimeout(()=>{if(pAudio){pAudio.pause();pAudio=null}},10000)}
function previewVol(soundKey,sliderId){const map={intro_sound:'introSound',background_sound:'backgroundSound',wrong_sound:'wrongSound',answer_sound:'answerSound',hurry_sound:'hurrySound',timeout_sound:'timeoutSound',question_sound:'questionSound',allwrong_sound:'allwrongSound',allcorrect_sound:'allcorrectSound'};const url=globalSounds[map[soundKey]];if(!url){showToast('Nicht hinterlegt',1);return}const vol=parseFloat(document.getElementById(sliderId).value);if(pAudio){pAudio.pause();pAudio=null}pAudio=new Audio(url);pAudio.volume=vol;const p=pAudio.play();if(p&&p.catch)p.catch(()=>{});setTimeout(()=>{if(pAudio){pAudio.pause();pAudio=null}},5000)}
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
function toggleMute(){bgMuted=!bgMuted;const a=document.getElementById('bg-audio'),b=document.getElementById('mute-btn');if(bgMuted){a.pause();b.innerHTML='&#128264;';b.classList.add('muted')}else{if(bgStarted)tryPlayBg(3);b.innerHTML='&#128266;';b.classList.remove('muted')}}
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
  case 'kicked':showToast(msg.payload.message,1);gameState=null;stopBg();sessionStorage.removeItem('h_invite');sessionStorage.removeItem('h_name');showScreen('home');break;
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
function submitAnswer(i){if(selectedAnswer>=0)return;selectedAnswer=i;playSound(globalSounds.answerSound,globalSounds.volAnswer);send('answer',{answer:i});renderOptions()}
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
  const pi=document.getElementById('l-playintro');if(pi)s.PlayIntro=pi.checked;}
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
  if(ph==='lobby'){showScreen('lobby');renderLobby(isHost)}
  else if(ph==='loading'){showScreen('loading')}else if(ph==='refill'){showScreen('refill')}
  else if(ph==='tutorial'){showScreen('tutorial');renderTutorial(isHost)}
  else if(ph==='intro'){showScreen('intro');playSound(globalSounds.introSound,globalSounds.volIntro);startBg()}
  else if(ph==='question'){showScreen('game');if(lastPhase!=='question'){selectedAnswer=-1;playSound(globalSounds.questionSound,globalSounds.volQuestion)}renderQuestion();renderGamePlayers(isHost);startTimer()}
  else if(ph==='results'){showScreen('game');renderResults();renderGamePlayers(isHost);clearTimer();
    if(gameState.allWrong&&gameState.players.filter(p=>p.alive&&p.connected).length>=2)playSound(globalSounds.allwrongSound,globalSounds.volAllwrong);
    else if(gameState.allCorrect)playSound(globalSounds.allcorrectSound,globalSounds.volAllcorrect);
    else if(gameState.someoneLost)playSound(globalSounds.wrongSound,globalSounds.volWrong)}
  else if(ph==='end'){clearTimer();showScreen('end');renderEnd(isHost);stopBg()}
  else if(ph==='error'){clearTimer();stopBg();showScreen('error');document.getElementById('error-msg').textContent=gameState.errorMsg||'Unbekannter Fehler'}
  document.getElementById('host-controls').style.display=(isHost&&(ph==='question'||ph==='results'))?'block':'none';lastPhase=ph}

function renderTutorial(isHost){document.getElementById('tutorial-content').innerHTML=tutorialHtml||'<p>Lade Tutorial...</p>';document.getElementById('tutorial-actions').innerHTML=isHost?'<button class="btn btn-p" onclick="skipTutorial()" style="max-width:280px">Weiter</button>':'<div style="color:var(--text2);font-size:.85rem">Warte auf den*die Host...</div>'}

function renderLobby(isHost){
  const link=location.origin+'/?join='+inviteCode;document.getElementById('invite-link').textContent=link;
  const qrEl=document.getElementById('qr-code');qrEl.innerHTML='';try{new QRCode(qrEl,{text:link,width:180,height:180,colorDark:'#e4e4ef',colorLight:'#0a0a0f'})}catch(e){}
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
    const oOpts=[2,3,4].map(n=>'<option value="'+n+'"'+(s.numOptions===n?' selected':'')+'>'+n+'</option>').join('');
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
      +'</div><div class="ig" style="margin-top:8px"><label><input type="checkbox" id="l-tutorial"'+(s.showTutorial?' checked':'')+' onchange="sendLS()"/>Tutorial anzeigen</label></div><div class="ig" style="margin-top:4px"><label><input type="checkbox" id="l-playintro"'+(s.playIntro!==false?' checked':'')+' onchange="sendLS()"/>Intro-Sound abspielen</label></div><div class="ig" style="margin-top:4px"><label><input type="checkbox" id="l-websearch"'+(s.webSearch?' checked':'')+' onchange="sendLS()"/>Internet-Recherche f\u00fcr Fragen</label></div>'
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
function renderQuestion(){const q=gameState.question;if(!q)return;document.getElementById('q-counter').textContent=qCounterText(q);document.getElementById('q-text').textContent=q.text;document.getElementById('results-info').style.display='none';renderOptions()}
function renderOptions(){const q=gameState.question;if(!q)return;const L=['A','B','C','D','E','F'];const alive=amIAlive();
  document.getElementById('q-options').innerHTML=q.options.map((o,i)=>{if(!alive)return '<button class="ob spectator"><span class="ol">'+L[i]+'</span> '+esc(o)+'</button>';const sel=selectedAnswer===i?' sel':'',dis=selectedAnswer>=0?' dis':'';return '<button class="ob'+sel+dis+'" '+(selectedAnswer<0?'onclick="submitAnswer('+i+')"':'')+'><span class="ol">'+L[i]+'</span> '+esc(o)+'</button>'}).join('')}
function renderResults(){const q=gameState.question,r=gameState.results;if(!q||!r)return;document.getElementById('q-counter').textContent=qCounterText(q);document.getElementById('q-text').textContent=q.text;const L=['A','B','C','D','E','F'];
  document.getElementById('q-options').innerHTML=q.options.map((o,i)=>{let c='ob dis';if(i===r.correctAnswer)c+=' correct';else if(selectedAnswer===i)c+=' wrong';return '<button class="'+c+'"><span class="ol">'+L[i]+'</span> '+esc(o)+'</button>'}).join('');
  const my=r.playerResults[myId],info=document.getElementById('results-info');info.style.display='block';
  if(my==='correct')info.innerHTML='<span style="color:var(--correct)">Richtige Antwort!</span>';else if(my==='wrong')info.innerHTML='<span style="color:var(--wrong)">Falsche Antwort – Zahn verloren!</span>';else if(my==='timeout'){info.innerHTML='<span style="color:var(--wrong)">Zeit abgelaufen – Zahn verloren!</span>'}else if(!amIAlive())info.innerHTML='<span style="color:var(--text2)">Du bist ausgeschieden.</span>';document.getElementById('q-timer').textContent='\u2014'}
function renderGamePlayers(isHost){const players=gameState.players||[],ph=gameState.phase;document.getElementById('game-players-title').textContent='Spieler*innen ('+players.filter(p=>p.alive).length+' aktiv)';
  document.getElementById('game-players').innerHTML=players.map(p=>{let ec='';if(!p.alive)ec=' eliminated';if(!p.connected)ec+=' disconnected';if(p.id===myId)ec+=' is-me';
    let status='';
    if(!p.alive)status='<span style="color:var(--wrong)">Ausgeschieden</span>';
    else if(!p.connected)status='<span style="color:var(--text2)">Getrennt</span>';
    else if(ph==='question'&&p.answered)status='<span style="color:var(--correct)">&#10003; Beantwortet</span>';
    else if(ph==='question')status='<span style="color:var(--text2)">&#8987; Wartet...</span>';
    else status='';
    return '<div class="pc'+ec+'"><div class="pn">'+esc(p.name)+(p.id===myId?' (Du)':'')+'</div><div class="ps">'+status+'</div>'+renderTeeth(p,ph==='results',true,true)+'</div>'}).join('')}
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
</script>
</body>
</html>` + ""
