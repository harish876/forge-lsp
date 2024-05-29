import * as path from "path";
import { workspace, ExtensionContext } from "vscode";

import {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
  TransportKind,
} from "vscode-languageclient/node";

let client: LanguageClient;


export function activate(context: ExtensionContext) {
  let serverOptions: ServerOptions

  serverOptions = {
    run: {
      command: "/home/harish/personal/forge-lsp/server/main",
      transport: TransportKind.stdio,

    },
    debug: {
      command: "/home/harish/personal/forge-lsp/server/tmp/main",
      transport: TransportKind.stdio,
    }
  };

  const clientOptions: LanguageClientOptions = {
    documentSelector: [
      { scheme: "file", language: "ini" },
      { scheme: "file", language: "python" }
    ],
    synchronize: {
      fileEvents: [
        workspace.createFileSystemWatcher('**/*.py'),
        workspace.createFileSystemWatcher('**/*.ini'),
      ]
    },
  };

  client = new LanguageClient(
    "Forge LSP Client",
    "Forge LSP Server",
    serverOptions,
    clientOptions
  );
  client.start();
}

export function deactivate(): Thenable<void> | undefined {
  if (!client) {
    return undefined;
  }
  return client.stop();
}
