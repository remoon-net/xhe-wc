import "./go_wasm_js/wasm_exec";
import { XheConnectInit as RawXheConnectInit } from "./index";
import defaultWasmUrl from "./xhe-wc.wasm?url";

export const XheConnectInit = (wasmUrl = defaultWasmUrl) => {
  return RawXheConnectInit(wasmUrl);
};
