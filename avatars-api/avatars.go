package main

import (
	"crypto/md5"
	"fmt"
	"html"
	"math"
	"strings"
)

// getInitials extracts initials from a name string
func getInitials(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "?"
	}
	parts := strings.Fields(name)
	initials := string([]rune(parts[0])[0])
	if len(parts) > 1 {
		initials += string([]rune(parts[1])[0])
	}
	return strings.ToUpper(initials)
}

// pickPart selects an element from a list based on a hash byte
func pickPart(list []string, b byte) string {
	return list[int(b)%len(list)]
}

// generateIdenticon creates a GitHub-style identicon
func generateIdenticon(input string) string {
	color, hash := generateColor(input)
	var rects strings.Builder
	gridSize := 5
	cellSize := 50
	for col := range 3 {
		for row := range 5 {
			byteIndex := (col * 5) + row
			if hash[byteIndex%16]%2 == 0 {
				x := col * cellSize
				y := row * cellSize
				fmt.Fprintf(&rects, `<rect x="%d" y="%d" width="%d" height="%d" fill="%s" />`, x, y, cellSize, cellSize, color)
				if col < 2 {
					mirrorX := (gridSize - 1 - col) * cellSize
					fmt.Fprintf(&rects, `<rect x="%d" y="%d" width="%d" height="%d" fill="%s" />`, mirrorX, y, cellSize, cellSize, color)
				}
			}
		}
	}

	return rects.String()
}

// getPlasmaValue generates plasma effect value for dithering
func getPlasmaValue(x, y int, width, height float64, hash [16]byte) float64 {
	u := float64(x) / width
	v := float64(y) / height

	seedX := float64(hash[0]) / 10.0
	seedY := float64(hash[1]) / 10.0
	freq := 3.0 + (float64(hash[2] % 5))

	value := math.Sin(u*freq+seedX) + math.Cos(v*freq+seedY) + math.Sin((u+v)*freq)

	return (value + 3.0) / 6.0
}

// generateDitheredAvatar creates a dithered plasma avatar
func generateDitheredAvatar(name string) string {
	primaryColor, hash := generateColor(name)
	secondaryColor := "#11011D"
	cols, rows := 32, 32
	pixelSize := 10

	var rects strings.Builder

	for y := 0; y < rows; y++ {
		for x := range cols {
			luminance := getPlasmaValue(x, y, float64(cols), float64(rows), hash)

			threshold := bayerMatrix4x4[y%4][x%4] / 17.0

			fill := secondaryColor

			if (luminance + 0.1) < threshold {
				fill = primaryColor
			}

			if fill == primaryColor {
				rects.WriteString(fmt.Sprintf(
					`<rect x="%d" y="%d" width="%d" height="%d" fill="%s" />`,
					x*pixelSize, y*pixelSize, pixelSize, pixelSize, fill,
				))
			}
		}
	}

	return fmt.Sprintf(`
	<svg xmlns="http://www.w3.org/2000/svg" width="100%%" height="100%%" viewBox="0 0 320 320">
		<rect width="100%%" height="100%%" fill="%s" />
		<g transform="translate(16, 16) scale(0.9)">
			%s
		</g>
	</svg>`, secondaryColor, rects.String())
}

// generateAsciiRobot creates an ASCII art robot avatar
func generateAsciiRobot(name string, size string) string {
	bgColor, hash := generateColor(name)

	head := pickPart(robotHeads, hash[0])
	eyes := pickPart(robotEyes, hash[1])
	body := pickPart(robotBodies, hash[2])
	legs := pickPart(robotLegs, hash[3])

	asciiArt := []string{
		head,
		eyes,
		body,
		legs,
	}

	var textBlock strings.Builder
	yPos := 60
	lineHeight := 24

	for _, line := range asciiArt {
		safeLine := strings.ReplaceAll(line, " ", "\u00A0")
		// Escape any characters that would break XML (e.g., '<', '>' or '&')
		safeLine = html.EscapeString(safeLine)

		fmt.Fprintf(&textBlock, `<tspan x="50%%" dy="%d" text-anchor="middle">%s</tspan>`,
			lineHeight, safeLine)
	}

	return fmt.Sprintf(`
	<svg xmlns="http://www.w3.org/2000/svg" width="%[1]s" height="%[1]s" viewBox="0 0 200 200">
		<rect width="100%%" height="100%%" fill="%[2]s" />
		<g transform="translate(10, 10) scale(0.9)">
			<text x="50%%" y="%[3]d" font-family="monospace" font-weight="bold" font-size="28" fill="white" letter-spacing="2">
				%[4]s
			</text>
		</g>
	</svg>`, size, bgColor, yPos, textBlock.String())
}

