<script>
  import { SelectDocument, SelectSaveFile, ConvertTxtToPDF, ConvertHtmlToPDF, ConvertDocToPDF } from '../../wailsjs/go/main/App.js'

  let inputPath = ''
  let inputName = ''
  let inputType = ''
  let converting = false
  let message = ''

  function detectType(path) {
    const ext = path.split('.').pop().toLowerCase()
    if (ext === 'txt') return 'txt'
    if (ext === 'html' || ext === 'htm') return 'html'
    if (ext === 'doc' || ext === 'docx') return 'word'
    return ''
  }

  async function selectFile(type) {
    try {
      const path = await SelectDocument(type)
      if (!path) return
      inputPath = path
      inputName = path.split(/[/\\]/).pop()
      inputType = detectType(path)
      message = ''
    } catch (e) {
      message = '选择文件失败: ' + e
    }
  }

  async function convert() {
    if (!inputPath) { message = '请先选择文件'; return }

    try {
      const baseName = inputName.replace(/\.[^.]+$/, '')
      const output = await SelectSaveFile(baseName + '.pdf', 'PDF文件', '*.pdf')
      if (!output) return

      converting = true
      message = '正在转换...'

      switch (inputType) {
        case 'txt':
          await ConvertTxtToPDF(inputPath, output)
          break
        case 'html':
          await ConvertHtmlToPDF(inputPath, output)
          break
        case 'word':
          await ConvertDocToPDF(inputPath, output)
          break
        default:
          throw '不支持的文件格式'
      }
      message = '转换完成! 已保存到: ' + output
    } catch (e) {
      message = '转换失败: ' + e
    } finally {
      converting = false
    }
  }

  const typeLabels = { txt: 'TXT 文本', html: 'HTML 网页', word: 'Word 文档' }
</script>

<div class="container">
  <div class="toolbar">
    <button class="btn" on:click={() => selectFile('txt')}>选择TXT</button>
    <button class="btn" on:click={() => selectFile('html')}>选择HTML</button>
    <button class="btn" on:click={() => selectFile('word')}>选择Word</button>
  </div>

  <div class="main-area">
    {#if inputPath}
      <div class="file-card">
        <div class="file-icon">
          {#if inputType === 'txt'}📝{:else if inputType === 'html'}🌐{:else if inputType === 'word'}📄{:else}📎{/if}
        </div>
        <div class="file-detail">
          <div class="file-name">{inputName}</div>
          <div class="file-type">{typeLabels[inputType] || '未知格式'}</div>
        </div>
        <div class="arrow">→</div>
        <div class="output-icon">📕</div>
        <div class="output-label">PDF</div>
      </div>

      {#if inputType === 'word'}
        <div class="notice">
          Word 转 PDF 需要安装 LibreOffice
        </div>
      {/if}

      <button class="btn success convert-btn" on:click={convert} disabled={converting}>
        {converting ? '转换中...' : '转换为 PDF'}
      </button>
    {:else}
      <div class="placeholder">
        <div class="icon">📑</div>
        <p>选择 TXT、HTML 或 Word 文件转换为 PDF</p>
        <p class="hint">支持 .txt .html .htm .doc .docx</p>
      </div>
    {/if}
  </div>

  {#if message}
    <div class="message" class:error={message.includes('失败')}>{message}</div>
  {/if}
</div>

<style>
  .container { display: flex; flex-direction: column; height: 100%; gap: 12px; }
  .toolbar { display: flex; gap: 10px; flex-shrink: 0; }
  .btn {
    padding: 8px 18px; border: none; border-radius: 6px;
    cursor: pointer; font-size: 14px; font-family: inherit;
    background: #2a2a4a; color: #ccc; transition: all 0.2s;
  }
  .btn:hover:not(:disabled) { background: #3a3a5a; }
  .btn:disabled { opacity: 0.5; cursor: not-allowed; }
  .btn.success { background: #1b6b3a; color: #fff; }
  .btn.success:hover:not(:disabled) { background: #238a4a; }
  .main-area {
    flex: 1; display: flex; flex-direction: column; align-items: center;
    justify-content: center; background: #16213e; border-radius: 10px; padding: 30px;
  }
  .file-card {
    display: flex; align-items: center; gap: 16px;
    background: #1a2744; padding: 24px 32px; border-radius: 12px; margin-bottom: 20px;
  }
  .file-icon, .output-icon { font-size: 36px; }
  .file-detail { min-width: 0; }
  .file-name { font-size: 16px; word-break: break-all; }
  .file-type { font-size: 12px; color: #888; margin-top: 4px; }
  .arrow { font-size: 24px; color: #e94560; margin: 0 8px; }
  .output-label { font-size: 16px; color: #e94560; font-weight: bold; }
  .notice {
    font-size: 12px; color: #d4a843; background: #2a2520; padding: 8px 14px;
    border-radius: 6px; margin-bottom: 16px;
  }
  .convert-btn { padding: 12px 40px; font-size: 16px; }
  .placeholder {
    display: flex; flex-direction: column; align-items: center; color: #666;
  }
  .placeholder .icon { font-size: 48px; margin-bottom: 12px; }
  .placeholder p { margin: 4px 0; }
  .placeholder .hint { font-size: 12px; color: #555; }
  .message {
    padding: 10px 14px; border-radius: 6px;
    background: #1b3a2a; color: #6fda8a; font-size: 13px; flex-shrink: 0;
  }
  .message.error { background: #3a1b1b; color: #da6f6f; }
</style>
