const CONFIG = {
  // REMOTE_TARGET: "a.clickyab.com/ads/",
  REMOTE_TARGET: "{{.URL}}",
  MAXIMUM_AD_IN_PAGE: 30,
  SELECTOR_CLASS: "clickyab-ad",
  ELEMENT_PROPERTY_PREFIX: "clickyab-",
  BANNER_TYPES: {
    'video': 'video',
    'native': 'native',
    'mobile': 'mobile'
  },
  BANNER_SIZES: {
    '120_600': 1,
    '160_600': 2,
    '300_250': 3,
    '336_280': 4,
    '468_60': 5,
    '728_90': 6,
    '120_240': 7,
    '320_50': 8,
    '800_440': 9,
    '300_600': 11,
    '970_90': 12,
    '970_250': 13,
    '250_250': 14,
    '300_1050': 15,
    '320_480': 16,
    '480_320': 17,
    '128_128': 18,
  }
};

export default CONFIG;
