// WebGL spinner showing a real human upper first molar (mesh from BodyParts3D,
// see molar-mesh.js for the licence). Replaces the earlier Zdog toy molar. Needs
// molar-mesh.js loaded first; falls back to nothing if WebGL is unavailable.
(function () {
  var VS = [
    'attribute vec3 aPos;attribute vec3 aNor;',
    'uniform mat4 uMV;uniform mat4 uP;',
    'varying vec3 vN;varying vec3 vV;varying float vY;',
    'void main(){vec4 p=uMV*vec4(aPos,1.0);vN=mat3(uMV)*aNor;vV=-p.xyz;vY=aPos.y;gl_Position=uP*p;}'
  ].join('\n');

  // Enamel crown blending into a yellower, matte root below the cemento-enamel
  // junction. Wrap lighting + soft rim fakes the translucency of enamel.
  var FS = [
    'precision mediump float;',
    'varying vec3 vN;varying vec3 vV;varying float vY;',
    'void main(){',
    ' vec3 n=normalize(vN);vec3 v=normalize(vV);',
    ' if(!gl_FrontFacing)n=-n;',
    ' vec3 L1=normalize(vec3(-0.5,0.8,0.7));vec3 L2=normalize(vec3(0.7,0.1,0.4));vec3 L3=normalize(vec3(0.0,-0.6,-0.8));',
    ' float cej=smoothstep(-0.06,0.16,vY);',
    ' vec3 enamel=vec3(0.95,0.93,0.86);vec3 root=vec3(0.88,0.8,0.64);',
    ' vec3 base=mix(root,enamel,cej);',
    ' float w=0.35;',
    ' float d1=max((dot(n,L1)+w)/(1.0+w),0.0);float d2=max((dot(n,L2)+w)/(1.0+w),0.0);float d3=max(dot(n,L3),0.0);',
    ' vec3 col=base*(0.22+0.78*d1+0.28*d2)+vec3(1.0,0.85,0.6)*0.18*d3;',
    ' float gloss=mix(18.0,70.0,cej);float ks=mix(0.08,0.45,cej);',
    ' vec3 h1=normalize(L1+v);vec3 h2=normalize(L2+v);',
    ' col+=ks*(pow(max(dot(n,h1),0.0),gloss)+0.4*pow(max(dot(n,h2),0.0),gloss));',
    ' float rim=pow(1.0-max(dot(n,v),0.0),3.0);',
    ' col+=rim*mix(vec3(0.25,0.2,0.12),vec3(0.35,0.36,0.4),cej);',
    ' col*=1.0-0.18*smoothstep(-0.4,-1.0,vY);',
    ' gl_FragColor=vec4(pow(col,vec3(0.95)),1.0);',
    '}'
  ].join('\n');

  function b64(s, Type) {
    var bin = atob(s), u8 = new Uint8Array(bin.length);
    for (var i = 0; i < bin.length; i++) u8[i] = bin.charCodeAt(i);
    return new Type(u8.buffer);
  }

  var meshData = null;
  function getMesh() {
    if (!meshData && window.MOLAR_MESH) {
      var m = window.MOLAR_MESH;
      meshData = { p: b64(m.p, Int16Array), n: b64(m.n, Int8Array), i: b64(m.i, Uint16Array) };
    }
    return meshData;
  }

  function persp(fov, aspect, near, far) {
    var f = 1 / Math.tan(fov / 2), nf = 1 / (near - far);
    return [f / aspect, 0, 0, 0, 0, f, 0, 0, 0, 0, (far + near) * nf, -1, 0, 0, 2 * far * near * nf, 0];
  }

  // model-view: scale, spin about y, tilt about x, push back along -z
  function modelView(spin, tilt, s, dist, dy) {
    var cy = Math.cos(spin), sy = Math.sin(spin), cx = Math.cos(tilt), sx = Math.sin(tilt);
    // R = Rx(tilt) * Ry(spin), column-major
    return [
      s * cy, s * sx * sy, -s * cx * sy, 0,
      0, s * cx, s * sx, 0,
      s * sy, -s * sx * cy, s * cx * cy, 0,
      0, dy, -dist, 1
    ];
  }

  function compile(gl, type, src) {
    var sh = gl.createShader(type);
    gl.shaderSource(sh, src);
    gl.compileShader(sh);
    if (!gl.getShaderParameter(sh, gl.COMPILE_STATUS)) throw new Error(gl.getShaderInfoLog(sh));
    return sh;
  }

  function setupGL(canvas, mesh) {
    var gl = canvas.getContext('webgl', { alpha: true, antialias: true, premultipliedAlpha: true });
    if (!gl) return null;
    var prog = gl.createProgram();
    gl.attachShader(prog, compile(gl, gl.VERTEX_SHADER, VS));
    gl.attachShader(prog, compile(gl, gl.FRAGMENT_SHADER, FS));
    gl.linkProgram(prog);
    gl.useProgram(prog);
    function attr(name, data, type, norm) {
      var b = gl.createBuffer();
      gl.bindBuffer(gl.ARRAY_BUFFER, b);
      gl.bufferData(gl.ARRAY_BUFFER, data, gl.STATIC_DRAW);
      var loc = gl.getAttribLocation(prog, name);
      gl.enableVertexAttribArray(loc);
      gl.vertexAttribPointer(loc, 3, type, norm, 0, 0);
    }
    attr('aPos', mesh.p, gl.SHORT, true);
    attr('aNor', mesh.n, gl.BYTE, true);
    var ib = gl.createBuffer();
    gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, ib);
    gl.bufferData(gl.ELEMENT_ARRAY_BUFFER, mesh.i, gl.STATIC_DRAW);
    gl.enable(gl.DEPTH_TEST);
    gl.clearColor(0, 0, 0, 0);
    return {
      gl: gl, count: mesh.i.length,
      uMV: gl.getUniformLocation(prog, 'uMV'),
      uP: gl.getUniformLocation(prog, 'uP')
    };
  }

  function mount(el) {
    var isMini = el.classList.contains('tooth-spinner-mini');
    var size = isMini ? 150 : 220;
    var dpr = Math.min(window.devicePixelRatio || 1, 2);

    var canvas = document.createElement('canvas');
    canvas.className = 'tooth-canvas';
    canvas.width = Math.round(size * dpr);
    canvas.height = Math.round(size * dpr);
    canvas.style.width = size + 'px';
    canvas.style.height = size + 'px';

    var stage = document.createElement('div');
    stage.className = 'tooth-stage';
    stage.appendChild(canvas);

    var orbit = document.createElement('div');
    orbit.className = 'orbit-rot';
    var html = '<div class="orbit-dot head"></div>';
    var trail = [[8, .85], [16, .7], [24, .58], [32, .48], [40, .4], [50, .32], [60, .25], [72, .18], [86, .12], [102, .07], [120, .04]];
    trail.forEach(function (c) {
      html += '<div class="orbit-dot" style="--a:-' + c[0] + 'deg;--o:' + c[1] + '"></div>';
    });
    orbit.innerHTML = html;
    stage.appendChild(orbit);

    el.innerHTML = '';
    el.appendChild(stage);

    var mesh = getMesh();
    var r = null;
    try { r = mesh && setupGL(canvas, mesh); } catch (e) { r = null; }
    if (!r) return;
    var gl = r.gl;
    gl.viewport(0, 0, canvas.width, canvas.height);
    gl.uniformMatrix4fv(r.uP, false, persp(0.5, 1, 0.1, 20));

    var spin = 0, last = 0;
    function frame(t) {
      if (canvas.offsetParent !== null) {
        var dt = last ? Math.min(t - last, 100) : 16;
        spin += dt * 0.0015;
        gl.uniformMatrix4fv(r.uMV, false, modelView(spin, 0.3, 0.78, 4.6, 0.02));
        gl.clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT);
        gl.drawElements(gl.TRIANGLES, r.count, gl.UNSIGNED_SHORT, 0);
      }
      last = t;
      requestAnimationFrame(frame);
    }
    requestAnimationFrame(frame);
  }

  function mountAll() {
    document.querySelectorAll('.tooth-mount').forEach(function (el) {
      if (!el.querySelector('.tooth-stage')) mount(el);
    });
  }

  window.Tooth3D = { mount: mount, mountAll: mountAll };
})();
