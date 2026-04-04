<template>
  <div class="app-layout">
    <!-- Drop Overlay -->
    <div v-show="isDragging" class="drop-overlay">
      <div class="drop-box">
        <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" class="drop-icon"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg>
        <h2>Release to open CSV</h2>
      </div>
    </div>

    <aside class="sidebar">
      <div class="sidebar-header" @click="handleManualOpenFile">
        <h2>Hardware Viz</h2>
        <button class="open-btn">
          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg>
          Open File
        </button>
      </div>

      <!-- File Explorer inside Sidebar -->
      <div class="file-explorer" v-if="parseResult && parseResult.files.length > 0">
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

      <div class="metrics-list" v-if="parseResult && categorizedMetrics.length > 0">
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

    <main class="chart-area-container">
      <div class="chart-area" v-if="parseResult">
        <div v-if="selectedIndices.length === 0" class="no-chart-state">
          <div class="placeholder">
            <div class="glow-orb"></div>
            <svg xmlns="http://www.w3.org/2000/svg" width="56" height="56" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round" class="placeholder-icon"><path d="M3 3v18h18"/><path d="m19 9-5 5-4-4-3 3"/></svg>
            <p>Select metrics from the left to view charts</p>
          </div>
        </div>
        
        <transition-group name="chart-list" tag="div" class="charts-wrapper" mode="out-in">
          <Chart
            v-for="(idx, i) in selectedIndices"
            :key="idx"
            :title="parseResult.headers[idx]"
            :times="parseResult.times"
            :data="parseResult.data[idx]"
            :colorLine="colors[i % colors.length]"
          />
        </transition-group>
      </div>
      
      <!-- Global Time Scrubber -->
      <div class="global-timeline" v-if="parseResult && selectedIndices.length > 0">
        <v-chart class="timeline-chart" :option="timelineOption" autoresize />
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue';
import { use } from 'echarts/core';
import { CanvasRenderer } from 'echarts/renderers';
import { LineChart, CustomChart } from 'echarts/charts';
import { GridComponent, DataZoomComponent } from 'echarts/components';
import VChart from 'vue-echarts';
import Chart from './components/Chart.vue';
import './style.css';

use([CanvasRenderer, LineChart, CustomChart, GridComponent, DataZoomComponent]);

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
    'CPU 處理器': [], 'GPU 顯示卡': [], 'Memory 記憶體': [],
    'Storage 磁碟': [], 'Network 網路': [], 'Power & Fan 電源與散熱': [], 'Others 其他': []
  };

  parseResult.value.headers.forEach((h, idx) => {
    const lower = h.toLowerCase();
    if (/cpu|核心|處理器|thread|clock|時脈/i.test(lower) && !/gpu/i.test(lower)) groups['CPU 處理器'].push(idx);
    else if (/gpu|顯示卡|video|vram|d3d/i.test(lower)) groups['GPU 顯示卡'].push(idx);
    else if (/記憶體|memory|ram|分頁檔|virtual memory/i.test(lower) && !/gpu/i.test(lower)) groups['Memory 記憶體'].push(idx);
    else if (/硬碟|讀取|寫入|s\.m\.a\.r\.t|drive|ssd|hdd|磁碟|read|write/i.test(lower)) groups['Storage 磁碟'].push(idx);
    else if (/網路|network|頻寬|lan|ethernet/i.test(lower)) groups['Network 網路'].push(idx);
    else if (/主機板|vcore|風扇|系統|fan|system|voltage|power|電壓|功耗|溫度|temp|w/i.test(lower)) groups['Power & Fan 電源與散熱'].push(idx);
    else groups['Others 其他'].push(idx);
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
  setTimeout(() => loadWorkspace(res), 50);
};

const handleManualOpenFile = async () => {
  try {
    if (window.go && window.go.main && window.go.main.App) {
      const res = await window.go.main.App.OpenFile();
      handleParseResult(res);
    }
  } catch (error) { alert("Failed to open: " + error); }
};

const openSpecificFile = async (path: string) => {
  if (path === parseResult.value?.loadedFile) return;
  try {
    if (window.go && window.go.main && window.go.main.App) {
      const res = await window.go.main.App.ParseCSV(path);
      handleParseResult(res);
    }
  } catch (error) { alert("Failed to load file: " + error); }
};

// --- DRAG AND DROP ---
onMounted(() => {
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

// --- GLOBAL TIMELINE (SCRUBBER) ---
// Dummy chart option just to render a global dataZoom slider mapped to hw-group
const timelineOption = computed(() => {
  if (!parseResult.value) return {};
  return {
    group: 'hw-group',
    backgroundColor: 'rgba(30, 20, 15, 0.95)',
    grid: { left: '8%', right: '5%', top: 5, bottom: 5, height: 0 },
    xAxis: { type: 'category', data: parseResult.value.times, show: false },
    yAxis: { type: 'value', show: false },
    dataZoom: [
      {
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
  width: 330px;
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
  padding: 32px 24px 24px;
  /* background: linear-gradient(180deg, rgba(245, 158, 11, 0.05), transparent); */
}

.sidebar-header h2 {
  font-size: 1.8rem;
  font-weight: 600;
  margin-bottom: 24px;
  background: linear-gradient(135deg, #fcd34d, #f97316, #ef4444);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  letter-spacing: -0.02em;
}

.open-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  width: 100%;
  padding: 14px;
  background: linear-gradient(135deg, var(--accent), #ea580c);
  color: white;
  border: 1px solid rgba(255,255,255,0.1);
  border-radius: 10px;
  font-size: 1.05rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s;
  box-shadow: 0 4px 15px var(--accent-glow);
}
.open-btn:hover { transform: translateY(-2px); box-shadow: 0 8px 25px var(--accent-glow); }

/* File Explorer */
.file-explorer {
  border-bottom: 1px solid var(--border-color);
  background: rgba(0,0,0,0.15);
  display: flex;
  flex-direction: column;
  max-height: 250px;
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

.charts-wrapper { display: flex; flex-direction: column; gap: 24px; padding-bottom: 40px; }
.chart-list-enter-active, .chart-list-leave-active { transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1); }
.chart-list-enter-from, .chart-list-leave-to { opacity: 0; transform: translateY(20px) scale(0.95); }

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
</style>
