// Obsidian-like tags graph using D3 v7
// - Builds a tag co-occurrence graph from /static/index.json
// - Zoom, pan, drag, hover highlights, click to open tag page

let width, height, svg, simulation, color, container, g, linkGroup, nodeGroup, labelGroup; let graphZoom; // d3 zoom behavior instance used for fit-to-view
// Visual mode flags
let monochrome = true; // show single-color graph
let monoColor = null; // will be derived from site CSS when initialized

// label spacing (pixels) and runtime sizes
const LABEL_OFFSET = 18;
let currentSizeScale = null;

// Dynamically load a script and return a Promise that resolves on load
function loadScript(src) {
    return new Promise((resolve, reject) => {
        try {
            const s = document.createElement('script');
            s.src = src;
            s.async = true;
            s.onload = () => resolve();
            s.onerror = (e) => reject(e);
            document.head.appendChild(s);
        } catch (e) {
            reject(e);
        }
    });
}

// Ensure D3 is available. If not present, attempt to load it from CDN.
async function ensureD3() {
    if (typeof d3 !== 'undefined') return;
    try {
        console.debug('graph: d3 not present, loading dynamically');
        await loadScript('https://d3js.org/d3.v7.min.js');
        console.debug('graph: d3 loaded dynamically');
    } catch (e) {
        console.error('graph: failed to load d3 dynamically', e);
    }
}

async function fetchData() {
    const candidates = ['/static/index.json', './static/index.json', '../static/index.json', '/rendered/static/index.json'];
    for (const path of candidates) {
        try {
            const response = await fetch(path);
            if (!response.ok) { console.debug(`graph: ${path} returned ${response.status}`); continue; }
            const json = await response.json();
            console.debug(`graph: loaded index from ${path}`);
            return json;
        } catch (err) {
            console.debug(`graph: fetch ${path} failed`, err);
            continue;
        }
    }
    console.error('graph: could not load index.json from known locations');
    return null;
}

function resize() {
    if (!container) return;
    width = container.clientWidth;
    height = container.clientHeight;
    if (svg) svg.attr('viewBox', '0 0 ' + width + ' ' + height);
}

// Utilities: convert between HEX and HSL and generate a palette based on the site primary color
function hexToRgb(hex) {
    hex = hex.replace('#', '');
    if (hex.length === 3) hex = hex.split('').map(c => c + c).join('');
    const bigint = parseInt(hex, 16);
    return { r: (bigint >> 16) & 255, g: (bigint >> 8) & 255, b: bigint & 255 };
}

function rgbToHsl(r, g, b) {
    r /= 255; g /= 255; b /= 255;
    const max = Math.max(r, g, b), min = Math.min(r, g, b);
    let h, s, l = (max + min) / 2;
    if (max === min) { h = s = 0; } else {
        const d = max - min;
        s = l > 0.5 ? d / (2 - max - min) : d / (max + min);
        switch (max) {
            case r: h = (g - b) / d + (g < b ? 6 : 0); break;
            case g: h = (b - r) / d + 2; break;
            case b: h = (r - g) / d + 4; break;
        }
        h *= 60;
    }
    return { h: Math.round(h), s: +s.toFixed(3), l: +l.toFixed(3) };
}

function hslToRgb(h, s, l) {
    h /= 360;
    let r, g, b;
    if (s === 0) { r = g = b = l; }
    else {
        const hue2rgb = (p, q, t) => {
            if (t < 0) t += 1;
            if (t > 1) t -= 1;
            if (t < 1 / 6) return p + (q - p) * 6 * t;
            if (t < 1 / 2) return q;
            if (t < 2 / 3) return p + (q - p) * (2 / 3 - t) * 6;
            return p;
        };
        const q = l < 0.5 ? l * (1 + s) : l + s - l * s;
        const p = 2 * l - q;
        r = hue2rgb(p, q, h + 1 / 3);
        g = hue2rgb(p, q, h);
        b = hue2rgb(p, q, h - 1 / 3);
    }
    return { r: Math.round(r * 255), g: Math.round(g * 255), b: Math.round(b * 255) };
}