// generateDotMatrix creates a dot matrix display avatar
func generateDotMatrix(name string, size string) string {
	mainColor, _ := generateColor(name)

	initials := getInitials(name)
	if len(initials) > 2 {
		initials = initials[:2]
	}

	var svgContent strings.Builder

	dotRadius := 4
	spacing := 12
	letterSpacing := 10
	startX := 20
	startY := 35

	currentX := startX

	for _, char := range initials {
		grid, ok := dotFont[char]
		if !ok {
			grid = dotFont['?']
		}

		for row := range 7 {
			for col := range 5 {
				cx := currentX + (col * spacing)
				cy := startY + (row * spacing)

				isLit := grid[row][col] == '1'

				if isLit {
					fmt.Fprintf(&svgContent, `<circle cx="%d" cy="%d" r="%d" fill="%s" opacity="0.4" />`,
						cx, cy, dotRadius+2, mainColor)
					fmt.Fprintf(&svgContent, "<circle cx=\"%d\" cy=\"%d\" r=\"%d\" fill=\"%s\" />\n", cx, cy, dotRadius, mainColor)
				} else {
					fmt.Fprintf(&svgContent, `<circle cx="%d" cy="%d" r="%d" fill="#333" opacity="0.3" />`,
						cx, cy, dotRadius)
				}
			}
		}
		currentX += (5 * spacing) + letterSpacing
	}

	return fmt.Sprintf(`
	<svg xmlns="http://www.w3.org/2000/svg" width="%[1]s" height="%[1]s" viewBox="0 0 170 170">
		<rect width="100%%" height="100%%" fill="#111111" />
		<g transform="translate(8.5, 8.5) scale(0.9)">
			%[2]s
		</g>
	</svg>`, size, svgContent.String())
}

// generateTerminalBlock creates a retro terminal block text avatar
func generateTerminalBlock(name string, size string) string {
	textColor, _ := generateColor(name)
	initials := getInitials(name)
	if len(initials) > 2 {
		initials = initials[:2]
	}

	var blocks strings.Builder

	blockSize := 20
	gap := 2

	startX := 40
	startY := 60
	letterSpacing := 20

	currentX := startX

	for _, char := range initials {
		grid, ok := blockFont[char]
		if !ok {
			grid = blockFont['?']
		}

		for row := range 5 {
			for col := range 5 {
				if grid[row][col] == '1' {
					x := currentX + (col * (blockSize + gap))
					y := startY + (row * (blockSize + gap))

					fmt.Fprintf(&blocks, `<rect x="%d" y="%d" width="%d" height="%d" fill="#000" opacity="0.5" />`,
						x+4, y+4, blockSize, blockSize)

					fmt.Fprintf(&blocks, `<rect x="%d" y="%d" width="%d" height="%d" fill="%s" />`,
						x, y, blockSize, blockSize, textColor)
				}
			}
		}
		currentX += (5 * (blockSize + gap)) + letterSpacing
	}
	cursorX := currentX
	fmt.Fprintf(&blocks, `<rect x="%d" y="%d" width="%d" height="%d" fill="%s" opacity="0.7" />`,
		cursorX, startY+(4*(blockSize+gap)), blockSize, blockSize, textColor)

	scanlines := `
	<defs>
		<pattern id="scanlines" patternUnits="userSpaceOnUse" width="10" height="4">
			<rect width="10" height="2" fill="#000" opacity="0.3" />
		</pattern>
	</defs>
	<rect width="100%" height="100%" fill="url(#scanlines)" />`

	return fmt.Sprintf(`
	<svg xmlns="http://www.w3.org/2000/svg" width="%[1]s" height="%[1]s" viewBox="0 0 350 350">
		<rect width="100%%" height="100%%" fill="#1a1b26" />
		%[3]s
		<g transform="translate(17.5, 17.5) scale(0.9)">
			%[2]s
		</g>
	</svg>`, size, blocks.String(), scanlines)
}

