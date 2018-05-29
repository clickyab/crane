export interface INative {
	template: string;
	html: string;
}

export interface INativeOptions {
	clickyab?: string;
	type?: string;
	fontFamily?: fontFamilies;
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
	orientation?: "vertical" | "horizontal";
	parent?: string;
	ref?: string;
	nostyle?: string;
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
	pt_8 = "0.8rem",
	pt_10 = "1.0rem",
	pt_12 = "1.2rem",
	pt_16 = "1.6rem",
}


export enum position {
	top = "top",
	bottom = "bottom"
}

export enum types {
	grid = "grid",
	grid4x = "grid",
	grid3x = "grid3x",
	single = "single",
	vertical = "vertical",
}