function rgbToHex({ r, g, b }) {
    return '#' + [r, g, b].map(x => x.toString(16).padStart(2, '0')).join('');
}

function generatePalette(baseHex, n) {
    try {
        const rgb = hexToRgb(baseHex);
        const hsl = rgbToHsl(rgb.r, rgb.g, rgb.b);
        const palette = [];
        // include base color as the first color for visual coherence
        palette.push(baseHex);
        // spread hues evenly around the wheel but keep saturation/lightness similar
        for (let i = 1; i < n; i++) {
            const hue = Math.round((hsl.h + (i * (360 / n))) % 360);
            // vary saturation slightly and keep lightness in readable range
            const sat = Math.max(0.45, Math.min(0.85, hsl.s * (0.9 + ((i % 3) - 1) * 0.03)));
            const light = Math.max(0.32, Math.min(0.7, hsl.l * (0.9 + ((i % 2) ? 0.04 : -0.03))));
            const rgb2 = hslToRgb(hue, sat, light);
            palette.push(rgbToHex(rgb2));
        }
        return palette.slice(0, n);
    } catch (e) {
        // fallback
        return d3.schemeTableau10;
    }
}

function setupSvg() {
    d3.select('#chart').html('');
    container = document.getElementById('chart');
    width = container.clientWidth || 600;
    height = container.clientHeight || 360;

    // derive colors from site CSS variables
    const rootStyles = getComputedStyle(document.documentElement);
    const base = (rootStyles.getPropertyValue('--color-primary') || '#0969da').trim();
    // prefer primary color for the monochrome style
    monoColor = (rootStyles.getPropertyValue('--color-primary') || '#0969da').trim();
    if (monochrome) {
        color = d3.scaleOrdinal([monoColor]);
    } else {
        color = d3.scaleOrdinal(generatePalette(base, 12));
    }

    svg = d3.select('#chart').append('svg')
        .attr('class', 'chart')
        .attr('width', '100%')
        .attr('height', '100%')
        .attr('viewBox', '0 0 ' + width + ' ' + height)
        .attr('preserveAspectRatio', 'xMidYMid meet');

    // wrapper group for zoom/pan
    g = svg.append('g');
    linkGroup = g.append('g').attr('class', 'links');
    nodeGroup = g.append('g').attr('class', 'nodes');
    labelGroup = g.append('g').attr('class', 'labels');

    // zoom & pan
    graphZoom = d3.zoom().scaleExtent([0.2, 6]).on('zoom', (event) => g.attr('transform', event.transform));
    svg.call(graphZoom).on('dblclick.zoom', null);

    window.addEventListener('resize', resize);
}

function buildGraphFromIndex(index) {
    // Compute tag counts and co-occurrence pairs
    const posts = Object.values(index || {});
    const tagCounts = {}; // tag -> count
    const pairs = {}; // 'a||b' -> count

    posts.forEach(post => {
        if (!post.Tags || !Array.isArray(post.Tags)) return;
        // Support both string tags or nested as Frontmatter.Tags depending on index format
        const tagsRaw = Array.isArray(post.Tags) ? post.Tags : (post.Frontmatter && Array.isArray(post.Frontmatter.Tags) ? post.Frontmatter.Tags : []);
        const tags = Array.from(new Set(tagsRaw.map(t => t.trim()).filter(Boolean)));
        tags.forEach(t => tagCounts[t] = (tagCounts[t] || 0) + 1);
        for (let i = 0; i < tags.length; i++) {
            for (let j = i + 1; j < tags.length; j++) {
                const a = tags[i], b = tags[j];
                const key = a < b ? `${a}||${b}` : `${b}||${a}`;
                pairs[key] = (pairs[key] || 0) + 1;
            }
        }
    });

    const nodes = Object.keys(tagCounts).map(id => ({ id, count: tagCounts[id] }));
    const links = Object.keys(pairs).map(k => {
        const [a, b] = k.split('||');
        return { source: a, target: b, value: pairs[k] };
    });

    return { nodes, links };
}