// generateBauhaus creates a Bauhaus-style geometric avatar
func generateBauhaus(name string, size string) string {
	hash := md5.Sum([]byte(name))
	bgIndex := int(hash[0]) % len(bauhausPalette)
	bgColor := bauhausPalette[bgIndex]

	var shapes strings.Builder

	numShapes := 3 + (int(hash[1]) % 3)

	for i := 0; i < numShapes; i++ {
		h1 := int(hash[i+2])
		h2 := int(hash[i+5])
		h3 := int(hash[i+8])

		color := bauhausPalette[(bgIndex+i+1)%len(bauhausPalette)]

		shapeType := h1 % 3

		x := h2 % 100
		y := h3 % 100
		w := 20 + (h1 % 60)

		opacity := 0.5 + (float64(h2%5) / 10.0)

		switch shapeType {
		case 0:
			fmt.Fprintf(&shapes, `<circle cx="%d" cy="%d" r="%d" fill="%s" opacity="%.2f" />`,
				x, y, w/2, color, opacity)
		case 1:
			rotation := 0
			if h1%2 == 0 {
				rotation = 45
			}
			fmt.Fprintf(&shapes, `<rect x="%d" y="%d" width="%d" height="%d" fill="%s" opacity="%.2f" transform="rotate(%d %d %d)" />`,
				x-w/2, y-w/2, w, w, color, opacity, rotation, x, y)
		case 2:
			p1 := fmt.Sprintf("%d,%d", x, y-(w/2))
			p2 := fmt.Sprintf("%d,%d", x-(w/2), y+(w/2))
			p3 := fmt.Sprintf("%d,%d", x+(w/2), y+(w/2))
			fmt.Fprintf(&shapes, `<polygon points="%s %s %s" fill="%s" opacity="%.2f" />`,
				p1, p2, p3, color, opacity)
		}
	}

	return fmt.Sprintf(`
	<svg xmlns="http://www.w3.org/2000/svg" width="%[1]s" height="%[1]s" viewBox="0 0 100 100">
		<rect width="100%%" height="100%%" fill="%[2]s" />
		<g transform="translate(5, 5) scale(0.9)">
			%[3]s
		</g>
	</svg>`, size, bgColor, shapes.String())
}

// generateRing creates a gradient ring avatar
func generateRing(name string, size string) string {
	hash := md5.Sum([]byte(name))
	c1, _ := generateColor(name)
	c2, _ := generateColor(name + "2")
	c3, _ := generateColor(name + "3")

	angle := int(hash[0]) % 360

	gradID := fmt.Sprintf("grad-%x", hash[0:3])

	return fmt.Sprintf(`
	<svg xmlns="http://www.w3.org/2000/svg" width="%[1]s" height="%[1]s" viewBox="0 0 100 100">
		<defs>
			<linearGradient id="%[5]s" x1="0%%" y1="0%%" x2="100%%" y2="100%%" gradientTransform="rotate(%[6]d .5 .5)">
				<stop offset="0%%" stop-color="%[2]s" />
				<stop offset="50%%" stop-color="%[3]s" />
				<stop offset="100%%" stop-color="%[4]s" />
			</linearGradient>
		</defs>
		<g transform="translate(5, 5) scale(0.9)">
			<circle cx="50" cy="50" r="50" fill="url(#%[5]s)" />
		</g>
	</svg>`, size, c1, c2, c3, gradID, angle)
}

