<script>
  import { SelectPDF, SelectOutputDir, ConvertPDFToImages, GetPDFPageCount, PreviewPDFPage } from '../../wailsjs/go/main/App.js'

  let pdfPath = ''
  let pdfName = ''
  let pageCount = 0
  let format = 'jpg'
  let dpi = 150
  let pageRange = 'all'
  let customRange = ''
  let outputDir = ''
  let converting = false
  let message = ''
  let previewSrc = ''
  let previewPage = 0

  async function selectPDF() {
    try {
      const path = await SelectPDF()
      if (!path) return
      pdfPath = path
      pdfName = path.split(/[/\\]/).pop()
      pageCount = await GetPDFPageCount(path)
      message = ''
      previewPage = 0
      await loadPreview()
    } catch (e) {
      message = '打开PDF失败: ' + e
    }
  }

  async function selectOutput() {
    try {
      const dir = await SelectOutputDir()
      if (dir) outputDir = dir
    } catch (e) {
      message = '选择目录失败: ' + e
    }
  }

  async function loadPreview() {
    if (!pdfPath) return
    try {
      previewSrc = await PreviewPDFPage(pdfPath, previewPage)
    } catch (e) {
      previewSrc = ''
    }
  }

  async function prevPage() {
    if (previewPage > 0) {
      previewPage--
      await loadPreview()
    }
  }

  async function nextPage() {
    if (previewPage < pageCount - 1) {
      previewPage++
      await loadPreview()
    }
  }

  async function convert() {
    if (!pdfPath) { message = '请先选择PDF文件'; return }
    if (!outputDir) { message = '请先选择输出目录'; return }
    converting = true
    message = '正在转换...'
    try {
      const range = pageRange === 'custom' ? customRange : 'all'
      const result = await ConvertPDFToImages(pdfPath, outputDir, format, dpi, range)
      message = `转换完成! 共导出 ${result.length} 张图片到: ${outputDir}`
    } catch (e) {
      message = '转换失败: ' + e
    } finally {
      converting = false
    }
  }
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
    <div class="preview-panel">
      {#if previewSrc}
        <img src={previewSrc} alt="PDF预览" class="preview-img" />
        <div class="page-nav">
          <button class="btn-sm" on:click={prevPage} disabled={previewPage <= 0}>←</button>
          <span>第 {previewPage + 1} / {pageCount} 页</span>
          <button class="btn-sm" on:click={nextPage} disabled={previewPage >= pageCount - 1}>→</button>
        </div>
      {:else}
        <div class="preview-placeholder">
          <div class="icon">📄</div>
          <p>选择PDF后预览</p>
        </div>
      {/if}
    </div>

    <div class="settings-panel">
      <div class="setting-group">
        <label>输出格式</label>
        <select bind:value={format}>
          <option value="jpg">JPG</option>
          <option value="png">PNG</option>
        </select>
      </div>

      <div class="setting-group">
        <label>DPI (分辨率)</label>
        <select bind:value={dpi}>
          <option value={72}>72 (低)</option>
          <option value={150}>150 (中)</option>
          <option value={300}>300 (高)</option>
        </select>
      </div>

      <div class="setting-group">
        <label>页码范围</label>
        <select bind:value={pageRange}>
          <option value="all">全部页面</option>
          <option value="custom">自定义</option>
        </select>
        {#if pageRange === 'custom'}
          <input type="text" bind:value={customRange} placeholder="例: 1-5, 8, 10-12" class="range-input" />
        {/if}
      </div>

      <div class="setting-group">
        <label>输出目录</label>
        <div class="dir-row">
          <button class="btn" on:click={selectOutput}>选择目录</button>
          {#if outputDir}
            <span class="dir-path">{outputDir}</span>
          {/if}
        </div>
      </div>

      <button class="btn success full-width" on:click={convert} disabled={converting || !pdfPath}>
        {converting ? '转换中...' : '开始转换'}
      </button>
    </div>
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
  .main-area { display: flex; gap: 16px; flex: 1; min-height: 0; }
  .preview-panel {
    flex: 1; display: flex; flex-direction: column; align-items: center;
    justify-content: center; background: #16213e; border-radius: 10px;
    padding: 16px; min-width: 0;
  }
  .preview-img { max-width: 100%; max-height: calc(100% - 40px); object-fit: contain; border-radius: 4px; }
  .page-nav { display: flex; align-items: center; gap: 12px; margin-top: 10px; color: #aaa; font-size: 13px; }
  .preview-placeholder {
    display: flex; flex-direction: column; align-items: center; color: #555;
  }
  .preview-placeholder .icon { font-size: 48px; margin-bottom: 8px; }
  .settings-panel {
    width: 260px; flex-shrink: 0; display: flex; flex-direction: column; gap: 16px;
  }
  .setting-group { display: flex; flex-direction: column; gap: 6px; }
  .setting-group label { font-size: 13px; color: #888; }
  select, .range-input {
    padding: 8px 12px; border-radius: 6px; border: 1px solid #333;
    background: #2a2a4a; color: #ccc; font-size: 14px; font-family: inherit;
    width: 100%; box-sizing: border-box;
  }
  .range-input { margin-top: 6px; }
  .dir-row { display: flex; gap: 8px; align-items: center; }
  .dir-path { font-size: 12px; color: #666; word-break: break-all; }
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
  .full-width { width: 100%; margin-top: auto; }
  .btn-sm {
    width: 28px; height: 28px; border: none; border-radius: 4px;
    background: #2a2a4a; color: #aaa; cursor: pointer; font-size: 14px;
    display: flex; align-items: center; justify-content: center; padding: 0;
  }
  .btn-sm:hover:not(:disabled) { background: #3a3a5a; color: #fff; }
  .btn-sm:disabled { opacity: 0.3; cursor: not-allowed; }
  .message {
    padding: 10px 14px; border-radius: 6px;
    background: #1b3a2a; color: #6fda8a; font-size: 13px; flex-shrink: 0;
  }
  .message.error { background: #3a1b1b; color: #da6f6f; }
</style>
