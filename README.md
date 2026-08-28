<p align="center">
  <img src="tray/assets/icon.png" alt="Look News" width="420" />
</p>

<h1 align="center">Look News</h1>

<p align="center">
  RSS/Atom/RDF news aggregator.<br>
  Filters by relevance with AI and delivers directly to the system tray.
</p>

<p align="center">
  <a href="https://github.com/clemilsonazevedo/look-news/releases/latest">
    <img src="https://img.shields.io/github/v/release/clemilsonazevedo/look-news" alt="Release">
  </a>
</p>

## Installation

**macOS / Linux**

```bash
curl -fsSL https://raw.githubusercontent.com/clemilsonazevedo/look-news/main/install.sh | bash
```

**Windows** (PowerShell)

```powershell
irm https://raw.githubusercontent.com/clemilsonazevedo/look-news/main/install.ps1 | iex
```

Done. Just open the app and it will appear as an icon in the system tray.

## How it works

Runs in the background. Click the icon to see the most relevant news from your sources.

<p align="center">
  <img src="tray/assets/tray-diagram.png" alt="System tray placement" width="600" />
</p>

## Sources to get started

| Source       | Feed                                      |
| ------------ | ----------------------------------------- |
| Huncoding    | `https://huncoding.com/feed.xml`          |
| TabNews      | `https://www.tabnews.com.br/recentes/rss` |
| Tecnoblog    | `https://tecnoblog.net/feed/`             |
| TechCrunch   | `https://techcrunch.com/feed/`            |
| SecurityWeek | `https://www.securityweek.com/feed/`      |

Add any RSS, Atom, or RDF feed. To find new ones, search for “RSS” on the website or ask an AI:

> Find RSS sources about [your topic]. Return the direct feed links.

## Structure

| Part     | Description                | Docs                             |
| -------- | -------------------------- | -------------------------------- |
| **Tray** | Electron app (system tray) | [tray/README.md](tray/README.md) |
| **API**  | Go backend + AI            | [api/README.md](api/README.md)   |

## License

MIT