// generateBeam creates a connected network nodes avatar
func generateBeam(name string, size string) string {
	hash := md5.Sum([]byte(name))
	bgColor := "#0a0a0a"
	accentColor, _ := generateColor(name)

	var svgContent strings.Builder

	type Point struct{ X, Y int }
	points := make([]Point, 6)

	for i := range 6 {
		points[i] = Point{
			X: 10 + int(hash[i])%80,
			Y: 10 + int(hash[i+6])%80,
		}
		fmt.Fprintf(&svgContent, `<circle cx="%d" cy="%d" r="3" fill="%s" />`, points[i].X, points[i].Y, accentColor)
	}
	for i := range 6 {
		for j := i + 1; j < 6; j++ {
			p1 := points[i]
			p2 := points[j]
			distSq := (p1.X-p2.X)*(p1.X-p2.X) + (p1.Y-p2.Y)*(p1.Y-p2.Y)

			if distSq < 3600 {
				opacity := 1.0 - (float64(distSq) / 3600.0)
				fmt.Fprintf(&svgContent, `<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="%s" stroke-width="1" opacity="%.2f" />`, p1.X, p1.Y, p2.X, p2.Y, accentColor, opacity)
			}
		}
	}

	return fmt.Sprintf(`
	<svg xmlns="http://www.w3.org/2000/svg" width="%[1]s" height="%[1]s" viewBox="0 0 100 100">
		<rect width="100%%" height="100%%" fill="%[2]s" />
		<g transform="translate(5, 5) scale(0.9)">
			%[3]s
		</g>
	</svg>`, size, bgColor, svgContent.String())
}

// generateMarble creates a marble texture avatar
func generateMarble(name string, size string) string {
	hash := md5.Sum([]byte(name))
	c1, _ := generateColor(name)
	c2, _ := generateColor(name + "x")

	freq := 0.005 + (float64(hash[0])/255.0)*0.02

	octaves := 1 + (int(hash[1]) % 4)

	return fmt.Sprintf(`
	<svg xmlns="http://www.w3.org/2000/svg" width="%[1]s" height="%[1]s" viewBox="0 0 100 100">
		<defs>
			<filter id="liquid">
				<feTurbulence type="fractalNoise" baseFrequency="%[4].4f" numOctaves="%[5]d" result="noise" />
				<feDiffuseLighting in="noise" lighting-color="white" surfaceScale="2">
					<feDistantLight azimuth="45" elevation="60" />
				</feDiffuseLighting>
			</filter>
			<linearGradient id="grad" x1="0%%" y1="0%%" x2="100%%" y2="100%%">
				<stop offset="0%%" stop-color="%[2]s" />
				<stop offset="100%%" stop-color="%[3]s" />
			</linearGradient>
		</defs>
		<g transform="translate(5, 5) scale(0.9)">
			<rect width="100%%" height="100%%" fill="url(#grad)" />
			<rect width="100%%" height="100%%" fill="transparent" filter="url(#liquid)" opacity="0.5" style="mix-blend-mode: overlay;" />
		</g>
	</svg>`, size, c1, c2, freq, octaves)
}

// generateGlitch creates a cyberpunk glitch effect avatar
func generateGlitch(name string, size string) string {
	hash := md5.Sum([]byte(name))
	initials := getInitials(name)

	bgColor := "#0f0f0f"

	var glitchLines strings.Builder
	for i := 0; i < 5; i++ {
		y := int(hash[i]) % 100
		h := int(hash[i+5])%5 + 1
		w := int(hash[i+10])%50 + 20
		x := int(hash[i+2]) % 80

		glitchLines.WriteString(fmt.Sprintf(
			`<rect x="%d" y="%d" width="%d" height="%d" fill="white" opacity="0.1" />`,
			x, y, w, h,
		))
	}

	return fmt.Sprintf(`
	<svg xmlns="http://www.w3.org/2000/svg" width="%[1]s" height="%[1]s" viewBox="0 0 100 100">
		<rect width="100%%" height="100%%" fill="%[2]s" />
		<g transform="translate(5, 5) scale(0.9)">
			<text x="48" y="55" dominant-baseline="middle" text-anchor="middle" font-family="Arial Black, sans-serif" font-weight="900" font-size="50" fill="#00ffff" opacity="0.8" style="mix-blend-mode: screen;">
				%[3]s
			</text>
			
			<text x="52" y="55" dominant-baseline="middle" text-anchor="middle" font-family="Arial Black, sans-serif" font-weight="900" font-size="50" fill="#ff0000" opacity="0.8" style="mix-blend-mode: screen;">
				%[3]s
			</text>
			
			<text x="50" y="55" dominant-baseline="middle" text-anchor="middle" font-family="Arial Black, sans-serif" font-weight="900" font-size="50" fill="#ffffff">
				%[3]s
			</text>
			
			%[4]s
		</g>
	</svg>`, size, bgColor, initials, glitchLines.String())
}

