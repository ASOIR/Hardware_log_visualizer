<template>
  <div class="chart-container">
    <v-chart class="chart" :option="chartOption" autoresize @ready="onChartReady" />
  </div>
</template>

<script setup lang="ts">
import { computed, provide } from 'vue';
import { use } from 'echarts/core';
import { CanvasRenderer } from 'echarts/renderers';
import { LineChart } from 'echarts/charts';
import { TitleComponent, TooltipComponent, GridComponent, DataZoomComponent, ToolboxComponent } from 'echarts/components';
import VChart, { THEME_KEY } from 'vue-echarts';
import * as echarts from 'echarts/core';

use([
  CanvasRenderer,
  LineChart,
  TitleComponent,
  TooltipComponent,
  GridComponent,
  DataZoomComponent,
  ToolboxComponent
]);

provide(THEME_KEY, 'dark');

const props = defineProps<{
  title: string;
  times: string[];
  data: number[];
  colorLine?: string;
}>();

echarts.connect('hw-group');

const onChartReady = (instance: any) => {
  // Activate box-zoom tool by default so left-drag creates a zoom box on load
  instance.dispatchAction({
    type: 'takeGlobalCursor',
    key: 'dataZoomSelect',
    dataZoomSelectActive: true
  });
};

const chartOption = computed(() => {
  const c = props.colorLine || '#f97316'; // Fallback to sunset orange
  return {
    group: 'hw-group',
    backgroundColor: 'transparent',
    title: {
      text: props.title,
      left: 'center',
      textStyle: { color: '#fffbeb', fontFamily: 'Outfit, sans-serif', fontSize: 15, fontWeight: 500 }
    },
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'rgba(44, 34, 30, 0.95)',
      borderColor: c,
      borderWidth: 1.5,
      padding: [8, 12],
      textStyle: { color: '#fff' }
    },
    toolbox: {
      show: true,
      feature: {
        dataZoom: { yAxisIndex: 'none', title: { zoom: '框選放大', back: '還原縮放' } }
      },
      iconStyle: { borderColor: '#a8a29e' },
      emphasis: { iconStyle: { borderColor: '#f59e0b' } }
    },
    grid: { left: '8%', right: '5%', bottom: '15%', top: '22%', containLabel: true },
    xAxis: {
      type: 'category',
      data: props.times,
      boundaryGap: false,
      axisLabel: { color: '#a8a29e', fontFamily: 'Inter' },
      axisLine: { lineStyle: { color: 'rgba(245, 158, 11, 0.2)' } }
    },
    yAxis: {
      type: 'value',
      axisLabel: { color: '#a8a29e', fontFamily: 'Inter' },
      splitLine: { lineStyle: { color: 'rgba(245, 158, 11, 0.1)', type: 'dashed' } },
      scale: true // Auto scale y-axis nicely based on current zoom data
    },
    dataZoom: [
      { 
        type: 'inside', 
        filterMode: 'filter' 
      }
    ],
    series: [
      {
        name: props.title,
        type: 'line',
        data: props.data,
        smooth: true, // curved lines
        showSymbol: false,
        lineStyle: { 
          color: c, 
          width: 2.5,
          shadowColor: 'rgba(0, 0, 0, 0.5)',
          shadowBlur: 10,
          shadowOffsetY: 5
        },
        itemStyle: { color: c },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: c.replace(')', ', 0.3)').replace('rgb', 'rgba') },
            { offset: 0.8, color: c.replace(')', ', 0.0)').replace('rgb', 'rgba') }
          ]) // Simple fallback if c is hex
        }
      }
    ]
  };
});
</script>

<style scoped>
.chart-container {
  width: 100%;
  height: 350px;
  background: var(--bg-card);
  backdrop-filter: var(--glass-blur);
  border: 1px solid var(--border-color);
  border-radius: 14px;
  padding: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.4), inset 0 1px 0 rgba(255, 255, 255, 0.05);
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1), box-shadow 0.3s ease, border-color 0.3s ease;
  position: relative;
  overflow: hidden;
}

.chart-container::before {
  content: '';
  position: absolute;
  top: 0; left: 0; right: 0;
  height: 100%;
  background: radial-gradient(circle at 50% 0%, rgba(245, 158, 11, 0.05), transparent 70%);
  pointer-events: none;
}

.chart-container:hover {
  transform: translateY(-4px);
  border-color: rgba(245, 158, 11, 0.4);
  box-shadow: 0 12px 30px rgba(0, 0, 0, 0.5), 0 0 20px rgba(245, 158, 11, 0.15);
  z-index: 2;
}

.chart {
  width: 100%;
  height: 100%;
  position: relative;
  z-index: 1;
}
</style>
