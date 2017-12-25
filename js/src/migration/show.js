// migrate to new version of showJs


if (window.clickyab_ad) {
  var slotDom = document.createElement("div");
  slotDom.setAttribute("class", "clickyab-ad");

  var keys = Object.keys(window.clickyab_ad);
  console.log(keys);
  for (var i = 0; i < keys.length; i++) {
    console.log(i, keys[i]);
    if (window.clickyab_ad[keys[i]]) {
      slotDom.setAttribute("clickyab-" + keys[i], window.clickyab_ad[keys[i]]);
    }
  }
  document.write(slotDom.outerHTML);

  if (!document.getElementById('clickyab-show-js-v2')){

    window["clickyabParams"] = {
      id: window.clickyab_ad["id"],
      domain: window.clickyab_ad["domain"]
    };

    setTimeout(function () {
      var scriptDom = document.createElement("script");
      scriptDom.setAttribute("id","clickyab-show-js-v2");
      scriptDom.setAttribute("src", "http://supplier.clickyab.ae/api/multi.js");
      document.body.appendChild(scriptDom);
    },0);
  }

}
