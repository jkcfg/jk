export class PrometheusRule {
  constructor(name, desc) {
    this.apiVersion = 'monitoring.coreos.com/v1';
    this.kind = 'PrometheusRule';
    this.metadata = Object.assign({}, (desc && desc.metadata) || {}, { name });
    this.spec = (desc && desc.spec) || undefined;
    Object.assign(this.metadata, { name });
  }
}
