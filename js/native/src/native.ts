import {corners, fontFamilies, fontSizes, INative, INativeOptions, position, types} from "./definitions/native";

export default class NativeComponent {
	private wrapperElement: HTMLElement;
	private nativeUrl: string;
	private randomSting: string;
	private style: string = `__STYLE_TEMPLATE__`;
	private customOptions: INativeOptions = {};
	private defaultOptions: INativeOptions = {
		clickyab: "_clickyab_",
		type: types.grid4x,
		fontFamily: fontFamilies.samim,
		count: "4",
		corners: corners.sharp,
		title: "مطالب از سراسر وب",
		horizontal: false,
		fontsize: fontSizes.pt_8,
		position: position.top,
		orientation: "horizontal",
		tid: "ـ",
		domain: "ـ",
		id: "ـ",
		titleBackGround: "",
		titleColor: "#000",
		parent: document.location.href,
		ref: document.referrer,
		nostyle: "false",
	};
	private options: INativeOptions;

	constructor(wrapper: HTMLElement, url: string) {
		this.wrapperElement = wrapper;
		this.nativeUrl = url;
		this.randomSting = Math.random().toString().replace("0.", "");

		if (!this.domainValidation(wrapper)) {
			return;
		}

		this.fillOptions();
		window.addEventListener("resize", this.addWrapperStyle.bind(this));
		this.options = Object.assign(this.defaultOptions, this.customOptions) as INativeOptions;
		if (!this.customOptions.type) {
			this.customOptions.type = this.getType();
		}
		this.reload();
	}

	public reload() {
		this.loadNativeJson()
			.then((data: INative) => {
				this.wrapperElement.innerHTML = this.compileHtml(data.html);
				this.addClass(this.options.type || "");
				this.addClass(this.options.corners || "");
				this.addClass(this.options.horizontal ? "horizontal-item" : "");
				this.addClass(`clickyab_${this.randomSting}`);
				this.addClass(this.options.orientation || "");
				this.addWrapperStyle();
			});
	}

	private getType() {
		let type = types.grid4x;
		if (this.options.orientation === "vertical") {
			type = types.vertical;
		} else if (this.options.count === "1") {
			type = types.single;
		} else if (this.options.count === "3" || this.options.count === "6" || this.options.count === "9") {
			type = types.grid3x;
		}
		return type;
	}

	private addWrapperStyle() {

		this.removeClass("xs");
		this.removeClass("sm");
		this.removeClass("md");
		this.removeClass("lg");
		this.removeClass("xl");


		const width = this.wrapperElement.offsetWidth;
		if (width <= 280) {
			this.addClass("xs");
		} else if (280 < width && width <= 480) {
			this.addClass("sm");
		} else if (480 < width && width <= 970) {
			this.addClass("md");
		} else if (970 < width && width < 1200) {
			this.addClass("lg");
		} else {
			this.addClass("xl");
		}
	}


	private domainValidation(element: HTMLElement): boolean {
		if (location.hostname === "localhost" ||
			location.hostname === "127.0.0.1" ||
			location.hostname.split(".").splice(-2).join(".") === "clickyab.com" ||
			location.hostname.split(".").splice(-2).join(".") === "clickyab.ae"
		) {
			return true;
		}
		try {
			const domain = element.getAttribute("data-domain") as string;
			const baseDomain = domain.split(":")[0].split(".").splice(-2).join(".");
			const currentDomain = document.location.hostname.split(":")[0].split(".").splice(-2).join(".");
			if (baseDomain !== currentDomain) {
				console.error("Current domain is not match with config. It also happens when current page's domain is not valid.");
			}
			return baseDomain === currentDomain;
		} catch (e) {
			console.error("Current domain is not match with config. It also happens when current page's domain is not valid.");
			return false;
		}
	}

	private addClass(className: string) {
		let classes: string[] = (this.wrapperElement.getAttribute("class") || "").split(" ");
		if (classes.findIndex((c: string) => (c === className)) === -1) {
			classes.push(className);
			this.wrapperElement.setAttribute("class", classes.join(" "));
		}
	}

