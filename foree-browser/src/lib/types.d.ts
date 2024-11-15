export type MobileNavigations = {
	subMenuTitle: string,
	defaultActive?: boolean,
	navigations: NavigationLink[]
}[]

export type Navigations = NavigationLink[]

export interface NavigationLink {
	title: string;
	href: string;
	icon: string;
}

export interface Document {
	slug: string;
	file: string;
	metadata: {
		title: string;
		[key: string]: any;
	};
	breadcrumbs: Array<{ title: string }>;
	body: string;
	sections: Section[];
	children: Document[];
	assets?: Record<string, string>;
	next: null | { slug: string; title: string };
	prev: null | { slug: string; title: string };
}

export interface Section {
	slug: string;
	title: string;
}

export interface BannerData {
	id: string;
	start: Date;
	end: Date;
	arrow: boolean;
	href: string;
	content: {
		lg?: string;
		sm?: string;
	};
}
