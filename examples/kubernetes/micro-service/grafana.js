class TimePicker {
  constructor(
    refresh_intervals = [
      '5s',
      '10s',
      '30s',
      '1m',
      '5m',
      '15m',
      '30m',
      '1h',
      '2h',
      '1d',
    ],
    time_options = [
      '5m',
      '15m',
      '1h',
      '6h',
      '12h',
      '24h',
      '2d',
      '7d',
      '30d',
    ],
  ) {
    Object.assign(this, {
      refresh_intervals,
      time_options,
    });
  }
}

export class Legend {
  constructor({
    show = true,
    values = false,
    min = false,
    max = false,
    current = false,
    total = false,
    avg = false,
    alignAsTable = false,
    rightSide = false,
    hideEmpty = undefined,
    hideZero = undefined,
    sort = undefined,
    sortDesc = undefined,
  } = {}) {
    Object.assign(this, {
      show,
      values,
      min,
      max,
      current,
      total,
      avg,
      alignAsTable,
      rightSide,
      hideEmpty,
      hideZero,
      sort,
      sortDesc,
    });
  }
}

export class YAxis {
  constructor({
    format = 'short',
    min = null,
    max = null,
    label = null,
    show = true,
    logBase = 1,
    decimals = undefined,
  } = {}) {
    Object.assign(this, {
      format,
      min,
      max,
      label,
      show,
      logBase,
      decimals,
    });
  }
}

export class XAxis {
  constructor({
    show = true,
    mode = 'time',
    name = null,
    values = undefined,
    buckets = null,
  } = {}) {
    Object.assign(this, {
      show,
      mode,
      name,
      values,
      buckets,
    });
  }
}

export class Prometheus {
  constructor(expr, {
    format = 'time_series',
    intervalFactor = 2,
    legendFormat = '',
    datasource = undefined,
    interval = undefined,
    instant = undefined,
  } = {}) {
    Object.assign(this, {
      datasource,
      expr,
      format,
      intervalFactor,
      legendFormat,
      interval,
      instant,
    });
  }
}


const targetMixin = {
  targets: [],

  addTarget(target) {
    target.refId = String.fromCharCode('A'.charCodeAt(0) + this.targets.length);
    this.targets.push(target);
    return this;
  },

  addTargets(targets) {
    targets.forEach(t => this.addTarget(t));
    return this;
  },
};

export class Graph {
  constructor(title, {
    span = undefined,
    fill = 1,
    linewidth = 1,
    description = undefined,
    min_span = undefined,
    lines = true,
    datasource = null,
    points = false,
    pointradius = 5,
    bars = false,
    height = undefined,
    nullPointMode = 'null',
    dashes = false,
    stack = false,
    repeat = null,
    repeatDirection = undefined,
    sort = 0,
    legend = new Legend(),
    aliasColors = {},
    thresholds = [],
    transparent = undefined,
    value_type = 'individual',
    yAxis = [
      new YAxis(),
      new YAxis(),
    ],
    xAxis = new XAxis(),
  } = {}) {
    Object.assign(this, {
      title,
      span,
      min_span,
      type: 'graph',
      datasource,
      targets: [],
      description,
      height,
      renderer: 'flot',
      yaxes: yAxis,
      xaxis: xAxis,
      lines,
      fill,
      linewidth,
      dashes,
      dashLength: 10,
      spaceLength: 10,
      points,
      pointradius,
      bars,
      stack,
      percentage: false,
      legend,
      nullPointMode,
      steppedLine: false,
      tooltip: {
        value_type,
        shared: true,
        sort,
      },
      timeFrom: null,
      timeShift: null,
      transparent,
      aliasColors,
      repeat,
      repeatDirection,
      seriesOverrides: [],
      thresholds,
      links: [],
    });
  }
}

Object.assign(Graph.prototype, targetMixin);

export const Tooltip = {
  DEFAULT: 0,
  SHARED_CROSSHAIR: 1,
  SHARED_TOOLTIP: 2,
};

export class Dashboard {
  constructor(title, {
    editable = false,
    style = 'dark',
    tags = [],
    time_from = 'now-6h',
    time_to = 'now',
    timezone = 'browser',
    refresh = '',
    timePicker = new TimePicker(),
    graphTooltip = Tooltip.DEFAULT,
    hideControls = false,
    schemaVersion = 16,
    uid = '',
    description = undefined,
  } = {}) {
    Object.assign(this, {
      annotations: {
        list: [],
      },
      uid,
      editable,
      description,
      gnetId: null,
      graphTooltip,
      hideControls,
      id: null,
      links: [],
      panels: [],
      refresh,
      schemaVersion,
      style,
      tags,
      time: {
        from: time_from,
        to: time_to,
      },
      timezone,
      timepicker: timePicker,
      title,
      version: 0,
    });
  }

  static numPanels(panels) {
    return panels.reduce((count, panel) => {
      if ('panels' in panel) {
        return count + Dashboard.numPanels(panel.panels);
      }
      return count + 1;
    }, 0, 0);
  }

  static setPanelsId(panels, startId) {
    panels.forEach((p) => {
      p.id = startId + 1;
      if ('panels' in p) {
        startId = Dashboard.setPanelsId(p.panels, startId);
      }
    });
    return startId;
  }

  addPanel(panel, {
    gridPos = {
      x: 0, y: 0, w: 24, h: 7,
    },
  } = {}) {
    const nextId = Dashboard.numPanels(this.panels) + 1;
    Dashboard.setPanelsId([panel], nextId);
    panel.gridPos = gridPos;
    this.panels.push(panel);
    return this;
  }
}