// generateSunset creates a procedural sunset scene avatar
func generateSunset(name string, size string) string {
	hash := md5.Sum([]byte(name))

	var skyTop, skyBot string
	mood := hash[0] % 3
	switch mood {
	case 0:
		skyTop, skyBot = "#3e1c6b", "#ff8a5c"
	case 1:
		skyTop, skyBot = "#29b6f6", "#fff9c4"
	default:
		skyTop, skyBot = "#0d1b2a", "#415a77"
	}
	sunX := 20 + (int(hash[1]) % 60)
	sunY := 20 + (int(hash[2]) % 30)
	sunColor := "#ffffff"
	if mood == 0 {
		sunColor = "#ffeb3b"
	}
	var mountains strings.Builder
	mountains.WriteString("M 0 100 L 0 60 ")

	for x := 0; x <= 100; x += 5 {
		y := 60.0 + math.Sin(float64(x)*0.1+float64(hash[3]))*15.0
		if x%10 == 0 {
			y -= 5
		}

		fmt.Fprintf(&mountains, "L %d %.2f ", x, y)
	}
	mountains.WriteString("L 100 100 Z")

	mountColor := "#1a1a1a"
	if mood == 1 {
		mountColor = "#4caf50"
	}

	return fmt.Sprintf(`
	<svg xmlns="http://www.w3.org/2000/svg" width="%[1]s" height="%[1]s" viewBox="0 0 100 100">
		<defs>
			<linearGradient id="sky" x1="0%%" y1="0%%" x2="0%%" y2="100%%">
				<stop offset="0%%" stop-color="%[2]s" />
				<stop offset="100%%" stop-color="%[3]s" />
			</linearGradient>
		</defs>
		<rect width="100%%" height="100%%" fill="url(#sky)" />
		<g transform="translate(5, 5) scale(0.9)">
			<circle cx="%d" cy="%d" r="8" fill="%s" opacity="0.9" />
			
			<path d="%s" fill="%s" opacity="0.9" />
		</g>
	</svg>`, size, skyTop, skyBot, sunX, sunY, sunColor, mountains.String(), mountColor)
}

// generateSmile creates a minimalist face avatar
func generateSmile(name string, size string) string {
	hash := md5.Sum([]byte(name))
	skinTones := []string{"#FFDFC4", "#F0C8C9", "#E5B99F", "#8D5524", "#C68642", "#FFDCB1", "#E0AC69", "#B9D2B1", "#A8C8E8"}
	skinColor := skinTones[int(hash[0])%len(skinTones)]
	eyeType := int(hash[1]) % 3
	mouthType := int(hash[2]) % 4
	hasBlush := int(hash[3])%2 == 0

	var features strings.Builder

	switch eyeType {
	case 0:
		features.WriteString(`<circle cx="35" cy="45" r="5" fill="#333" /><circle cx="65" cy="45" r="5" fill="#333" />`)
	case 1:
		features.WriteString(`<path d="M 30 45 Q 35 40 40 45" stroke="#333" stroke-width="3" fill="none" />`)
		features.WriteString(`<path d="M 60 45 Q 65 40 70 45" stroke="#333" stroke-width="3" fill="none" />`)
	case 2:
		features.WriteString(`<circle cx="35" cy="45" r="5" fill="#333" />`)
		features.WriteString(`<rect x="60" y="44" width="10" height="2" fill="#333" />`)
	}

	switch mouthType {
	case 0:
		features.WriteString(`<path d="M 35 65 Q 50 75 65 65" stroke="#333" stroke-width="3" fill="none" stroke-linecap="round" />`)
	case 1:
		features.WriteString(`<path d="M 35 65 Q 50 80 65 65 Z" fill="#fff" stroke="#333" stroke-width="2" />`)
	case 2:
		features.WriteString(`<line x1="40" y1="70" x2="60" y2="70" stroke="#333" stroke-width="3" stroke-linecap="round" />`)
	case 3:
		features.WriteString(`<circle cx="50" cy="70" r="6" stroke="#333" stroke-width="3" fill="none" />`)
	}

	if hasBlush {
		features.WriteString(`<circle cx="30" cy="55" r="5" fill="#ff0000" opacity="0.2" />`)
		features.WriteString(`<circle cx="70" cy="55" r="5" fill="#ff0000" opacity="0.2" />`)
	}

	return fmt.Sprintf(`
	<svg xmlns="http://www.w3.org/2000/svg" width="%[1]s" height="%[1]s" viewBox="0 0 100 100">
		<g transform="translate(5, 5) scale(0.9)">
			<circle cx="50" cy="50" r="45" fill="%[2]s" />
			%[3]s
		</g>
	</svg>`, size, skinColor, features.String())
}

