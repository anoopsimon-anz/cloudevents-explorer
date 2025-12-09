package templates

const Base64Modal = `<div id="base64Modal" style="display: none; position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); z-index: 1000; align-items: center; justify-content: center;">
    <div style="background: white; border-radius: 8px; max-width: 900px; width: 90%; max-height: 90vh; overflow: hidden; display: flex; flex-direction: column;">
        <div style="padding: 20px; border-bottom: 1px solid #dadce0; display: flex; justify-content: space-between; align-items: center;">
            <h2 style="font-size: 20px; font-weight: 500; color: #202124;">Base64 Encoder/Decoder</h2>
            <button onclick="closeBase64Tool()" style="background: none; border: none; font-size: 24px; cursor: pointer; color: #5f6368;">&times;</button>
        </div>
        <div style="display: flex; flex: 1; overflow: hidden;">
            <div style="flex: 1; padding: 20px; border-right: 1px solid #dadce0; display: flex; flex-direction: column;">
                <label style="font-size: 13px; color: #5f6368; font-weight: 500; margin-bottom: 8px;">Input Text:</label>
                <textarea id="base64Input" style="flex: 1; font-family: 'Monaco', monospace; font-size: 13px; border: 1px solid #dadce0; border-radius: 4px; padding: 12px; resize: none;" placeholder="Enter text or Base64 string here"></textarea>
                <div style="display: flex; gap: 8px; margin-top: 12px;">
                    <button onclick="encodeBase64()" style="flex: 1; background: #1a73e8; color: white; border: none; padding: 10px 20px; border-radius: 4px; cursor: pointer; font-weight: 500;">Encode to Base64</button>
                    <button onclick="decodeBase64()" style="flex: 1; background: #188038; color: white; border: none; padding: 10px 20px; border-radius: 4px; cursor: pointer; font-weight: 500;">Decode from Base64</button>
                </div>
            </div>
            <div style="flex: 1; padding: 20px; display: flex; flex-direction: column;">
                <label style="font-size: 13px; color: #5f6368; font-weight: 500; margin-bottom: 8px;">Output:</label>
                <textarea id="base64Output" readonly style="flex: 1; font-family: 'Monaco', monospace; font-size: 13px; border: 1px solid #dadce0; border-radius: 4px; padding: 12px; resize: none; background: #f8f9fa;"></textarea>
                <button onclick="copyOutput()" style="margin-top: 12px; background: white; color: #5f6368; border: 1px solid #dadce0; padding: 10px 20px; border-radius: 4px; cursor: pointer; font-weight: 500;">Copy to Clipboard</button>
            </div>
        </div>
    </div>
</div>`

const Base64ModalJS = `function openBase64Tool() {
    document.getElementById('base64Modal').style.display = 'flex';
}

function closeBase64Tool() {
    document.getElementById('base64Modal').style.display = 'none';
    document.getElementById('base64Input').value = '';
    document.getElementById('base64Output').value = '';
}

function encodeBase64() {
    const input = document.getElementById('base64Input').value;
    const output = document.getElementById('base64Output');

    if (!input) {
        output.value = 'Error: Please enter some text to encode';
        return;
    }

    try {
        const encoded = btoa(unescape(encodeURIComponent(input)));
        output.value = encoded;
    } catch (e) {
        output.value = 'Error: Failed to encode - ' + e.message;
    }
}

function decodeBase64() {
    const input = document.getElementById('base64Input').value;
    const output = document.getElementById('base64Output');

    if (!input) {
        output.value = 'Error: Please enter a Base64 string to decode';
        return;
    }

    try {
        const decoded = decodeURIComponent(escape(atob(input)));
        output.value = decoded;
    } catch (e) {
        output.value = 'Error: Invalid Base64 string - ' + e.message;
    }
}

function copyOutput() {
    const output = document.getElementById('base64Output');
    if (!output.value || output.value.startsWith('Error:')) {
        return;
    }
    output.select();
    document.execCommand('copy');

    const btn = event.target;
    const originalText = btn.textContent;
    btn.textContent = 'Copied!';
    btn.style.background = '#188038';
    btn.style.color = 'white';
    setTimeout(function() {
        btn.textContent = originalText;
        btn.style.background = 'white';
        btn.style.color = '#5f6368';
    }, 2000);
}

// Close modal on outside click
document.getElementById('base64Modal')?.addEventListener('click', function(e) {
    if (e.target === this) {
        closeBase64Tool();
    }
});`
