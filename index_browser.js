require("./go_wasm_js/wasm_exec");

const go = new Go();

let wasmUrl = globalThis["XheConnectWasm"];
if (!wasmUrl) {
  wasmUrl = "https://unpkg.com/xhe-wc/xhe-wc.wasm";
}

if (!WebAssembly.instantiateStreaming) {
  WebAssembly.instantiateStreaming = async (resp, importObject) => {
    const source = await (await resp).arrayBuffer();
    return await WebAssembly.instantiate(source, importObject);
  };
}

exports.XheConnectInit = Promise.resolve(1)
  .then(() => WebAssembly.instantiateStreaming(fetch(wasmUrl), go.importObject))
  .then(({ instance }) => {
    return { process: go.run(instance) };
  });
