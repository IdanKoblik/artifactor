---
layout: home
hero:
  name: Artifactor
  text: Package Version Management
  tagline: Store, retrieve, and manage versioned build artifacts with a simple REST API.
  image:
    src: /artifactor-logo.png
    alt: Artifactor
  actions:
    - theme: brand
      text: Get Started
      link: /getting-started
    - theme: alt
      text: API Reference
      link: /api
    - theme: alt
      text: View on GitHub
      link: https://github.com/IdanKoblik/artifactor
features:
  - title: Token-based Auth
    details: Every request is authenticated via an X-Api-Token header. Admin tokens can manage users; per-product permissions control who can upload, download, or delete.
  - title: Product Lifecycle
    details: Create products, grant token access, and manage multiple named versions — all through a clean REST API.
  - title: Artifact Storage
    details: Upload arbitrary files as versioned artifacts. SHA-256 checksums are computed and stored automatically.
---
