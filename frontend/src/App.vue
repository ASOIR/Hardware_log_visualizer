<template>
  <div class="app-layout">
    <!-- Drop Overlay -->
    <div v-show="isDragging" class="drop-overlay">
      <div class="drop-box">
        <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="drop-icon"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg>
        <h2>Release to open CSV</h2>
      </div>
    </div>

    <!-- Loading Overlay -->
    <div v-if="isLoading" class="loading-overlay">
      <div class="spinner-container">
        <div class="glow-orb loading-orb"></div>
        <div class="spinner"></div>
        <h3>Loading Data</h3>
        <p>Parsing dataset...</p>
      </div>
    </div>

    <!-- Error Modal -->
    <div v-if="errorMessage" class="error-modal-overlay">
      <div class="error-modal-box">
        <div class="error-modal-header">
          <svg class="error-icon" xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
          <h3>Error</h3>
        </div>
        <p class="error-modal-body">{{ errorMessage }}</p>
        <button class="error-modal-btn" @click="errorMessage = ''">Close</button>
      </div>
    </div>

    <!-- Dock Toolbar -->
    <nav class="dock">
      <div class="dock-top">
        <button v-if="parseResult" class="dock-btn" :class="{ active: activeSidebarTab === 'metrics' }" @click="activeSidebarTab = 'metrics'" title="Metrics">
          <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="12 2 2 7 12 12 22 7 12 2"/><polyline points="2 17 12 22 22 17"/><polyline points="2 12 12 17 22 12"/></svg>
        </button>
        <button v-if="parseResult" class="dock-btn" :class="{ active: activeSidebarTab === 'files' }" @click="activeSidebarTab = 'files'" title="Files">
          <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/></svg>
        </button>
      </div>
      <div class="dock-bottom">
        <button class="dock-btn primary" @click="handleManualOpenFile" title="Open File">
          <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg>
        </button>
      </div>
    </nav>

    <aside class="sidebar">
      <div class="sidebar-header">
        <h2>Hardware Viz</h2>
      </div>

      <!-- File Explorer inside Sidebar -->
      <div class="file-explorer" v-if="parseResult && activeSidebarTab === 'files' && parseResult.files.length > 0">
        <div class="explorer-header">
          <h3>Directory Logs</h3>
        </div>
        <div class="explorer-filter">
           <input type="text" v-model="fileSearch" placeholder="Filter files..." />
        </div>
        <div class="file-list">
          <div 
            v-for="file in filteredFiles" 
            :key="file"
            class="file-item" 
            :class="{ active: file === parseResult.loadedFile }"
            @click="openSpecificFile(file)"
            :title="getBasename(file)"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14.5 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7.5L14.5 2z"/><polyline points="14 2 14 8 20 8"/></svg>
            <span>{{ getBasename(file) }}</span>
          </div>
        </div>
      </div>

      <div class="metrics-list" v-if="parseResult && activeSidebarTab === 'metrics' && categorizedMetrics.length > 0">
        <div class="metrics-title-row">
          <h3>Categories</h3>
          <span class="badge">{{ parseResult.headers.length }} total</span>
        </div>
        
        <div class="category" v-for="category in categorizedMetrics" :key="category.name">
          <div class="category-header" @click="toggleCategory(category.name)" :class="{ 'is-open': openCategories[category.name] }">
            <div class="cat-left">
              <svg class="chevron" :class="{'rotate-90': openCategories[category.name]}" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 18 15 12 9 6"/></svg>
              <span>{{ category.name }}</span>
            </div>
            <span class="category-badge">{{ category.indices.length }}</span>
          </div>
          
          <transition name="slide-fade">
            <div class="category-content" v-show="openCategories[category.name]">
              <label 
                class="metric-item" 
                v-for="index in category.indices" 
                :key="index"
                :class="{ selected: selectedIndices.includes(index) }"
              >
                <div class="custom-checkbox">
                  <input type="checkbox" v-model="selectedIndices" :value="index" />
                  <div class="checkmark"></div>
                </div>
                <span class="metric-name" :title="parseResult.headers[index]">{{ parseResult.headers[index] }}</span>
              </label>
            </div>
          </transition>
        </div>
      </div>
      
      <div v-else-if="!parseResult" class="empty-state">
        <div class="empty-icon-wrap">
          <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M14.5 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7.5L14.5 2z"/><polyline points="14 2 14 8 20 8"/></svg>
        </div>
        <p>Drag & Drop a CSV or click Open File</p>
      </div>
    </aside>

    <main class="chart-area-container" ref="chartAreaContainerRef">
      <div class="chart-area" v-if="parseResult">
        <div v-if="selectedIndices.length === 0" class="no-chart-state">
          <div class="placeholder">
            <div class="glow-orb"></div>
            <svg xmlns="http://www.w3.org/2000/svg" width="56" height="56" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round" class="placeholder-icon"><path d="M3 3v18h18"/><path d="m19 9-5 5-4-4-3 3"/></svg>
            <p>Select metrics from the left to view charts</p>
          </div>
        </div>
        
        <div class="mega-chart-container" v-show="selectedIndices.length > 0">
          <v-chart 
            ref="megaChartRef"
            class="mega-chart" 
            :style="{ height: megaChartHeight + 'px' }" 
            :option="megaChartOption" 
            :init-options="{ renderer: 'canvas' }" 
            @datazoom="onMegaZoom"
            @restore="onRestore"
            autoresize 
          />
        </div>
      </div>
      
      <!-- Fixed Global Timeline Scrubber -->
      <div class="global-timeline" v-if="parseResult && selectedIndices.length > 0">
        <v-chart ref="timelineRef" :key="parseResult.loadedFile" class="timeline-chart" :option="timelineOption" @datazoom="onTimelineZoom" @restore="onRestore" autoresize />
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue';
import { use } from 'echarts/core';
import { CanvasRenderer } from 'echarts/renderers';
import { LineChart, CustomChart } from 'echarts/charts';
import { GridComponent, DataZoomComponent, TooltipComponent, TitleComponent, ToolboxComponent } from 'echarts/components';
import VChart from 'vue-echarts';
import './style.css';
// Removed unused echarts default import

