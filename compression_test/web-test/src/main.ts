import "./style.css";

interface TestResult {
  filename: string;
  size: number;
  decompressTime: number;
  parseTime: number;
  compressionType: string;
}

// Import test files
const TEST_FILES = {
  jmdict: {
    gzip: new URL("./test_data/jmdict.json.gz", import.meta.url).href,
    brotli: new URL("./test_data/jmdict.json.br", import.meta.url).href,
    raw: new URL("./test_data/jmdict.json", import.meta.url).href,
  },
  jmnedict: {
    gzip: new URL("./test_data/jmnedict.json.gz", import.meta.url).href,
    brotli: new URL("./test_data/jmnedict.json.br", import.meta.url).href,
    raw: new URL("./test_data/jmnedict.json", import.meta.url).href,
  },
  kanjidic: {
    gzip: new URL("./test_data/kanjidic.json.gz", import.meta.url).href,
    brotli: new URL("./test_data/kanjidic.json.br", import.meta.url).href,
    raw: new URL("./test_data/kanjidic.json", import.meta.url).href,
  },
};

document.querySelector<HTMLDivElement>("#app")!.innerHTML = `
  <div>
    <h1>Dictionary Decompression Test</h1>
    <div class="card">
      <div class="test-controls">
        <button id="runAllTests">Run All Tests</button>
        <button id="clearResults">Clear Results</button>
      </div>
      <div class="results" id="results">
        <h2>Results</h2>
        <div id="resultsList"></div>
      </div>
    </div>
  </div>
`;

const runAllButton = document.querySelector<HTMLButtonElement>("#runAllTests")!;
const clearButton = document.querySelector<HTMLButtonElement>("#clearResults")!;
const resultsList = document.querySelector<HTMLDivElement>("#resultsList")!;

async function fetchAndDecompress(
  url: string,
  compressionType: string
): Promise<TestResult> {
  const startFetch = performance.now();
  const response = await fetch(url);
  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }

  const startDecompress = performance.now();
  const arrayBuffer = await response.arrayBuffer();
  let text: string;

  if (compressionType === "gzip" || compressionType === "brotli") {
    // Browser will automatically decompress gzip and brotli
    text = new TextDecoder().decode(arrayBuffer);
  } else {
    // Raw JSON
    text = new TextDecoder().decode(arrayBuffer);
  }
  const decompressTime = performance.now() - startDecompress;

  const startParse = performance.now();
  JSON.parse(text);
  const parseTime = performance.now() - startParse;

  return {
    filename: url.split("/").pop() || url,
    size: arrayBuffer.byteLength,
    decompressTime,
    parseTime,
    compressionType,
  };
}

function formatBytes(bytes: number): string {
  if (bytes === 0) return "0 Bytes";
  const k = 1024;
  const sizes = ["Bytes", "KB", "MB", "GB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
}

function displayResult(result: TestResult) {
  const resultElement = document.createElement("div");
  resultElement.className = "result-item";
  resultElement.innerHTML = `
    <h3>${result.filename}</h3>
    <p>Compression: ${result.compressionType}</p>
    <p>File Size: ${formatBytes(result.size)}</p>
    <p>Decompression Time: ${result.decompressTime.toFixed(2)}ms</p>
    <p>Parse Time: ${result.parseTime.toFixed(2)}ms</p>
    <p>Total Time: ${(result.decompressTime + result.parseTime).toFixed(
      2
    )}ms</p>
  `;
  resultsList.appendChild(resultElement);
}

async function runAllTests() {
  resultsList.innerHTML = "";

  for (const [dict, formats] of Object.entries(TEST_FILES)) {
    const dictHeader = document.createElement("div");
    dictHeader.className = "dict-header";
    dictHeader.innerHTML = `<h2>${dict}</h2>`;
    resultsList.appendChild(dictHeader);

    for (const [format, url] of Object.entries(formats)) {
      try {
        const result = await fetchAndDecompress(url, format);
        displayResult(result);
      } catch (error) {
        console.error(`Error processing ${dict} ${format}:`, error);
        const errorElement = document.createElement("div");
        errorElement.className = "result-item error";
        errorElement.innerHTML = `
          <h3>${dict} - ${format}</h3>
          <p class="error-message">Error: ${
            error instanceof Error ? error.message : "Unknown error"
          }</p>
        `;
        resultsList.appendChild(errorElement);
      }
    }
  }
}

runAllButton.addEventListener("click", runAllTests);
clearButton.addEventListener("click", () => {
  resultsList.innerHTML = "";
});
