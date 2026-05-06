import json
import re
import urllib.parse

import requests

API_BASE = "http://localhost:8080/api"


def has_ascii_letters(text):
    return any(("a" <= ch <= "z") or ("A" <= ch <= "Z") for ch in str(text or ""))


def translate_text(text):
    raw = str(text or "").strip()
    if not raw:
        return ""
    chunks = re.split(r"(?<=[.;!?])\s+", raw)
    translated_parts = []
    for chunk in chunks:
        chunk = chunk.strip()
        if not chunk:
            continue
        url = (
            "https://translate.googleapis.com/translate_a/single"
            f"?client=gtx&sl=auto&tl=zh-CN&dt=t&q={urllib.parse.quote(chunk)}"
        )
        response = requests.get(url, timeout=30)
        response.raise_for_status()
        data = response.json()
        translated_parts.append("".join(part[0] for part in data[0] if part and part[0]))
    translated = "".join(translated_parts)
    translated = re.sub(r"\s+", " ", translated).strip()
    return translated


def main():
    items = requests.get(f"{API_BASE}/visa/free-countries", timeout=30).json()["data"]
    updated = 0
    skipped = 0

    for item in items:
        next_stay = item.get("stay", "")
        next_note = item.get("note", "")
        changed = False

        if has_ascii_letters(next_stay):
            next_stay = translate_text(next_stay)
            changed = True

        if has_ascii_letters(next_note):
            next_note = translate_text(next_note)
            changed = True

        if not changed:
            skipped += 1
            continue

        payload = {
            "name": item["name"],
            "code": item["code"],
            "flag": item.get("flag", ""),
            "region": item.get("region", ""),
            "city": item.get("city", ""),
            "policyType": item.get("policyType", ""),
            "stay": next_stay,
            "note": next_note,
            "mapX": item.get("mapX", 0),
            "mapY": item.get("mapY", 0),
            "enabled": item.get("enabled", True),
            "supportedCountryId": item.get("supportedCountryId", 0),
            "supportedCountryName": item.get("supportedCountryName", ""),
            "supportedVisaId": item.get("supportedVisaId", 0),
            "supportedVisaName": item.get("supportedVisaName", ""),
            "keywords": item.get("keywords", []),
        }
        requests.put(f"{API_BASE}/visa/free-countries/{item['id']}", json=payload, timeout=30).raise_for_status()
        updated += 1

    print(json.dumps({"updated": updated, "skipped": skipped, "total": len(items)}, ensure_ascii=False))


if __name__ == "__main__":
    main()
