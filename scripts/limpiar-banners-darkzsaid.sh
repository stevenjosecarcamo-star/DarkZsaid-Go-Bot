#!/bin/bash

python3 <<'PY'
from pathlib import Path

carpeta = Path("/etc/ssh_banners")

for p in carpeta.glob("*.banner"):
    lines = p.read_text(errors="ignore").splitlines()

    # Si la línea 3 existe y tiene el logo W/braille o font monospace, se borra
    if len(lines) >= 3:
        linea3 = lines[2]
        if "monospace" in linea3 or any('\u2800' <= ch <= '\u28ff' for ch in linea3):
            del lines[2]
            p.write_text("\n".join(lines) + "\n")
            print("Limpiado:", p)
PY
