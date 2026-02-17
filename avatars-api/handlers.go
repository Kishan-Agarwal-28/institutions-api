package main

import (
	"fmt"
	"net/http"
	"strings"
)

// documentationHandler serves the API documentation HTML page
func documentationHandler(w http.ResponseWriter, r *http.Request) {
	const apiDocsHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Avatar Generator API - Documentation</title>
	<style>
		:root {
			--bg-color: #f8f9fa;
			--text-color: #212529;
			--accent-color: #007bff;
			--code-bg: #e9ecef;
			--card-bg: #ffffff;
			--border-color: #dee2e6;
			--shadow: 0 4px 6px rgba(0,0,0,0.1);
		}
		@media (prefers-color-scheme: dark) {
			:root {
				--bg-color: #0d1117;
				--text-color: #e6edf3;
				--accent-color: #58a6ff;
				--code-bg: #161b22;
				--card-bg: #161b22;
				--border-color: #30363d;
				--shadow: 0 4px 6px rgba(0,0,0,0.4);
			}
		}
		* { margin: 0; padding: 0; box-sizing: border-box; }
		body { 
			font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; 
			line-height: 1.6; 
			color: var(--text-color); 
			background: var(--bg-color); 
		}
		.header {
			background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
			color: white;
			padding: 3rem 2rem;
			text-align: center;
		}
		.header h1 { font-size: 2.5rem; margin-bottom: 0.5rem; }
		.header p { font-size: 1.2rem; opacity: 0.9; }
		.container { 
			max-width: 1200px; 
			margin: 0 auto; 
			padding: 2rem; 
		}
		.section { 
			background: var(--card-bg); 
			padding: 2rem; 
			border-radius: 8px; 
			box-shadow: var(--shadow); 
			margin-bottom: 2rem;
			border: 1px solid var(--border-color);
		}
		h2 { 
			color: var(--accent-color); 
			margin-bottom: 1rem; 
			padding-bottom: 0.5rem;
			border-bottom: 2px solid var(--accent-color);
		}
		h3 { margin: 1.5rem 0 1rem; color: var(--text-color); }
		.endpoint { 
			background: var(--code-bg); 
			padding: 1rem; 
			border-radius: 5px; 
			font-family: monospace; 
			font-size: 1rem; 
			margin: 1rem 0;
			border: 1px solid var(--border-color);
			overflow-x: auto;
		}
		.method { 
			color: #fff; 
			background: #28a745; 
			padding: 4px 10px; 
			border-radius: 4px; 
			margin-right: 10px; 
			font-weight: bold; 
			font-size: 0.9rem;
		}
		code { 
			background: var(--code-bg); 
			padding: 3px 6px; 
			border-radius: 4px; 
			font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, monospace; 
			font-size: 0.9em;
		}
		pre { 
			background: var(--code-bg); 
			padding: 1rem; 
			border-radius: 5px; 
			overflow-x: auto; 
			border: 1px solid var(--border-color);
			margin: 1rem 0;
		}
		table { 
			width: 100%; 
			border-collapse: collapse; 
			margin: 1rem 0; 
		}
		th, td { 
			text-align: left; 
			padding: 12px; 
			border-bottom: 1px solid var(--border-color); 
		}
		th { 
			background: var(--code-bg); 
			font-weight: 600;
		}
		.badge { 
			display: inline-block; 
			padding: 3px 10px; 
			border-radius: 12px; 
			font-size: 0.85em; 
			font-weight: bold; 
		}
		.required { background: #dc3545; color: white; }
		.optional { background: #6c757d; color: white; }
		.avatar-grid {
			display: grid;
			grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
			gap: 1.5rem;
			margin: 2rem 0;
		}
		.avatar-card {
			background: var(--card-bg);
			border: 1px solid var(--border-color);
			border-radius: 8px;
			padding: 1rem;
			text-align: center;
			transition: transform 0.2s, box-shadow 0.2s;
		}
		.avatar-card:hover {
			transform: translateY(-4px);
			box-shadow: 0 8px 12px rgba(0,0,0,0.15);
		}
		.avatar-card img {
			width: 150px;
			height: 150px;
			border-radius: 8px;
			margin-bottom: 0.5rem;
			background: var(--code-bg);
		}
		.avatar-card h4 {
			margin: 0.5rem 0;
			color: var(--accent-color);
		}
		.avatar-card code {
			font-size: 0.8rem;
			word-break: break-all;
		}
		.footer { 
			text-align: center; 
			padding: 2rem; 
			color: #6c757d; 
			font-size: 0.9rem;
		}
		@media (max-width: 768px) {
			.header h1 { font-size: 2rem; }
			.header p { font-size: 1rem; }
			.container { padding: 1rem; }
			.section { padding: 1rem; }
			.avatar-grid {
				grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
				gap: 1rem;
			}
			.avatar-card img {
				width: 120px;
				height: 120px;
			}
		}
	</style>
</head>
<body>
	<div class="header">
		<h1>🎨 Avatar Generator API</h1>
		<p>Generate beautiful, unique avatars with a simple HTTP request</p>
	</div>

	<div class="container">
		<div class="section">
			<h2>📖 Overview</h2>
			<p>A high-performance JSON API that generates unique, deterministic SVG avatars based on a name input. Perfect for user profiles, comments, and any application needing consistent, personalized avatars.</p>
		</div>

		<div class="section">
			<h2>🚀 Quick Start</h2>
			<h3>Base URL</h3>
			<div class="endpoint">
				<span class="method">GET</span> /avatars/api/generate-avatar
			</div>

			<h3>Example Request</h3>
			<pre><code>GET /avatars/api/generate-avatar?name=John%20Doe&type=avatar&size=200</code></pre>

			<h3>Parameters</h3>
			<table>
				<thead>
					<tr>
						<th>Parameter</th>
						<th>Type</th>
						<th>Required</th>
						<th>Default</th>
						<th>Description</th>
					</tr>
				</thead>
				<tbody>
					<tr>
						<td><code>name</code></td>
						<td>String</td>
						<td><span class="badge optional">Optional</span></td>
						<td>"User"</td>
						<td>Name used to generate the avatar. Same name always produces the same avatar.</td>
					</tr>
					<tr>
						<td><code>type</code></td>
						<td>String</td>
						<td><span class="badge optional">Optional</span></td>
						<td>"avatar"</td>
						<td>Avatar style. See available types below.</td>
					</tr>
					<tr>
						<td><code>size</code></td>
						<td>Integer</td>
						<td><span class="badge optional">Optional</span></td>
						<td>100</td>
						<td>Size of the avatar in pixels (width and height).</td>
					</tr>
					<tr>
						<td><code>color</code></td>
						<td>String</td>
						<td><span class="badge optional">Optional</span></td>
						<td>Auto</td>
						<td>Hex color code (e.g., #FF5733). If not provided, color is generated from name.</td>
					</tr>
				</tbody>
			</table>
		</div>

		<div class="section">
			<h2>🎭 Avatar Types</h2>
			<p>All examples use the name "Alex Morgan" for consistency. Each type generates a unique, deterministic design.</p>
			
			<div class="avatar-grid">
				<div class="avatar-card">
					<img src="/avatars/api/generate-avatar?name=Alex%20Morgan&type=avatar&size=150" alt="Avatar">
					<h4>avatar</h4>
					<p>Classic circular avatar with initials</p>
					<code>type=avatar</code>
				</div>

				<div class="avatar-card">
					<img src="/avatars/api/generate-avatar?name=Alex%20Morgan&type=gravatar&size=150" alt="Gravatar">
					<h4>gravatar</h4>
					<p>GitHub-style identicon</p>
					<code>type=gravatar</code>
				</div>

				<div class="avatar-card">
					<img src="/avatars/api/generate-avatar?name=Alex%20Morgan&type=dither&size=150" alt="Dither">
					<h4>dither</h4>
					<p>Retro dithered plasma effect</p>
					<code>type=dither</code>
				</div>

				<div class="avatar-card">
					<img src="/avatars/api/generate-avatar?name=Alex%20Morgan&type=ascii&size=150" alt="ASCII">
					<h4>ascii</h4>
					<p>Procedural ASCII robot</p>
					<code>type=ascii</code>
				</div>

				<div class="avatar-card">
					<img src="/avatars/api/generate-avatar?name=Alex%20Morgan&type=dotmatrix&size=150" alt="Dot Matrix">
					<h4>dotmatrix</h4>
					<p>LED dot matrix display</p>
					<code>type=dotmatrix</code>
				</div>

				<div class="avatar-card">
					<img src="/avatars/api/generate-avatar?name=Alex%20Morgan&type=terminal&size=150" alt="Terminal">
					<h4>terminal</h4>
					<p>Retro terminal block text</p>
					<code>type=terminal</code>
				</div>

				<div class="avatar-card">
					<img src="/avatars/api/generate-avatar?name=Alex%20Morgan&type=bauhaus&size=150" alt="Bauhaus">
					<h4>bauhaus</h4>
					<p>Geometric Bauhaus design</p>
					<code>type=bauhaus</code>
				</div>

				<div class="avatar-card">
					<img src="/avatars/api/generate-avatar?name=Alex%20Morgan&type=ring&size=150" alt="Ring">
					<h4>ring</h4>
					<p>Gradient ring pattern</p>
					<code>type=ring</code>
				</div>

				<div class="avatar-card">
					<img src="/avatars/api/generate-avatar?name=Alex%20Morgan&type=beam&size=150" alt="Beam">
					<h4>beam</h4>
					<p>Connected network nodes</p>
					<code>type=beam</code>
				</div>

				<div class="avatar-card">
					<img src="/avatars/api/generate-avatar?name=Alex%20Morgan&type=marble&size=150" alt="Marble">
					<h4>marble</h4>
					<p>Marble texture effect</p>
					<code>type=marble</code>
				</div>

				<div class="avatar-card">
					<img src="/avatars/api/generate-avatar?name=Alex%20Morgan&type=glitch&size=150" alt="Glitch">
					<h4>glitch</h4>
					<p>Cyberpunk glitch effect</p>
					<code>type=glitch</code>
				</div>

				<div class="avatar-card">
					<img src="/avatars/api/generate-avatar?name=Alex%20Morgan&type=sunset&size=150" alt="Sunset">
					<h4>sunset</h4>
					<p>Procedural sunset scene</p>
					<code>type=sunset</code>
				</div>

				<div class="avatar-card">
					<img src="/avatars/api/generate-avatar?name=Alex%20Morgan&type=smile&size=150" alt="Smile">
					<h4>smile</h4>
					<p>Minimalist face with expressions</p>
					<code>type=smile</code>
				</div>

				<div class="avatar-card">
					<img src="/avatars/api/generate-avatar?name=Alex%20Morgan&type=circuit&size=150" alt="Circuit">
					<h4>circuit</h4>
					<p>Circuit board pattern</p>
					<code>type=circuit</code>
				</div>

				<div class="avatar-card">
					<img src="/avatars/api/generate-avatar?name=Alex%20Morgan&type=pixel&size=150" alt="Pixel">
					<h4>pixel</h4>
					<p>Isometric pixel art cube</p>
					<code>type=pixel</code>
				</div>

				<div class="avatar-card">
					<img src="/avatars/api/generate-avatar?name=Alex%20Morgan&type=constellation&size=150" alt="Constellation">
					<h4>constellation</h4>
					<p>Real constellation star maps</p>
					<code>type=constellation</code>
				</div>
			</div>
		</div>

		<div class="section">
			<h2>💡 Usage Examples</h2>
			
			<h3>HTML Image Tag</h3>
			<pre><code>&lt;img src="/avatars/api/generate-avatar?name=Jane%20Smith&type=avatar&size=100" alt="Avatar"&gt;</code></pre>

			<h3>With Custom Color</h3>
			<pre><code>&lt;img src="/avatars/api/generate-avatar?name=John&type=gravatar&color=%23FF5733" alt="Avatar"&gt;</code></pre>

			<h3>Large Size</h3>
			<pre><code>&lt;img src="/avatars/api/generate-avatar?name=Sarah&type=dotmatrix&size=500" alt="Avatar"&gt;</code></pre>

			<h3>JavaScript Fetch</h3>
			<pre><code>fetch('/avatars/api/generate-avatar?name=Bob%20Johnson&type=beam')
  .then(response => response.text())
  .then(svg => {
    document.getElementById('avatar').innerHTML = svg;
  });</code></pre>
		</div>

		<div class="section">
			<h2>⚡ Features</h2>
			<ul style="list-style: none; padding: 0;">
				<li style="padding: 0.5rem 0;">✅ <strong>Deterministic</strong> - Same name always generates the same avatar</li>
				<li style="padding: 0.5rem 0;">✅ <strong>SVG Format</strong> - Scalable to any size without quality loss</li>
				<li style="padding: 0.5rem 0;">✅ <strong>16 Unique Styles</strong> - From classic to creative designs</li>
				<li style="padding: 0.5rem 0;">✅ <strong>Fast & Lightweight</strong> - Generated on-the-fly, no storage needed</li>
				<li style="padding: 0.5rem 0;">✅ <strong>CORS Enabled</strong> - Use from any domain</li>
				<li style="padding: 0.5rem 0;">✅ <strong>Rate Limited</strong> - 100 requests per minute per IP</li>
				<li style="padding: 0.5rem 0;">✅ <strong>Dark Mode Support</strong> - Responsive documentation</li>
			</ul>
		</div>

		<div class="section">
			<h2>📝 Response Format</h2>
			<p>All avatars are returned as SVG (Scalable Vector Graphics) with the content type:</p>
			<div class="endpoint">Content-Type: image/svg+xml; charset=utf-8</div>
			
			<h3>Status Codes</h3>
			<table>
				<thead>
					<tr>
						<th>Code</th>
						<th>Description</th>
					</tr>
				</thead>
				<tbody>
					<tr>
						<td><code>200 OK</code></td>
						<td>Request successful, avatar returned</td>
					</tr>
					<tr>
						<td><code>400 Bad Request</code></td>
						<td>Invalid avatar type specified</td>
					</tr>
					<tr>
						<td><code>429 Too Many Requests</code></td>
						<td>Rate limit exceeded (100 req/min)</td>
					</tr>
					<tr>
						<td><code>500 Internal Server Error</code></td>
						<td>Server-side error occurred</td>
					</tr>
				</tbody>
			</table>
		</div>

		<div class="section">
			<h2>🎨 Color Generation</h2>
			<p>When no color is specified, avatars use the <strong>OKLCH color space</strong> for perceptually uniform, accessible colors. The algorithm:</p>
			<ul style="margin-left: 2rem; margin-top: 1rem;">
				<li>Generates a deterministic hash from the input name</li>
				<li>Maps hash values to hue, chroma, and lightness</li>
				<li>Ensures sufficient contrast for readability</li>
				<li>Produces consistent colors across all sessions</li>
			</ul>
		</div>
	</div>

	<div class="footer">
		<p>Avatar Generator API v1.0 • Built with Go & SVG • MIT License</p>
	</div>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(apiDocsHTML))
}

// generateAvatarHandler handles avatar generation requests
func generateAvatarHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	avatarType := query.Get("type")
	size := query.Get("size")
	color := query.Get("color")
	name := query.Get("name")
	if size == "" {
		size = "100"
	}
	if avatarType == "" {
		avatarType = "avatar"
	}
	if name == "" {
		name = "User"
	}

	if color == "" {
		color, _ = generateColor(name)
	}

	cacheKey := fmt.Sprintf("%s:%s:%s:%s", name, avatarType, size, color)

	if cache != nil {
		if cachedSVG, found := cache.Get(cacheKey); found {
			w.Header().Set("X-Cache", "HIT")
			w.Header().Set("Cache-Control", "public, max-age=86400")
			w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(cachedSVG))
			return
		}
	}

	var avatarContent string

	switch avatarType {
	case "avatar":
		initials := getInitials(name)
		avatarContent = fmt.Sprintf(`
		<svg xmlns="http://www.w3.org/2000/svg" width="%[1]s" height="%[1]s" viewBox="0 0 100 100">
			<g transform="translate(5, 5) scale(0.9)">
				<circle cx="50" cy="50" r="50" fill="%[2]s" />
				<text x="50" y="55" dominant-baseline="middle" text-anchor="middle" font-family="Arial, sans-serif" font-size="40" fill="#ffffff">%[3]s</text>
			</g>
		</svg>`, size, color, initials)
	case "gravatar":
		identiconRects := generateIdenticon(name)
		avatarContent = fmt.Sprintf(`
		<svg xmlns="http://www.w3.org/2000/svg" width="%[1]s" height="%[1]s" viewBox="0 0 250 250">
			<rect width="100%%" height="100%%" fill="#11011D" />
			<g transform="translate(20, 10) scale(0.8)">
				%[3]s
			</g>	
		</svg>`, size, color, identiconRects)
	case "dither":
		svgBody := generateDitheredAvatar(name)
		avatarContent = strings.Replace(svgBody, `width="100%" height="100%"`, fmt.Sprintf(`width="%s" height="%s"`, size, size), 1)
	case "ascii":
		avatarContent = generateAsciiRobot(name, size)
	case "dotmatrix":
		avatarContent = generateDotMatrix(name, size)
	case "terminal":
		avatarContent = generateTerminalBlock(name, size)
	case "bauhaus":
		avatarContent = generateBauhaus(name, size)
	case "ring":
		avatarContent = generateRing(name, size)
	case "beam":
		avatarContent = generateBeam(name, size)
	case "marble":
		avatarContent = generateMarble(name, size)
	case "glitch":
		avatarContent = generateGlitch(name, size)
	case "sunset":
		avatarContent = generateSunset(name, size)
	case "smile":
		avatarContent = generateSmile(name, size)
	case "circuit":
		avatarContent = generateCircuit(name, size)
	case "pixel":
		avatarContent = generatePixel(name, size)
	case "constellation":
		loadGeoJSON()
		avatarContent = generateGeoJSONAvatar(name, size)
	default:
		http.Error(w, "Invalid avatar type. Use 'avatar' or 'gravatar'.", http.StatusBadRequest)
		return
	}

	// Store in cache
	if cache != nil {
		cache.Add(cacheKey, avatarContent)
	}

	w.Header().Set("X-Cache", "MISS")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(avatarContent))
}
