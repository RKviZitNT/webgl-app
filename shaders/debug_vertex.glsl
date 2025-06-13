#version 100
precision highp float;

attribute vec2 aPosition;
uniform vec2 uResolution;

void main() {
    vec2 normalizedPos = vec2(
        (aPosition.x / uResolution.x) * 2.0 - 1.0,
        1.0 - (aPosition.y / uResolution.y) * 2.0
    );

    gl_Position = vec4(normalizedPos, 0.0, 1.0);
}