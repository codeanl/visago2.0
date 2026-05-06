import re

import requests
from bs4 import BeautifulSoup

API_BASE = "http://localhost:8080/api"

TYPE_MAP = {
    "visa_free_access": "免签",
    "visa_on_arrival": "落地签",
    "electronic_travel_authorisation": "ETA",
    "visa_online": "电子签",
}

REGION_MAP = {
    "Asia": "亚洲",
    "Europe": "欧洲",
    "Africa": "非洲",
    "Americas": "美洲",
    "Oceania": "大洋洲",
}

SPECIAL_CODE_META = {
    "RE": {"name": "留尼汪", "flag": "🇫🇷", "region": "非洲", "city": "圣但尼", "latlng": [-20.8789, 55.4481]},
    "PF": {"name": "法属波利尼西亚", "flag": "🇵🇫", "region": "大洋洲", "city": "帕皮提", "latlng": [-17.5516, -149.5585]},
    "SJ": {"name": "斯瓦尔巴", "flag": "🇳🇴", "region": "欧洲", "city": "朗伊尔城", "latlng": [78.2232, 15.6469]},
    "VG": {"name": "英属维尔京群岛", "flag": "🇻🇬", "region": "美洲", "city": "罗德城", "latlng": [18.4207, -64.64]},
    "SX": {"name": "荷属圣马丁", "flag": "🇸🇽", "region": "美洲", "city": "菲利普斯堡", "latlng": [18.0425, -63.0548]},
    "BQ": {"name": "荷兰加勒比区", "flag": "🇳🇱", "region": "美洲", "city": "克拉伦代克", "latlng": [12.2019, -68.2624]},
    "GS": {"name": "南乔治亚和南桑威奇群岛", "flag": "🇬🇸", "region": "美洲", "city": "格里特维肯", "latlng": [-54.4296, -36.5879]},
    "SH": {"name": "圣赫勒拿", "flag": "🇸🇭", "region": "非洲", "city": "詹姆斯敦", "latlng": [-15.965, -5.7089]},
    "NC": {"name": "新喀里多尼亚", "flag": "🇳🇨", "region": "大洋洲", "city": "努美阿", "latlng": [-22.2558, 166.4505]},
}


def fetch_html(url):
    headers = {"User-Agent": "Mozilla/5.0 (compatible; CodexBot/1.0; +https://openai.com)"}
    response = requests.get(url, headers=headers, timeout=40)
    response.raise_for_status()
    return response.text


def clean_text(text):
    text = re.sub(r"\[\s*\d+\s*\]", "", text)
    text = re.sub(r"\s+", " ", text.replace("\xa0", " ")).strip()
    return text


def load_restcountries():
    data = requests.get(
        "https://restcountries.com/v3.1/all?fields=cca2,translations,capital,region,latlng,flag,name",
        timeout=40,
    ).json()
    by_code = {}
    by_english = {}
    for item in data:
        zh = item.get("translations", {}).get("zho", {})
        native_zh = item.get("name", {}).get("nativeName", {}).get("zho", {})
        zh_name = zh.get("common") or native_zh.get("common") or zh.get("official") or native_zh.get("official") or item.get("name", {}).get("common")
        item["_zh_name"] = zh_name
        code = (item.get("cca2") or "").upper()
        by_code[code] = item
        by_english[item.get("name", {}).get("common")] = item
    return by_code, by_english


def load_wiki_details():
    headers = {"User-Agent": "Mozilla/5.0 (compatible; CodexBot/1.0; +https://openai.com)"}
    html = requests.get("https://en.wikipedia.org/wiki/Visa_requirements_for_Chinese_citizens", headers=headers, timeout=40).text
    soup = BeautifulSoup(html, "html.parser")
    table = soup.find_all("table")[3]
    rows = {}
    for tr in table.find_all("tr")[1:]:
        cells = tr.find_all(["td", "th"])
        vals = [clean_text(c.get_text(" ", strip=True)) for c in cells]
        if len(vals) >= 4:
            rows[vals[0].strip().lower()] = {
                "stay": vals[2],
                "note": vals[3],
            }
    return rows


def extract_stay_from_note(note):
    text = clean_text(note).lower()
    patterns = [
        r"up to (\d+)\s+days",
        r"for (\d+)\s+days",
        r"(\d+)\s+days",
        r"(\d+)\s+months",
        r"(\d+)\s+month",
        r"(\d+)\s+years",
        r"(\d+)\s+year",
    ]
    for pattern in patterns:
        match = re.search(pattern, text)
        if match:
            raw = f"{match.group(1)} {'days' if 'day' in pattern else 'months' if 'month' in pattern else 'years'}"
            return normalize_stay(raw)
    return ""


def lonlat_to_xy(latlng):
    lat, lon = latlng
    x = round((lon + 180.0) / 360.0 * 100.0, 1)
    y = round((90.0 - lat) / 180.0 * 100.0, 1)
    return x, y


def load_existing():
    items = requests.get(f"{API_BASE}/visa/free-countries", timeout=20).json()["data"]
    return {(item["code"] or "").upper(): item for item in items}


def load_supported_visa_map():
    countries = requests.get(f"{API_BASE}/visa/countries", timeout=20).json()["data"]
    visas = requests.get(f"{API_BASE}/visa/country-visas", timeout=20).json()["data"]
    country_by_id = {int(item["id"]): item for item in countries}
    result = {}
    for visa in visas:
        if not visa.get("visaFree"):
            continue
        country = country_by_id.get(int(visa["countryId"]))
        if not country:
            continue
        result[country["code"].upper()] = {
            "countryId": int(country["id"]),
            "countryName": country["name"],
            "visaId": int(visa["id"]),
            "visaName": visa["name"],
        }
    return result


