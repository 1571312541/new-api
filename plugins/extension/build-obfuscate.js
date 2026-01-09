/**
 * 插件代码混淆加密和压缩脚本
 * 
 * 使用方法:
 * 1. 安装依赖: npm install javascript-obfuscator terser html-minifier-terser cheerio
 * 2. 运行: node build-obfuscate.js
 */

const fs = require('fs');
const path = require('path');

// 动态导入模块
async function main() {
    const JavaScriptObfuscator = require('javascript-obfuscator');
    const { minify: terserMinify } = require('terser');
    const { minify: htmlMinify } = require('html-minifier-terser');
    const cheerio = require('cheerio');

    // 配置
    const config = {
        // JavaScript 混淆配置 (高强度)
        obfuscator: {
            compact: true,
            controlFlowFlattening: true,
            controlFlowFlatteningThreshold: 0.75,
            deadCodeInjection: true,
            deadCodeInjectionThreshold: 0.4,
            debugProtection: false, // VSCode 插件不建议开启
            disableConsoleOutput: false, // 保留 console 用于调试
            identifierNamesGenerator: 'hexadecimal',
            log: false,
            numbersToExpressions: true,
            renameGlobals: false, // 不重命名全局变量，避免破坏 VSCode API
            selfDefending: false, // VSCode 环境不需要
            simplify: true,
            splitStrings: true,
            splitStringsChunkLength: 10,
            stringArray: true,
            stringArrayCallsTransform: true,
            stringArrayEncoding: ['base64'],
            stringArrayIndexShift: true,
            stringArrayRotate: true,
            stringArrayShuffle: true,
            stringArrayWrappersCount: 2,
            stringArrayWrappersChainedCalls: true,
            stringArrayWrappersParametersMaxCount: 4,
            stringArrayWrappersType: 'function',
            stringArrayThreshold: 0.75,
            transformObjectKeys: true,
            unicodeEscapeSequence: false
        },
        // Terser 压缩配置
        terser: {
            compress: {
                drop_console: false, // 保留 console
                drop_debugger: true,
                passes: 2
            },
            mangle: {
                reserved: ['vscode', 'acquireVsCodeApi', 'exports', 'module', 'require']
            },
            format: {
                comments: false
            }
        },
        // HTML 压缩配置
        htmlMinifier: {
            collapseWhitespace: true,
            removeComments: true,
            removeRedundantAttributes: true,
            removeScriptTypeAttributes: true,
            removeStyleLinkTypeAttributes: true,
            useShortDoctype: true,
            minifyCSS: true,
            minifyJS: false // 我们单独处理 JS
        }
    };

    // 文件路径
    const files = {
        extensionJs: path.join(__dirname, 'out', 'extension.js'),
        customPanelHtml: path.join(__dirname, 'out', 'custom-panel.html'),
        customPanelHtmlSrc: path.join(__dirname, 'common-webviews', 'custom-panel.html')
    };

    // 备份目录
    const backupDir = path.join(__dirname, 'backup');
    if (!fs.existsSync(backupDir)) {
        fs.mkdirSync(backupDir, { recursive: true });
    }

    console.log('🔐 开始混淆加密和压缩...\n');

    // 1. 处理 extension.js
    await processExtensionJs(files.extensionJs, config, backupDir, JavaScriptObfuscator, terserMinify);

    // 2. 处理 custom-panel.html (out 目录)
    await processHtmlFile(files.customPanelHtml, config, backupDir, JavaScriptObfuscator, terserMinify, htmlMinify, cheerio);

    // 3. 处理 custom-panel.html (common-webviews 目录)
    await processHtmlFile(files.customPanelHtmlSrc, config, backupDir, JavaScriptObfuscator, terserMinify, htmlMinify, cheerio);

    console.log('\n✅ 混淆加密和压缩完成！');
    console.log(`📁 原始文件已备份到: ${backupDir}`);
}

