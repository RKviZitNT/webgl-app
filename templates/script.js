function showScreen(screenId) {
    document.querySelectorAll('.screen').forEach(el => {
        el.classList.remove('visible');
    });
    document.getElementById(screenId).classList.add('visible');
}

function copyLobbyCode() {
    const cpsc = 'Copied!'
    const codeElement = document.getElementById('lobby_code');
    const code = codeElement.textContent;
    
    if (code != cpsc) {
        navigator.clipboard.writeText(code).then(() => {
            const originalText = codeElement.textContent;
            codeElement.textContent = cpsc;
            setTimeout(() => {
                codeElement.textContent = originalText;
            }, 1000);
        }).catch(err => {
            console.error('Failed to copy code:', err);
        });
    }
}

function updateOwnerControls(isOwner) {
    const startBtn = document.getElementById('start_button');
    startBtn.style.display = isOwner ? 'block' : 'none';
}

function switchStartButtonState(isEnabled) {
    const startBtn = document.getElementById('start_button');
    startBtn.disabled = !isEnabled;
    
    if (!isEnabled) {
        startBtn.classList.add('disabled');
    } else {
        startBtn.classList.remove('disabled');
    }
}

function resizeCanvas() {
    const canvas = document.getElementById('game_canvas');
    if (!canvas) return;
    canvas.width = window.innerWidth;
    canvas.height = window.innerHeight;
}

async function loadShader(path) {
    try {
        const shader = await fetch(path);
        if (!shader.ok) throw new Error(`HTTP error! status: ${shader.status}`);
        return await shader.text();
    } catch (error) {
        console.error("Error loading shaders:", error);
        throw error;
    }
}

async function loadImage(path) {
    try {
        const img = new Image();
        img.src = path;
        await new Promise((resolve, reject) => {
            img.onload = resolve
            img.onerror = () => reject(new Error(`Failed to load image: ${path}`))
        });
        return img
    } catch(error) {
        console.error("Image load error:", error);
        throw error;
    }
}

async function init() {
    const go = new Go();
    const result = await WebAssembly.instantiateStreaming(
        fetch('main.wasm'), 
        go.importObject
    );
    
    await go.run(result.instance);
}

init().catch(err => console.error("Initialization failed:", err));