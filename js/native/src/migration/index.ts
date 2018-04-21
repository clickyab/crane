declare let clickyab_native: any;

let element = document.getElementsByClassName(clickyab_native.selector);
Object.keys(clickyab_native).forEach((option) => {
	if (option !== "selector") {
		if (element.item(0)) element.item(0).setAttribute(`data-${option}`, clickyab_native[option]);
	}
});

element.item(0).setAttribute("class", "clickyab-native");

if (!window["clickyab_native_migration"]) {
	window["clickyab_native_migration"] = true;
	let scriptTag = document.createElement("script");
	scriptTag.type = "text/javascript";
	scriptTag.src = "../build/browser/index.cjs.js";
	document.body.appendChild(scriptTag);
}
