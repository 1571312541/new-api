/**
 * extension.js 专用压缩脚本
 * 
 * 由于 extension.js 文件过大且包含特殊字符，使用纯 Terser 压缩
 * 不进行混淆（已经是打包后的代码，变量名已经很短）
 * 
 * 使用方法: node build-extension-minify.js
 */

const fs = require('fs');
const path = require('path');

async function main() {
    const { minify } = require('terser');
    
    const filePath = path.join(__dirname, 'out', 'extension.js');
    const backupDir = path.join(__dirname, 'backup');
    
    if (!fs.existsSync(filePath)) {
        console.log('❌ extension.js 不存在');
        return;
    }
    
    console.log('🗜️ 开始压缩 extension.js...\n');
    
    const originalCode = fs.readFileSync(filePath, 'utf8');
    const originalSize = Buffer.byteLength(originalCode, 'utf8');
    
    // 备份
    if (!fs.existsSync(backupDir)) {
        fs.mkdirSync(backupDir, { recursive: true });
    }
    const backupPath = path.join(backupDir, `extension.js.${Date.now()}.bak`);
    fs.writeFileSync(backupPath, originalCode);
    console.log(`💾 已备份到: ${backupPath}`);
    
    try {
        // 使用 Terser 压缩，配置更宽松以处理特殊字符
        const result = await minify(originalCode, {
            compress: {
                drop_console: false,
                drop_debugger: true,
                passes: 2,
                dead_code: true,
                unused: true,
            },
            mangle: {
                reserved: ['vscode', 'acquireVsCodeApi', 'exports', 'module', 'require', 'global', 'process'],
                keep_fnames: false,
                keep_classnames: false,
            },
            format: {
                comments: false,
                ascii_only: true,  // 避免 URI 编码问题
            },
            sourceMap: false,
        });
        
        if (result.code) {
            const finalSize = Buffer.byteLength(result.code, 'utf8');
            fs.writeFileSync(filePath, result.code);
            
            const reduction = ((1 - finalSize / originalSize) * 100).toFixed(1);
            console.log(`\n✅ 压缩完成!`);
            console.log(`   原始大小: ${formatSize(originalSize)}`);
            console.log(`   压缩后:   ${formatSize(finalSize)}`);
            console.log(`   减少:     ${reduction}%`);
        } else {
            console.log('❌ 压缩失败: 无输出');
        }
    } catch (error) {
        console.error(`❌ 压缩失败: ${error.message}`);
        
        // 如果 Terser 也失败，尝试简单的空白压缩
        console.log('\n🔄 尝试简单压缩...');
        try {
            // 移除多余空白和注释
            let simpleMinified = originalCode
                .replace(/\/\*[\s\S]*?\*\//g, '')  // 移除块注释
                .replace(/\/\/.*$/gm, '')           // 移除行注释
                .replace(/\s+/g, ' ')               // 压缩空白
                .replace(/\s*([{};,:])\s*/g, '$1'); // 移除符号周围空白
            
            const finalSize = Buffer.byteLength(simpleMinified, 'utf8');
            fs.writeFileSync(filePath, simpleMinified);
            
            const reduction = ((1 - finalSize / originalSize) * 100).toFixed(1);
            console.log(`✅ 简单压缩完成!`);
            console.log(`   原始大小: ${formatSize(originalSize)}`);
            console.log(`   压缩后:   ${formatSize(finalSize)}`);
            console.log(`   减少:     ${reduction}%`);
        } catch (e) {
            console.error(`❌ 简单压缩也失败: ${e.message}`);
        }
    }
}

function formatSize(bytes) {
    if (bytes < 1024) return bytes + ' B';
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB';
    return (bytes / (1024 * 1024)).toFixed(2) + ' MB';
}

main().catch(console.error);

