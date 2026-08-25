// Zdog-based 3D molar spinner. Replaces the flat rotating SVG (which vanished to
// a line at 90°). Real geometry: crown = bottom dome + cylinder + top dome, with
// four cusps on top and two roots hanging below.
(function () {
  var TAU = Zdog.TAU;

  var COL = {
    enamel:      '#f5ebd2',
    enamelBack:  '#c9b98f',
    cusp:        '#fbf3d9',
    cuspBack:    '#d5c69c',
    root:        '#e8d9b0',
    rootBack:    '#b5a577',
    fissure:     '#8b7a55'
  };

  function buildTooth(illo) {
    var tooth = new Zdog.Group({ addTo: illo });

    // Roots — two cones tapering downward. Base sits buried inside the crown
    // capsule so the flat base disc is hidden by the crown's curved surface.
    // Cone default: base at origin, apex at +z=length. rotate.x = -TAU/4
    // turns +z into +y (down in screen coords).
    [-1, 1].forEach(function (sign) {
      new Zdog.Cone({
        addTo: tooth,
        diameter: 26,
        length: 60,
        stroke: false,
        color: COL.root,
        backface: COL.rootBack,
        translate: { x: sign * 14, y: 24 },
        rotate: { x: -TAU / 4 }
      });
    });

    // Crown — a single capsule (Shape with thick stroke). This is one continuous
    // primitive, no seams from cylinder/hemisphere junctions.
    new Zdog.Shape({
      addTo: tooth,
      path: [{ y: -14 }, { y: 14 }],
      stroke: 78,
      color: COL.enamel,
      backface: COL.enamelBack
    });

    // Four cusps as spheres (Shape with stroke, no path length = sphere) sitting
    // on top of the crown. Overlap with the crown gives an organic bump.
    var cusps = [
      { x: -16, z: -16 },
      { x:  16, z: -16 },
      { x: -16, z:  16 },
      { x:  16, z:  16 }
    ];
    cusps.forEach(function (c) {
      new Zdog.Shape({
        addTo: tooth,
        stroke: 28,
        translate: { x: c.x, y: -46, z: c.z },
        color: COL.cusp,
        backface: COL.cuspBack
      });
    });

    return tooth;
  }

  function mount(el) {
    var isMini = el.classList.contains('tooth-spinner-mini');
    var size = isMini ? 150 : 220;

    // Zdog reads canvas.width/height as CSS pixels and internally scales the
    // backing store by devicePixelRatio (see setSizeCanvas). So we pass CSS size.
    var canvas = document.createElement('canvas');
    canvas.className = 'tooth-canvas';
    canvas.width = size;
    canvas.height = size;

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

    var illo = new Zdog.Illustration({
      element: canvas,
      zoom: size / 260,
      rotate: { x: -0.32 },
      translate: { y: -30 }
    });

    var tooth = buildTooth(illo);

    function frame() {
      if (canvas.offsetParent !== null) {
        tooth.rotate.y += 0.025;
        illo.updateRenderGraph();
      }
      requestAnimationFrame(frame);
    }
    frame();
  }

  function mountAll() {
    document.querySelectorAll('.tooth-mount').forEach(function (el) {
      if (!el.querySelector('.tooth-stage')) mount(el);
    });
  }

  window.Tooth3D = { mount: mount, mountAll: mountAll };
})();
