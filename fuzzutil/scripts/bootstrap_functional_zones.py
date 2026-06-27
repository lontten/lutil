#!/usr/bin/env python3
"""Bootstrap functional_zones_national.json from Wikipedia/hunan source texts."""
import json
import re
import os
import unicodedata

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
WIKI = os.environ.get(
    "WIKI_HTZ",
    os.path.join(os.path.dirname(__file__), "..", "..", "..", "agent-tools", "683083bb-d989-4c84-85d9-d3d4ea9edad7.txt"),
)
HUNAN = os.environ.get(
    "HUNAN_ETDZ",
    os.path.join(os.path.dirname(__file__), "..", "..", "..", "agent-tools", "712c6e65-9e4b-4906-a6ac-33fd71524202.txt"),
)
OUT = os.path.join(ROOT, "data", "functional_zones_national.json")

# zone full name -> admin district (民政部区县名)
DISTRICT = {
    "中关村科技园区": "海淀区",
    "天津滨海高新技术产业开发区": "滨海新区",
    "石家庄高新技术产业开发区": "裕华区",
    "保定高新技术产业开发区": "竞秀区",
    "唐山高新技术产业开发区": "路北区",
    "燕郊高新技术产业开发区": "三河市",
    "承德高新技术产业开发区": "双桥区",
    "太原高新技术产业开发区": "小店区",
    "长治高新技术产业开发区": "潞州区",
    "包头稀土高新技术产业开发区": "稀土高新区",
    "呼和浩特金山高新技术产业开发区": "赛罕区",
    "鄂尔多斯高新技术产业开发区": "康巴什区",
    "沈阳高新技术产业开发区": "浑南区",
    "大连高新技术产业园区": "甘井子区",
    "鞍山高新技术产业开发区": "立山区",
    "营口高新技术产业开发区": "西市区",
    "辽阳高新技术产业开发区": "宏伟区",
    "本溪高新技术产业开发区": "明山区",
    "锦州高新技术产业开发区": "太和区",
    "阜新高新技术产业开发区": "细河区",
    "长春高新技术产业开发区": "朝阳区",
    "长春净月高新技术产业开发区": "净月区",
    "吉林高新技术产业开发区": "丰满区",
    "延吉高新技术产业开发区": "延吉市",
    "通化医药高新技术产业开发区": "东昌区",
    "哈尔滨高新技术产业开发区": "松北区",
    "大庆高新技术产业开发区": "萨尔图区",
    "齐齐哈尔高新技术产业开发区": "龙沙区",
    "上海张江高新技术产业开发区": "浦东新区",
    "上海紫竹高新技术产业开发区": "闵行区",
    "南京高新技术产业开发区": "玄武区",
    "苏州高新技术产业开发区": "虎丘区",
    "昆山高新技术产业开发区": "昆山市",
    "无锡高新技术产业开发区": "新吴区",
    "江阴高新技术产业开发区": "江阴市",
    "常州高新技术产业开发区": "新北区",
    "武进高新技术产业开发区": "武进区",
    "泰州医药高新技术产业开发区": "医药高新区",
    "徐州高新技术产业开发区": "铜山区",
    "苏州工业园区": "虎丘区",
    "南通高新技术产业开发区": "崇川区",
    "镇江高新技术产业开发区": "丹徒区",
    "盐城高新技术产业开发区": "盐都区",
    "连云港高新技术产业开发区": "海州区",
    "扬州高新技术产业开发区": "邗江区",
    "常熟高新技术产业开发区": "常熟市",
    "宿迁高新技术产业开发区": "宿城区",
    "淮安高新技术产业开发区": "清江浦区",
    "杭州高新技术产业开发区": "滨江区",
    "萧山临江高新技术产业开发区": "萧山区",
    "宁波高新技术产业开发区": "鄞州区",
    "绍兴高新技术产业开发区": "越城区",
    "温州高新技术产业开发区": "龙湾区",
    "衢州高新技术产业开发区": "柯城区",
    "湖州莫干山高新技术产业开发区": "德清县",
    "嘉兴高新技术产业开发区": "秀洲区",
    "合肥高新技术产业开发区": "蜀山区",
    "蚌埠高新技术产业开发区": "禹会区",
    "芜湖高新技术产业开发区": "弋江区",
    "马鞍山慈湖高新技术产业开发区": "花山区",
    "铜陵狮子山高新技术产业开发区": "铜官区",
    "淮南高新技术产业开发区": "田家庵区",
    "滁州高新技术产业开发区": "南谯区",
    "福州高新技术产业开发区": "马尾区",
    "厦门火炬高技术产业开发区": "湖里区",
    "泉州高新技术产业开发区": "丰泽区",
    "莆田高新技术产业开发区": "城厢区",
    "漳州高新技术产业开发区": "龙文区",
    "三明高新技术产业开发区": "三元区",
    "龙岩高新技术产业开发区": "新罗区",
    "南昌高新技术产业开发区": "青山湖区",
    "新余高新技术产业开发区": "渝水区",
    "景德镇高新技术产业开发区": "昌江区",
    "鹰潭高新技术产业开发区": "月湖区",
    "抚州高新技术产业开发区": "临川区",
    "吉安高新技术产业开发区": "吉州区",
    "赣州高新技术产业开发区": "章贡区",
    "九江共青城高新技术产业开发区": "共青城市",
    "宜春丰城高新技术产业开发区": "丰城市",
    "济南高新技术产业开发区": "历下区",
    "威海火炬高技术产业开发区": "环翠区",
    "青岛高新技术产业开发区": "城阳区",
    "潍坊高新技术产业开发区": "奎文区",
    "淄博高新技术产业开发区": "张店区",
    "济宁高新技术产业开发区": "任城区",
    "烟台高新技术产业开发区": "莱山区",
    "临沂高新技术产业开发区": "罗庄区",
    "泰安高新技术产业开发区": "泰山区",
    "枣庄高新技术产业开发区": "薛城区",
    "德州高新技术产业开发区": "德城区",
    "莱芜高新技术产业开发区": "莱城区",
    "黄河三角洲农业高新技术产业示范区": "东营区",
    "郑州高新技术产业开发区": "中原区",
    "洛阳高新技术产业开发区": "涧西区",
    "安阳高新技术产业开发区": "龙安区",
    "南阳高新技术产业开发区": "卧龙区",
    "新乡高新技术产业开发区": "红旗区",
    "平顶山高新技术产业开发区": "卫东区",
    "焦作高新技术产业开发区": "山阳区",
    "信阳高新技术产业开发区": "平桥区",
    "武汉东湖新技术开发区": "洪山区",
    "襄阳高新技术产业开发区": "襄州区",
    "宜昌高新技术产业开发区": "西陵区",
    "孝感高新技术产业开发区": "孝南区",
    "荆门高新技术产业开发区": "掇刀区",
    "仙桃高新技术产业开发区": "仙桃市",
    "随州高新技术产业开发区": "曾都区",
    "黄冈高新技术产业开发区": "黄州区",
    "咸宁高新技术产业开发区": "咸安区",
    "荆州高新技术产业开发区": "荆州区",
    "黄石大冶湖高新技术产业开发区": "大冶市",
    "潜江高新技术产业开发区": "潜江市",
    "长沙高新技术产业开发区": "岳麓区",
    "株洲高新技术产业开发区": "天元区",
    "湘潭高新技术产业开发区": "岳塘区",
    "益阳高新技术产业开发区": "赫山区",
    "衡阳高新技术产业开发区": "蒸湘区",
    "郴州高新技术产业开发区": "苏仙区",
    "常德高新技术产业开发区": "武陵区",
    "怀化高新技术产业开发区": "鹤城区",
    "广州高新技术产业开发区": "黄埔区",
    "深圳市高新技术产业园区": "南山区",
    "中山火炬高技术产业开发区": "火炬开发区",
    "佛山高新技术产业开发区": "禅城区",
    "惠州仲恺高新技术产业开发区": "仲恺区",
    "珠海高新技术产业开发区": "香洲区",
    "东莞松山湖高新技术产业开发区": "大岭山镇",
    "肇庆高新技术产业开发区": "肇庆高新区",
    "江门高新技术产业开发区": "江海区",
    "源城高新技术产业开发区": "源城区",
    "清远高新技术产业开发区": "清城区",
    "汕头高新技术产业开发区": "龙湖区",
    "湛江高新技术产业开发区": "坡头区",
    "茂名高新技术产业开发区": "茂南区",
    "南宁高新技术产业开发区": "西乡塘区",
    "桂林高新技术产业开发区": "七星区",
    "柳州高新技术产业开发区": "柳东区",
    "北海高新技术产业开发区": "银海区",
    "海口高新技术产业开发区": "秀英区",
    "重庆高新技术产业开发区": "九龙坡区",
    "璧山高新技术产业开发区": "璧山区",
    "荣昌高新技术产业开发区": "荣昌区",
    "永川高新技术产业开发区": "永川区",
    "成都高新技术产业开发区": "武侯区",
    "绵阳高新技术产业开发区": "涪城区",
    "自贡高新技术产业开发区": "自流井区",
    "内江高新技术产业开发区": "东兴区",
    "乐山高新技术产业开发区": "市中区",
    "泸州高新技术产业开发区": "江阳区",
    "攀枝花钒钛高新技术产业开发区": "东区",
    "德阳高新技术产业开发区": "旌阳区",
    "贵阳高新技术产业开发区": "观山湖区",
    "安顺高新技术产业开发区": "西秀区",
    "遵义高新技术产业开发区": "红花岗区",
    "昆明高新技术产业开发区": "五华区",
    "玉溪高新技术产业开发区": "红塔区",
    "楚雄高新技术产业开发区": "楚雄市",
    "西安高新技术产业开发区": "雁塔区",
    "宝鸡高新技术产业开发区": "渭滨区",
    "杨凌农业高新技术产业示范区": "杨陵区",
    "渭南高新技术产业开发区": "临渭区",
    "榆林高新技术产业开发区": "榆阳区",
    "咸阳高新技术产业开发区": "秦都区",
    "安康高新技术产业开发区": "汉滨区",
    "兰州高新技术产业开发区": "城关区",
    "白银高新技术产业开发区": "白银区",
    "青海高新技术产业开发区": "城东区",
    "银川高新技术产业开发区": "金凤区",
    "石嘴山高新技术产业开发区": "大武口区",
    "乌鲁木齐高新技术产业开发区": "新市区",
    "昌吉高新技术产业开发区": "昌吉市",
    "新疆生产建设兵团石河子高新技术产业开发区": "石河子市",
    "克拉玛依高新技术产业开发区": "克拉玛依区",
    "河北雄安高新技术产业开发区": "雄县",
    # 2024+ additions
    "广州花都经济技术开发区": "花都区",
    "贵溪经济技术开发区": "贵溪市",
    "涪陵经济技术开发区": "涪陵区",
    "沈阳金融商贸经济技术开发区": "沈河区",
    # 国家级经开区（行政归属）
    "北京经济技术开发区": "大兴区",
    "天津经济技术开发区": "滨海新区",
    "广州经济技术开发区": "黄埔区",
    "苏州工业园区": "虎丘区",
    "昆山经济技术开发区": "昆山市",
    "南京经济技术开发区": "栖霞区",
    "杭州经济技术开发区": "钱塘区",
    "成都经济技术开发区": "龙泉驿区",
    "武汉经济技术开发区": "蔡甸区",
    "郑州经济技术开发区": "管城回族区",
    "长沙经济技术开发区": "长沙县",
    "合肥经济技术开发区": "蜀山区",
    "大连经济技术开发区": "金州区",
    "青岛经济技术开发区": "黄岛区",
    "宁波经济技术开发区": "北仑区",
    "厦门海沧台商投资区": "海沧区",
    "洋浦经济开发区": "洋浦区",
    "上海漕河泾新兴技术开发区": "徐汇区",
}

