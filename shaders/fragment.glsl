precision mediump float;

uniform sampler2D u_texture;
varying vec2 v_texCoord;

void main() {
    gl_FragColor = texture2D(u_texture, v_texCoord);
}

// precision mediump float;

// void main() {
//     gl_FragColor = vec4(0.2, 0.6, 1.0, 1.0);
// }
