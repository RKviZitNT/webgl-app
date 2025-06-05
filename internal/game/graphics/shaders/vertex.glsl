attribute vec2 a_position;
uniform vec2 u_offset;

void main() {
    vec2 pos = a_position + u_offset;
    gl_Position = vec4(pos, 0.0, 1.0);
}