use([CanvasRenderer, LineChart, CustomChart, GridComponent, DataZoomComponent, TooltipComponent, TitleComponent, ToolboxComponent]);

declare global {
  interface Window {
    go: any;
    runtime: any;
  }
}

interface ParseResult {
  headers: string[];
  times: string[];
  data: number[][];
  files: string[];
  loadedFile: string;
}

const parseResult = ref<ParseResult | null>(null);
const selectedIndices = ref<number[]>([]);
const openCategories = ref<Record<string, boolean>>({});
const fileSearch = ref("");
const isDragging = ref(false);
const activeSidebarTab = ref<'metrics' | 'files'>('metrics');
const isLoading = ref(false);
const errorMessage = ref("");

const megaChartRef = ref<any>(null);
const timelineRef = ref<any>(null);

let isSyncing = false;

const onMegaZoom = (e: any) => {
  if (isSyncing || !timelineRef.value) return;
  isSyncing = true;
  timelineRef.value.dispatchAction({
    type: 'dataZoom',
    dataZoomIndex: 0,
    start: e.start ?? (e.batch && e.batch.length > 0 ? e.batch[0].start : undefined),
    end: e.end ?? (e.batch && e.batch.length > 0 ? e.batch[0].end : undefined),
    startValue: e.startValue ?? (e.batch && e.batch.length > 0 ? e.batch[0].startValue : undefined),
    endValue: e.endValue ?? (e.batch && e.batch.length > 0 ? e.batch[0].endValue : undefined)
  });
  window.requestAnimationFrame(() => isSyncing = false);
};

const onTimelineZoom = (e: any) => {
  if (isSyncing || !megaChartRef.value) return;
  isSyncing = true;
  megaChartRef.value.dispatchAction({
    type: 'dataZoom',
    dataZoomIndex: 0,
    start: e.start ?? (e.batch && e.batch.length > 0 ? e.batch[0].start : undefined),
    end: e.end ?? (e.batch && e.batch.length > 0 ? e.batch[0].end : undefined),
    startValue: e.startValue ?? (e.batch && e.batch.length > 0 ? e.batch[0].startValue : undefined),
    endValue: e.endValue ?? (e.batch && e.batch.length > 0 ? e.batch[0].endValue : undefined)
  });
  window.requestAnimationFrame(() => isSyncing = false);
};

