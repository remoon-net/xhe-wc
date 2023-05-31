require("./wrtc-polyfill");
require("./go_wasm_js/wasm_exec_node");
const go = new Go();
const fs = require("node:fs/promises");
const path = require("path");
const defaultWasmUrl = path.join(module.path, "xhe-wc.wasm");

exports.XheConnectInit = async (wasmUrl = defaultWasmUrl) => {
  const b = await fs.readFile(wasmUrl);
  const { instance } = await WebAssembly.instantiate(b, go.importObject);
  return { process: go.run(instance) };
};
