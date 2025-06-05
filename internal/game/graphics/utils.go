package graphics

import (
	"embed"
	"fmt"
)

//go:embed shaders/*.glsl
var shaderFS embed.FS

var (
	VertexShaderSrc   string
	FragmentShaderSrc string
)

func LoadShaders() error {
	vSrc, err := shaderFS.ReadFile("shaders/vertex.glsl")
	if err != nil {
		return fmt.Errorf("failed to read vertex shader: %w", err)
	}

	fSrc, err := shaderFS.ReadFile("shaders/fragment.glsl")
	if err != nil {
		return fmt.Errorf("failed to read fragment shader: %w", err)
	}

	VertexShaderSrc = string(vSrc)
	FragmentShaderSrc = string(fSrc)
	return nil
}
