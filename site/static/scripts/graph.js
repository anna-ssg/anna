let width;
let height;
let color;
let svg;
let simulation;

async function fetchData() {
    try {
        const response = await fetch('/static/index.json');
        if (!response.ok) { throw new Error('Failed to fetch data'); }
        return await response.json();
    } catch (error) {
        console.error('Error:', error);
        return null;
    }
}

function setupGraph() {
    d3.select("#chart").html("");
    width = 400;
    height = 400;
    color = d3.scaleOrdinal(d3.schemeCategory10);
}

function setupSimulation(nodes, links) {
    simulation = d3.forceSimulation(nodes)
        .force("link", d3.forceLink(links).id(d => d.id).distance(50))
        .force("charge", d3.forceManyBody().strength(-300))
        .force("center", d3.forceCenter(width / 2, height / 2))
        .on("tick", ticked);
}

function drawLinks(links) {
    svg.append("g")
        .selectAll("line")
        .data(links)
        .enter().append("line")
        .attr("stroke", "#999")
        .attr("stroke-opacity", 1)
        .attr("stroke-width", d => Math.sqrt(d.value));
}

function drawNodes(nodes) {
    const node = svg.append("g")
        .selectAll("circle")
        .data(nodes)
        .enter().append("circle")
        .attr("r", 10)
        .attr("fill", d => d.group === 0 ? "red" : color(d.group))
        .attr("stroke", "#fff")
        .attr("stroke-width", 2)
        .call(d3.drag()
            .on("start", dragstarted)
            .on("drag", dragged)
            .on("end", dragended)
        )
        .on("click", onClick);

    node.append("title")
        .text(d => d.id);
}

function drawLabels(nodes) {
    const label = svg.append("g")
        .selectAll("text")
        .data(nodes)
        .enter().append("text")
        .text(d => d.id)
        .style("font-size", "12px")
        .attr("dx", 12)
        .attr("dy", ".35em")
        .attr("fill", "white");

    label.append("title")
        .text(d => d.id);
}

function ticked() {
    svg.selectAll("line")
        .attr("x1", d => d.source.x)
        .attr("y1", d => d.source.y)
        .attr("x2", d => d.target.x)
        .attr("y2", d => d.target.y);

    svg.selectAll("circle")
        .attr("cx", d => d.x)
        .attr("cy", d => d.y);

    svg.selectAll("text")
        .attr("x", d => d.x)
        .attr("y", d => d.y);
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

async function onClick(event, d) {
    svg.selectAll("*").remove();
    const url = `/tags/${d.id}.html`;
    window.location.href = url;
}

async function init() {
    const data = await fetchData();
    if (data) {
        setupGraph();
        drawGraph(data, "index");
    }
}

async function drawGraph(data, str) {
    const links = [];
    const nodesMap = {};

    // Iterate through each post to create nodes and links
    Object.keys(data).forEach(key => {
        const post = data[key];
        const tags = post.Tags;
        if (tags) {
            tags.forEach(tag => {
                nodesMap[tag] = { id: tag, group: 1 };
                const linkId = `Center-${tag}`;
                const existingLink = links.find(d => d.id === linkId);
                if (!existingLink) {
                    links.push({ source: str, target: tag, id: linkId, value: 1 });
                } else {
                    existingLink.value++;
                }
            });
        }
    });

    nodesMap[str] = { id: str, group: 0 };
    const nodes = Object.values(nodesMap);

    svg = d3.select("#chart").append("svg")
        .attr("width", width)
        .attr("height", height);

    setupSimulation(nodes, links);
    drawLinks(links);
    drawNodes(nodes);
    drawLabels(nodes);
}

init();