function highlightNeighbors(node, nodesById, links) {
    const neighbors = new Set([node.id]);
    links.forEach(l => {
        if (l.source.id === node.id) neighbors.add(l.target.id);
        if (l.target.id === node.id) neighbors.add(l.source.id);
    });

    nodeGroup.selectAll('circle').classed('node--faded', d => !neighbors.has(d.id));
    labelGroup.selectAll('text').classed('node--faded', d => !neighbors.has(d.id));
    linkGroup.selectAll('line').classed('link--active', d => (d.source.id === node.id || d.target.id === node.id));
}

function clearHighlight() {
    nodeGroup.selectAll('circle').classed('node--faded', false);
    labelGroup.selectAll('text').classed('node--faded', false);
    linkGroup.selectAll('line').classed('link--active', false);
}

function showEmptyMessage(elOrId, text) {
    const el = typeof elOrId === 'string' ? document.getElementById(elOrId) : elOrId;
    if (!el) return;
    // remove existing svg children
    d3.select(el).selectAll('svg').remove();
    // clear and show subtle message
    d3.select(el).selectAll('.graph-empty').remove();
    const msg = document.createElement('div');
    msg.className = 'graph-empty';
    msg.textContent = text;
    el.appendChild(msg);
}

function showLoading(elOrId, text) {
    const el = typeof elOrId === 'string' ? document.getElementById(elOrId) : elOrId;
    if (!el) return;
    d3.select(el).selectAll('.graph-empty, .graph-loading').remove();
    const msg = document.createElement('div');
    msg.className = 'graph-loading';
    msg.textContent = text;
    el.appendChild(msg);
}

function hideLoading(elOrId) {
    const el = typeof elOrId === 'string' ? document.getElementById(elOrId) : elOrId;
    if (!el) return;
    d3.select(el).selectAll('.graph-loading').remove();
}

function draw(graph) {
    if (simulation) simulation.stop();

    const nodes = graph.nodes.map(d => Object.assign({}, d));
    const links = graph.links.map(d => Object.assign({}, d));

    // scales
    // update color palette to match site theme and number of nodes
    const rootStyles = getComputedStyle(document.documentElement);
    const base = (rootStyles.getPropertyValue('--color-primary') || '#0969da').trim();
    const palette = generatePalette(base, Math.max(8, nodes.length));
    color.range(palette).domain(nodes.map(d => d.id)); // deterministic mapping per render

    const sizeScale = d3.scaleSqrt().domain(d3.extent(nodes, d => d.count || 1)).range([6, 18]);
    const linkWidth = d3.scaleLinear().domain(d3.extent(links, d => d.value || 1)).range([0.6, 4]);
    // expose the sizeScale so tick handlers can offset labels
    currentSizeScale = sizeScale;

    // simulation (tighter packing)
    simulation = d3.forceSimulation(nodes)
        .force('link', d3.forceLink(links).id(d => d.id).distance(d => 36 / (d.value || 1)).strength(0.9))
        .force('charge', d3.forceManyBody().strength(-42))
        .force('center', d3.forceCenter(width / 2, height / 2))
        .force('x', d3.forceX(width / 2).strength(0.08))
        .force('y', d3.forceY(height / 2).strength(0.08))
        .force('collision', d3.forceCollide().radius(d => (sizeScale(d.count || 1) + 8)))
        .on('tick', ticked)
        .on('end', () => {
            // fit the graph to view when simulation settles
            try { fitToView(nodes); } catch (e) { /* ignore */ }
        });

    // links
    const link = linkGroup.selectAll('line').data(links, d => `${d.source}-${d.target}`);
    link.join(
        enter => enter.append('line').attr('class', 'link').attr('stroke-width', d => Math.max(0.6, linkWidth(d.value))),
        update => update.attr('stroke-width', d => Math.max(0.6, linkWidth(d.value))),
        exit => exit.remove()
    );

    // nodes
    const node = nodeGroup.selectAll('circle').data(nodes, d => d.id);
    const nodeEnter = node.join(
        enter => enter.append('circle')
            .attr('class', 'node')
            .attr('r', d => sizeScale(d.count || 1))
            .attr('fill', d => monochrome ? monoColor : color(d.id))
            .attr('stroke', 'var(--color-text)')
            .attr('stroke-width', 1.6)
            .call(d3.drag().on('start', dragstarted).on('drag', dragged).on('end', dragended))
            .on('mouseover', function (event, d) { highlightNeighbors(d); })
            .on('mouseout', () => clearHighlight())
            .on('click', function (event, d) {
                // navigate to tag page
                const name = encodeURIComponent(d.id);
                window.location.href = `/tags/${name}.html`;
            }),
        update => update.attr('r', d => sizeScale(d.count || 1)),
        exit => exit.remove()
    );

    nodeEnter.append('title').text(d => `${d.id} (${d.count || 0})`);

    // labels — use CSS colors/sizes and make clickable
    const label = labelGroup.selectAll('text').data(nodes, d => d.id);
    label.join(
        enter => enter.append('text')
            .attr('class', 'node-label')
            .attr('text-anchor', 'middle')
            .style('cursor', 'pointer')
            .text(d => d.id)
            .on('click', function (event, d) { const name = encodeURIComponent(d.id); window.location.href = `/tags/${name}.html`; }),
        update => update.text(d => d.id),
        exit => exit.remove()
    );

    // build index for faster checks
    const nodesById = new Map(nodes.map(d => [d.id, d]));

    simulation.nodes(nodes);
    simulation.force('link').links(links);

    // initial subtle centering and gentle settle for smoother layout
    simulation.alpha(0.6).restart();
}