# Extra aliases not auto-generated
EXTRA = {
    "中关村科技园区": ["中关村科技园"],
    "武汉东湖新技术开发区": ["武汉东湖高新区", "东湖新技术开发区", "中国光谷", "光谷"],
    "上海张江高新技术产业开发区": ["张江高新区", "张江高科技园区"],
    "苏州工业园区": ["苏州工业园"],
    "北京经济技术开发区": ["北京经开区", "亦庄经济技术开发区", "亦庄开发区"],
    "广州经济技术开发区": ["广州经开区"],
    "郑州高新技术产业开发区": ["郑州高新区"],
    "深圳市高新技术产业园区": ["深圳高新区", "深圳高新技术产业园区"],
    "东莞松山湖高新技术产业开发区": ["松山湖高新区", "东莞松山湖"],
}

HITEC_SUFFIXES = [
    "高新技术产业开发区",
    "高新技术产业园区",
    "新技术产业开发区",
    "高技术产业开发区",
    "火炬高技术产业开发区",
    "农业高新技术产业示范区",
    "科技园区",
    "工业园区",
]
ETDZ_SUFFIXES = [
    "经济技术开发区",
    "经济开发区",
    "出口加工区",
    "台商投资区",
    "融侨经济技术开发区",
    "招商局经济技术开发区",
]