// generateCircuit creates a circuit board pattern avatar
func generateCircuit(name string, size string) string {
	hash := md5.Sum([]byte(name))

	boardColors := []string{"#004d40", "#1a237e", "#212121", "#1b5e20"}
	bgColor := boardColors[int(hash[0])%len(boardColors)]
	traceColor := "#ffd700"

	var traces strings.Builder

	for i := range 5 {
		x1 := 10 + (int(hash[i]) % 80)
		y1 := 10 + (int(hash[i+5]) % 80)
		x2 := 10 + (int(hash[i+2]) % 80)
		y2 := 10 + (int(hash[i+7]) % 80)

		fmt.Fprintf(&traces, `<path d="M %d %d L %d %d L %d %d" stroke="%s" stroke-width="2" fill="none" opacity="0.8" />`,
			x1, y1, x2, y1, x2, y2, traceColor)
		traces.WriteString(fmt.Sprintf(`<circle cx="%d" cy="%d" r="3" fill="%s" />`, x1, y1, traceColor))
		traces.WriteString(fmt.Sprintf(`<circle cx="%d" cy="%d" r="3" fill="%s" />`, x2, y2, traceColor))
	}

	return fmt.Sprintf(`
	<svg xmlns="http://www.w3.org/2000/svg" width="%[1]s" height="%[1]s" viewBox="0 0 100 100">
		<rect width="100%%" height="100%%" fill="%[2]s" />
		<g transform="translate(5, 5) scale(0.9)">
			%[3]s
		</g>
	</svg>`, size, bgColor, traces.String())
}

// generatePixel creates an isometric pixel art cube avatar
func generatePixel(name string, size string) string {
	baseHex, hash := generateColor(name)
	var topPattern strings.Builder
	if hash[0]%2 == 0 {
		topPattern.WriteString(`<path d="M 50 30 L 70 40 L 50 50 L 30 40 Z" fill="rgba(255,255,255,0.3)" />`)
	}

	return fmt.Sprintf(`
	<svg xmlns="http://www.w3.org/2000/svg" width="%[1]s" height="%[1]s" viewBox="0 0 100 100">
		<rect width="100%%" height="100%%" fill="#f0f0f0" />
		<g transform="translate(5, 5) scale(0.9)">
			<path d="M 20 35 L 50 50 L 50 80 L 20 65 Z" fill="%[2]s" />
			
			<path d="M 50 50 L 80 35 L 80 65 L 50 80 Z" fill="%[2]s" />
			<path d="M 50 50 L 80 35 L 80 65 L 50 80 Z" fill="black" opacity="0.2" />
			
			<path d="M 50 20 L 80 35 L 50 50 L 20 35 Z" fill="%[2]s" />
			<path d="M 50 20 L 80 35 L 50 50 L 20 35 Z" fill="white" opacity="0.3" />
			
			%[3]s
		</g>
	</svg>`, size, baseHex, topPattern.String())
}
