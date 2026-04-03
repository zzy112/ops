<script>
  import { SelectPDF, SelectSaveFile, SelectOutputDir, ConvertPDFToTxt, ConvertPDFToHtml, ConvertPDFToDocx, GetPDFPageCount } from '../../wailsjs/go/main/App.js'

  let pdfPath = ''
  let pdfName = ''
  let pageCount = 0
  let outputFormat = 'txt'
  let converting = false
  let message = ''

  async function selectPDF() {
    try {
      const path = await SelectPDF()
      if (!path) return
      pdfPath = path
      pdfName = path.split(/[/\\]/).pop()
      pageCount = await GetPDFPageCount(path)
      message = ''
    } catch (e) {
      message = '打开PDF失败: ' + e
    }
  }

  async function convert() {
    if (!pdfPath) { message = '请先选择PDF文件'; return }

    try {
      const baseName = pdfName.replace(/\.pdf$/i, '')
      let output = ''

      switch (outputFormat) {
        case 'txt':
          output = await SelectSaveFile(baseName + '.txt', '文本文件', '*.txt')
          if (!output) return
          converting = true; message = '正在转换...'
          await ConvertPDFToTxt(pdfPath, output)
          break
        case 'html':
          output = await SelectSaveFile(baseName + '.html', 'HTML文件', '*.html')
          if (!output) return
          converting = true; message = '正在转换...'
          await ConvertPDFToHtml(pdfPath, output)
          break
        case 'docx':
          output = await SelectSaveFile(baseName + '.docx', 'Word文档', '*.docx')
          if (!output) return
          converting = true; message = '正在转换...'
          await ConvertPDFToDocx(pdfPath, output)
          break
      }
      message = '转换完成! 已保存到: ' + output
    } catch (e) {
      message = '转换失败: ' + e
    } finally {
      converting = false
    }
  }

  const formatLabels = { txt: 'TXT 文本', html: 'HTML 网页', docx: 'Word 文档' }
  const formatIcons = { txt: '📝', html: '🌐', docx: '📄' }
</script>

<div class="container">
  <div class="toolbar">
    <button class="btn primary" on:click={selectPDF}>选择PDF</button>
    {#if pdfName}
      <span class="pdf-name">{pdfName}</span>
      <span class="page-count">({pageCount} 页)</span>
    {/if}
  </div>

  <div class="main-area">
    {#if pdfPath}
      <div class="format-section">
        <div class="section-label">选择输出格式</div>
        <div class="format-cards">
          {#each ['txt', 'html', 'docx'] as fmt}
            <button class="format-card" class:selected={outputFormat === fmt}
              on:click={() => outputFormat = fmt}>
              <div class="fmt-icon">{formatIcons[fmt]}</div>
              <div class="fmt-label">{formatLabels[fmt]}</div>
            </button>
          {/each}
        </div>
      </div>

      {#if outputFormat === 'docx'}
        <div class="notice">
          PDF 转 Word 仅提取文本内容，不保留原始排版
        </div>
      {/if}

      <button class="btn success convert-btn" on:click={convert} disabled={converting}>
        {converting ? '转换中...' : '开始转换'}
      </button>
    {:else}
      <div class="placeholder">
        <div class="icon">📕</div>
        <p>选择 PDF 文件转换为其他格式</p>
        <p class="hint">支持转换为 TXT、HTML、Word</p>
      </div>
    {/if}
  </div>

  {#if message}
    <div class="message" class:error={message.includes('失败')}>{message}</div>
  {/if}
</div>

<style>
  .container { display: flex; flex-direction: column; height: 100%; gap: 12px; }
  .toolbar { display: flex; gap: 10px; align-items: center; flex-shrink: 0; }
  .pdf-name { color: #aaa; font-size: 14px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; max-width: 300px; }
  .page-count { color: #666; font-size: 13px; }
  .btn {
    padding: 8px 18px; border: none; border-radius: 6px;
    cursor: pointer; font-size: 14px; font-family: inherit;
    background: #2a2a4a; color: #ccc; transition: all 0.2s;
  }
  .btn:hover:not(:disabled) { background: #3a3a5a; }
  .btn:disabled { opacity: 0.5; cursor: not-allowed; }
  .btn.primary { background: #0f3460; color: #fff; }
  .btn.primary:hover:not(:disabled) { background: #1a4a7a; }
  .btn.success { background: #1b6b3a; color: #fff; }
  .btn.success:hover:not(:disabled) { background: #238a4a; }
  .main-area {
    flex: 1; display: flex; flex-direction: column; align-items: center;
    justify-content: center; background: #16213e; border-radius: 10px; padding: 30px;
  }
  .format-section { text-align: center; margin-bottom: 20px; }
  .section-label { font-size: 14px; color: #888; margin-bottom: 14px; }
  .format-cards { display: flex; gap: 16px; }
  .format-card {
    width: 110px; padding: 20px 16px; border-radius: 10px;
    background: #1a2744; border: 2px solid transparent;
    cursor: pointer; transition: all 0.2s; text-align: center;
    color: #ccc; font-family: inherit;
  }
  .format-card:hover { background: #1e2f52; border-color: #333; }
  .format-card.selected { border-color: #e94560; background: #1e2040; }
  .fmt-icon { font-size: 32px; margin-bottom: 8px; }
  .fmt-label { font-size: 13px; }
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