def parse_wiki_hitech(path):
    if not os.path.isfile(path):
        return []
    text = open(path, encoding="utf-8").read()
    # Table section: from "## 名单" until "## 注释"
    start = text.find("## 名单")
    end = text.find("## 注释")
    if start >= 0 and end > start:
        text = text[start:end]
    patterns = [
        r"\[([^\]]+高新技术产业开发区)\]",
        r"\[([^\]]+高新技术产业园区)\]",
        r"\[([^\]]+新技术产业开发区)\]",
        r"\[([^\]]+火炬高技术产业开发区)\]",
        r"\[([^\]]+农业高新技术产业示范区)\]",
        r"\[([^\]]+科技园区)\]",
        r"\[([^\]]+工业园区)\]",
    ]
    names = []
    skip = {"北京市", "天津市", "上海市", "重庆市", "河北省", "山西省", "内蒙古自治区",
            "辽宁省", "吉林省", "黑龙江省", "江苏省", "浙江省", "安徽省", "福建省",
            "江西省", "山东省", "河南省", "湖北省", "湖南省", "广东省", "广西壮族自治区",
            "海南省", "四川省", "贵州省", "云南省", "陕西省", "甘肃省", "青海省",
            "宁夏回族自治区", "新疆维吾尔自治区", "新疆生产建设兵团"}
    for pat in patterns:
        for m in re.finditer(pat, text):
            n = m.group(1)
            if n in skip or len(n) < 4:
                continue
            if n not in names:
                names.append(n)
    return names


