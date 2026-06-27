const koData = [
    { id: 1, kanji: "東風解凍", en: "East wind melts the ice", startMonth: 2, startDay: 4 },
    { id: 2, kanji: "黄鶯睍睆", en: "Bush warblers begin to sing", startMonth: 2, startDay: 9 },
    { id: 3, kanji: "魚上氷", en: "Fish emerge from beneath the ice", startMonth: 2, startDay: 14 },
    { id: 4, kanji: "土脉潤起", en: "Rain moistens the soil", startMonth: 2, startDay: 19 },
    { id: 5, kanji: "霞始靆", en: "Mist begins to linger", startMonth: 2, startDay: 24 },
    { id: 6, kanji: "草木萌動", en: "Grass and trees begin to sprout", startMonth: 3, startDay: 1 },
    { id: 7, kanji: "蟄虫啓戸", en: "Hibernating insects surface", startMonth: 3, startDay: 6 },
    { id: 8, kanji: "桃始笑", en: "Peach blossoms begin to bloom", startMonth: 3, startDay: 11 },
    { id: 9, kanji: "菜虫化蝶", en: "Caterpillars become butterflies", startMonth: 3, startDay: 16 },
    { id: 10, kanji: "雀始巣", en: "Sparrows begin to nest", startMonth: 3, startDay: 21 },
    { id: 11, kanji: "櫻始開", en: "Cherry blossoms begin to bloom", startMonth: 3, startDay: 26 },
    { id: 12, kanji: "雷乃発声", en: "Thunder roars in the distance", startMonth: 3, startDay: 31 },
    { id: 13, kanji: "玄鳥至", en: "Swallows return from the south", startMonth: 4, startDay: 5 },
    { id: 14, kanji: "鴻雁北", en: "Wild geese fly north", startMonth: 4, startDay: 10 },
    { id: 15, kanji: "虹始見", en: "Rainbows begin to appear", startMonth: 4, startDay: 15 },
    { id: 16, kanji: "葭始生", en: "Reeds begin to sprout", startMonth: 4, startDay: 20 },
    { id: 17, kanji: "霜止出苗", en: "Frost ends, rice seedlings appear", startMonth: 4, startDay: 25 },
    { id: 18, kanji: "牡丹華", en: "Peonies bloom", startMonth: 4, startDay: 30 },
    { id: 19, kanji: "蛙始鳴", en: "Frogs begin to croak", startMonth: 5, startDay: 6 },
    { id: 20, kanji: "蚯蚓出", en: "Earthworms surface from the ground", startMonth: 5, startDay: 11 },
    { id: 21, kanji: "竹笋生", en: "Bamboo shoots sprout", startMonth: 5, startDay: 16 },
    { id: 22, kanji: "蚕起食桑", en: "Silkworms feed on mulberry leaves", startMonth: 5, startDay: 21 },
    { id: 23, kanji: "紅花栄", en: "Safflowers bloom in vibrant red", startMonth: 5, startDay: 26 },
    { id: 24, kanji: "麦秋至", en: "Wheat ripens and turns golden", startMonth: 5, startDay: 31 },
    { id: 25, kanji: "螳螂生", en: "Praying mantises hatch", startMonth: 6, startDay: 6 },
    { id: 26, kanji: "腐草為蛍", en: "Rotten grass transforms into fireflies", startMonth: 6, startDay: 11 },
    { id: 27, kanji: "梅子黄", en: "Plums turn yellow on the branch", startMonth: 6, startDay: 16 },
    { id: 28, kanji: "乃東枯", en: "Prunella flowers wither", startMonth: 6, startDay: 21 },
    { id: 29, kanji: "菖蒲華", en: "Irises bloom along the water", startMonth: 6, startDay: 26 },
    { id: 30, kanji: "半夏生", en: "Crow-dipper sprouts in the shade", startMonth: 7, startDay: 1 },
    { id: 31, kanji: "温風至", en: "Warm winds begin to blow", startMonth: 7, startDay: 7 },
    { id: 32, kanji: "蓮始開", en: "Lotus flowers begin to bloom", startMonth: 7, startDay: 12 },
    { id: 33, kanji: "鷹乃学習", en: "Young hawks learn to fly", startMonth: 7, startDay: 17 },
    { id: 34, kanji: "桐始結花", en: "Paulownia trees begin to flower", startMonth: 7, startDay: 23 },
    { id: 35, kanji: "土潤溽暑", en: "Earth is damp and humid with heat", startMonth: 7, startDay: 28 },
    { id: 36, kanji: "大雨時行", en: "Heavy rains occasionally fall", startMonth: 8, startDay: 2 },
    { id: 37, kanji: "涼風至", en: "Cool winds begin to arrive", startMonth: 8, startDay: 7 },
    { id: 38, kanji: "寒蝉鳴", en: "Evening cicadas begin to sing", startMonth: 8, startDay: 12 },
    { id: 39, kanji: "蒙霧升降", en: "Thick mist descends upon the land", startMonth: 8, startDay: 17 },
    { id: 40, kanji: "綿柍開", en: "Cotton bolls begin to open", startMonth: 8, startDay: 23 },
    { id: 41, kanji: "天地始粛", en: "The heat finally begins to subside", startMonth: 8, startDay: 28 },
    { id: 42, kanji: "禾乃登", en: "Rice ripens in the fields", startMonth: 9, startDay: 2 },
    { id: 43, kanji: "草露白", en: "Dew glistens white on the grass", startMonth: 9, startDay: 8 },
    { id: 44, kanji: "鶺鴒鳴", en: "Wagtails chirp in the fields", startMonth: 9, startDay: 13 },
    { id: 45, kanji: "玄鳥去", en: "Swallows depart for warmer lands", startMonth: 9, startDay: 18 },
    { id: 46, kanji: "雷乃収声", en: "Thunder ceases its rumble", startMonth: 9, startDay: 23 },
    { id: 47, kanji: "蟄虫坏戸", en: "Insects hide and seal their doors", startMonth: 9, startDay: 28 },
    { id: 48, kanji: "水始涸", en: "Water begins to recede and dry", startMonth: 10, startDay: 3 },
    { id: 49, kanji: "鴻雁来", en: "Wild geese arrive from the north", startMonth: 10, startDay: 8 },
    { id: 50, kanji: "菊花開", en: "Chrysanthemums begin to bloom", startMonth: 10, startDay: 13 },
    { id: 51, kanji: "蟋蟀在戸", en: "Crickets chirp at the doorstep", startMonth: 10, startDay: 18 },
    { id: 52, kanji: "霜始降", en: "Frost begins to fall from the sky", startMonth: 10, startDay: 23 },
    { id: 53, kanji: "霎時施", en: "Light rain falls from time to time", startMonth: 10, startDay: 28 },
    { id: 54, kanji: "楓蔦黄", en: "Maple leaves and ivy turn yellow", startMonth: 11, startDay: 2 },
    { id: 55, kanji: "山茶始開", en: "Camellias begin to bloom", startMonth: 11, startDay: 7 },
    { id: 56, kanji: "地始凍", en: "The ground begins to freeze", startMonth: 11, startDay: 12 },
    { id: 57, kanji: "金盞香", en: "Daffodils begin to bloom", startMonth: 11, startDay: 17 },
    { id: 58, kanji: "虹蔵不見", en: "Rainbows hide and are no longer seen", startMonth: 11, startDay: 22 },
    { id: 59, kanji: "朔風払葉", en: "North wind sweeps fallen leaves", startMonth: 11, startDay: 27 },
    { id: 60, kanji: "橘始黄", en: "Citrus trees begin to turn yellow", startMonth: 12, startDay: 2 },
    { id: 61, kanji: "閉塞成冬", en: "Cold sets in and winter truly begins", startMonth: 12, startDay: 7 },
    { id: 62, kanji: "熊蟄穴", en: "Bears retreat to their dens to hibernate", startMonth: 12, startDay: 12 },
    { id: 63, kanji: "鱖魚群", en: "Salmon gather in schools upstream", startMonth: 12, startDay: 17 },
    { id: 64, kanji: "乃東生", en: "Prunella begins to sprout anew", startMonth: 12, startDay: 22 },
    { id: 65, kanji: "麋角解", en: "Elk shed their antlers", startMonth: 12, startDay: 27 },
    { id: 66, kanji: "雪下出麦", en: "Snow falls, yet wheat continues to grow", startMonth: 1, startDay: 1 },
    { id: 67, kanji: "芹乃栄", en: "Parsley flourishes in the cold", startMonth: 1, startDay: 6 },
    { id: 68, kanji: "水泉動", en: "Spring water begins to stir beneath the ice", startMonth: 1, startDay: 11 },
    { id: 69, kanji: "雉始鳴", en: "Pheasants begin to call", startMonth: 1, startDay: 16 },
    { id: 70, kanji: "款冬華", en: "Butterbur flowers bloom through the frost", startMonth: 1, startDay: 20 },
    { id: 71, kanji: "水沢腹堅", en: "Ice thickens on streams and rivers", startMonth: 1, startDay: 25 },
    { id: 72, kanji: "鶏始乳", en: "Hens begin to lay eggs", startMonth: 1, startDay: 30 },
];

const daysInMonth = [31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31];

function dayOfYear(m: number, d: number): number {
    let doy = 0;
    for (let i = 0; i < m - 1; i++) doy += daysInMonth[i];
    return doy + d;
}

function findKo() {
    const now = new Date();
    const today = dayOfYear(now.getMonth() + 1, now.getDate());
    for (let i = 0; i < koData.length; i++) {
        const start = dayOfYear(koData[i].startMonth, koData[i].startDay);
        const end = koData[i + 1]
            ? dayOfYear(koData[i + 1].startMonth, koData[i + 1].startDay) - 1
            : 365;
        if (today >= start && today <= end) return koData[i];
    }
    return koData[koData.length - 1];
}

const ko = findKo();
const kanjiEl = document.getElementById("ko-kanji");
const enEl = document.getElementById("ko-en");
if (kanjiEl) kanjiEl.textContent = ko.kanji;
if (enEl) enEl.textContent = ko.en;