def default_note(name, policy_type):
    if policy_type == "免签":
        return f"{name} 当前支持中国护照免签入境，出行前请再次确认停留时长、返程机票与住宿要求。"
    if policy_type == "落地签":
        return f"{name} 当前支持中国护照落地签入境，建议提前确认口岸、护照有效期及现场缴费要求。"
    if policy_type == "ETA":
        return f"{name} 当前需要提前申请电子旅行授权（ETA），建议出发前在线完成登记并保存审批结果。"
    return f"{name} 当前支持中国护照电子签证入境，建议出发前提前在线完成申请并保存签发结果。"


def normalize_stay(raw):
    text = clean_text(raw)
    if not text:
        return ""
    replacements = [
        (" days", "天"),
        (" day", "天"),
        (" months", "个月"),
        (" month", "个月"),
        (" year", "年"),
        (" years", "年"),
        ("within any 180-day period", "每180天内"),
        ("within any 180 day period", "每180天内"),
        ("within any 1 calendar year", "每1个日历年内"),
    ]
    for old, new in replacements:
        text = text.replace(old, new)
    text = text.replace(" / ", " / ")
    text = text.replace(".", "")
    return text


def build_dataset():
    html = fetch_html("https://passportfactory.com/zh/passport/china")
    soup = BeautifulSoup(html, "html.parser")
    existing = load_existing()
    by_code, by_english = load_restcountries()
    wiki_details = load_wiki_details()

    dataset = {}
    for category in soup.select(".visa-category"):
        filter_value = (category.get("data-category") or "").strip()
        policy_type = TYPE_MAP.get(filter_value)
        if not policy_type:
            continue
        for item in category.select("li.visa-country"):
            code = (item.get("data-country-code") or "").strip().upper()
            english_name = (item.get("data-country-name") or "").strip()
            if not code:
                continue

            meta = by_code.get(code)
            if meta:
                zh_name = meta["_zh_name"]
                flag = meta.get("flag") or ""
                region = REGION_MAP.get(meta.get("region"), meta.get("region") or "")
                city = (meta.get("capital") or [""])[0]
                x, y = lonlat_to_xy(meta.get("latlng") or [0, 0])
            else:
                special = SPECIAL_CODE_META.get(code)
                if not special:
                    continue
                zh_name = special["name"]
                flag = special["flag"]
                region = special["region"]
                city = special["city"]
                x, y = lonlat_to_xy(special["latlng"])

            old = existing.get(f"{code}-{policy_type}".upper(), {})
            wiki = wiki_details.get(english_name.lower(), {})
            wiki_stay = normalize_stay(wiki.get("stay", "")) or extract_stay_from_note(wiki.get("note", ""))
            old_stay = old.get("stay") or (existing.get(code, {}).get("stay") if policy_type == "免签" else "")
            stay = wiki_stay or old_stay or f"以{policy_type}政策为准"
            note = old.get("note") or default_note(zh_name, policy_type)
            mapping = None

            dataset[f"{code}-{policy_type}".upper()] = {
                "name": zh_name,
                "code": f"{code}-{policy_type}",
                "flag": flag,
                "region": region,
                "city": city,
                "policyType": policy_type,
                "stay": stay,
                "note": note[:255],
                "mapX": x,
                "mapY": y,
                "enabled": True,
                "keywords": [item for item in [zh_name, code, city, region, policy_type] if item],
                "mappingCode": code,
            }
    return dataset


def sync_database():
    dataset = build_dataset()
    existing = load_existing()
    supported_map = load_supported_visa_map()
    created = 0
    updated = 0

    for key, item in dataset.items():
        mapping = supported_map.get(item["mappingCode"])
        payload = {
            "name": item["name"],
            "code": item["code"],
            "flag": item["flag"],
            "region": item["region"],
            "city": item["city"],
            "policyType": item["policyType"],
            "stay": item["stay"],
            "note": item["note"],
            "mapX": item["mapX"],
            "mapY": item["mapY"],
            "enabled": item["enabled"],
            "supportedCountryId": mapping["countryId"] if mapping else 0,
            "supportedCountryName": mapping["countryName"] if mapping else "",
            "supportedVisaId": mapping["visaId"] if mapping else 0,
            "supportedVisaName": mapping["visaName"] if mapping else "",
            "keywords": item["keywords"],
        }
        if key in existing:
            requests.put(f"{API_BASE}/visa/free-countries/{existing[key]['id']}", json=payload, timeout=20).raise_for_status()
            updated += 1
        else:
            requests.post(f"{API_BASE}/visa/free-countries", json=payload, timeout=20).raise_for_status()
            created += 1

    target_keys = set(dataset.keys())
    for code, row in existing.items():
        if code not in target_keys:
            requests.delete(f"{API_BASE}/visa/free-countries/{row['id']}", timeout=20).raise_for_status()

    counts = {"免签": 0, "落地签": 0, "ETA": 0, "电子签": 0}
    for item in dataset.values():
        counts[item["policyType"]] += 1
    return {"created": created, "updated": updated, "total": len(dataset), "counts": counts}


if __name__ == "__main__":
    print(sync_database())
