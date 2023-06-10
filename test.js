const { XheConnectInit } = require("./");
const { v4: uuidv4 } = require("uuid");

const device_id = uuidv4();

const server = require("http")
  .createServer((req, res) => {
    res.end("hello world");
  })
  .listen(0);

let port = server.address().port;

XheConnectInit().then(async () => {
  let xhe = await XheConnect({
    PrivateKey: "CFp1j2epz2sUJ8ovPeUgiStto8EOAAnmvGWp+SSECmI=",
    Addrs: ["192.168.4.1/24"],
    Link: `https://test:test@signaler.slive.fun?t=${device_id}`,
    Peers: [
      {
        PublicKey: "TUpLBfCKwL0joxi+nOsE3+wgxjhlIZdtlcftup/lRik=",
        AllowedIPs: ["192.168.4.2/32"],
      },
    ],
  });
  let server = await xhe.ListenTCP();
  server.Serve().catch(() => {
    // donothing
  });
  if (!server.ServeReady()) {
    throw new Error("server is not ready");
  }
  await server.ReverseProxy("/", `http://127.0.0.1:${port}/`);
  await server.HandleEval("/xhe-eval");

  let code = await new Promise((rl, rj) => {
    const g = require("child_process").spawn(
      `go`,
      [`run`, `./golang_call`, `-id`, device_id],
      {
        stdio: "inherit",
      }
    );
    g.on("exit", (code) => {
      rl(code);
    });
    g.on("error", (err) => {
      rj(err);
    });
  });
  process.exit(code);
});

// hanlde performance.markResourceTiming is not a function
process.on("uncaughtException", (err) => {
  // console.error(err);
});
