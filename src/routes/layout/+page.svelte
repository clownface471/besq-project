<script lang="ts">
  import { onMount } from "svelte";
  import { fade, scale } from "svelte/transition";

  type ChartShape = "rectangle" | "oval" | "circle" | "triangle" | "rhombus";
  type ArrowStyle = "solid" | "dashed" | "dotted";
  type ArrowHead = "triangle" | "diamond" | "circle" | "none";

  interface Chart {
    id: string;
    x: number;
    y: number;
    width: number;
    height: number;
    shape: ChartShape;
    title: string;
    content: string;
    rotation: number;
    locked: boolean;
  }

  interface Arrow {
    id: string;
    fromId: string;
    toId: string;
    style: ArrowStyle;
    head: ArrowHead;
    curvature: number;
    label: string;
    color: string;
  }

  const STORAGE_KEY = "chart-editor-data";

  let charts: Chart[] = [
    {
      id: "1",
      x: 100,
      y: 150,
      width: 180,
      height: 100,
      shape: "rectangle",
      title: "Start",
      content: "Begin here",
      rotation: 0,
      locked: false,
    },
    {
      id: "2",
      x: 400,
      y: 150,
      width: 180,
      height: 100,
      shape: "oval",
      title: "Process",
      content: "Processing",
      rotation: 0,
      locked: false,
    },
    {
      id: "3",
      x: 700,
      y: 150,
      width: 140,
      height: 140,
      shape: "circle",
      title: "End",
      content: "Finish",
      rotation: 0,
      locked: false,
    },
  ];

  let arrows: Arrow[] = [
    {
      id: "a1",
      fromId: "1",
      toId: "2",
      style: "solid",
      head: "triangle",
      curvature: 0,
      label: "",
      color: "#6366f1",
    },
    {
      id: "a2",
      fromId: "2",
      toId: "3",
      style: "dashed",
      head: "diamond",
      curvature: 20,
      label: "Go to",
      color: "#8b5cf6",
    },
  ];

  let selectedId: string | null = null;
  let isDragging = false;
  let dragOffset = { x: 0, y: 0 };
  let isConnecting = false;
  let connectStartId: string | null = null;
  let mousePos = { x: 0, y: 0 };
  let popupPos = { x: 0, y: 0 };
  let isResizing = false;
  let resizeHandle = "";
  let canvasEl: HTMLDivElement | null = null;
  let showArrowPopup = false;
  let showShapeMenu = false;
  let saveStatus: "saved" | "saving" | "unsaved" = "saved";
  let lastSavedTime: Date | null = null;

  const GRID_SIZE = 20; // For choppy movement

  // Storage functions
  function saveToStorage(): void {
    saveStatus = "saving";
    try {
      const data = { charts, arrows, timestamp: new Date().toISOString() };
      localStorage.setItem(STORAGE_KEY, JSON.stringify(data));
      lastSavedTime = new Date();
      saveStatus = "saved";
    } catch (e) {
      console.error("Failed to save to localStorage:", e);
      saveStatus = "unsaved";
    }
  }

  function loadFromStorage(): { charts: Chart[]; arrows: Arrow[] } | null {
    try {
      const data = localStorage.getItem(STORAGE_KEY);
      if (data) {
        const parsed = JSON.parse(data);
        // Validate data structure
        if (
          parsed.charts &&
          Array.isArray(parsed.charts) &&
          parsed.arrows &&
          Array.isArray(parsed.arrows)
        ) {
          return { charts: parsed.charts, arrows: parsed.arrows };
        }
      }
    } catch (e) {
      console.error("Failed to load from localStorage:", e);
    }
    return null;
  }

  function resetToDefaults(): void {
    charts = [
      {
        id: "1",
        x: 100,
        y: 150,
        width: 180,
        height: 100,
        shape: "rectangle",
        title: "Start",
        content: "Begin here",
        rotation: 0,
        locked: false,
      },
      {
        id: "2",
        x: 400,
        y: 150,
        width: 180,
        height: 100,
        shape: "oval",
        title: "Process",
        content: "Processing",
        rotation: 0,
        locked: false,
      },
      {
        id: "3",
        x: 700,
        y: 150,
        width: 140,
        height: 140,
        shape: "circle",
        title: "End",
        content: "Finish",
        rotation: 0,
        locked: false,
      },
    ];
    arrows = [
      {
        id: "a1",
        fromId: "1",
        toId: "2",
        style: "solid",
        head: "triangle",
        curvature: 0,
        label: "",
        color: "#6366f1",
      },
      {
        id: "a2",
        fromId: "2",
        toId: "3",
        style: "dashed",
        head: "diamond",
        curvature: 20,
        label: "Go to",
        color: "#8b5cf6",
      },
    ];
    saveToStorage();
  }

  function triggerSave(): void {
    saveToStorage();
  }

  function generateId(): string {
    return Math.random().toString(36).substr(2, 9);
  }

  function snapToGrid(value: number): number {
    return Math.round(value / GRID_SIZE) * GRID_SIZE;
  }

  function addChart() {
    const id = generateId();
    const newChart: Chart = {
      id,
      x: snapToGrid(100 + Math.random() * 300),
      y: snapToGrid(100 + Math.random() * 200),
      width: 180,
      height: 100,
      shape: "rectangle",
      title: "New Chart",
      content: "Click to edit",
      rotation: 0,
      locked: false,
    };
    charts = [...charts, newChart];
    selectedId = id;
    saveToStorage();
  }

  function addChartShape(shape: ChartShape) {
    const id = generateId();
    let width = 180,
      height = 100;
    if (shape === "circle") {
      width = 120;
      height = 120;
    }
    if (shape === "triangle") {
      width = 140;
      height = 120;
    }
    if (shape === "rhombus") {
      width = 160;
      height = 100;
    }

    const newChart: Chart = {
      id,
      x: snapToGrid(100 + Math.random() * 200),
      y: snapToGrid(100 + Math.random() * 200),
      width,
      height,
      shape,
      title: `${shape.charAt(0).toUpperCase() + shape.slice(1)}`,
      content: "Content here",
      rotation: 0,
      locked: false,
    };
    charts = [...charts, newChart];
    selectedId = id;
    showShapeMenu = false;
    saveToStorage();
  }

  function removeChart(id: string) {
    charts = charts.filter((c) => c.id !== id);
    arrows = arrows.filter((a) => a.fromId !== id && a.toId !== id);
    if (selectedId === id) selectedId = null;
    saveToStorage();
  }

  function addArrow(fromId: string, toId: string) {
    if (fromId === toId) return;
    const exists = arrows.find((a) => a.fromId === fromId && a.toId === toId);
    if (!exists) {
      const newArrow: Arrow = {
        id: generateId(),
        fromId,
        toId,
        style: "solid",
        head: "triangle",
        curvature: 0,
        label: "",
        color: "#6366f1",
      };
      arrows = [...arrows, newArrow];
      saveToStorage();
    }
  }

  function removeArrow(id: string) {
    arrows = arrows.filter((a) => a.id !== id);
    showArrowPopup = false;
    saveToStorage();
  }

  function startDrag(e: MouseEvent, chartId: string) {
    if (isConnecting || isResizing) return;
    e.stopPropagation();
    const chart = charts.find((c) => c.id === chartId);
    if (chart) {
      selectedId = chartId;
      isDragging = true;
      dragOffset = {
        x: e.clientX - chart.x,
        y: e.clientY - chart.y,
      };
    }
  }

  function handleMouseMove(e: MouseEvent) {
    mousePos = { x: e.clientX, y: e.clientY };

    if (isDragging && selectedId && !isConnecting) {
      const chart = charts.find((c) => c.id === selectedId);
      if (chart) {
        chart.x = snapToGrid(e.clientX - dragOffset.x);
        chart.y = snapToGrid(e.clientY - dragOffset.y);
        charts = charts;
        saveStatus = "unsaved";
      }
    }

    if (isResizing && selectedId) {
      const chart = charts.find((c) => c.id === selectedId);
      if (chart && canvasEl) {
        const rect = canvasEl.getBoundingClientRect();
        if (resizeHandle === "se") {
          chart.width = Math.max(
            80,
            Math.round((e.clientX - chart.x) / 10) * 10,
          );
          chart.height = Math.max(
            60,
            Math.round((e.clientY - chart.y) / 10) * 10,
          );
        } else if (resizeHandle === "sw") {
          const newWidth = Math.max(
            80,
            Math.round((chart.x + chart.width - e.clientX) / 10) * 10,
          );
          if (newWidth !== chart.width) {
            chart.x = chart.x + chart.width - newWidth;
            chart.width = newWidth;
          }
          chart.height = Math.max(
            60,
            Math.round((e.clientY - chart.y) / 10) * 10,
          );
        } else if (resizeHandle === "ne") {
          chart.width = Math.max(
            80,
            Math.round((e.clientX - chart.x) / 10) * 10,
          );
          const newHeight = Math.max(
            60,
            Math.round((chart.y + chart.height - e.clientY) / 10) * 10,
          );
          if (newHeight !== chart.height) {
            chart.y = chart.y + chart.height - newHeight;
            chart.height = newHeight;
          }
        } else if (resizeHandle === "nw") {
          const newWidth = Math.max(
            80,
            Math.round((chart.x + chart.width - e.clientX) / 10) * 10,
          );
          if (newWidth !== chart.width) {
            chart.x = chart.x + chart.width - newWidth;
            chart.width = newWidth;
          }
          const newHeight = Math.max(
            60,
            Math.round((chart.y + chart.height - e.clientY) / 10) * 10,
          );
          if (newHeight !== chart.height) {
            chart.y = chart.y + chart.height - newHeight;
            chart.height = newHeight;
          }
        }
        charts = charts;
        saveStatus = "unsaved";
      }
    }
  }

  function handleMouseUp() {
    if (isDragging || isResizing) {
      saveToStorage();
    }
    isDragging = false;
    isResizing = false;
    resizeHandle = "";
    if (isConnecting) {
      isConnecting = false;
      connectStartId = null;
    }
  }

  function startConnect(chartId: string) {
    isConnecting = true;
    connectStartId = chartId;
  }

  function endConnect(chartId: string) {
    if (isConnecting && connectStartId) {
      addArrow(connectStartId, chartId);
      isConnecting = false;
      connectStartId = null;
    }
  }

  function handleCanvasClick(e: MouseEvent) {
    if (e.target === canvasEl) {
      selectedId = null;
      showArrowPopup = false;
    }
  }

  function getChartCenter(chartId: string): { x: number; y: number } | null {
    const chart = charts.find((c) => c.id === chartId);
    if (chart) {
      return {
        x: chart.x + chart.width / 2,
        y: chart.y + chart.height / 2,
      };
    }
    return null;
  }

  function getChartEdgePoint(
    chartId: string,
    targetX: number,
    targetY: number,
  ): { x: number; y: number } | null {
    const chart = charts.find((c) => c.id === chartId);
    if (!chart) return null;

    const centerX = chart.x + chart.width / 2;
    const centerY = chart.y + chart.height / 2;
    const dx = targetX - centerX;
    const dy = targetY - centerY;
    const distance = Math.sqrt(dx * dx + dy * dy);

    if (distance === 0) return { x: centerX, y: centerY };

    const nx = dx / distance;
    const ny = dy / distance;

    // Find intersection with chart rectangle edges
    let intersectX = centerX;
    let intersectY = centerY;
    let tMax = Infinity;

    // Right edge
    if (nx !== 0) {
      const t = (chart.x + chart.width - centerX) / nx;
      if (t > 0 && t < tMax) {
        const y = centerY + t * ny;
        if (y >= chart.y && y <= chart.y + chart.height) {
          tMax = t;
          intersectX = chart.x + chart.width;
          intersectY = y;
        }
      }
    }

    // Left edge
    if (nx !== 0) {
      const t = (chart.x - centerX) / nx;
      if (t > 0 && t < tMax) {
        const y = centerY + t * ny;
        if (y >= chart.y && y <= chart.y + chart.height) {
          tMax = t;
          intersectX = chart.x;
          intersectY = y;
        }
      }
    }

    // Bottom edge
    if (ny !== 0) {
      const t = (chart.y + chart.height - centerY) / ny;
      if (t > 0 && t < tMax) {
        const x = centerX + t * nx;
        if (x >= chart.x && x <= chart.x + chart.width) {
          tMax = t;
          intersectX = x;
          intersectY = chart.y + chart.height;
        }
      }
    }

    // Top edge
    if (ny !== 0) {
      const t = (chart.y - centerY) / ny;
      if (t > 0 && t < tMax) {
        const x = centerX + t * nx;
        if (x >= chart.x && x <= chart.x + chart.width) {
          tMax = t;
          intersectX = x;
          intersectY = chart.y;
        }
      }
    }

    return { x: intersectX, y: intersectY };
  }

  function getArrowPath(arrow: Arrow): string {
    const fromCenter = getChartCenter(arrow.fromId);
    const toCenter = getChartCenter(arrow.toId);
    if (!fromCenter || !toCenter) return "";

    const from = getChartEdgePoint(arrow.fromId, toCenter.x, toCenter.y);
    const to = getChartEdgePoint(arrow.toId, fromCenter.x, fromCenter.y);
    if (!from || !to) return "";

    const dx = to.x - from.x;
    const dy = to.y - from.y;
    const distance = Math.sqrt(dx * dx + dy * dy);

    if (distance === 0) return "";

    // Normalize and create offset for arrow heads
    const offset = 10;
    const nx = dx / distance;
    const ny = dy / distance;

    const startX = from.x + nx * offset;
    const startY = from.y + ny * offset;
    const endX = to.x - nx * offset;
    const endY = to.y - ny * offset;

    if (arrow.curvature === 0) {
      return `M ${startX} ${startY} L ${endX} ${endY}`;
    }

    // Curved path
    const midX = (startX + endX) / 2;
    const midY = (startY + endY) / 2;
    const perpX = -ny;
    const perpY = nx;
    const cpX = midX + perpX * arrow.curvature;
    const cpY = midY + perpY * arrow.curvature;

    return `M ${startX} ${startY} Q ${cpX} ${cpY} ${endX} ${endY}`;
  }

  function isSelected(id: string): boolean {
    return selectedId === id;
  }

  function handleKeyDown(e: KeyboardEvent) {
    if (e.key === "Delete" || e.key === "Backspace") {
      if (selectedId && !showArrowPopup) {
        const isChart = charts.find((c) => c.id === selectedId);
        if (isChart) {
          removeChart(selectedId);
        } else {
          const isArrow = arrows.find((a) => a.id === selectedId);
          if (isArrow) removeArrow(selectedId);
        }
      }
    }
    if (e.key === "Escape") {
      selectedId = null;
      isConnecting = false;
      connectStartId = null;
      showArrowPopup = false;
      showShapeMenu = false;
    }
  }

  onMount(() => {
    // Load saved data from localStorage
    const savedData = loadFromStorage();
    if (savedData) {
      charts = savedData.charts;
      arrows = savedData.arrows;
      lastSavedTime = new Date();
      saveStatus = "saved";
    }

    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  });

  function updateArrow(arrowId: string, updates: Partial<Arrow>) {
    arrows = arrows.map((a) => (a.id === arrowId ? { ...a, ...updates } : a));
    saveToStorage();
  }

  function updateChart(chartId: string, updates: Partial<Chart>) {
    charts = charts.map((c) => (c.id === chartId ? { ...c, ...updates } : c));
    saveToStorage();
  }

  function renderChartShape(chart: Chart) {
    const baseClass = "w-full h-full flex flex-col items-center justify-center";

    switch (chart.shape) {
      case "oval":
        return `<div class="${baseClass} rounded-full" style="background: white; border: 2px solid #e5e7eb;">`;
      case "circle":
        return `<div class="${baseClass} rounded-full" style="background: white; border: 2px solid #e5e7eb;">`;
      case "triangle":
        return `<div class="${baseClass}" style="clip-path: polygon(50% 0%, 100% 100%, 0% 100%); background: white; border: 2px solid #e5e7eb;">`;
      case "rhombus":
        return `<div class="${baseClass}" style="transform: rotate(45deg); background: white; border: 2px solid #e5e7eb;">`;
      default:
        return `<div class="${baseClass} rounded-lg" style="background: white; border: 2px solid #e5e7eb;">`;
    }
  }
