// migrate to new version of showJs


if (window.clickyab_ad) {

  var slotDom = document.createElement("div");
  slotDom.setAttribute("class", "clickyab-ad");

  var keys = Object.keys(window.clickyab_ad);
  for (var i = 0; i < keys.length; i++) {
    if (window.clickyab_ad[keys[i]]) {
      slotDom.setAttribute("clickyab-" + keys[i], window.clickyab_ad[keys[i]]);
    }
  }
  document.write(slotDom.outerHTML);

  window.addEventListener('DOMContentLoaded', function () {
    window.totalOfClickyabShowAd = window.totalOfClickyabShowAd || window.calculateAdsCount();
    window.countOfClickyabShowAd = window.countOfClickyabShowAd ? window.countOfClickyabShowAd + 1 : 1;
    var time = 1500;
    if (window.totalOfClickyabShowAd === window.countOfClickyabShowAd) {
      time = 0;
    }

    window["clickyabParams"] = {
      id: window.clickyab_ad["id"],
      domain: window.clickyab_ad["domain"]
    };


    setTimeout(function () {
      // if (!window.injectClickyabMultiJs) {
      window.injectClickyabMultiJs = true;
      var scriptDom = document.createElement("script");
      scriptDom.setAttribute("id", "clickyab-show-js-v2");
      scriptDom.setAttribute("src", "http://supplier.clickyab.ae/api/multi.js");
      document.body.appendChild(scriptDom);
      // }
    }, time);


  }, false);

//
  window.calculateAdsCount = function () {
    var element = document.getElementsByClassName("clickyab-ad");
    return element.length;
  }
}
