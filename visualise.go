package catrace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"
)

// VisualiseOptions controls the HTML graph output.
type VisualiseOptions struct {
	// Title is shown in the browser tab and as the graph heading.
	// Defaults to the kernel's first state name or "Markov Kernel".
	Title string

	// MinEdge omits directed edges whose transition probability is below this
	// threshold. Reduces visual clutter on dense kernels. Default: 0.01.
	MinEdge float64

	// Width and Height of the SVG canvas in pixels. Defaults: 960 × 680.
	Width  int
	Height int

	// StationaryTol and StationaryMaxIter are passed to Stationary.
	// If stationary computation fails (e.g. reducible chain), all nodes
	// are rendered at equal size.
	StationaryTol     float64
	StationaryMaxIter int
}

func (o *VisualiseOptions) withDefaults() VisualiseOptions {
	out := VisualiseOptions{
		Title:             "Markov Kernel",
		MinEdge:           0.01,
		Width:             960,
		Height:            680,
		StationaryTol:     1e-12,
		StationaryMaxIter: 5000,
	}
	if o == nil {
		return out
	}
	if o.Title != "" {
		out.Title = o.Title
	}
	if o.MinEdge > 0 {
		out.MinEdge = o.MinEdge
	}
	if o.Width > 0 {
		out.Width = o.Width
	}
	if o.Height > 0 {
		out.Height = o.Height
	}
	if o.StationaryTol > 0 {
		out.StationaryTol = o.StationaryTol
	}
	if o.StationaryMaxIter > 0 {
		out.StationaryMaxIter = o.StationaryMaxIter
	}
	return out
}

// gNode is the JSON representation of a graph node.
type gNode struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Pi    float64 `json:"pi"`  // stationary probability (0 if unavailable)
	Class int     `json:"cls"` // communicating class index; -1 = transient
}

// gEdge is the JSON representation of a directed edge.
type gEdge struct {
	Source int     `json:"source"`
	Target int     `json:"target"`
	Value  float64 `json:"value"` // transition probability
}

// ToHTML generates a self-contained HTML file that renders the kernel as an
// interactive force-directed graph, similar to Obsidian's graph view.
//
// Node radius is proportional to the stationary distribution π. Edge stroke
// width and opacity are proportional to the transition probability. Directed
// edges are drawn as curved arcs with arrowheads so that A→B and B→A are
// visually distinct. Self-loops are rendered as small arcs above the node.
// Nodes are coloured by communicating class; transient nodes are grey.
//
// The returned bytes are a self-contained UTF-8 HTML document. Write them to
// a .html file and open in any modern browser. D3.js v7 is loaded from CDN;
// an internet connection is required on first open (after that, browser cache
// handles it).
//
// Example:
//
//	html, err := kernel.ToHTML(&catrace.VisualiseOptions{Title: "Q kernel"})
//	if err != nil { log.Fatal(err) }
//	os.WriteFile("graph.html", html, 0644)
func (k *Kernel) ToHTML(opts *VisualiseOptions) ([]byte, error) {
	if k == nil || k.P == nil {
		return nil, fmt.Errorf("nil kernel")
	}
	o := opts.withDefaults()
	n := k.NumStates()

	// --- stationary distribution ---
	pi := make([]float64, n)
	if stat, err := k.Stationary(o.StationaryTol, o.StationaryMaxIter); err == nil {
		copy(pi, stat)
	} else {
		// fallback: uniform
		for i := range pi {
			pi[i] = 1.0 / float64(n)
		}
	}

	// --- communicating classes ---
	classIdx := make([]int, n) // -1 = transient
	for i := range classIdx {
		classIdx[i] = -1
	}
	if cd, err := k.Classes(1e-9); err == nil {
		for ri, comp := range cd.Recurrent {
			for _, s := range comp {
				classIdx[s] = ri
			}
		}
	}

	// --- build graph data ---
	nodes := make([]gNode, n)
	for i := 0; i < n; i++ {
		nodes[i] = gNode{
			ID:    i,
			Name:  k.StateNames[i],
			Pi:    pi[i],
			Class: classIdx[i],
		}
	}

	var edges []gEdge
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			v := k.P.At(i, j)
			if v >= o.MinEdge {
				edges = append(edges, gEdge{Source: i, Target: j, Value: v})
			}
		}
	}

	nodesJSON, err := json.Marshal(nodes)
	if err != nil {
		return nil, fmt.Errorf("marshalling nodes: %w", err)
	}
	edgesJSON, err := json.Marshal(edges)
	if err != nil {
		return nil, fmt.Errorf("marshalling edges: %w", err)
	}

	tmpl, err := template.New("graph").Parse(htmlTemplate)
	if err != nil {
		return nil, fmt.Errorf("parsing template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, map[string]any{
		"Title":     o.Title,
		"Width":     o.Width,
		"Height":    o.Height,
		"NodesJSON": string(nodesJSON),
		"EdgesJSON": string(edgesJSON),
	})
	if err != nil {
		return nil, fmt.Errorf("executing template: %w", err)
	}
	return buf.Bytes(), nil
}

