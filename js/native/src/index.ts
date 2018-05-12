import NativeComponent from "./native";

migration();

const elements = document.getElementsByClassName("clickyab-native");

for (let i = 0; i < elements.length; i++) {
	setTimeout(() => {
		const wrapper: HTMLElement = elements.item(i) as HTMLElement;
		if (!wrapper.getAttribute("cy-key")) {
			wrapper.setAttribute("cy-key", Math.random().toString(36).substring(7));
			// const url = "//supplier.clickyab.com/api/get/native?tid=__tid__&i=__id__&d=__domain__&count=__count__&t=__type__&ref=__ref__&parent=__parent__";
			const url = "{{.URL}}?tid=__tid__&i=__id__&d=__domain__&count=__count__&t=__type__&ref=__ref__&parent=__parent__";

			const nativeComponent = new NativeComponent(wrapper, url);
		}
	}, i * 100);
}


function migration() {
	if (!window["clickyab_native_migration"]) {
		window["clickyab_native_migration"] = {};
	}
	const params = window["clickyab_native"];
	const elements = document.getElementsByClassName(params.selector);
	if (elements.item(0) && !window["clickyab_native_migration"][params.selector]) {
		const element = elements.item(0);
		Object.keys(params).forEach(key => {
			element.setAttribute(`data-${key}`, params[key]);
			element.setAttribute("class", element.getAttribute("class") + " clickyab-native");
		});
	}
	window["clickyab_native_migration"][params.selector] = true;
}