def parse_hunan_etdz(path):
    if not os.path.isfile(path):
        return []
    text = open(path, encoding="utf-8").read()
    start = text.find("219家国家级经济技术开发区名单")
    end = text.find("145家国家级高新技术产业开发区名单")
    if start < 0:
        return []
    if end < 0 or end <= start:
        end = text.find("2 、国家级高新技术产业开发区")
    if end < 0 or end <= start:
        end = len(text)
    block = text[start:end]
    names = []
    keywords = ("经济技术开发区", "经济开发区", "出口加工区", "台商投资区", "产业园区")
    for line in block.splitlines():
        line = line.strip()
        if not line or line in ("名称", "数量", "省份", "219家国家级经济技术开发区名单"):
            continue
        if line.isdigit():
            continue
        if not any(k in line for k in keywords):
            continue
        for part in re.split(r"[、，,\s]+", line):
            part = part.strip()
            if not part or part in ("名称", "数量", "省份"):
                continue
            if any(k in part for k in keywords) or part.endswith("开发区"):
                if part not in names:
                    names.append(part)
    return names


def city_short(name):
    for suf in HITEC_SUFFIXES + ETDZ_SUFFIXES:
        if name.endswith(suf):
            return name[: -len(suf)]
    return name


def auto_aliases(name, ztype):
    aliases = [name]
    short = city_short(name)
    if len(short) >= 2:
        if ztype == "hitech":
            if name.endswith("高新技术产业开发区"):
                aliases.append(short + "高新区")
            elif name.endswith("高新技术产业园区"):
                aliases.append(short + "高新区")
            elif name.endswith("新技术产业开发区"):
                aliases.append(short + "高新区")
        elif ztype == "etdz":
            if "经济技术开发区" in name:
                aliases.append(short + "经开区")
            elif name.endswith("经济开发区"):
                aliases.append(short + "经开区")
    return list(dict.fromkeys(aliases))