	private removeClass(className: string) {
		let classes: string[] = (this.wrapperElement.getAttribute("class") || "").split(" ");
		const indexOfClass = classes.findIndex((c: string) => (c === className));
		if (indexOfClass > -1) {
			classes.splice(indexOfClass, 1);
			this.wrapperElement.setAttribute("class", classes.join(" "));
		}
	}

	private fillOptions() {
		this.customOptions = {};
		Object.keys(this.defaultOptions).forEach(key => {
			const optionKey = key.replace(/([A-Z])/g, (g) => `-${g[0].toLowerCase()}`);
			if (this.wrapperElement.getAttribute(`data-${optionKey}`)) {
				this.customOptions[key] = this.wrapperElement.getAttribute(`data-${optionKey}`);
			}
		});
	}

	private loadNativeJson(): Promise<INative> {
		return new Promise((resolve, reject) => {
			const xhr = new XMLHttpRequest();
			const url = this.compiler(this.nativeUrl, Object.assign(this.defaultOptions, this.customOptions));
			xhr.onreadystatechange = function () {
				if (this.readyState === 4 && this.status === 200) {
					let nativeObj = JSON.parse(xhr.responseText);
					resolve(nativeObj as INative);
				}
			};

			xhr.onerror = () => {
				reject();
				throw new Error(`Failed to load Native Json from ${url}`);
			};

			xhr.open("GET", url);
			xhr.send();
		});
	}

	private compileHtml(htmlData: string): string {
		let html = htmlData.replace(new RegExp("_clickyab_", "ig"), `clickyab_${this.randomSting}`);

		let styleTag = this.style && this.customOptions.nostyle !== "true" ? `<style>
				${this.compiler(this.style, this.options)
			.replace(new RegExp("_clickyab_", "ig"), `clickyab_${this.randomSting}`)}
				</style>` : "";
		return ` ${this.getHeader()}${styleTag} <div class="cy-items">${html}</div>`;
	}