const onRestore = () => {
  if (isSyncing || !megaChartRef.value || !timelineRef.value) return;
  isSyncing = true;
  megaChartRef.value.dispatchAction({ type: 'restore' });
  timelineRef.value.dispatchAction({ type: 'restore' });
  window.requestAnimationFrame(() => isSyncing = false);
};

const colors = [
  '#f59e0b', '#f97316', '#ea580c', '#ef4444', 
  '#dc2626', '#b91c1c', '#fcd34d', '#fb923c'
];

// --- FILE PATH UTILS ---
const getBasename = (path: string) => {
  return path.split(/[\\/]/).pop() || path;
};

// --- FILE EXPLORER ---
const filteredFiles = computed(() => {
  if (!parseResult.value) return [];
  if (!fileSearch.value) return parseResult.value.files;
  const q = fileSearch.value.toLowerCase();
  return parseResult.value.files.filter(f => getBasename(f).toLowerCase().includes(q));
});

// --- METRIC GROUPS ---
const categorizedMetrics = computed(() => {
  if (!parseResult.value) return [];
  
  const groups: Record<string, number[]> = {
    '處理器相關 (CPU)': [],
    '顯示卡相關 (GPU)': [],
    '系統與主機板記憶體 (System & Memory)': [],
    '散熱、溫度與功耗 (Power & Thermals)': [],
    '磁碟相關 (Disk)': [],
    '網路相關 (Network)': [],
    '其他 (Others)': []
  };

  parseResult.value.headers.forEach((h, idx) => {
    if (/溫度|℃|功耗|功率|\[W\]|\[V\]|電壓|\[A\]|電流|風扇|\[RPM\]|PPT|TDC|EDC|Power|Temp|Thermal|Pump|VID|Vcore/i.test(h)) {
      groups['散熱、溫度與功耗 (Power & Thermals)'].push(idx);
    } else if (/磁碟|讀取活動|寫入活動|總活動|讀取速率|寫入速率|讀取總計|寫入總計|主機寫入|主機讀取/i.test(h)) {
      groups['磁碟相關 (Disk)'].push(idx);
    } else if (/下載|上傳/i.test(h)) {
      groups['網路相關 (Network)'].push(idx);
    } else if (/GPU|顯示記憶體|PCIe 連結速度|影格率|Frame Time/i.test(h)) {
      groups['顯示卡相關 (GPU)'].push(idx);
    } else if (/Core \d|Core C0|Core C1|Core C6|CPU|L3|VDDCR|SOC|FCLK|Infinity/i.test(h) && !/GPU/i.test(h)) {
      groups['處理器相關 (CPU)'].push(idx);
    } else if (/記憶體|核心|Tcas|Trcd|Trp|Tras|Trc|Trfc|Command Rate|UCLK|DRAM|系統|MOS|晶片組|VSB|VTT|VBAT|VREF|SPD|VDD|VPP|VOUT|VIN|PMIC|效能下降|效能限制|PCI Express Error|Error|Replay|DLLP|TLP|LCRC|NAKs|Recovery|錯誤總數|V2|T15|\+12V|\+5V|\+3.3V|3VSB|AVSB/i.test(h)) {
      groups['系統與主機板記憶體 (System & Memory)'].push(idx);
    } else {
      if (/CPU|Core/i.test(h)) groups['處理器相關 (CPU)'].push(idx);
      else if (/GPU|Video|Graphics/i.test(h)) groups['顯示卡相關 (GPU)'].push(idx);
      else if (/Disk|Read|Write/i.test(h)) groups['磁碟相關 (Disk)'].push(idx);
      else if (/Network|LAN|Download|Upload/i.test(h)) groups['網路相關 (Network)'].push(idx);
      else groups['其他 (Others)'].push(idx);
    }
  });

  return Object.entries(groups).filter(([_, indices]) => indices.length > 0).map(([name, indices]) => ({ name, indices }));
});

