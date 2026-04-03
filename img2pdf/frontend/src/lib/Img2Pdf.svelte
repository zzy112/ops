<script>
  import { SelectImages, SelectOutputFile, ConvertImagesToPDF, GetImageInfo } from '../../wailsjs/go/main/App.js'
  import { EventsOn } from '../../wailsjs/runtime/runtime.js'

  let files = []
  let pageSize = 'A4'
  let converting = false
  let message = ''
  let dragOver = false

  async function addImages() {
    try {
      const selected = await SelectImages()
      if (selected && selected.length > 0) {
        files = [...files, ...selected]
      }
    } catch (e) {
      message = '选择文件失败: ' + e
    }
  }

  function removeFile(index) {
    files = files.filter((_, i) => i !== index)
  }

  function moveUp(index) {
    if (index <= 0) return
    const arr = [...files]
    ;[arr[index - 1], arr[index]] = [arr[index], arr[index - 1]]
    files = arr
  }

  function moveDown(index) {
    if (index >= files.length - 1) return
    const arr = [...files]
    ;[arr[index], arr[index + 1]] = [arr[index + 1], arr[index]]
    files = arr
  }

  async function convert() {
    if (files.length === 0) {
      message = '请先添加图片'
      return
    }
    try {
      const output = await SelectOutputFile()
      if (!output) return
      converting = true
      message = '正在转换...'
      const paths = files.map(f => f.path)
      await ConvertImagesToPDF(paths, output, pageSize)
      message = '转换完成! 已保存到: ' + output
    } catch (e) {
      message = '转换失败: ' + e
    } finally {
      converting = false
    }
  }

  function formatSize(bytes) {
    if (bytes < 1024) return bytes + ' B'
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
    return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
  }

  async function handleDrop(e) {
    e.preventDefault()
    dragOver = false
    const items = e.dataTransfer.files
    for (let i = 0; i < items.length; i++) {
      const path = items[i].path || items[i].name
      if (path) {
        try {
          const info = await GetImageInfo(path)
          if (info) files = [...files, info]
        } catch (err) {
          // skip unsupported files
        }
      }
    }
  }

  function handleDragOver(e) {
    e.preventDefault()
    dragOver = true
  }

  function handleDragLeave() {
    dragOver = false
  }
</script>

<div class="container">
  <div class="toolbar">
    <button class="btn primary" on:click={addImages}>添加图片</button>
    <select bind:value={pageSize}>
      <option value="A4">A4 页面</option>
      <option value="original">原始尺寸</option>
    </select>
    <button class="btn success" on:click={convert} disabled={converting || files.length === 0}>
      {converting ? '转换中...' : '生成PDF'}
    </button>
    <button class="btn" on:click={() => { files = []; message = '' }} disabled={files.length === 0}>清空</button>
  </div>

  {#if message}
    <div class="message" class:error={message.includes('失败')}>{message}</div>
  {/if}

  <div class="drop-zone" class:drag-over={dragOver}
    on:drop={handleDrop} on:dragover={handleDragOver} on:dragleave={handleDragLeave}>
    {#if files.length === 0}
      <div class="placeholder">
        <div class="icon">📁</div>
        <p>拖拽图片到这里，或点击"添加图片"按钮</p>
        <p class="hint">支持 JPG、PNG、BMP、TIFF、WebP、GIF</p>
      </div>
    {:else}
      <div class="file-list">
        {#each files as file, i}
          <div class="file-item">
            <span class="file-index">{i + 1}</span>
            {#if file.thumbnail}
              <img src={file.thumbnail} alt={file.name} class="thumb" />
            {:else}
              <div class="thumb-placeholder">IMG</div>
            {/if}
            <div class="file-info">
              <div class="file-name">{file.name}</div>
              <div class="file-size">{formatSize(file.size)}</div>
            </div>
            <div class="file-actions">
              <button class="btn-sm" on:click={() => moveUp(i)} disabled={i === 0}>↑</button>
              <button class="btn-sm" on:click={() => moveDown(i)} disabled={i === files.length - 1}>↓</button>
              <button class="btn-sm danger" on:click={() => removeFile(i)}>✕</button>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </div>
</div>

<style>
  .container { display: flex; flex-direction: column; height: 100%; gap: 12px; }
  .toolbar { display: flex; gap: 10px; align-items: center; flex-shrink: 0; }
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
  select {
    padding: 8px 12px; border-radius: 6px; border: 1px solid #333;
    background: #2a2a4a; color: #ccc; font-size: 14px; font-family: inherit;
  }
  .message {
    padding: 10px 14px; border-radius: 6px;
    background: #1b3a2a; color: #6fda8a; font-size: 13px; flex-shrink: 0;
  }
  .message.error { background: #3a1b1b; color: #da6f6f; }
  .drop-zone {
    flex: 1; border: 2px dashed #333; border-radius: 10px;
    overflow-y: auto; transition: all 0.2s; min-height: 0;
  }
  .drop-zone.drag-over { border-color: #e94560; background: rgba(233,69,96,0.05); }
  .placeholder {
    display: flex; flex-direction: column; align-items: center;
    justify-content: center; height: 100%; color: #666;
  }
  .placeholder .icon { font-size: 48px; margin-bottom: 12px; }
  .placeholder p { margin: 4px 0; }
  .placeholder .hint { font-size: 12px; color: #555; }
  .file-list { padding: 8px; }
  .file-item {
    display: flex; align-items: center; gap: 12px;
    padding: 8px 12px; border-radius: 8px; margin-bottom: 6px;
    background: #16213e; transition: background 0.15s;
  }
  .file-item:hover { background: #1a2744; }
  .file-index { color: #666; font-size: 13px; min-width: 24px; text-align: center; }
  .thumb { width: 50px; height: 50px; object-fit: cover; border-radius: 4px; }
  .thumb-placeholder {
    width: 50px; height: 50px; background: #2a2a4a; border-radius: 4px;
    display: flex; align-items: center; justify-content: center;
    color: #666; font-size: 11px;
  }
  .file-info { flex: 1; min-width: 0; }
  .file-name { font-size: 14px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
  .file-size { font-size: 12px; color: #666; margin-top: 2px; }
  .file-actions { display: flex; gap: 4px; }
  .btn-sm {
    width: 28px; height: 28px; border: none; border-radius: 4px;
    background: #2a2a4a; color: #aaa; cursor: pointer; font-size: 13px;
    display: flex; align-items: center; justify-content: center;
    transition: all 0.15s; padding: 0;
  }
  .btn-sm:hover:not(:disabled) { background: #3a3a5a; color: #fff; }
  .btn-sm:disabled { opacity: 0.3; cursor: not-allowed; }
  .btn-sm.danger:hover:not(:disabled) { background: #5a2020; color: #ff6b6b; }
</style>
