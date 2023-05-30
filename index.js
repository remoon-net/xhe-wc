require("./wrtc-polyfill");
require("./go_wasm_js/wasm_exec_node");
const go = new Go();

const path = require("path");
const wasmUrl = path.join(module.path, "xhe-wc.wasm");
const b = fs.readFileSync(wasmUrl);

exports.XheConnectInit = Promise.resolve(1)
  .then(() => WebAssembly.instantiate(b, go.importObject))
  .then(({ instance }) => {
    return { process: go.run(instance) };
  });