function fitToView(nodes) {
    if (!svg || !graphZoom) return;
    if (!nodes || nodes.length === 0) return;

    const xs = nodes.map(d => Number.isFinite(d.x) ? d.x : 0);
    const ys = nodes.map(d => Number.isFinite(d.y) ? d.y : 0);
    const minX = Math.min(...xs), maxX = Math.max(...xs);
    const minY = Math.min(...ys), maxY = Math.max(...ys);

    const padding = 40;
    const w = Math.max(1, maxX - minX);
    const h = Math.max(1, maxY - minY);
    const scale = Math.max(0.2, Math.min(6, Math.min((width - padding) / w, (height - padding) / h)));

    const tx = (width / 2) - ((minX + maxX) / 2) * scale;
    const ty = (height / 2) - ((minY + maxY) / 2) * scale;

    const transform = d3.zoomIdentity.translate(tx, ty).scale(scale);
    svg.transition().duration(700).call(graphZoom.transform, transform);
}

function ticked() {
    linkGroup.selectAll('line')
        .attr('x1', d => Number.isFinite(d.source.x) ? d.source.x : width / 2)
        .attr('y1', d => Number.isFinite(d.source.y) ? d.source.y : height / 2)
        .attr('x2', d => Number.isFinite(d.target.x) ? d.target.x : width / 2)
        .attr('y2', d => Number.isFinite(d.target.y) ? d.target.y : height / 2);

    nodeGroup.selectAll('circle')
        .attr('cx', d => {
            const x = Number.isFinite(d.x) ? d.x : width / 2;
            d.x = Math.max(6, Math.min(width - 6, x));
            return d.x;
        })
        .attr('cy', d => {
            const y = Number.isFinite(d.y) ? d.y : height / 2;
            d.y = Math.max(6, Math.min(height - 6, y));
            return d.y;
        });

    labelGroup.selectAll('text')
        .attr('x', d => Number.isFinite(d.x) ? d.x : width / 2)
        .attr('y', d => Number.isFinite(d.y) ? d.y : height / 2);
}

function dragstarted(event) {
    if (!event.active) simulation.alphaTarget(0.3).restart();
    event.subject.fx = event.subject.x;
    event.subject.fy = event.subject.y;
}

function dragged(event) {
    event.subject.fx = event.x;
    event.subject.fy = event.y;
}