</script>

<svelte:window on:mousemove={handleMouseMove} on:mouseup={handleMouseUp} />

<!-- svelte-ignore a11y_no_noninteractive_tabindex -->
<!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
<div 

  bind:this={canvasEl}
  class="w-screen h-screen bg-gray-50 overflow-hidden relative cursor-crosshair"
  role="region"
  tabindex="0"
  on:click={handleCanvasClick}
  on:keydown={(e) => {
    if (e.key === "Escape") handleCanvasClick({} as MouseEvent);
  }}
  style="background-image: radial-gradient(circle, #d1d5db 1px, transparent 1px); background-size: {GRID_SIZE}px {GRID_SIZE}px;"
>
  <!-- SVG Layer for Arrows -->
  <svg class="absolute inset-0 w-full h-full z-0">
    <defs>
      {#each arrows as arrow (arrow.id)}
        {#if arrow.head === "triangle"}
          <marker
            id="arrowhead-{arrow.id}"
            markerWidth="12"
            markerHeight="8"
            refX="11"
            refY="4"
            orient="auto"
          >
            <polygon points="0 0, 12 4, 0 8" fill={arrow.color} />
          </marker>
        {:else if arrow.head === "diamond"}
          <marker
            id="arrowhead-{arrow.id}"
            markerWidth="12"
            markerHeight="12"
            refX="11"
            refY="6"
            orient="auto"
          >
            <polygon points="0 6, 6 0, 12 6, 6 12" fill={arrow.color} />
          </marker>
        {:else if arrow.head === "circle"}
          <marker
            id="arrowhead-{arrow.id}"
            markerWidth="12"
            markerHeight="12"
            refX="11"
            refY="6"
            orient="auto"
          >
            <circle cx="6" cy="6" r="5" fill={arrow.color} />
          </marker>
        {/if}
      {/each}
    </defs>

    {#each arrows as arrow (arrow.id)}
      <g
        aria-label="Connector from {charts.find((c) => c.id === arrow.fromId)
          ?.title || 'unknown'} to {charts.find((c) => c.id === arrow.toId)
          ?.title || 'unknown'}"
        on:click|stopPropagation={() => {
          selectedId = arrow.id;
          popupPos = { x: mousePos.x, y: mousePos.y };
          showArrowPopup = true;
        }}
        role="button"
        tabindex="0"
        on:keydown={(e) => {
          if (e.key === "Enter") {
            selectedId = arrow.id;
            popupPos = { x: mousePos.x, y: mousePos.y };
            showArrowPopup = true;
          }
        }}
        class="cursor-pointer"
      >
        <path
          d={getArrowPath(arrow)}
          stroke={arrow.color}
          stroke-width={isSelected(arrow.id) ? 3 : 2}
          fill="none"
          stroke-dasharray={arrow.style === "dashed"
            ? "10,5"
            : arrow.style === "dotted"
              ? "3,3"
              : "none"}
          class="transition-all duration-200"
          marker-end="url(#arrowhead-{arrow.id})"
          style="pointer-events: auto;"
        />
      </g>
      {#if arrow.label}
        {@const from = getChartCenter(arrow.fromId)}
        {@const to = getChartCenter(arrow.toId)}
        {#if from && to}
          <text
            x={(from.x + to.x) / 2}
            y={(from.y + to.y) / 2 - 10}
            text-anchor="middle"
            fill={arrow.color}
            font-size="12"
            font-weight="500"
          >
            {arrow.label}
          </text>
        {/if}
      {/if}
    {/each}

    <!-- Connecting Line Preview -->
    {#if isConnecting && connectStartId}
      {@const from = getChartCenter(connectStartId)}
      {#if from}
        <line
          x1={from.x}
          y1={from.y}
          x2={mousePos.x}
          y2={mousePos.y}
          stroke="#8b5cf6"
          stroke-width="2"
          stroke-dasharray="5,5"
          opacity="0.6"
        />
      {/if}
    {/if}
  </svg>

  <!-- Charts -->
  {#each charts as chart (chart.id)}
    <div
      class="absolute cursor-move select-none transition-shadow duration-150 z-10"
      class:border-blue-500={isSelected(chart.id) && !isConnecting}
      class:border-purple-500={isConnecting && connectStartId === chart.id}
      class:border-gray-200={!isSelected(chart.id) &&
        !(isConnecting && connectStartId === chart.id)}
      class:shadow-lg={isSelected(chart.id)}
      class:shadow-md={!isSelected(chart.id)}
      role="button"
      tabindex="0"
      aria-label="{chart.title} - {chart.shape} shape"
      style="
        left: {chart.x}px; 
        top: {chart.y}px; 
        width: {chart.width}px; 
        height: {chart.height}px;
        border-width: 2px;
        border-style: solid;
      "
      on:mousedown={(e) => startDrag(e, chart.id)}
      on:mouseup={() => endConnect(chart.id)}
      on:keydown={(e) => {
        if (e.key === "Enter") selectedId = chart.id;
      }}
    >
      <!-- Shape Container -->
      <div
        class="w-full h-full bg-white flex flex-col"
        style:border-radius={chart.shape === "rectangle"
          ? "8px"
          : chart.shape === "oval"
            ? "9999px"
            : chart.shape === "circle"
              ? "9999px"
              : "0"}
        style:clip-path={chart.shape === "triangle"
          ? "polygon(50% 0%, 100% 100%, 0% 100%)"
          : chart.shape === "rhombus"
            ? "polygon(50% 0%, 100% 50%, 50% 100%, 0% 50%)"
            : "none"}
      >
        <!-- Header (only for rectangle) -->
        {#if chart.shape === "rectangle"}
          <div
            class="h-7 bg-gray-50 border-b border-gray-100 flex items-center justify-between px-2 shrink-0 rounded-t-lg"
          >
            <input
              bind:value={chart.title}
              class="text-xs font-semibold text-gray-700 bg-transparent border-none focus:ring-0 outline-none flex-1"
              placeholder="Title"
            />
            <div class="flex items-center gap-0.5">
              <button
                class="p-1 text-gray-400 hover:text-purple-500 hover:bg-purple-50 rounded transition-colors"
                title="Connect arrow"
                on:click|stopPropagation={() => startConnect(chart.id)}
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="12"
                  height="12"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                >
                  <path d="M5 12h14M12 5l7 7-7 7" />
                </svg>
              </button>
              <button
                class="p-1 text-gray-400 hover:text-red-500 hover:bg-red-50 rounded transition-colors"
                title="Remove chart"
                on:click|stopPropagation={() => removeChart(chart.id)}
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="12"
                  height="12"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                >
                  <path d="M18 6L6 18M6 6l12 12" />
                </svg>
              </button>
            </div>
          </div>
        {/if}

        <!-- Content -->
        <div
          class="flex-1 p-2 flex flex-col items-center justify-center"
          class:p-3={chart.shape !== "rectangle"}
        >
          {#if chart.shape !== "rectangle"}
            <input
              bind:value={chart.title}
              class="text-xs font-semibold text-gray-700 bg-transparent border-none focus:ring-0 outline-none text-center w-full"
              placeholder="Title"
            />
          {/if}
          <textarea
            bind:value={chart.content}
            class="w-full flex-1 text-xs text-gray-600 bg-transparent border-none focus:ring-0 resize-none outline-none text-center mt-1"
            placeholder="Content..."
            style:min-height={chart.shape === "rectangle" ? "60px" : "40px"}
          ></textarea>
        </div>

        <!-- Shape-specific buttons -->
        {#if chart.shape !== "rectangle"}
          <div class="absolute top-1 right-1 flex gap-0.5 transition-opacity">
            <button
              class="p-1 text-gray-400 hover:text-purple-500 hover:bg-purple-50 rounded bg-white/90"
              title="Connect"
              on:click|stopPropagation={() => startConnect(chart.id)}
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="12"
                height="12"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <path d="M5 12h14M12 5l7 7-7 7" />
              </svg>
            </button>
            <button
              class="p-1 text-gray-400 hover:text-red-500 hover:bg-red-50 rounded bg-white/90"
              title="Remove"
              on:click|stopPropagation={() => removeChart(chart.id)}
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="12"
                height="12"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <path d="M18 6L6 18M6 6l12 12" />
              </svg>
            </button>
          </div>
        {/if}
      </div>

      <!-- Resize Handles -->
      {#if isSelected(chart.id)}
        <button
          class="absolute -bottom-1 -right-1 w-4 h-4 bg-blue-500 rounded-full cursor-se-resize"
          aria-label="Resize south-east"
          on:mousedown|stopPropagation={() => {
            isResizing = true;
            resizeHandle = "se";
          }}
        ></button>
        <button
          class="absolute -bottom-1 -left-1 w-4 h-4 bg-blue-500 rounded-full cursor-sw-resize"
          aria-label="Resize south-west"
          on:mousedown|stopPropagation={() => {
            isResizing = true;
            resizeHandle = "sw";
          }}
        ></button>
        <button
          class="absolute -top-1 -right-1 w-4 h-4 bg-blue-500 rounded-full cursor-ne-resize"
          aria-label="Resize north-east"
          on:mousedown|stopPropagation={() => {
            isResizing = true;
            resizeHandle = "ne";
          }}
        ></button>
        <button
          class="absolute -top-1 -left-1 w-4 h-4 bg-blue-500 rounded-full cursor-nw-resize"
          aria-label="Resize north-west"
          on:mousedown|stopPropagation={() => {
            isResizing = true;
            resizeHandle = "nw";
          }}
        ></button>
      {/if}
    </div>
  {/each}

  <!-- Arrow Popup -->
  {#if showArrowPopup && selectedId}
    {@const arrow = arrows.find((a) => a.id === selectedId)}
    {@const fromChart = arrow
      ? charts.find((c) => c.id === arrow.fromId)
      : null}
    {@const toChart = arrow ? charts.find((c) => c.id === arrow.toId) : null}
    {#if arrow && fromChart && toChart}
      <!-- svelte-ignore a11y_interactive_supports_focus -->
      <div 

        class="absolute bg-white rounded-lg shadow-xl border border-gray-200 p-4 z-50 w-64"
        role="dialog"
        tabindex="0"
        style="left: {popupPos.x + 10}px; top: {popupPos.y + 10}px;"
        transition:scale={{ duration: 150 }}
        on:click|stopPropagation
        on:keydown={(e) => {
          if (e.key === "Escape") showArrowPopup = false;
        }}
      >
        <div class="flex items-center justify-between mb-3">
          <h3 class="font-semibold text-gray-800">Arrow Properties</h3>
          <button
            class="text-gray-400 hover:text-red-500"
            aria-label="Delete arrow"
            on:click={() => {
              if (selectedId) removeArrow(selectedId);
            }}
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="16"
              height="16"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <path
                d="M3 6h18M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"
              />
            </svg>
          </button>
        </div>

        <div class="space-y-3">
          <!-- From/To Display -->
          <div class="flex items-center gap-2 text-xs">
            <span class="bg-blue-100 text-blue-700 px-2 py-1 rounded"
              >{fromChart.title}</span
            >
            <span class="text-gray-400">→</span>
            <span class="bg-green-100 text-green-700 px-2 py-1 rounded"
              >{toChart.title}</span
            >
          </div>

          <!-- Change Connection -->
          <div class="flex gap-1">
            <select
              class="flex-1 text-xs border border-gray-200 rounded px-2 py-1"
              value={arrow.fromId}
              on:change={(e) =>
                updateArrow(arrow.id, { fromId: e.currentTarget.value })}
            >
              {#each charts as chart}
                <option value={chart.id}>{chart.title}</option>
              {/each}
            </select>
            <select
              class="flex-1 text-xs border border-gray-200 rounded px-2 py-1"
              value={arrow.toId}
              on:change={(e) =>
                updateArrow(arrow.id, { toId: e.currentTarget.value })}
            >
              {#each charts as chart}
                <option value={chart.id}>{chart.title}</option>
              {/each}
            </select>
          </div>

          <!-- Style Options -->
          <div class="grid grid-cols-3 gap-2">
            <div>
              <label for="style-select" class="text-xs text-gray-500 block mb-1"
                >Style</label
              >
              <select
                id="style-select"
                class="w-full text-xs border border-gray-200 rounded px-2 py-1"
                value={arrow.style}
                on:change={(e) =>
                  updateArrow(arrow.id, {
                    style: e.currentTarget.value as ArrowStyle,
                  })}
              >
                <option value="solid">Solid</option>
                <option value="dashed">Dashed</option>
                <option value="dotted">Dotted</option>
              </select>
            </div>
            <div>
              <label for="head-select" class="text-xs text-gray-500 block mb-1"
                >Head</label
              >
              <select
                id="head-select"
                class="w-full text-xs border border-gray-200 rounded px-2 py-1"
                value={arrow.head}
                on:change={(e) =>
                  updateArrow(arrow.id, {
                    head: e.currentTarget.value as ArrowHead,
                  })}
              >
                <option value="triangle">Triangle</option>
                <option value="diamond">Diamond</option>
                <option value="circle">Circle</option>
                <option value="none">None</option>
              </select>
            </div>
            <div>
              <label for="curve-input" class="text-xs text-gray-500 block mb-1"
                >Curve</label
              >
              <input
                id="curve-input"
                type="range"
                min="-50"
                max="50"
                value={arrow.curvature}
                on:input={(e) =>
                  updateArrow(arrow.id, {
                    curvature: parseInt(e.currentTarget.value),
                  })}
                class="w-full"
              />
            </div>
          </div>

          <!-- Label & Color -->
          <div class="flex gap-2">
            <input
              type="text"
              placeholder="Label"
              bind:value={arrow.label}
              class="flex-1 text-xs border border-gray-200 rounded px-2 py-1"
            />
            <input
              type="color"
              bind:value={arrow.color}
              on:input={(e) =>
                updateArrow(arrow.id, { color: e.currentTarget.value })}
              class="w-8 h-8 rounded cursor-pointer"
            />
          </div>
        </div>
      </div>
    {/if}
  {/if}

  <!-- Toolbar -->
  <div
    class="absolute top-4 left-1/2 -translate-x-1/2 flex items-center gap-2 bg-white rounded-lg shadow-lg border border-gray-200 p-2 z-50"
  >
    <div class="relative">
      <button
        class="flex items-center gap-2 px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 transition-colors font-medium text-sm"
        on:click={() => (showShapeMenu = !showShapeMenu)}
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="16"
          height="16"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <rect x="3" y="3" width="18" height="18" rx="2" />
          <path d="M12 8v8M8 12h8" />
        </svg>
        Add Chart ▾
      </button>

      {#if showShapeMenu}
        <div
          class="absolute top-full left-0 mt-1 bg-white rounded-lg shadow-xl border border-gray-200 p-2 flex gap-1"
        >
          <button
            class="p-2 hover:bg-gray-50 rounded transition-colors"
            title="Rectangle"
            on:click={() => addChartShape("rectangle")}
          >
            <div
              class="w-8 h-6 bg-white border-2 border-gray-300 rounded"
            ></div>
          </button>
          <button
            class="p-2 hover:bg-gray-50 rounded transition-colors"
            title="Oval"
            on:click={() => addChartShape("oval")}
          >
            <div
              class="w-8 h-6 bg-white border-2 border-gray-300 rounded-full"
            ></div>
          </button>
          <button
            class="p-2 hover:bg-gray-50 rounded transition-colors"
            title="Circle"
            on:click={() => addChartShape("circle")}
          >
            <div
              class="w-6 h-6 bg-white border-2 border-gray-300 rounded-full"
            ></div>
          </button>
          <button
            class="p-2 hover:bg-gray-50 rounded transition-colors"
            title="Triangle"
            on:click={() => addChartShape("triangle")}
          >
            <div
              class="w-0 h-0 border-l-4 border-r-4 border-b-8 border-l-transparent border-r-transparent border-b-gray-300"
            ></div>
          </button>
          <button
            class="p-2 hover:bg-gray-50 rounded transition-colors"
            title="Rhombus"
            on:click={() => addChartShape("rhombus")}
          >
            <div
              class="w-6 h-6 bg-white border-2 border-gray-300 rotate-45"
            ></div>
          </button>
        </div>
      {/if}
    </div>

    <div class="w-px h-6 bg-gray-200"></div>

    <div class="flex items-center gap-1 text-xs text-gray-500">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        width="14"
        height="14"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
      >
        <path d="M5 12h14M12 5l7 7-7 7" />
      </svg>
      Connect
    </div>

    <div class="w-px h-6 bg-gray-200"></div>

    <div class="flex items-center gap-1 text-xs text-gray-500">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        width="14"
        height="14"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
      >
        <path
          d="M3 6h18M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"
        />
      </svg>
      Delete
    </div>

    <div class="w-px h-6 bg-gray-200"></div>

    <!-- Save/Load/Reset Buttons -->
    <div class="flex items-center gap-1">
      <button
        class="flex items-center gap-1 px-3 py-1.5 text-xs font-medium text-gray-600 hover:text-green-600 hover:bg-green-50 rounded transition-colors"
        title="Save now"
        on:click={triggerSave}
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="14"
          height="14"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <path
            d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"
          />
          <polyline points="17 21 17 13 7 13 7 21" />
          <polyline points="7 3 7 8 15 8" />
        </svg>
        Save
      </button>

      <div class="flex items-center gap-1 px-2 py-1 bg-gray-50 rounded">
        {#if saveStatus === "saved"}
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="12"
            height="12"
            viewBox="0 0 24 24"
            fill="none"
            stroke="#10b981"
            stroke-width="2"
          >
            <polyline points="20 6 9 17 4 12" />
          </svg>
          <span class="text-xs text-green-600">Saved</span>
        {:else if saveStatus === "saving"}
          <svg
            class="animate-spin"
            xmlns="http://www.w3.org/2000/svg"
            width="12"
            height="12"
            viewBox="0 0 24 24"
            fill="none"
            stroke="#6366f1"
            stroke-width="2"
          >
            <circle cx="12" cy="12" r="10" stroke-opacity="0.25" />
            <path d="M12 2a10 10 0 0 1 10 10" stroke-linecap="round" />
          </svg>
          <span class="text-xs text-indigo-600">Saving...</span>
        {:else}
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="12"
            height="12"
            viewBox="0 0 24 24"
            fill="none"
            stroke="#f59e0b"
            stroke-width="2"
          >
            <circle cx="12" cy="12" r="10" />
            <line x1="12" y1="8" x2="12" y2="12" />
            <line x1="12" y1="16" x2="12.01" y2="16" />
          </svg>
          <span class="text-xs text-amber-600">Unsaved</span>
        {/if}
      </div>

      <button
        class="flex items-center gap-1 px-3 py-1.5 text-xs font-medium text-gray-600 hover:text-blue-600 hover:bg-blue-50 rounded transition-colors"
        title="Load from storage"
        on:click={() => {
          const saved = loadFromStorage();
          if (saved) {
            charts = saved.charts;
            arrows = saved.arrows;
            saveStatus = "saved";
          }
        }}
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="14"
          height="14"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
          <polyline points="7 10 12 15 17 10" />
          <line x1="12" y1="15" x2="12" y2="3" />
        </svg>
        Load
      </button>

      <button
        class="flex items-center gap-1 px-3 py-1.5 text-xs font-medium text-gray-600 hover:text-red-600 hover:bg-red-50 rounded transition-colors"
        title="Reset to defaults"
        on:click={resetToDefaults}
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="14"
          height="14"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
        >
          <polyline points="1 4 1 10 7 10" />
          <path d="M3.51 15a9 9 0 1 0 2.13-9.36L1 10" />
        </svg>
        Reset
      </button>
    </div>
  </div>

  <!-- Instructions -->
  <div
    class="absolute bottom-4 left-4 bg-white rounded-lg shadow-lg border border-gray-200 p-3 z-50"
  >
    <p class="text-xs text-gray-500 font-medium mb-2">Controls:</p>
    <p class="text-xs text-gray-400">• Drag charts (snapped to grid)</p>
    <p class="text-xs text-gray-400">• Resize rectangles with corner handles</p>
    <p class="text-xs text-gray-400">• Click connect icon, then target</p>
    <p class="text-xs text-gray-400">• Click arrow for properties popup</p>
    <p class="text-xs text-gray-400">• Change arrow style, head, curve</p>
    <p class="text-xs text-gray-400">• Change arrow connection endpoints</p>
  </div>
</div>

<style>
  :global(body) {
    margin: 0;
    padding: 0;
    overflow: hidden;
  }

  :global(input[type="range"]) {
    appearance: none;
    -webkit-appearance: none;
    appearance: none;
    height: 6px;
    background: #e5e7eb;
    border-radius: 3px;
    outline: none;
  }

  :global(input[type="range"]::-webkit-slider-thumb) {
    appearance: none;
    -webkit-appearance: none;
    appearance: none;
    width: 14px;
    height: 14px;
    background: #6366f1;
    border-radius: 50%;
    cursor: pointer;
  }

  :global(input[type="color"]) {
    appearance: none;
    -webkit-appearance: none;
    appearance: none;
    border: none;
    padding: 0;
  }
</style>