async function processExtensionJs(filePath, config, backupDir, JavaScriptObfuscator, terserMinify) {
    if (!fs.existsSync(filePath)) {
        console.log(`⚠️ 文件不存在: ${filePath}`);
        return;
    }

    console.log(`📄 处理 extension.js...`);
    const originalCode = fs.readFileSync(filePath, 'utf8');
    const originalSize = Buffer.byteLength(originalCode, 'utf8');

    // 备份
    const backupPath = path.join(backupDir, `extension.js.${Date.now()}.bak`);
    fs.writeFileSync(backupPath, originalCode);
    console.log(`   💾 已备份到: ${backupPath}`);

    try {
        // 先用 Terser 压缩
        console.log('   🗜️ 压缩中...');
        const minified = await terserMinify(originalCode, config.terser);
        
        // 再用 javascript-obfuscator 混淆
        console.log('   🔒 混淆中...');
        const obfuscated = JavaScriptObfuscator.obfuscate(minified.code, config.obfuscator);
        
        const finalCode = obfuscated.getObfuscatedCode();
        const finalSize = Buffer.byteLength(finalCode, 'utf8');

        fs.writeFileSync(filePath, finalCode);
        console.log(`   ✅ 完成! ${formatSize(originalSize)} → ${formatSize(finalSize)} (${((1 - finalSize/originalSize) * 100).toFixed(1)}% 减少)`);
    } catch (error) {
        console.error(`   ❌ 处理失败: ${error.message}`);
    }
}

async function processHtmlFile(filePath, config, backupDir, JavaScriptObfuscator, terserMinify, htmlMinify, cheerio) {
    if (!fs.existsSync(filePath)) {
        console.log(`⚠️ 文件不存在: ${filePath}`);
        return;
    }

    const fileName = path.basename(filePath);
    const dirName = path.basename(path.dirname(filePath));
    console.log(`📄 处理 ${dirName}/${fileName}...`);

    const originalHtml = fs.readFileSync(filePath, 'utf8');
    const originalSize = Buffer.byteLength(originalHtml, 'utf8');

    // 备份
    const backupPath = path.join(backupDir, `${dirName}-${fileName}.${Date.now()}.bak`);
    fs.writeFileSync(backupPath, originalHtml);
    console.log(`   💾 已备份到: ${backupPath}`);

    try {
        const $ = cheerio.load(originalHtml, { decodeEntities: false });

        // 处理所有内联 script 标签
        const scripts = $('script:not([src])');
        console.log(`   📝 发现 ${scripts.length} 个内联脚本`);

        for (let i = 0; i < scripts.length; i++) {
            const script = $(scripts[i]);
            const jsCode = script.html();
            
            if (jsCode && jsCode.trim().length > 50) {
                try {
                    // 压缩
                    const minified = await terserMinify(jsCode, {
                        ...config.terser,
                        mangle: {
                            reserved: ['vscode', 'acquireVsCodeApi', 'postMessage', 'addEventListener']
                        }
                    });
                    
                    // 混淆 (对 HTML 内的 JS 使用较轻的混淆)
                    const lightObfuscatorConfig = {
                        ...config.obfuscator,
                        controlFlowFlattening: false,
                        deadCodeInjection: false,
                        splitStrings: false,
                        stringArrayThreshold: 0.5
                    };
                    
                    const obfuscated = JavaScriptObfuscator.obfuscate(minified.code, lightObfuscatorConfig);
                    script.html(obfuscated.getObfuscatedCode());
                } catch (e) {
                    console.log(`   ⚠️ 脚本 ${i + 1} 混淆失败，保持原样`);
                }
            }
        }

        // 压缩 HTML
        let finalHtml = $.html();
        finalHtml = await htmlMinify(finalHtml, config.htmlMinifier);
        
        const finalSize = Buffer.byteLength(finalHtml, 'utf8');
        fs.writeFileSync(filePath, finalHtml);
        console.log(`   ✅ 完成! ${formatSize(originalSize)} → ${formatSize(finalSize)} (${((1 - finalSize/originalSize) * 100).toFixed(1)}% 减少)`);
    } catch (error) {
        console.error(`   ❌ 处理失败: ${error.message}`);
    }
}

function formatSize(bytes) {
    if (bytes < 1024) return bytes + ' B';
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB';
    return (bytes / (1024 * 1024)).toFixed(2) + ' MB';
}

main().catch(console.error);

