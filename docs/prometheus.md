---
layout: page
title: Prometheus metrics
hide_hero: true
show_sidebar: false
menubar: docs-menu
---

# Metrics

The port `9096` is bound by Fed for our admin service. This HTTP server has endpoints for Prometheus metrics (`GET /metrics`), readiness checks (`GET /ready`), and liveness checks (`GET /live`).