const toggleCategory = (name: string) => openCategories.value[name] = !openCategories.value[name];

// --- WORKSPACE MEMORY ---
const saveWorkspace = () => {
  if (parseResult.value) {
    const names = selectedIndices.value.map(i => parseResult.value!.headers[i]);
    localStorage.setItem('hwviz_workspace_names', JSON.stringify(names));
  }
};
watch(selectedIndices, saveWorkspace, { deep: true });

const loadWorkspace = (result: ParseResult) => {
  const saved = localStorage.getItem('hwviz_workspace_names');
  let restored = false;
  if (saved) {
    try {
      const parsedNames = JSON.parse(saved);
      if (Array.isArray(parsedNames) && parsedNames.length > 0) {
        const toSelect: number[] = [];
        result.headers.forEach((h, idx) => {
          if (parsedNames.includes(h)) toSelect.push(idx);
        });
        if (toSelect.length > 0) {
          selectedIndices.value = toSelect;
          // expand corresponding categories
          openCategories.value = {};
          categorizedMetrics.value.forEach(cat => {
            if (cat.indices.some(i => toSelect.includes(i))) {
              openCategories.value[cat.name] = true;
            }
          });
          restored = true;
        }
      }
    } catch(e) {}
  }
  
  if (!restored) {
    // Default: expand first two, select nothing
    selectedIndices.value = [];
    openCategories.value = {};
    const cat = categorizedMetrics.value;
    if (cat.length > 0) openCategories.value[cat[0].name] = true;
    if (cat.length > 1) openCategories.value[cat[1].name] = true;
  }
};

// --- FILE LOADING ---
const handleParseResult = (res: ParseResult) => {
  parseResult.value = res;
  fileSearch.value = "";
  localStorage.setItem('hwviz_last_file', res.loadedFile);
  setTimeout(() => loadWorkspace(res), 50);
};

const handleManualOpenFile = async () => {
  try {
    if (window.go && window.go.main && window.go.main.App) {
      isLoading.value = true;
      const res = await window.go.main.App.OpenFile();
      if(res) handleParseResult(res);
    }
  } catch (error) { errorMessage.value = "Failed to open: " + error; }
  finally { isLoading.value = false; }
};

const openSpecificFile = async (path: string) => {
  if (path === parseResult.value?.loadedFile) return;
  try {
    if (window.go && window.go.main && window.go.main.App) {
      isLoading.value = true;
      const res = await window.go.main.App.ParseCSV(path);
      if(res) handleParseResult(res);
    }
  } catch (error) { errorMessage.value = "Failed to load file: " + error; }
  finally { isLoading.value = false; }
};

const chartAreaContainerRef = ref<HTMLElement | null>(null);
const chartContainerWidth = ref(1200);

// --- ON MOUNT / INIT ---
onMounted(() => {
  const resizeObserver = new ResizeObserver(entries => {
    if (entries.length > 0) {
      chartContainerWidth.value = entries[0].contentRect.width;
    }
  });

  if (chartAreaContainerRef.value) {
    resizeObserver.observe(chartAreaContainerRef.value);
  }

  const lastFile = localStorage.getItem('hwviz_last_file');
  if (lastFile) {
    openSpecificFile(lastFile);
  }

  // Wails Native Event for Native Desktop Drops
  if (window.runtime && window.runtime.EventsOn) {
    window.runtime.EventsOn("file-dropped", (path: string) => {
      openSpecificFile(path);
    });
  }

  // HTML5 Drag fallback (for visual overlay inside the browser window)
  window.addEventListener('dragover', (e) => { e.preventDefault(); isDragging.value = true; });
  window.addEventListener('dragleave', (e) => {
    if (e.clientX <= 0 || e.clientY <= 0 || e.clientX >= window.innerWidth || e.clientY >= window.innerHeight) {
      isDragging.value = false;
    }
  });
  window.addEventListener('drop', (e) => {
    e.preventDefault();
    isDragging.value = false;
  });
});