def infer_district(name):
    if name in DISTRICT:
        return DISTRICT[name]
    short = city_short(name)
    if short.endswith("市") or short.endswith("县") or short.endswith("区"):
        return short
    # prefecture-level city default: append 区 for main urban
    if len(short) >= 2 and not short.endswith("市"):
        return short + "区"
    return short + "区"


def main():
    hitech = parse_wiki_hitech(WIKI)
    etdz = parse_hunan_etdz(HUNAN)

    # Manual additions (2024 upgrades, missing from stale sources)
    manual_hitech = [
        "河北雄安高新技术产业开发区",
        "克拉玛依高新技术产业开发区",
        "信阳高新技术产业开发区",
        "滁州高新技术产业开发区",
        "遵义高新技术产业开发区",
    ]
    manual_etdz = [
        "广州花都经济技术开发区",
        "贵溪经济技术开发区",
        "涪陵经济技术开发区",
        "沈阳金融商贸经济技术开发区",
    ]
    for n in manual_hitech:
        if n not in hitech:
            hitech.append(n)
    for n in manual_etdz:
        if n not in etdz:
            etdz.append(n)

    zones = []
    for name in hitech:
        ztype = "hitech"
        district = infer_district(name)
        aliases = auto_aliases(name, ztype)
        for a in EXTRA.get(name, []):
            if a not in aliases:
                aliases.append(a)
        zones.append({
            "zone_id": re.sub(r"[^a-z0-9]+", "_", name.lower())[:40] or "zone",
            "type": ztype,
            "name": name,
            "district": district,
            "aliases": aliases,
        })
    for name in etdz:
        ztype = "etdz"
        district = infer_district(name)
        if name in DISTRICT:
            district = DISTRICT[name]
        elif name == "北京经济技术开发区":
            district = "大兴区"
        elif name == "广州经济技术开发区":
            district = "黄埔区"
        elif name == "苏州工业园区":
            district = "虎丘区"
        aliases = auto_aliases(name, ztype)
        for a in EXTRA.get(name, []):
            if a not in aliases:
                aliases.append(a)
        zones.append({
            "zone_id": re.sub(r"[^a-z0-9]+", "_", name.lower())[:40] or "zone",
            "type": ztype,
            "name": name,
            "district": district,
            "aliases": aliases,
        })

    os.makedirs(os.path.dirname(OUT), exist_ok=True)
    with open(OUT, "w", encoding="utf-8") as f:
        json.dump({"zones": zones}, f, ensure_ascii=False, indent=2)
    print(f"wrote {len(zones)} zones ({len(hitech)} hitech + {len(etdz)} etdz) -> {OUT}")


if __name__ == "__main__":
    main()