function dragended(event) {
    if (!event.active) simulation.alphaTarget(0);
    event.subject.fx = null;
    event.subject.fy = null;
}

async function setupMini(id, indexData) {
    const el = document.getElementById(id);
    if (!el) return;
    showLoading(el, 'Loading graph…');

    // small compact graph (no labels) for sidebar
    const graph = buildGraphFromIndex(indexData);

    // If index has no tags, show friendly message
    if (!graph.nodes || graph.nodes.length === 0) {
        showEmptyMessage(el, 'Graph unavailable — no tags indexed');
        return;
    }

    const w = el.clientWidth || 200;
    const h = el.clientHeight || 140;

    d3.select(el).selectAll('.graph-empty, .graph-loading').remove();

    // ensure mini uses site palette
    const rootStyles = getComputedStyle(document.documentElement);
    const base = (rootStyles.getPropertyValue('--color-primary') || '#0969da').trim();
    const miniPalette = generatePalette(base, Math.max(6, graph.nodes.length));
    color = color || d3.scaleOrdinal(miniPalette).domain(graph.nodes.map(d => d.id));

    const miniSvg = d3.select(el).append('svg')
        .attr('class', 'chart mini')
        .attr('width', '100%')
        .attr('height', '100%')
        .attr('viewBox', '0 0 ' + w + ' ' + h)
        .attr('preserveAspectRatio', 'xMinYMid meet');

    const mg = miniSvg.append('g');
    const linksG = mg.append('g').attr('class', 'links');
    const nodesG = mg.append('g').attr('class', 'nodes');
    const labelG = mg.append('g').attr('class', 'labels');

    // add zoom to mini graph as well
    const miniZoom = d3.zoom().scaleExtent([0.4, 6]).on('zoom', (event) => mg.attr('transform', event.transform));
    miniSvg.call(miniZoom).on('dblclick.zoom', null);

    // ensure nodes have initial positions for mini graph to avoid NaN
    graph.nodes.forEach(n => {
        if (!Number.isFinite(n.x)) n.x = w / 2 + (Math.random() - 0.5) * 40;
        if (!Number.isFinite(n.y)) n.y = h / 2 + (Math.random() - 0.5) * 24;
    });

    // scales for mini graph
    const sizeScale = d3.scaleSqrt().domain(d3.extent(graph.nodes, d => d.count || 1)).range([3, 8]);
    const linkWidth = d3.scaleLinear().domain(d3.extent(graph.links, d => d.value || 1)).range([0.4, 2]);

    const sim = d3.forceSimulation(graph.nodes)
        .force('link', d3.forceLink(graph.links).id(d => d.id).distance(d => 28 / (d.value || 1)).strength(0.9))
        .force('charge', d3.forceManyBody().strength(-28))
        .force('center', d3.forceCenter(w / 2, h / 2))
        .force('x', d3.forceX(w / 2).strength(0.08))
        .force('y', d3.forceY(h / 2).strength(0.08))
        .force('collision', d3.forceCollide().radius(d => sizeScale(d.count || 1) + 6))
        .on('tick', tickMini)
        .on('end', () => { hideLoading(el); });

    const lk = linksG.selectAll('line').data(graph.links).join('line').attr('class', 'link').attr('stroke-width', d => Math.max(0.4, linkWidth(d.value)));
    const nd = nodesG.selectAll('circle').data(graph.nodes).join('circle').attr('class', 'node')
        .attr('r', d => sizeScale(d.count || 1)).attr('fill', d => monochrome ? monoColor : color(d.id)).attr('stroke', 'var(--color-text)').attr('stroke-width', 0.8)
        .on('click', (e, d) => { const name = encodeURIComponent(d.id); window.location.href = `/tags/${name}.html`; });
    // initialize positions for rendered nodes so the first tick isn't NaN
    nodesG.selectAll('circle').each(function (d) { if (!Number.isFinite(d.x)) { d.x = w / 2 + (Math.random() - 0.5) * 40; } if (!Number.isFinite(d.y)) { d.y = h / 2 + (Math.random() - 0.5) * 20; } });

    // labels for mini graph
    const miniLabels = labelG.selectAll('text').data(graph.nodes).join(
        enter => enter.append('text')
            .attr('class', 'node-label')
            .attr('text-anchor', 'middle')
            .attr('font-size', 10)
            .style('pointer-events', 'auto')
            .style('cursor', 'pointer')
            .text(d => d.id)
            .on('click', (e, d) => { const name = encodeURIComponent(d.id); window.location.href = `/tags/${name}.html`; }),
        update => update.text(d => d.id),
        exit => exit.remove()
    );

    function tickMini() {
        lk.attr('x1', d => Number.isFinite(d.source.x) ? d.source.x : w / 2)
            .attr('y1', d => Number.isFinite(d.source.y) ? d.source.y : h / 2)
            .attr('x2', d => Number.isFinite(d.target.x) ? d.target.x : w / 2)
            .attr('y2', d => Number.isFinite(d.target.y) ? d.target.y : h / 2);

        nd.attr('cx', d => { d.x = Math.max(6, Math.min(w - 6, (Number.isFinite(d.x) ? d.x : w / 2))); return d.x; })
            .attr('cy', d => { d.y = Math.max(6, Math.min(h - 6, (Number.isFinite(d.y) ? d.y : h / 2))); return d.y; });

        miniLabels.attr('x', d => d.x)
            .attr('y', d => d.y - (sizeScale(d.count || 1) + LABEL_OFFSET));
    }

    // simple hover highlight
    nd.on('mouseover', function (event, d) {
        lk.classed('link--active', l => l.source.id === d.id || l.target.id === d.id);
        nd.classed('node--faded', n => n.id !== d.id && !graph.links.some(l => (l.source.id === d.id && l.target.id === n.id) || (l.target.id === d.id && l.source.id === n.id)));
    }).on('mouseout', function () {
        lk.classed('link--active', false);
        nd.classed('node--faded', false);
    });
}

