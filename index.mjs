if (!WebAssembly.instantiateStreaming) {
  WebAssembly.instantiateStreaming = async (resp, importObject) => {
    const source = await (await resp).arrayBuffer();
    return await WebAssembly.instantiate(source, importObject);
  };
}

import { version } from "./package.json";
const defaultWasmUrl = `https://unpkg.com/xhe-wc@${version}/xhe-wc.wasm`;

export async function XheConnectInit(wasmUrl = defaultWasmUrl) {
  const go = new Go();
  const { instance } = await WebAssembly.instantiateStreaming(
    fetch(wasmUrl),
    go.importObject
  );
  return { process: go.run(instance) };
}
