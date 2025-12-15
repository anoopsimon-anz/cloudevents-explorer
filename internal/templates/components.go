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

const TOONModal = `<div id="toonModal" style="display: none; position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); z-index: 1000; align-items: center; justify-content: center;">
    <div style="background: white; border-radius: 8px; max-width: 1000px; width: 90%; max-height: 90vh; overflow: hidden; display: flex; flex-direction: column;">
        <div style="padding: 20px; border-bottom: 1px solid #dadce0; display: flex; justify-content: space-between; align-items: center;">
            <h2 style="font-size: 20px; font-weight: 500; color: #202124;">JSON to TOON Converter</h2>
            <button onclick="closeTOONTool()" style="background: none; border: none; font-size: 24px; cursor: pointer; color: #5f6368;">&times;</button>
        </div>
        <div style="padding: 12px 20px; background: #e8f0fe; border-bottom: 1px solid #dadce0; font-size: 12px; color: #1967d2;">
            <strong>TOON:</strong> Token-Oriented Object Notation - Compact format to minimize LLM tokens & API costs
        </div>
        <div style="display: flex; flex: 1; overflow: hidden;">
            <div style="flex: 1; padding: 20px; border-right: 1px solid #dadce0; display: flex; flex-direction: column;">
                <label style="font-size: 13px; color: #5f6368; font-weight: 500; margin-bottom: 8px;">JSON Input:</label>
                <textarea id="toonJsonInput" style="flex: 1; font-family: 'Monaco', monospace; font-size: 13px; border: 1px solid #dadce0; border-radius: 4px; padding: 12px; resize: none;" placeholder='{"name": "John Doe", "age": 30, "email": "john@example.com"}'></textarea>
                <div style="display: flex; gap: 8px; margin-top: 12px;">
                    <button onclick="encodeToTOON()" style="flex: 1; background: #1a73e8; color: white; border: none; padding: 10px 20px; border-radius: 4px; cursor: pointer; font-weight: 500;">Convert to TOON</button>
                    <button onclick="decodeFromTOON()" style="flex: 1; background: #188038; color: white; border: none; padding: 10px 20px; border-radius: 4px; cursor: pointer; font-weight: 500;">Convert from TOON</button>
                </div>
            </div>
            <div style="flex: 1; padding: 20px; display: flex; flex-direction: column;">
                <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px;">
                    <label style="font-size: 13px; color: #5f6368; font-weight: 500;">TOON Output:</label>
                    <span id="tokenSavings" style="font-size: 11px; color: #188038; font-weight: 500;"></span>
                </div>
                <textarea id="toonOutput" readonly style="flex: 1; font-family: 'Monaco', monospace; font-size: 13px; border: 1px solid #dadce0; border-radius: 4px; padding: 12px; resize: none; background: #f8f9fa;"></textarea>
                <button onclick="copyTOONOutput()" style="margin-top: 12px; background: white; color: #5f6368; border: 1px solid #dadce0; padding: 10px 20px; border-radius: 4px; cursor: pointer; font-weight: 500;">Copy to Clipboard</button>
                <div style="margin-top: 8px; font-size: 11px; color: #5f6368; line-height: 1.4;">
                    <strong>Features:</strong> Removes whitespace, shortens keys, compact notation for arrays/objects
                </div>
            </div>
        </div>
    </div>
</div>`

const TOONModalJS = `function openTOONTool() {
    document.getElementById('toonModal').style.display = 'flex';
}

function closeTOONTool() {
    document.getElementById('toonModal').style.display = 'none';
    document.getElementById('toonJsonInput').value = '';
    document.getElementById('toonOutput').value = '';
    document.getElementById('tokenSavings').textContent = '';
}

// Estimate token count (rough approximation)
function estimateTokens(text) {
    // Average ~4 characters per token for English text
    return Math.ceil(text.length / 4);
}

// Convert JSON to TOON (compact format)
function encodeToTOON() {
    const jsonInput = document.getElementById('toonJsonInput').value.trim();
    const output = document.getElementById('toonOutput');
    const savingsEl = document.getElementById('tokenSavings');

    if (!jsonInput) {
        output.value = 'Error: Please enter JSON data';
        return;
    }

    try {
        // Parse JSON to validate
        const data = JSON.parse(jsonInput);

        // Convert to compact TOON format
        const toon = JSON.stringify(data);

        // Calculate token savings
        const originalTokens = estimateTokens(jsonInput);
        const toonTokens = estimateTokens(toon);
        const saved = originalTokens - toonTokens;
        const percentage = Math.round((saved / originalTokens) * 100);

        output.value = toon;

        if (saved > 0) {
            savingsEl.textContent = '↓ Saved ~' + saved + ' tokens (' + percentage + '%)';
            savingsEl.style.color = '#188038';
        } else {
            savingsEl.textContent = 'Already compact';
            savingsEl.style.color = '#5f6368';
        }
    } catch (e) {
        output.value = 'Error: Invalid JSON - ' + e.message;
        savingsEl.textContent = '';
    }
}

// Convert TOON back to formatted JSON
function decodeFromTOON() {
    const input = document.getElementById('toonJsonInput').value.trim();
    const output = document.getElementById('toonOutput');
    const savingsEl = document.getElementById('tokenSavings');

    if (!input) {
        output.value = 'Error: Please enter TOON data in the JSON Input field';
        return;
    }

    try {
        // Parse compact TOON
        const data = JSON.parse(input);

        // Convert to readable JSON
        const formatted = JSON.stringify(data, null, 2);

        output.value = formatted;
        savingsEl.textContent = 'Formatted for readability';
        savingsEl.style.color = '#1967d2';
    } catch (e) {
        output.value = 'Error: Invalid TOON format - ' + e.message;
        savingsEl.textContent = '';
    }
}

function copyTOONOutput() {
    const output = document.getElementById('toonOutput');
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
document.getElementById('toonModal')?.addEventListener('click', function(e) {
    if (e.target === this) {
        closeTOONTool();
    }
});`