async function init() {
    // Ensure D3 is available (load dynamically if necessary)
    await ensureD3();
    if (typeof d3 === 'undefined') {
        console.error('graph: d3 is not available; aborting graph render');
        if (document.getElementById('chart')) showEmptyMessage('chart', 'Graph unavailable — failed to load D3 library');
        if (document.getElementById('mini-graph')) showEmptyMessage('mini-graph', 'Graph unavailable — failed to load D3 library');
        return;
    }

    const data = await fetchData();

    // If index fetch failed entirely, show messages in available graph containers
    if (!data) {
        if (document.getElementById('chart')) showEmptyMessage('chart', 'Graph unavailable — index missing');
        if (document.getElementById('mini-graph')) showEmptyMessage('mini-graph', 'Graph unavailable');
        return;
    }

    const graphData = buildGraphFromIndex(data);
    if (!graphData.nodes || graphData.nodes.length === 0) {
        if (document.getElementById('chart')) showEmptyMessage('chart', 'Graph unavailable — no tags found');
        if (document.getElementById('mini-graph')) showEmptyMessage('mini-graph', 'Graph unavailable');
        return;
    }

    // full chart on tags page
    if (document.getElementById('chart')) {
        setupSvg();
        draw(graphData);
    }

    // mini graph in sidebar
    if (document.getElementById('mini-graph')) {
        // ensure color scale is initialized from site palette (use monochrome if enabled)
        const rootStyles = getComputedStyle(document.documentElement);
        const base = (rootStyles.getPropertyValue('--color-primary') || '#0969da').trim();
        monoColor = (rootStyles.getPropertyValue('--color-primary') || '#0969da').trim();
        const miniPalette = generatePalette(base, Math.max(6, graphData.nodes.length));
        color = color || (monochrome ? d3.scaleOrdinal([monoColor]).domain(graphData.nodes.map(d => d.id)) : d3.scaleOrdinal(miniPalette).domain(graphData.nodes.map(d => d.id)));
        setupMini('mini-graph', data);
    }
}

init();
