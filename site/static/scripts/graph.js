let width, height, color, svg, simulation;

async function fetchData() {
    try {
        const response = await fetch('/static/index.json');
        if (!response.ok) throw new Error('Failed to fetch data');
        return await response.json();
    } catch (error) {
        console.error('Error:', error);
        return null;
    }
}

function setupGraph() {
    d3.select("#chart").html("");
    width = height = 400;
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
    svg.append("g").selectAll("line").data(links).enter().append("line")
        .attr("class", "link") // Apply class for links from parent theme
        .attr("stroke", "var(--color-text-dim)").attr("stroke-opacity", 1)
        .attr("stroke-width", d => Math.sqrt(d.value));
}

function drawNodes(nodes) {
    const node = svg.append("g").selectAll("circle").data(nodes).enter().append("circle")
        .attr("r", 10).attr("class", d => d.group === 0 ? "root-node" : "child-node") // Apply classes for nodes from parent theme
        .attr("fill", d => d.group === 0 ? "var(--color-primary)" : color(d.group))
        .attr("stroke", "var(--color-text)").attr("stroke-width", 2)
        .call(d3.drag().on("start", dragstarted).on("drag", dragged).on("end", dragended))
        .on("click", onClick);

    node.append("title").text(d => {
        //Removing the url prefix and file extension

        if (d.id.indexOf("/") != -1) {
            d.id = d.id.replace(/.*\//, "")
        }
        d.id = d.id.replace(/.html/, "")
        d.id;
    })
}

function drawLabels(nodes) {
    const label = svg.append("g").selectAll("text").data(nodes).enter().append("text")
        .text(d => d.id).style("font-size", "12px")
        .attr("class", "node-label") // Apply class for labels from parent theme
        .attr("dx", 12).attr("dy", ".35em").attr("fill", "var(--color-text)");

    label.append("title").text(d => d.id);
}

function ticked() {
    svg.selectAll("line").attr("x1", d => d.source.x).attr("y1", d => d.source.y)
        .attr("x2", d => d.target.x).attr("y2", d => d.target.y);

    svg.selectAll("circle").attr("cx", d => d.x).attr("cy", d => d.y);
    svg.selectAll("text").attr("x", d => d.x).attr("y", d => d.y);
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
    const data = await fetchData();
    const tag = d.id;
    svg.selectAll("*").remove();

    // fish out common posts based on a tag
    const commonPosts = Object.entries(data)
        .filter(([key, post]) => post.Tags && post.Tags.includes(tag))
        .map(([key, post]) => ({ filename: key, title: post.Frontmatter.Title }));

    if (commonPosts.length === 0 && d.group !== 0) {
        nodeName = "posts/" + d.id + ".html" // Adding url prefix and extension for posts
        window.location.href = `/${nodeName}`; // If it's a leaf node, automatically redirect
    }
    else {
        setupGraph();
        drawGraph(commonPosts, tag);
    }
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
    let nodes;

    if (Array.isArray(data)) {
        links.push(...data.map(post => ({
            source: str, target: post.filename, id: `${str}-${post.filename}`, value: 1
        })));
        nodes = [{ id: str, group: 0 }, ...data.map(post => ({ id: post.filename, group: 1 }))];
    } else {
        const nodesMap = {};
        Object.keys(data).forEach(key => {
            const post = data[key];
            const tags = post.Tags;
            if (tags) {
                tags.forEach(tag => {
                    nodesMap[tag] = { id: tag, group: 1 };
                    const linkId = `Center-${tag}`;
                    const existingLink = links.find(d => d.id === linkId);
                    if (!existingLink) links.push({ source: str, target: tag, id: linkId, value: 1 });
                    else existingLink.value++;
                });
            }
        });
        nodesMap[str] = { id: str, group: 0 };
        nodes = Object.values(nodesMap);
    }

    svg = d3.select("#chart").append("svg").attr("class", "chart").style("font-family", "agave, monospace").attr("width", width).attr("height", height);
    setupSimulation(nodes, links);
    drawLinks(links);
    drawNodes(nodes);
    drawLabels(nodes);
}

init();
