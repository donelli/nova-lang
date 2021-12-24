import { exec, execSync } from 'child_process';
import { open, openSync, readFile, readSync } from 'fs';
import * as vscode from 'vscode';

let diagnosticCollection: vscode.DiagnosticCollection;

export function activate(context: vscode.ExtensionContext) {
	
	diagnosticCollection = vscode.languages.createDiagnosticCollection('recital');
	
	let disposableOpen = vscode.workspace.onDidOpenTextDocument(document => {
		loadDiagnostics(document);
	});

	let disposableChange = vscode.workspace.onDidSaveTextDocument(document => {
		loadDiagnostics(document);
	});
	
	context.subscriptions.push(disposableOpen, disposableChange, diagnosticCollection);
}

export function deactivate() {}

function loadDiagnostics(doc: vscode.TextDocument) {

	if (doc.languageId !== 'recital') {
		return;
	}
	
	diagnosticCollection.set(doc.uri, undefined);
	
	const uri = doc.uri.fsPath;
	
	execSync(`E:\\Projects\\recital\\rt.exe parse ${uri} -json E:\\Projects\\recital\\out.json`);
	
	readFile('E:\\Projects\\recital\\out.json', 'utf-8' , (err, data) => {
		if (err) {
			vscode.window.showErrorMessage(err.message);
			return;
		}

		const diagnostic = JSON.parse(data);
		const diag = [];
		
		console.log(diagnostic);
		
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

		diagnosticCollection.set(doc.uri, diag);

		console.log(diag);
		
		});
	
}