// --- MEGA CHART SETUP ---
const megaChartHeight = computed(() => {
  const count = selectedIndices.value.length;
  if (count === 0) return 0;
  const colCount = Math.max(1, Math.floor(chartContainerWidth.value / 600));
  const rowCount = Math.ceil(count / colCount);
  return rowCount * 220 + 40; // 220px per row + 40px padding
});

const megaChartOption = computed(() => {
  if (!parseResult.value || selectedIndices.value.length === 0) return {};

  const indices = selectedIndices.value;
  const count = indices.length;
  
  const titles: any[] = [];
  const grids: any[] = [];
  const xAxes: any[] = [];
  const yAxes: any[] = [];
  const series: any[] = [];

  const colCount = Math.max(1, Math.floor(chartContainerWidth.value / 600));
  const rowCount = Math.ceil(count / colCount);

  const chartHeightPx = 220;
  const topOffset = 35;
  const gridHeight = 160;
  
  indices.forEach((dataIndex, i) => {
    const col = i % colCount;
    const row = Math.floor(i / colCount);

    const span = 100 / colCount;
    const leftPadding = 6;
    const rightPadding = 4;
    const left = col * span + leftPadding;
    const width = span - leftPadding - rightPadding;

    titles.push({
      text: parseResult.value!.headers[dataIndex],
      left: `${left}%`,
      top: row * chartHeightPx + 5,
      textStyle: { color: colors[i % colors.length], fontSize: 13, fontWeight: 500 }
    });

    grids.push({
      left: `${left}%`,
      width: `${width}%`,
      top: row * chartHeightPx + topOffset,
      height: gridHeight,
      containLabel: false
    });

    const isBottomMostInCol = row === rowCount - 1 || (row === rowCount - 2 && col >= count % colCount && count % colCount !== 0);

    xAxes.push({
      type: 'category',
      data: parseResult.value!.times,
      gridIndex: i,
      axisLabel: { show: isBottomMostInCol, color: '#a8a29e' },
      axisTick: { show: isBottomMostInCol },
      axisLine: { lineStyle: { color: '#444' } }
    });

    yAxes.push({
      type: 'value',
      gridIndex: i,
      splitLine: { lineStyle: { color: 'rgba(255,255,255,0.05)' } },
      axisLabel: { color: '#a8a29e', margin: 12 },
      axisPointer: { 
        show: true, 
        label: { show: true, precision: 2, backgroundColor: colors[i % colors.length] }
      },
      scale: true
    });

    series.push({
      name: parseResult.value!.headers[dataIndex],
      type: 'line',
      data: parseResult.value!.data[dataIndex],
      xAxisIndex: i,
      yAxisIndex: i,
      itemStyle: { color: colors[i % colors.length] },
      lineStyle: { width: 1.5, opacity: 0.9 },
      showSymbol: false,
      sampling: 'lttb',
      large: true,
      largeThreshold: 2000,
      hoverAnimation: false,
      emphasis: { 
        focus: 'series',
        scale: true,
        label: {
          show: true,
          position: 'top',
          color: '#ffffff',
          backgroundColor: colors[i % colors.length],
          padding: [4, 6],
          borderRadius: 4,
          formatter: (params: any) => Number(params.value).toFixed(2),
          fontWeight: 'bold',
          fontSize: 12
        }
      }
    });
  });

  const xAxisIndices = Array.from({length: count}, (_, i) => i);

  return {
    backgroundColor: 'transparent',
    title: titles,
    grid: grids,
    xAxis: xAxes,
    yAxis: yAxes,
    series: series,
    tooltip: { 
      trigger: 'axis',
      showContent: false,
      axisPointer: { 
        type: 'cross', 
        snap: true, 
        animation: false,
        crossStyle: { color: '#f59e0b', type: 'dashed' },
        label: { backgroundColor: 'rgba(30, 20, 15, 0.95)', color: '#e5e7eb' }
      }
    },
    axisPointer: {
      link: [{ xAxisIndex: 'all' }]
    },
    toolbox: {
      feature: {
        dataZoom: { yAxisIndex: 'none' }
      }
    },
    dataZoom: [
      {
        id: 'global-zoom',
        type: 'inside',
        xAxisIndex: xAxisIndices,
        filterMode: 'filter'
      }
    ]
  };
});

