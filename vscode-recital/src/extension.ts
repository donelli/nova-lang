import { exec, execSync } from 'child_process';
import { open, openSync, readFile, readSync } from 'fs';
import * as vscode from 'vscode';

export function activate(context: vscode.ExtensionContext) {
	
	let diagnosticCollection: vscode.DiagnosticCollection = vscode.languages.createDiagnosticCollection('recital');;
	
	let disposable = vscode.commands.registerCommand('donelli-recital.run', () => {
		
		const uri = vscode.window.activeTextEditor?.document.uri.fsPath;
		
		execSync(`E:\\Projects\\recital\\rt.exe parse ${uri} -json E:\\Projects\\recital\\out.json`);
		
		readFile('E:\\Projects\\recital\\out.json', 'utf-8' , (err, data) => {
			if (err) {
			  vscode.window.showErrorMessage(err.message);
			  return;
			}

			const diagnostic = JSON.parse(data);
			const diag = [];
			
			console.log(diagnostic);
			
			// diagnosticCollection.clear();
			
			for (let err of diagnostic.errors) {
				let start = new vscode.Position(err.range.start.row, err.range.start.column);
				let end = new vscode.Position(err.range.end.row, err.range.end.column);
				diag.push(new vscode.Diagnostic(new vscode.Range(start, end), err.type + ' ' + err.message, vscode.DiagnosticSeverity.Error));
			}

			for (let warn of diagnostic.warnings) {
				let start = new vscode.Position(warn.range.start.row, warn.range.start.column);
				let end = new vscode.Position(warn.range.end.row, warn.range.end.column);
				diag.push(new vscode.Diagnostic(new vscode.Range(start, end), warn.type + ' ' + warn.message, vscode.DiagnosticSeverity.Warning));
			}

			diagnosticCollection.set(vscode.window.activeTextEditor!.document.uri, diag);

			console.log(diag);
			
		 });
		
	});

	context.subscriptions.push(disposable, diagnosticCollection);
}

// this method is called when your extension is deactivated
export function deactivate() {}
