require("./go_wasm_js/wasm_exec");

const go = new Go();

if (!WebAssembly.instantiateStreaming) {
  WebAssembly.instantiateStreaming = async (resp, importObject) => {
    const source = await (await resp).arrayBuffer();
    return await WebAssembly.instantiate(source, importObject);
  };
}

exports.XheConnectInit = async (
  wasmUrl = "https://unpkg.com/xhe-wc/xhe-wc.wasm"
) => {
  const { instance } = await WebAssembly.instantiateStreaming(
    fetch(wasmUrl),
    go.importObject
  );
  return { process: go.run(instance) };
};