// --- GLOBAL TIMELINE (SCRUBBER) ---
const timelineOption = computed(() => {
  if (!parseResult.value) return {};
  return {
    backgroundColor: 'rgba(30, 20, 15, 0.95)',
    grid: { left: '6%', right: '5%', top: 5, bottom: 5, height: 0 },
    xAxis: { type: 'category', data: parseResult.value.times, show: false },
    yAxis: { type: 'value', show: false },
    dataZoom: [
      {
        id: 'global-zoom',
        type: 'slider',
        show: true,
        xAxisIndex: [0],
        filterMode: 'filter',
        bottom: 0,
        height: 38,
        borderColor: 'transparent',
        backgroundColor: 'rgba(0,0,0,0.3)',
        fillerColor: 'rgba(245, 158, 11, 0.25)',
        textStyle: { color: '#a8a29e' },
        handleStyle: { color: '#f59e0b', shadowBlur: 3, shadowColor: 'rgba(0,0,0,0.6)' },
        moveHandleStyle: { color: '#ea580c' }
      }
    ],
    series: [{ type: 'line', data: parseResult.value.times.map(()=>0), showSymbol: false, lineStyle: {opacity: 0} }]
  };
});
</script>

<style scoped>
.app-layout {
  display: flex;
  width: 100vw;
  height: 100vh;
  position: relative;
}

/* Drop Overlay */
.drop-overlay {
  position: absolute;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(15, 10, 8, 0.85);
  backdrop-filter: blur(10px);
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
}
.drop-box {
  width: 500px;
  height: 300px;
  border: 3px dashed var(--accent);
  border-radius: 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: var(--accent);
  box-shadow: 0 0 50px rgba(245, 158, 11, 0.2);
  animation: pulse-drop 2s infinite alternate;
}
.drop-icon { margin-bottom: 20px; }
@keyframes pulse-drop {
  0% { transform: scale(0.98); box-shadow: 0 0 20px rgba(245, 158, 11, 0.1); }
  100% { transform: scale(1.02); box-shadow: 0 0 60px rgba(245, 158, 11, 0.4); }
}

.sidebar {
  width: 270px;
  background: var(--bg-card);
  backdrop-filter: var(--glass-blur);
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  height: 100%;
  box-shadow: 4px 0 24px rgba(0,0,0,0.2);
  z-index: 10;
}

.sidebar-header {
  padding: 24px 24px 16px;
  /* background: linear-gradient(180deg, rgba(245, 158, 11, 0.05), transparent); */
}

.sidebar-header h2 {
  font-size: 1.5rem;
  font-weight: 600;
  background: linear-gradient(135deg, #fcd34d, #f97316, #ef4444);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  letter-spacing: -0.02em;
}

/* Dock */
.dock {
  width: 68px;
  background: rgba(15, 10, 8, 0.95);
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 24px 0;
  z-index: 15;
}
.dock-top, .dock-bottom {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}
.dock-btn {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.05);
  color: var(--text-muted);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
}
.dock-btn:hover {
  background: rgba(255, 255, 255, 0.08);
  color: #fff;
  transform: translateY(-2px);
}
.dock-btn.active {
  background: rgba(245, 158, 11, 0.15);
  border-color: rgba(245, 158, 11, 0.4);
  color: var(--accent);
  box-shadow: 0 4px 15px rgba(245, 158, 11, 0.2);
}
.dock-btn.primary {
  background: linear-gradient(135deg, var(--accent), #ea580c);
  color: white;
  border: none;
  box-shadow: 0 4px 15px var(--accent-glow);
}
.dock-btn.primary:hover {
  box-shadow: 0 8px 25px var(--accent-glow);
  transform: translateY(-2px);
}

/* Loading Overlay */
.loading-overlay {
  position: absolute;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(15, 10, 8, 0.85);
  backdrop-filter: blur(10px);
  z-index: 99999;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
}
.spinner-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  position: relative;
  text-align: center;
}
.spinner {
  width: 60px;
  height: 60px;
  border: 4px solid rgba(245, 158, 11, 0.1);
  border-top: 4px solid var(--accent);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 24px;
  z-index: 2;
}
.loading-orb {
  width: 120px;
  height: 120px;
  z-index: 1;
}
@keyframes spin { 0% { transform: rotate(0deg); } 100% { transform: rotate(360deg); } }
.spinner-container h3 {
  font-size: 1.4rem;
  font-weight: 600;
  color: white;
  margin-bottom: 8px;
  z-index: 2;
}
.spinner-container p {
  color: var(--text-muted);
  z-index: 2;
}