// htmlTemplate is the self-contained D3.js v7 force-directed graph template.
// Template variables: Title, Width, Height, NodesJSON, EdgesJSON.
const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<title>{{.Title}}</title>
<style>
  * { box-sizing: border-box; margin: 0; padding: 0; }
  body { background: #1a1a2e; color: #e0e0e0; font-family: "Inter", "Segoe UI", sans-serif; overflow: hidden; }
  #title { position: absolute; top: 14px; left: 18px; font-size: 15px; font-weight: 600; color: #a0a8c0; letter-spacing: 0.04em; }
  #legend { position: absolute; bottom: 14px; left: 18px; font-size: 11px; color: #666; line-height: 1.7; }
  #tooltip {
    position: absolute; pointer-events: none; display: none;
    background: rgba(20,20,40,0.93); border: 1px solid #334;
    border-radius: 6px; padding: 7px 11px; font-size: 12px;
    color: #ccd; line-height: 1.6; white-space: pre;
  }
  .node circle { stroke: #1a1a2e; stroke-width: 1.5px; cursor: grab; }
  .node circle:active { cursor: grabbing; }
  .node text { font-size: 11px; fill: #dde; pointer-events: none; text-anchor: middle; dominant-baseline: central; }
  .link { fill: none; }
  .self-loop { fill: none; }
</style>
</head>
<body>
<div id="title">{{.Title}}</div>
<div id="legend">
  Node radius ∝ stationary distribution π<br>
  Edge width &amp; opacity ∝ transition probability<br>
  Scroll to zoom · drag nodes · hover for values
</div>
<div id="tooltip"></div>
<svg id="svg" width="{{.Width}}" height="{{.Height}}"></svg>

<script src="https://d3js.org/d3.v7.min.js"></script>
<script>
const W = {{.Width}}, H = {{.Height}};
const rawNodes = {{.NodesJSON}};
const rawEdges = {{.EdgesJSON}};

// Class colour palette (Obsidian-adjacent)
const palette = [
  "#7c6af7","#4ecdc4","#f7b731","#fc5c65","#45aaf2",
  "#a55eea","#26de81","#fd9644","#eb3b5a","#2bcbba"
];
const transientColour = "#4a4a6a";

function nodeColour(d) {
  if (d.cls < 0) return transientColour;
  return palette[d.cls % palette.length];
}

// Node radius: map pi to [8, 28]
const piVals = rawNodes.map(d => d.pi);
const piMin = Math.min(...piVals), piMax = Math.max(...piVals);
const piRange = piMax - piMin || 1;
function nodeR(pi) { return 8 + 20 * (pi - piMin) / piRange; }

// Self-loops and regular edges separated
const selfLoops = rawEdges.filter(e => e.source === e.target);
const dirEdges  = rawEdges.filter(e => e.source !== e.target);

const svg = d3.select("#svg");
const g   = svg.append("g");

// Zoom/pan
svg.call(d3.zoom().scaleExtent([0.15, 5]).on("zoom", ev => g.attr("transform", ev.transform)));

// Arrow markers — one per opacity bucket to avoid re-coloring SVG markers
const markerColour = "#8899bb";
g.append("defs").append("marker")
  .attr("id", "arrow")
  .attr("viewBox", "0 -4 8 8")
  .attr("refX", 8).attr("refY", 0)
  .attr("markerWidth", 6).attr("markerHeight", 6)
  .attr("orient", "auto")
  .append("path")
    .attr("d", "M0,-4L8,0L0,4")
    .attr("fill", markerColour);

// Edge layer
const linkLayer = g.append("g").attr("class", "links");
const loopLayer = g.append("g").attr("class", "loops");
const nodeLayer = g.append("g").attr("class", "nodes");

const tooltip = d3.select("#tooltip");

// Nodes map for lookup
const nodeMap = Object.fromEntries(rawNodes.map(d => [d.id, d]));

// --- Draw directed edges ---
const link = linkLayer.selectAll("path")
  .data(dirEdges)
  .join("path")
    .attr("class", "link")
    .attr("stroke", markerColour)
    .attr("stroke-width", d => 0.8 + 5 * d.value)
    .attr("stroke-opacity", d => 0.25 + 0.65 * d.value)
    .attr("marker-end", "url(#arrow)")
    .on("mouseover", (ev, d) => {
      const sn = nodeMap[d.source]?.name ?? d.source;
      const tn = nodeMap[d.target]?.name ?? d.target;
      tooltip.style("display","block")
        .text("P(" + sn + " → " + tn + ") = " + d.value.toFixed(4));
    })
    .on("mousemove", ev => tooltip.style("left",(ev.pageX+14)+"px").style("top",(ev.pageY-28)+"px"))
    .on("mouseout", () => tooltip.style("display","none"));

// --- Draw self-loops ---
const loop = loopLayer.selectAll("path")
  .data(selfLoops)
  .join("path")
    .attr("class", "self-loop")
    .attr("stroke", markerColour)
    .attr("stroke-width", d => 0.8 + 4 * d.value)
    .attr("stroke-opacity", d => 0.25 + 0.65 * d.value)
    .on("mouseover", (ev, d) => {
      const sn = nodeMap[d.source]?.name ?? d.source;
      tooltip.style("display","block")
        .text("P(" + sn + " → " + sn + ") = " + d.value.toFixed(4));
    })
    .on("mousemove", ev => tooltip.style("left",(ev.pageX+14)+"px").style("top",(ev.pageY-28)+"px"))
    .on("mouseout", () => tooltip.style("display","none"));

// --- Draw nodes ---
const node = nodeLayer.selectAll("g")
  .data(rawNodes)
  .join("g")
    .attr("class","node")
    .call(d3.drag()
      .on("start", (ev,d) => { if (!ev.active) sim.alphaTarget(0.3).restart(); d.fx=d.x; d.fy=d.y; })
      .on("drag",  (ev,d) => { d.fx=ev.x; d.fy=ev.y; })
      .on("end",   (ev,d) => { if (!ev.active) sim.alphaTarget(0); d.fx=null; d.fy=null; }));

node.append("circle")
  .attr("r", d => nodeR(d.pi))
  .attr("fill", nodeColour)
  .on("mouseover", (ev, d) => {
    const outEdges = dirEdges.filter(e => e.source === d.id);
    const lines = ["State: " + d.name, "π = " + d.pi.toFixed(5), "out-edges: " + outEdges.length];
    tooltip.style("display","block").text(lines.join("\n"));
  })
  .on("mousemove", ev => tooltip.style("left",(ev.pageX+14)+"px").style("top",(ev.pageY-28)+"px"))
  .on("mouseout", () => tooltip.style("display","none"));

node.append("text")
  .attr("dy", d => nodeR(d.pi) + 13)
  .text(d => d.name);

// --- Force simulation ---
const sim = d3.forceSimulation(rawNodes)
  .force("link", d3.forceLink(dirEdges)
    .id(d => d.id)
    .distance(d => 80 + 120 * (1 - d.value))
    .strength(d => 0.3 + 0.5 * d.value))
  .force("charge", d3.forceManyBody().strength(-350))
  .force("center", d3.forceCenter(W/2, H/2))
  .force("collision", d3.forceCollide().radius(d => nodeR(d.pi) + 12))
  .on("tick", ticked);

// Path generator for curved directed edges
// Two edges A→B and B→A get opposite curvature; single-direction edges get slight curvature.
const edgePairs = new Set(
  dirEdges.map(e => [Math.min(e.source,e.target), Math.max(e.source,e.target)].join(","))
);
function isBidirectional(e) {
  return dirEdges.some(f => f.source === e.target && f.target === e.source);
}

function edgePath(d) {
  const s = rawNodes[d.source], t = rawNodes[d.target];
  const dx = t.x - s.x, dy = t.y - s.y;
  const dist = Math.sqrt(dx*dx + dy*dy) || 1;
  const tr = nodeR(t.pi) + 7; // target radius + arrow gap
  // unit vector
  const ux = dx/dist, uy = dy/dist;
  // end point: just outside target circle
  const ex = t.x - ux*tr, ey = t.y - uy*tr;
  // start point: just outside source circle
  const sr = nodeR(s.pi);
  const sx = s.x + ux*sr, sy = s.y + uy*sr;
  // curvature: bidirectional edges curve opposite ways
  if (isBidirectional(d)) {
    const curve = dist * 0.25;
    // perpendicular
    const px = -uy*curve, py = ux*curve;
    const sign = d.source < d.target ? 1 : -1;
    const cx = (sx+ex)/2 + sign*px, cy = (sy+ey)/2 + sign*py;
    return "M"+sx+","+sy+" Q"+cx+","+cy+" "+ex+","+ey;
  }
  // single direction: very slight curve to avoid overlapping node labels
  const curve = dist * 0.08;
  const px = -uy*curve, py = ux*curve;
  const cx = (sx+ex)/2 + px, cy = (sy+ey)/2 + py;
  return "M"+sx+","+sy+" Q"+cx+","+cy+" "+ex+","+ey;
}

function selfLoopPath(d) {
  const n = rawNodes[d.source];
  const r = nodeR(n.pi);
  const lx = n.x, ly = n.y - r;
  const lw = r * 1.6;
  return "M"+(lx-lw/2)+","+(ly)+" C"+(lx-lw)+","+(ly-lw*1.2)+" "+(lx+lw)+","+(ly-lw*1.2)+" "+(lx+lw/2)+","+ly;
}

function ticked() {
  link.attr("d", edgePath);
  loop.attr("d", selfLoopPath);
  node.attr("transform", d => "translate("+d.x+","+d.y+")");
}
</script>
</body>
</html>`
