# AIP Food Lookup

**A fast, reliable lookup tool for AIP (Autoimmune Protocol) food compatibility.**  
Identify whether a food item aligns with the AIP diet directly from the command line or integration.

---

## Table of Contents

- [Overview](#overview)  
- [Features](#features)  
- [Quick Start](#quick-start)  
- [Project Structure](#project-structure)  
- [Usage Examples](#usage-examples)  
- [Environment & Configuration](#environment--configuration)  
- [Contributing](#contributing)  
- [License](#license)  
- [Contact](#contact)

---

## Overview

AIP Food Lookup is a lightweight, user-friendly tool designed to quickly check whether a food complies with the AIP diet. Built in Go, it’s ideal for command-line use or integration into other workflows. It’s open source and freely available under the MIT License.:contentReference[oaicite:0]{index=0}

---

## Features

- ✅ Lookup food status (AIP-compliant or not)  
- 🆓 Lightweight and fast — written in Go  
- 📦 Simple CLI usage; easy to integrate into other tools or scripts  
- 🆓 Free and open source (MIT License):contentReference[oaicite:1]{index=1}  

---

## Quick Start

```bash
# Clone the repository
git clone https://github.com/CalypsoSys/aip_food_lookup.git
cd aip_food_lookup

# Build the project
go build -o aipfoodlookup

# Run a lookup (example: "tomato")
./aipfoodlookup tomato
