import NativeComponent from "./native";


const elements = document.getElementsByClassName("clickyab-native");

for (let i = 0; i < elements.length; i++) {
	const wrapper: HTMLElement = elements.item(i) as HTMLElement;
	const url = "{{.URL}}?tid=__tid__&i=__id__&d=__domain__&count=__count__";

	const nativeComponent = new NativeComponent(wrapper, url);
}

