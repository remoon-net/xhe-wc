interface Device {
  PrivateKey: string;
  ListenPort?: number; // not work
  Address: string;
}

interface Peer {
  PublicKey: string;
  AllowedIPs: string[];
  PresharedKey?: string;
  Endpoint?: string;
  PersistentKeepalive?: string;
}

interface Config extends Device {
  Peers: Peer[];
}

interface Options {
  signaler?: string;
  ices?: string;
  logger?: string;
}

interface XheConnect {
  (config: Config, options?: Options): Promise<XheWireguard>;
}

interface XheWireguard {
  ListenTCP(port?: number): Promise<TCPServer>;
  IpcGet(): Promise<string>;
}

interface TCPServer {
  Serve(): Promise<void>;
  Close(): Promise<void>;
  ServeReady(): boolean;
  ReverseProxy(path: string, remote: string): Promise<void>;
}

declare global {
  var XheConnect: XheConnect;
}

export const XheConnectInit: Promise<any>;
