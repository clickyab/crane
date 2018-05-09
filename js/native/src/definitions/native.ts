export interface INative {
	template: string;
	html: string;
}

export interface INativeOptions {
	clickyab ?: string;
	type ?: "grid";
	fontFamily ?: fontFamilies;
	count?: string;
	corners?: corners;
	title?: string;
	horizontal?: boolean;
	fontsize?: fontSizes;
	position?: position;
	domain?: string;
	id?: string;
	tid?: string;
	titleColor?: string;
	titleBackGround?: string;
	orientation ?: "vertical" | "horizontal";
}

export enum fontFamilies {
	sahel = "sahel",
	samim = "samim",
	vazir = "vazir",
	behdad = "behdad",
	nazanin = "nazanin"
}

export enum corners {
	sharp = "sharp",
	carve = "carve",
}

export enum fontSizes {
	pt_12 = "12pt",
	pt_14 = "14pt",
	pt_16 = "16pt",
	pt_18 = "18pt",
}


export enum position {
	top = "top",
	bottom = "bottom"
}