/* File Explorer */
.file-explorer {
  flex: 1;
  border-bottom: 1px solid var(--border-color);
  background: rgba(0,0,0,0.15);
  display: flex;
  flex-direction: column;
}
.explorer-header h3 {
  padding: 16px 24px 10px;
  font-size: 0.8rem;
  font-weight: 600;
  text-transform: uppercase;
  color: var(--text-muted);
}
.explorer-filter {
  padding: 0 24px 12px;
}
.explorer-filter input {
  width: 100%;
  padding: 8px 12px;
  background: rgba(255,255,255,0.05);
  border: 1px solid rgba(255,255,255,0.1);
  border-radius: 6px;
  color: white;
  font-family: inherit;
  font-size: 0.85rem;
  outline: none;
  transition: border-color 0.2s;
}
.explorer-filter input:focus { border-color: var(--accent); }
.file-list {
  overflow-y: auto;
  padding-bottom: 10px;
}
.file-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 24px;
  font-size: 0.85rem;
  color: #d6d3d1;
  cursor: pointer;
  transition: background 0.2s;
}
.file-item:hover { background: rgba(255,255,255,0.05); }
.file-item.active { background: rgba(245, 158, 11, 0.15); color: #fff; border-left: 3px solid var(--accent); }

/* Metrics */
.metrics-list {
  flex: 1;
  overflow-y: auto;
  padding: 20px 0;
}
.metrics-title-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px 16px;
}
.metrics-title-row h3 {
  font-size: 0.8rem;
  font-weight: 600;
  text-transform: uppercase;
  color: var(--text-muted);
}
.badge {
  background: rgba(245, 158, 11, 0.15);
  color: var(--accent);
  padding: 3px 8px;
  border-radius: 12px;
  font-size: 0.75rem;
}

