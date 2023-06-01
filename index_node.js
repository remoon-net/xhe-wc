require("./wrtc-polyfill");
require("./go_wasm_js/wasm_exec_node");
const fs = require("node:fs/promises");
const path = require("path");
const defaultWasmPath = path.join(module.path, "xhe-wc.wasm");

const { XheConnectInit } = require("./index.cjs");

exports.XheConnectInit = async (wasmPath = defaultWasmPath) => {
  const wasmUrl = await getWasmUrl(wasmPath);
  return XheConnectInit(wasmUrl);
};

/**
 * @param {string} wasmPath
 */
async function getWasmUrl(wasmPath) {
  let buf = await fs.readFile(wasmPath);
  return "data:application/wasm;base64," + buf.toString("base64");
}
