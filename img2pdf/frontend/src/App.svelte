<script>
  let activeTab = 'img2pdf'
  import Img2Pdf from './lib/Img2Pdf.svelte'
  import Pdf2Img from './lib/Pdf2Img.svelte'
  import Doc2Pdf from './lib/Doc2Pdf.svelte'
  import Pdf2Doc from './lib/Pdf2Doc.svelte'
</script>

<main>
  <div class="tabs">
    <button class="tab" class:active={activeTab === 'img2pdf'} on:click={() => activeTab = 'img2pdf'}>
      图片转PDF
    </button>
    <button class="tab" class:active={activeTab === 'pdf2img'} on:click={() => activeTab = 'pdf2img'}>
      PDF转图片
    </button>
    <button class="tab" class:active={activeTab === 'doc2pdf'} on:click={() => activeTab = 'doc2pdf'}>
      文档转PDF
    </button>
    <button class="tab" class:active={activeTab === 'pdf2doc'} on:click={() => activeTab = 'pdf2doc'}>
      PDF转文档
    </button>
  </div>
  <div class="content">
    {#if activeTab === 'img2pdf'}
      <Img2Pdf />
    {:else if activeTab === 'pdf2img'}
      <Pdf2Img />
    {:else if activeTab === 'doc2pdf'}
      <Doc2Pdf />
    {:else}
      <Pdf2Doc />
    {/if}
  </div>
</main>

<style>
  :global(body) {
    margin: 0;
    padding: 0;
    font-family: "Microsoft YaHei", "PingFang SC", sans-serif;
    background: #1a1a2e;
    color: #eee;
    overflow: hidden;
    height: 100vh;
  }
  main {
    display: flex;
    flex-direction: column;
    height: 100vh;
  }
  .tabs {
    display: flex;
    background: #16213e;
    border-bottom: 2px solid #0f3460;
    padding: 0;
    flex-shrink: 0;
  }
  .tab {
    flex: 1;
    padding: 14px 0;
    border: none;
    background: transparent;
    color: #8899aa;
    font-size: 15px;
    cursor: pointer;
    transition: all 0.2s;
    font-family: inherit;
  }
  .tab:hover {
    color: #ccc;
    background: #1a2744;
  }
  .tab.active {
    color: #e94560;
    border-bottom: 3px solid #e94560;
    background: #1a1a2e;
  }
  .content {
    flex: 1;
    overflow: auto;
    padding: 20px;
  }
</style>
