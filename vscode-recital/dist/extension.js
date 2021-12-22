/******/ (() => { // webpackBootstrap
/******/ 	"use strict";
/******/ 	var __webpack_modules__ = ([
/* 0 */,
/* 1 */
/***/ ((module) => {

module.exports = require("child_process");

/***/ }),
/* 2 */
/***/ ((module) => {

module.exports = require("fs");

/***/ }),
/* 3 */
/***/ ((module) => {

module.exports = require("vscode");

/***/ })
/******/ 	]);
/************************************************************************/
/******/ 	// The module cache
/******/ 	var __webpack_module_cache__ = {};
/******/ 	
/******/ 	// The require function
/******/ 	function __webpack_require__(moduleId) {
/******/ 		// Check if module is in cache
/******/ 		var cachedModule = __webpack_module_cache__[moduleId];
/******/ 		if (cachedModule !== undefined) {
/******/ 			return cachedModule.exports;
/******/ 		}
/******/ 		// Create a new module (and put it into the cache)
/******/ 		var module = __webpack_module_cache__[moduleId] = {
/******/ 			// no module.id needed
/******/ 			// no module.loaded needed
/******/ 			exports: {}
/******/ 		};
/******/ 	
/******/ 		// Execute the module function
/******/ 		__webpack_modules__[moduleId](module, module.exports, __webpack_require__);
/******/ 	
/******/ 		// Return the exports of the module
/******/ 		return module.exports;
/******/ 	}
/******/ 	
/************************************************************************/
var __webpack_exports__ = {};
// This entry need to be wrapped in an IIFE because it need to be isolated against other modules in the chunk.
(() => {
var exports = __webpack_exports__;

Object.defineProperty(exports, "__esModule", ({ value: true }));
exports.deactivate = exports.activate = void 0;
const child_process_1 = __webpack_require__(1);
const fs_1 = __webpack_require__(2);
const vscode = __webpack_require__(3);
function activate(context) {
    let diagnosticCollection = vscode.languages.createDiagnosticCollection('recital');
    ;
    let disposable = vscode.commands.registerCommand('donelli-recital.run', () => {
        const uri = vscode.window.activeTextEditor?.document.uri.fsPath;
        (0, child_process_1.execSync)(`E:\\Projects\\recital\\rt.exe parse ${uri} -json E:\\Projects\\recital\\out.json`);
        (0, fs_1.readFile)('E:\\Projects\\recital\\out.json', 'utf-8', (err, data) => {
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
            diagnosticCollection.set(vscode.window.activeTextEditor.document.uri, diag);
            console.log(diag);
        });
    });
    context.subscriptions.push(disposable, diagnosticCollection);
}
exports.activate = activate;
// this method is called when your extension is deactivated
function deactivate() { }
exports.deactivate = deactivate;

})();

module.exports = __webpack_exports__;
/******/ })()
;
//# sourceMappingURL=extension.js.map