/* Accordion */
.category { margin-bottom: 4px; }
.category-header { display: flex; align-items: center; justify-content: space-between; padding: 12px 24px; cursor: pointer; background: rgba(255, 255, 255, 0.01); transition: 0.2s; }
.category-header:hover { background: rgba(255, 255, 255, 0.05); }
.category-header.is-open { background: rgba(245, 158, 11, 0.05); border-bottom: 1px solid rgba(245, 158, 11, 0.1); }
.cat-left { display: flex; align-items: center; gap: 10px; color: #fbbf24; font-weight: 500; font-size: 0.95rem; }
.chevron { transition: transform 0.3s ease; color: var(--text-muted); }
.rotate-90 { transform: rotate(90deg); color: #fbbf24; }
.category-badge { font-size: 0.75rem; color: var(--text-muted); background: rgba(0, 0, 0, 0.3); padding: 2px 6px; border-radius: 6px; }

.category-content { overflow: hidden; background: rgba(0, 0, 0, 0.15); }
.slide-fade-enter-active { transition: all 0.3s ease-out; max-height: 2000px; }
.slide-fade-leave-active { transition: all 0.2s ease-in; }
.slide-fade-enter-from, .slide-fade-leave-to { opacity: 0; transform: translateY(-10px); }

/* Checkboxes */
.metric-item { display: flex; align-items: center; padding: 10px 24px 10px 42px; cursor: pointer; transition: 0.2s; border-left: 3px solid transparent; }
.metric-item:hover { background: rgba(255, 255, 255, 0.03); }
.metric-item.selected { background: rgba(245, 158, 11, 0.08); border-left: 3px solid var(--accent); }

.custom-checkbox { position: relative; width: 16px; height: 16px; margin-right: 12px; flex-shrink: 0; }
.custom-checkbox input { opacity: 0; width: 0; height: 0; position: absolute; }
.checkmark { position: absolute; top: 0; left: 0; width: 100%; height: 100%; background: rgba(0,0,0,0.4); border: 1px solid rgba(255,255,255,0.15); border-radius: 4px; transition: 0.2s; }
.metric-item:hover .checkmark { border-color: rgba(255,255,255,0.3); }
.custom-checkbox input:checked ~ .checkmark { background-color: var(--accent); border-color: var(--accent); box-shadow: 0 0 10px var(--accent-glow); }
.checkmark:after { content: ""; position: absolute; display: none; left: 5px; top: 1px; width: 3px; height: 7px; border: solid white; border-width: 0 2px 2px 0; transform: rotate(45deg); }
.custom-checkbox input:checked ~ .checkmark:after { display: block; animation: check-pop 0.2s; }
@keyframes check-pop { 0% { transform: rotate(45deg) scale(0); } 50% { transform: rotate(45deg) scale(1.2); } 100% { transform: rotate(45deg) scale(1); } }

.metric-name { font-size: 0.9rem; color: #e2e8f0; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; user-select: none; transition: 0.2s; }
.metric-item.selected .metric-name { color: #fff; font-weight: 500; }

.empty-state { flex: 1; display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 32px; text-align: center; color: var(--text-muted); }
.empty-icon-wrap { width: 64px; height: 64px; border-radius: 50%; background: rgba(255,255,255,0.03); display: flex; align-items: center; justify-content: center; margin-bottom: 20px; }

/* Main Area */
.chart-area-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  position: relative;
  overflow: hidden;
}

.chart-area {
  flex: 1;
  overflow-y: auto;
  padding: 32px;
  scroll-behavior: smooth;
  position: relative;
}

.mega-chart-container {
  width: 100%;
  min-height: 100%;
}
.mega-chart {
  width: 100%;
}

.no-chart-state { position: absolute; top: 0; left: 0; right: 0; bottom: 0; display: flex; align-items: center; justify-content: center; }
.placeholder { display: flex; flex-direction: column; align-items: center; text-align: center; color: var(--text-muted); opacity: 0.8; position: relative; }
.glow-orb { position: absolute; width: 150px; height: 150px; background: var(--accent); filter: blur(100px); opacity: 0.15; top: 50%; left: 50%; transform: translate(-50%, -50%); animation: pulse 4s infinite alternate; }
@keyframes pulse { 0% { opacity: 0.1; transform: translate(-50%, -50%) scale(0.8); } 100% { opacity: 0.2; transform: translate(-50%, -50%) scale(1.2); } }
.placeholder-icon { margin-bottom: 20px; filter: drop-shadow(0 0 10px rgba(0,0,0,0.5)); z-index: 1; }
.placeholder p { font-size: 1.1rem; z-index: 1; }

/* Global Timeline */
.global-timeline {
  height: 50px;
  width: 100%;
  border-top: 1px solid var(--border-color);
  background: rgba(15, 10, 8, 0.95);
  backdrop-filter: blur(20px);
  box-shadow: 0 -4px 20px rgba(0,0,0,0.3);
  z-index: 20;
}
.timeline-chart {
  width: 100%;
  height: 100%;
}

/* Error Modal */
.error-modal-overlay {
  position: absolute;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(15, 10, 8, 0.7);
  backdrop-filter: blur(5px);
  z-index: 100000;
  display: flex;
  align-items: center;
  justify-content: center;
}
.error-modal-box {
  background: rgba(30, 20, 20, 0.95);
  border: 1px solid rgba(239, 68, 68, 0.3);
  padding: 24px;
  border-radius: 12px;
  max-width: 400px;
  text-align: center;
  box-shadow: 0 10px 40px rgba(239, 68, 68, 0.2);
}
.error-modal-header {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  color: #ef4444;
  margin-bottom: 16px;
}
.error-modal-body {
  color: #d1d5db;
  font-size: 0.95rem;
  margin-bottom: 24px;
  word-break: break-all;
}
.error-modal-btn {
  background: #ef4444;
  color: white;
  border: none;
  padding: 8px 24px;
  border-radius: 6px;
  cursor: pointer;
  transition: 0.2s;
}
.error-modal-btn:hover {
  background: #dc2626;
}
</style>
