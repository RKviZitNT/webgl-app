#version 100
precision highp float;

varying vec2 vTexCoords;

uniform sampler2D uTexture;

void main() {
    gl_FragColor = texture2D(uTexture, vTexCoords);
}