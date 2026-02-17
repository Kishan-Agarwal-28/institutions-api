package main

import (
	"crypto/md5"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strings"
)

// GeoJSON data structures
type GeoJSON struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

type Feature struct {
	ID         string     `json:"id"`
	Properties Properties `json:"properties"`
	Geometry   Geometry   `json:"geometry"`
}

type Properties struct {
	Name string `json:"n"`
}

type Geometry struct {
	Type        string        `json:"type"`
	Coordinates [][][]float64 `json:"coordinates"`
}

var constellationData GeoJSON

//go:embed constellations.json
var constellationFile embed.FS

// loadGeoJSON loads constellation data from embedded JSON file
func loadGeoJSON() {
	data, err := constellationFile.ReadFile("constellations.json")
	if err != nil {
		log.Fatalf("Failed to read embedded json: %v", err)
	}

	err = json.Unmarshal(data, &constellationData)
	if err != nil {
		log.Fatalf("Failed to parse json: %v", err)
	}

	fmt.Printf("✅ Loaded %d constellations\n", len(constellationData.Features))
}

// normalizeGeoJSON normalizes constellation coordinates to fit in a 100x100 viewBox
func normalizeGeoJSON(coords [][][]float64) [][][]float64 {
	minX, maxX := math.MaxFloat64, -math.MaxFloat64
	minY, maxY := math.MaxFloat64, -math.MaxFloat64

	for _, line := range coords {
		for _, point := range line {
			x, y := point[0], point[1]
			if x < minX {
				minX = x
			}
			if x > maxX {
				maxX = x
			}
			if y < minY {
				minY = y
			}
			if y > maxY {
				maxY = y
			}
		}
	}

	width := maxX - minX
	height := maxY - minY
	if width == 0 {
		width = 1
	}
	if height == 0 {
		height = 1
	}
	scaleX := 60.0 / width
	scaleY := 60.0 / height
	scale := math.Min(scaleX, scaleY)

	offsetX := (100.0 - (width * scale)) / 2.0
	offsetY := (100.0 - (height * scale)) / 2.0

	normalized := make([][][]float64, len(coords))
	for i, line := range coords {
		normalized[i] = make([][]float64, len(line))
		for j, point := range line {
			newX := ((point[0] - minX) * scale) + offsetX
			newY := 100 - (((point[1] - minY) * scale) + offsetY)
			normalized[i][j] = []float64{newX, newY}
		}
	}
	return normalized
}

// generateGeoJSONAvatar creates a constellation-based avatar
func generateGeoJSONAvatar(name string, size string) string {
	hash := md5.Sum([]byte(name))
	idx := int(hash[0]) % len(constellationData.Features)
	feature := constellationData.Features[idx]
	lines := normalizeGeoJSON(feature.Geometry.Coordinates)

	bgStart := "#1e1b4b"
	bgEnd := "#020617"
	var svgContent strings.Builder

	// Add background stars
	for i := range 40 {
		x := (int(hash[i%16]) * (i + 3)) % 100
		y := (int(hash[(i+2)%16]) * (i + 5)) % 100
		op := float64((int(hash[i%16])%5)+1) / 10.0
		fmt.Fprintf(&svgContent, `<circle cx="%d" cy="%d" r="0.4" fill="white" opacity="%.1f" />`, x, y, op)
	}

	// Draw constellation lines and nodes
	for _, line := range lines {
		for i := 0; i < len(line)-1; i++ {
			p1 := line[i]
			p2 := line[i+1]

			fmt.Fprintf(&svgContent, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#93c5fd" stroke-width="0.5" opacity="0.8" />`,
				p1[0], p1[1], p2[0], p2[1])

			fmt.Fprintf(&svgContent, `<circle cx="%.1f" cy="%.1f" r="1.5" fill="white" />`, p1[0], p1[1])
			fmt.Fprintf(&svgContent, `<circle cx="%.1f" cy="%.1f" r="3" fill="#38bdf8" opacity="0.2" />`, p1[0], p1[1])
		}

		lastP := line[len(line)-1]
		fmt.Fprintf(&svgContent, `<circle cx="%.1f" cy="%.1f" r="1.5" fill="white" />`, lastP[0], lastP[1])
	}

	return fmt.Sprintf(`
	<svg xmlns="http://www.w3.org/2000/svg" width="%[1]s" height="%[1]s" viewBox="0 0 100 100">
		<defs>
			<radialGradient id="grad" cx="50%%" cy="50%%" r="80%%">
				<stop offset="0%%" stop-color="%[2]s" />
				<stop offset="100%%" stop-color="%[3]s" />
			</radialGradient>
		</defs>
		<rect width="100%%" height="100%%" fill="url(#grad)" />
		%[4]s
		<text x="50" y="90" text-anchor="middle" font-family="Times New Roman" font-weight="bold" font-size="6" fill="#7dd3fc" letter-spacing="0.5">
			%[5]s
		</text>
	</svg>`, size, bgStart, bgEnd, svgContent.String(), strings.ToUpper(feature.Properties.Name))
}
