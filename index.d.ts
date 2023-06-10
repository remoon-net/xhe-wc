export interface Device {
  PrivateKey: string;
  ListenPort?: number; // not work
  /**
   * @deprecated 转而使用 Addrs 字段
   */
  Address?: string;
  Addrs: string[];
  Link?: string;
}

export interface Peer {
  PublicKey: string;
  AllowedIPs: string[];
  PresharedKey?: string;
  Endpoint?: string;
  PersistentKeepalive?: string;
}

export interface Config extends Device {
  Peers: Peer[];
}

export interface Options {
  /**
   * @deprecated 转而使用 Config.Link
   */
  signaler?: string;
  ices?: string;
  logger?: string;
}

export interface XheConnect {
  (config: Config, options?: Options): Promise<XheWireguard>;
}

export interface XheWireguard {
  ListenTCP(port?: number): Promise<TCPServer>;
  IpcGet(): Promise<string>;
}

export interface TCPServer {
  Serve(): Promise<void>;
  Close(): Promise<void>;
  ServeReady(): boolean;
  ReverseProxy(path: string, remote: string): Promise<void>;
  HandleEval(path: string): void;
}

declare global {
  var XheConnect: XheConnect;
}

export const XheConnectInit: (wasmUrl?: string) => Promise<any>;