	private getHeader() {
		let templ = `
			<div class="cy-title-holder">
				<div class="cy-title ">${this.options.title}</div>
				<div class="cy-logo">
					<a rel="nofollow" target="_blank" href="https://www.clickyab.com/?ref=icon" class="cy-logo-container">
						<img src="data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiPz48c3ZnIHdpZHRoPSIxNXB4IiBoZWlnaHQ9IjE4cHgiIHZpZXdCb3g9IjAgMCAxNSAxOCIgdmVyc2lvbj0iMS4xIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHhtbG5zOnhsaW5rPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5L3hsaW5rIj4gICAgICAgIDx0aXRsZT5Hcm91cDwvdGl0bGU+ICAgIDxkZXNjPkNyZWF0ZWQgd2l0aCBTa2V0Y2guPC9kZXNjPiAgICA8ZGVmcz48L2RlZnM+ICAgIDxnIGlkPSJQYWdlLTIiIHN0cm9rZT0ibm9uZSIgc3Ryb2tlLXdpZHRoPSIxIiBmaWxsPSJub25lIiBmaWxsLXJ1bGU9ImV2ZW5vZGQiPiAgICAgICAgPGcgaWQ9IkRlc2t0b3AtSEQiIHRyYW5zZm9ybT0idHJhbnNsYXRlKC0xMTcuMDAwMDAwLCAtMjA2LjAwMDAwMCkiIGZpbGwtcnVsZT0ibm9uemVybyI+ICAgICAgICAgICAgPGcgaWQ9Ikdyb3VwLTIiIHRyYW5zZm9ybT0idHJhbnNsYXRlKDExNy4wMDAwMDAsIDIwMi4wMDAwMDApIj4gICAgICAgICAgICAgICAgPGcgaWQ9ImNsaWNreWFiLWVuIiB0cmFuc2Zvcm09InRyYW5zbGF0ZSgwLjAwMDAwMCwgMi4wMDAwMDApIj4gICAgICAgICAgICAgICAgICAgIDxnIGlkPSJHcm91cCIgdHJhbnNmb3JtPSJ0cmFuc2xhdGUoMC4wMDAwMDAsIDIuMDAwMDAwKSI+ICAgICAgICAgICAgICAgICAgICAgICAgPGcgaWQ9IlNoYXBlIiBmaWxsPSIjNDFCNkU2Ij4gICAgICAgICAgICAgICAgICAgICAgICAgICAgPHBhdGggaWQ9InVwTG9nbyIgZD0iTTEzLjM0NzE0NTUsMy41MTgxODE4MiBDMTEuOTk5Nzc2MSwxLjI3OTMzODg0IDkuNjE3Njg2NTcsMC4wMzcxOTAwODI2IDcuMTgzNDg4ODEsMC4wMzcxOTAwODI2IEw3LjE4MzQ4ODgxLDAuMDUyMDY2MTE1NyBMNy4xODM0ODg4MSwwLjA1MjA2NjExNTcgTDcuMTgzNDg4ODEsMC4wMzcxOTAwODI2IEw3LjE4MzQ4ODgxLDAuMDM3MTkwMDgyNiBDNS45MjU0NDc3NiwwLjAzNzE5MDA4MjYgNC42NDUwNzQ2MywwLjM3MTkwMDgyNiAzLjQ5MTI1LDEuMDYzNjM2MzYgQzAuMDg5MzI4MzU4MiwzLjEwOTA5MDkxIC0xLjAxMjM4ODA2LDcuNTEyMzk2NjkgMS4wMzQ3MjAxNSwxMC45MTE1NzAyIEMyLjM4MjA4OTU1LDEzLjE1MDQxMzIgNC43NjQxNzkxLDE0LjM5MjU2MiA3LjE5ODM3Njg3LDE0LjM5MjU2MiBDOC40NTY0MTc5MSwxNC4zOTI1NjIgOS43MzY3OTEwNCwxNC4wNTc4NTEyIDEwLjg5MDYxNTcsMTMuMzY2MTE1NyBDMTQuMjg1MDkzMywxMS4zMjgwOTkyIDE1LjM4NjgwOTcsNi45MTczNTUzNyAxMy4zNDcxNDU1LDMuNTE4MTgxODIgWiI+PC9wYXRoPiAgICAgICAgICAgICAgICAgICAgICAgIDwvZz4gICAgICAgICAgICAgICAgICAgICAgICA8cGF0aCBkPSJNMTAuNzQxNzM1MSw2LjEyODkyNTYyIEMxMC42ODk2MjY5LDYuMDQ3MTA3NDQgMTAuNjMwMDc0Niw1Ljk3MjcyNzI3IDEwLjU1NTYzNDMsNS45MDU3ODUxMiBDMTAuNTEwOTcwMSw1Ljg2ODU5NTA0IDEwLjQ1ODg2MTksNS44MzE0MDQ5NiAxMC4zOTkzMDk3LDUuODAxNjUyODkgTDEwLjM5MTg2NTcsNS44MTY1Mjg5MyBMMTAuMzkxODY1Nyw1LjgxNjUyODkzIEwxMC4zOTkzMDk3LDUuODAxNjUyODkgTDEwLjM4NDQyMTYsNS43OTQyMTQ4OCBMMTAuMzg0NDIxNiw1Ljc5NDIxNDg4IEw1Ljc2OTEyMzEzLDMuMzM5NjY5NDIgTDUuNzYxNjc5MSwzLjMzMjIzMTQgTDUuNzYxNjc5MSwzLjMzMjIzMTQgQzUuNzMxOTAyOTksMy4zMTczNTUzNyA1LjcwMjEyNjg3LDMuMzAyNDc5MzQgNS42NzIzNTA3NSwzLjI5NTA0MTMyIEM1LjYzNTEzMDYsMy4yODAxNjUyOSA1LjU5NzkxMDQ1LDMuMjY1Mjg5MjYgNS41NTMyNDYyNywzLjI1Nzg1MTI0IEM1LjMxNTAzNzMxLDMuMTk4MzQ3MTEgNS4wNjkzODQzMywzLjIzNTUzNzE5IDQuODYwOTUxNDksMy4zNjE5ODM0NyBDNC42NTI1MTg2NiwzLjQ4ODQyOTc1IDQuNTAzNjM4MDYsMy42ODkyNTYyIDQuNDQ0MDg1ODIsMy45MjcyNzI3MyBDNC40MzY2NDE3OSwzLjk3MTkwMDgzIDQuNDI5MTk3NzYsNC4wMDkwOTA5MSA0LjQyMTc1MzczLDQuMDUzNzE5MDEgQzQuNDIxNzUzNzMsNC4wODM0NzEwNyA0LjQxNDMwOTcsNC4xMTMyMjMxNCA0LjQxNDMwOTcsNC4xNDI5NzUyMSBMNC40MTQzMDk3LDQuMTQyOTc1MjEgTDQuNDE0MzA5Nyw5LjM3OTMzODg0IEw0LjQyMTc1MzczLDkuMzc5MzM4ODQgTDQuNDIxNzUzNzMsOS4zOTQyMTQ4OCBDNC40MjE3NTM3Myw5LjQ1MzcxOTAxIDQuNDI5MTk3NzYsOS41MjA2NjExNiA0LjQ0NDA4NTgyLDkuNTgwMTY1MjkgQzQuNDY2NDE3OTEsOS42NzY4NTk1IDQuNTAzNjM4MDYsOS43NjYxMTU3IDQuNTU1NzQ2MjcsOS44NTUzNzE5IEM0LjY4MjI5NDc4LDEwLjA2MzYzNjQgNC44ODMyODM1OCwxMC4yMTIzOTY3IDUuMTIxNDkyNTQsMTAuMjcxOTAwOCBDNS4zNTk3MDE0OSwxMC4zMzE0MDUgNS42MTI3OTg1MSwxMC4yOTQyMTQ5IDUuODIxMjMxMzQsMTAuMTY3NzY4NiBDNS45MjU0NDc3NiwxMC4xMDA4MjY0IDYuMDIyMjIwMTUsMTAuMDE5MDA4MyA2LjA4OTIxNjQyLDkuOTIyMzE0MDUgQzYuMTExNTQ4NTEsOS44ODUxMjM5NyA2LjE0MTMyNDYzLDkuODQ3OTMzODggNi4xNTYyMTI2OSw5LjgxMDc0MzggTDYuMTYzNjU2NzIsOS43OTU4Njc3NyBMNi4xNjM2NTY3Miw5Ljc5NTg2Nzc3IEM2LjE3ODU0NDc4LDkuNzczNTUzNzIgNi4xODU5ODg4MSw5Ljc0MzgwMTY1IDYuMjAwODc2ODcsOS43MTQwNDk1OSBMNi4yMDA4NzY4Nyw5LjcxNDA0OTU5IEw2Ljc1OTE3OTEsOC4zMTU3MDI0OCBMOC4yMzMwOTcwMSwxMC43NzAyNDc5IEw4LjI1NTQyOTEsMTAuOCBMOC4yNTU0MjkxLDEwLjggQzguNDM0MDg1ODIsMTEuMDY3NzY4NiA4LjczMTg0NzAxLDExLjIxNjUyODkgOS4wMjk2MDgyMSwxMS4yMTY1Mjg5IEM5LjE5MzM3Njg3LDExLjIxNjUyODkgOS4zNTcxNDU1MiwxMS4xNzE5MDA4IDkuNDk4NTgyMDksMTEuMDkwMDgyNiBDOS45MjI4OTE3OSwxMC44MzcxOTAxIDEwLjA2NDMyODQsMTAuMjg2Nzc2OSA5LjgzMzU2MzQzLDkuODU1MzcxOSBMOS44MzM1NjM0Myw5Ljg1NTM3MTkgTDguMzM3MzEzNDMsNy4zNzEwNzQzOCBMOS44NDEwMDc0Niw3LjUzNDcxMDc0IEw5Ljg0MTAwNzQ2LDcuNTM0NzEwNzQgQzkuODcwNzgzNTgsNy41NDIxNDg3NiA5LjkwODAwMzczLDcuNTQyMTQ4NzYgOS45MzAzMzU4Miw3LjU0MjE0ODc2IEw5Ljk0NTIyMzg4LDcuNTQyMTQ4NzYgQzkuOTg5ODg4MDYsNy41NDIxNDg3NiAxMC4wMjcxMDgyLDcuNTQyMTQ4NzYgMTAuMDcxNzcyNCw3LjUzNDcxMDc0IEMxMC4xOTgzMjA5LDcuNTE5ODM0NzEgMTAuMzA5OTgxMyw3LjQ3NTIwNjYxIDEwLjQyMTY0MTgsNy40MDgyNjQ0NiBDMTAuNjMwMDc0Niw3LjI4MTgxODE4IDEwLjc3ODk1NTIsNy4wODA5OTE3NCAxMC44Mzg1MDc1LDYuODQyOTc1MjEgQzEwLjkwNTUwMzcsNi41ODI2NDQ2MyAxMC44NjgyODM2LDYuMzM3MTkwMDggMTAuNzQxNzM1MSw2LjEyODkyNTYyIFoiIGlkPSJTaGFwZSIgZmlsbD0iI0ZGRkZGRiI+PC9wYXRoPiAgICAgICAgICAgICAgICAgICAgICAgIDxwYXRoIGlkPSJkb3duTG9nbyIgZD0iTTEzLjk0MjY2NzksMTYuMzA0MTMyMiBMMTIuODMzNTA3NSwxNC40NTk1MDQxIEMxMi42ODQ2MjY5LDE0LjIwNjYxMTYgMTIuNDM4OTczOSwxNC4wMjgwOTkyIDEyLjE1NjEwMDcsMTMuOTUzNzE5IEMxMS44NjU3ODM2LDEzLjg3OTMzODggMTEuNTc1NDY2NCwxMy45MjM5NjY5IDExLjMyMjM2OTQsMTQuMDgwMTY1MyBDMTAuODAxMjg3MywxNC4zOTI1NjIgMTAuNjMwMDc0NiwxNS4wNzY4NTk1IDEwLjk0MjcyMzksMTUuNTk3NTIwNyBMMTIuMDUxODg0MywxNy40NDIxNDg4IEMxMi4yNTI4NzMxLDE3Ljc2OTQyMTUgMTIuNjE3NjMwNiwxNy45Nzc2ODYgMTMuMDA0NzIwMSwxNy45Nzc2ODYgTDEzLjAwNDcyMDEsMTcuOTc3Njg2IEMxMy4yMDU3MDksMTcuOTc3Njg2IDEzLjM5OTI1MzcsMTcuOTI1NjE5OCAxMy41NzA0NjY0LDE3LjgyMTQ4NzYgQzEzLjgyMzU2MzQsMTcuNjcyNzI3MyAxNC4wMDIyMjAxLDE3LjQyNzI3MjcgMTQuMDc2NjYwNCwxNy4xNDQ2MjgxIEMxNC4xMzYyMTI3LDE2Ljg1NDU0NTUgMTQuMDkxNTQ4NSwxNi41NTcwMjQ4IDEzLjk0MjY2NzksMTYuMzA0MTMyMiBaIiBmaWxsPSIjNDFCNkU2Ij48L3BhdGg+ICAgICAgICAgICAgICAgICAgICA8L2c+ICAgICAgICAgICAgICAgIDwvZz4gICAgICAgICAgICA8L2c+ICAgICAgICA8L2c+ICAgIDwvZz48L3N2Zz4=">
					</a>
				</div>
			</div>
		`;


		return templ;
	}

	private compiler(text, options): string {
		let re = /__(.+?)__/g,
			reExp = /(^( )?(var|if|for|else|switch|case|break|{|}|;))(.*)?/g,
			code = "with(obj) { var r=[];\n",
			cursor = 0,
			result,
			match;
		let add = function (line, js = false) {
			js ? (code += line.match(reExp) ? `${line}\n` : `r.push(${line} ? ${line} : "");\n`) :
				(code += line !== "" ? `r.push('${line.replace(/"/g, "\"")}');\n` : "");
			return add;
		};
		while (match = re.exec(text)) {
			add(text.slice(cursor, match.index))(match[1], true);
			cursor = match.index + match[0].length;
		}
		add(text.substr(cursor, text.length - cursor));
		code = (code + `return r.join("");}`).replace(/[\r\t\n]/g, " ");
		try {
			result = new Function("obj", code).apply(options, [options]);
		}
		catch (err) {
			console.error("'" + err.message + "'", " in \n\nCode:\n", code, "\n");
		}
		return result;
	}

}
